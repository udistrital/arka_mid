package catalogoElementosHelper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/cuentasContablesHelper"
	"github.com/udistrital/arka_mid/helpers/movimientosArkaHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

func GetInfoSubgrupo(subgrupoId int) (detalleSubgrupo map[string]interface{}, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetInfoSubgrupo - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	if subgrupoId <= 0 {
		err := fmt.Errorf("subgrupoId MUST be > 0 - Got: %d", subgrupoId)
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetInfoSubgrupo - subgrupoId <= 0",
			"err":     err,
			"status":  "400",
		}
		panic(outputError)
	}

	var detalles []map[string]interface{}

	urlcrud := "http://" + beego.AppConfig.String("catalogoElementosService") + "subgrupo?limit=-1"
	urlcrud += "&query=Activo:True&Id:" + strconv.Itoa(int(subgrupoId))

	if response, err := request.GetJsonTest(urlcrud, &detalles); err == nil && response.StatusCode == 200 { // (2) error servicio caido
		// fmt.Println(cuentasSubgrupo[0])
		if detalles[0]["Id"].(int) != 0 {
			return detalles[0], nil
		} else {
			err = fmt.Errorf("Cuenta no existe")
			outputError = map[string]interface{}{
				"funcion": "GetInfoSubgrupo - request.GetJsonTest(urlcrud, &cuentasSubgrupo)",
				"err":     err,
				"status:": "502",
			}
			return nil, outputError
		}
	} else {
		if err == nil {
			err = fmt.Errorf("Undesired Status Code: %d", response.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetInfoSubgrupo - request.GetJsonTest(urlcrud, &cuentasSubgrupo)",
			"err":     err,
			"status:": "502",
		}
		return nil, outputError
	}
}

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
			if strings.Contains(fm.CodigoAbreviacion, "ENT_") || fm.CodigoAbreviacion == "SAL" || fm.CodigoAbreviacion == "BJ_HT" {
				movs = append(movs, fm)
			} else if detalle.Depreciacion && fm.CodigoAbreviacion == "DEP" {
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

func GetDetalleSubgrupo(subgrupoId int) (subgrupo []*models.DetalleSubgrupo, outputError map[string]interface{}) {
	if subgrupoId > 0 {

		defer func() {
			if err := recover(); err != nil {
				outputError = map[string]interface{}{
					"funcion": "/GetDetalleSubgrupo - Unhandled Error!",
					"err":     err,
					"status":  "500",
				}
				panic(outputError)
			}
		}()

		urlSubgrupo := "http://" + beego.AppConfig.String("catalogoElementosService") + "detalle_subgrupo?"
		urlSubgrupo += "query=Activo:true,SubgrupoId__Id:" + strconv.Itoa(subgrupoId)
		urlSubgrupo += "&fields=SubgrupoId,TipoBienId,Depreciacion,ValorResidual,VidaUtil&sortby=Id&order=desc"
		if response, err := request.GetJsonTest(urlSubgrupo, &subgrupo); err == nil {
			if response.StatusCode == 200 {
				return subgrupo, nil
			} else {
				err := fmt.Errorf("Undesired Status: %s", response.Status)
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "GetDetalleSubgrupo - request.GetJsonTest(urlSubgrupo, &subgrupo) / response.StatusCode == 200",
					"err":     err,
					"status":  "502",
				}
				return nil, outputError
			}
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetUnidad - request.GetJsonTest(urlSubgrupo, &subgrupo)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
	} else {
		err := fmt.Errorf("subgrupoId MUST be greater than 0")
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetDetalleSubgrupo",
			"err":     err,
			"status":  "400",
		}
		return nil, outputError
	}
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
