package terceros

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

var tercerosMID = beego.AppConfig.String("tercerosMidService")

func GetTercerosByTipo(tipo string, id int, terceros interface{}) (outputError map[string]interface{}) {

	funcion := "GetTercerosByTipo"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + tercerosMID + "tipo/" + tipo
	if id > 0 {
		urlcrud += "/" + strconv.Itoa(id)
	}

	if err := request.GetJson(urlcrud, &terceros); err != nil {
		logs.Error(urlcrud + ", " + err.Error())
		eval := " - request.GetJson(urlcrud, &terceros)"
		return errorctrl.Error(funcion+eval, err, "502")
	}

	return
}
