package bajasHelper

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// Post Crea registro de baja
func Post(baja *models.TrSoporteMovimiento) (bajaR *models.Movimiento, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("Post - Unhandled Error!", "500")

	var consecutivo models.Consecutivo
	ctxConsecutivo, _ := beego.AppConfig.Int("contxtBajaCons")
	outputError = consecutivos.Get(ctxConsecutivo, "Registro Baja Arka", &consecutivo)
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
