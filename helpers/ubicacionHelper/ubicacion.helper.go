package ubicacionHelper

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
	// "github.com/udistrital/arka_mid/helpers/utilsHelper"
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

func GetAsignacionSedeDependencia(Id string) (Relacion map[string]interface{}, err error) {

	var ubicacion []map[string]interface{}
	relacion := make(map[string]interface{}, 0)

	url2 := "http://"+beego.AppConfig.String("oikos2Service")+"asignacion_espacio_fisico_dependencia?query=Id:" + Id
			
	if _, err := request.GetJsonTest(url2, &ubicacion); err == nil { // (2) error servicio caido
		
		if keys := len(ubicacion[0]); keys != 0 {

			return ubicacion[0], nil

		} else {
			return relacion, nil
		}
	} else {
		panic(err.Error())
		return nil, err
	}
	return ubicacion[0], nil
}
