package oikos

import (
	"fmt"
	"regexp"

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
	err := request.GetJson(urlcrud, &asignaciones)
	if err != nil {
		logs.Error(err)
		eval := "request.GetJsonTest(urlcrud, &asignaciones)"
		outputError = errorctrl.Error(funcion+eval, err, "502")
	}

	return
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

func GetAllEspacioFisicoCampo(payload string) (espacios []models.EspacioFisicoCampo, outputError map[string]interface{}) {

	funcion := "GetAllEspacioFisicoCampo - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := basePath + "espacio_fisico_campo?" + payload
	err := request.GetJson(urlcrud, &espacios)
	if err != nil {
		logs.Error(err)
		eval := `request.GetJson(urlcrud, &espacios)`
		outputError = errorctrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetSedeEspacioFisico
func GetSedeEspacioFisico(espacioFisico models.EspacioFisico) (sede models.EspacioFisico, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GetSedeEspacioFisico - Unhandled Error", "500")

	rgxp := regexp.MustCompile(`\d.*`)
	codigoSede := espacioFisico.CodigoAbreviacion
	codigoSede = codigoSede[0:2] + rgxp.ReplaceAllString(codigoSede[2:], "")

	payload := "query=TipoEspacioFisicoId__Nombre:SEDE,CodigoAbreviacion:" + codigoSede
	sede_, outputError := GetAllEspacioFisico(payload)
	if outputError != nil {
		return
	}

	if len(sede_) > 0 {
		sede = sede_[0]
	}

	return
}

// GetDependenciaById consulta controlador dependencia/{id} del api oikos_crud
func GetDependenciaById(id int) (dependencia *models.Dependencia, outputError map[string]interface{}) {

	funcion := "GetDependenciaById - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := basePath + "dependencia/" + fmt.Sprint(id)
	err := request.GetJson(urlcrud, &dependencia)
	if err != nil {
		logs.Error(err)
		eval := "request.GetJson(urlcrud, &dependencia)"
		outputError = errorctrl.Error(funcion+eval, err, "502")
	}

	return
}
