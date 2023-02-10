package entradaHelper

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// UpdateEntrada Consulta el tipo de movimiento y completa el detalle de una entrada que se quiere actualizar
func UpdateEntrada(data *models.TransaccionEntrada, movimientoId int, resultado *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("UpdateEntrada - Unhandled Error!", "500")

	mov, outputError := movimientosArka.GetAllMovimiento("limit=1&query=Id:" + strconv.Itoa(movimientoId))
	if outputError != nil || len(mov) != 1 || mov[0].EstadoMovimientoId.Nombre != "Entrada Rechazada" {
		return outputError
	}

	resultado.Movimiento = *mov[0]
	resultado.Movimiento.Observacion = data.Observacion
	resultado.Movimiento.Activo = true

	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&resultado.Movimiento.EstadoMovimientoId.Id, "Entrada En Trámite")
	if outputError != nil {
		return
	}

	outputError = movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&resultado.Movimiento.FormatoTipoMovimientoId.Id, data.FormatoTipoMovimientoId)
	if outputError != nil {
		return
	}

	outputError = crearDetalleEntrada(data.Detalle, &resultado.Movimiento.Detalle)
	if outputError != nil {
		return
	}

	outputError = getConsecutivoEntrada(&resultado.Movimiento, false)
	if outputError != nil {
		return
	}

	_, outputError = movimientosArka.PutMovimiento(&resultado.Movimiento, movimientoId)
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

		outputError = movimientosArka.PostSoporteMovimiento(&soporte)
		if outputError != nil {
			return
		}
	}

	return
}
