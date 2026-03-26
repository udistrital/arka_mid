package salidaHelper

import (
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

func RechazarSalida(id int) (salida *models.Movimiento, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("RechazarSalida - Unhandled Error!", "500")

	salida, outputError = movimientosArka.GetMovimientoById(id)
	if outputError != nil || salida.EstadoMovimientoId.Nombre != "Salida En Trámite" {
		return
	}

	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&salida.EstadoMovimientoId.Id, "Salida Rechazada")
	if outputError != nil {
		return
	}

	outputError = movimientosArka.PutMovimiento(salida, salida.Id)

	return
}
