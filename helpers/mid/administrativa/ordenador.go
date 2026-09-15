package administrativa

import (
	"context"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	requestV2 "github.com/udistrital/utils_oas/v2/request"
)

var administrativa_amazon, _ = beego.AppConfig.String("administrativaService")

func GetOrdenadores(ctx context.Context, id int, ordenadores interface{}) (outputError map[string]interface{}) {

	funcion := "GetOrdenadores - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := administrativa_amazon + "ordenadores"
	if id > 0 {
		urlcrud += "/" + strconv.Itoa(id)
	}

	_, err := requestV2.GetWithContext(ctx, urlcrud, ordenadores)
	if err != nil {
		logs.Error(urlcrud + ", " + err.Error())
		eval := "requestV2.GetWithContext(ctx, urlcrud, ordenadores)"
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}
