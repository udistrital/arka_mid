package asientoContable

import (
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// GenerarMovimientosContables Genera los movimientos contables para una serie de cuentas y valores
func GenerarMovimientosContables(totales map[int]float64, detalleCuentas map[string]models.CuentaContable, cuentasSubgrupo map[int]models.CuentaSubgrupo,
	parDebito, parCredito, terceroIdCr, terceroIdDb int, descripcion string, ajuste bool, movimientos *[]*models.MovimientoTransaccion) (outputError map[string]interface{}) {

	funcion := "GenerarMovimientosContables"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	if ajuste {
		parDebito, parCredito = parCredito, parDebito
	}

	for sg, valor := range totales {
		ctaCr := detalleCuentas[cuentasSubgrupo[sg].CuentaCreditoId]
		movCr := CreaMovimiento(valor, descripcion, terceroIdCr, &ctaCr, parCredito)
		ctaDb := detalleCuentas[cuentasSubgrupo[sg].CuentaDebitoId]
		movDb := CreaMovimiento(valor, descripcion, terceroIdDb, &ctaDb, parDebito)
		*movimientos = append(*movimientos, movCr, movDb)
	}

	return

}
