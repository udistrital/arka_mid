package administrativa

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetProveedorById Retorna los datos de un proveedor a partir del Id como proveedor
//
// Deprecated: Traer de terceros_crud o terceros_mid
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
//
// Deprecated: Traer de terceros_crud o terceros_mid
func GetProveedorByDoc(docNum string) (proveedor *models.Proveedor, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetProveedorByDoc - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	if docNum != "" { // (1) error parametro
		var proveedores []*models.Proveedor
		urlProveedor := "http://" + beego.AppConfig.String("administrativaService") + "informacion_proveedor?query=NumDocumento:" + docNum
		if response, err := request.GetJsonTest(urlProveedor, &proveedores); err == nil && response.StatusCode == 200 { // (2) error servicio caido
			status := "500"
			if len(proveedores) == 1 && proveedores[0].Id > 0 {
				return proveedores[0], nil
			} else if len(proveedores) == 0 || proveedores[0].Id == 0 {
				err = fmt.Errorf("Proveedor con Doc.Num.: '%s' no registrado", docNum)
				status = "404"
			} else { // len(proveedores) > 1
				n := len(proveedores)
				s := ""
				if n >= 10 {
					s = " (o más)"
				}
				err = fmt.Errorf("Proveedor con Doc.Num.: '%s' registrado %d%s veces", docNum, n, s)
				status = "409"
			}
			logs.Warn(err)
			outputError = map[string]interface{}{
				"funcion": "GetProveedorByDoc - len(proveedores) == 0 || proveedores[0].Id == 0",
				"err":     err,
				"status":  status,
			}
			return nil, outputError
		} else {
			if err == nil {
				err = fmt.Errorf("Undesired Status Code: %d", response.StatusCode)
			}
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetProveedorByDoc - request.GetJsonTest(urlProveedor, &proveedores)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
	} else {
		err := fmt.Errorf("No se especificó un documento")
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetProveedorByDoc - docNum != \"\"",
			"err":     err,
			"status":  "400",
		}
		return nil, outputError
	}
}
