package bajasHelper

import (
	"net/url"
	"strconv"
	"time"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/depreciacionHelper"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// AprobarBajas Aprobación masiva de bajas: transacciones contables, actualización de movmientos y registro de novedades
func AprobarBajas(data *models.TrRevisionBaja, response *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("AprobarBajas - Unhandled Error!", "500")

	var (
		movBj, movCr int
		terceroUD    int
		bajas        []*models.Movimiento
	)

	if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&movBj, "BJ_HT"); err != nil {
		return err
	}

	if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&movCr, "CRR"); err != nil {
		return err
	}

	if UD, err := terceros.GetTerceroUD(); err != nil {
		return err
	} else if UD == 0 {
		response.Error = "No se pudo consultar el tercero para asociar a la transacción contable. Contacte soporte."
		return
	} else {
		terceroUD = UD
	}

	var (
		bufferCuentas     = make(map[string]models.CuentaContable)
		detalleSubgrupos  = make(map[int]models.DetalleSubgrupo)
		detalleMediciones = make(map[int]models.FormatoDepreciacion)
		detalleBajas      = make(map[int]models.FormatoBaja)
	)

	if bajas_, err := movimientosArka.GetAllMovimiento(payloadBajas(data.Bajas)); err != nil {
		return err
	} else if len(bajas_) == len(data.Bajas) {
		bajas = bajas_
	} else {
		response.Error = "No se pudo consultar el detalle de las bajas a aprobar. Contacte soporte."
		return
	}

	for _, baja := range bajas {
		var detalle models.FormatoBaja
		if err := utilsHelper.Unmarshal(baja.Detalle, &detalle); err != nil {
			return err
		}

		detalleBajas[baja.Id] = detalle
	}

	for _, baja := range bajas {

		var (
			ids                 []int
			detalleBaja         models.FormatoBaja
			elementosMovimiento []*models.ElementosMovimiento
			elementosActa       = make(map[int]models.Elemento)
			bajas               []*models.Elemento
			mediciones          []*models.Elemento
			transaccion         models.TransaccionMovimientos
		)

		detalleBaja = detalleBajas[baja.Id]
		if elementos_, err := movimientosArka.GetAllElementosMovimiento(payloadElementosMovimiento(detalleBaja.Elementos)); err != nil {
			return err
		} else if len(elementos_) == len(detalleBaja.Elementos) {
			elementosMovimiento = elementos_
			for _, el := range elementos_ {
				ids = append(ids, el.ElementoActaId)
			}
		} else {
			response.Error = "No se pudo consultar el detalle de los elementos. Contacte soporte."
			return
		}

		if elementos_, err := actaRecibido.GetAllElemento(payloadElementos(ids), fieldsElementos(), "", "", "", "-1"); err != nil {
			return err
		} else if len(elementos_) == len(detalleBaja.Elementos) {
			for _, el := range elementos_ {
				elementosActa[el.Id] = *el
			}
		} else {
			response.Error = "No se pudo consultar el detalle de los elementos. Contacte soporte."
			return
		}

		for _, el := range elementosMovimiento {

			acta := elementosActa[el.ElementoActaId]
			var novedad *models.NovedadElemento

			if _, ok := detalleSubgrupos[acta.SubgrupoCatalogoId]; !ok {
				if detalle, err := catalogoElementos.GetAllDetalleSubgrupo(getPayloadDetalleSubgrupo(acta.SubgrupoCatalogoId)); err != nil {
					return err
				} else if len(detalle) == 1 {
					detalleSubgrupos[acta.SubgrupoCatalogoId] = *detalle[0]
				} else {
					response.Error = "No se pudo la parametrización de los elementos. Contacte soporte."
					return
				}
			}

			subgrupo := detalleSubgrupos[acta.SubgrupoCatalogoId]
			if novedad_, err := movimientosArka.GetAllNovedadElemento(payloadNovedad(el.Id)); err != nil {
				return err
			} else if len(novedad_) == 1 {
				novedad = novedad_[0]
				if _, ok := detalleMediciones[novedad.MovimientoId.Id]; !ok {
					var detalle models.FormatoDepreciacion
					if err := utilsHelper.Unmarshal(novedad.MovimientoId.Detalle, &detalle); err != nil {
						return err
					}
					detalleMediciones[novedad.MovimientoId.Id] = detalle
				}
			}

			if subgrupo.Amortizacion || subgrupo.Depreciacion {

				var (
					ref           time.Time
					valorPresente float64
					valorResidual float64
					vidaUtil      float64
				)

				if novedad != nil {
					if novedad.ValorLibros-novedad.ValorResidual > 0 {
						ref, _ = time.Parse("2006-01-02", detalleMediciones[novedad.MovimientoId.Id].FechaCorte)
						valorPresente = novedad.ValorLibros
						valorResidual = novedad.ValorResidual
						vidaUtil = novedad.VidaUtil
						acta.ValorTotal = novedad.ValorLibros
					} else if novedad.ValorLibros > 0 { // Ya fue depreciado totalmente pero queda el valor residual
						baja_ := acta
						baja_.ValorTotal = novedad.ValorLibros
						bajas = append(bajas, &acta)
						continue
					} else { // Ya fue depreciado totalmente. No hay afectaciones contables.
						continue
					}
				} else { // No hay novedad pero se debe depreciar
					ref = el.MovimientoId.FechaModificacion
					valorPresente = el.ValorTotal
					valorResidual = el.ValorResidual
					vidaUtil = el.VidaUtil
				}

				valorMedicion, _ := depreciacionHelper.CalculaDp(
					valorPresente,
					valorResidual,
					vidaUtil,
					ref.AddDate(0, 0, 1),
					baja.FechaCreacion.Local())

				if valorMedicion > 0 {
					medicion_ := acta
					medicion_.ValorTotal = valorMedicion
					acta.ValorTotal -= valorMedicion

					if subgrupo.Depreciacion || subgrupo.Amortizacion {
						mediciones = append(mediciones, &medicion_)
					}
				}
				bajas = append(bajas, &acta)
			} else {
				bajas = append(bajas, &acta)
			}

		}

		if msg, err := asientoContable.CalcularMovimientosContables(bajas, descBaja(), 0, movBj, terceroUD, terceroUD, bufferCuentas, detalleSubgrupos, &transaccion.Movimientos); err != nil || msg != "" {
			response.Error = msg
			return err
		}

		if msg, err := asientoContable.CalcularMovimientosContables(mediciones, descMovCr(), 0, movCr, terceroUD, terceroUD, bufferCuentas, detalleSubgrupos, &transaccion.Movimientos); err != nil || msg != "" {
			response.Error = msg
			return err
		}

		if len(transaccion.Movimientos) > 0 {
			if msg, err := asientoContable.CreateTransaccionContable(getTipoComprobanteBajas(), "Baja de elementos almacén.", &transaccion); err != nil || msg != "" {
				response.Error = msg
				return err
			}

			transaccion.ConsecutivoId = detalleBaja.ConsecutivoId
			if _, err := movimientosContables.PostTrContable(&transaccion); err != nil {
				return err
			}
		}

		data_ := data
		data_.Bajas = []int{baja.Id}
		if ids_, err := movimientosArka.PutRevision(data_); err != nil {
			return err
		} else {
			ids = ids_
		}

	}

	return
}

