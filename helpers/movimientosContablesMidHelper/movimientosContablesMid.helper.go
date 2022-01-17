package movimientosContablesMidHelper

import (
	"strconv"

	"github.com/astaxie/beego"
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
