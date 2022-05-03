package configuracion

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	// "github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	modelsConfiguracion "github.com/udistrital/configuracion_api/models"
	e "github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

func GetParametros(query utilsHelper.Query, parametros *[]modelsConfiguracion.Parametro) (outputError map[string]interface{}) {
	const funcion = "GetParametros - "
	defer e.ErrorControlFunction(funcion+"unhandled error!", fmt.Sprint(http.StatusInternalServerError))

	qString := query.Encode()
	urlParametros := "http://" + beego.AppConfig.String("configuracionCrud") + "parametro?" + qString
	// logs.Debug(urlParametros)
	var respuestaApi interface{}
	if resp, err := request.GetJsonTest(urlParametros, &respuestaApi); err != nil || resp.StatusCode != http.StatusOK {
		const contexto = "request.GetJsonTest(urlParametros, &respuestaApi)"
		if err == nil {
			err = fmt.Errorf("undesired status code: %d", resp.StatusCode)
		}
		outputError = e.Error(funcion+contexto, err, fmt.Sprint(http.StatusBadGateway))
		return
	}
	// formatdata.JsonPrint(respuestaApi)
	if err := formatdata.FillStruct(respuestaApi, &parametros); err != nil {
		const contexto = "formatdata.FillStruct(respuestaApi, &parametros)"
		outputError = e.Error(funcion+contexto, err, fmt.Sprint(http.StatusInternalServerError))
	}
	return
}