func GetTerceroIdEncargado(elementoId int, terceroId *int) (outputError map[string]interface{}) {

	funcion := "GetTerceroIdEncargado"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var historial models.Historial
	if historial_, err := movimientosArka.GetHistorialElemento(elementoId, true); err != nil {
		return err
	} else {
		historial = *historial_
	}

	if tercero, _, err := GetEncargado(&historial); err != nil {
		return err
	} else {
		*terceroId = tercero
	}

	return
}

// movimientosContablesBaja Genera los tres movimientos contables para un elemenento dado de baja.
func movimientosContablesBaja(cuentasBj, cuentasDp, cuentasAm map[int]models.CuentasSubgrupo,
	gasto, medicion float64, subgrupo, credito, debito, terceroUD int,
	detalleCuentas map[string]models.CuentaContable, movimientos *[]*models.MovimientoTransaccion) {

	if medicion > 0 {
		ctaDp := detalleCuentas[cuentasDp[subgrupo].CuentaDebitoId]
		movDp := asientoContable.CreaMovimiento(medicion, descMovCr(), terceroUD, &ctaDp, debito)
		*movimientos = append(*movimientos, movDp)
	}

	if gasto > 0 {
		ctaGasto := detalleCuentas[cuentasBj[subgrupo].CuentaDebitoId]
		movGasto := asientoContable.CreaMovimiento(gasto, descMovGasto(), terceroUD, &ctaGasto, debito)
		*movimientos = append(*movimientos, movGasto)
	}

	if gasto+medicion > 0 {
		ctaInventario := detalleCuentas[cuentasBj[subgrupo].CuentaCreditoId]
		movInventario := asientoContable.CreaMovimiento(gasto+medicion, descMovInventario(), terceroUD, &ctaInventario, credito)
		*movimientos = append(*movimientos, movInventario)
	}

	return

}

func descMovInventario() string {
	return "Movimiento a cuenta de inventario"
}

func descMovGasto() string {
	return "Movimiento a cuenta de gasto"
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
	return "fields=Detalle,Id,FechaCreacion&limit=-1&query=Id__in:" + url.QueryEscape(utilsHelper.ArrayToString(ids, "|"))
}

func payloadElementosMovimiento(elementos []int) string {
	return "limit=-1&fields=Id,ElementoActaId,ValorTotal,ValorResidual,VidaUtil,MovimientoId&sortby=ElementoActaId&order=desc&query=Id__in:" +
		url.QueryEscape(utilsHelper.ArrayToString(elementos, "|"))
}

func payloadElementos(elementos []int) string {
	return "Id__in:" + utilsHelper.ArrayToString(elementos, "|")
}

func fieldsElementos() string {
	return "Id,SubgrupoCatalogoId,TipoBienId,ValorUnitario,ValorTotal"
}

func getPayloadDetalleSubgrupo(id int) string {
	return "limit=1&fields=SubgrupoId,TipoBienId,Amortizacion,Depreciacion&sortby=FechaCreacion&order=desc&query=Activo:true,SubgrupoId__Id:" +
		strconv.Itoa(id)
}

func payloadNovedad(id int) string {
	return "limit=1&sortby=FechaCreacion&order=desc&query=Activo:true,ElementoMovimientoId__Id:" + strconv.Itoa(id)
}
