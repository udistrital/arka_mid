package bajasHelper

import (
	"net/url"
	"time"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/catalogoElementosHelper"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	crudTerceros "github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/depreciacionHelper"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// AprobarBajas Aprobación masiva de bajas: transacciones contables, actualización de movmientos y registro de novedades
func AprobarBajas(data *models.TrRevisionBaja) (ids []int, outputError map[string]interface{}) {

	funcion := "AprobarBajas"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		bajas               map[int]models.Movimiento
		elementosMovimiento map[int]models.ElementosMovimiento
		novedades           map[int]models.NovedadElemento
		elementosActa       map[int]models.Elemento
		cuentasBaja         map[int]models.CuentaSubgrupo
		cuentasDp           map[int]models.CuentaSubgrupo
		cuentasAm           map[int]models.CuentaSubgrupo
		detalleCuentas      map[string]models.CuentaContable
		detalleSubgrupos    map[int]models.DetalleSubgrupo
		detalleMediciones   map[int]models.FormatoDepreciacion
		detalleBajas        map[int]models.FormatoBaja
		movBj, movDp, movAm int
		parDebito           int
		parCredito          int
		comprobanteID       string
		query               string
		terceroUD           int
	)

	if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&movBj, "BJ_HT"); err != nil {
		return nil, err
	}
	if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&movDp, "DEP"); err != nil {
		return nil, err
	}
	if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&movAm, "AMT"); err != nil {
		return nil, err
	}

	if db_, cr_, err := parametros.GetParametrosDebitoCredito(); err != nil {
		return nil, err
	} else {
		parDebito = db_
		parCredito = cr_
	}

	if err := cuentasContables.GetComprobante(getTipoComprobanteBajas(), &comprobanteID); err != nil {
		return nil, err
	}

	query = "query=TipoDocumentoId__Nombre:NIT,Numero:" + crudTerceros.GetDocUD()
	if terceroUD_, err := crudTerceros.GetAllDatosIdentificacion(query); err != nil {
		return nil, err
	} else {
		terceroUD = terceroUD_[0].TerceroId.Id
	}

	detalleCuentas = make(map[string]models.CuentaContable)
	cuentasBaja, cuentasDp, cuentasAm = make(map[int]models.CuentaSubgrupo), make(map[int]models.CuentaSubgrupo), make(map[int]models.CuentaSubgrupo)

	// Paso 1: Consulta los movimientos
	query = "fields=Detalle,Id,FechaCreacion&limit=-1&query=Id__in:"
	query += url.QueryEscape(utilsHelper.ArrayToString(data.Bajas, "|"))
	if bajas_, err := movimientosArka.GetAllMovimiento(query); err != nil {
		return nil, err
	} else {

		bajas = make(map[int]models.Movimiento)
		detalleBajas = make(map[int]models.FormatoBaja)
		for _, baja := range bajas_ {

			var detalle models.FormatoBaja
			if err := utilsHelper.Unmarshal(baja.Detalle, &detalle); err != nil {
				return nil, err
			}

			bajas[baja.Id] = *baja
			detalleBajas[baja.Id] = detalle

			ids = append(ids, detalle.Elementos...)
		}
	}

	// Paso 2: Consulta los elementos
	query = "limit=-1&fields=Id,ElementoActaId,ValorTotal,ValorResidual,VidaUtil,MovimientoId&sortby=ElementoActaId&order=desc"
	query += "&query=Id__in:" + url.QueryEscape(utilsHelper.ArrayToString(ids, "|"))
	if elementos_, err := movimientosArka.GetAllElementosMovimiento(query); err != nil {
		return nil, err
	} else {
		elementosMovimiento = make(map[int]models.ElementosMovimiento)
		for _, el := range elementos_ {
			elementosMovimiento[el.Id] = *el
		}
	}

	// Paso 3: Consulta las novedades
	query = "limit=-1&sortby=MovimientoId,FechaCreacion&order=asc,asc&query=Activo:true,ElementoMovimientoId__Id__in:"
	query += utilsHelper.ArrayToString(ids, "|")
	if novedades_, err := movimientosArka.GetAllNovedadElemento(query); err != nil {
		return nil, err
	} else {
		novedades = make(map[int]models.NovedadElemento)
		for _, nov := range novedades_ {
			novedades[nov.ElementoMovimientoId.Id] = *nov
		}

	}

	// Paso 4: Consulta los elementos del acta
	ids = []int{}
	for _, el := range elementosMovimiento {
		ids = append(ids, el.ElementoActaId)
	}

	fields := "Id,SubgrupoCatalogoId"
	query = "Id__in:" + utilsHelper.ArrayToString(ids, "|")
	if elementos_, err := actaRecibido.GetAllElemento(query, fields, "", "", "", "-1"); err != nil {
		return nil, err
	} else {
		elementosActa = make(map[int]models.Elemento)
		for _, el_ := range elementos_ {
			elementosActa[el_.Id] = *el_
		}

	}

	ids = []int{}
	for _, el := range elementosActa {
		ids = append(ids, el.SubgrupoCatalogoId)
	}

	detalleSubgrupos = make(map[int]models.DetalleSubgrupo)
	if err := catalogoElementosHelper.GetDetalleSubgrupos(ids, detalleSubgrupos); err != nil {
		return nil, err
	}

	// Paso 5: Calcula el valor de la transaccion contable para cada baja (si aplica)
	detalleMediciones = make(map[int]models.FormatoDepreciacion)
	for _, baja := range bajas {

		movimientos := make([]*models.MovimientoTransaccion, 0)
		for _, el := range detalleBajas[baja.Id].Elementos {
			var (
				sg            int
				detalleSg     models.DetalleSubgrupo
				ref           time.Time
				valorPresente float64
				valorResidual float64
				vidaUtil      float64
				depreciacion  float64
				amortizacion  float64
				gasto         float64
			)

			if val, ok := elementosActa[elementosMovimiento[el].ElementoActaId]; ok {
				sg = val.SubgrupoCatalogoId
			}

			if val, ok := detalleSubgrupos[sg]; ok {
				detalleSg = val
			}

			if nov, ok := novedades[el]; ok && nov.ValorLibros > 0 {

				if _, ok := detalleMediciones[nov.MovimientoId.Id]; !ok {
					if dt_, err := depreciacionHelper.GetDetalleDepreciacion(nov.MovimientoId.Detalle); err != nil {
						return nil, err
					} else {
						detalleMediciones[nov.MovimientoId.Id] = *dt_
					}
				}

				ref, _ = time.Parse("2006-01-02", detalleMediciones[nov.MovimientoId.Id].FechaCorte)
				valorPresente = nov.ValorLibros
				valorResidual = nov.ValorResidual
				vidaUtil = nov.VidaUtil

			} else if !ok && elementosMovimiento[el].ValorTotal > 0 {

				ref = elementosMovimiento[el].MovimientoId.FechaModificacion
				valorPresente = elementosMovimiento[el].ValorTotal
				valorResidual = elementosMovimiento[el].ValorResidual
				vidaUtil = elementosMovimiento[el].VidaUtil

			} else {
				continue
			}

			if detalleSg.Amortizacion || detalleSg.Depreciacion {
				valorMedicion, _ := depreciacionHelper.CalculaDp(
					valorPresente,
					valorResidual,
					vidaUtil,
					ref.AddDate(0, 0, 1),
					baja.FechaCreacion)

				if valorMedicion > 0 {
					if detalleSg.Depreciacion {
						depreciacion = valorMedicion
						asientoContable.GetInfoContableSubgrupos(movDp, []int{sg}, cuentasDp, detalleCuentas)
					} else if detalleSg.Amortizacion {
						amortizacion = valorMedicion
						asientoContable.GetInfoContableSubgrupos(movAm, []int{sg}, cuentasAm, detalleCuentas)
					}
				}

				if valorPresente-valorMedicion > 0 {
					gasto = valorPresente - valorMedicion
					asientoContable.GetInfoContableSubgrupos(movBj, []int{sg}, cuentasBaja, detalleCuentas)
				}

			} else {
				gasto = valorPresente
				asientoContable.GetInfoContableSubgrupos(movBj, []int{sg}, cuentasBaja, detalleCuentas)
			}

			movimientosContablesBaja(cuentasBaja, cuentasDp, cuentasAm, gasto, depreciacion, amortizacion, sg, parCredito, parDebito, terceroUD, detalleCuentas, &movimientos)

		}

		if len(movimientos) > 0 {
			transaccion := *new(models.TransaccionMovimientos)

			if comprobanteID != "" {
				etiquetas := *new(models.Etiquetas)
				etiquetas.ComprobanteId = comprobanteID
				if err := utilsHelper.Marshal(etiquetas, &transaccion.Etiquetas); err != nil {
					return nil, err
				}
			} else {
				transaccion.Etiquetas = ""
			}

			transaccion.Activo = true
			transaccion.ConsecutivoId = detalleBajas[baja.Id].ConsecutivoId
			transaccion.Descripcion = "Baja de elementos"
			transaccion.FechaTransaccion = time.Now()
			transaccion.Movimientos = movimientos

			if _, err := movimientosContables.PostTrContable(&transaccion); err != nil {
				return nil, err
			}
		}

		data_ := data
		data_.Bajas = []int{baja.Id}
		if ids_, err := movimientosArka.PutRevision(data_); err != nil {
			return nil, err
		} else {
			ids = ids_
		}

	}

	return ids, nil
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
func movimientosContablesBaja(cuentasBj, cuentasDp, cuentasAm map[int]models.CuentaSubgrupo,
	gasto, depreciacion, amortizacion float64, subgrupo, credito, debito, terceroUD int,
	detalleCuentas map[string]models.CuentaContable, movimientos *[]*models.MovimientoTransaccion) {

	var medicion float64

	if depreciacion > 0 {
		medicion = depreciacion
		ctaDp := detalleCuentas[cuentasDp[subgrupo].CuentaDebitoId]
		movDp := asientoContable.CreaMovimiento(depreciacion, descMovDp(), terceroUD, &ctaDp, debito)
		*movimientos = append(*movimientos, movDp)
	} else if amortizacion > 0 {
		medicion = amortizacion
		ctaAm := detalleCuentas[cuentasAm[subgrupo].CuentaDebitoId]
		movAm := asientoContable.CreaMovimiento(amortizacion, descMovAm(), terceroUD, &ctaAm, debito)
		*movimientos = append(*movimientos, movAm)
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

func descMovDp() string {
	return "Depreciación restante en baja de elementos"
}

func descMovAm() string {
	return "Amortización restante en baja de elementos"
}

func getTipoComprobanteBajas() string {
	return "H23"
}
