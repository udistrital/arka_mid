package oikos

import (
	"regexp"
	"strconv"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

func GetSedeDependenciaUbicacion(ubicacionId int) (resultado *models.DetalleSedeDependencia, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GetSedeDependenciaUbicacion - Unhandled Error!", "500")

	resultado = new(models.DetalleSedeDependencia)

	payload := "query=Id:" + strconv.Itoa(ubicacionId)
	ubicacion, outputError := GetAllAsignacion(payload)
	if outputError != nil || len(ubicacion) == 0 {
		return
	}

	resultado.Dependencia = ubicacion[0].DependenciaId
	resultado.Ubicacion = &ubicacion[0]

	rgxp := regexp.MustCompile(`\d.*`)
	strSede := ubicacion[0].EspacioFisicoId.CodigoAbreviacion
	strSede = strSede[0:2] + rgxp.ReplaceAllString(strSede[2:], "")

	payload = "?query=CodigoAbreviacion:" + strSede
	sede_, outputError := GetAllEspacioFisico(payload)
	if outputError != nil || len(sede_) == 0 {
		return
	}

	resultado.Sede = &sede_[0]
	return

}
