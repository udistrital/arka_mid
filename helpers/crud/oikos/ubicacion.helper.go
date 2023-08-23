package oikos

import (
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

	sede, outputError := GetSedeEspacioFisico(*resultado.Ubicacion.EspacioFisicoId)
	if outputError != nil {
		return
	}

	resultado.Sede = &sede
	return
}
