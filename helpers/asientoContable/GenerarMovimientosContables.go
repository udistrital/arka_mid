package asientoContable

import (
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// GenerarMovimientosContables Genera los movimientos contables para una serie de cuentas y valores
func GenerarMovimientosContables(totales map[int]float64, detalleCuentas map[string]models.CuentaContable, cuentasSubgrupo map[int]models.CuentaSubgrupo,
	parDebito, parCredito, terceoId int, ajuste bool, movimientos *[]*models.MovimientoTransaccion) (outputError map[string]interface{}) {

	funcion := "GenerarMovimientosContables"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	if ajuste {
		parDebito, parCredito = parCredito, parDebito
	}

	for sg, valor := range totales {
		ctaCr := detalleCuentas[cuentasSubgrupo[sg].CuentaCreditoId]
		ctaDb := detalleCuentas[cuentasSubgrupo[sg].CuentaCreditoId]
		movDb := CreaMovimiento(valor, "", terceoId, &ctaDb, parDebito)
		movCr := CreaMovimiento(valor, "", terceoId, &ctaCr, parCredito)
		*movimientos = append(*movimientos, movCr, movDb)
	}

	return

}
