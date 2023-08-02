package administrativa

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	"github.com/udistrital/arka_mid/utils_oas/request"
)

var administrativa_amazon, _ = beego.AppConfig.String("administrativaService")

func GetOrdenadores(id int, ordenadores interface{}) (outputError map[string]interface{}) {

	funcion := "GetOrdenadores - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + administrativa_amazon + "ordenadores"
	if id > 0 {
		urlcrud += "/" + strconv.Itoa(id)
	}

	if err := request.GetJson(urlcrud, &ordenadores); err != nil {
		logs.Error(urlcrud + ", " + err.Error())
		eval := "request.GetJson(urlcrud, &ordenadores)"
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}
