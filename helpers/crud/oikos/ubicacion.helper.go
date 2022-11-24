package oikos

import (
	"regexp"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
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
	relacion := make(map[string]interface{})

	urlcrud := "http://" + beego.AppConfig.String("oikosService") + "asignacion_espacio_fisico_dependencia?query=Id:" + Id

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

	defer errorctrl.ErrorControlFunction("GetSedeDependenciaUbicacion - Unhandled Error!", "500")

	var (
		payload   string
		ubicacion []*models.AsignacionEspacioFisicoDependencia
	)

	resultado := new(models.DetalleSedeDependencia)

	payload = "?query=Id:" + strconv.Itoa(ubicacionId)
	if ubicacion_, err := GetAllAsignacion(payload); err != nil {
		return nil, err
	} else {
		ubicacion = ubicacion_
	}

	resultado.Dependencia = ubicacion[0].DependenciaId
	resultado.Ubicacion = ubicacion[0]

	rgxp := regexp.MustCompile("\\d.*")
	strSede := ubicacion[0].EspacioFisicoId.CodigoAbreviacion
	strSede = rgxp.ReplaceAllString(strSede, "")

	payload = "?query=CodigoAbreviacion:" + strSede
	if sede_, err := GetAllEspacioFisico(payload); err != nil {
		return nil, err
	} else {
		resultado.Sede = sede_[0]
	}

	return resultado, nil
}
