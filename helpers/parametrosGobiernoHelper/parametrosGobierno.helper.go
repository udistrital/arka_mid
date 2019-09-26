package parametrosGobiernoHelper

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetIva ...
func GetIva(ivaId int) (iva []*models.ParametrosGobierno, outputError map[string]interface{}) {
	if ivaId != 0 { // (1) error parametro

		if response, err := request.GetJsonTest("http://"+beego.AppConfig.String("parametrosGobiernoService")+"vigencia_impuesto?query=Id:"+strconv.Itoa(ivaId), &iva); err == nil { // (2) error servicio caido
			if response.StatusCode == 200 { // (3) error estado de la solicitud
				return iva, nil
			} else {
				logs.Info("Error (3) estado de la solicitud")
				outputError = map[string]interface{}{"Function": "GetIva:GetIva", "Error": response.Status}
				return nil, outputError
			}
		} else {
			logs.Info("Error (2) servicio caido")
			outputError = map[string]interface{}{"Function": "GetIva", "Error": err}
			return nil, outputError
		}
	} else {
		logs.Info("Error (1) Parametro")
		outputError = map[string]interface{}{"Function": "FuncionalidadMidController:GetIva", "Error": "null parameter"}
		return nil, outputError
	}
}
