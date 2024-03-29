package asientoContable

import (
	"strconv"
	"strings"

	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
)

// CalcularMovimientosContables Calcula los movimientos contables dados los valores y parametrización correspondiente de cada elemento.
func CalcularMovimientosContables(elementos []*models.Elemento, dsc string, movId, sMovId, terceroCr, terceroDb int,
	cuentas map[string]models.CuentaContable, subgrupos map[int]models.DetalleSubgrupo, movimientos *[]*models.MovimientoTransaccion) (
	errMsg string, outputError map[string]interface{}) {

	if subgrupos == nil {
		subgrupos = make(map[int]models.DetalleSubgrupo)
	}

	var parCr int
	var parDb int
	var uvt float64 = 1

	var payload = "limit=1&fields=TipoBienId,Amortizacion,Depreciacion,SubgrupoId&sortby=Id&order=desc&query=Activo:true,SubgrupoId__Id:"

	if cuentas == nil {
		cuentas = make(map[string]models.CuentaContable)
	}

	// if uvt_, err := parametros.GetUVTByVigencia(time.Now().Year()); err != nil {
	// 	return "", err
	// } else if uvt_ == 0 {
	// 	return "No se pudo consultar el valor del UVT. Intente más tarde o contacte soporte.", nil
	// } else {
	// 	uvt = uvt_
	// }

	if db_, cr_, err := parametros.GetParametrosDebitoCredito(); err != nil {
		return "", err
	} else {
		parCr = cr_
		parDb = db_
	}

	tiposBien := make(map[int]models.TipoBien)
	cuentasSgTb := make(map[int]map[int]models.CuentasSubgrupo)
	totalesCr := make(map[string]float64)
	totalesDb := make(map[string]float64)
	var actasConflicto []int
	var subgruposConflicto []string

	for _, el := range elementos {

		if el.ValorTotal == 0 {
			continue
		}

		if el.SubgrupoCatalogoId <= 0 {
			return "No se pudo determinar la clase de los elementos. Revise el detalle del acta de recibido o contacte soporte.", nil
		}

		if _, ok := subgrupos[el.SubgrupoCatalogoId]; !ok {
			sg, outputError := catalogoElementos.GetAllDetalleSubgrupo(payload + strconv.Itoa(el.SubgrupoCatalogoId))
			if outputError != nil {
				return "", outputError
			} else if len(sg) == 0 {
				return "No se pudo consultar la parametrización de las clases. Contacte soporte.", nil
			}

			subgrupos[el.SubgrupoCatalogoId] = *sg[0]

		}

		if el.TipoBienId == 0 {
			tb, outputError := catalogoElementos.GetTipoBienIdByValor(subgrupos[el.SubgrupoCatalogoId].TipoBienId.Id, el.ValorUnitario/uvt, tiposBien)
			if outputError != nil {
				return "", outputError
			} else if tb == 0 {
				return "No se pudo establecer el tipo de bien de los elementos. Revise la parametrización de los tipos de bien.", nil
			}

			el.TipoBienId = tb

		} else {
			if _, ok := tiposBien[el.TipoBienId]; !ok {
				var tipoBien models.TipoBien
				outputError = catalogoElementos.GetTipoBienById(el.TipoBienId, &tipoBien)
				if outputError != nil {
					return
				}
				tiposBien[el.TipoBienId] = tipoBien
			}

			if tiposBien[el.TipoBienId].TipoBienPadreId.Id != subgrupos[el.SubgrupoCatalogoId].TipoBienId.Id {
				var elemento models.Elemento
				outputError = actaRecibido.GetElementoById(el.Id, &elemento)
				if outputError != nil {
					return "", outputError
				}

				exists := containsInt(actasConflicto, elemento.ActaRecibidoId.Id)
				if !exists {
					actasConflicto = append(actasConflicto, elemento.ActaRecibidoId.Id)
				}
				continue
			}
		}

		if _, ok := cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId]; !ok {
			cst, outputError := catalogoElementos.GetAllCuentasSubgrupo(payloadCuentas(el.SubgrupoCatalogoId, el.TipoBienId, movId, sMovId))
			if outputError != nil {
				return "", outputError
			} else if len(cst) == 1 {
				if cuentasSgTb[el.SubgrupoCatalogoId] == nil {
					cuentasSgTb[el.SubgrupoCatalogoId] = make(map[int]models.CuentasSubgrupo)
				}
				cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId] = *cst[0]
			} else {
				exists := containsString(subgruposConflicto, subgrupos[el.SubgrupoCatalogoId].SubgrupoId.Codigo)
				if !exists {
					subgruposConflicto = append(subgruposConflicto, subgrupos[el.SubgrupoCatalogoId].SubgrupoId.Codigo)
				}
				continue
			}
		}

		if _, ok := cuentas[cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId].CuentaCreditoId]; !ok {
			cr, outputError := cuentasContables.GetCuentaContable(cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId].CuentaCreditoId)
			if outputError != nil {
				return "", outputError
			} else if cr != nil {
				cuentas[cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId].CuentaCreditoId] = *cr
			} else {
				return "No se pudo encontrar la cuenta contable. Contacte soporte", nil
			}
		}

		if _, ok := cuentas[cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId].CuentaDebitoId]; !ok {
			db, outputError := cuentasContables.GetCuentaContable(cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId].CuentaDebitoId)
			if outputError != nil {
				return "", outputError
			} else if db != nil {
				cuentas[cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId].CuentaDebitoId] = *db
			} else {
				return "No se pudo encontrar la cuenta contable. Contacte soporte", nil
			}

		}

		totalesCr[cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId].CuentaCreditoId] += el.ValorTotal
		totalesDb[cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId].CuentaDebitoId] += el.ValorTotal

	}

	if len(actasConflicto) > 0 {
		return "El tipo bien asignado manualmente no corresponde a la clase correspondiente. Revise las siguientes actas: " + utilsHelper.ArrayToString(actasConflicto, ", "), nil
	} else if len(subgruposConflicto) > 0 {
		return "No se pudo establecer la parametrización contable de las siguientes clases: " + strings.Join(subgruposConflicto, ", "), nil
	}

	for cta, val := range totalesCr {
		var movimiento models.MovimientoTransaccion
		fillMovimiento(val, dsc, terceroCr, parCr, cuentas[cta], &movimiento)
		*movimientos = append(*movimientos, &movimiento)
	}

	for cta, val := range totalesDb {
		var movimiento models.MovimientoTransaccion
		fillMovimiento(val, dsc, terceroDb, parDb, cuentas[cta], &movimiento)
		*movimientos = append(*movimientos, &movimiento)
	}

	return

}

func fillMovimiento(valor float64, dsc string, terceroId, tipoMov int, cuenta models.CuentaContable, movimiento *models.MovimientoTransaccion) {

	if cuenta.RequiereTercero {
		movimiento.TerceroId = &terceroId
	} else {
		movimiento.TerceroId = nil
	}

	movimiento.CuentaId = cuenta.Id
	movimiento.NombreCuenta = cuenta.Nombre
	movimiento.TipoMovimientoId = tipoMov
	movimiento.Valor = valor
	movimiento.Descripcion = dsc
	movimiento.Activo = true

}

func payloadCuentas(sg, tb, mov, sMov int) string {
	return "fields=CuentaDebitoId,CuentaCreditoId&query=Activo:true,SubgrupoId__Id:" +
		strconv.Itoa(sg) + ",TipoBienId__Id:" + strconv.Itoa(tb) + ",TipoMovimientoId:" + strconv.Itoa(mov) +
		",SubtipoMovimientoId:" + strconv.Itoa(sMov)
}

func containsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
