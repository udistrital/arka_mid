package ubicacionHelper

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetUbicacion ...
func GetUbicacion(espacioFisicoId int) (espacioFisico []*models.EspacioFisico, outputError map[string]interface{}) {
	if espacioFisicoId != 0 { // (1) error parametro
		if response, err := request.GetJsonTest("http://"+beego.AppConfig.String("oikos2Service")+"espacio_fisico?query=Id:"+strconv.Itoa(int(espacioFisicoId)), &espacioFisico); err == nil { // (2) error servicio caido
			if response.StatusCode == 200 { // (3) error estado de la solicitud
				return espacioFisico, nil
			} else {
				logs.Info("Error (3) estado de la solicitud")
				outputError = map[string]interface{}{"Function": "GetUbicacion:GetUbicacion", "Error": response.Status}
				return nil, outputError
			}
		} else {
			logs.Debug(err)
			logs.Info("Error (2) servicio caido")
			outputError = map[string]interface{}{"Function": "GetUbicacion", "Error": err}
			return nil, outputError
		}
	} else {
		logs.Info("Error (1) Parametro")
		outputError = map[string]interface{}{"Function": "FuncionalidadMidController:GetUbicacion", "Error": "null parameter"}
		return nil, outputError
	}
}
