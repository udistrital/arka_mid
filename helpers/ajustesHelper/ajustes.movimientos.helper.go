package ajustesHelper

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	crudActas "github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// separarElementosPorModificacion Separa los elementos según se deba modificar Subgrupo, Valores, Misceláneos o Mediciones posteriores
// msc: Cambios miscelaneos, elementos a los que unicamente se les debe ajustar nombre, marca, serie, unidad.
// vls: Cambios a valores, elementos a los que se les debe cambiar el valor total.
// sg: Cambia el subgrupo del elemento. Se ajusta la placa de acuerdo al nuevo subgrupo.
// mp: Cambian los parametros de las mediciones posteriores. vida util o valor residual.
func separarElementosPorModificacion(originales []*models.Elemento,
	actualizados []*models.DetalleElemento_,
	mediciones bool) (
	msc, vls, sg, mp []*models.DetalleElemento_,
	outputError map[string]interface{}) {

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

// calcularAjusteMovimiento Calcula la transacción contable generada a partir de los elementos y el cambio de cada uno.
// actualizarVl: Elementos para actualizar los montos de las transacciones contables.
// actualizarSg: Elementos para actualizar el subgrupo y por tanto, pueden cambiar las cuentas.
func calcularAjusteMovimiento(originales []*models.Elemento,
	actualizarVl, actualizarSg []*models.DetalleElemento_,
	movimientoId, proveedorId int,
	consecutivo, tipoMovimiento string) (movimientos []*models.MovimientoTransaccion,
	outputError map[string]interface{}) {

	funcion := "calcularAjusteMovimiento"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		ids             []int
		movDebito       int
		movCredito      int
		cuentasSubgrupo map[int]*models.CuentasSubgrupo
		detalleCuenta   map[string]*models.CuentaContable
	)

	detalleCuenta = make(map[string]*models.CuentaContable)
	if db_, cr_, err := parametros.GetParametrosDebitoCredito(); err != nil {
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

			if detalleCuenta_, err := fillCuentas(detalleCuenta,
				[]string{cuentasSubgrupo[originales[idx].SubgrupoCatalogoId].CuentaCreditoId,
					cuentasSubgrupo[el.SubgrupoCatalogoId].CuentaCreditoId,
					cuentasSubgrupo[originales[idx].SubgrupoCatalogoId].CuentaDebitoId,
					cuentasSubgrupo[el.SubgrupoCatalogoId].CuentaDebitoId}); err != nil {
				return nil, err
			} else {
				detalleCuenta = detalleCuenta_
			}

			movimientos = append(movimientos,
				generaTrContable(originales[idx].ValorTotal, el.ValorTotal,
					consecutivo,
					tipoMovimiento,
					movDebito,
					movCredito,
					originales[idx].SubgrupoCatalogoId,
					el.SubgrupoCatalogoId,
					proveedorId,
					cuentasSubgrupo,
					detalleCuenta)...)
		}

	}

	for _, el := range actualizarVl {
		if idx := findElementoInArrayE(originales, el.Id); idx > -1 {

			if detalleCuenta_, err := fillCuentas(detalleCuenta,
				[]string{cuentasSubgrupo[originales[idx].SubgrupoCatalogoId].CuentaCreditoId,
					cuentasSubgrupo[el.SubgrupoCatalogoId].CuentaCreditoId,
					cuentasSubgrupo[originales[idx].SubgrupoCatalogoId].CuentaDebitoId,
					cuentasSubgrupo[el.SubgrupoCatalogoId].CuentaDebitoId}); err != nil {
			} else {
				detalleCuenta = detalleCuenta_
			}

			movimientos = append(movimientos,
				generaTrContable(originales[idx].ValorTotal, el.ValorTotal,
					consecutivo,
					tipoMovimiento,
					movDebito,
					movCredito,
					0,
					el.SubgrupoCatalogoId,
					proveedorId,
					cuentasSubgrupo,
					detalleCuenta)...)

		}

	}

	return

}

// submitUpdates Actualiza los registros relacionados a las novedades y elementos
func submitUpdates(elementosActa []*models.Elemento,
	elementosMovimiento []*models.ElementosMovimiento,
	novedades []*models.NovedadElemento) (outputError map[string]interface{}) {

	for _, el := range elementosActa {
		if _, err := crudActas.PutElemento(el, el.Id); err != nil {
			return err
		}
	}

	for _, el := range elementosMovimiento {
		if _, err := movimientosArka.PutElementosMovimiento(el, el.Id); err != nil {
			return err
		}
	}

	for _, nv := range novedades {
		if _, err := movimientosArka.PutNovedadElemento(nv, nv.Id); err != nil {
			return err
		}
	}

	return nil

}

