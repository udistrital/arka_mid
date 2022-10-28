package depreciacionHelper

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// GetCierre Consulta la infomación de un cierre y la transacción contable correspondiente
func GetCierre(id int, detalle_ *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GetCierre - Unhandled Error!", "500")

	var (
		detalle     models.FormatoDepreciacion
		transaccion models.TransaccionMovimientos
		cuentas     map[string]models.CuentaContable
	)

	if mov_, err := movimientosArka.GetAllMovimiento("limit=1&query=Id:" + strconv.Itoa(id)); err != nil {
		return err
	} else if len(mov_) == 1 {
		detalle_.Movimiento = *mov_[0]
	}

	if err := utilsHelper.Unmarshal(detalle_.Movimiento.Detalle, &detalle); err != nil {
		return err
	}

	if detalle_.Movimiento.EstadoMovimientoId.Nombre == "Cierre En Curso" || detalle_.Movimiento.EstadoMovimientoId.Nombre == "Cierre Rechazado" {
		if err := calcularCierre(detalle.FechaCorte, nil, nil, &transaccion, detalle_); err != nil {
			return err
		}
	} else if detalle_.Movimiento.EstadoMovimientoId.Nombre == "Cierre Aprobado" && detalle.ConsecutivoId > 0 {
		if tr, err := movimientosContables.GetTransaccion(detalle.ConsecutivoId, "consecutivo", true); err != nil {
			return err
		} else {
			transaccion = *tr
		}
	}

	if detalleContable, err := asientoContable.GetDetalleContable(transaccion.Movimientos, cuentas); err != nil {
		return err
	} else if len(detalleContable) > 0 {
		trContable := models.InfoTransaccionContable{
			Movimientos: detalleContable,
			Concepto:    dscTransaccionCierre(),
			Fecha:       transaccion.FechaTransaccion,
		}
		detalle_.TransaccionContable = trContable
	}

	return
}
