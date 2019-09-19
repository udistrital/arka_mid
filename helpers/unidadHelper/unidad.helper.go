package unidadHelper

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
)

// GetUnidad ...
func GetUnidad(unidadId int) (unidad interface{}, outputError map[string]interface{}) {
	if unidadId != 0 { // (1) error parametro
		if response, err := request.GetJsonTest("http://"+beego.AppConfig.String("unidadService")+strconv.Itoa(unidadId)+
			",SoporteActaId.ActaRecibidoId.Activo:True&limit=-1", &unidad); err == nil { // (2) error servicio caido
			if response.StatusCode == 200 { // (3) error estado de la solicitud
				return unidad, nil
			} else {
				logs.Info("Error (3) estado de la solicitud")
				outputError = map[string]interface{}{"Function": "GetAllActasRecibido:GetAllActasRecibido", "Error": response.Status}
				return outputError, nil
			}
		} else {
			logs.Info("Error (2) servicio caido")
			outputError = map[string]interface{}{"Function": "GetAllActasRecibido", "Error": err}
			return outputError, nil
		}
	} else {
		logs.Info("Error (1) Parametro")
		outputError = map[string]interface{}{"Function": "FuncionalidadMidController:getUserAgora", "Error": "null parameter"}
		return nil, outputError
	}
}
