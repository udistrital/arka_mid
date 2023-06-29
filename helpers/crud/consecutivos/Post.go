package consecutivos

import (
	"strings"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	"github.com/udistrital/arka_mid/utils_oas/request"
)

var ConsecutivosCRUD, _ = beego.AppConfig.String("consecutivosService")

// Post post controlador consecutivo del api consecutivos_crud
func Post(consecutivo interface{}) (outputError map[string]interface{}) {

	funcion := "Post"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + ConsecutivosCRUD + "consecutivo"
	response := new(models.RespuestaAPI1Interface)
	if err := request.SendJson(urlcrud, "POST", &response, &consecutivo); err != nil {
		if strings.Contains(err.Error(), "invalid character") {
			logs.Error(err)
			outputError = Post(consecutivo)
			return
		} else {
			logs.Error(urlcrud + ", " + err.Error())
			eval := ` - request.SendJson(urlcrud, "POST", &response, &consecutivo)`
			return errorCtrl.Error(funcion+eval, err, "502")
		}
	}

	if !response.Success {
		err := response.Message
		logs.Error(err)
		eval := ` - request.SendJson(urlcrud, "POST", &response, &consecutivo)`
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	outputError = utilsHelper.FillStruct(response.Data, &consecutivo)
	return

}
