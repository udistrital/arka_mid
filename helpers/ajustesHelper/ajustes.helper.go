package ajustesHelper

import (
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/catalogoElementosHelper"
	"github.com/udistrital/arka_mid/helpers/cuentasContablesHelper"
	"github.com/udistrital/arka_mid/helpers/entradaHelper"
	"github.com/udistrital/arka_mid/helpers/movimientosArkaHelper"
	"github.com/udistrital/arka_mid/helpers/movimientosContablesMidHelper"
	"github.com/udistrital/arka_mid/helpers/parametrosHelper"
	"github.com/udistrital/arka_mid/helpers/salidaHelper"
	"github.com/udistrital/arka_mid/helpers/tercerosHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
)

func PostAjuste(trContable *models.PreTrAjuste) (movimiento *models.Movimiento, outputError map[string]interface{}) {

	funcion := "PostAjuste"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var query string
	movimiento = new(models.Movimiento)
	detalle := new(models.FormatoAjuste)

	ctxConsecutivo, _ := beego.AppConfig.Int("contxtAjusteCons")
	if consecutivo, _, err := utilsHelper.GetConsecutivo("%05.0f", ctxConsecutivo, "Ajuste Contable Arka"); err != nil {
		return nil, err
	} else {
		consecutivo = utilsHelper.FormatConsecutivo(getTipoComprobanteAjustes()+"-", consecutivo, fmt.Sprintf("%s%04d", "-", time.Now().Year()))
		detalle.Consecutivo = consecutivo
		detalle.PreTrAjuste = trContable
	}

	if jsonData, err := json.Marshal(detalle); err != nil {
		logs.Error(err)
		eval := " - json.Marshal(detalle)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	} else {
		movimiento.Detalle = string(jsonData[:])
	}

	query = "query=Nombre:" + url.QueryEscape("Ajuste Contable")
	if fm, err := movimientosArkaHelper.GetAllFormatoTipoMovimiento(query); err != nil {
		return nil, err
	} else {
		movimiento.FormatoTipoMovimientoId = fm[0]
	}

	if sm, err := movimientosArkaHelper.GetAllEstadoMovimiento(url.QueryEscape("Ajuste En Trámite")); err != nil {
		return nil, err
	} else {
		movimiento.EstadoMovimientoId = sm[0]
	}

	movimiento.Activo = true

	if res, err := movimientosArkaHelper.PostMovimiento(movimiento); err != nil {
		return nil, err
	} else {
		return res, nil
	}

}

// GetDetalleAjuste Consulta los detalles de un ajuste contable
func GetDetalleAjuste(id int) (Ajuste *models.DetalleAjuste, outputError map[string]interface{}) {

	funcion := "GetDetalleAjuste"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		movimiento         *models.Movimiento
		detalle            *models.FormatoAjuste
		movimientos        []*models.PreMovAjuste
		parametroCreditoId int
		parametroDebitoId  int
	)

	Ajuste = new(models.DetalleAjuste)

	if movimiento_, err := movimientosArkaHelper.GetMovimientoById(id); err != nil {
		return nil, err
	} else {
		movimiento = movimiento_
	}

	if err := json.Unmarshal([]byte(movimiento.Detalle), &detalle); err != nil {
		logs.Error(err)
		eval := " - json.Unmarshal([]byte(movimiento.Detalle), &detalle)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	if db_, cr_, err := parametrosHelper.GetParametrosDebitoCredito(); err != nil {
		return nil, err
	} else {
		parametroDebitoId = db_
		parametroCreditoId = cr_
	}

	if detalle.PreTrAjuste != nil && detalle.TrContableId == 0 {
		movimientos = detalle.PreTrAjuste.Movimientos
	} else if detalle.PreTrAjuste == nil && detalle.TrContableId > 0 {
		if tr, err := movimientosContablesMidHelper.GetTransaccion(detalle.TrContableId, "consecutivo", true); err != nil {
			return nil, err
		} else {
			for _, mov := range tr.Movimientos {
				mov_ := new(models.PreMovAjuste)
				mov_.Cuenta = mov.CuentaId
				mov_.TerceroId = *mov.TerceroId
				mov_.Descripcion = mov.Descripcion

				if mov.TipoMovimientoId == parametroCreditoId {
					mov_.Credito = mov.Valor
					mov_.Debito = 0
				} else if mov.TipoMovimientoId == parametroDebitoId {
					mov_.Debito = mov.Valor
					mov_.Credito = 0
				}

				movimientos = append(movimientos, mov_)
			}
		}
	}

	movs := make([]*models.DetalleMovimientoContable, 0)
	for _, mov := range movimientos {
		mov_ := new(models.DetalleMovimientoContable)
		var cta *models.DetalleCuenta

		if ctaCr_, err := cuentasContablesHelper.GetCuentaContable(mov.Cuenta); err != nil {
			return nil, err
		} else {
			if err := formatdata.FillStruct(ctaCr_, &cta); err != nil {
				logs.Error(err)
				eval := " - formatdata.FillStruct(ctaCr_, &ctaCr)"
				return nil, errorctrl.Error(funcion+eval, err, "500")
			}
			mov_.Cuenta = cta
		}

		if mov.TerceroId > 0 {
			if tercero_, err := tercerosHelper.GetTerceroById(mov.TerceroId); err != nil {
				return nil, err
			} else {
				mov_.TerceroId = tercero_
			}
		}
		mov_.Credito = mov.Credito
		mov_.Debito = mov.Debito
		mov_.Descripcion = mov.Descripcion
		movs = append(movs, mov_)
	}

	Ajuste.TrContable = movs
	Ajuste.Movimiento = movimiento

	return Ajuste, nil

}

