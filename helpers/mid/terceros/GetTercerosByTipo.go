package terceros

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	"github.com/udistrital/arka_mid/utils_oas/request"
)

var tercerosMID, _ = beego.AppConfig.String("tercerosMidService")

func GetTercerosByTipo(tipo string, id int, terceros interface{}) (outputError map[string]interface{}) {

	funcion := "GetTercerosByTipo"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + tercerosMID + "tipo/" + tipo
	if id > 0 {
		urlcrud += "/" + strconv.Itoa(id)
	}

	if err := request.GetJson(urlcrud, &terceros); err != nil {
		logs.Error(urlcrud + ", " + err.Error())
		eval := " - request.GetJson(urlcrud, &terceros)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}
