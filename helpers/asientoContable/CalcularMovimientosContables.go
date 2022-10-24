package asientoContable

import (
	"strconv"
	"time"

	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/arka_mid/models"
)

// CalcularMovimientosContables Calcula los movimientos contables dados los valores y parametrización correspondiente de cada elemento.
func CalcularMovimientosContables(elementos []*models.Elemento, dsc string, movId, terceroCr, terceroDb int, cuentas map[string]models.CuentaContable, movimientos *[]*models.MovimientoTransaccion) (
	errMsg string, outputError map[string]interface{}) {

	subgrupos := make(map[int]models.DetalleSubgrupo)
	tiposBien := make(map[int]models.TipoBien)
	cuentasSgTb := make(map[int]map[int]models.CuentaSubgrupo)
	totalesCr := make(map[string]float64)
	totalesDb := make(map[string]float64)

	var parCr int
	var parDb int
	var uvt float64
	var payload = "limit=1&fields=TipoBienId&sortby=Id&order=desc&query=Activo:true,SubgrupoId__Id:"

	if cuentas == nil {
		cuentas = make(map[string]models.CuentaContable)
	}

	if uvt_, err := parametros.GetUVTByVigencia(time.Now().Year()); err != nil {
		return "", err
	} else if uvt_ == 0 {
		return "No se pudo consultar el valor del UVT. Intente más tarde o contacte soporte.", nil
	} else {
		uvt = uvt_
	}

	if db_, cr_, err := parametros.GetParametrosDebitoCredito(); err != nil {
		return "", err
	} else {
		parCr = cr_
		parDb = db_
	}

	for _, el := range elementos {

		if el.SubgrupoCatalogoId <= 0 {
			return "No se pudo determinar la clase de los elementos. Revise el detalle del acta de recibido o contacte soporte.", nil
		}

		if el.TipoBienId == 0 {
			if _, ok := subgrupos[el.SubgrupoCatalogoId]; !ok {
				if sg, err := catalogoElementos.GetAllDetalleSubgrupo(payload + strconv.Itoa(el.SubgrupoCatalogoId)); err != nil {
					return "", err
				} else if len(sg) == 1 {
					subgrupos[el.SubgrupoCatalogoId] = *sg[0]
				} else {
					return "No se pudo consultar la parametrización de las clases. Contacte soporte.", nil
				}
			}

			if tb, err := catalogoElementos.GetTipoBienIdByValor(subgrupos[el.SubgrupoCatalogoId].TipoBienId.Id, int(el.ValorUnitario/uvt), tiposBien); err != nil {
				return "", err
			} else if tb == 0 {
				return "No se pudo establecer el tipo de bien de los elementos. Contacte soporte.", nil
			} else {
				el.TipoBienId = tb
			}
		}

		if _, ok := cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId]; !ok {
			if cst, err := catalogoElementos.GetAllCuentasSubgrupo(payloadCuentas(el.SubgrupoCatalogoId, el.TipoBienId, movId)); err != nil {
				return "", err
			} else if len(cst) == 1 {
				if cuentasSgTb[el.SubgrupoCatalogoId] == nil {
					cuentasSgTb[el.SubgrupoCatalogoId] = make(map[int]models.CuentaSubgrupo)
				}
				cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId] = *cst[0]
			} else {
				return "No se pudo establecer la parametrización contable.", nil
			}
		}

		if _, ok := cuentas[cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId].CuentaCreditoId]; !ok {
			if cr, err := cuentasContables.GetCuentaContable(cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId].CuentaCreditoId); err != nil {
				return "", err
			} else if cr != nil {
				cuentas[cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId].CuentaCreditoId] = *cr
			} else {
				return "No se pudo encontrar la cuenta contable. Contacte soporte", nil
			}
		}

		if _, ok := cuentas[cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId].CuentaDebitoId]; !ok {
			if db, err := cuentasContables.GetCuentaContable(cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId].CuentaDebitoId); err != nil {
				return "", err
			} else if db != nil {
				cuentas[cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId].CuentaDebitoId] = *db
			} else {
				return "No se pudo encontrar la cuenta contable. Contacte soporte", nil
			}

		}

		totalesCr[cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId].CuentaCreditoId] += el.ValorTotal
		totalesDb[cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId].CuentaDebitoId] += el.ValorTotal

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

	return
}

func payloadCuentas(sg, tb, mov int) string {
	return "fields=CuentaDebitoId,CuentaCreditoId&query=Activo:true,SubgrupoId__Id:" +
		strconv.Itoa(sg) + ",TipoBienId__Id:" + strconv.Itoa(tb) + ",SubtipoMovimientoId:" + strconv.Itoa(mov)
}