// AprobarAjuste Realiza la transacción contable correspondiente
func AprobarAjuste(id int) (movimiento *models.Movimiento, outputError map[string]interface{}) {

	funcion := "AprobarAjuste"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		detalle            *models.FormatoAjuste
		parametroCreditoId int
		parametroDebitoId  int
	)

	if movimiento_, err := movimientosArkaHelper.GetMovimientoById(id); err != nil {
		return nil, err
	} else {
		movimiento = movimiento_
	}

	if err := json.Unmarshal([]byte(movimiento.Detalle), &detalle); err != nil {
		logs.Error(err)
		eval := " - json.Unmarshal([]byte(movimiento.Detalle), &detalle)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	if db_, cr_, err := parametrosHelper.GetParametrosDebitoCredito(); err != nil {
		return nil, err
	} else {
		parametroDebitoId = db_
		parametroCreditoId = cr_
	}

	movs := make([]*models.MovimientoTransaccion, 0)
	for _, mov := range detalle.PreTrAjuste.Movimientos {
		mov_ := new(models.MovimientoTransaccion)
		var cta *models.DetalleCuenta

		if ctaCr_, err := cuentasContablesHelper.GetCuentaContable(mov.Cuenta); err != nil {
			return nil, err
		} else {
			if err := formatdata.FillStruct(ctaCr_, &cta); err != nil {
				logs.Error(err)
				eval := " - formatdata.FillStruct(ctaCr_, &ctaCr)"
				return nil, errorctrl.Error(funcion+eval, err, "500")
			}
			mov_.CuentaId = cta.Id
			mov_.NombreCuenta = cta.Nombre
		}

		if mov.TerceroId > 0 {
			mov_.TerceroId = &mov.TerceroId
		}

		if mov.Credito > 0 {
			mov_.TipoMovimientoId = parametroCreditoId
			mov_.Valor = mov.Credito
		} else if mov.Debito > 0 {
			mov_.TipoMovimientoId = parametroDebitoId
			mov_.Valor = mov.Debito
		}

		mov_.Activo = true
		movs = append(movs, mov_)
	}

	transaccion := new(models.TransaccionMovimientos)

	if _, consecutivoId_, err := utilsHelper.GetConsecutivo("%05.0f", 1, "CNTB"); err != nil {
		return nil, outputError
	} else {
		transaccion.ConsecutivoId = consecutivoId_
	}

	transaccion.Movimientos = movs
	transaccion.FechaTransaccion = time.Now()
	transaccion.Activo = true
	transaccion.Etiquetas = ""
	transaccion.Descripcion = ""

	if tr, err := movimientosContablesMidHelper.PostTrContable(transaccion); err != nil {
		return nil, err
	} else {
		detalle.TrContableId = tr.ConsecutivoId
		detalle.PreTrAjuste = nil
		detalle.RazonRechazo = ""
	}

	if sm, err := movimientosArkaHelper.GetAllEstadoMovimiento(url.QueryEscape("Ajuste Aprobado")); err != nil {
		return nil, err
	} else {
		movimiento.EstadoMovimientoId = sm[0]
	}

	if jsonData, err := json.Marshal(detalle); err != nil {
		logs.Error(err)
		eval := " - json.Marshal(detalle)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	} else {
		movimiento.Detalle = string(jsonData[:])
	}

	if movimiento_, err := movimientosArkaHelper.PutMovimiento(movimiento, movimiento.Id); err != nil {
		return nil, err
	} else {
		movimiento = movimiento_
	}

	return movimiento, nil
}

