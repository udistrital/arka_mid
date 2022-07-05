package salidaHelper

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

func RechazarSalida(salida *models.Movimiento) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("RechazarSalida - Unhandled Error!", "500")

	query := "limit=1&query=Id:" + strconv.Itoa(salida.Id)
	if mov, err := movimientosArka.GetAllMovimiento(query); err != nil {
		return err
	} else if len(mov) == 1 && mov[0].EstadoMovimientoId.Nombre == "Salida En Trámite" {
		*salida = *mov[0]
	} else {
		return
	}

	if err := movimientosArka.GetEstadoMovimientoIdByNombre(&salida.EstadoMovimientoId.Id, "Salida Rechazada"); err != nil {
		return err
	}

	if _, err := movimientosArka.PutMovimiento(salida, salida.Id); err != nil {
		return err
	}

	return
}
