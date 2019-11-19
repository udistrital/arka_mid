package catalogoElementosHelper

import (
	"strconv"
	"fmt"
	"github.com/udistrital/arka_mid/helpers/cuentasContablesHelper"
	"encoding/json"
	"github.com/udistrital/arka_mid/models"
	"strings"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
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

	var	urlcrud         string
	var	cuentasSubgrupo []*models.CuentasGrupo
	var	cuentaCredito   *models.CuentaContable
	var	cuentaDebito    *models.CuentaContable
	

	urlcrud = "http://" + beego.AppConfig.String("catalogoElementosService") + "cuentas_subgrupo?query=SubgrupoId.Id:" + strconv.Itoa(int(subgrupoId)) + ",Activo:True&limit=-1"

	if response, err := request.GetJsonTest(urlcrud, &cuentasSubgrupo); err == nil { // (2) error servicio caido
		if response.StatusCode == 200  { // (3) error estado de la solicitud
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
	var	urlcrud         string
	urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "tipo_movimiento?query=Activo:true&limit=-1"

	if response, err := request.GetJsonTest(urlcrud, &movimientos); err == nil { // (2) error servicio caido
		if response.StatusCode == 200  { // (3) error estado de la solicitud
			var data []map[string]interface{}
			if jsonString, err := json.Marshal(movimientos["Body"]); err == nil {
				if err2 := json.Unmarshal(jsonString, &data); err2 == nil {
					for _, movimiento := range data {
						if number := strings.Index(fmt.Sprintf("%v",movimiento["Acronimo"]),"arka"); number != -1 {
							Movimientos_Arka = append(Movimientos_Arka, map[string]interface{}{
								"Id": 					movimiento["Id"],
      							"Nombre":				movimiento["Nombre"],
      							"Descripcion": 			movimiento["Descripcion"],
      							"Acronimo": 			movimiento["Acronimo"],
      							"Activo": 				movimiento["Activo"],
      							"FechaCreacion": 		movimiento["FechaCreacion"],
      							"FechaModificacion": 	movimiento["FechaModificacion"],
      							"Parametros": 			movimiento["Parametros"],
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
