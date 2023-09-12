package cuentasContables

import (
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// GetCuentaContable Consulta controlador nodo_cuenta_contable/{UUID}
func GetDetalleCuentasContables(cuentas []string, detalleCuentas map[string]models.CuentaContable) (outputError map[string]interface{}) {

	funcion := "GetDetalleCuentasContables"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	for _, cuenta := range cuentas {
		if _, ok := detalleCuentas[cuenta]; !ok {
			if cta, err := GetCuentaContable(cuenta); err != nil {
				return err
			} else {
				detalleCuentas[cuenta] = *cta
			}
		}
	}

	return

}
