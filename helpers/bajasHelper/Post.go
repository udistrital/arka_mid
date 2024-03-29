package bajasHelper

import (
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// Post Crea registro de baja
func Post(baja *models.TrSoporteMovimiento) (bajaR *models.Movimiento, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("Post - Unhandled Error!", "500")

	var consecutivo models.Consecutivo
	outputError = consecutivos.Get("contxtBajaCons", "Registro Baja Arka", &consecutivo)
	if outputError != nil {
		return
	}

	baja.Movimiento.Consecutivo = utilsHelper.String(consecutivos.Format("%05d", getTipoComprobanteBajas(), &consecutivo))
	baja.Movimiento.ConsecutivoId = &consecutivo.Id

	// Crea registro en api movimientos_arka_crud
	outputError = movimientosArka.PostMovimiento(baja.Movimiento)
	if outputError != nil {
		return
	}

	// Crea registro en table soporte_movimiento si es necesario
	baja.Soporte.MovimientoId = baja.Movimiento
	outputError = movimientosArka.PostSoporteMovimiento(baja.Soporte)
	if outputError != nil {
		return
	}

	return baja.Movimiento, nil
}
