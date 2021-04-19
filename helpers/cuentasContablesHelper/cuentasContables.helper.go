package cuentasContablesHelper

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/models"
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

	urlcrud := "http://" + beego.AppConfig.String("cuentasContablesService") + "nodo_cuenta_contable/" + cuentaContableId
	// logs.Debug("urlcrud:", urlcrud)

	var data models.RespuestaAPI2obj
	if resp, err := request.GetJsonTest(urlcrud, &data); err == nil && resp.StatusCode == 200 && data.Code == 200 {
		return data.Body, nil
	} else {
		if err == nil {
			if resp.StatusCode != 200 {
				err = fmt.Errorf("Undesired Status Code: %d", resp.StatusCode)
			} else {
				err = fmt.Errorf("Undesired Status Code: %d - in Body: %d", resp.StatusCode, data.Code)
			}
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
