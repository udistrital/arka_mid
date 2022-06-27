package oikos

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	e "github.com/udistrital/utils_oas/errorctrl"
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
	defer e.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := basePath + "dependencia/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &dependencia); err != nil {
		eval := "request.GetJson(urlcrud, &dependencia)"
		return e.Error(funcion+eval, err, "502")
	}

	return nil
}

func GetDependencia(query, fields, sortby, order string, limit, offset int, dependencias interface{}) (outputError interface{}) {
	const funcion = "GetDependencias - "
	defer e.ErrorControlFunction(funcion+"unhandled error!", fmt.Sprint(http.StatusInternalServerError))

	urlOikosDependencia := "http://" + beego.AppConfig.String("oikosService") + "dependencia?"
	urlOikosDependencia += utilsHelper.EncodeUrl(query, fields, sortby, order, fmt.Sprint(offset), fmt.Sprint(limit))
	logs.Debug("urlOikosDependencia:", urlOikosDependencia)
	if _, err := request.GetJsonTest(urlOikosDependencia, &dependencias); err != nil {
		logs.Error(err)
		outputError = e.Error(funcion+"request.GetJsonTest(urlOikosDependencia, &dependencias)", err, fmt.Sprint(http.StatusBadGateway))
	}
	return
}
