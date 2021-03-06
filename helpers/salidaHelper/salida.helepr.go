package salidaHelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	// "reflect"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/models"
	// "github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

type Consecutivo struct {
	Id          int
	ContextoId  int
	Year        int
	Consecutivo int
	Descripcion string
	Activo      bool
}

// AsignarPlaca Transacción para asignar las placas
func AsignarPlaca(m *models.Elemento) (resultado map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/ModificarPlaca", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	fmt.Printf("entro a asignar")
	fmt.Printf("%+v\n", m)
	year, month, day := time.Now().Date()
	//	fecstring := fmt.Sprintf("%4d", year) + fmt.Sprintf("%02d", int(month)) + fmt.Sprintf("%02d", day)

	consec := Consecutivo{0, 0, year, 0, "Placas", true}
	var (
		res       map[string]interface{} // models.SalidaGeneral
		fecstring string
	)

	apiCons := "http://" + beego.AppConfig.String("consecutivosService") + "consecutivo"
	putElemento := "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento/" + fmt.Sprintf("%d", m.Id)

	// Inserta salida en Movimientos ARKA
	// AsignarPlaca Transacción para asignar las placas
	if err := request.SendJson(apiCons, "POST", &res, &consec); err == nil {
		resultado, _ := res["Data"].(map[string]interface{})
		fmt.Printf("%+v\n", &resultado)
		fmt.Printf("%05.0f\n", resultado["Consecutivo"])
		fecstring := fmt.Sprintf("%4d", year) + fmt.Sprintf("%02d", int(month)) + fmt.Sprintf("%02d", day) + fmt.Sprintf("%05.0f", resultado["Consecutivo"])
		fmt.Printf(fecstring)
		m.Placa = fecstring
		if err := request.SendJson(putElemento, "PUT", &resultado, &m); err == nil {
			return resultado, nil
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "/AsignarPlaca",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/AsignarPlaca",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	fmt.Println("day: ", fecstring)
	return nil, outputError
}

// AddEntrada Transacción para registrar la información de una salida
func AddSalida(m *models.SalidaGeneral) (resultado []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/AddSalida", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var (
		res  map[string][](map[string]interface{}) // models.SalidaGeneral
		resM map[string]interface{}
	)

	movArka := "http://" + beego.AppConfig.String("movimientosArkaService") + "tr_salida"
	movKronos := "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo"

	// Inserta salida en Movimientos ARKA
	if err := request.SendJson(movArka, "POST", &res, &m); err == nil {

		// fmt.Printf("len(res): %v - len(res[\"Salidas\"]) %v\n", len(res), len(res["Salidas"]))
		// formatdata.JsonPrint(res["Salidas"])

		for _, salidaTr := range res["Salidas"] {

			// fmt.Printf("salidaTr[\"Elementos\"] T: %T -- salidaTr[\"Salida\"] T: %T\n", salidaTr["Elementos"], salidaTr["Salida"])
			// formatdata.JsonPrint(salidaTr["Elementos"])
			// formatdata.JsonPrint(salidaTr)

			if dataSalida, ok := salidaTr["Salida"].(map[string]interface{}); ok {
				if salidaID, ok := dataSalida["Id"].(float64); ok {
					procesoExterno := int64(salidaID)
					// logs.Debug(procesoExterno)

					var tipo models.TipoMovimiento

					if procesoExterno == 9 {
						tipo.Id = 16
					} else {
						tipo.Id = 22
					}
					movimientosKronos := models.MovimientoProcesoExterno{
						TipoMovimientoId: &tipo,
						ProcesoExterno:   procesoExterno,
						Activo:           true,
					}
					// fmt.Printf("movimientosKronos (%T): %v\n", movimientosKronos, movimientosKronos)

					// formatdata.JsonPrint(movimientosKronos)

					// Inserta salida en Movimientos KRONOS
					//*
					if err2 := request.SendJson(movKronos, "POST", &resM, &movimientosKronos); err2 == nil {
						salidaTr["MovimientosKronos"] = resM["Body"]
						resultado = append(resultado, salidaTr)
					} else {
						logs.Error(err2)
						outputError = map[string]interface{}{
							"funcion": "/AddSalida",
							"err":     err2,
							"status":  "502",
						}
						return nil, outputError
					}
					// */

				} else {
					logs.Error("carajo5")
				}

			} else {
				logs.Error("carajo4")
			}

		}

	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/AddSalida",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	return resultado, nil
}

func GetSalida(id int) (Salida map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetSalidas", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "tr_salida/" + strconv.Itoa(id)
	var salida_ map[string]interface{}
	if _, err := request.GetJsonTest(urlcrud, &salida_); err == nil {

		var data_ []map[string]interface{}
		if jsonString, err := json.Marshal(salida_["Elementos"]); err == nil {

			if err2 := json.Unmarshal(jsonString, &data_); err2 == nil {

				for i, elemento := range data_ {

					var elemento_ []map[string]interface{}
					urlcrud_ := "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento?query=Id:" + fmt.Sprintf("%v", elemento["ElementoActaId"]) + "&fields=Id,Nombre,TipoBienId,Marca,Serie,Placa,SubgrupoCatalogoId"
					if _, err := request.GetJsonTest(urlcrud_, &elemento_); err == nil {
						var subgrupo_ map[string]interface{}
						urlcrud_2 := "http://" + beego.AppConfig.String("catalogoElementosService") + "subgrupo/" + fmt.Sprintf("%v", elemento_[0]["SubgrupoCatalogoId"])
						if _, err := request.GetJsonTest(urlcrud_2, &subgrupo_); err == nil {
							data_[i]["Nombre"] = elemento_[0]["Nombre"]
							data_[i]["TipoBienId"] = elemento_[0]["TipoBienId"]
							data_[i]["SubgrupoCatalogoId"] = subgrupo_
							data_[i]["Marca"] = elemento_[0]["Marca"]
							data_[i]["Serie"] = elemento_[0]["Serie"]
							data_[i]["Placa"] = elemento_[0]["Placa"]

						} else {
							logs.Error(err)
							outputError = map[string]interface{}{
								"funcion": "/GetSalida - request.GetJsonTest(urlcrud_2, &subgrupo_)",
								"err":     err,
								"status":  "502",
							}
							return nil, outputError
						}
					} else {
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "/GetSalida - request.GetJsonTest(urlcrud_, &elemento_)",
							"err":     err,
							"status":  "502",
						}
						return nil, outputError
					}

					if _, err := request.GetJsonTest(urlcrud, &salida_); err != nil {
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "/GetSalida - request.GetJsonTest(urlcrud, &salida_)",
							"err":     err,
							"status":  "502",
						}
						return nil, outputError
					}

				}

			} else {
				logs.Error(err2)
				outputError = map[string]interface{}{
					"funcion": "/GetSalida - json.Unmarshal(jsonString, &data_)",
					"err":     err2,
					"status":  "500",
				}
				return nil, outputError
			}
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "/GetSalida - json.Marshal(salida_[\"Elementos\"])",
				"err":     err,
				"status":  "500",
			}
			return nil, outputError
		}

		if salida__, err := TraerDetalle(salida_["Salida"]); err == nil {

			Salida_final := map[string]interface{}{
				"Elementos": data_,
				"Salida":    salida__,
			}
			return Salida_final, nil

		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "/GetSalida - TraerDetalle(salida_[\"Salida\"])",
				"err":     err,
				"status":  "500",
			}
			return nil, outputError
		}

	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/GetSalida - request.GetJsonTest(urlcrud, &salida_)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}

