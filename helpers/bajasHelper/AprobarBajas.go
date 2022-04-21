package bajasHelper

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/catalogoElementosHelper"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	crudMovimientosArka "github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
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
		bajas                             map[int]models.Movimiento
		elementosMovimiento               map[int]models.ElementosMovimiento
		novedades                         map[int]models.NovedadElemento
		elementosActa                     map[int]models.Elemento
		cuentasBaja                       map[int]models.CuentaSubgrupo
		cuentasDp                         map[int]models.CuentaSubgrupo
		cuentasAm                         map[int]models.CuentaSubgrupo
		detalleCuentas                    map[string]models.CuentaContable
		detalleSubgrupos                  map[int]models.DetalleSubgrupo
		detalleMediciones                 map[int]models.FormatoDepreciacion
		detalleBajas                      map[int]models.FormatoBaja
		movBj, movDp, movAm               int
		parDebito                         int
		parCredito                        int
		totalesDp, totalesAm, totalesBaja map[int]float64
		comprobanteID                     string
	)

	if err := crudMovimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&movBj, "BJ_HT"); err != nil {
		return nil, err
	}
	if err := crudMovimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&movDp, "DEP"); err != nil {
		return nil, err
	}
	if err := crudMovimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&movAm, "AMT"); err != nil {
		return nil, err
	}

	if db_, cr_, err := parametros.GetParametrosDebitoCredito(); err != nil {
		return nil, err
	} else {
		parDebito = db_
		parCredito = cr_
	}

	if err := cuentasContables.GetComprobante(tipoComprobanteBaja(), &comprobanteID); err != nil {
		return nil, err
	}

	detalleCuentas = make(map[string]models.CuentaContable)
	cuentasBaja, cuentasDp, cuentasAm = make(map[int]models.CuentaSubgrupo), make(map[int]models.CuentaSubgrupo), make(map[int]models.CuentaSubgrupo)

	// Paso 1: Consulta los movimientos
	query := "fields=Detalle,Id,FechaCreacion&limit=-1&query=Id__in:"
	query += url.QueryEscape(utilsHelper.ArrayToString(data.Bajas, "|"))
	if bajas_, err := crudMovimientosArka.GetAllMovimiento(query); err != nil {
		return nil, err
	} else {

		bajas = make(map[int]models.Movimiento)
		detalleBajas = make(map[int]models.FormatoBaja)
		for _, baja := range bajas_ {

			var detalle models.FormatoBaja
			if err := json.Unmarshal([]byte(baja.Detalle), &detalle); err != nil {
				logs.Error(err)
				eval := " - json.Unmarshal([]byte(mov.Detalle), &detalle)"
				return nil, errorctrl.Error(funcion+eval, err, "500")
			}

			bajas[baja.Id] = *baja
			detalleBajas[baja.Id] = detalle

			ids = append(ids, detalle.Elementos...)
		}
	}

	// Paso 2: Consulta los elementos
	query = "limit=-1&fields=Id,ElementoActaId,ValorTotal,ValorResidual,VidaUtil,MovimientoId&sortby=ElementoActaId&order=desc"
	query += "&query=Id__in:" + url.QueryEscape(utilsHelper.ArrayToString(ids, "|"))
	if elementos_, err := crudMovimientosArka.GetAllElementosMovimiento(query); err != nil {
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
	if novedades_, err := crudMovimientosArka.GetAllNovedadElemento(query); err != nil {
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
		totalesBaja, totalesAm, totalesDp = make(map[int]float64), make(map[int]float64), make(map[int]float64)
		movimientos := make([]*models.MovimientoTransaccion, 0)
		for _, el := range detalleBajas[baja.Id].Elementos {
			var sg int
			var ref time.Time
			var valorPresente, valorResidual, vidaUtil float64
			var detalleSg models.DetalleSubgrupo

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
						utilsHelper.FillMapTotales(totalesDp, sg, valorMedicion)
					} else if detalleSg.Amortizacion {
						utilsHelper.FillMapTotales(totalesAm, sg, valorMedicion)
					}
				}

				utilsHelper.FillMapTotales(totalesBaja, sg, valorPresente-valorMedicion)

			} else {
				utilsHelper.FillMapTotales(totalesBaja, sg, valorPresente)
			}

		}

		if len(totalesBaja) > 0 {
			ids = []int{}
			for sg := range totalesBaja {
				ids = append(ids, sg)
			}

			asientoContable.GetInfoContableSubgrupos(movBj, ids, cuentasBaja, detalleCuentas)
			asientoContable.GenerarMovimientosContables(totalesBaja, detalleCuentas, cuentasBaja, parDebito, parCredito, 10000, descMovBaja(), false, &movimientos)
		}

		if len(totalesDp) > 0 {
			ids = []int{}
			for sg := range totalesDp {
				ids = append(ids, sg)
			}

			asientoContable.GetInfoContableSubgrupos(movDp, ids, cuentasDp, detalleCuentas)
			asientoContable.GenerarMovimientosContables(totalesDp, detalleCuentas, cuentasDp, parDebito, parCredito, 10000, descMovDp(), false, &movimientos)
		}

		if len(totalesAm) > 0 {
			ids = []int{}
			for sg := range totalesAm {
				ids = append(ids, sg)
			}

			asientoContable.GetInfoContableSubgrupos(movAm, ids, cuentasAm, detalleCuentas)
			asientoContable.GenerarMovimientosContables(totalesAm, detalleCuentas, cuentasAm, parDebito, parCredito, 10000, descMovAm(), false, &movimientos)
		}

		if len(movimientos) > 0 {
			transaccion := *new(models.TransaccionMovimientos)

			if comprobanteID != "" {
				etiquetas := *new(models.Etiquetas)
				etiquetas.ComprobanteId = comprobanteID
				if jsonData, err := json.Marshal(etiquetas); err != nil {
					logs.Error(err)
					eval := " - json.Marshal(etiquetas)"
					return nil, errorctrl.Error(funcion+eval, err, "500")
				} else {
					transaccion.Etiquetas = string(jsonData[:])
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

func descMovBaja() string {
	return "Baja de elementos"
}

func descMovDp() string {
	return "Depreciación restante en baja de elementos"
}

func descMovAm() string {
	return "Amortización restante en baja de elementos"
}

func tipoComprobanteBaja() string {
	return "H21"
}
