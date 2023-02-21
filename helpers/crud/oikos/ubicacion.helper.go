package oikos

import (
	"regexp"
	"strconv"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

func GetSedeDependenciaUbicacion(ubicacionId int) (DetalleUbicacion *models.DetalleSedeDependencia, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GetSedeDependenciaUbicacion - Unhandled Error!", "500")

	var (
		payload   string
		ubicacion []models.AsignacionEspacioFisicoDependencia
	)

	resultado := new(models.DetalleSedeDependencia)

	payload = "query=Id:" + strconv.Itoa(ubicacionId)
	if ubicacion_, err := GetAllAsignacion(payload); err != nil {
		return nil, err
	} else {
		ubicacion = ubicacion_
	}

	resultado.Dependencia = ubicacion[0].DependenciaId
	resultado.Ubicacion = &ubicacion[0]

	rgxp := regexp.MustCompile(`\d.*`)
	strSede := ubicacion[0].EspacioFisicoId.CodigoAbreviacion
	strSede = strSede[0:2] + rgxp.ReplaceAllString(strSede[2:], "")

	payload = "?query=CodigoAbreviacion:" + strSede
	if sede_, err := GetAllEspacioFisico(payload); err != nil {
		return nil, err
	} else {
		resultado.Sede = &sede_[0]
	}

	return resultado, nil
}
