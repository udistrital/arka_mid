package administrativa

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetUnidad ...
func GetUnidad(unidadId int) (unidad []*models.Unidad, outputError map[string]interface{}) {
	if unidadId > 0 { // (1) error parametro

		defer func() {
			if err := recover(); err != nil {
				outputError = map[string]interface{}{
					"funcion": "/GetUnidad - Unhandled Error!",
					"err":     err,
					"status":  "500",
				}
				panic(outputError)
			}
		}()

		var unidadAux *models.Unidad

		urlUnidad := "http://" + beego.AppConfig.String("administrativaService") + "unidad/" + strconv.Itoa(unidadId)
		if response, err := request.GetJsonTest(urlUnidad, &unidadAux); err == nil { // (2) error servicio caido
			if response.StatusCode == 200 { // (3) error estado de la solicitud
				unidad = append(unidad, unidadAux)
				return unidad, nil
			} else {
				err := fmt.Errorf("Undesired Status: %s", response.Status)
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "GetUnidad - request.GetJsonTest(urlUnidad, &unidadAux) / response.StatusCode == 200",
					"err":     err,
					"status":  "500", // Error (3) estado de la solicitud
				}
				return nil, outputError
			}
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetUnidad - request.GetJsonTest(urlUnidad, &unidadAux)",
				"err":     err,
				"status":  "502", // Error (2) servicio caido
			}
			return nil, outputError
		}
	} else {
		err := fmt.Errorf("unidadId MUST be greater than 0")
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetUnidad",
			"err":     err,
			"status":  "400", // null parameter
		}
		return nil, outputError
	}
}

func GetUnidades(respuesta interface{}) (outputError map[string]interface{}) {
	urlUnidad := "http://" + beego.AppConfig.String("AdministrativaService") + "unidad?limit=-1"
	if _, err := request.GetJsonTest(urlUnidad, &respuesta); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetUnidades - request.GetJsonTest(urlUnidad, &Unidades)",
			"err":     err,
			"status":  "502",
		}
	}
	return
}
