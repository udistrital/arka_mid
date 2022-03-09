package ajustesHelper

import (
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/catalogoElementosHelper"
	"github.com/udistrital/arka_mid/helpers/entradaHelper"
	"github.com/udistrital/arka_mid/helpers/movimientosArkaHelper"
	"github.com/udistrital/arka_mid/helpers/parametrosHelper"
	"github.com/udistrital/arka_mid/helpers/salidaHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// GenerarAjusteAutomatico Genera transacción contable, actualiza elementos y novedades como consecuencia de actualizar una serie de elementos de un acta
func GenerarAjusteAutomatico(elementos []*models.DetalleElemento_) (resultado []*models.MovimientoTransaccion, outputError map[string]interface{}) {

	funcion := "GenerarAjusteAutomatico"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		query                string
		idsEl                []int
		entrada              *models.Movimiento
		orgActa              []*models.Elemento
		elementosSalida      map[int]*models.ElementosPorActualizarSalida
		updateVls            []*models.DetalleElemento_
		updateSg             []*models.DetalleElemento_
		updateMp             []*models.DetalleElemento_
		movimientos          []*models.MovimientoTransaccion
		tipoMovimientoSalida int
		novedades            []*models.NovedadElemento
		actualizados         []*models.ElementosMovimiento
	)

	for _, el := range elementos {
		idsEl = append(idsEl, el.Id)
	}

	query = "Id__in:" + utilsHelper.ArrayToString(idsEl, "|")
	if elementos_, err := actaRecibido.GetAllElemento(query, "", "Id", "desc", "0", "-1"); err != nil {
		return nil, err
	} else {
		orgActa = elementos_
	}

	if entrada_, err := movimientosArkaHelper.GetEntradaByActa(orgActa[0].ActaRecibidoId.Id); err != nil {
		return nil, err
	} else if entrada_ == nil {
		return nil, nil
	} else {
		entrada = entrada_
	}

	if _, vls, sg, mp, err := separarElementosPorModificacion(orgActa, elementos, entrada.EstadoMovimientoId.Nombre == "Entrada Con Salida"); err != nil {
		return nil, err
	} else {
		updateVls = vls
		updateSg = sg
		updateMp = mp
	}

	if (len(updateVls)+len(updateSg) > 0) && (entrada.EstadoMovimientoId.Nombre == "Entrada Aprobada" || entrada.EstadoMovimientoId.Nombre == "Entrada Con Salida") {
		var proveedorId int
		var consecutivo string

		query = "Activo:true,ActaRecibidoId__Id:" + strconv.Itoa(orgActa[0].ActaRecibidoId.Id)
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

		if movsEntrada, err := calcularAjusteMovimiento(orgActa, updateVls, updateSg, entrada.FormatoTipoMovimientoId.Id, consecutivo, proveedorId, "Entrada"); err != nil {
			return nil, err
		} else {
			movimientos = append(movimientos, movsEntrada...)
		}
	}

	if entrada.EstadoMovimientoId.Nombre == "Entrada Con Salida" {

		query = "limit=-1&sortby=MovimientoId,ElementoActaId&order=desc,desc&query=ElementoActaId__in:" + utilsHelper.ArrayToString(idsEl, "|")
		if elementos_, err := movimientosArkaHelper.GetAllElementosMovimiento(query); err != nil {
			return nil, err
		} else {
			if elementosSalida_, updateMp_, actualizados_, err := separarElementosPorSalida(elementos_, updateVls, updateSg, updateMp); err != nil {
				return nil, err
			} else {
				actualizados = actualizados_
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
			idsEl = []int{}
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
			idsEl = append(idsEl, el.Id)
		}

		for _, el := range elms.UpdateVls {
			idsEl = append(idsEl, el.Id)
		}

		if movsSalida, err := calcularAjusteMovimiento(orgActa, elms.UpdateVls, elms.UpdateSg, tipoMovimientoSalida, consecutivo, funcionario, "Salida"); err != nil {
			return nil, err
		} else {
			movimientos = append(movimientos, movsSalida...)
		}

	}

	if len(idsEl) > 0 {
		query = "limit=-1&sortby=MovimientoId,FechaCreacion&order=asc,asc&query=ElementoMovimientoId__ElementoActaId__in:" + utilsHelper.ArrayToString(idsEl, "|")
		if novedades_, err := movimientosArkaHelper.GetAllNovedadElemento(query); err != nil {
			return nil, err
		} else {
			novedadesMedicion := separarNovedadesPorElemento(novedades_)

			if movimientos_, novedades_, err := calcularAjusteMediciones(novedadesMedicion, updateSg, updateVls, updateMp, orgActa); err != nil {
				return nil, err
			} else {
				novedades = novedades_
				movimientos = append(movimientos, movimientos_...)
			}
		}
	}

	return movimientos, nil
}

// separarElementosPorModificacion Separa los elementos según se deba modificar Subgrupo, Valores, Misceláneos o Mediciones posteriores
func separarElementosPorModificacion(originales []*models.Elemento, actualizados []*models.DetalleElemento_, mediciones bool) (msc, vls, sg, mp []*models.DetalleElemento_, outputError map[string]interface{}) {

	funcion := "separarElementosPorModificacion"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	msc = make([]*models.DetalleElemento_, 0)
	vls = make([]*models.DetalleElemento_, 0)
	sg = make([]*models.DetalleElemento_, 0)
	mp = make([]*models.DetalleElemento_, 0)

	for _, el_ := range originales {
		if idx := findElementoInArrayD(actualizados, el_.Id); idx > -1 {
			if msc_, vls_, sg_, err := determinarDeltaActa(el_, actualizados[idx]); err != nil {
				return nil, nil, nil, nil, err
			} else if msc_ {
				msc = append(msc, actualizados[idx])
			} else if vls_ {
				vls = append(vls, actualizados[idx])
			} else if sg_ {
				sg = append(sg, actualizados[idx])
			} else if mediciones {
				mp = append(mp, actualizados[idx])
			}
		}
	}

	return msc, vls, sg, mp, nil

}

// calcularAjusteMovimiento Calcula la transacción contable generada a partir de los elementos y el cambio de cada uno
func calcularAjusteMovimiento(originales []*models.Elemento, actualizarVl, actualizarSg []*models.DetalleElemento_, movimientoId int, consecutivo string, proveedorId int, tipoMovimiento string) (movimientos []*models.MovimientoTransaccion, outputError map[string]interface{}) {

	funcion := "calcularAjusteMovimiento"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		ids             []int
		movDebito       int
		movCredito      int
		cuentasSubgrupo map[int]*models.CuentaSubgrupo
		detalleCuenta   map[string]*models.CuentaContable
	)

	cuentasSubgrupo = make(map[int]*models.CuentaSubgrupo)
	detalleCuenta = make(map[string]*models.CuentaContable)
	if db_, cr_, err := parametrosHelper.GetParametrosDebitoCredito(); err != nil {
		return nil, err
	} else {
		movDebito = db_
		movCredito = cr_
	}

	for _, el := range originales {
		ids = append(ids, el.SubgrupoCatalogoId)
	}

	for _, el := range actualizarSg {
		ids = append(ids, el.SubgrupoCatalogoId)
	}

	if cuentasSg, err := getCuentasByMovimientoSubgrupos(movimientoId, ids); err != nil {
		return nil, err
	} else {
		cuentasSubgrupo = cuentasSg
	}

	for _, el := range actualizarSg {
		if idx := findElementoInArrayE(originales, el.Id); idx > -1 {

			if detalleCuenta_, err := fillCuentas(detalleCuenta, []string{cuentasSubgrupo[originales[idx].SubgrupoCatalogoId].CuentaCreditoId,
				cuentasSubgrupo[el.SubgrupoCatalogoId].CuentaCreditoId, cuentasSubgrupo[originales[idx].SubgrupoCatalogoId].CuentaDebitoId, cuentasSubgrupo[el.SubgrupoCatalogoId].CuentaDebitoId}); err != nil {
				return nil, err
			} else {
				detalleCuenta = detalleCuenta_
			}

			movimientos = append(movimientos, generaTrContable(el.ValorTotal-originales[idx].ValorTotal, consecutivo, tipoMovimiento,
				movDebito, movCredito, originales[idx].SubgrupoCatalogoId, el.SubgrupoCatalogoId, proveedorId, cuentasSubgrupo, detalleCuenta)...)

		}

	}

	for _, el := range actualizarVl {
		if idx := findElementoInArrayE(originales, el.Id); idx > -1 {

			if detalleCuenta_, err := fillCuentas(detalleCuenta, []string{cuentasSubgrupo[originales[idx].SubgrupoCatalogoId].CuentaCreditoId,
				cuentasSubgrupo[el.SubgrupoCatalogoId].CuentaCreditoId, cuentasSubgrupo[originales[idx].SubgrupoCatalogoId].CuentaDebitoId, cuentasSubgrupo[el.SubgrupoCatalogoId].CuentaDebitoId}); err != nil {
			} else {
				detalleCuenta = detalleCuenta_
			}

			movimientos = append(movimientos, generaTrContable(el.ValorTotal-originales[idx].ValorTotal, consecutivo, tipoMovimiento,
				movDebito, movCredito, 0, el.SubgrupoCatalogoId, proveedorId, cuentasSubgrupo, detalleCuenta)...)

		}

	}

	return movimientos, nil

}

