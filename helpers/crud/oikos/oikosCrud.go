package oikos

import (
	"regexp"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	"github.com/udistrital/arka_mid/utils_oas/request"
)

var basePath, _ = beego.AppConfig.String("oikosService")

func GetAllAsignacion(payload string) (asignaciones []models.AsignacionEspacioFisicoDependencia, outputError map[string]interface{}) {

	funcion := "GetAllAsignacion - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := "http://" + basePath + "asignacion_espacio_fisico_dependencia?" + payload
	_, err := request.GetJsonTest(urlcrud, &asignaciones)
	if err != nil {
		logs.Info(urlcrud)
		eval := "request.GetJsonTest(urlcrud, &asignaciones)"
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

func GetAllEspacioFisico(payload string) (espacios []models.EspacioFisico, outputError map[string]interface{}) {

	funcion := "GetAllEspacioFisico - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := "http://" + basePath + "espacio_fisico?" + payload
	err := request.GetJson(urlcrud, &espacios)
	if err != nil {
		logs.Error(err)
		eval := `request.GetJson(urlcrud, &espacios)`
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

func GetAllEspacioFisicoCampo(payload string) (espacios []models.EspacioFisicoCampo, outputError map[string]interface{}) {

	funcion := "GetAllEspacioFisicoCampo - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := "http://" + basePath + "espacio_fisico_campo?" + payload
	err := request.GetJson(urlcrud, &espacios)
	if err != nil {
		logs.Error(err)
		eval := `request.GetJson(urlcrud, &espacios)`
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetSedeEspacioFisico
func GetSedeEspacioFisico(espacioFisico models.EspacioFisico) (sede models.EspacioFisico, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("GetSedeEspacioFisico - Unhandled Error", "500")

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
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := "http://" + basePath + "dependencia/" + strconv.Itoa(id)
	err := request.GetJson(urlcrud, &dependencia)
	if err != nil {
		eval := "request.GetJson(urlcrud, &dependencia)"
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}
