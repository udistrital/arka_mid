package ajustesHelper

import (
	"time"

	"github.com/udistrital/arka_mid/helpers/depreciacionHelper"
	"github.com/udistrital/arka_mid/helpers/parametrosHelper"
	"github.com/udistrital/arka_mid/helpers/tercerosHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

const queryUD string = "query=TipoDocumentoId__Nombre:NIT,Numero:"

// calcularAjusteMediciones Vuelve a generar las novedades y calcula las transacciones contables según las modificaciones que hayan afectado mediciones posteriores aprobadas
func calcularAjusteMediciones(novedades map[int][]*models.NovedadElemento,
	sg, vls, mp []*models.DetalleElemento_,
	org []*models.Elemento) (movimientos []*models.MovimientoTransaccion,
	novedades_ []*models.NovedadElemento, outputError map[string]interface{}) {

	funcion := "calcularAjusteMediciones"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		cuentasSubgrupo map[int]*models.CuentaSubgrupo
		bufferCtas      map[string]*models.CuentaContable
		movDebito       int
		movCredito      int
		terceroUD       int
	)

	detalleMd := make(map[int]*models.FormatoDepreciacion)
	novedadesNuevas := make(map[int][]*models.NovedadElemento)

	if db_, cr_, err := parametrosHelper.GetParametrosDebitoCredito(); err != nil {
		return nil, nil, err
	} else {
		movDebito = db_
		movCredito = cr_
	}

	if terceroUD_, err := tercerosHelper.GetAllDatosIdentificacion(queryUD + tercerosHelper.GetDocUD()); err != nil {
		return nil, nil, err
	} else {
		terceroUD = terceroUD_[0].TerceroId.Id
	}

	if cuentasSg, cuentas, err := consultaCuentasMp(novedades, sg, vls, mp, org); err != nil {
		return nil, nil, err
	} else {
		cuentasSubgrupo = cuentasSg
		bufferCtas = cuentas
	}

	for key, nv := range novedades {

		var nuevo *models.DetalleElemento_
		var sgOrg int
		if idx := findElementoInArrayD(sg, nv[0].ElementoMovimientoId.ElementoActaId); idx > -1 {
			nuevo = sg[idx]
			if idx := findElementoInArrayE(org, nv[0].ElementoMovimientoId.ElementoActaId); idx > -1 {
				sgOrg = org[idx].SubgrupoCatalogoId
			}
		}

		if idx := findElementoInArrayD(vls, nv[0].ElementoMovimientoId.ElementoActaId); idx > -1 {
			nuevo = vls[idx]
		}

		if idx := findElementoInArrayD(mp, nv[0].ElementoMovimientoId.ElementoActaId); idx > -1 {
			nuevo = mp[idx]
		}

		for idx, nv_ := range nv {

			if detalleMd[nv_.MovimientoId.Id] == nil {
				if dt, err := depreciacionHelper.GetDetalleDepreciacion(nv_.MovimientoId.Detalle); err != nil {
					return nil, nil, err
				} else {
					detalleMd[nv_.MovimientoId.Id] = dt
				}
			}

			var fCorte time.Time
			var dpOrg, dpNvo, deltaT float64
			var novedadNueva *models.NovedadElemento
			fCorte, _ = time.Parse("2006-01-02", detalleMd[nv_.MovimientoId.Id].FechaCorte)
			if idx == 0 {
				dpOrg, _ = depreciacionHelper.CalculaDp(
					nv_.ElementoMovimientoId.ValorTotal,
					nv_.ElementoMovimientoId.ValorResidual,
					nv_.ElementoMovimientoId.VidaUtil,
					nv_.ElementoMovimientoId.MovimientoId.FechaModificacion,
					fCorte)
				dpNvo, deltaT = depreciacionHelper.CalculaDp(
					nuevo.ValorTotal,
					nuevo.ValorResidual,
					nuevo.VidaUtil,
					nv_.ElementoMovimientoId.MovimientoId.FechaModificacion,
					fCorte)
				novedadNueva = generarNovedad(
					nuevo.ValorTotal-dpNvo,
					nuevo.ValorResidual,
					nuevo.VidaUtil-deltaT,
					nv_)
			} else {
				ref, _ := time.Parse("2006-01-02", detalleMd[novedadesNuevas[key][idx-1].MovimientoId.Id].FechaCorte)
				dpOrg, _ = depreciacionHelper.CalculaDp(
					nv_.ValorLibros,
					nv_.ValorResidual,
					nv_.VidaUtil,
					ref.AddDate(0, 0, 1),
					fCorte)
				dpNvo, deltaT = depreciacionHelper.CalculaDp(
					novedadesNuevas[key][idx-1].ValorLibros,
					nuevo.ValorResidual,
					novedadesNuevas[key][idx-1].VidaUtil,
					ref.AddDate(0, 0, 1),
					fCorte)
				novedadNueva = generarNovedad(
					novedadesNuevas[key][idx-1].ValorLibros-dpNvo,
					nuevo.ValorResidual,
					novedadesNuevas[key][idx-1].VidaUtil-deltaT,
					nv_)
			}

			novedadesNuevas[key] = append(novedadesNuevas[key], novedadNueva)

			movimientos = append(movimientos,
				generaTrContable(dpOrg, dpNvo,
					detalleMd[nv_.MovimientoId.Id].FechaCorte,
					nv_.MovimientoId.FormatoTipoMovimientoId.Nombre,
					movDebito,
					movCredito,
					sgOrg,
					nuevo.SubgrupoCatalogoId,
					terceroUD,
					cuentasSubgrupo, bufferCtas)...)
		}
	}

	for _, nv := range novedadesNuevas {
		for _, nv_ := range nv {
			novedades_ = append(novedades_, nv_)
		}
	}

	return movimientos, novedades_, nil

}

