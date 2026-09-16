package cuentasContables

import (
	"context"

	"github.com/udistrital/arka_mid/models"
	errorCtrl "github.com/udistrital/utils_oas/v2/errorctrl"
)

// GetCuentaContable Consulta controlador nodo_cuenta_contable/{UUID}
func GetDetalleCuentasContables(ctx context.Context, cuentas []string, detalleCuentas map[string]models.CuentaContable) (outputError map[string]interface{}) {

	funcion := "GetDetalleCuentasContables"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	for _, cuenta := range cuentas {
		if _, ok := detalleCuentas[cuenta]; !ok {
			if cta, err := GetCuentaContable(ctx, cuenta); err != nil {
				return err
			} else {
				detalleCuentas[cuenta] = *cta
			}
		}
	}

	return

}
