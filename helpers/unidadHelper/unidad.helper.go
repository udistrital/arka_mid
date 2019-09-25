package unidadHelper

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetUnidad ...
func GetUnidad(unidadId int) (unidad []*models.Unidad, outputError map[string]interface{}) {
	if unidadId != 0 { // (1) error parametro
		var unidadAux *models.Unidad

		if response, err := request.GetJsonTest("http://"+beego.AppConfig.String("unidadService")+strconv.Itoa(unidadId), &unidadAux); err == nil { // (2) error servicio caido
			if response.StatusCode == 200 { // (3) error estado de la solicitud
				unidad = append(unidad, unidadAux)
				return unidad, nil
			} else {
				logs.Info("Error (3) estado de la solicitud")
				outputError = map[string]interface{}{"Function": "GetUnidad:GetUnidad", "Error": response.Status}
				return nil, outputError
			}
		} else {
			logs.Info("Error (2) servicio caido")
			outputError = map[string]interface{}{"Function": "GetUnidad", "Error": err}
			return nil, outputError
		}
	} else {
		logs.Info("Error (1) Parametro")
		outputError = map[string]interface{}{"Function": "FuncionalidadMidController:GetUnidad", "Error": "null parameter"}
		return nil, outputError
	}
}