// consultaCuentasMp Consulta las cuentas asignadas a cada subgrupo y su detalle según el tipo de novedad
func consultaCuentasMp(novedades map[int][]*models.NovedadElemento,
	sg, vls, mp []*models.DetalleElemento_,
	org []*models.Elemento) (
	ctasSg map[int]*models.CuentaSubgrupo,
	ctas map[string]*models.CuentaContable,
	outputError map[string]interface{}) {

	funcion := "consultaCuentasMp"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var idsD, idsA []int
	var idD, idA int
	var ctasD, ctasA map[int]*models.CuentaSubgrupo

	for _, nv := range novedades {
		var ids []int
		if len(sg) > 0 {
			if idx := findElementoInArrayD(sg, nv[0].ElementoMovimientoId.ElementoActaId); idx > -1 {
				ids = append(ids, sg[idx].SubgrupoCatalogoId)
				if idx := findElementoInArrayE(org, nv[0].ElementoMovimientoId.ElementoActaId); idx > -1 {
					ids = append(ids, org[idx].SubgrupoCatalogoId)
				}
			}
		} else if len(vls) > 0 {
			if idx := findElementoInArrayD(vls, nv[0].ElementoMovimientoId.ElementoActaId); idx > -1 {
				ids = append(ids, vls[idx].SubgrupoCatalogoId)
			}
		} else if len(mp) > 0 {
			if idx := findElementoInArrayD(mp, nv[0].ElementoMovimientoId.ElementoActaId); idx > -1 {
				ids = append(ids, mp[idx].SubgrupoCatalogoId)
			}
		}

		if nv[0].MovimientoId.FormatoTipoMovimientoId.Nombre == "Depreciación" {
			if idD == 0 {
				idD = nv[0].MovimientoId.FormatoTipoMovimientoId.Id
			}
			idsD = append(idsD, ids...)
		} else if nv[0].MovimientoId.FormatoTipoMovimientoId.Nombre == "Amortizacion" {
			if idA == 0 {
				idA = nv[0].MovimientoId.FormatoTipoMovimientoId.Id
			}
			idsA = append(idsA, ids...)
		}

	}

	if idD > 0 {
		if ctas, err := getCuentasByMovimientoSubgrupos(idD, idsD); err != nil {
			return nil, nil, err
		} else {
			ctasD = ctas
		}
	}

	if idA > 0 {
		if ctas, err := getCuentasByMovimientoSubgrupos(idA, idsA); err != nil {
			return nil, nil, err
		} else {
			ctasA = ctas
		}
	}

	ctasSg = joinMaps(ctasD, ctasA)

	idsCtas := make([]string, 0)
	for _, ctas := range ctasSg {
		idsCtas = append(idsCtas, ctas.CuentaDebitoId, ctas.CuentaCreditoId)
	}

	ctas = make(map[string]*models.CuentaContable)
	if detalleCuenta_, err := fillCuentas(ctas, idsCtas); err != nil {
		return nil, nil, err
	} else {
		ctas = detalleCuenta_
	}

	return ctasSg, ctas, nil

}

// separarNovedadesPorElemento Separa las novedades por elementos
func separarNovedadesPorElemento(novedades []*models.NovedadElemento) (novedades_ map[int][]*models.NovedadElemento) {

	novedades_ = make(map[int][]*models.NovedadElemento, 0)
	for _, nv := range novedades {
		novedades_[nv.ElementoMovimientoId.Id] = append(novedades_[nv.ElementoMovimientoId.Id], nv)
	}

	return novedades_

}

// generarNovedad Actualiza los valores afectados de una novedad al hacer un ajuste a un elemento
func generarNovedad(libros, residual, util float64, novedad *models.NovedadElemento) *models.NovedadElemento {

	novedad.ValorLibros = libros
	novedad.VidaUtil = util
	novedad.ValorResidual = residual

	return novedad

}
