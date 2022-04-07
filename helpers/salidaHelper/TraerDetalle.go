package salidaHelper

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/utils_oas/request"
)

func TraerDetalle(salida interface{}) (salida_ map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "TraerDetalle - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	sedeVacia := map[string]interface{}{
		"Id": 0,
	}
	ubicacionVacia := map[string]interface{}{
		"DependenciaId":   0,
		"EspacioFisicoId": 0,
	}

	if jsonString, err := json.Marshal(salida); err == nil {

		var data map[string]interface{}
		if err := json.Unmarshal(jsonString, &data); err == nil {
			str := fmt.Sprintf("%v", data["Detalle"])

			var data2 map[string]interface{}
			if err := json.Unmarshal([]byte(str), &data2); err == nil {

				urlcrud3 := "http://" + beego.AppConfig.String("oikos2Service") + "asignacion_espacio_fisico_dependencia"
				urlcrud3 += "?query=Id:" + fmt.Sprintf("%v", data2["ubicacion"])

				var tercero []map[string]interface{}
				var ubicacion []map[string]interface{}
				var sede []map[string]interface{}
				if data2["ubicacion"] != nil {
					if _, err := request.GetJsonTest(urlcrud3, &ubicacion); err == nil {

						var ubicacion2 map[string]interface{}
						if jsonString3, err := json.Marshal(ubicacion[0]["EspacioFisicoId"]); err == nil {
							if err2 := json.Unmarshal(jsonString3, &ubicacion2); err2 == nil {
								str2 := fmt.Sprintf("%v", ubicacion2["CodigoAbreviacion"])
								rgxp := regexp.MustCompile("[0-9]")
								str2 = rgxp.ReplaceAllString(str2, "")

								urlcrud4 := "http://" + beego.AppConfig.String("oikos2Service") + "espacio_fisico?query=CodigoAbreviacion:" + str2
								if _, err := request.GetJsonTest(urlcrud4, &sede); err != nil {
									logs.Error(err)
									outputError = map[string]interface{}{
										"funcion": "TraerDetalle - request.GetJsonTest(urlcrud4, &sede)",
										"err":     err,
										"status":  "502",
									}
									return nil, outputError
								}

							} else {
								logs.Error(err2)
								outputError = map[string]interface{}{
									"funcion": "TraerDetalle - json.Unmarshal(jsonString3, &ubicacion2)",
									"err":     err2,
									"status":  "500",
								}
								return nil, outputError
							}
						} else {
							logs.Error(err)
							outputError = map[string]interface{}{
								"funcion": "TraerDetalle - json.Marshal(ubicacion[0][\"EspacioFisicoId\"])",
								"err":     err,
								"status":  "500",
							}
							return nil, outputError
						}

					} else {
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "TraerDetalle - request.GetJsonTest(urlcrud3, &ubicacion)",
							"err":     err,
							"status":  "502",
						}
						return nil, outputError
					}
				} else {
					sede = append(sede, sedeVacia)
					ubicacion = append(ubicacion, ubicacionVacia)
				}

				Salida2 := map[string]interface{}{
					"Id":                      data["Id"],
					"Observacion":             data["Observacion"],
					"Sede":                    sede[0],
					"Dependencia":             ubicacion[0]["DependenciaId"],
					"Ubicacion":               ubicacion[0]["EspacioFisicoId"],
					"FechaCreacion":           data["FechaCreacion"],
					"FechaModificacion":       data["FechaModificacion"],
					"Activo":                  data["Activo"],
					"MovimientoPadreId":       data["MovimientoPadreId"],
					"FormatoTipoMovimientoId": data["FormatoTipoMovimientoId"],
					"EstadoMovimientoId":      data["EstadoMovimientoId"].(map[string]interface{})["Id"],
					"Consecutivo":             data2["consecutivo"],
					"ConsecutivoId":           data2["ConsecutivoContableId"],
				}

				if data2["funcionario"] != nil {

					urlcrud2 := "http://" + beego.AppConfig.String("tercerosService") + "tercero/?query=Id:" + fmt.Sprintf("%v", data2["funcionario"]) + "&fields=Id,NombreCompleto"
					if _, err := request.GetJsonTest(urlcrud2, &tercero); err != nil {
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "TraerDetalle - request.GetJsonTest(urlcrud3, &ubicacion)",
							"err":     err,
							"status":  "502",
						}
						return nil, outputError
					}

					Salida2["Funcionario"] = tercero[0]

				}

				return Salida2, nil

			} else {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "TraerDetalle - json.Unmarshal([]byte(str), &data2)",
					"err":     err,
					"status":  "500",
				}
				return nil, outputError
			}

		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "TraerDetalle - json.Unmarshal(jsonString, &data)",
				"err":     err,
				"status":  "500",
			}
			return nil, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "TraerDetalle - json.Marshal(salida)",
			"err":     err,
			"status":  "400",
		}
		return nil, outputError
	}
}
