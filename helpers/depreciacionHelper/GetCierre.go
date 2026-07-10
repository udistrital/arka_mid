package depreciacionHelper

import (
	"strings"

	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// GetCierre Consulta la infomación de un cierre y la transacción contable correspondiente
func GetCierre(id int, detalle_ *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("GetCierre - Unhandled Error!", "500")

	mov, outputError := movimientosArka.GetMovimientoById(id)
	if outputError != nil || mov.FormatoTipoMovimientoId.CodigoAbreviacion != "CRR" || !strings.HasPrefix(mov.EstadoMovimientoId.Nombre, "Cierre ") {
		return
	}

	detalle_.Movimiento = *mov
	var (
		detalle     models.FormatoDepreciacion
		transaccion *models.TransaccionMovimientos
	)
	if err := utilsHelper.Unmarshal(detalle_.Movimiento.Detalle, &detalle); err != nil {
		return err
	}

	detalle_.Error = detalle.CalculoError

	if detalle_.Movimiento.EstadoMovimientoId.Nombre == "Cierre Aprobado" && detalle_.Movimiento.ConsecutivoId != nil && *detalle_.Movimiento.ConsecutivoId > 0 {
		transaccion, outputError = movimientosContables.GetTransaccion(*detalle_.Movimiento.ConsecutivoId, "consecutivo", true)
		if outputError != nil {
			return
		}
		detalle_.TransaccionContable.Concepto = dscTransaccionCierre()
		detalle_.TransaccionContable.Fecha = transaccion.FechaTransaccion
		if detalle.PreviewContable != nil {
			detalle_.TransaccionContable.Movimientos = detalle.PreviewContable.Movimientos
		}
		return nil
	}

	if detalle.CalculoListo && detalle.PreviewContable != nil {
		detalle_.TransaccionContable = *detalle.PreviewContable
	}

	return
}
