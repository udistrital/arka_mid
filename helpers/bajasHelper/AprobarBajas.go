package bajasHelper

import (
	"net/url"
	"strconv"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/depreciacionHelper"
	"github.com/udistrital/arka_mid/helpers/inventarioHelper"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// AprobarBajas Aprobación masiva de bajas: transacciones contables, actualización de movmientos y registro de novedades
func AprobarBajas(data *models.TrRevisionBaja, response *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("AprobarBajas - Unhandled Error!", "500")

	var movBj, movCr int

	outputError = movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&movBj, "BJ_HT")
	if outputError != nil {
		return
	}

	outputError = movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&movCr, "CRR")
	if outputError != nil {
		return
	}

	terceroUD, outputError := terceros.GetTerceroUD()
	if outputError != nil {
		return
	} else if terceroUD == 0 {
		response.Error = "No se pudo consultar el tercero para asociar a la transacción contable. Contacte soporte."
		return
	}

	var (
		bufferCuentas    = make(map[string]models.CuentaContable)
		detalleSubgrupos = make(map[int]models.DetalleSubgrupo)
	)

	bajas, outputError := movimientosArka.GetAllMovimiento(payloadBajas(data.Bajas))
	if outputError != nil {
		return
	} else if len(bajas) != len(data.Bajas) {
		response.Error = "No se pudo consultar el detalle de las bajas a aprobar. Contacte soporte."
		return
	}

	for _, baja := range bajas {

		var (
			detalleBaja models.FormatoBaja
			bajas       []*models.Elemento
			mediciones  []*models.Elemento
			transaccion models.TransaccionMovimientos
		)

		outputError = utilsHelper.Unmarshal(baja.Detalle, &detalleBaja)
		if outputError != nil {
			return
		}

		for _, el := range detalleBaja.Elementos {

			historial, err := movimientosArka.GetHistorialElemento(el, true)
			if err != nil {
				return err
			} else if historial == nil {
				response.Error = "No se pudo la parametrización de los elementos. Contacte soporte."
				return
			}

			valor, residual, vidaUtil, ref, err := inventarioHelper.GetUltimoValor(*historial)
			if err != nil {
				return err
			}

			if valor-residual <= 0 && valor <= 0 {
				continue
			}

			var elementoActa models.Elemento
			outputError = actaRecibido.GetElementoById(*historial.Elemento.ElementoActaId, &elementoActa)
			if outputError != nil {
				return
			}

			if _, ok := detalleSubgrupos[elementoActa.SubgrupoCatalogoId]; !ok {
				if detalle, err := catalogoElementos.GetAllDetalleSubgrupo(getPayloadDetalleSubgrupo(elementoActa.SubgrupoCatalogoId)); err != nil {
					return err
				} else if len(detalle) == 1 {
					detalleSubgrupos[elementoActa.SubgrupoCatalogoId] = *detalle[0]
				} else {
					response.Error = "No se pudo la parametrización de los elementos. Contacte soporte."
					return
				}
			}

			elementoActa.ValorTotal = valor
			subgrupo := detalleSubgrupos[elementoActa.SubgrupoCatalogoId]
			if !subgrupo.Amortizacion && !subgrupo.Depreciacion || valor-residual == 0 {
				bajas = append(bajas, &elementoActa)
				continue
			}

			valorMedicion, _ := depreciacionHelper.CalculaDp(
				valor,
				residual,
				vidaUtil,
				ref.AddDate(0, 0, 1).UTC(),
				baja.FechaCreacion.UTC())

			if valorMedicion > 0 {
				medicion_ := elementoActa
				medicion_.ValorTotal = valorMedicion
				elementoActa.ValorTotal -= valorMedicion
				mediciones = append(mediciones, &medicion_)
			}

			if elementoActa.ValorTotal > 0 {
				bajas = append(bajas, &elementoActa)
			}
		}

		response.Error, outputError = asientoContable.CalcularMovimientosContables(bajas, descBaja(), 0, movBj, terceroUD, terceroUD, bufferCuentas, detalleSubgrupos, &transaccion.Movimientos)
		if outputError != nil || response.Error != "" {
			return
		}

		response.Error, outputError = asientoContable.CalcularMovimientosContables(mediciones, descMovCr(), 0, movCr, terceroUD, terceroUD, bufferCuentas, detalleSubgrupos, &transaccion.Movimientos)
		if outputError != nil || response.Error != "" {
			return
		}

		if len(transaccion.Movimientos) > 0 {
			response.Error, outputError = asientoContable.CreateTransaccionContable(getTipoComprobanteBajas(), "Baja de elementos almacén.", &transaccion)
			if outputError != nil || response.Error != "" {
				return
			}

			transaccion.ConsecutivoId = *baja.ConsecutivoId
			_, outputError = movimientosContables.PostTrContable(&transaccion)
			if outputError != nil {
				return
			}
		}

		data_ := data
		data_.Bajas = []int{baja.Id}
		_, outputError = movimientosArka.PutRevision(data_)
		if outputError != nil {
			return
		}

	}

	return
}

func descMovCr() string {
	return "Depreciación restante en baja de elementos"
}

func descBaja() string {
	return "Baja de elementos de almacén."
}

func getTipoComprobanteBajas() string {
	return "H23"
}

func payloadBajas(ids []int) string {
	return "fields=ConsecutivoId,Detalle,Id,FechaCreacion&limit=-1&query=Id__in:" + url.QueryEscape(utilsHelper.ArrayToString(ids, "|"))
}

func getPayloadDetalleSubgrupo(id int) string {
	return "limit=1&fields=SubgrupoId,TipoBienId,Amortizacion,Depreciacion&sortby=FechaCreacion&order=desc&query=Activo:true,SubgrupoId__Id:" +
		strconv.Itoa(id)
}
