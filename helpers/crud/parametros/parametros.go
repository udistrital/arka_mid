package parametros

import (
	"context"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	errorCtrl "github.com/udistrital/utils_oas/v2/errorctrl"
	requestV2 "github.com/udistrital/utils_oas/v2/request"
)

var basePath, _ = beego.AppConfig.String("parametrosService")

// GetAllParametro query controlador parametro del api parametros_crud
func GetAllParametro(ctx context.Context, query string) (parametros []*models.Parametro, outputError map[string]interface{}) {

	funcion := "GetAllParametro"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "parametro?" + query
	response := new(models.RespuestaAPI1Arr)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &response); err != nil {
		eval := " - requestV2.GetWithContext(ctx, urlcrud, &response)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	} else {
		outputError = utilsHelper.FillStruct(response.Data, &parametros)
	}
	return
}

// GetParametroById query controlador parametro/{id} del api parametros_crud
func GetParametroById(ctx context.Context, id int, parametro interface{}) (outputError map[string]interface{}) {

	funcion := "GetAllParametro - "
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "parametro/" + strconv.Itoa(id)
	response := new(models.RespuestaAPI1Interface)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, response); err != nil {
		eval := "requestV2.GetWithContext(ctx, urlcrud, response)"
		return errorCtrl.Error(funcion+eval, err, "502")
	} else if !response.Success {
		eval := "requestV2.GetWithContext(ctx, urlcrud, response)"
		return errorCtrl.Error(funcion+eval, response.Message, response.Status)
	} else {
		outputError = utilsHelper.FillStruct(response.Data, &parametro)
	}
	return
}

// GetAllParametro query controlador parametro del api parametros_crud
func GetAllParametroPeriodo(ctx context.Context, payload string, parametros *[]models.ParametroPeriodo) (outputError map[string]interface{}) {

	funcion := "GetAllParametroPeriodo - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := basePath + "parametro_periodo?" + payload
	response := new(models.RespuestaAPI1Arr)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &response); err != nil {
		eval := "requestV2.GetWithContext(ctx, urlcrud, &response)"
		return errorCtrl.Error(funcion+eval, err, "502")
	} else {
		outputError = utilsHelper.FillStruct(response.Data, &parametros)
	}
	return
}
