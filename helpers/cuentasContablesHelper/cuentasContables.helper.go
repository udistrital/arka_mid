package cuentasContablesHelper

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/utils_oas/request"
)

// GetCuentaContable ...
func GetCuentaContable(cuentaContableId string) (cuentaContable map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetCuentaContable - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var (
		urlcrud string
	)
	urlcrud = "http://" + beego.AppConfig.String("cuentasContablesService") + "cuenta_contable/" + cuentaContableId
	logs.Debug("urlcrud:", urlcrud)

	if response, err := request.GetJsonTest(urlcrud, &cuentaContable); err == nil && response.StatusCode == 200 {
		return cuentaContable, nil
	} else {
		if err == nil {
			err = fmt.Errorf("Undesired Status Code: %d", response.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetCuentaContable - request.GetJsonTest(urlcrud, &cuentaContable)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}
