package ajustesHelper

import (
	"math"
	"net/url"
	"strconv"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// generaTrContable Dado un valor, subgrupo nuevo y original genera la transacción contable.
func generaTrContable(vInicial, vNuevo float64,
	consecutivo, tipoMedicion string,
	db, cr, sgOriginal, sgNuevo, tercero int,
	ctasSg map[int]*models.CuentasSubgrupo,
	ctas map[string]*models.CuentaContable) (movimientos []*models.MovimientoTransaccion) {

	dsc := getDescripcionMovContable(tipoMedicion, consecutivo)
	if sgOriginal > 0 {
		if ctasSg[sgOriginal] != nil && ctasSg[sgNuevo] != nil {
			if ctasSg[sgOriginal].CuentaCreditoId != ctasSg[sgNuevo].CuentaCreditoId {
				movimientoR := asientoContable.CreaMovimiento(vInicial, dsc, tercero, ctas[ctasSg[sgOriginal].CuentaCreditoId], db)
				movimiento := asientoContable.CreaMovimiento(vNuevo, dsc, tercero, ctas[ctasSg[sgNuevo].CuentaCreditoId], cr)
				movimientos = append(movimientos, movimientoR, movimiento)
			} else if vInicial != vNuevo {
				tipoMovimiento := cr
				if vNuevo-vInicial < 0 {
					tipoMovimiento = db
				}

				movimiento := asientoContable.CreaMovimiento(math.Abs(vNuevo-vInicial), dsc, tercero, ctas[ctasSg[sgNuevo].CuentaCreditoId], tipoMovimiento)
				movimientos = append(movimientos, movimiento)
			}

			if ctasSg[sgOriginal].CuentaDebitoId != ctasSg[sgNuevo].CuentaDebitoId {
				movimientoR := asientoContable.CreaMovimiento(vInicial, dsc, tercero, ctas[ctasSg[sgOriginal].CuentaDebitoId], cr)
				movimiento := asientoContable.CreaMovimiento(vNuevo, dsc, tercero, ctas[ctasSg[sgNuevo].CuentaDebitoId], db)
				movimientos = append(movimientos, movimientoR, movimiento)
			} else if vInicial != vNuevo {
				tipoMovimiento := db
				if vNuevo-vInicial < 0 {
					tipoMovimiento = cr
				}

				movimiento := asientoContable.CreaMovimiento(math.Abs(vNuevo-vInicial), dsc, tercero, ctas[ctasSg[sgNuevo].CuentaDebitoId], tipoMovimiento)
				movimientos = append(movimientos, movimiento)
			}
		}

	} else if ctasSg[sgNuevo] != nil {

		tipoMovimientoC := cr
		tipoMovimientoD := db

		if vNuevo-vInicial < 0 {
			tipoMovimientoC = db
			tipoMovimientoD = cr
		}

		movimientoC := asientoContable.CreaMovimiento(math.Abs(vNuevo-vInicial), dsc, tercero, ctas[ctasSg[sgNuevo].CuentaCreditoId], tipoMovimientoC)
		movimientoD := asientoContable.CreaMovimiento(math.Abs(vNuevo-vInicial), dsc, tercero, ctas[ctasSg[sgNuevo].CuentaDebitoId], tipoMovimientoD)
		movimientos = append(movimientos, movimientoC, movimientoD)

	}

	return movimientos

}

// getCuentasByMovimientoSubgrupos Retorna las cuentas de cada subgrupo en una estructura para fácil acceso
func getCuentasByMovimientoSubgrupos(movimientoId int, subgrupos []int) (
	cuentasSubgrupo map[int]*models.CuentasSubgrupo, outputError map[string]interface{}) {

	funcion := "getCuentasByMovimientoSubgrupos"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	cuentasSubgrupo = make(map[int]*models.CuentasSubgrupo)

	query := "limit=-1&fields=CuentaDebitoId,CuentaCreditoId,SubgrupoId&sortby=Id&order=desc&"
	query += "query=SubtipoMovimientoId:" + strconv.Itoa(movimientoId) + ",Activo:true,SubgrupoId__Id__in:"
	query += url.QueryEscape(utilsHelper.ArrayToString(subgrupos, "|"))
	if cuentas_, err := catalogoElementos.GetAllCuentasSubgrupo(query); err != nil {
		return nil, err
	} else {
		for _, cuenta := range cuentas_ {
			cuentasSubgrupo[cuenta.SubgrupoId.Id] = cuenta
		}
	}

	return cuentasSubgrupo, nil

}

func joinMaps(map1, map2 map[int]*models.CuentasSubgrupo) map[int]*models.CuentasSubgrupo {

	if len(map1) > 0 {
		for sg, ctas := range map2 {
			map1[sg] = ctas
		}
		return map1
	} else if len(map2) > 0 {
		for sg, ctas := range map1 {
			map2[sg] = ctas
		}
		return map2
	}

	return nil

}

// fillCuentas Consulta el detalle de una serie de cuentas
func fillCuentas(cuentas map[string]*models.CuentaContable, cuentas_ []string) (cuentasCompletas map[string]*models.CuentaContable, outputError map[string]interface{}) {

	funcion := "fillCuentas"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	for _, id := range cuentas_ {
		if _, ok := cuentas[id]; !ok {
			if cta_, err := cuentasContables.GetCuentaContable(id); err != nil {
				return nil, err
			} else {
				cuentas[id] = cta_
			}
		}
	}

	return cuentas, nil

}

// findElementoInArray Retorna la posicion en que se encuentra el id específicado
func findElementoInArrayD(elementos []*models.DetalleElemento_, id int) (i int) {
	for i, el_ := range elementos {
		if int(el_.Id) == id {
			return i
		}
	}
	return -1
}

// findElementoInArray Retorna la posicion en que se encuentra el id específicado
func findElementoInArrayE(elementos []*models.Elemento, id int) (i int) {
	for i, el_ := range elementos {
		if int(el_.Id) == id {
			return i
		}
	}
	return -1
}

// findElementoInArray Retorna la posicion en que se encuentra el id específicado
func findElementoInArrayEM(elementos []*models.DetalleElemento, id int) (i int) {
	for i, el_ := range elementos {
		if int(el_.Id) == id {
			return i
		}
	}
	return -1
}

func getTipoComprobanteAjustes() string {
	return "N39"
}

func getDescripcionMovContable(tipoMovimiento, consecutivo string) string {
	return "Ajuste contable " + tipoMovimiento + " " + consecutivo
}
