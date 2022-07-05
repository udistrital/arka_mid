package depreciacionHelper

import (
	"strconv"
	"time"

	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// AprobarDepreciacion Registra las novedades para los elementos depreciados y realiza la transaccion contable
func AprobarDepreciacion(id int, resultado *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("AprobarDepreciacion - Unhandled Error!", "500")

	var (
		movimiento        models.Movimiento
		detalle           *models.FormatoDepreciacion
		elementos         []int
		transaccionCierre models.TransaccionCierre
	)

	if mov_, err := movimientosArka.GetAllMovimiento("limit=1&query=Id:" + strconv.Itoa(id)); err != nil {
		return err
	} else if len(mov_) == 1 && mov_[0].EstadoMovimientoId.Nombre == "Cierre En Curso" {
		movimiento = *mov_[0]
	} else {
		return
	}

	if err := utilsHelper.Unmarshal(movimiento.Detalle, &detalle); err != nil {
		return err
	}

	var transaccion models.TransaccionMovimientos
	if err := calcularCierre(detalle.FechaCorte, &elementos, &transaccion, resultado); err != nil {
		return err
	} else if resultado.Error != "" || len(elementos) == 0 || len(transaccion.Movimientos) == 0 {
		return
	}

	transaccion.Activo = true
	transaccion.ConsecutivoId = detalle.ConsecutivoId
	transaccion.FechaTransaccion = time.Now()
	transaccion.Descripcion = descAsiento()

	if _, err := movimientosContables.PostTrContable(&transaccion); err != nil {
		return err
	}

	transaccionCierre = models.TransaccionCierre{
		MovimientoId:         id,
		ElementoMovimientoId: elementos,
	}
	if err := movimientosArka.AprobarCierre(&transaccionCierre, &resultado.Movimiento); err != nil {
		return err
	}

	return
}
