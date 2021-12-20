package catalogoElementosHelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/cuentasContablesHelper"
	"github.com/udistrital/arka_mid/models"
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

// GetCuentasContablesGrupo ...
func GetCuentasContablesSubgrupo(subgrupoId int) (cuentasSubgrupoTransaccion []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetCuentasContablesSubgrupo - Unhandled Error!",
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
			"funcion": "GetCuentasContablesSubgrupo - subgrupoId <= 0",
			"err":     err,
			"status":  "400",
		}
	}

	var cuentasSubgrupo []*models.CuentaSubgrupo

	urlcrud := "http://" + beego.AppConfig.String("catalogoElementosService") + "cuentas_subgrupo?limit=-1"
	urlcrud += "&query=Activo:True,SubgrupoId.Id:" + strconv.Itoa(int(subgrupoId))
	// logs.Debug("urlcrud:", urlcrud)

	if response, err := request.GetJsonTest(urlcrud, &cuentasSubgrupo); err == nil && response.StatusCode == 200 { // (2) error servicio caido
		// fmt.Println(cuentasSubgrupo[0])
		if cuentasSubgrupo[0].Id != 0 {
			for _, cuenta := range cuentasSubgrupo {
				// logs.Debug("CuentaCreditoId:", cuenta.CuentaCreditoId, " - CuentaDebitoId:", cuenta.CuentaDebitoId)
				cuentaCredito, _ := cuentasContablesHelper.GetCuentaContable(cuenta.CuentaCreditoId)
				cuentaDebito, _ := cuentasContablesHelper.GetCuentaContable(cuenta.CuentaDebitoId)

				cuentasSubgrupoTransaccion = append(cuentasSubgrupoTransaccion, map[string]interface{}{
					"Id":                  cuenta.Id,
					"CuentaCreditoId":     cuentaCredito,
					"CuentaDebitoId":      cuentaDebito,
					"SubtipoMovimientoId": cuenta.SubtipoMovimientoId,
					"FechaCreacion":       cuenta.FechaCreacion,
					"FechaModificacion":   cuenta.FechaModificacion,
					"Activo":              cuenta.Activo,
					"SubgrupoId":          cuenta.SubgrupoId,
				})
			}
			return cuentasSubgrupoTransaccion, nil
		} else {
			cuentasSubgrupoTransaccion = append(cuentasSubgrupoTransaccion, map[string]interface{}{})
			return cuentasSubgrupoTransaccion, nil
		}
	} else {
		if err == nil {
			err = fmt.Errorf("Undesired Status Code: %d", response.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetCuentasContablesSubgrupo - request.GetJsonTest(urlcrud, &cuentasSubgrupo)",
			"err":     err,
			"status:": "502",
		}
		return nil, outputError
	}
}

func GetMovimientosKronos() (Movimientos_Arka []map[string]interface{}, outputError map[string]interface{}) {

	step := "0"

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetMovimientosKronos - Unhandled Error! - after step:" + step,
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var movimientos map[string]interface{}

	urlcrud := "http://" + beego.AppConfig.String("movimientosKronosService") + "tipo_movimiento?query=Activo:true&limit=-1"
	// logs.Debug("urlcrud:", urlcrud)
	if resp, err := request.GetJsonTest(urlcrud, &movimientos); err == nil && resp.StatusCode == 200 { // (2) error servicio caido
		step = "1"
		var data []map[string]interface{}
		if jsonString, err := json.Marshal(movimientos["Body"]); err == nil {
			step = "2"
			if err2 := json.Unmarshal(jsonString, &data); err2 == nil {
				step = "3"
				for _, movimiento := range data {
					if number := strings.Index(fmt.Sprintf("%v", movimiento["Acronimo"]), "arka"); number != -1 {
						Movimientos_Arka = append(Movimientos_Arka, map[string]interface{}{
							"Id":                movimiento["Id"],
							"Nombre":            movimiento["Nombre"],
							"Descripcion":       movimiento["Descripcion"],
							"Acronimo":          movimiento["Acronimo"],
							"Activo":            movimiento["Activo"],
							"FechaCreacion":     movimiento["FechaCreacion"],
							"FechaModificacion": movimiento["FechaModificacion"],
							"Parametros":        movimiento["Parametros"],
						})
					}
				}
				step = "4"
				return Movimientos_Arka, nil
			} else {
				logs.Error(err2)
				outputError = map[string]interface{}{
					"funcion": "GetMovimientosKronos - json.Unmarshal(jsonString, &data)",
					"err":     err,
					"status":  "500",
				}
				return nil, outputError
			}
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetMovimientosKronos - json.Marshal(movimientos[\"Body\"])",
				"err":     err,
				"status":  "500",
			}
			return nil, outputError
		}
	} else {
		if err == nil {
			err = fmt.Errorf("Undesired Status Code: %d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetMovimientosKronos - request.GetJsonTest(urlcrud, &movimientos)",
			"err":     err,
			"status:": "502",
		}
		return nil, outputError
	}
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
		urlSubgrupo += "&fields=SubgrupoId,TipoBienId&sortby=Id&order=desc"
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
