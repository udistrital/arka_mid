package parametrosGobiernoHelper

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetIva ...
func GetIva(ivaId int) (iva []*models.ParametrosGobierno, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/GetIva - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	if ivaId > 0 { // (1) error parametro

		urlParametroIVA := "http://" + beego.AppConfig.String("parametrosGobiernoService") + "vigencia_impuesto?query=Id:" + strconv.Itoa(ivaId)
		if response, err := request.GetJsonTest(urlParametroIVA, &iva); err == nil { // (2) error servicio caido
			if response.StatusCode == 200 { // (3) error estado de la solicitud
				return iva, nil
			} else {
				err := fmt.Errorf("Undesired Status: %s", response.Status)
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "GetIva - request.GetJsonTest(urlParametroIVA, &iva)",
					"err":     err,
					"status":  "500",
				}
				return nil, outputError
			}
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetIva - request.GetJsonTest(urlParametroIVA, &iva)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
	} else {
		err := fmt.Errorf("ivaId MUST be greater than 0")
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "FuncionalidadMidController:GetIva",
			"err":     err,
			"status":  "400",
		}
		return nil, outputError
	}
}
