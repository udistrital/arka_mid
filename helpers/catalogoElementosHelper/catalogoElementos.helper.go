package catalogoElementosHelper

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/cuentasContablesHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetCatalogoById ...
func GetCatalogoById(catalogoId int) (catalogo *[]map[string]interface{}, outputError map[string]interface{}) {
	var (
		urlcrud string
	)

	urlcrud = "http://" + beego.AppConfig.String("catalogoElementosService") + "tr_catalogo/" + strconv.Itoa(int(catalogoId))

	if response, err := request.GetJsonTest(urlcrud, &catalogo); err == nil { // (2) error servicio caido
		if response.StatusCode == 200 { // (3) error estado de la solicitud
			logs.Debug(catalogo)
			return catalogo, nil
		} else {
			logs.Info("Error (3) estado de la solicitud")
			outputError = map[string]interface{}{"Function": "GetCatalogoById:GetCatalogoById", "Error": response.Status}
			return nil, outputError
		}
	} else {
		logs.Info("Error (2) servicio caido")
		outputError = map[string]interface{}{"Function": "GetCatalogoById", "Error": err}
		return nil, outputError
	}
}

// GetCuentasContablesGrupo ...
func GetCuentasContablesSubgrupo(subgrupoId int) (cuentasSubgrupoTransaccion []map[string]interface{}, outputError map[string]interface{}) {

	var urlcrud string
	var cuentasSubgrupo []*models.CuentasGrupo
	var cuentaCredito *models.CuentaContable
	var cuentaDebito *models.CuentaContable

	urlcrud = "http://" + beego.AppConfig.String("catalogoElementosService") + "cuentas_subgrupo?query=SubgrupoId.Id:" + strconv.Itoa(int(subgrupoId)) + ",Activo:True&limit=-1"

	if response, err := request.GetJsonTest(urlcrud, &cuentasSubgrupo); err == nil { // (2) error servicio caido
		if response.StatusCode == 200 { // (3) error estado de la solicitud
			fmt.Println(cuentasSubgrupo[0])
			if cuentasSubgrupo[0].Id != 0 {
				for _, cuenta := range cuentasSubgrupo {
					cuentaCredito, outputError = cuentasContablesHelper.GetCuentaContable(cuenta.CuentaCreditoId)
					cuentaDebito, outputError = cuentasContablesHelper.GetCuentaContable(cuenta.CuentaDebitoId)

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
			logs.Info("Error (3) estado de la solicitud")
			outputError = map[string]interface{}{"Function": "GetCuentasContablesGrupo:GetCuentasContablesGrupo", "Error": response.Status}
			return nil, outputError
		}
	} else {
		logs.Info("Error (2) servicio caido")
		outputError = map[string]interface{}{"Function": "GetCuentasContablesGrupo", "Error": err}
		return nil, outputError
	}
}

func GetMovimientosKronos() (Movimientos_Arka []map[string]interface{}, outputError map[string]interface{}) {

	var movimientos map[string]interface{}
	var urlcrud string
	urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "tipo_movimiento?query=Activo:true&limit=-1"

	if response, err := request.GetJsonTest(urlcrud, &movimientos); err == nil { // (2) error servicio caido
		if response.StatusCode == 200 { // (3) error estado de la solicitud
			var data []map[string]interface{}
			if jsonString, err := json.Marshal(movimientos["Body"]); err == nil {
				if err2 := json.Unmarshal(jsonString, &data); err2 == nil {
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
					return Movimientos_Arka, nil
				} else {
					logs.Info("Error (5) estado de la solicitud")
					outputError = map[string]interface{}{"Function": "GetCuentasContablesGrupo:GetCuentasContablesGrupo", "Error": response.Status}
					return nil, outputError
				}
			} else {
				logs.Info("Error (4) estado de la solicitud")
				outputError = map[string]interface{}{"Function": "GetCuentasContablesGrupo:GetCuentasContablesGrupo", "Error": response.Status}
				return nil, outputError
			}
		} else {
			logs.Info("Error (3) estado de la solicitud")
			outputError = map[string]interface{}{"Function": "GetCuentasContablesGrupo:GetCuentasContablesGrupo", "Error": response.Status}
			return nil, outputError
		}
	} else {
		logs.Info("Error (2) servicio caido")
		outputError = map[string]interface{}{"Function": "GetCuentasContablesGrupo", "Error": err}
		return nil, outputError
	}
}

func GetTipoMovimiento(arreglosubgrupos []interface{}) (subgrupos []models.SubgrupoCuentasMovimiento, outputError map[string]interface{}) {
	var urlcrud, urlcatalogo string
	urlcatalogo = "http://" + beego.AppConfig.String("catalogoElementosService") + "tr_cuentas_subgrupo"
	var cuentasasociadas []models.SubgrupoCuentasModelo
	var arreglocuentas []models.CuentasGrupoMovimiento
	var tipomovimiento map[string]interface{}

	if outputError := request.SendJson(urlcatalogo, "POST", &cuentasasociadas, arreglosubgrupos); outputError == nil {
		for _, subgrupocuentas := range cuentasasociadas {

			for _, cuenta := range subgrupocuentas.CuentasAsociadas {
				cuentaaso := models.CuentasGrupoMovimiento{
					Id:                cuenta.Id,
					CuentaCreditoId:   cuenta.CuentaCreditoId,
					CuentaDebitoId:    cuenta.CuentaDebitoId,
					FechaCreacion:     cuenta.FechaCreacion,
					FechaModificacion: cuenta.FechaModificacion,
					Activo:            cuenta.Activo,
					SubgrupoId:        cuenta.SubgrupoId,
				}
				urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "tipo_movimiento/" + strconv.Itoa(cuenta.SubtipoMovimientoId)
				if response, outputError := request.GetJsonTest(urlcrud, &tipomovimiento); outputError == nil {
					if response.StatusCode == 200 {
						//fmt.Println("tipo movimiento", tipomovimiento.Id, "  ", tipomovimiento.Nombre)
						tipomov := tipomovimiento["Body"].(map[string]interface{})
						layout := "2006-01-02T15:04:05.000Z"
						a, _ := time.Parse(layout, tipomov["FechaCreacion"].(string))
						b, _ := time.Parse(layout, tipomov["FechaModificacion"].(string))

						cuentaaso.SubtipoMovimientoId = &models.TipoMovimiento{
							Id:                int(tipomov["Id"].(float64)),
							Nombre:            tipomov["Nombre"].(string),
							Descripcion:       tipomov["Descripcion"].(string),
							Acronimo:          tipomov["Acronimo"].(string),
							Activo:            tipomov["Activo"].(bool),
							FechaCreacion:     a,
							FechaModificacion: b,
							Parametros:        tipomov["Parametros"].(string),
						}

					}
				} else {
					return nil, map[string]interface{}{"Function": "GetTipoMovimiento", "Error": outputError}
				}
				arreglocuentas = append(arreglocuentas, cuentaaso)
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

		}
	} else {
		return nil, map[string]interface{}{"Function": "GetTipoMovimiento", "Error": outputError}
	}
	return subgrupos, outputError
	//}

	//if response, err := request.GetJsonTest(urlcrud, &movimientos); err == nil { // (2) error servicio caido
	//	if response.StatusCode == 200 { // (3) error estado de la solicitud
	//		var data []map[string]interface{}
	//		if jsonString, err := json.Marshal(movimientos["Body"]); err == nil {
	//			if err2 := json.Unmarshal(jsonString, &data); err2 == nil {
	//				for _, movimiento := range data {
	//				}
	//			}
	//		}
	//	}
	//}
}
