package ajustesHelper

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
)

func PostAjuste(trContable *models.PreTrAjuste) (movimiento *models.Movimiento, outputError map[string]interface{}) {

	funcion := "PostAjuste"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		query       string
		consecutivo models.Consecutivo
	)
	movimiento = new(models.Movimiento)
	detalle := new(models.FormatoAjuste)

	ctxConsecutivo, _ := beego.AppConfig.Int("contxtAjusteCons")
	if err := consecutivos.Get(ctxConsecutivo, "Ajuste Contable Arka", &consecutivo); err != nil {
		return nil, err
	}

	detalle.Consecutivo = consecutivos.Format("%05d", getTipoComprobanteAjustes(), &consecutivo)
	detalle.ConsecutivoId = consecutivo.Id
	detalle.PreTrAjuste = trContable

	if jsonData, err := json.Marshal(detalle); err != nil {
		logs.Error(err)
		eval := " - json.Marshal(detalle)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	} else {
		movimiento.Detalle = string(jsonData[:])
	}

	query = "query=Nombre:" + url.QueryEscape("Ajuste Contable")
	if fm, err := movimientosArka.GetAllFormatoTipoMovimiento(query); err != nil {
		return nil, err
	} else {
		movimiento.FormatoTipoMovimientoId = fm[0]
	}

	if sm, err := movimientosArka.GetAllEstadoMovimiento("query=Nombre:" + url.QueryEscape("Ajuste En Trámite")); err != nil {
		return nil, err
	} else {
		movimiento.EstadoMovimientoId = sm[0]
	}

	movimiento.Activo = true

	if err := movimientosArka.PostMovimiento(movimiento); err != nil {
		return nil, err
	}

	return movimiento, nil

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

	if movimiento_, err := movimientosArka.GetMovimientoById(id); err != nil {
		return nil, err
	} else {
		movimiento = movimiento_
	}

	if err := json.Unmarshal([]byte(movimiento.Detalle), &detalle); err != nil {
		logs.Error(err)
		eval := " - json.Unmarshal([]byte(movimiento.Detalle), &detalle)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	if db_, cr_, err := parametros.GetParametrosDebitoCredito(); err != nil {
		return nil, err
	} else {
		parametroDebitoId = db_
		parametroCreditoId = cr_
	}

	if detalle.PreTrAjuste != nil && detalle.TrContableId == 0 {
		movimientos = detalle.PreTrAjuste.Movimientos
	} else if detalle.PreTrAjuste == nil && detalle.TrContableId > 0 {
		if tr, err := movimientosContables.GetTransaccion(detalle.TrContableId, "consecutivo", true); err != nil {
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

		if ctaCr_, err := cuentasContables.GetCuentaContable(mov.Cuenta); err != nil {
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
			if tercero, err := terceros.GetNombreTerceroById(mov.TerceroId); err != nil {
				return nil, err
			} else {
				mov_.TerceroId = tercero
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

	if movimiento_, err := movimientosArka.GetMovimientoById(id); err != nil {
		return nil, err
	} else {
		movimiento = movimiento_
	}

	if err := json.Unmarshal([]byte(movimiento.Detalle), &detalle); err != nil {
		logs.Error(err)
		eval := " - json.Unmarshal([]byte(movimiento.Detalle), &detalle)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	if db_, cr_, err := parametros.GetParametrosDebitoCredito(); err != nil {
		return nil, err
	} else {
		parametroDebitoId = db_
		parametroCreditoId = cr_
	}

	movs := make([]*models.MovimientoTransaccion, 0)
	for _, mov := range detalle.PreTrAjuste.Movimientos {
		mov_ := new(models.MovimientoTransaccion)
		var cta *models.DetalleCuenta

		if ctaCr_, err := cuentasContables.GetCuentaContable(mov.Cuenta); err != nil {
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

	transaccion.ConsecutivoId = detalle.ConsecutivoId
	transaccion.Movimientos = movs
	transaccion.FechaTransaccion = time.Now()
	transaccion.Activo = true
	transaccion.Etiquetas = ""
	transaccion.Descripcion = ""

	if tr, err := movimientosContables.PostTrContable(transaccion); err != nil {
		return nil, err
	} else {
		detalle.TrContableId = tr.ConsecutivoId
		detalle.PreTrAjuste = nil
		detalle.RazonRechazo = ""
	}

	if sm, err := movimientosArka.GetAllEstadoMovimiento("query=Nombre:" + url.QueryEscape("Ajuste Aprobado")); err != nil {
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

	if movimiento_, err := movimientosArka.PutMovimiento(movimiento, movimiento.Id); err != nil {
		return nil, err
	} else {
		movimiento = movimiento_
	}

	return movimiento, nil
}
