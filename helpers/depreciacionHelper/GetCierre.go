package depreciacionHelper

import (
	"strings"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
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
	var transaccion = new(models.TransaccionMovimientos)
	if detalle_.Movimiento.EstadoMovimientoId.Nombre == "Cierre En Curso" || detalle_.Movimiento.EstadoMovimientoId.Nombre == "Cierre Rechazado" {
		outputError = calcularCierre(detalle_.Movimiento.FechaCorte.UTC().Format("2006-01-02"), nil, transaccion, detalle_)
	} else if detalle_.Movimiento.EstadoMovimientoId.Nombre == "Cierre Aprobado" && detalle_.Movimiento.ConsecutivoId != nil && *detalle_.Movimiento.ConsecutivoId > 0 {
		transaccion, outputError = movimientosContables.GetTransaccion(*detalle_.Movimiento.ConsecutivoId, "consecutivo", true)
	}

	if outputError != nil {
		return
	}

	detalle_.TransaccionContable.Concepto = dscTransaccionCierre()
	detalle_.TransaccionContable.Fecha = transaccion.FechaTransaccion
	detalle_.TransaccionContable.Movimientos, outputError = asientoContable.GetDetalleContable(transaccion.Movimientos, nil)

	return
}
