package entradaHelper

import (
	"context"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// UpdateEntrada Consulta el tipo de movimiento y completa el detalle de una entrada que se quiere actualizar
func UpdateEntrada(ctx context.Context, data *models.TransaccionEntrada, movimientoId int, resultado *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("UpdateEntrada - Unhandled Error!", "500")

	mov, outputError := movimientosArka.GetMovimientoById(ctx, movimientoId)
	if outputError != nil || mov.EstadoMovimientoId.Nombre != "Entrada Rechazada" {
		return outputError
	}

	resultado.Movimiento = *mov
	resultado.Movimiento.Observacion = data.Observacion
	resultado.Movimiento.Activo = true

	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(ctx, &resultado.Movimiento.EstadoMovimientoId.Id, "Entrada En Trámite")
	if outputError != nil {
		return
	}

	outputError = movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(ctx, &resultado.Movimiento.FormatoTipoMovimientoId.Id, data.FormatoTipoMovimientoId)
	if outputError != nil {
		return
	}

	outputError = crearDetalleEntrada(data.Detalle, &resultado.Movimiento.Detalle)
	if outputError != nil {
		return
	}

	outputError = getConsecutivoEntrada(ctx, &resultado.Movimiento, false)
	if outputError != nil {
		return
	}

	outputError = movimientosArka.PutMovimiento(ctx, &resultado.Movimiento, movimientoId)
	if outputError != nil {
		return
	}

	// Crea registro en table soporte_movimiento si es necesario
	if data.SoporteMovimientoId > 0 {
		soporte := models.SoporteMovimiento{
			DocumentoId:  data.SoporteMovimientoId,
			Activo:       true,
			MovimientoId: &models.Movimiento{Id: resultado.Movimiento.Id},
		}

		outputError = movimientosArka.PostSoporteMovimiento(ctx, &soporte)
		if outputError != nil {
			return
		}
	}

	return
}
