package ajustesHelper

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/asientoContable"
	crudActas "github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/salidaHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// GenerarAjusteAutomatico Genera transacción contable, actualiza elementos y novedades como consecuencia de actualizar una serie de elementos de un acta
func GenerarAjusteAutomatico(elementos []*models.DetalleElemento_) (resultado *models.DetalleAjusteAutomatico, outputError map[string]interface{}) {

	funcion := "GenerarAjusteAutomatico"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		ids                  []int
		query                string
		tipoMovimientoSalida int
		entrada              *models.Movimiento
		orgiginalesActa      []*models.Elemento
		updateSg             []*models.DetalleElemento_
		updateVls            []*models.DetalleElemento_
		updateMp             []*models.DetalleElemento_
		updateMsc            []*models.DetalleElemento_
		elementosSalida      map[int]*models.ElementosPorActualizarSalida
		movimientos          []*models.MovimientoTransaccion
		nuevosNovedades      []*models.NovedadElemento
		nuevosMovArka        []*models.ElementosMovimiento
		nuevosActa           []*models.Elemento
	)

	resultado = new(models.DetalleAjusteAutomatico)

	for _, el := range elementos {
		ids = append(ids, el.Id)
	}

	query = "Id__in:" + utilsHelper.ArrayToString(ids, "|")
	if elementos_, err := crudActas.GetAllElemento(query, "", "Id", "desc", "0", "-1"); err != nil {
		return nil, err
	} else {
		orgiginalesActa = elementos_
	}

	if entrada_, err := movimientosArka.GetEntradaByActa(orgiginalesActa[0].ActaRecibidoId.Id); err != nil {
		return nil, err
	} else if entrada_ == nil {
		return nil, nil
	} else {
		entrada = entrada_
	}

	if msc, vls, sg, mp, err := separarElementosPorModificacion(orgiginalesActa, elementos, entrada.EstadoMovimientoId.Nombre == "Entrada Con Salida"); err != nil {
		return nil, err
	} else {
		updateMsc = msc
		updateVls = vls
		updateSg = sg
		updateMp = mp
	}

	if (len(updateVls)+len(updateSg) > 0) && (entrada.EstadoMovimientoId.Nombre == "Entrada Aprobada" || entrada.EstadoMovimientoId.Nombre == "Entrada Con Salida") {
		var proveedorId int
		var consecutivo string

		query = "Activo:true,ActaRecibidoId__Id:" + strconv.Itoa(orgiginalesActa[0].ActaRecibidoId.Id)
		if ha, err := crudActas.GetAllHistoricoActa(query, "", "FechaCreacion", "desc", "", "-1"); err != nil {
			return nil, err
		} else {
			proveedorId = ha[0].ProveedorId
		}

		if entrada.Consecutivo != nil {
			consecutivo = *entrada.Consecutivo
		}

		if movsEntrada, err := calcularAjusteMovimiento(orgiginalesActa, updateVls, updateSg, entrada.FormatoTipoMovimientoId.Id, proveedorId, consecutivo, "Entrada"); err != nil {
			return nil, err
		} else {
			movimientos = append(movimientos, movsEntrada...)
		}
	}

	if entrada.EstadoMovimientoId.Nombre == "Entrada Con Salida" {

		query = "limit=-1&sortby=MovimientoId,ElementoActaId&order=desc,desc&query=ElementoActaId__in:" + utilsHelper.ArrayToString(ids, "|")
		if elementos_, err := movimientosArka.GetAllElementosMovimiento(query); err != nil {
			return nil, err
		} else {
			if elementosSalida_, updateMp_, actualizados_, err := separarElementosPorSalida(elementos_, updateVls, updateSg, updateMp); err != nil {
				return nil, err
			} else {
				nuevosMovArka = actualizados_
				elementosSalida = elementosSalida_
				updateMp = updateMp_
			}

			if len(elementosSalida) > 0 {
				query = "query=CodigoAbreviacion:SAL"
				if fm, err := movimientosArka.GetAllFormatoTipoMovimiento(query); err != nil {
					return nil, err
				} else {
					tipoMovimientoSalida = fm[0].Id
				}
			}
			ids = []int{}
		}
	}

	for _, elms := range elementosSalida {

		funcionario, err := salidaHelper.GetInfoSalida(elms.Salida.Detalle)
		if err != nil {
			return nil, err
		}

		for _, el := range elms.UpdateSg {
			ids = append(ids, el.Id)
		}

		for _, el := range elms.UpdateVls {
			ids = append(ids, el.Id)
		}

		if movsSalida, err := calcularAjusteMovimiento(orgiginalesActa, elms.UpdateVls, elms.UpdateSg, tipoMovimientoSalida, funcionario, *elms.Salida.Consecutivo, "Salida"); err != nil {
			return nil, err
		} else {
			movimientos = append(movimientos, movsSalida...)
		}

	}

	for _, el := range updateMp {
		ids = append(ids, el.Id)
	}

	if len(ids) > 0 {
		query = "limit=-1&sortby=MovimientoId,FechaCreacion&order=asc,asc&query=ElementoMovimientoId__ElementoActaId__in:" + utilsHelper.ArrayToString(ids, "|")
		if novedades_, err := movimientosArka.GetAllNovedadElemento(query); err != nil {
			return nil, err
		} else {
			novedadesMedicion := separarNovedadesPorElemento(novedades_)

			if movimientos_, novedades_, err := calcularAjusteMediciones(novedadesMedicion, updateSg, updateVls, updateMp, orgiginalesActa); err != nil {
				return nil, err
			} else {
				nuevosNovedades = novedades_
				movimientos = append(movimientos, movimientos_...)
			}
		}
	}

	if len(updateSg)+len(updateVls)+len(updateMsc) > 0 {
		if nuevos, err := generarNuevosActa(append(updateSg, (append(updateVls, updateMsc...))...)); err != nil {
			return nil, err
		} else {
			nuevosActa = nuevos
		}
	} else if len(updateMp) == 0 {
		return resultado, nil
	}

	if err := submitUpdates(nuevosActa, nuevosMovArka, nuevosNovedades); err != nil {
		return nil, err
	}

	if rs, tr, err := generarMovimientoAjuste(updateSg, updateVls, updateMsc, updateMp, movimientos); err != nil {
		return nil, err
	} else {
		resultado.Movimiento = rs
		if tr != nil && tr.Movimientos != nil && len(tr.Movimientos) > 0 {
			if tr_, err := asientoContable.GetDetalleContable(tr.Movimientos, nil); err != nil {
				return nil, err
			} else {
				resultado.TrContable = tr_
			}
		}
	}

	if elementos_, err := fillElementos(append(updateSg, (append(updateVls, append(updateMsc, updateMp...)...))...)); err != nil {
		return nil, err
	} else {
		resultado.Elementos = elementos_
	}

	return resultado, nil

}

