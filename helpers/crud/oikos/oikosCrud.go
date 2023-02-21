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

func GetAllAsignacion(payload string) (asignaciones []models.AsignacionEspacioFisicoDependencia, outputError map[string]interface{}) {

	funcion := "GetAllAsignacion - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := "http://" + beego.AppConfig.String("oikosService") + "asignacion_espacio_fisico_dependencia?" + payload
	if _, err := request.GetJsonTest(urlcrud, &asignaciones); err != nil {
		eval := "request.GetJsonTest(urlcrud, &asignaciones)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}

	return asignaciones, nil
}

func GetAllEspacioFisico(payload string) (espacios []models.EspacioFisico, outputError map[string]interface{}) {

	funcion := "GetAllEspacioFisico - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := basePath + "espacio_fisico?" + payload
	err := request.GetJson(urlcrud, &espacios)
	if err != nil {
		logs.Error(err)
		eval := `request.GetJson(urlcrud, &espacios)`
		outputError = errorctrl.Error(funcion+eval, err, "502")
	}

	return
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
