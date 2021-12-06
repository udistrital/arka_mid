package catalogoElementosHelper

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

func GetAllCuentasSubgrupo(query string) (elementos []*models.CuentaSubgrupo, outputError map[string]interface{}) {

	funcion := "GetAllCuentasSubgrupo"
	defer errorctrl.ErrorControlFunction(funcion, "500")

	urlcrud := "http://" + beego.AppConfig.String("catalogoElementosService") + "cuentas_subgrupo?" + query
	if err := request.GetJson(urlcrud, &elementos); err != nil {
		logs.Error(err)
		funcion += " - request.GetJson(urlcrud, &elementos)"
		return nil, errorctrl.Error(funcion, err, "500")
	}

	return elementos, nil
}
