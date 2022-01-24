package cuentasContablesHelper

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

// GetCuentaContable Consulta controlador nodo_cuenta_contable/{UUID}
func GetCuentaContable(cuentaContableId string) (cuentaContable *models.CuentaContable, outputError map[string]interface{}) {

	funcion := "GetCuentaContable"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("cuentasContablesService") + "nodo_cuenta_contable/" + cuentaContableId
	var data models.RespuestaAPI2obj
	if err := request.GetJson(urlcrud, &data); err != nil || data.Code != 200 {
		if data.Message == "document-no-found" {
			return nil, nil
		}
		eval := " - request.GetJson(urlcrud, &response)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	} else {
		if err := formatdata.FillStruct(data.Body, &cuentaContable); err != nil {
			logs.Error(err)
			eval := " - formatdata.FillStruct(data.Body, &cuentaContable)"
			return nil, errorctrl.Error(funcion+eval, err, "500")
		}
		return cuentaContable, nil
	}
}
