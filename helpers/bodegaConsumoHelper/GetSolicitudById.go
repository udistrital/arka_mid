package bodegaConsumoHelper

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	// "github.com/udistrital/utils_oas/formatdata"
	crudTerceros "github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/utils_oas/request"
)

//GetTerceroById trae el nombre de un encargado por su id
func GetSolicitudById(id int) (Solicitud map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetSolicitudById - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var solicitud_ []map[string]interface{}
	var elementos___ []map[string]interface{}

	// url := "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento/" + fmt.Sprintf("%v", id) + ""
	url := "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento?query=Id:" + fmt.Sprintf("%v", id) + ""
	// logs.Debug(url)
	if res, err := request.GetJsonTest(url, &solicitud_); err == nil && res.StatusCode == 200 {

		// logs.Debug("solicitud_:")
		// formatdata.JsonPrint(solicitud_)
		// fmt.Println("")

		// TO-DO: Arreglar el CRUD! No debería retornar un arreglo con un elemento vacío ([{}])
		// Por máximo debería retornar el arreglo vacío! (sin el objeto vacío, [])
		// (Y uno de los siguientes estados: 204 o 404)
		if len(solicitud_) == 0 || len(solicitud_[0]) == 0 {
			err := fmt.Errorf("Movimiento %d no encontrado", id)
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "/GetSolicitudById",
				"err":     err,
				"status":  "404",
			}
			return nil, outputError
		}

		str := fmt.Sprintf("%v", solicitud_[0]["Detalle"])
		// logs.Debug(fmt.Sprintf("str: %s", str))
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(str), &data); err == nil {

			// logs.Debug("data:", data)
			if tercero, err := crudTerceros.GetNombreTerceroById(fmt.Sprintf("%v", data["Funcionario"])); err == nil {
				solicitud_[0]["Funcionario"] = tercero
			} else {
				return nil, err
			}
			var data_ []map[string]interface{}
			if jsonString, err := json.Marshal(data["Elementos"]); err == nil {
				if err2 := json.Unmarshal(jsonString, &data_); err2 == nil {

					for _, elementos := range data_ {
						// logs.Debug("k:", k, "- elementos:", elementos)

						if Elemento__, err := TraerElementoSolicitud(elementos); err == nil {
							Elemento__["Cantidad"] = elementos["Cantidad"]
							// fmt.Println(elementos["CantidadAprobada"])
							if elementos["CantidadAprobada"] != nil {
								Elemento__["CantidadAprobada"] = elementos["CantidadAprobada"]
							} else {
								Elemento__["CantidadAprobada"] = 0
							}

							elementos___ = append(elementos___, Elemento__)
						}
					}
					Solicitud = map[string]interface{}{
						"Solicitud": solicitud_,
						"Elementos": elementos___,
					}

					return Solicitud, nil

				} else {
					logs.Error(err2)
					outputError = map[string]interface{}{
						"funcion": "/GetSolicitudById - json.Marshal(data[\"Elementos\"])",
						"err":     err2,
						"status":  "500",
					}
					return nil, outputError
				}

			} else {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "/GetSolicitudById - json.Marshal(data[\"Elementos\"])",
					"err":     err,
					"status":  "500",
				}
				return nil, outputError
			}

		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "/GetSolicitudById - json.Unmarshal([]byte(str), &data)",
				"err":     err,
				"status":  "500",
			}
			return nil, outputError
		}

	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/GetSolicitudById - request.GetJsonTest(url, &solicitud_)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}
