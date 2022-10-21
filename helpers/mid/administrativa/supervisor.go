package administrativa

import (
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

func GetSupervisor(id int, supervisores interface{}) (outputError map[string]interface{}) {
	funcion := "GetSupervisor"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + administrativa_amazon + "supervisor_contrato"
	if id > 0 {
		urlcrud += "/" + strconv.Itoa(id)
	}
	if err := request.GetJson(urlcrud, &supervisores); err != nil {
		logs.Error(urlcrud + ", " + err.Error())
		eval := " - request.GetJson(urlcrud, &supervisores)"
		return errorctrl.Error(funcion+eval, err, "502")
	}

	return
}
