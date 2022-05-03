package consecutivos

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

var ConsecutivosCRUD = "http://" + beego.AppConfig.String("consecutivosService")

// Post post controlador consecutivo del api consecutivos_crud
func Post(consecutivo interface{}) (outputError map[string]interface{}) {

	funcion := "Post"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := ConsecutivosCRUD + "consecutivo"
	response := new(models.RespuestaAPI1Interface)
	if err := request.SendJson(urlcrud, "POST", &response, &consecutivo); err != nil {
		if strings.Contains(err.Error(), "invalid character") {
			logs.Error(err)
			outputError = Post(consecutivo)
			return
		} else {
			logs.Error(urlcrud + ", " + err.Error())
			eval := ` - request.SendJson(urlcrud, "POST", &response, &consecutivo)`
			return errorctrl.Error(funcion+eval, err, "502")
		}
	}

	if !response.Success {
		err := response.Message
		logs.Error(err)
		eval := ` - request.SendJson(urlcrud, "POST", &response, &consecutivo)`
		return errorctrl.Error(funcion+eval, err, "502")
	}

	if err := formatdata.FillStruct(response.Data, &consecutivo); err != nil {
		logs.Error(err)
		eval := " - formatdata.FillStruct(response.Data, &consecutivo)"
		return errorctrl.Error(funcion+eval, err, "500")
	}

	return

}
