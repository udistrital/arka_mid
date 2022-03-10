package ajustesHelper

import (
	"fmt"
	"strconv"

	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/entradaHelper"
	"github.com/udistrital/arka_mid/helpers/movimientosArkaHelper"
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
	if elementos_, err := actaRecibido.GetAllElemento(query, "", "Id", "desc", "0", "-1"); err != nil {
		return nil, err
	} else {
		orgiginalesActa = elementos_
	}

	if entrada_, err := movimientosArkaHelper.GetEntradaByActa(orgiginalesActa[0].ActaRecibidoId.Id); err != nil {
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
		if ha, err := actaRecibido.GetAllHistoricoActa(query, "", "FechaCreacion", "desc", "", "-1"); err != nil {
			return nil, err
		} else {
			proveedorId = ha[0].ProveedorId
		}

		if cs, err := entradaHelper.GetConsecutivoEntrada(entrada.Detalle); err != nil {
			return nil, err
		} else {
			consecutivo = cs
		}

		if movsEntrada, err := calcularAjusteMovimiento(orgiginalesActa, updateVls, updateSg, entrada.FormatoTipoMovimientoId.Id, consecutivo, proveedorId, "Entrada"); err != nil {
			return nil, err
		} else {
			movimientos = append(movimientos, movsEntrada...)
		}
	}

	if entrada.EstadoMovimientoId.Nombre == "Entrada Con Salida" {

		query = "limit=-1&sortby=MovimientoId,ElementoActaId&order=desc,desc&query=ElementoActaId__in:" + utilsHelper.ArrayToString(ids, "|")
		if elementos_, err := movimientosArkaHelper.GetAllElementosMovimiento(query); err != nil {
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
				if fm, err := movimientosArkaHelper.GetAllFormatoTipoMovimiento(query); err != nil {
					return nil, err
				} else {
					tipoMovimientoSalida = fm[0].Id
				}
			}
			ids = []int{}
		}
	}

	for _, elms := range elementosSalida {

		var funcionario int
		var consecutivo string

		if func_, cons_, err := salidaHelper.GetInfoSalida(elms.Salida.Detalle); err != nil {
			return nil, err
		} else {
			funcionario = func_
			consecutivo = cons_
		}

		for _, el := range elms.UpdateSg {
			ids = append(ids, el.Id)
		}

		for _, el := range elms.UpdateVls {
			ids = append(ids, el.Id)
		}

		if movsSalida, err := calcularAjusteMovimiento(orgiginalesActa, elms.UpdateVls, elms.UpdateSg, tipoMovimientoSalida, consecutivo, funcionario, "Salida"); err != nil {
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
		if novedades_, err := movimientosArkaHelper.GetAllNovedadElemento(query); err != nil {
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
		resultado.TrContable = tr
	}

	if elementos_, err := fillElementos(append(updateSg, (append(updateVls, append(updateMsc, updateMp...)...))...)); err != nil {
		return nil, err
	} else {
		resultado.Elementos = elementos_
	}

	return resultado, nil

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
	if elsMov, outputError = movimientosArkaHelper.GetAllElementosMovimiento(query); outputError != nil {
		return nil, outputError
	}

	for _, el := range elsMov {
		if idx := findElementoInArrayEM(elsActa, el.ElementoActaId); idx > -1 {
			var elemento_ *models.DetalleElemento__
			if elemento_, outputError = fillElemento(elsActa[idx], el); outputError != nil {
				return nil, outputError
			}

			elementos = append(elementos, elemento_)
		}
	}

	return elementos, nil

}
