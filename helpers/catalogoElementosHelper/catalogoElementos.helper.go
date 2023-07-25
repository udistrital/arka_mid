package catalogoElementosHelper

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// GetCuentasContablesSubgrupo ...
func GetCuentasContablesSubgrupo(subgrupoId, movimientoId int, cuentas *[]models.DetalleCuentasSubgrupo) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GetCuentasContablesSubgrupo - Unhandled Error!", "500")

	var (
		query         string
		ctas          []models.CuentasSubgrupo
		movs          []*models.FormatoTipoMovimiento
		tiposBien     []models.TipoBien
		formatoSalida models.FormatoTipoMovimiento
	)

	query = "limit=1&sortby=FechaCreacion&order=desc&fields=Id,Depreciacion,Amortizacion,TipoBienId" +
		"&query=Activo:true,SubgrupoId__Id:" + strconv.Itoa(subgrupoId)
	detalle, outputError := catalogoElementos.GetAllDetalleSubgrupo(query)
	if outputError != nil || len(detalle) != 1 {
		return
	}

	query = "limit=-1&sortby=LimiteSuperior&order=asc&fields=Id,Nombre,BodegaConsumo" +
		"&query=Activo:true,TipoBienPadreId__Id:" + strconv.Itoa(detalle[0].TipoBienId.Id)
	outputError = catalogoElementos.GetAllTipoBien(query, &tiposBien)
	if outputError != nil || len(tiposBien) == 0 {
		return
	}

	query = "limit=-1&sortby=CodigoAbreviacion&order=asc&query=Activo:true,CodigoAbreviacion:SAL&fields=Id,CodigoAbreviacion,Nombre"
	tipos, outputError := movimientosArka.GetAllFormatoTipoMovimiento(query)
	if outputError != nil || len(tipos) != 1 {
		return
	} else {
		formatoSalida = *tipos[0]
	}

	query = "limit=-1&sortby=CodigoAbreviacion&order=asc&fields=Id,CodigoAbreviacion,Nombre&query=Activo:true"
	if movimientoId > 0 {
		query += ",Id:" + strconv.Itoa(movimientoId)
	}

	tipos, outputError = movimientosArka.GetAllFormatoTipoMovimiento(query)
	if outputError != nil {
		return
	}

	for _, fm := range tipos {
		if (strings.Contains(fm.CodigoAbreviacion, "ENT_") || fm.CodigoAbreviacion == "BJ_HT") && !strings.Contains(fm.CodigoAbreviacion, "KDX") {
			movs = append(movs, fm)
		} else if fm.CodigoAbreviacion == "CRR" && (detalle[0].Depreciacion || detalle[0].Amortizacion) {
			movs = append(movs, fm)
		}
	}

	outputError = catalogoElementos.GetTrCuentasSubgrupo(subgrupoId, movimientoId, &ctas)
	if outputError != nil {
		return
	}

	detalleCtas := make(map[string]models.DetalleCuenta)
	for _, fm := range movs {
		for _, tb := range tiposBien {
			if tb.BodegaConsumo && (fm.CodigoAbreviacion == "CRR" || fm.CodigoAbreviacion == "BJ_HT") {
				continue
			}

			outputError = fillCuentaSubgrupo(subgrupoId, cuentas, tb, models.FormatoTipoMovimiento{Id: 0}, *fm, ctas, detalleCtas)
			if outputError != nil {
				return
			}

			if !strings.Contains(fm.CodigoAbreviacion, "ENT_") {
				continue
			}

			outputError = fillCuentaSubgrupo(subgrupoId, cuentas, tb, *fm, formatoSalida, ctas, detalleCtas)
			if outputError != nil {
				return
			}
		}
	}

	return
}

func fillCuentaSubgrupo(sgId int, cFinales *[]models.DetalleCuentasSubgrupo, tb models.TipoBien, mov, sMov models.FormatoTipoMovimiento,
	ctasSg []models.CuentasSubgrupo, cuentas map[string]models.DetalleCuenta) (
	outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("fillCuentaSubgrupo - Unhandled Error!", "500")

	var dCta models.DetalleCuentasSubgrupo
	dCta.SubgrupoId = sgId
	dCta.TipoMovimientoId = &mov
	dCta.SubtipoMovimientoId = &sMov
	dCta.TipoBienId = tb

	if idx := findInArray(ctasSg, mov.Id, sMov.Id, tb.Id); idx > -1 {
		dCta.Id = ctasSg[idx].Id
		dCta.CuentaCreditoId = new(models.DetalleCuenta)
		dCta.CuentaDebitoId = new(models.DetalleCuenta)

		err := findCuentaSubgrupo(dCta.CuentaCreditoId, ctasSg[idx].CuentaCreditoId, cuentas)
		if err != nil {
			return err
		}

		err = findCuentaSubgrupo(dCta.CuentaDebitoId, ctasSg[idx].CuentaDebitoId, cuentas)
		if err != nil {
			return err
		}
	}

	*cFinales = append(*cFinales, dCta)
	return
}

func findCuentaSubgrupo(ctaSg *models.DetalleCuenta, cuentaId string, cuentas map[string]models.DetalleCuenta) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("findCuentaSubgrupo - Unhandled Error!", "500")

	if val, ok := cuentas[cuentaId]; ok {
		*ctaSg = val
		return
	}

	cta, outputError := cuentasContables.GetCuentaContable(cuentaId)
	if outputError != nil {
		return
	} else if cta != nil {
		var dcta models.DetalleCuenta
		outputError = utilsHelper.FillStruct(cta, &dcta)
		if outputError != nil {
			return
		}

		*ctaSg = dcta
		cuentas[cuentaId] = dcta
	}

	return
}

// GetCuentasByMovimientoSubgrupos Consulta las cuentas para una serie de subgrupos y las almacena en una estructura de fácil acceso
func GetCuentasByMovimientoAndSubgrupos(movimientoId int, subgrupos []int, cuentasSubgrupo map[int]models.CuentasSubgrupo) (
	outputError map[string]interface{}) {

	funcion := "GetCuentasByMovimientoSubgrupos"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var subgrupos_ []int
	for _, sg := range subgrupos {
		if _, ok := cuentasSubgrupo[sg]; !ok {
			subgrupos_ = append(subgrupos_, sg)
		}
	}

	if len(subgrupos_) == 0 {
		return
	}

	query := "limit=-1&fields=CuentaDebitoId,CuentaCreditoId,SubgrupoId&sortby=Id&order=desc&" +
		"query=Activo:true,SubtipoMovimientoId:" + strconv.Itoa(movimientoId) +
		",SubgrupoId__Id__in:" + url.QueryEscape(utilsHelper.ArrayToString(subgrupos_, "|"))
	if cuentas_, err := catalogoElementos.GetAllCuentasSubgrupo(query); err != nil {
		return err
	} else {
		for _, cuenta := range cuentas_ {
			cuentasSubgrupo[cuenta.SubgrupoId.Id] = *cuenta
		}

	}

	return

}

// findIdInArray Retorna la posicion en que se encuentra el id específicado
func findInArray(cuentasSg []models.CuentasSubgrupo, movimientoId, sMovimientoId, tipoBienId int) (i int) {
	for i, cuentaSg := range cuentasSg {
		if cuentaSg.TipoMovimientoId == movimientoId && cuentaSg.SubtipoMovimientoId == sMovimientoId && cuentaSg.TipoBienId.Id == tipoBienId {
			return i
		}
	}
	return -1
}
