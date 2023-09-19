package parametros

import (
	"strconv"

	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	"github.com/udistrital/arka_mid/utils_oas/request"
)

var basePath, _ = beego.AppConfig.String("parametrosService")

// GetAllParametro query controlador parametro del api parametros_crud
func GetAllParametro(query string) (parametros []*models.Parametro, outputError map[string]interface{}) {

	funcion := "GetAllParametro"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "parametro?" + query
	response := new(models.RespuestaAPI1Arr)
	if err := request.GetJson(urlcrud, &response); err != nil {
		eval := " - request.GetJson(urlcrud, &response)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	} else {
		outputError = utilsHelper.FillStruct(response.Data, &parametros)
	}
	return
}

// GetParametroById query controlador parametro/{id} del api parametros_crud
func GetParametroById(id int, parametro interface{}) (outputError map[string]interface{}) {

	funcion := "GetAllParametro - "
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "parametro/" + strconv.Itoa(id)
	response := new(models.RespuestaAPI1Interface)
	if err := request.GetJson(urlcrud, &response); err != nil {
		eval := "request.GetJson(urlcrud, &response)"
		return errorCtrl.Error(funcion+eval, err, "502")
	} else if !response.Success {
		eval := "request.GetJson(urlcrud, &response)"
		return errorCtrl.Error(funcion+eval, response.Message, response.Status)
	} else {
		outputError = utilsHelper.FillStruct(response.Data, &parametro)
	}
	return
}

// GetAllParametro query controlador parametro del api parametros_crud
func GetAllParametroPeriodo(payload string, parametros *[]models.ParametroPeriodo) (outputError map[string]interface{}) {

	funcion := "GetAllParametroPeriodo - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "parametro_periodo?" + payload
	response := new(models.RespuestaAPI1Arr)
	if err := request.GetJson(urlcrud, &response); err != nil {
		eval := "request.GetJson(urlcrud, &response)"
		return errorCtrl.Error(funcion+eval, err, "502")
	} else {
		outputError = utilsHelper.FillStruct(response.Data, &parametros)
	}
	return
}
