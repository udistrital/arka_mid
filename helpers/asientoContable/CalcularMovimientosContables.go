package asientoContable

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
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

	payloadDetalleSubgrupoBase := "limit=1&fields=TipoBienId,Amortizacion,Depreciacion,SubgrupoId&sortby=Id&order=desc&query=Activo:true,SubgrupoId__Id:"

	if db_, cr_, err := parametros.GetParametrosDebitoCredito(); err != nil {
		logs.Error("CalcularMovimientosContables -> GetParametrosDebitoCredito err=%v", err)
		return "", err
	} else {
		parCr = cr_
		parDb = db_
		logs.Info("CalcularMovimientosContables -> parametros contables parCr=%d parDb=%d", parCr, parDb)
	}

	tiposBien := make(map[int]models.TipoBien)
	cuentasSg := make(map[int]models.CuentasSubgrupo)
	totalesCr := make(map[string]float64)
	totalesDb := make(map[string]float64)

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
			payloadDetalleSubgrupo := payloadDetalleSubgrupoBase + strconv.Itoa(el.SubgrupoCatalogoId)
			logs.Info("CalcularMovimientosContables -> DEBUG helper=GetAllDetalleSubgrupo payload=%s", payloadDetalleSubgrupo)

			sg, outputError := catalogoElementos.GetAllDetalleSubgrupo(payloadDetalleSubgrupo)
			logs.Info("CalcularMovimientosContables -> DEBUG GetAllDetalleSubgrupo outputError=%v len=%d", outputError, len(sg))
			if outputError != nil {
				return "", outputError
			}
			if len(sg) == 0 {
				logs.Error("CalcularMovimientosContables -> sin detalle subgrupo para SubgrupoCatalogoId=%d", el.SubgrupoCatalogoId)
				return "No se pudo consultar la parametrización de las clases. Contacte soporte.", nil
			}

			subgrupos[el.SubgrupoCatalogoId] = *sg[0]
			logs.Info("CalcularMovimientosContables -> DEBUG detalleSubgrupoSeleccionado=%+v", subgrupos[el.SubgrupoCatalogoId])
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
			logs.Info("CalcularMovimientosContables -> DEBUG helper=GetTipoBienIdByValor tipoBienPadre=%d valorUnitario=%v uvt=%v",
				detalleSg.TipoBienId.Id, el.ValorUnitario, uvt)

			tb, outputError := catalogoElementos.GetTipoBienIdByValor(detalleSg.TipoBienId.Id, el.ValorUnitario/uvt, tiposBien)
			logs.Info("CalcularMovimientosContables -> DEBUG GetTipoBienIdByValor tb=%d outputError=%v", tb, outputError)
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
				logs.Info("CalcularMovimientosContables -> DEBUG helper=GetTipoBienById TipoBienId=%d", el.TipoBienId)

				outputError = catalogoElementos.GetTipoBienById(el.TipoBienId, &tipoBien)
				logs.Info("CalcularMovimientosContables -> DEBUG GetTipoBienById outputError=%v tipoBien=%+v",
					outputError, tipoBien)
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
				logs.Warn("CalcularMovimientosContables -> se omite validación de tipo de bien. TipoBienId=%d TipoBienPadreId=%d esperado=%d elemento=%+v",
					el.TipoBienId, tiposBien[el.TipoBienId].TipoBienPadreId.Id, detalleSg.TipoBienId.Id, *el)
			}
		}

		if _, ok := cuentasSg[el.SubgrupoCatalogoId]; !ok {
			payloadCtas := payloadCuentas(el.SubgrupoCatalogoId, movId, sMovId)
			logs.Info("CalcularMovimientosContables -> DEBUG criterios parametrizacion: SubgrupoCatalogoId=%d", el.SubgrupoCatalogoId)
			logs.Info("CalcularMovimientosContables -> DEBUG helper=GetAllCuentasSubgrupo payload=%s", payloadCtas)

			cst, outputError := catalogoElementos.GetAllCuentasSubgrupo(payloadCtas)
			logs.Info("CalcularMovimientosContables -> DEBUG GetAllCuentasSubgrupo outputError=%v len=%d", outputError, len(cst))
			if outputError != nil {
				return "", outputError
			}

			if len(cst) == 0 {
				logs.Error("CalcularMovimientosContables -> no hay parametrización contable para subgrupo=%d payload=%s",
					el.SubgrupoCatalogoId, payloadCtas)
				return mensajeFaltaParametrizacionSubgrupo(detalleSg), nil
			}

			for idx, cfg := range cst {
				if cfg != nil {
					logs.Info("CalcularMovimientosContables -> DEBUG cst[%d]=%+v", idx, *cfg)
				} else {
					logs.Warn("CalcularMovimientosContables -> DEBUG cst[%d]=nil", idx)
				}
			}

			if len(cst) > 1 {
				logs.Warn("CalcularMovimientosContables -> se encontraron %d parametrizaciones para subgrupo=%d. Se tomará la primera.",
					len(cst), el.SubgrupoCatalogoId)
			}

			cuentasSg[el.SubgrupoCatalogoId] = *cst[0]
			logs.Info("CalcularMovimientosContables -> DEBUG parametrizacionSeleccionada=%+v", cuentasSg[el.SubgrupoCatalogoId])
		}

		cuentaCfg := cuentasSg[el.SubgrupoCatalogoId]

		if cuentaCfg.CuentaCreditoId == "" {
			logs.Error("CalcularMovimientosContables -> CuentaCreditoId vacío. subgrupo=%d config=%+v",
				el.SubgrupoCatalogoId, cuentaCfg)
			return mensajeFaltaParametrizacionSubgrupo(detalleSg), nil
		}

		if cuentaCfg.CuentaDebitoId == "" {
			logs.Error("CalcularMovimientosContables -> CuentaDebitoId vacío. subgrupo=%d config=%+v",
				el.SubgrupoCatalogoId, cuentaCfg)
			return mensajeFaltaParametrizacionSubgrupo(detalleSg), nil
		}

		logs.Info("CalcularMovimientosContables -> DEBUG cuentaCreditoID=%s cuentaDebitoID=%s",
			cuentaCfg.CuentaCreditoId, cuentaCfg.CuentaDebitoId)

		if _, ok := cuentas[cuentaCfg.CuentaCreditoId]; !ok {
			logs.Info("CalcularMovimientosContables -> DEBUG helper=GetCuentaContable tipo=credito id=%s", cuentaCfg.CuentaCreditoId)

			cr, outputError := cuentasContables.GetCuentaContable(cuentaCfg.CuentaCreditoId)
			logs.Info("CalcularMovimientosContables -> DEBUG GetCuentaContable credito outputError=%v cuenta=%+v", outputError, cr)
			if outputError != nil {
				return "", outputError
			}
			if cr != nil {
				cuentas[cuentaCfg.CuentaCreditoId] = *cr
			} else {
				logs.Error("CalcularMovimientosContables -> cuenta crédito no encontrada=%s", cuentaCfg.CuentaCreditoId)
				return mensajeFaltaParametrizacionSubgrupo(detalleSg), nil
			}
		}

		if _, ok := cuentas[cuentaCfg.CuentaDebitoId]; !ok {
			logs.Info("CalcularMovimientosContables -> DEBUG helper=GetCuentaContable tipo=debito id=%s", cuentaCfg.CuentaDebitoId)

			db, outputError := cuentasContables.GetCuentaContable(cuentaCfg.CuentaDebitoId)
			logs.Info("CalcularMovimientosContables -> DEBUG GetCuentaContable debito outputError=%v cuenta=%+v", outputError, db)
			if outputError != nil {
				return "", outputError
			}
			if db != nil {
				cuentas[cuentaCfg.CuentaDebitoId] = *db
			} else {
				logs.Error("CalcularMovimientosContables -> cuenta débito no encontrada=%s", cuentaCfg.CuentaDebitoId)
				return mensajeFaltaParametrizacionSubgrupo(detalleSg), nil
			}
		}

		totalesCr[cuentaCfg.CuentaCreditoId] += el.ValorTotal
		totalesDb[cuentaCfg.CuentaDebitoId] += el.ValorTotal

		logs.Info("CalcularMovimientosContables -> acumulados crédito=%+v", totalesCr)
		logs.Info("CalcularMovimientosContables -> acumulados débito=%+v", totalesDb)
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

func payloadCuentas(sg, movId, sMovId int) string {
	payload := "limit=1&sortby=Id&order=desc&fields=CuentaDebitoId,CuentaCreditoId&query=Activo:true,SubgrupoId__Id:" +
		strconv.Itoa(sg)
	if sMovId > 0 {
		payload += ",SubtipoMovimientoId:" + strconv.Itoa(sMovId)
	} else if movId > 0 {
		payload += ",TipoMovimientoId:" + strconv.Itoa(movId)
	}

	logs.Info("payloadCuentas -> sg=%d movId=%d sMovId=%d payload=%s", sg, movId, sMovId, payload)
	return payload
}

func mensajeFaltaParametrizacionSubgrupo(detalleSg models.DetalleSubgrupo) string {
	if detalleSg.SubgrupoId != nil {
		return "Debe parametrizar las cuentas activas para el subgrupo " + detalleSg.SubgrupoId.Codigo + " " + detalleSg.SubgrupoId.Nombre
	}

	return "Debe parametrizar las cuentas activas para el subgrupo solicitado."
}
