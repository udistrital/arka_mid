package asientoContable

import (
	"time"

	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

var getComprobanteCreateTransaccionContable = cuentasContables.GetComprobante

// CreateTransaccionContable Consulta el tipo de comprobante y completa otros datos de la transacción contable
func CreateTransaccionContable(tipoComprobante, dsc string, transaccion *models.TransaccionMovimientos) (msg string, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("CreateTransaccionContable - Unhandled Error!", "500")

	var comprobanteID string

	if tipoComprobante == "" {
		return "No se pudo consultar el comprobante contable. Contacte soporte.", nil
	}

	if err := getComprobanteCreateTransaccionContable(tipoComprobante, &comprobanteID); err != nil {
		return "", err
	}

	if comprobanteID == "" {
		return "No se pudo consultar el comprobante contable. Contacte soporte.", nil
	}

	etiquetas := *new(models.Etiquetas)
	etiquetas.ComprobanteId = comprobanteID
	if err := utilsHelper.Marshal(etiquetas, &transaccion.Etiquetas); err != nil {
		return "", err
	}

	transaccion.Descripcion = dsc
	if transaccion.FechaTransaccion.IsZero() {
		transaccion.FechaTransaccion = time.Now()
	}
	transaccion.Activo = true

	return

}
