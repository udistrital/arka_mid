package terceros

import (
	"context"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	errorCtrl "github.com/udistrital/utils_oas/v2/errorctrl"
	requestV2 "github.com/udistrital/utils_oas/v2/request"
)

var tercerosMID, _ = beego.AppConfig.String("tercerosMidService")

func GetTercerosByTipo(ctx context.Context, tipo string, id int, terceros interface{}) (outputError map[string]interface{}) {

	funcion := "GetTercerosByTipo"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := tercerosMID + "tipo/" + tipo
	if id > 0 {
		urlcrud += "/" + strconv.Itoa(id)
	}

	if _, err := requestV2.GetWithContext(ctx, urlcrud, terceros); err != nil {
		logs.Error(urlcrud + ", " + err.Error())
		eval := " - requestV2.GetWithContext(ctx, urlcrud, terceros)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}
