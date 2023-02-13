package depreciacionHelper

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// calcularCierre Calcula la transacción contable que se generará una vez se liquide el cierre a una fecha determinada
func calcularCierre(fechaCorte string, cuentas map[string]models.CuentaContable, transaccion *models.TransaccionMovimientos, resultado *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("calcularCierre - Unhandled Error!", "500")

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

	subgrupos := make(map[int]models.DetalleSubgrupo)
	for _, val := range infoCorte {
		if val.DeltaValor == 0 {
			continue
		}

		payload = "Id:" + strconv.Itoa(val.ElementoActaId)
		if elemento, err := actaRecibido.GetAllElemento(payload, "Id,ValorUnitario,ValorTotal,SubgrupoCatalogoId,TipoBienId", "", "", "", ""); err != nil {
			return err
		} else if len(elemento) == 1 {

			payload = "limit=1&fields=TipoBienId,Amortizacion,Depreciacion&sortby=Id&order=desc&query=Activo:true,SubgrupoId__Id:"
			if _, ok := subgrupos[elemento[0].SubgrupoCatalogoId]; !ok {
				if sg, err := catalogoElementos.GetAllDetalleSubgrupo(payload + strconv.Itoa(elemento[0].SubgrupoCatalogoId)); err != nil {
					return err
				} else if len(sg) == 1 {
					subgrupos[elemento[0].SubgrupoCatalogoId] = *sg[0]
				} else {
					resultado.Error = "No se pudo consultar la parametrización de las clases. Contacte soporte."
					return
				}
			}

			elemento[0].ValorTotal = val.DeltaValor
			if subgrupos[elemento[0].SubgrupoCatalogoId].Depreciacion || subgrupos[elemento[0].SubgrupoCatalogoId].Amortizacion {
				elementos_ = append(elementos_, elemento...)
			}

		} else {
			resultado.Error = "No se pudo consultar el detalle de los elementos. Contacte soporte."
			return
		}

	}

	if len(elementos_) == 0 {
		return
	}

	resultado.Error, outputError = asientoContable.CalcularMovimientosContables(elementos_, getDescripcionMovmientoCierre(), 0, formtatoCrr, terceroUD, terceroUD, cuentas, subgrupos, &transaccion.Movimientos)

	return
}

func getTipoComprobanteCierre() string {
	return "H22"
}

func getDescripcionMovmientoCierre() string {
	return "Mediciones posteriores"
}
