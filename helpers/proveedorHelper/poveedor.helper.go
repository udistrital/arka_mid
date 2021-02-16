package proveedorHelper

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetProveedorById Retorna los datos de un proveedor a partir del Id como proveedor
func GetProveedorById(proveedorId int) (proveedor []*models.Proveedor, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/GetProveedorById - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	if proveedorId > 0 { // (1) error parametro

		urlProveedor := "http://" + beego.AppConfig.String("administrativaService") + "informacion_proveedor?query=Id:" + strconv.Itoa(proveedorId)
		if response, err := request.GetJsonTest(urlProveedor, &proveedor); err == nil { // (2) error servicio caido
			if response.StatusCode == 200 { // (3) error estado de la solicitud
				return proveedor, nil
			} else {
				err := fmt.Errorf("Undesired Status: %s", response.Status)
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "GetProveedorById - request.GetJsonTest(urlProveedor, &proveedor)",
					"err":     err,
					"status":  "500", // Error (3) estado de la solicitud
				}
				return nil, outputError
			}
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetProveedorById - request.GetJsonTest(urlProveedor, &proveedor)",
				"err":     err,
				"status":  "502", // Error (2) servicio caido
			}
			return nil, outputError
		}
	} else {
		err := fmt.Errorf("proveedorId MUST be greater than 0")
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetProveedorById - proveedorId > 0",
			"err":     err,
			"status":  "400", // (1) error parametro
		}
		return nil, outputError
	}
}

// GetProveedorByDoc Retorna los datos de un proveedor a partir del # de documento
func GetProveedorByDoc(docNum string) (proveedor []*models.Proveedor, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetProveedorByDoc", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	if docNum != "" { // (1) error parametro
		if response, err := request.GetJsonTest("http://"+beego.AppConfig.String("administrativaService")+"informacion_proveedor?query=NumDocumento:"+docNum, &proveedor); err == nil { // (2) error servicio caido
			if response.StatusCode == 200 { // (3) error estado de la solicitud
				return proveedor, nil
			} else {
				outputError = map[string]interface{}{"funcion": "GetProveedorByDoc", "err": "Error (3) estado de la solicitud", "status": response.Status}
				logs.Error(outputError)
				return nil, outputError
			}
		} else {
			outputError = map[string]interface{}{"funcion": "GetProveedorByDoc", "err": err, "status": "502"}
			logs.Error(outputError)
			return nil, outputError
		}
	} else {
		outputError = map[string]interface{}{"funcion": "GetProveedorByDoc", "err": "Error (1) Parametro", "status": "400"}
		logs.Error(outputError)
		return nil, outputError
	}
}
