package administrativa

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	"github.com/udistrital/arka_mid/utils_oas/request"
)

func GetSupervisor(id int, supervisores interface{}) (outputError map[string]interface{}) {
	funcion := "GetSupervisor"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + administrativa_amazon + "supervisor_contrato"
	if id > 0 {
		urlcrud += "/" + strconv.Itoa(id)
	}
	if err := request.GetJson(urlcrud, &supervisores); err != nil {
		logs.Error(urlcrud + ", " + err.Error())
		eval := " - request.GetJson(urlcrud, &supervisores)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

func GetAllDependenciaSIC(payload string, dependencias *[]interface{}) (outputError map[string]interface{}) {
	funcion := "GetAllDependenciaSIC - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + administrativa_amazon + "dependencia_SIC?" + payload
	if err := request.GetJson(urlcrud, &dependencias); err != nil {
		logs.Error(urlcrud + ", " + err.Error())
		eval := "request.GetJson(urlcrud, &supervisores)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}
