package asientoContable

import (
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
)

// CalcularMovimientosContables Calcula los movimientos contables dados los valores y parametrización correspondiente de cada elemento.
func CalcularMovimientosContables(
	elementos []*models.Elemento,
	dsc string,
	movId, sMovId, terceroCr, terceroDb int,
	cuentas map[string]models.CuentaContable,
	subgrupos map[int]models.DetalleSubgrupo,
	movimientos *[]*models.MovimientoTransaccion,
) (errMsg string, outputError map[string]interface{}) {

	logs.Info("==== INICIO CalcularMovimientosContables ====")
	logs.Info("CalcularMovimientosContables -> len(elementos)=%d dsc=%q movId=%d sMovId=%d terceroCr=%d terceroDb=%d",
		len(elementos), dsc, movId, sMovId, terceroCr, terceroDb)

	if movimientos == nil {
		logs.Error("CalcularMovimientosContables -> movimientos nil")
		return "No se recibió el contenedor de movimientos contables. Contacte soporte.", nil
	}

	if subgrupos == nil {
		subgrupos = make(map[int]models.DetalleSubgrupo)
		logs.Info("CalcularMovimientosContables -> subgrupos inicializado")
	}

	if cuentas == nil {
		cuentas = make(map[string]models.CuentaContable)
		logs.Info("CalcularMovimientosContables -> cuentas inicializado")
	}

	var parCr int
	var parDb int
	var uvt float64 = 1

	payloadDetalleSubgrupo := "limit=1&fields=TipoBienId,Amortizacion,Depreciacion,SubgrupoId&sortby=Id&order=desc&query=Activo:true,SubgrupoId__Id:"

	// if uvt_, err := parametros.GetUVTByVigencia(time.Now().Year()); err != nil {
	// 	return "", err
	// } else if uvt_ == 0 {
	// 	return "No se pudo consultar el valor del UVT. Intente más tarde o contacte soporte.", nil
	// } else {
	// 	uvt = uvt_
	// }

	if db_, cr_, err := parametros.GetParametrosDebitoCredito(); err != nil {
		logs.Error("CalcularMovimientosContables -> GetParametrosDebitoCredito err=%v", err)
		return "", err
	} else {
		parCr = cr_
		parDb = db_
		logs.Info("CalcularMovimientosContables -> parametros contables parCr=%d parDb=%d", parCr, parDb)
	}

	tiposBien := make(map[int]models.TipoBien)
	cuentasSgTb := make(map[int]map[int]models.CuentasSubgrupo)
	totalesCr := make(map[string]float64)
	totalesDb := make(map[string]float64)
	var actasConflicto []int
	var subgruposConflicto []string

	for i, el := range elementos {
		if el == nil {
			logs.Error("CalcularMovimientosContables -> elemento nil en posición=%d", i)
			return "Se encontró un elemento nulo en la entrada. Contacte soporte.", nil
		}

		logs.Info("CalcularMovimientosContables -> elemento %d/%d: %+v", i+1, len(elementos), *el)

		if el.ValorTotal == 0 {
			logs.Info("CalcularMovimientosContables -> se omite elemento Id=%d por ValorTotal=0", el.Id)
			continue
		}

		if el.SubgrupoCatalogoId <= 0 {
			logs.Error("CalcularMovimientosContables -> SubgrupoCatalogoId inválido. elemento=%+v", *el)
			return "No se pudo determinar la clase de los elementos. Revise el detalle del acta de recibido o contacte soporte.", nil
		}

		if _, ok := subgrupos[el.SubgrupoCatalogoId]; !ok {
			logs.Info("CalcularMovimientosContables -> consultando detalle subgrupo para SubgrupoCatalogoId=%d", el.SubgrupoCatalogoId)

			sg, outputError := catalogoElementos.GetAllDetalleSubgrupo(payloadDetalleSubgrupo + strconv.Itoa(el.SubgrupoCatalogoId))
			logs.Info("CalcularMovimientosContables -> GetAllDetalleSubgrupo outputError=%v len=%d", outputError, len(sg))
			if outputError != nil {
				return "", outputError
			}
			if len(sg) == 0 {
				logs.Error("CalcularMovimientosContables -> sin detalle subgrupo para SubgrupoCatalogoId=%d", el.SubgrupoCatalogoId)
				return "No se pudo consultar la parametrización de las clases. Contacte soporte.", nil
			}

			subgrupos[el.SubgrupoCatalogoId] = *sg[0]
			logs.Info("CalcularMovimientosContables -> detalle subgrupo cargado=%+v", subgrupos[el.SubgrupoCatalogoId])
		}

		detalleSg := subgrupos[el.SubgrupoCatalogoId]

		if detalleSg.TipoBienId == nil {
			logs.Error("CalcularMovimientosContables -> TipoBienId nil en detalle subgrupo. SubgrupoCatalogoId=%d detalle=%+v",
				el.SubgrupoCatalogoId, detalleSg)
			return "La clase del elemento no tiene TipoBien parametrizado. Contacte soporte.", nil
		}

		if detalleSg.SubgrupoId == nil {
			logs.Error("CalcularMovimientosContables -> SubgrupoId nil en detalle subgrupo. SubgrupoCatalogoId=%d detalle=%+v",
				el.SubgrupoCatalogoId, detalleSg)
			return "La clase del elemento no tiene Subgrupo parametrizado. Contacte soporte.", nil
		}

		if el.TipoBienId == 0 {
			logs.Info("CalcularMovimientosContables -> TipoBienId=0, calculando por valor. tipoBienPadre=%d valorUnitario=%v uvt=%v",
				detalleSg.TipoBienId.Id, el.ValorUnitario, uvt)

			tb, outputError := catalogoElementos.GetTipoBienIdByValor(detalleSg.TipoBienId.Id, el.ValorUnitario/uvt, tiposBien)
			logs.Info("CalcularMovimientosContables -> GetTipoBienIdByValor tb=%d outputError=%v", tb, outputError)
			if outputError != nil {
				return "", outputError
			}
			if tb == 0 {
				logs.Error("CalcularMovimientosContables -> no se pudo determinar tipo de bien por valor. elemento=%+v", *el)
				return "No se pudo establecer el tipo de bien de los elementos. Revise la parametrización de los tipos de bien.", nil
			}

			el.TipoBienId = tb
			logs.Info("CalcularMovimientosContables -> TipoBienId asignado automáticamente=%d", el.TipoBienId)

		} else {
			if _, ok := tiposBien[el.TipoBienId]; !ok {
				var tipoBien models.TipoBien
				outputError = catalogoElementos.GetTipoBienById(el.TipoBienId, &tipoBien)
				logs.Info("CalcularMovimientosContables -> GetTipoBienById TipoBienId=%d outputError=%v tipoBien=%+v",
					el.TipoBienId, outputError, tipoBien)
				if outputError != nil {
					return "", outputError
				}
				tiposBien[el.TipoBienId] = tipoBien
			}

			if tiposBien[el.TipoBienId].TipoBienPadreId == nil {
				logs.Error("CalcularMovimientosContables -> TipoBienPadreId nil. TipoBienId=%d tipoBien=%+v",
					el.TipoBienId, tiposBien[el.TipoBienId])
				return "El tipo de bien no tiene tipo padre parametrizado. Contacte soporte.", nil
			}

			if tiposBien[el.TipoBienId].TipoBienPadreId.Id != detalleSg.TipoBienId.Id {
				logs.Error("CalcularMovimientosContables -> conflicto tipo bien. TipoBienId=%d TipoBienPadreId=%d esperado=%d",
					el.TipoBienId, tiposBien[el.TipoBienId].TipoBienPadreId.Id, detalleSg.TipoBienId.Id)

				var elemento models.Elemento
				outputError = actaRecibido.GetElementoById(el.Id, &elemento)
				logs.Info("CalcularMovimientosContables -> GetElementoById Id=%d outputError=%v elemento=%+v",
					el.Id, outputError, elemento)
				if outputError != nil {
					return "", outputError
				}

				if elemento.ActaRecibidoId == nil {
					logs.Error("CalcularMovimientosContables -> ActaRecibidoId nil en elemento consultado. elemento=%+v", elemento)
					return "No se pudo validar el acta asociada al elemento. Contacte soporte.", nil
				}

				exists := containsInt(actasConflicto, elemento.ActaRecibidoId.Id)
				if !exists {
					actasConflicto = append(actasConflicto, elemento.ActaRecibidoId.Id)
					logs.Info("CalcularMovimientosContables -> acta agregada a conflicto=%d", elemento.ActaRecibidoId.Id)
				}
				continue
			}
		}

		if cuentasSgTb[el.SubgrupoCatalogoId] == nil {
			cuentasSgTb[el.SubgrupoCatalogoId] = make(map[int]models.CuentasSubgrupo)
			logs.Info("CalcularMovimientosContables -> mapa cuentasSgTb inicializado para subgrupo=%d", el.SubgrupoCatalogoId)
		}

		if _, ok := cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId]; !ok {
			payload := payloadCuentas(el.SubgrupoCatalogoId, el.TipoBienId, movId, sMovId)
			logs.Info("CalcularMovimientosContables -> consultando cuentas subgrupo payload=%s", payload)

			cst, outputError := catalogoElementos.GetAllCuentasSubgrupo(payload)
			logs.Info("CalcularMovimientosContables -> GetAllCuentasSubgrupo outputError=%v len=%d", outputError, len(cst))
			if outputError != nil {
				return "", outputError
			}
			if len(cst) == 1 {
				cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId] = *cst[0]
				logs.Info("CalcularMovimientosContables -> cuentas subgrupo-tipo bien=%+v", cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId])
			} else {
				exists := containsString(subgruposConflicto, detalleSg.SubgrupoId.Codigo)
				if !exists {
					subgruposConflicto = append(subgruposConflicto, detalleSg.SubgrupoId.Codigo)
					logs.Info("CalcularMovimientosContables -> subgrupo en conflicto agregado=%s", detalleSg.SubgrupoId.Codigo)
				}
				continue
			}
		}

		cuentaCfg := cuentasSgTb[el.SubgrupoCatalogoId][el.TipoBienId]

		if cuentaCfg.CuentaCreditoId == "" {
			logs.Error("CalcularMovimientosContables -> CuentaCreditoId vacío. subgrupo=%d tipoBien=%d config=%+v",
				el.SubgrupoCatalogoId, el.TipoBienId, cuentaCfg)
			return "La parametrización contable no tiene cuenta crédito. Contacte soporte.", nil
		}

		if cuentaCfg.CuentaDebitoId == "" {
			logs.Error("CalcularMovimientosContables -> CuentaDebitoId vacío. subgrupo=%d tipoBien=%d config=%+v",
				el.SubgrupoCatalogoId, el.TipoBienId, cuentaCfg)
			return "La parametrización contable no tiene cuenta débito. Contacte soporte.", nil
		}

		if _, ok := cuentas[cuentaCfg.CuentaCreditoId]; !ok {
			cr, outputError := cuentasContables.GetCuentaContable(cuentaCfg.CuentaCreditoId)
			logs.Info("CalcularMovimientosContables -> GetCuentaContable crédito=%s outputError=%v cuenta=%+v",
				cuentaCfg.CuentaCreditoId, outputError, cr)
			if outputError != nil {
				return "", outputError
			}
			if cr != nil {
				cuentas[cuentaCfg.CuentaCreditoId] = *cr
			} else {
				logs.Error("CalcularMovimientosContables -> cuenta crédito no encontrada=%s", cuentaCfg.CuentaCreditoId)
				return "No se pudo encontrar la cuenta contable. Contacte soporte", nil
			}
		}

		if _, ok := cuentas[cuentaCfg.CuentaDebitoId]; !ok {
			db, outputError := cuentasContables.GetCuentaContable(cuentaCfg.CuentaDebitoId)
			logs.Info("CalcularMovimientosContables -> GetCuentaContable débito=%s outputError=%v cuenta=%+v",
				cuentaCfg.CuentaDebitoId, outputError, db)
			if outputError != nil {
				return "", outputError
			}
			if db != nil {
				cuentas[cuentaCfg.CuentaDebitoId] = *db
			} else {
				logs.Error("CalcularMovimientosContables -> cuenta débito no encontrada=%s", cuentaCfg.CuentaDebitoId)
				return "No se pudo encontrar la cuenta contable. Contacte soporte", nil
			}
		}

		totalesCr[cuentaCfg.CuentaCreditoId] += el.ValorTotal
		totalesDb[cuentaCfg.CuentaDebitoId] += el.ValorTotal

		logs.Info("CalcularMovimientosContables -> acumulados crédito=%+v", totalesCr)
		logs.Info("CalcularMovimientosContables -> acumulados débito=%+v", totalesDb)
	}

	if len(actasConflicto) > 0 {
		logs.Error("CalcularMovimientosContables -> actasConflicto=%v", actasConflicto)
		return "El tipo bien asignado manualmente no corresponde a la clase correspondiente. Revise las siguientes actas: " + utilsHelper.ArrayToString(actasConflicto, ", "), nil
	}

	if len(subgruposConflicto) > 0 {
		logs.Error("CalcularMovimientosContables -> subgruposConflicto=%v", subgruposConflicto)
		return "No se pudo establecer la parametrización contable de las siguientes clases: " + strings.Join(subgruposConflicto, ", "), nil
	}

	for cta, val := range totalesCr {
		var movimiento models.MovimientoTransaccion
		fillMovimiento(val, dsc, terceroCr, parCr, cuentas[cta], &movimiento)
		*movimientos = append(*movimientos, &movimiento)
		logs.Info("CalcularMovimientosContables -> movimiento crédito agregado=%+v", movimiento)
	}

	for cta, val := range totalesDb {
		var movimiento models.MovimientoTransaccion
		fillMovimiento(val, dsc, terceroDb, parDb, cuentas[cta], &movimiento)
		*movimientos = append(*movimientos, &movimiento)
		logs.Info("CalcularMovimientosContables -> movimiento débito agregado=%+v", movimiento)
	}

	logs.Info("CalcularMovimientosContables -> total movimientos generados=%d", len(*movimientos))
	logs.Info("==== FIN CalcularMovimientosContables ====")
	return
}

func fillMovimiento(valor float64, dsc string, terceroId, tipoMov int, cuenta models.CuentaContable, movimiento *models.MovimientoTransaccion) {
	logs.Info("fillMovimiento -> valor=%v dsc=%q terceroId=%d tipoMov=%d cuenta=%+v", valor, dsc, terceroId, tipoMov, cuenta)

	if movimiento == nil {
		logs.Error("fillMovimiento -> movimiento nil")
		return
	}

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

	logs.Info("fillMovimiento -> movimiento construido=%+v", *movimiento)
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
