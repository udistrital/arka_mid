package asientoContable

import (
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// GenerarMovimientosContables Genera los movimientos contables para una serie de cuentas y valores
func GenerarMovimientosContables(totales map[int]float64, detalleCuentas map[string]models.CuentaContable, cuentasSubgrupo map[int]models.CuentasSubgrupo,
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

// ConstruirMovimientosContables Genera los movimientos contables para una serie de cuentas y valores
func ConstruirMovimientosContables(totales map[int]float64, detalleCuentas map[string]models.CuentaContable, cuentasSubgrupo map[int]models.CuentasSubgrupo,
	terceroIdCr, terceroIdDb int, descripcion string, ajuste bool, movimientos *[]*models.MovimientoTransaccion) (
	err string, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("ConstruirMovimientosContables - Unhandled Error!", "500")

	var parDebito, parCredito int
	if db_, cr_, err := parametros.GetParametrosDebitoCredito(); err != nil {
		return "", err
	} else {
		parDebito = db_
		parCredito = cr_
	}

	if ajuste {
		parDebito, parCredito = parCredito, parDebito
	}

	for sg, valor := range totales {
		if val, ok := cuentasSubgrupo[sg]; ok && val.CuentaCreditoId != "" {
			if val_, ok_ := detalleCuentas[val.CuentaCreditoId]; ok_ && val_.Id != "" {
				movCr := CreaMovimiento(valor, descripcion, terceroIdCr, &val_, parCredito)
				*movimientos = append(*movimientos, movCr)
			} else {
				if subgrupo, err := catalogoElementos.GetSubgrupoById(sg); err != nil {
					return "", err
				} else {
					return "Debe parametrizar las cuentas del subgrupo " + subgrupo.Codigo + " " + subgrupo.Nombre, nil
				}
			}
		} else {
			if subgrupo, err := catalogoElementos.GetSubgrupoById(sg); err != nil {
				return "", err
			} else {
				return "Debe parametrizar las cuentas del subgrupo " + subgrupo.Codigo + " " + subgrupo.Nombre, nil
			}
		}

		if val, ok := cuentasSubgrupo[sg]; ok && val.CuentaDebitoId != "" {
			if val_, ok_ := detalleCuentas[val.CuentaDebitoId]; ok_ && val_.Id != "" {
				movCr := CreaMovimiento(valor, descripcion, terceroIdCr, &val_, parDebito)
				*movimientos = append(*movimientos, movCr)
			} else {
				if subgrupo, err := catalogoElementos.GetSubgrupoById(sg); err != nil {
					return "", err
				} else {
					return "Debe parametrizar las cuentas del subgrupo " + subgrupo.Codigo + " " + subgrupo.Nombre, nil
				}
			}
		} else {
			if subgrupo, err := catalogoElementos.GetSubgrupoById(sg); err != nil {
				return "", err
			} else {
				return "Debe parametrizar las cuentas del subgrupo " + subgrupo.Codigo + " " + subgrupo.Nombre, nil
			}
		}
	}

	return

}
