package ubicacionHelper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/utilsHelper"
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

func GetAsignacionSedeDependencia(Id string) (Relacion map[string]interface{}, err error) {

	var ubicacion []map[string]interface{}
	relacion := make(map[string]interface{}, 0)

	url2 := "http://" + beego.AppConfig.String("oikos2Service") + "asignacion_espacio_fisico_dependencia?query=Id:" + Id

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

func GetSedeDependenciaUbicacion(Id string) (Sede map[string]interface{}, Dependencia map[string]interface{}, Ubicacion map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/GetSedeDependenciaUbicacion",
				"err":     err,
				"status":  "404",
			}
			panic(outputError)
		}
	}()

	var Ubicacion_ []map[string]interface{}

	url2 := "http://" + beego.AppConfig.String("oikos2Service") + "asignacion_espacio_fisico_dependencia?query=Id:" + Id

	if _, err := request.GetJsonTest(url2, &Ubicacion_); err == nil { // (2) error servicio caido

		if data, err := utilsHelper.ConvertirInterfaceMap(Ubicacion_[0]["DependenciaId"]); err == nil {
			Dependencia = data
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "/GetSedeDependenciaUbicacion",
				"err":     err,
				"status":  "500",
			}
			return nil, nil, nil, outputError
		}

		if data, err := utilsHelper.ConvertirInterfaceMap(Ubicacion_[0]["EspacioFisicoId"]); err == nil {
			Ubicacion = data

			str2 := fmt.Sprintf("%v", data["CodigoAbreviacion"])
			z := strings.Split(str2, "")
			var sede []map[string]interface{}
			urlcrud4 := "http://" + beego.AppConfig.String("oikos2Service") + "espacio_fisico?query=CodigoAbreviacion:" + z[0] + z[1] + z[2] + z[3]

			if _, err := request.GetJsonTest(urlcrud4, &sede); err == nil {
				Sede = sede[0]
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "/GetSedeDependenciaUbicacion",
					"err":     err,
					"status":  "404",
				}
				return nil, nil, nil, outputError
			}
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "/GetSedeDependenciaUbicacion",
				"err":     err,
				"status":  "500",
			}
			return nil, nil, nil, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/GetSedeDependenciaUbicacion",
			"err":     err,
			"status":  "404",
		}
		return nil, nil, nil, outputError
	}
	return Sede, Dependencia, Ubicacion, nil

}
