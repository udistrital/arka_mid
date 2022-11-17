package parametros

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

var basePath = "http://" + beego.AppConfig.String("parametrosService")

// GetAllParametro query controlador parametro del api parametros_crud
func GetAllParametro(query string) (parametros []*models.Parametro, outputError map[string]interface{}) {

	funcion := "GetAllParametro"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("parametrosService") + "parametro?" + query
	response := new(models.RespuestaAPI1Arr)
	if err := request.GetJson(urlcrud, &response); err != nil {
		eval := " - request.GetJson(urlcrud, &response)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	} else {
		if err := formatdata.FillStruct(response.Data, &parametros); err != nil {
			logs.Error(err)
			eval := " - formatdata.FillStruct(response.Data, &parametros)"
			return nil, errorctrl.Error(funcion+eval, err, "500")
		}
	}
	return parametros, nil
}

// GetParametroById query controlador parametro/{id} del api parametros_crud
func GetParametroById(id int, parametro interface{}) (outputError map[string]interface{}) {

	funcion := "GetAllParametro - "
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "parametro/" + strconv.Itoa(id)
	response := new(models.RespuestaAPI1Interface)
	if err := request.GetJson(urlcrud, &response); err != nil {
		eval := "request.GetJson(urlcrud, &response)"
		return errorctrl.Error(funcion+eval, err, "502")
	} else if !response.Success {
		eval := "request.GetJson(urlcrud, &response)"
		return errorctrl.Error(funcion+eval, response.Message, response.Status)
	} else {
		if err := formatdata.FillStruct(response.Data, &parametro); err != nil {
			logs.Error(err)
			eval := "formatdata.FillStruct(response.Data, &parametro)"
			return errorctrl.Error(funcion+eval, err, "500")
		}
	}
	return
}

// GetAllParametro query controlador parametro del api parametros_crud
func GetAllParametroPeriodo(payload string, parametros *[]models.ParametroPeriodo) (outputError map[string]interface{}) {

	funcion := "GetAllParametroPeriodo - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("parametrosService") + "parametro_periodo?" + payload
	response := new(models.RespuestaAPI1Arr)
	if err := request.GetJson(urlcrud, &response); err != nil {
		eval := "request.GetJson(urlcrud, &response)"
		return errorctrl.Error(funcion+eval, err, "502")
	} else {
		if err := formatdata.FillStruct(response.Data, &parametros); err != nil {
			logs.Error(err)
			eval := "formatdata.FillStruct(response.Data, &parametros)"
			return errorctrl.Error(funcion+eval, err, "500")
		}
	}

	return
}
