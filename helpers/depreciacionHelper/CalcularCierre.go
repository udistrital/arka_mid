package depreciacionHelper

import (
	"strconv"
	"sync"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

const maxWorkersConsultaElementosCierre = 5

var getAllElementoDepreciacionCierre = actaRecibido.GetAllElemento

type elementoCierreDetalle struct {
	deltaValor float64
	elemento   *models.Elemento
}

type resultadoElementoCierre struct {
	detalle     *elementoCierreDetalle
	errorMsg    string
	outputError map[string]interface{}
}

// calcularCierre Calcula la transacción contable que se generará una vez se liquide el cierre a una fecha determinada
func calcularCierre(fechaCorte string, cuentas map[string]models.CuentaContable, transaccion *models.TransaccionMovimientos, resultado *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("calcularCierre - Unhandled Error!", "500")

	var (
		formtatoCrr int
		elementos_  []*models.Elemento
		payload     string
	)

	outputError = movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&formtatoCrr, "CRR")
	if outputError != nil {
		return
	}

	infoCorte, outputError := movimientosArka.GetCorteDepreciacion(fechaCorte)
	if outputError != nil {
		return
	}

	if len(infoCorte) == 0 {
		return
	}

	terceroUD, outputError := terceros.GetTerceroUD()
	if outputError != nil {
		return
	} else if terceroUD == 0 {
		resultado.Error = "No se pudo consultar el tercero para asociar a la transacción contable. Contacte soporte."
		return
	}

	detallesElementos, errMsg, outputError := consultarElementosParaCierre(infoCorte)
	if outputError != nil {
		return outputError
	}
	if errMsg != "" {
		resultado.Error = errMsg
		return
	}

	subgrupos := make(map[int]models.DetalleSubgrupo)
	for _, detalleElemento := range detallesElementos {
		elemento := detalleElemento.elemento
		if elemento == nil {
			continue
		}

		payload = "limit=1&fields=TipoBienId,Amortizacion,Depreciacion,SubgrupoId&sortby=Id&order=desc&query=Activo:true,SubgrupoId__Id:"
		if _, ok := subgrupos[elemento.SubgrupoCatalogoId]; !ok {
			if sg, err := catalogoElementos.GetAllDetalleSubgrupo(payload + strconv.Itoa(elemento.SubgrupoCatalogoId)); err != nil {
				return err
			} else if len(sg) == 1 {
				subgrupos[elemento.SubgrupoCatalogoId] = *sg[0]
			} else {
				resultado.Error = "No se pudo consultar la parametrización de las clases. Contacte soporte."
				return
			}
		}

		elemento.ValorTotal = detalleElemento.deltaValor
		if subgrupos[elemento.SubgrupoCatalogoId].Depreciacion || subgrupos[elemento.SubgrupoCatalogoId].Amortizacion {
			elementos_ = append(elementos_, elemento)
		}
	}

	if len(elementos_) == 0 {
		return
	}

	resultado.Error, outputError = asientoContable.CalcularMovimientosContables(elementos_, getDescripcionMovmientoCierre(), 0, formtatoCrr, terceroUD, terceroUD, cuentas, subgrupos, &transaccion.Movimientos)

	return
}

func consultarElementosParaCierre(infoCorte []models.DepreciacionElemento) (detalles []elementoCierreDetalle, errMsg string, outputError map[string]interface{}) {
	pendientes := make([]models.DepreciacionElemento, 0, len(infoCorte))
	for _, val := range infoCorte {
		if val.DeltaValor != 0 {
			pendientes = append(pendientes, val)
		}
	}

	if len(pendientes) == 0 {
		return nil, "", nil
	}

	workers := maxWorkersConsultaElementosCierre
	if len(pendientes) < workers {
		workers = len(pendientes)
	}

	jobs := make(chan models.DepreciacionElemento, len(pendientes))
	results := make(chan resultadoElementoCierre, len(pendientes))

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for val := range jobs {
				results <- consultarElementoCierre(val)
			}
		}()
	}

	for _, val := range pendientes {
		jobs <- val
	}
	close(jobs)

	wg.Wait()
	close(results)

	detalles = make([]elementoCierreDetalle, 0, len(pendientes))
	for result := range results {
		if result.outputError != nil && outputError == nil {
			outputError = result.outputError
		}
		if result.errorMsg != "" && errMsg == "" {
			errMsg = result.errorMsg
		}
		if result.detalle != nil {
			detalles = append(detalles, *result.detalle)
		}
	}

	return detalles, errMsg, outputError
}

func consultarElementoCierre(val models.DepreciacionElemento) resultadoElementoCierre {
	payload := "Id:" + strconv.Itoa(val.ElementoActaId)
	elementos, err := getAllElementoDepreciacionCierre(payload, "Id,ValorUnitario,ValorTotal,SubgrupoCatalogoId,TipoBienId,Activo", "", "", "", "")
	if err != nil {
		return resultadoElementoCierre{outputError: err}
	}
	if len(elementos) != 1 {
		return resultadoElementoCierre{errorMsg: "No se pudo consultar el detalle de los elementos. Contacte soporte."}
	}
	if !elementoValidoParaDepreciacion(elementos[0]) {
		return resultadoElementoCierre{}
	}

	return resultadoElementoCierre{
		detalle: &elementoCierreDetalle{
			deltaValor: val.DeltaValor,
			elemento:   elementos[0],
		},
	}
}

func elementoValidoParaDepreciacion(elemento *models.Elemento) bool {
	if elemento == nil {
		return false
	}

	if !elemento.Activo {
		return false
	}

	if elemento.TipoBienId <= 0 {
		return false
	}

	return true
}

func getTipoComprobanteCierre() string {
	return "H22"
}

func getDescripcionMovmientoCierre() string {
	return "Mediciones posteriores"
}
