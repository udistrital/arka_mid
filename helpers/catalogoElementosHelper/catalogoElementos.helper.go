package catalogoElementosHelper

import (
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/cuentasContablesHelper"
	"github.com/udistrital/arka_mid/helpers/movimientosArkaHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
)

// GetCuentasContablesSubgrupo ...
func GetCuentasContablesSubgrupo(subgrupoId int) (cuentas []*models.DetalleCuentasSubgrupo, outputError map[string]interface{}) {

	funcion := "GetCuentasContablesSubgrupo"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		query   string
		ctas    []*models.CuentasSubgrupo
		movs    []*models.FormatoTipoMovimiento
		detalle *models.DetalleSubgrupo
	)

	query = "limit=1&sortby=FechaCreacion&order=desc&query=Activo:true,SubgrupoId__Id:" + strconv.Itoa(subgrupoId)
	if detalle_, err := GetAllDetalleSubgrupo(query); err != nil {
		return nil, err
	} else if len(detalle_) == 0 {
		return nil, nil
	} else {
		detalle = detalle_[0]
	}

	query = "limit=-1&sortby=CodigoAbreviacion&order=asc&query=Activo:true"
	if movs_, err := movimientosArkaHelper.GetAllFormatoTipoMovimiento(query); err != nil {
		return nil, err
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

	if cuentas_, err := GetTrCuentasSubgrupo(subgrupoId); err != nil {
		return nil, err
	} else {
		ctas = cuentas_
	}

	if len(ctas) > 0 {
		detalleCtas := make(map[string]*models.DetalleCuenta)
		for _, fm := range movs {
			dCta := new(models.DetalleCuentasSubgrupo)
			dCta.SubtipoMovimientoId = fm
			dCta.SubgrupoId = subgrupoId
			if idx := FindInArray(ctas, fm.Id); idx > -1 {
				dCta.Id = ctas[idx].Id
				if val, ok := detalleCtas[ctas[idx].CuentaCreditoId]; ok {
					dCta.CuentaCreditoId = val
				} else {
					if cta, err := cuentasContablesHelper.GetCuentaContable(ctas[idx].CuentaCreditoId); err != nil {
						return nil, err
					} else if cta != nil {
						var cdt *models.DetalleCuenta
						if err := formatdata.FillStruct(cta, &cdt); err != nil {
							logs.Error(err)
							eval := " - formatdata.FillStruct(cta, &cdt)"
							return nil, errorctrl.Error(funcion+eval, err, "500")
						} else {
							dCta.CuentaCreditoId = cdt
							detalleCtas[ctas[idx].CuentaCreditoId] = cdt
						}
					}
				}

				if val, ok := detalleCtas[ctas[idx].CuentaDebitoId]; ok {
					dCta.CuentaDebitoId = val
				} else {
					if cta, err := cuentasContablesHelper.GetCuentaContable(ctas[idx].CuentaDebitoId); err != nil {
						return nil, err
					} else if cta != nil {
						var dbt *models.DetalleCuenta
						if err := formatdata.FillStruct(cta, &dbt); err != nil {
							logs.Error(err)
							eval := " - formatdata.FillStruct(cta, &dbt)"
							return nil, errorctrl.Error(funcion+eval, err, "500")
						} else {
							dCta.CuentaDebitoId = dbt
							detalleCtas[ctas[idx].CuentaDebitoId] = dbt
						}
					}
				}
			}
			cuentas = append(cuentas, dCta)
		}
		return cuentas, nil
	} else {
		for _, fm := range movs {
			dCta := new(models.DetalleCuentasSubgrupo)
			dCta.SubtipoMovimientoId = fm
			dCta.SubgrupoId = subgrupoId
			cuentas = append(cuentas, dCta)
		}
	}

	return cuentas, nil
}

// findIdInArray Retorna la posicion en que se encuentra el id específicado
func FindInArray(cuentasSg []*models.CuentasSubgrupo, movimientoId int) (i int) {
	for i, cuentaSg := range cuentasSg {
		if int(cuentaSg.SubtipoMovimientoId) == movimientoId {
			return i
		}
	}
	return -1
}
