package cuentasContablesHelper

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetCuentaContable ...
func GetCuentaContable(cuentaContableId int) (cuentaContable *models.CuentaContable, outputError map[string]interface{}) {
	var (
		urlcrud string
	)

	urlcrud = "http://" + beego.AppConfig.String("cuentasContablesService") + "cuenta_contable/" + strconv.Itoa(int(cuentaContableId))

	if response, err := request.GetJsonTest(urlcrud, &cuentaContable); err == nil { // (2) error servicio caido
		if response.StatusCode == 200 { // (3) error estado de la solicitud
			return cuentaContable, nil
		} else {
			logs.Info("Error (3) estado de la solicitud")
			outputError = map[string]interface{}{"Function": "GetCuentasContablesGrupo:GetCuentasContablesGrupo", "Error": response.Status}
			return nil, outputError
		}
	} else {
		logs.Info("Error (2) servicio caido")
		outputError = map[string]interface{}{"Function": "GetCuentasContablesGrupo", "Error": err}
		return nil, outputError
	}
}