// GetAjusteAutomatico Consulta el detalle de los elementos y la transacción contable asociada a un ajuste.
func GetAjusteAutomatico(movimientoId int) (ajuste *models.DetalleAjusteAutomatico, outputError map[string]interface{}) {

	funcion := "GetAjusteAutomatico"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		ids           []int
		query         string
		movimiento    *models.Movimiento
		detalle       *models.FormatoAjusteAutomatico
		elementosActa []*models.DetalleElemento
		elementosMov  []*models.ElementosMovimiento
		elementos     []*models.DetalleElemento__
	)

	ajuste = new(models.DetalleAjusteAutomatico)

	if movimiento, outputError = movimientosArka.GetMovimientoById(movimientoId); outputError != nil {
		return nil, outputError
	}

	if err := json.Unmarshal([]byte(movimiento.Detalle), &detalle); err != nil {
		logs.Error(err)
		eval := " - json.Unmarshal([]byte(movimiento.Detalle), &detalle)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	if elementosActa, outputError = actaRecibido.GetElementos(0, detalle.Elementos); outputError != nil {
		return nil, outputError
	}

	for _, e := range elementosActa {
		ids = append(ids, e.Id)
	}

	query = "limit=-1&sortby=Id&order=desc&query=ElementoActaId__in:" + utilsHelper.ArrayToString(ids, "|")
	if elementosMov, outputError = movimientosArka.GetAllElementosMovimiento(query); outputError != nil {
		return nil, outputError
	}

	for _, el := range elementosActa {
		var idx int
		var elemento_ *models.DetalleElemento__
		detalle := new(models.ElementosMovimiento)

		if idx = utilsHelper.FindElementoInArrayElementosMovimiento(elementosMov, el.Id); idx > -1 {
			detalle = elementosMov[idx]
		} else {
			detalle.ValorResidual = 0
			detalle.VidaUtil = 0
		}

		if elemento_, outputError = utilsHelper.FillElemento(el, detalle); outputError != nil {
			return nil, outputError
		}

		elementos = append(elementos, elemento_)

	}

	if movimiento.ConsecutivoId != nil && *movimiento.ConsecutivoId > 0 {
		if tr, err := movimientosContables.GetTransaccion(*movimiento.ConsecutivoId, "consecutivo", true); err != nil {
			return nil, err
		} else if len(tr.Movimientos) > 0 {
			if detalleContable, err := asientoContable.GetDetalleContable(tr.Movimientos, nil); err != nil {
				return nil, err
			} else {
				ajuste.TrContable = detalleContable
			}
		}
	}

	ajuste.Movimiento = movimiento
	ajuste.Elementos = elementos

	return ajuste, nil

}

// GetDetalleElementosActa Genera transacción contable, actualiza elementos y novedades como consecuencia de actualizar una serie de elementos de un acta
func GetDetalleElementosActa(actaRecibidoId int) (elementos []*models.DetalleElemento__, outputError map[string]interface{}) {

	funcion := "GetDetalleElementosActa"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		ids     []int
		query   string
		elsActa []*models.DetalleElemento
		elsMov  []*models.ElementosMovimiento
	)

	if elsActa, outputError = actaRecibido.GetElementos(actaRecibidoId, []int{}); outputError != nil {
		return nil, outputError
	} else if len(elsActa) == 0 {
		return nil, nil
	}

	for _, e := range elsActa {
		ids = append(ids, e.Id)
	}

	query = "limit=-1&sortby=Id&order=desc&query=ElementoActaId__in:" + utilsHelper.ArrayToString(ids, "|")
	if elsMov, outputError = movimientosArka.GetAllElementosMovimiento(query); outputError != nil {
		return nil, outputError
	}

	if len(elsMov) > 0 {
		for _, el := range elsMov {
			if idx := findElementoInArrayEM(elsActa, *el.ElementoActaId); idx > -1 {
				var elemento_ *models.DetalleElemento__
				if elemento_, outputError = utilsHelper.FillElemento(elsActa[idx], el); outputError != nil {
					return nil, outputError
				}

				elementos = append(elementos, elemento_)
			}
		}
	} else {
		if elementos, outputError = generarNuevos(elsActa); outputError != nil {
			return nil, outputError
		}
	}

	return elementos, nil

}
