package administrativa

import (
	"context"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	errorCtrl "github.com/udistrital/utils_oas/v2/errorctrl"
	requestV2 "github.com/udistrital/utils_oas/v2/request"
)

func GetSupervisor(ctx context.Context, id int, supervisores interface{}) (outputError map[string]interface{}) {
	funcion := "GetSupervisor"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := administrativa_amazon + "supervisor_contrato"
	if id > 0 {
		urlcrud += "/" + strconv.Itoa(id)
	}
	if _, err := requestV2.GetWithContext(ctx, urlcrud, supervisores); err != nil {
		logs.Error(urlcrud + ", " + err.Error())
		eval := " - requestV2.GetWithContext(ctx, urlcrud, supervisores)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

func GetSupervisorByQuery(ctx context.Context, payload string, supervisores interface{}) (outputError map[string]interface{}) {
	funcion := "GetSupervisorByQuery"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := administrativa_amazon + "supervisor_contrato"
	if payload != "" {
		urlcrud += "?" + payload
	}
	if _, err := requestV2.GetWithContext(ctx, urlcrud, supervisores); err != nil {
		logs.Error(urlcrud + ", " + err.Error())
		eval := " - requestV2.GetWithContext(ctx, urlcrud, supervisores)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

func GetAllDependenciaSIC(ctx context.Context, payload string, dependencias *[]interface{}) (outputError map[string]interface{}) {
	funcion := "GetAllDependenciaSIC - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := administrativa_amazon + "dependencia_SIC?" + payload
	if _, err := requestV2.GetWithContext(ctx, urlcrud, dependencias); err != nil {
		logs.Error(urlcrud + ", " + err.Error())
		eval := "requestV2.GetWithContext(ctx, urlcrud, dependencias)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}
