package salidaHelper

import (
	"context"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/models"
	errorCtrl "github.com/udistrital/utils_oas/v2/errorctrl"
)

func RechazarSalida(ctx context.Context, id int) (salida *models.Movimiento, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("RechazarSalida - Unhandled Error!", "500")

	salida, outputError = movimientosArka.GetMovimientoById(ctx, id)
	if outputError != nil || salida.EstadoMovimientoId.Nombre != "Salida En Trámite" {
		return
	}

	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(ctx, &salida.EstadoMovimientoId.Id, "Salida Rechazada")
	if outputError != nil {
		return
	}

	outputError = movimientosArka.PutMovimiento(ctx, salida, salida.Id)

	return
}