func GenerarAjuste(elementos []*models.DetalleElemento_) (resultado []*models.MovimientoTransaccion, outputError map[string]interface{}) {

	funcion := "GenerarAjuste"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		query                string
		idsEl                []int
		entrada              *models.Movimiento
		orgActa              []*models.Elemento
		elementosSalida      map[int]*models.ElementosPorActualizarSalida
		// updateMsc            []*models.DetalleElemento_
		updateVls            []*models.DetalleElemento_
		updateSg             []*models.DetalleElemento_
		updateMp             []*models.DetalleElemento_
		movimientos          []*models.MovimientoTransaccion
		tipoMovimientoSalida int
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
		// updateMsc = msc
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
			if elementosSalida_, updateMp_, err := separarElementosPorSalida(elementos_, updateVls, updateSg, updateMp); err != nil {
				return nil, err
			} else {
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

		if movsSalida, err := calcularAjusteMovimiento(orgActa, elms.UpdateVls, elms.UpdateSg, tipoMovimientoSalida, consecutivo, funcionario, "Salida"); err != nil {
			return nil, err
		} else {
			movimientos = append(movimientos, movsSalida...)
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
			} else if msc_ != nil {
				msc = append(msc, msc_)
			} else if vls_ != nil {
				vls = append(vls, vls_)
			} else if sg_ != nil {
				sg = append(sg, sg_)
			} else if mediciones {
				mp = append(mp, actualizados[idx])
			}
		}
	}

	return msc, vls, sg, mp, nil

}

