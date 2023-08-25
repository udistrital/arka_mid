package salidaHelper

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

func RechazarSalida(salida *models.Movimiento) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("RechazarSalida - Unhandled Error!", "500")

	query := "limit=1&query=Id:" + strconv.Itoa(salida.Id)
	if mov, err := movimientosArka.GetAllMovimiento(query); err != nil {
		return err
	} else if len(mov) == 1 && mov[0].EstadoMovimientoId.Nombre == "Salida En Trámite" {
		*salida = *mov[0]
	} else {
		return
	}

	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&salida.EstadoMovimientoId.Id, "Salida Rechazada")
	if outputError != nil {
		return
	}

	_, outputError = movimientosArka.PutMovimiento(salida, salida.Id)

	return
}
