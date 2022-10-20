package administrativa

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

var administrativa_amazon = beego.AppConfig.String("administrativaService")

func GetOrdenadores(id int, ordenadores interface{}) (outputError map[string]interface{}) {

	funcion := "GetTercerosByTipo"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + administrativa_amazon + "ordenadores"
	if id > 0 {
		urlcrud += "/" + strconv.Itoa(id)
	}

	if err := request.GetJson(urlcrud, &ordenadores); err != nil {
		logs.Error(urlcrud + ", " + err.Error())
		eval := " - request.GetJson(urlcrud, &ordenadores)"
		return errorctrl.Error(funcion+eval, err, "502")
	}

	return
}