// separarElementosPorSalida Separa los elementos según el tipo de ajuste de cada uno y los agrupa según la salida. Además retorna los elementos actualizados
func separarElementosPorSalida(elementos []*models.ElementosMovimiento, updateVls, updateSg, updateMp []*models.DetalleElemento_) (elementosSalidas map[int]*models.ElementosPorActualizarSalida, pendientes_ []*models.DetalleElemento_, actualizados []*models.ElementosMovimiento, outputError map[string]interface{}) {

	funcion := "separarElementosPorSalida"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	elementosSalidas = make(map[int]*models.ElementosPorActualizarSalida)
	for _, el := range elementos {

		if len(updateMp) > 0 {
			if idx := findElementoInArrayD(updateMp, el.ElementoActaId); idx > -1 {
				if updateMp[idx].ValorResidual == el.ValorResidual && updateMp[idx].VidaUtil == el.VidaUtil {
					updateMp = append(updateMp[:idx], updateMp[idx+1:]...)
				} else {
					elemento_ := creaNuevoElementoMovimiento(updateMp[idx], el)
					actualizados = append(actualizados, elemento_)
				}
				continue
			}
		} else if el.MovimientoId.EstadoMovimientoId.Nombre != "Salida Aprobada" {
			continue
		}

		if len(updateSg) > 0 {
			if idx := findElementoInArrayD(updateSg, el.ElementoActaId); idx > -1 {
				if elementosSalidas[el.MovimientoId.Id] == nil {
					elementosSalidas[el.MovimientoId.Id] = new(models.ElementosPorActualizarSalida)
					elementosSalidas[el.MovimientoId.Id].Salida = el.MovimientoId
				} else {
					elemento_ := creaNuevoElementoMovimiento(updateSg[idx], el)
					actualizados = append(actualizados, elemento_)
				}

				elementosSalidas[el.MovimientoId.Id].UpdateSg = append(elementosSalidas[el.MovimientoId.Id].UpdateSg, updateSg[idx])
			}
		} else if len(updateVls) > 0 {
			if idx := findElementoInArrayD(updateVls, el.ElementoActaId); idx > -1 {
				if elementosSalidas[el.MovimientoId.Id] == nil {
					elementosSalidas[el.MovimientoId.Id] = new(models.ElementosPorActualizarSalida)
					elementosSalidas[el.MovimientoId.Id].Salida = el.MovimientoId
				} else {
					elemento_ := creaNuevoElementoMovimiento(updateVls[idx], el)
					actualizados = append(actualizados, elemento_)
				}

				elementosSalidas[el.MovimientoId.Id].UpdateVls = append(elementosSalidas[el.MovimientoId.Id].UpdateVls, updateVls[idx])
			}
		}

	}
	return elementosSalidas, updateMp, actualizados, nil

}