func GetSalidas() (Salidas []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/GetSalidas - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento?query=FormatoTipoMovimientoId.CodigoAbreviacion__contains:SAL,FormatoTipoMovimientoId.Descripcion__contains:guardar,Activo:true&limit=-1"

	var salidas_ []map[string]interface{}
	if resp, err := request.GetJsonTest(urlcrud, &salidas_); err == nil && resp.StatusCode == 200 {
		logs.Info(fmt.Sprintf("#Salidas %d:  %v", len(salidas_), salidas_))

		if len(salidas_) == 0 || len(salidas_[0]) == 0 {
			err := errors.New("There's currently no outs records")
			logs.Warn(err)
			outputError = map[string]interface{}{
				"funcion": "/GetSalidas",
				"err":     err,
				"status":  "200", // TODO: Debería ser un 204 pero el cliente (Angular) se ofende... (hay que hacer varios ajustes)
			}
			return nil, outputError
		}

		for _, salida := range salidas_ {
			// fmt.Println("Salidas: ", salida)
			if salida__, err := TraerDetalle(salida); err == nil {

				Salidas = append(Salidas, salida__)

			} else {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "/GetSalidas",
					"err":     err,
					"status":  "502",
				}
				return nil, err
			}
		}
		return Salidas, nil

	} else if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/GetSalidas -request.GetJsonTest(urlcrud, &salidas_)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	} else {
		err := fmt.Errorf("Undesired Status Code: %d", resp.StatusCode)
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/GetSalidas - request.GetJsonTest(urlcrud, &salidas_)",
			"err":     err,
			"status":  "502", // (2) error servicio caido
		}
		return nil, outputError
	}
}

