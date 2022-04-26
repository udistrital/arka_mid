package cuentasContables

import (
	"strconv"

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

// GetTipoComprobante Consulta controlador nodo_cuenta_contable/{UUID}
func GetTipoComprobante(tipoDocumento string) (tipoComprobante *models.TipoComprobante, outputError map[string]interface{}) {

	funcion := "GetCuentaContable"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var data models.RespuestaAPI2arr
	urlcrud := "http://" + beego.AppConfig.String("cuentasContablesService") + "tipo_comprobante"
	if err := request.GetJson(urlcrud, &data); err != nil || data.Code != 200 {
		eval := " - request.GetJson(urlcrud, &response)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}

	var tiposComprobante []*models.TipoComprobante
	if err := formatdata.FillStruct(data.Body, &tiposComprobante); err != nil {
		logs.Error(err)
		eval := " - formatdata.FillStruct(data.Body, &tiposComprobante)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	for _, tipoComprobante := range tiposComprobante {
		if tipoComprobante.TipoDocumento == tipoDocumento {
			return tipoComprobante, nil
		}

	}

	return nil, nil
}

// GetComprobante Retorna el id para un tipo de comprobante dato
func GetComprobante(tipoDocumento string, id *string) (outputError map[string]interface{}) {

	funcion := "GetComprobante"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		data         models.RespuestaAPI2arr
		comprobantes []*models.Comprobante
	)

	urlcrud := "http://" + beego.AppConfig.String("cuentasContablesService") + "comprobante"
	if err := request.GetJson(urlcrud, &data); err != nil || data.Code != 200 {
		eval := " - request.GetJson(urlcrud, &data)"
		return errorctrl.Error(funcion+eval, err, "502")
	}

	if err := formatdata.FillStruct(data.Body, &comprobantes); err != nil {
		logs.Error(err)
		eval := " - formatdata.FillStruct(data.Body, &comprobantes)"
		return errorctrl.Error(funcion+eval, err, "500")
	}

	for _, comprobante := range comprobantes {
		if comprobante.TipoComprobante.TipoDocumento+strconv.Itoa(comprobante.Numero) == tipoDocumento {
			*id = comprobante.Id
			return
		}
	}

	return

}