// separarElementosPorSalida Separa los elementos según el tipo de ajuste de cada uno y los agrupa según la salida
func separarElementosPorSalida(elementos []*models.ElementosMovimiento, updateVls, updateSg, updateMp []*models.DetalleElemento_) (elementosSalidas map[int]*models.ElementosPorActualizarSalida, pendientes_ []*models.DetalleElemento_, outputError map[string]interface{}) {

	funcion := "separarElementosPorModificacion"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	elementosSalidas = make(map[int]*models.ElementosPorActualizarSalida)
	for _, el := range elementos {

		if len(updateMp) > 0 {
			if idx := findElementoInArrayD(updateMp, el.ElementoActaId); idx > -1 {
				if updateMp[idx].ValorResidual == el.ValorResidual && updateMp[idx].VidaUtil == el.VidaUtil {
					updateMp = append(updateMp[:idx], updateMp[idx+1:]...)
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

				elementosSalidas[el.MovimientoId.Id].UpdateSg = append(elementosSalidas[el.MovimientoId.Id].UpdateSg, updateSg[idx])
				updateSg = append(updateSg[:idx], updateSg[idx+1:]...)
			}
		} else if len(updateVls) > 0 {
			if idx := findElementoInArrayD(updateVls, el.ElementoActaId); idx > -1 {
				if elementosSalidas[el.MovimientoId.Id] == nil {
					elementosSalidas[el.MovimientoId.Id] = new(models.ElementosPorActualizarSalida)
					elementosSalidas[el.MovimientoId.Id].Salida = el.MovimientoId
				}

				elementosSalidas[el.MovimientoId.Id].UpdateVls = append(elementosSalidas[el.MovimientoId.Id].UpdateVls, updateVls[idx])
				updateVls = append(updateVls[:idx], updateVls[idx+1:]...)
			}
		}

	}
	return elementosSalidas, updateMp, nil

}

// determinarDeltaActa Separa elementos según el ajuste
func determinarDeltaActa(org *models.Elemento, nvo *models.DetalleElemento_) (msc, vls, sg *models.DetalleElemento_, outputError map[string]interface{}) {

	funcion := "determinarDeltaActa"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	if org.SubgrupoCatalogoId != nvo.SubgrupoCatalogoId {

		urlcrud := "fields=TipoBienId&sortby=Id&order=desc&query=Activo:true,SubgrupoId__Id:" + strconv.Itoa(nvo.SubgrupoCatalogoId)
		if detalleSubgrupo_, err := catalogoElementosHelper.GetAllDetalleSubgrupo(urlcrud); err != nil {
			return nil, nil, nil, err
		} else if len(detalleSubgrupo_) == 0 {
			err := "len(detalleSubgrupo_) = 0"
			eval := " - catalogoElementosHelper.GetAllDetalleSubgrupo(urlcrud)"
			return nil, nil, nil, errorctrl.Error(funcion+eval, err, "500")
		} else {
			if detalleSubgrupo_[0].TipoBienId.NecesitaPlaca && nvo.Placa != "" {
				ctxPlaca, _ := beego.AppConfig.Int("contxtPlaca")
				if placa_, _, err := utilsHelper.GetConsecutivo("%05.0f", ctxPlaca, "Registro Placa Arka"); err != nil {
					return nil, nil, nil, err
				} else {
					year, month, day := time.Now().Date()
					nvo.Placa = utilsHelper.FormatConsecutivo(fmt.Sprintf("%04d%02d%02d", year, month, day), placa_, "")
				}
			} else if !detalleSubgrupo_[0].TipoBienId.NecesitaPlaca && nvo.Placa != "" {
				nvo.Placa = ""
			}

		}

		nvo.Activo = true
		sg = nvo

	} else if org.ValorTotal != nvo.ValorTotal {
		nvo.Activo = true
		vls = nvo

	} else if org.Nombre != nvo.Nombre || org.Marca != nvo.Marca ||
		org.Serie != nvo.Serie || org.UnidadMedida != nvo.UnidadMedida ||
		org.Cantidad != nvo.Cantidad || org.ValorUnitario != nvo.ValorUnitario ||
		org.Subtotal != nvo.Subtotal || org.Descuento != nvo.Descuento ||
		org.PorcentajeIvaId != nvo.PorcentajeIvaId ||
		org.ValorIva != nvo.ValorIva || org.ValorFinal != nvo.ValorFinal {

		nvo.Activo = true
		msc = nvo

	}

	return msc, vls, sg, nil

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

	query := "limit=-1&fields=CuentaDebitoId,CuentaCreditoId,SubgrupoId&sortby=Id&order=desc&"
	query += "query=SubtipoMovimientoId:" + strconv.Itoa(movimientoId) + ",Activo:true,SubgrupoId__Id__in:"
	query += url.QueryEscape(utilsHelper.ArrayToString(ids, "|"))
	if cuentas_, err := catalogoElementosHelper.GetAllCuentasSubgrupo(query); err != nil {
		return nil, err
	} else {
		for _, cuenta := range cuentas_ {
			cuentasSubgrupo[cuenta.SubgrupoId.Id] = cuenta
		}
	}

	dsc := getDescripcionMovContable(tipoMovimiento, consecutivo)
	for _, el := range actualizarSg {
		if idx := findElementoInArrayE(originales, el.Id); idx > -1 {

			if detalleCuenta_, err := fillCuentas(detalleCuenta, []string{cuentasSubgrupo[originales[idx].SubgrupoCatalogoId].CuentaCreditoId,
				cuentasSubgrupo[el.SubgrupoCatalogoId].CuentaCreditoId, cuentasSubgrupo[originales[idx].SubgrupoCatalogoId].CuentaDebitoId, cuentasSubgrupo[el.SubgrupoCatalogoId].CuentaDebitoId}); err != nil {
				return nil, err
			} else {
				detalleCuenta = detalleCuenta_
			}

			if cuentasSubgrupo[el.SubgrupoCatalogoId].CuentaCreditoId != cuentasSubgrupo[originales[idx].SubgrupoCatalogoId].CuentaCreditoId {

				movimientoR := asientoContable.CreaMovimiento(originales[idx].ValorTotal, dsc, proveedorId, detalleCuenta[cuentasSubgrupo[originales[idx].SubgrupoCatalogoId].CuentaCreditoId], movDebito)
				movimiento := asientoContable.CreaMovimiento(el.ValorTotal, dsc, proveedorId, detalleCuenta[cuentasSubgrupo[el.SubgrupoCatalogoId].CuentaCreditoId], movCredito)
				movimientos = append(movimientos, movimientoR, movimiento)

			} else if el.ValorTotal != originales[idx].ValorTotal {

				tipoMovimiento := movCredito
				if el.ValorTotal < originales[idx].ValorTotal {
					tipoMovimiento = movDebito
				}

				movimiento := asientoContable.CreaMovimiento(math.Abs(el.ValorTotal-originales[idx].ValorTotal), dsc, proveedorId, detalleCuenta[cuentasSubgrupo[el.SubgrupoCatalogoId].CuentaCreditoId], tipoMovimiento)
				movimientos = append(movimientos, movimiento)

			}

			if cuentasSubgrupo[el.SubgrupoCatalogoId].CuentaDebitoId != cuentasSubgrupo[originales[idx].SubgrupoCatalogoId].CuentaDebitoId {

				movimientoR := asientoContable.CreaMovimiento(originales[idx].ValorTotal, dsc, proveedorId, detalleCuenta[cuentasSubgrupo[originales[idx].SubgrupoCatalogoId].CuentaDebitoId], movCredito)
				movimiento := asientoContable.CreaMovimiento(el.ValorTotal, dsc, proveedorId, detalleCuenta[cuentasSubgrupo[el.SubgrupoCatalogoId].CuentaDebitoId], movDebito)
				movimientos = append(movimientos, movimientoR, movimiento)

			} else if el.ValorTotal != originales[idx].ValorTotal {

				tipoMovimiento := movDebito
				if el.ValorTotal < originales[idx].ValorTotal {
					tipoMovimiento = movCredito
				}

				movimiento := asientoContable.CreaMovimiento(math.Abs(el.ValorTotal-originales[idx].ValorTotal), dsc, proveedorId, detalleCuenta[cuentasSubgrupo[el.SubgrupoCatalogoId].CuentaDebitoId], tipoMovimiento)
				movimientos = append(movimientos, movimiento)

			}
		}

	}

	for _, el := range actualizarVl {
		if idx := findElementoInArrayE(originales, el.Id); idx > -1 {

			if detalleCuenta_, err := fillCuentas(detalleCuenta, []string{cuentasSubgrupo[originales[idx].SubgrupoCatalogoId].CuentaCreditoId,
				cuentasSubgrupo[el.SubgrupoCatalogoId].CuentaCreditoId, cuentasSubgrupo[originales[idx].SubgrupoCatalogoId].CuentaDebitoId, cuentasSubgrupo[el.SubgrupoCatalogoId].CuentaDebitoId}); err != nil {
			} else {
				detalleCuenta = detalleCuenta_
			}

			if el.ValorTotal != originales[idx].ValorTotal {

				tipoMovimientoC := movCredito
				tipoMovimientoD := movDebito

				if el.ValorTotal < originales[idx].ValorTotal {
					tipoMovimientoC = movDebito
					tipoMovimientoD = movCredito
				}

				movimientoC := asientoContable.CreaMovimiento(math.Abs(el.ValorTotal-originales[idx].ValorTotal), dsc, proveedorId, detalleCuenta[cuentasSubgrupo[el.SubgrupoCatalogoId].CuentaCreditoId], tipoMovimientoC)
				movimientoD := asientoContable.CreaMovimiento(math.Abs(el.ValorTotal-originales[idx].ValorTotal), dsc, proveedorId, detalleCuenta[cuentasSubgrupo[el.SubgrupoCatalogoId].CuentaDebitoId], tipoMovimientoD)
				movimientos = append(movimientos, movimientoC, movimientoD)
			}

		}

	}

	return movimientos, nil

}

// fillCuentas Consulta el detalle de una serie de cuentas
func fillCuentas(cuentas map[string]*models.CuentaContable, cuentas_ []string) (cuentasCompletas map[string]*models.CuentaContable, outputError map[string]interface{}) {

	funcion := "fillCuentas"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	for _, id := range cuentas_ {
		if _, ok := cuentas[id]; !ok {
			if cta_, err := cuentasContablesHelper.GetCuentaContable(id); err != nil {
				return nil, err
			} else {
				cuentas[id] = cta_
			}
		}
	}

	return cuentas, nil

}

// findElementoInArray Retorna la posicion en que se encuentra el id específicado
func findElementoInArrayD(elementos []*models.DetalleElemento_, id int) (i int) {
	for i, el_ := range elementos {
		if int(el_.Id) == id {
			return i
		}
	}
	return -1
}

// findElementoInArray Retorna la posicion en que se encuentra el id específicado
func findElementoInArrayE(elementos []*models.Elemento, id int) (i int) {
	for i, el_ := range elementos {
		if int(el_.Id) == id {
			return i
		}
	}
	return -1
}

func getTipoComprobanteAjustes() string {
	return "N20"
}

func getDescripcionMovContable(tipoMovimiento, consecutivo string) string {
	return "Ajuste contable " + tipoMovimiento + " " + consecutivo
}
