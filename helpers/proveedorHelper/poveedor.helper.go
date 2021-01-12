package proveedorHelper

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetProveedorId ...
func GetProveedorById(proveedorId int) (proveedor []*models.Proveedor, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetProveedorById", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	if proveedorId != 0 { // (1) error parametro

		if response, err := request.GetJsonTest("http://"+beego.AppConfig.String("administrativaService")+"informacion_proveedor?query=Id:"+strconv.Itoa(proveedorId), &proveedor); err == nil { // (2) error servicio caido
			if response.StatusCode == 200 { // (3) error estado de la solicitud
				return proveedor, nil
			} else {
				logs.Info("Error (3) estado de la solicitud")
				outputError = map[string]interface{}{"funcion": "GetProveedorById", "err": response.Status, "status": response.Status}
				return nil, outputError
			}
		} else {
			logs.Debug(err)
			logs.Info("Error (2) servicio caido")
			outputError = map[string]interface{}{"funcion": "GetProveedorById", "err": err, "status": "502"}
			return nil, outputError
		}
	} else {
		logs.Info("Error (1) Parametro")
		outputError = map[string]interface{}{"funcion": "GetProveedorById", "err": "null parameter", "status": "400"}
		return nil, outputError
	}
}
