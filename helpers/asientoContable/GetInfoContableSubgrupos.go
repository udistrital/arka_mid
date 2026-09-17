package asientoContable

import (
	"context"
	"github.com/udistrital/arka_mid/helpers/catalogoElementosHelper"
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/models"
)

// GetInfoContableSubgrupos Consulta las cuentas contables para una lista de subgrupos y un tipo de movimiento así como el detalle de las cuentas contables
func GetInfoContableSubgrupos(ctx context.Context, tipoMovimiento int, subgrupos []int, cuentasSubgrupo map[int]models.CuentasSubgrupo, detalleCuentas map[string]models.CuentaContable) (outputError map[string]interface{}) {

	if err := catalogoElementosHelper.GetCuentasByMovimientoAndSubgrupos(ctx, tipoMovimiento, subgrupos, cuentasSubgrupo); err != nil {
		return err
	}

	idsCuentas := []string{}
	for _, cta := range cuentasSubgrupo {
		idsCuentas = append(idsCuentas, cta.CuentaCreditoId)
		idsCuentas = append(idsCuentas, cta.CuentaDebitoId)
	}
	if len(idsCuentas) == 0 {
		return
	}

	if err := cuentasContables.GetDetalleCuentasContables(ctx, idsCuentas, detalleCuentas); err != nil {
		return err
	}

	return
}
