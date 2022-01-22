package catalogoElementosHelper

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/movimientosArkaHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

// GetCatalogoById ...
func GetCatalogoById(catalogoId int) (catalogo *[]map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetCatalogoById - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	if catalogoId <= 0 {
		err := errors.New("catalogoId MUST be > 0")
		logs.Error(err)
		return nil, map[string]interface{}{
			"funcion": "GetCatalogoById - catalogoId <= 0",
			"err":     err,
			"status":  "400",
		}
	}

	urlcrud := "http://" + beego.AppConfig.String("catalogoElementosService") + "tr_catalogo/" + strconv.Itoa(catalogoId)
	if response, err := request.GetJsonTest(urlcrud, &catalogo); err == nil && response.StatusCode == 200 { // (2) error servicio caido
		return catalogo, nil
	} else {
		if err == nil {
			err = fmt.Errorf("Undesired Status Code: %d", response.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetCatalogoById - request.GetJsonTest(urlcrud, &catalogo)",
			"err":     err,
			"status:": "502",
		}
		return nil, outputError
	}
}

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
		for _, fm := range movs {
			dCta := new(models.DetalleCuentasSubgrupo)
			dCta.SubtipoMovimientoId = fm
			dCta.SubgrupoId = subgrupoId
				cuentas = append(cuentas, dCta)
		}
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

//GetTipoMovimiento funcion para traer cuenta asociadas a subgrupos por lo tanto crea sus propias estructuras como subgrupoCuentasModelo
func GetTipoMovimiento(arreglosubgrupos []models.SubgrupoCuentasModelo) (subgrupos []models.SubgrupoCuentasMovimiento, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetTipoMovimiento - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var urlcatalogo, urlcuenta string
	var arreglocuentas []models.CuentasGrupoMovimiento
	var subgrupocatalogo models.SubgrupoCuentasModelo
	var cuentareal map[string]interface{}
	for _, subgrupocuentas := range arreglosubgrupos {
		urlcatalogo = "http://" + beego.AppConfig.String("catalogoElementosService") + "tr_cuentas_subgrupo/" + strconv.Itoa(subgrupocuentas.Id)
		if response, err := request.GetJsonTest(urlcatalogo, &subgrupocatalogo); err == nil && response.StatusCode == 200 {

			for _, cuenta := range subgrupocatalogo.CuentasAsociadas {

				cuentaaso := models.CuentasGrupoMovimiento{
					Id:                  cuenta.Id,
					FechaCreacion:       cuenta.FechaCreacion,
					FechaModificacion:   cuenta.FechaModificacion,
					Activo:              cuenta.Activo,
					SubgrupoId:          cuenta.SubgrupoId,
					SubtipoMovimientoId: cuenta.SubtipoMovimientoId,
				}

				urlcuenta = "http://" + beego.AppConfig.String("cuentasContablesService") + "cuenta_contable/" + strconv.Itoa(cuenta.CuentaCreditoId)
				if response, err := request.GetJsonTest(urlcuenta, &cuentareal); err == nil && response.StatusCode == 200 {
					cuentaaso.CuentaCreditoId = cuentareal["Codigo"].(string)
				} else if err != nil {
					logs.Error(err)
					outputError = map[string]interface{}{
						"funcion": "GetTipoMovimiento - request.GetJsonTest(urlcatalogo, &subgrupocatalogo)",
						"err":     err,
						"status":  "502",
					}
					return nil, outputError
				}
				arreglocuentas = append(arreglocuentas, cuentaaso)

			}

		} else if err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetTipoMovimiento - request.GetJsonTest(urlcatalogo, &subgrupocatalogo)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}

		subgrupos = append(subgrupos, models.SubgrupoCuentasMovimiento{
			Id:                subgrupocuentas.Id,
			Nombre:            subgrupocuentas.Nombre,
			Descripcion:       subgrupocuentas.Descripcion,
			FechaCreacion:     subgrupocuentas.FechaCreacion,
			FechaModificacion: subgrupocuentas.FechaModificacion,
			Activo:            subgrupocuentas.Activo,
			Codigo:            subgrupocuentas.Codigo,
			CuentasAsociadas:  arreglocuentas,
		})

	} //hasta aca va forr

	return subgrupos, nil
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
