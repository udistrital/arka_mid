package asientoContable

import (
	"time"

	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// CreateTransaccionContable Consulta el tipo de comprobante y completa otros datos de la transacción contable
func CreateTransaccionContable(tipoComprobante, dsc string, transaccion *models.TransaccionMovimientos) (msg string, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("CreateTransaccionContable - Unhandled Error!", "500")

	var comprobanteID string

	if tipoComprobante == "" {
		return "No se pudo consultar el comprobante contable. Contacte soporte.", nil
	}

	if err := cuentasContables.GetComprobante(tipoComprobante, &comprobanteID); err != nil {
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
	transaccion.FechaTransaccion = time.Now()
	transaccion.Activo = true

	return

}
