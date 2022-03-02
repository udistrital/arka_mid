package parametros

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

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
