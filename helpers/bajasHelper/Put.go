package bajasHelper

import (
	"context"
	"strconv"

	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// Put Actualiza información de baja
func Put(ctx context.Context, baja *models.TrSoporteMovimiento, bajaId int) (bajaR *models.Movimiento, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("Put - Unhandled Error!", "500")

	// Actualiza registro en api movimientos_arka_crud
	outputError = movimientosArka.PutMovimiento(ctx, baja.Movimiento, bajaId)
	if outputError != nil {
		return
	}

	// Actualiza el documento soporte en la tabla soporte_movimiento
	var soporte models.SoporteMovimiento
	soportes, outputError := movimientosArka.GetAllSoporteMovimiento(ctx, "limit=1&sortby=Id&order=desc&query=Activo:true,MovimientoId__Id:"+strconv.Itoa(bajaId))
	if outputError != nil {
		return
	} else if len(soportes) == 1 {
		soporte = soportes[0]
		soporte.DocumentoId = baja.Soporte.DocumentoId
		_, outputError = movimientosArka.PutSoporteMovimiento(ctx, &soporte, soporte.Id)
	} else {
		soporte.MovimientoId = &models.Movimiento{Id: bajaId}
		soporte.Activo = true
		soporte.DocumentoId = baja.Soporte.DocumentoId
		outputError = movimientosArka.PostSoporteMovimiento(ctx, &soporte)
	}

	return
}
