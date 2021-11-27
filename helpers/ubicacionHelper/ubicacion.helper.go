package ubicacionHelper

import (
	"regexp"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

func GetAsignacionSedeDependencia(Id string) (Relacion map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetAsignacionSedeDependencia - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var ubicacion []map[string]interface{}
	relacion := make(map[string]interface{}, 0)

	urlcrud := "http://" + beego.AppConfig.String("oikos2Service") + "asignacion_espacio_fisico_dependencia?query=Id:" + Id

	if _, err := request.GetJsonTest(urlcrud, &ubicacion); err == nil { // (2) error servicio caido

		if keys := len(ubicacion[0]); keys != 0 {

			return ubicacion[0], nil

		} else {
			return relacion, nil
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAsignacionSedeDependencia",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}

func GetSedeDependenciaUbicacion(ubicacionId int) (DetalleUbicacion *models.DetalleSedeDependencia, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetSedeDependenciaUbicacion - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()
	var (
		urlcrud   string
		ubicacion []*models.AsignacionEspacioFisicoDependencia
	)
	resultado := new(models.DetalleSedeDependencia)

	urlcrud = "?query=Id:" + strconv.Itoa(ubicacionId)
	if ubicacion_, err := GetAllAsignacion(urlcrud); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetSedeDependenciaUbicacion - GetAllAsignacion(urlcrud)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	} else {
		ubicacion = ubicacion_
	}

	resultado.Dependencia = ubicacion[0].DependenciaId
	resultado.Ubicacion = ubicacionId

	if espFisico, err := utilsHelper.ConvertirInterfaceMap(ubicacion[0].EspacioFisicoId); err != nil {
		return nil, err
	} else {
		rgxp := regexp.MustCompile("[0-9]")
		strSede := espFisico["CodigoAbreviacion"].(string)
		strSede = rgxp.ReplaceAllString(strSede, "")
		urlcrud = "?query=CodigoAbreviacion:" + strSede
	}

	if sede_, err := GetAllEspacioFisico(urlcrud); err != nil {
		return nil, err
	} else {
		resultado.Sede = sede_[0]
	}

	return resultado, nil
}
