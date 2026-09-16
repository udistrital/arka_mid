package oikos

import (
	"context"
	"strconv"

	"github.com/udistrital/arka_mid/models"
	errorCtrl "github.com/udistrital/utils_oas/v2/errorctrl"
)

func GetSedeDependenciaUbicacion(ctx context.Context, ubicacionId int) (resultado *models.DetalleSedeDependencia, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("GetSedeDependenciaUbicacion - Unhandled Error!", "500")

	resultado = new(models.DetalleSedeDependencia)

	payload := "query=Id:" + strconv.Itoa(ubicacionId)
	ubicacion, outputError := GetAllAsignacion(ctx, payload)
	if outputError != nil || len(ubicacion) == 0 {
		return
	}

	resultado.Dependencia = ubicacion[0].DependenciaId
	resultado.Ubicacion = &ubicacion[0]

	sede, outputError := GetSedeEspacioFisico(ctx, *resultado.Ubicacion.EspacioFisicoId)
	if outputError != nil {
		return
	}

	resultado.Sede = &sede
	return
}
