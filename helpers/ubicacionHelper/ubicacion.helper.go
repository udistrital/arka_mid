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
		sede      []*models.EspacioFisico
	)
	resultado := new(models.DetalleSedeDependencia)

	urlcrud = "http://" + beego.AppConfig.String("oikos2Service") + "asignacion_espacio_fisico_dependencia"
	urlcrud += "?query=Id:" + strconv.Itoa(ubicacionId)

	if _, err := request.GetJsonTest(urlcrud, &ubicacion); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetSedeDependenciaUbicacion - request.GetJsonTest(urlcrud, &ubicacion)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	resultado.Ubicacion = ubicacionId
	resultado.Dependencia = ubicacion[0].DependenciaId

	if espFisico, err := utilsHelper.ConvertirInterfaceMap(ubicacion[0].EspacioFisicoId); err != nil {
		return nil, err
	} else {
		rgxp := regexp.MustCompile("[0-9]")
		strSede := espFisico["CodigoAbreviacion"].(string)
		strSede = rgxp.ReplaceAllString(strSede, "")
		urlcrud = "http://" + beego.AppConfig.String("oikos2Service") + "espacio_fisico?query=CodigoAbreviacion:" + strSede
	}

	if _, err := request.GetJsonTest(urlcrud, &sede); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "TraerDetalle - request.GetJsonTest(urlcrud4, &sede)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	resultado.Sede = sede[0]

	return resultado, nil
}