// separarElementosPorSalida Separa los elementos según el tipo de ajuste de cada uno y los agrupa según la salida. Además retorna los elementos actualizados
// updateVl: Elementos para actualizar los montos de las transacciones contables.
// updateMp: Elementos para actualizar los parametros iniciales de las mediciones posteriores.
// updateSg: Elementos para actualizar el subgrupo y por tanto, pueden cambiar las cuentas.
func separarElementosPorSalida(elementos []*models.ElementosMovimiento,
	updateVls, updateSg, updateMp []*models.DetalleElemento_) (
	elementosSalidas map[int]*models.ElementosPorActualizarSalida,
	pendientes_ []*models.DetalleElemento_,
	actualizados []*models.ElementosMovimiento,
	outputError map[string]interface{}) {

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
				}
				elemento_ := creaNuevoElementoMovimiento(updateSg[idx], el)
				actualizados = append(actualizados, elemento_)

				elementosSalidas[el.MovimientoId.Id].UpdateSg = append(elementosSalidas[el.MovimientoId.Id].UpdateSg, updateSg[idx])
			}
		} else if len(updateVls) > 0 {
			if idx := findElementoInArrayD(updateVls, el.ElementoActaId); idx > -1 {
				if elementosSalidas[el.MovimientoId.Id] == nil {
					elementosSalidas[el.MovimientoId.Id] = new(models.ElementosPorActualizarSalida)
					elementosSalidas[el.MovimientoId.Id].Salida = el.MovimientoId
				}
				elemento_ := creaNuevoElementoMovimiento(updateVls[idx], el)
				actualizados = append(actualizados, elemento_)

				elementosSalidas[el.MovimientoId.Id].UpdateVls = append(elementosSalidas[el.MovimientoId.Id].UpdateVls, updateVls[idx])
			}
		}

	}
	return elementosSalidas, updateMp, actualizados, nil

}

// generarMovimientoAjuste Crea el registro del movimiento de inventario y contable resultantes del ajuste
func generarMovimientoAjuste(sg, vls, msc, mp []*models.DetalleElemento_,
	movContables []*models.MovimientoTransaccion) (movimiento *models.Movimiento,
	trContable *models.TransaccionMovimientos, outputError map[string]interface{}) {

	funcion := "generarMovimientoAjuste"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var consecutivo models.Consecutivo
	movimiento = new(models.Movimiento)
	detalle := new(models.FormatoAjusteAutomatico)
	query := "query=Nombre:" + url.QueryEscape("Ajuste Automático")
	if fm, err := movimientosArka.GetAllFormatoTipoMovimiento(query); err != nil {
		return nil, nil, err
	} else {
		movimiento.FormatoTipoMovimientoId = fm[0]
	}

	if sm, err := movimientosArka.GetAllEstadoMovimiento(url.QueryEscape("Ajuste Aprobado")); err != nil {
		return nil, nil, err
	} else {
		movimiento.EstadoMovimientoId = sm[0]
	}

	var ids []int
	for _, el := range append(sg, append(vls, append(msc, mp...)...)...) {
		ids = append(ids, el.Id)
	}

	ctxConsecutivo, _ := beego.AppConfig.Int("contxtAjusteCons")
	if err := consecutivos.Get(ctxConsecutivo, "Ajuste automático Arka", &consecutivo); err != nil {
		return nil, nil, err
	}

	detalle.Consecutivo = consecutivos.Format("%05d", getTipoComprobanteAjustes(), &consecutivo)
	detalle.ConsecutivoId = consecutivo.Id
	detalle.Elementos = ids

	if len(movContables) > 0 {
		trContable = new(models.TransaccionMovimientos)
		trContable.Movimientos = movContables
		trContable.ConsecutivoId = consecutivo.Id
		trContable.FechaTransaccion = time.Now()
		trContable.Activo = true
		trContable.Etiquetas = ""
		trContable.Descripcion = "Ajuste contable almacén"

		if _, err := movimientosContables.PostTrContable(trContable); err != nil {
			return nil, nil, err
		}
	}

	if jsonData, err := json.Marshal(detalle); err != nil {
		logs.Error(err)
		eval := " - json.Marshal(detalle)"
		return nil, nil, errorctrl.Error(funcion+eval, err, "500")
	} else {
		movimiento.Detalle = string(jsonData[:])
	}

	movimiento.Activo = true

	if err := movimientosArka.PostMovimiento(movimiento); err != nil {
		return nil, nil, err
	}

	return movimiento, trContable, nil

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
