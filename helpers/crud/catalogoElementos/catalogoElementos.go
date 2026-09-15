package catalogoElementos

import (
	"context"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	requestV2 "github.com/udistrital/utils_oas/v2/request"
)

var basePath, _ = beego.AppConfig.String("catalogoElementosService")

// GetAllCuentasSubgrupo query controlador cuentas_subgrupo del api catalogo_elementos_crud
func GetAllCuentasSubgrupo(ctx context.Context, query string) (elementos []*models.CuentasSubgrupo, outputError map[string]interface{}) {

	funcion := "GetAllCuentasSubgrupo"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "cuentas_subgrupo?" + query
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &elementos); err != nil {
		logs.Error(err)
		eval := " - requestV2.GetWithContext(ctx, urlcrud, &elementos)"
		return nil, errorCtrl.Error(funcion+eval, err, "500")
	}

	return elementos, nil
}

// GetTrCuentasSubgrupo query controlador cuentas_subgrupo del api catalogo_elementos_crud
func GetTrCuentasSubgrupo(ctx context.Context, id, movimientoId int, cuentas *[]models.CuentasSubgrupo) (outputError map[string]interface{}) {

	funcion := "GetTrCuentasSubgrupo - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := basePath + "tr_cuentas_subgrupo/" + strconv.Itoa(id)
	if movimientoId > 0 {
		urlcrud += "?movimientoId=" + strconv.Itoa(movimientoId)
	}

	if _, err := requestV2.GetWithContext(ctx, urlcrud, cuentas); err != nil {
		logs.Error(urlcrud, err)
		eval := "requestV2.GetWithContext(ctx, urlcrud, cuentas)"
		return errorCtrl.Error(funcion+eval, err, "500")
	}

	return
}

// GetAllDetalleSubgrupo query controlador cuentas_subgrupo del api catalogo_elementos_crud
func GetAllDetalleSubgrupo(ctx context.Context, query string) (detalle []*models.DetalleSubgrupo, outputError map[string]interface{}) {

	funcion := "GetAllDetalleSubgrupo"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "detalle_subgrupo?" + query
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &detalle); err != nil {
		logs.Error(err)
		eval := " - requestV2.GetWithContext(ctx, urlcrud, &detalle)"
		return nil, errorCtrl.Error(funcion+eval, err, "500")
	}

	return detalle, nil
}

// GetAllTipoBien query controlador tipo_bien del api catalogo_elementos_crud
func GetAllTipoBien(ctx context.Context, query string, tiposBien *[]models.TipoBien) (outputError map[string]interface{}) {

	funcion := "GetAllTipoBien - "
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "tipo_bien?" + query
	if _, err := requestV2.GetWithContext(ctx, urlcrud, tiposBien); err != nil {
		logs.Error(urlcrud, err)
		eval := "requestV2.GetWithContext(ctx, urlcrud, tiposBien)"
		return errorCtrl.Error(funcion+eval, err, "500")
	}

	return
}

// GetTipoBienById query controlador tipo_bien/{id} del api catalogo_elementos_crud
func GetTipoBienById(ctx context.Context, id int, tipoBien *models.TipoBien) (outputError map[string]interface{}) {

	funcion := "GetTipoBienById - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := basePath + "tipo_bien/" + strconv.Itoa(id)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &tipoBien); err != nil {
		logs.Error(err)
		eval := "requestV2.GetWithContext(ctx, urlcrud, &tipoBien)"
		return errorCtrl.Error(funcion+eval, err, "500")
	}

	return
}

// GetSubgrupoById Consulta controlador subgrupo/{id} del api catalogo_elementos_crud
func GetSubgrupoById(ctx context.Context, id int) (subgrupo models.Subgrupo, outputError map[string]interface{}) {

	funcion := "GetSubgrupoById - "
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "subgrupo/" + strconv.Itoa(id)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &subgrupo); err != nil {
		logs.Error(err)
		eval := "requestV2.GetWithContext(ctx, urlcrud, &subgrupo)"
		outputError = errorCtrl.Error(funcion+eval, err, "500")
	}

	return
}

// GetAllElemento Consulta controlador elemento del api catalogo_elementos_crud
func GetAllElemento(ctx context.Context, payload string, elementos *[]models.ElementoCatalogo) (outputError map[string]interface{}) {

	funcion := "GetAllElemento - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := basePath + "elemento?" + payload
	if _, err := requestV2.GetWithContext(ctx, urlcrud, elementos); err != nil {
		logs.Error(err)
		eval := "requestV2.GetWithContext(ctx, urlcrud, elementos)"
		return errorCtrl.Error(funcion+eval, err, "500")
	}

	return
}
