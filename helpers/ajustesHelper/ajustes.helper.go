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
	movimientoscontablesmidHelper "github.com/udistrital/arka_mid/helpers/movimientosContablesMidHelper"
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

func getTipoComprobanteAjustes() string {
	return "N20"
}
