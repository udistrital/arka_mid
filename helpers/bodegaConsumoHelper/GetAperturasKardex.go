package bodegaConsumoHelper

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	// "github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

func GetAperturasKardex() (Elementos []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetAperturasKardex - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var Elementos___ []map[string]interface{}
	url := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento?limit=-1"
	url += "&query=MovimientoId.FormatoTipoMovimientoId.CodigoAbreviacion:AP_KDX"
	// logs.Debug("url:", url)
	if res, err := request.GetJsonTest(url, &Elementos___); err == nil && res.StatusCode == 200 {
		// fmt.Println("Elementos___", Elementos___)

		if len(Elementos___) == 0 {
			return nil, nil
		}

		if keys := len(Elementos___[0]); keys != 0 {
			for _, elemento := range Elementos___ {

				// fmt.Println("Elemento", elemento)
				var data map[string]interface{}
				if jsonString, err := json.Marshal(elemento["MovimientoId"]); err == nil {
					// fmt.Println("Movimiento", jsonString)

					if err2 := json.Unmarshal(jsonString, &data); err2 == nil {
						// fmt.Println("DetalleMovimiento", data)

						str := fmt.Sprintf("%v", data["Detalle"])
						var data2 map[string]interface{}
						if err := json.Unmarshal([]byte(str), &data2); err == nil {
							// fmt.Println("Detalle", data2)

							var elemento_catalogo []map[string]interface{}
							url3 := "http://" + beego.AppConfig.String("catalogoElementosService") + "elemento?"
							url3 += "query=Id:" + fmt.Sprintf("%v", elemento["ElementoCatalogoId"])
							// logs.Debug("url3:", url3)
							if res, err := request.GetJsonTest(url3, &elemento_catalogo); err == nil && res.StatusCode == 200 {

								Elemento := map[string]interface{}{
									"MetodoValoracion":   data2["Metodo_Valoracion"],
									"CantidadMinima":     data2["Cantidad_Minima"],
									"CantidadMaxima":     data2["Cantidad_Maxima"],
									"FechaCreacion":      elemento["FechaCreacion"],
									"Observaciones":      data["Observacion"],
									"Id":                 data["Id"],
									"MovimientoPadreId":  data["MovimientoPadreId"],
									"ElementoCatalogoId": elemento_catalogo[0],
								}

								Elementos = append(Elementos, Elemento)

							} else {
								if err == nil {
									err = fmt.Errorf("Undesired Status Code: %d", res.StatusCode)
								}
								logs.Error(err)
								outputError = map[string]interface{}{
									"funcion": "GetAperturasKardex - request.GetJsonTest(url3, &elemento_catalogo)",
									"err":     err,
									"status":  "502",
								}
								return nil, outputError
							}

						} else {
							logs.Error(err)
							outputError = map[string]interface{}{
								"funcion": "GetAperturasKardex - json.Unmarshal([]byte(str), &data2)",
								"err":     err,
								"status":  "500",
							}
							return nil, outputError
						}

					} else {
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "GetAperturasKardex - json.Unmarshal(jsonString, &data)",
							"err":     err,
							"status":  "500",
						}
						return nil, outputError
					}

				} else {
					logs.Error(err)
					outputError = map[string]interface{}{
						"funcion": "GetAperturasKardex - json.Marshal(elemento[\"MovimientoId\"])",
						"err":     err,
						"status":  "500",
					}
					return nil, outputError
				}

			}

			return Elementos, nil
		} else {
			return Elementos___, nil
		}

	} else {
		if err == nil {
			err = fmt.Errorf("Undesired Status Code: %d", res.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAperturasKardex - request.GetJsonTest(url, &Elementos___)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}
