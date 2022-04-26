package catalogoElementos

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

// GetAllCuentasSubgrupo query controlador cuentas_subgrupo del api catalogo_elementos_crud
func GetAllCuentasSubgrupo(query string) (elementos []*models.CuentaSubgrupo, outputError map[string]interface{}) {

	funcion := "GetAllCuentasSubgrupo"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("catalogoElementosService") + "cuentas_subgrupo?" + query
	if err := request.GetJson(urlcrud, &elementos); err != nil {
		logs.Error(err)
		eval := " - request.GetJson(urlcrud, &elementos)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	return elementos, nil
}

// GetTrCuentasSubgrupo query controlador cuentas_subgrupo del api catalogo_elementos_crud
func GetTrCuentasSubgrupo(id int) (cuentas []*models.CuentasSubgrupo, outputError map[string]interface{}) {

	funcion := "GetTrCuentasSubgrupo"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("catalogoElementosService") + "tr_cuentas_subgrupo/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &cuentas); err != nil {
		logs.Error(err)
		eval := " - request.GetJson(urlcrud, &cuentas)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	return cuentas, nil
}

// GetAllDetalleSubgrupo query controlador cuentas_subgrupo del api catalogo_elementos_crud
func GetAllDetalleSubgrupo(query string) (detalle []*models.DetalleSubgrupo, outputError map[string]interface{}) {

	funcion := "GetAllDetalleSubgrupo"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("catalogoElementosService") + "detalle_subgrupo?" + query
	if err := request.GetJson(urlcrud, &detalle); err != nil {
		logs.Error(err)
		eval := " - request.GetJson(urlcrud, &detalle)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	return detalle, nil
}

// GetSubgrupoById Consulta controlador subgrupo/{id} del api catalogo_elementos_crud
func GetSubgrupoById(id int) (subgrupo *models.Subgrupo, outputError map[string]interface{}) {

	funcion := "GetSubgrupoById"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("catalogoElementosService") + "subgrupo/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &subgrupo); err != nil {
		logs.Error(err)
		eval := " - request.GetJson(urlcrud, &subgrupo)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	return subgrupo, nil
}
