package catalogoElementosHelper

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	crudCatalogo "github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	crudMovimientosArka "github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
)

// GetCuentasContablesSubgrupo ...
func GetCuentasContablesSubgrupo(subgrupoId int, cuentas *[]models.DetalleCuentasSubgrupo) (outputError map[string]interface{}) {

	funcion := "GetCuentasContablesSubgrupo - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	var (
		query     string
		ctas      []models.CuentasSubgrupo
		movs      []*models.FormatoTipoMovimiento
		detalle   models.DetalleSubgrupo
		tiposBien []models.TipoBien
	)

	query = "limit=1&sortby=FechaCreacion&order=desc&query=Activo:true,SubgrupoId__Id:" + strconv.Itoa(subgrupoId)
	query += "&fields=Id,Depreciacion,Amortizacion,TipoBienId"
	if detalle_, err := crudCatalogo.GetAllDetalleSubgrupo(query); err != nil {
		return err
	} else if len(detalle_) == 1 {
		detalle = *detalle_[0]
	} else {
		return
	}

	query = "limit=-1&sortby=LimiteSuperior&order=asc&query=Activo:true,TipoBienPadreId__Id:" + strconv.Itoa(detalle.TipoBienId.Id)
	query += "&fields=Id,Nombre"
	if err := catalogoElementos.GetAllTipoBien(query, &tiposBien); err != nil {
		return err
	} else if len(tiposBien) == 0 {
		return
	}

	query = "limit=-1&sortby=CodigoAbreviacion&order=asc&query=Activo:true"
	query += "&fields=Id,CodigoAbreviacion,Nombre"
	if movs_, err := crudMovimientosArka.GetAllFormatoTipoMovimiento(query); err != nil {
		return err
	} else {
		for _, fm := range movs_ {
			if (strings.Contains(fm.CodigoAbreviacion, "ENT_") || fm.CodigoAbreviacion == "SAL" || fm.CodigoAbreviacion == "BJ_HT") && !strings.Contains(fm.CodigoAbreviacion, "KDX") {
				movs = append(movs, fm)
			} else if fm.CodigoAbreviacion == "DEP" && detalle.Depreciacion {
				movs = append(movs, fm)
			} else if fm.CodigoAbreviacion == "AMT" && detalle.Amortizacion {
				movs = append(movs, fm)
			}
		}
	}

	if err := crudCatalogo.GetTrCuentasSubgrupo(subgrupoId, &ctas); err != nil {
		return err
	}

	if len(ctas) > 0 {
		detalleCtas := make(map[string]*models.DetalleCuenta)
		for _, fm := range movs {
			for _, tb := range tiposBien {
				var dCta models.DetalleCuentasSubgrupo
				dCta.SubtipoMovimientoId = fm
				dCta.SubgrupoId = subgrupoId
				dCta.TipoBienId = tb
				if idx := findInArray(ctas, fm.Id, tb.Id); idx > -1 {
					dCta.Id = ctas[idx].Id
					if val, ok := detalleCtas[ctas[idx].CuentaCreditoId]; ok {
						dCta.CuentaCreditoId = val
					} else {
						if cta, err := cuentasContables.GetCuentaContable(ctas[idx].CuentaCreditoId); err != nil {
							return err
						} else if cta != nil {
							var cdt *models.DetalleCuenta
							if err := formatdata.FillStruct(cta, &cdt); err != nil {
								logs.Error(err)
								eval := " - formatdata.FillStruct(cta, &cdt)"
								return errorctrl.Error(funcion+eval, err, "500")
							} else {
								dCta.CuentaCreditoId = cdt
								detalleCtas[ctas[idx].CuentaCreditoId] = cdt
							}
						}
					}

					if val, ok := detalleCtas[ctas[idx].CuentaDebitoId]; ok {
						dCta.CuentaDebitoId = val
					} else {
						if cta, err := cuentasContables.GetCuentaContable(ctas[idx].CuentaDebitoId); err != nil {
							return err
						} else if cta != nil {
							var dbt *models.DetalleCuenta
							if err := formatdata.FillStruct(cta, &dbt); err != nil {
								logs.Error(err)
								eval := " - formatdata.FillStruct(cta, &dbt)"
								return errorctrl.Error(funcion+eval, err, "500")
							} else {
								dCta.CuentaDebitoId = dbt
								detalleCtas[ctas[idx].CuentaDebitoId] = dbt
							}
						}
					}
				}
				*cuentas = append(*cuentas, dCta)
			}
		}
	} else {
		for _, fm := range movs {
			for _, tb := range tiposBien {
				var dCta models.DetalleCuentasSubgrupo
				dCta.SubtipoMovimientoId = fm
				dCta.SubgrupoId = subgrupoId
				dCta.TipoBienId = tb
				*cuentas = append(*cuentas, dCta)
			}
		}
	}

	return nil
}

// GetCuentasByMovimientoSubgrupos Consulta las cuentas para una serie de subgrupos y las almacena en una estructura de fácil acceso
func GetCuentasByMovimientoAndSubgrupos(movimientoId int, subgrupos []int, cuentasSubgrupo map[int]models.CuentaSubgrupo) (
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

	query := "limit=-1&fields=CuentaDebitoId,CuentaCreditoId,SubgrupoId&sortby=Id&order=desc&"
	query += "query=Activo:true,SubtipoMovimientoId:" + strconv.Itoa(movimientoId)
	query += ",SubgrupoId__Id__in:" + url.QueryEscape(utilsHelper.ArrayToString(subgrupos_, "|"))
	if cuentas_, err := crudCatalogo.GetAllCuentasSubgrupo(query); err != nil {
		return err
	} else {
		for _, cuenta := range cuentas_ {
			cuentasSubgrupo[cuenta.SubgrupoId.Id] = *cuenta
		}

	}

	return

}

// findIdInArray Retorna la posicion en que se encuentra el id específicado
func findInArray(cuentasSg []models.CuentasSubgrupo, movimientoId, tipoBienId int) (i int) {
	for i, cuentaSg := range cuentasSg {
		if cuentaSg.SubtipoMovimientoId == movimientoId && cuentaSg.TipoBienId.Id == tipoBienId {
			return i
		}
	}
	return -1
}