func TraerDetalle(salida interface{}) (salida_ map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/TraerDetalle", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	if jsonString, err := json.Marshal(salida); err == nil {

		var data map[string]interface{}
		if err2 := json.Unmarshal(jsonString, &data); err2 == nil {
			// fmt.Println("Salida: ", data)
			str := fmt.Sprintf("%v", data["Detalle"])

			var data2 map[string]interface{}
			if err := json.Unmarshal([]byte(str), &data2); err == nil {
				// fmt.Println("Detalle Salida: ", data2)
				urlcrud3 := "http://" + beego.AppConfig.String("oikos2Service") + "asignacion_espacio_fisico_dependencia?query=Id:" + fmt.Sprintf("%v", data2["ubicacion"])

				var tercero []map[string]interface{}
				var ubicacion []map[string]interface{}
				var sede []map[string]interface{}

				if _, err := request.GetJsonTest(urlcrud3, &ubicacion); err == nil {

					var ubicacion2 map[string]interface{}
					if jsonString3, err := json.Marshal(ubicacion[0]["EspacioFisicoId"]); err == nil {
						if err2 := json.Unmarshal(jsonString3, &ubicacion2); err2 == nil {
							str2 := fmt.Sprintf("%v", ubicacion2["CodigoAbreviacion"])

							z := strings.Split(str2, "")

							urlcrud4 := "http://" + beego.AppConfig.String("oikos2Service") + "espacio_fisico?query=CodigoAbreviacion:" + z[0] + z[1] + z[2] + z[3]

							if _, err := request.GetJsonTest(urlcrud4, &sede); err != nil {
								logs.Error(err)
								outputError = map[string]interface{}{
									"funcion": "/TraerDetalle - request.GetJsonTest(urlcrud4, &sede)",
									"err":     err,
									"status":  "502",
								}
								return nil, outputError
							}

						} else {
							logs.Error(err2)
							outputError = map[string]interface{}{
								"funcion": "/TraerDetalle - json.Unmarshal(jsonString3, &ubicacion2)",
								"err":     err2,
								"status":  "500",
							}
							return nil, outputError
						}
					} else {
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "/TraerDetalle - json.Marshal(ubicacion[0][\"EspacioFisicoId\"])",
							"err":     err,
							"status":  "500",
						}
						return nil, outputError
					}

				} else {
					logs.Error(err)
					outputError = map[string]interface{}{
						"funcion": "/TraerDetalle - request.GetJsonTest(urlcrud3, &ubicacion)",
						"err":     err,
						"status":  "502",
					}
					return nil, outputError
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
					"EstadoMovimientoId":      data["EstadoMovimientoId"],
				}

				if data2["funcionario"] != nil {

					urlcrud2 := "http://" + beego.AppConfig.String("tercerosService") + "tercero/?query=Id:" + fmt.Sprintf("%v", data2["funcionario"]) + "&fields=Id,NombreCompleto"
					if _, err := request.GetJsonTest(urlcrud2, &tercero); err != nil {
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "/TraerDetalle - request.GetJsonTest(urlcrud3, &ubicacion)",
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
					"funcion": "/TraerDetalle - json.Unmarshal([]byte(str), &data2)",
					"err":     err,
					"status":  "500",
				}
				return nil, outputError
			}

		} else {
			logs.Error(err2)
			outputError = map[string]interface{}{
				"funcion": "/TraerDetalle - json.Unmarshal(jsonString, &data)",
				"err":     err2,
				"status":  "500",
			}
			return nil, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/TraerDetalle - json.Marshal(salida)",
			"err":     err,
			"status":  "500",
		}
		return nil, outputError
	}
}
