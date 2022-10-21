package movimientosContables

import (
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

// GetTransaccion query controlador transaccion_movimientos del api movimientos_contables_mid
func GetTransaccion(id int, criteria string, detail bool) (transaccion *models.TransaccionMovimientos, outputError map[string]interface{}) {

	funcion := "GetTransaccion"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("movimientosContablesmidService") + "transaccion_movimientos/" + criteria + "/" + strconv.Itoa(id)
	if detail {
		urlcrud += "?detailed=true"
	}
	if err := request.GetJson(urlcrud, &transaccion); err != nil {
		eval := " - request.GetJson(urlcrud, &response)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}

	return transaccion, nil
}

// PostTrContable post controlador transaccion_movimientos/transaccion_movimientos/ del api movimientos_contables_mid
func PostTrContable(tr *models.TransaccionMovimientos) (tr_ *models.TransaccionMovimientos, outputError map[string]interface{}) {

	funcion := "PostTrContable"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error", "500")

	var resp map[string]interface{}
	urlcrud := "http://" + beego.AppConfig.String("movimientosContablesmidService") + "transaccion_movimientos"
	if err := request.SendJson(urlcrud, "POST", &resp, &tr); err != nil {
		eval := ` - request.SendJson(urlcrud, "POST", &resp, &tr)`
		return nil, errorctrl.Error(funcion+eval, err, "502")
	} else if !resp["Success"].(bool) {
		if strings.Contains(resp["Data"].(string), "invalid character") {
			logs.Error(resp["Data"])
			tr_, outputError = PostTrContable(tr)
		} else {
			logs.Info(resp["Data"])
			eval := ` - request.SendJson(urlcrud, "POST", &resp, &tr)`
			return nil, errorctrl.Error(funcion+eval, resp["Data"].(map[string]interface{})["err"], "502")
		}
	}

	return tr, nil
}