// determinarDeltaActa Separa elementos según el ajuste
func determinarDeltaActa(org *models.Elemento, nvo *models.DetalleElemento_) (msc, vls, sg bool, outputError map[string]interface{}) {

	funcion := "determinarDeltaActa"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	if org.SubgrupoCatalogoId != nvo.SubgrupoCatalogoId {

		urlcrud := "fields=TipoBienId&sortby=Id&order=desc&query=Activo:true,SubgrupoId__Id:" + strconv.Itoa(nvo.SubgrupoCatalogoId)
		if detalleSubgrupo_, err := catalogoElementosHelper.GetAllDetalleSubgrupo(urlcrud); err != nil {
			return false, false, false, err
		} else if len(detalleSubgrupo_) == 0 {
			err := "len(detalleSubgrupo_) = 0"
			eval := " - catalogoElementosHelper.GetAllDetalleSubgrupo(urlcrud)"
			return false, false, false, errorctrl.Error(funcion+eval, err, "500")
		} else {
			if detalleSubgrupo_[0].TipoBienId.NecesitaPlaca && nvo.Placa != "" {
				ctxPlaca, _ := beego.AppConfig.Int("contxtPlaca")
				if placa_, _, err := utilsHelper.GetConsecutivo("%05.0f", ctxPlaca, "Registro Placa Arka"); err != nil {
					return false, false, false, err
				} else {
					year, month, day := time.Now().Date()
					nvo.Placa = utilsHelper.FormatConsecutivo(fmt.Sprintf("%04d%02d%02d", year, month, day), placa_, "")
				}
			} else if !detalleSubgrupo_[0].TipoBienId.NecesitaPlaca && nvo.Placa != "" {
				nvo.Placa = ""
			}

		}

		nvo.Activo = true
		sg = true

	} else if org.ValorTotal != nvo.ValorTotal {
		nvo.Activo = true
		vls = true

	} else if org.Nombre != nvo.Nombre || org.Marca != nvo.Marca ||
		org.Serie != nvo.Serie || org.UnidadMedida != nvo.UnidadMedida ||
		org.Cantidad != nvo.Cantidad || org.ValorUnitario != nvo.ValorUnitario ||
		org.Subtotal != nvo.Subtotal || org.Descuento != nvo.Descuento ||
		org.PorcentajeIvaId != nvo.PorcentajeIvaId ||
		org.ValorIva != nvo.ValorIva || org.ValorFinal != nvo.ValorFinal {

		nvo.Activo = true
		msc = true

	}

	return msc, vls, sg, nil

}

func creaNuevoElementoMovimiento(nuevo *models.DetalleElemento_, org *models.ElementosMovimiento) *models.ElementosMovimiento {

	org.SaldoCantidad = float64(nuevo.Cantidad)
	org.SaldoValor = nuevo.ValorTotal
	org.Unidad = float64(nuevo.Cantidad)
	org.ValorUnitario = nuevo.ValorTotal / float64(nuevo.Cantidad)
	org.ValorTotal = nuevo.ValorTotal
	org.VidaUtil = nuevo.VidaUtil
	org.ValorResidual = nuevo.ValorResidual

	return org

}
