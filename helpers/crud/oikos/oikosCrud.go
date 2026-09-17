package oikos

import (
	"context"
	"regexp"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/models"
	errorCtrl "github.com/udistrital/utils_oas/v2/errorctrl"
	requestV2 "github.com/udistrital/utils_oas/v2/request"
)

var basePath, _ = beego.AppConfig.String("oikosService")

func GetAllAsignacion(ctx context.Context, payload string) (asignaciones []models.AsignacionEspacioFisicoDependencia, outputError map[string]interface{}) {

	funcion := "GetAllAsignacion - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := basePath + "asignacion_espacio_fisico_dependencia?" + payload
	_, err := requestV2.GetWithContext(ctx, urlcrud, &asignaciones)
	if err != nil {
		logs.Info(urlcrud)
		eval := "requestV2.GetWithContext(ctx, urlcrud, &asignaciones)"
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

func GetAllEspacioFisico(ctx context.Context, payload string) (espacios []models.EspacioFisico, outputError map[string]interface{}) {

	funcion := "GetAllEspacioFisico - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := basePath + "espacio_fisico?" + payload
	_, err := requestV2.GetWithContext(ctx, urlcrud, &espacios)
	if err != nil {
		logs.Error(err)
		eval := "requestV2.GetWithContext(ctx, urlcrud, &espacios)"
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

func GetAllEspacioFisicoCampo(ctx context.Context, payload string) (espacios []models.EspacioFisicoCampo, outputError map[string]interface{}) {

	funcion := "GetAllEspacioFisicoCampo - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := basePath + "espacio_fisico_campo?" + payload
	_, err := requestV2.GetWithContext(ctx, urlcrud, &espacios)
	if err != nil {
		logs.Error(err)
		eval := "requestV2.GetWithContext(ctx, urlcrud, &espacios)"
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetSedeEspacioFisico
func GetSedeEspacioFisico(ctx context.Context, espacioFisico models.EspacioFisico) (sede models.EspacioFisico, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("GetSedeEspacioFisico - Unhandled Error", "500")

	rgxp := regexp.MustCompile(`\d.*`)
	codigoSede := espacioFisico.CodigoAbreviacion
	codigoSede = codigoSede[0:2] + rgxp.ReplaceAllString(codigoSede[2:], "")

	payload := "query=TipoEspacioFisicoId__Nombre:SEDE,CodigoAbreviacion:" + codigoSede
	sede_, outputError := GetAllEspacioFisico(ctx, payload)
	if outputError != nil {
		return
	}

	if len(sede_) > 0 {
		sede = sede_[0]
	}

	return
}

// GetDependenciaById consulta controlador dependencia/{id} del api oikos_crud
func GetDependenciaById(ctx context.Context, id int) (dependencia *models.Dependencia, outputError map[string]interface{}) {

	funcion := "GetDependenciaById - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := basePath + "dependencia/" + strconv.Itoa(id)
	_, err := requestV2.GetWithContext(ctx, urlcrud, dependencia)
	if err != nil {
		eval := "requestV2.GetWithContext(ctx, urlcrud, dependencia)"
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}
