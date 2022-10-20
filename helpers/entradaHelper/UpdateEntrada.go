package entradaHelper

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// UpdateEntrada Consulta el tipo de movimiento y completa el detalle de una entrada que se quiere actualizar
func UpdateEntrada(data *models.TransaccionEntrada, movimientoId int, resultado *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	funcion := "UpdateEntrada - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	var (
		tipoMovimiento   int
		estadoMovimiento int
		detalle          string
		query            string
		consecutivo      models.ConsecutivoMovimiento
	)

	if data.Detalle.ActaRecibidoId <= 0 {
		err := "Se debe indicar un acta de recibido válida."
		return errorctrl.Error(funcion, err, "400")
	}

	query = "limit=1&query=Id:" + strconv.Itoa(movimientoId)
	if mov, err := movimientosArka.GetAllMovimiento(query); err != nil {
		return err
	} else if len(mov) == 1 && mov[0].EstadoMovimientoId.Nombre == "Entrada Rechazada" {
		*&resultado.Movimiento = *mov[0]
	} else {
		return
	}

	if err := movimientosArka.GetEstadoMovimientoIdByNombre(&estadoMovimiento, "Entrada En Trámite"); err != nil {
		return err
	}

	if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&tipoMovimiento, data.FormatoTipoMovimientoId); err != nil {
		return err
	}

	if err := utilsHelper.Unmarshal(resultado.Movimiento.Detalle, &consecutivo); err != nil {
		return err
	}

	if err := crearDetalleEntrada(&data.Detalle, false, &consecutivo, &detalle); err != nil {
		return err
	}

	resultado.Movimiento = models.Movimiento{
		Id:                      movimientoId,
		Observacion:             data.Observacion,
		Detalle:                 detalle,
		Activo:                  true,
		FormatoTipoMovimientoId: &models.FormatoTipoMovimiento{Id: tipoMovimiento},
		EstadoMovimientoId:      &models.EstadoMovimiento{Id: estadoMovimiento},
	}

	if _, err := movimientosArka.PutMovimiento(&resultado.Movimiento, movimientoId); err != nil {
		return err
	}

	// Crea registro en table soporte_movimiento si es necesario
	if data.SoporteMovimientoId > 0 {
		soporte := models.SoporteMovimiento{
			DocumentoId:  data.SoporteMovimientoId,
			Activo:       true,
			MovimientoId: &models.Movimiento{Id: resultado.Movimiento.Id},
		}

		if err := movimientosArka.PostSoporteMovimiento(&soporte); err != nil {
			return err
		}

	}

	return
}
