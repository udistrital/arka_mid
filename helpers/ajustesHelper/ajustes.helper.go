package ajustesHelper

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/cuentasContablesHelper"
	"github.com/udistrital/arka_mid/helpers/movimientosArkaHelper"
	"github.com/udistrital/arka_mid/helpers/movimientosContablesMidHelper"
	"github.com/udistrital/arka_mid/helpers/parametrosHelper"
	"github.com/udistrital/arka_mid/helpers/tercerosHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
)

type PreMovAjuste struct {
	Cuenta      string
	Debito      float64
	Credito     float64
	Descripcion string
	TerceroId   int
}

type PreTrAjuste struct {
	Descripcion string
	Movimientos []*PreMovAjuste
}

type FormatoAjuste struct {
	PreTrAjuste  *PreTrAjuste
	Consecutivo  string
	RazonRechazo string
	TrContableId int
}

type DetalleCuenta struct {
	Codigo          string
	Nombre          string
	RequiereTercero bool
}

type DetalleMovimientoContable struct {
	Cuenta      *DetalleCuenta
	Debito      float64
	Credito     float64
	Descripcion string
	TerceroId   *models.Tercero
}

type DetalleAjuste struct {
	Movimiento *models.Movimiento
	TrContable []*DetalleMovimientoContable
}

func PostAjuste(trContable *PreTrAjuste) (movimiento *models.Movimiento, outputError map[string]interface{}) {

	funcion := "PostAjuste"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var query string
	movimiento = new(models.Movimiento)
	detalle := new(FormatoAjuste)

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
		eval := " - jsonData, err := json.Marshal(detalle)"
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
func GetDetalleAjuste(id int) (Ajuste *DetalleAjuste, outputError map[string]interface{}) {

	funcion := "GetDetalleAjuste"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		movimiento         *models.Movimiento
		detalle            *FormatoAjuste
		movimientos        []*PreMovAjuste
		parametroCreditoId int
		parametroDebitoId  int
	)

	Ajuste = new(DetalleAjuste)

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

	if par_, err := parametrosHelper.GetAllParametro("query=CodigoAbreviacion:MCD"); err != nil {
		return nil, err
	} else {
		parametroDebitoId = par_[0].Id
	}

	if par_, err := parametrosHelper.GetAllParametro("query=CodigoAbreviacion:MCC"); err != nil {
		return nil, err
	} else {
		parametroCreditoId = par_[0].Id
	}

	if detalle.PreTrAjuste != nil && detalle.TrContableId == 0 {
		movimientos = detalle.PreTrAjuste.Movimientos
	} else if detalle.PreTrAjuste == nil && detalle.TrContableId > 0 {
		if tr, err := movimientosContablesMidHelper.GetTransaccion(detalle.TrContableId, "consecutivo", true); err != nil {
			return nil, err
		} else {
			for _, mov := range tr.Movimientos {
				mov_ := new(PreMovAjuste)
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

	movs := make([]*DetalleMovimientoContable, 0)
	for _, mov := range movimientos {
		mov_ := new(DetalleMovimientoContable)
		var cta *DetalleCuenta

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
		detalle            *FormatoAjuste
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

	if par_, err := parametrosHelper.GetAllParametro("query=CodigoAbreviacion:MCD"); err != nil {
		return nil, err
	} else {
		parametroDebitoId = par_[0].Id
	}

	if par_, err := parametrosHelper.GetAllParametro("query=CodigoAbreviacion:MCC"); err != nil {
		return nil, err
	} else {
		parametroCreditoId = par_[0].Id
	}

	movs := make([]*models.MovimientoTransaccion, 0)
	for _, mov := range detalle.PreTrAjuste.Movimientos {
		mov_ := new(models.MovimientoTransaccion)
		var cta *DetalleCuenta

		if ctaCr_, err := cuentasContablesHelper.GetCuentaContable(mov.Cuenta); err != nil {
			return nil, err
		} else {
			if err := formatdata.FillStruct(ctaCr_, &cta); err != nil {
				logs.Error(err)
				eval := " - formatdata.FillStruct(ctaCr_, &ctaCr)"
				return nil, errorctrl.Error(funcion+eval, err, "500")
			}
			mov_.CuentaId = cta.Codigo
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

	if resp, err := cuentasContablesHelper.PostTrContable(transaccion); err != nil || !resp.Success {
		if err == nil {
			eval := " - cuentasContablesHelper.PostTrContable(transaccion)"
			return nil, errorctrl.Error(funcion+eval, resp.Data, resp.Status)
		}
		return nil, err
	} else {
		detalle.TrContableId = transaccion.ConsecutivoId
		detalle.PreTrAjuste = nil
	}

	if sm, err := movimientosArkaHelper.GetAllEstadoMovimiento(url.QueryEscape("Ajuste Aprobado")); err != nil {
		return nil, err
	} else {
		movimiento.EstadoMovimientoId = sm[0]
	}

	if jsonData, err := json.Marshal(detalle); err != nil {
		logs.Error(err)
		eval := " - jsonData, err := json.Marshal(detalle)"
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

func getTipoComprobanteAjustes() string {
	return "N20"
}
