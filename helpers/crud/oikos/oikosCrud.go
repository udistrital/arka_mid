package oikos

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

var basePath = "http://" + beego.AppConfig.String("oikosService")

func GetAllAsignacion(query string) (asignacion []*models.AsignacionEspacioFisicoDependencia, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetAllAsignacion - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var ubicacion []*models.AsignacionEspacioFisicoDependencia

	urlcrud := "http://" + beego.AppConfig.String("oikosService") + "asignacion_espacio_fisico_dependencia" + query
	if _, err := request.GetJsonTest(urlcrud, &ubicacion); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllAsignacion - request.GetJsonTest(urlcrud, &ubicacion)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	return ubicacion, nil
}

func GetAllEspacioFisico(query string) (espacio []*models.EspacioFisico, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetAllEspacioFisico - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var espacioFisico []*models.EspacioFisico

	urlcrud := "http://" + beego.AppConfig.String("oikosService") + "espacio_fisico" + query
	if _, err := request.GetJsonTest(urlcrud, &espacioFisico); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllEspacioFisico - request.GetJsonTest(urlcrud, &espacioFisico)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	return espacioFisico, nil
}

// GetDependenciaById consulta controlador dependencia/{id} del api oikos_crud
func GetDependenciaById(id int, dependencia *models.Dependencia) (outputError map[string]interface{}) {

	funcion := "GetDependenciaById - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := basePath + "dependencia/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &dependencia); err != nil {
		eval := "request.GetJson(urlcrud, &dependencia)"
		return errorctrl.Error(funcion+eval, err, "502")
	}

	return nil
}
