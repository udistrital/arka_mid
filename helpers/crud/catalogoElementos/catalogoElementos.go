package catalogoElementos

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	"github.com/udistrital/arka_mid/utils_oas/request"
)

var basePath, _ = beego.AppConfig.String("catalogoElementosService")

// GetAllCuentasSubgrupo query controlador cuentas_subgrupo del api catalogo_elementos_crud
func GetAllCuentasSubgrupo(query string) (elementos []*models.CuentasSubgrupo, outputError map[string]interface{}) {

	funcion := "GetAllCuentasSubgrupo"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "cuentas_subgrupo?" + query
	if err := request.GetJson(urlcrud, &elementos); err != nil {
		logs.Error(err)
		eval := " - request.GetJson(urlcrud, &elementos)"
		return nil, errorCtrl.Error(funcion+eval, err, "500")
	}

	return elementos, nil
}

// GetTrCuentasSubgrupo query controlador cuentas_subgrupo del api catalogo_elementos_crud
func GetTrCuentasSubgrupo(id, movimientoId int, cuentas *[]models.CuentasSubgrupo) (outputError map[string]interface{}) {

	funcion := "GetTrCuentasSubgrupo - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "tr_cuentas_subgrupo/" + strconv.Itoa(id)
	if movimientoId > 0 {
		urlcrud += "?movimientoId=" + strconv.Itoa(movimientoId)
	}

	if err := request.GetJson(urlcrud, &cuentas); err != nil {
		logs.Error(urlcrud, err)
		eval := "request.GetJson(urlcrud, &cuentas)"
		return errorCtrl.Error(funcion+eval, err, "500")
	}

	return
}

// GetAllDetalleSubgrupo query controlador cuentas_subgrupo del api catalogo_elementos_crud
func GetAllDetalleSubgrupo(query string) (detalle []*models.DetalleSubgrupo, outputError map[string]interface{}) {

	funcion := "GetAllDetalleSubgrupo"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "detalle_subgrupo?" + query
	if err := request.GetJson(urlcrud, &detalle); err != nil {
		logs.Error(err)
		eval := " - request.GetJson(urlcrud, &detalle)"
		return nil, errorCtrl.Error(funcion+eval, err, "500")
	}

	return detalle, nil
}

// GetAllTipoBien query controlador tipo_bien del api catalogo_elementos_crud
func GetAllTipoBien(query string, tiposBien *[]models.TipoBien) (outputError map[string]interface{}) {

	funcion := "GetAllTipoBien - "
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "tipo_bien?" + query
	if err := request.GetJson(urlcrud, &tiposBien); err != nil {
		logs.Error(urlcrud, err)
		eval := "request.GetJson(urlcrud, &detalle)"
		return errorCtrl.Error(funcion+eval, err, "500")
	}

	return
}

// GetTipoBienById query controlador tipo_bien/{id} del api catalogo_elementos_crud
func GetTipoBienById(id int, tipoBien *models.TipoBien) (outputError map[string]interface{}) {

	funcion := "GetTipoBienById - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "tipo_bien/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &tipoBien); err != nil {
		logs.Error(err)
		eval := "request.GetJson(urlcrud, &tipoBien)"
		return errorCtrl.Error(funcion+eval, err, "500")
	}

	return
}

// GetSubgrupoById Consulta controlador subgrupo/{id} del api catalogo_elementos_crud
func GetSubgrupoById(id int) (subgrupo models.Subgrupo, outputError map[string]interface{}) {

	funcion := "GetSubgrupoById - "
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "subgrupo/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &subgrupo); err != nil {
		logs.Error(err)
		eval := "request.GetJson(urlcrud, &subgrupo)"
		outputError = errorCtrl.Error(funcion+eval, err, "500")
	}

	return
}

// GetAllElemento Consulta controlador elemento del api catalogo_elementos_crud
func GetAllElemento(payload string, elementos *[]models.ElementoCatalogo) (outputError map[string]interface{}) {

	funcion := "GetAllElemento - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "elemento?" + payload
	if err := request.GetJson(urlcrud, &elementos); err != nil {
		logs.Error(err)
		eval := "request.GetJson(urlcrud, &elementos)"
		return errorCtrl.Error(funcion+eval, err, "500")
	}

	return
}
