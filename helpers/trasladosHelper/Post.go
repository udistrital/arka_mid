package trasladoshelper

import (
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// Post Crea registro de traslado en estado en trámite
func Post(traslado *models.Movimiento) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("Post - Unhandled Error!", "500")

	var consecutivo models.Consecutivo
	outputError = consecutivos.Get("contxtAjusteCons", "Registro Traslado Arka", &consecutivo)
	if outputError != nil {
		return
	}

	traslado.Consecutivo = utilsHelper.String(consecutivos.Format("%05d", getTipoComprobanteTraslados(), &consecutivo))
	traslado.ConsecutivoId = &consecutivo.Id

	outputError = movimientosArka.PostMovimiento(traslado)

	return
}
