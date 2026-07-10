package depreciacionHelper

import (
	"strconv"
	"sync"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

const (
	maxWorkersConsultaElementosCierre = 8
	loteElementosCierreSize           = 200
)

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

type resultadoLoteElementosCierre struct {
	elementos   []*models.Elemento
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
	idsElementoActa := make([]int, 0, len(infoCorte))
	for _, val := range infoCorte {
		if val.DeltaValor != 0 {
			pendientes = append(pendientes, val)
			if val.ElementoActaId > 0 {
				idsElementoActa = append(idsElementoActa, val.ElementoActaId)
			}
		}
	}

	if len(pendientes) == 0 {
		return nil, "", nil
	}
	if len(idsElementoActa) == 0 {
		return nil, "No se pudo consultar el detalle de los elementos. Contacte soporte.", nil
	}

	elementosPorID, outputError := consultarElementosCierrePorLotes(idsElementoActa)
	if outputError != nil {
		return nil, "", outputError
	}

	detalles = make([]elementoCierreDetalle, 0, len(pendientes))
	for _, val := range pendientes {
		elemento, ok := elementosPorID[val.ElementoActaId]
		if !ok {
			errMsg = "No se pudo consultar el detalle de los elementos. Contacte soporte."
			continue
		}
		if !elementoValidoParaDepreciacion(elemento) {
			continue
		}

		detalles = append(detalles, elementoCierreDetalle{
			deltaValor: val.DeltaValor,
			elemento:   elemento,
		})
	}

	return detalles, errMsg, nil
}

func consultarElementosCierrePorLotes(ids []int) (map[int]*models.Elemento, map[string]interface{}) {
	ids = utilsHelper.RemoveDuplicateInt(ids)
	if len(ids) == 0 {
		return map[int]*models.Elemento{}, nil
	}

	lotes := chunkInts(ids, loteElementosCierreSize)
	workers := maxWorkersConsultaElementosCierre
	if len(lotes) < workers {
		workers = len(lotes)
	}

	jobs := make(chan []int, len(lotes))
	results := make(chan resultadoLoteElementosCierre, len(lotes))

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for lote := range jobs {
				payload := "Id__in:" + utilsHelper.ArrayToString(lote, "|")
				elementos, err := getAllElementoDepreciacionCierre(payload, "Id,ValorUnitario,ValorTotal,SubgrupoCatalogoId,TipoBienId,Activo", "", "", "", "-1")
				results <- resultadoLoteElementosCierre{
					elementos:   elementos,
					outputError: err,
				}
			}
		}()
	}

	for _, lote := range lotes {
		jobs <- lote
	}
	close(jobs)

	wg.Wait()
	close(results)

	elementosPorID := make(map[int]*models.Elemento, len(ids))
	for result := range results {
		if result.outputError != nil {
			return nil, result.outputError
		}
		for _, elemento := range result.elementos {
			if elemento == nil || elemento.Id <= 0 {
				continue
			}
			elementosPorID[elemento.Id] = elemento
		}
	}

	return elementosPorID, nil
}

func chunkInts(ids []int, size int) [][]int {
	if size <= 0 || len(ids) == 0 {
		return nil
	}

	lotes := make([][]int, 0, (len(ids)+size-1)/size)
	for start := 0; start < len(ids); start += size {
		end := start + size
		if end > len(ids) {
			end = len(ids)
		}

		lote := make([]int, end-start)
		copy(lote, ids[start:end])
		lotes = append(lotes, lote)
	}

	return lotes
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
