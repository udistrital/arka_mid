package entradaHelper

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/tercerosHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"

	"github.com/udistrital/arka_mid/models"
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

// AddEntrada Transacción para registrar la información de una entrada
func AddEntrada(data models.Movimiento) map[string]interface{} {
	var (
		urlcrud      string
		res          map[string]interface{}
		resA         map[string]interface{}
		resM         map[string]interface{}
		resS         map[string]interface{}
		actaRecibido []models.TransaccionActaRecibido
		resultado    map[string]interface{}
	)

	detalleJSON := map[string]interface{}{}
	if err := json.Unmarshal([]byte(data.Detalle), &detalleJSON); err != nil {
		panic(err.Error())
	}

	year, _, _ := time.Now().Date()
	consec := Consecutivo{0, 1, year, 0, "Entradas", true}
	apiCons := "http://" + beego.AppConfig.String("consecutivosService") + "consecutivo"
	if err := request.SendJson(apiCons, "POST", &res, &consec); err == nil {
		resultado, _ := res["Data"].(map[string]interface{})
		numeroentrada := fmt.Sprintf("%05.0f", resultado["Consecutivo"]) + "-" + strconv.Itoa(year)
		vconsecutivo := detalleJSON["consecutivo"].(string) + "-" + numeroentrada
		detalleJSON["consecutivo"] = vconsecutivo
	} else {
		logs.Error(err)
		panic(err.Error())
	}
	var jsonData []byte
	jsonData, err1 := json.Marshal(detalleJSON)
	if err1 != nil {
		logs.Error(err1)
		panic(err1.Error())
	}
	data.Detalle = string(jsonData[:])

	// Solicita información acta

	actaRecibidoId := int(detalleJSON["acta_recibido_id"].(float64))

	if data.Id > 0 { // Si desde el cliente se envía el id del movimiento, se hace el put
		fmt.Println("Editar Entrada")
		urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento/" + strconv.Itoa(int(data.Id))

		if err := request.SendJson(urlcrud, "PUT", &res, &data); err == nil {

			urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo?query=ProcesoExterno:" + strconv.Itoa(int(data.Id))

			var data0 map[string]interface{}
			if _, err := request.GetJsonTest(urlcrud, &data0); err == nil {
				var data1 map[string]interface{}
				if jsonString, err := json.Marshal(data0); err == nil {
					if err := json.Unmarshal(jsonString, &data1); err == nil {
						var data2 = data1["Body"]
						var data3 []map[string]interface{}
						if jsonString1, err := json.Marshal(data2); err == nil {
							if err := json.Unmarshal(jsonString1, &data3); err == nil {

								urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo/" + strconv.Itoa(int(data.Id))

								procesoExterno := data.Id
								var formatoArka map[string]interface{}
								var idMovArka int

								if jsonString, err := json.Marshal(res["FormatoTipoMovimientoId"]); err == nil {
									if err := json.Unmarshal(jsonString, &formatoArka); err == nil {
										idMovArka = int(formatoArka["Id"].(float64))
									} else {
										panic(err.Error())
									}
								} else {
									panic(err.Error())
								}

								tipo := models.TipoMovimiento{Id: data.IdTipoMovimiento}
								movimientosKronos := models.MovimientoProcesoExterno{
									Id:                       int(data3[0]["Id"].(float64)),
									TipoMovimientoId:         &tipo,
									ProcesoExterno:           int64(procesoExterno),
									Activo:                   true,
									MovimientoProcesoExterno: idMovArka,
								}
								if err = request.SendJson(urlcrud, "PUT", &resM, &movimientosKronos); err == nil {
									resultado = resM
								} else {
									panic(err.Error())
								}
							} else {
								panic(err.Error())
							}
						} else {
							panic(err.Error())
						}
					} else {
						panic(err.Error())
					}
				} else {
					panic(err.Error())
				}
			} else {
				panic(err.Error())
			}
		} else {
			panic(err.Error())
		}
		return resultado

	} else { // Si desde el cliente NO se envía el id del movimiento, se hace el POST
		fmt.Println("Registrar entrada")
		urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/"

		// Solicita información acta

		if err := request.GetJson(urlcrud+strconv.Itoa(int(actaRecibidoId)), &actaRecibido); err == nil {
			// Envia información entrada
			urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento"

			if err = request.SendJson(urlcrud, "POST", &res, &data); err == nil {
				// Si la entrada tiene soportes
				if data.SoporteMovimientoId != 0 {
					// Envia información soporte (Si tiene)
					urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "soporte_movimiento"

					idEntrada := int(res["Id"].(float64))

					movimientoEntrada := models.Movimiento{Id: idEntrada}
					soporteMovimiento := models.SoporteMovimiento{
						DocumentoId:  data.SoporteMovimientoId,
						Activo:       true,
						MovimientoId: &movimientoEntrada,
					}

					if err = request.SendJson(urlcrud, "POST", &resS, &soporteMovimiento); err != nil {
						panic(err.Error())
					}
				}

				// Envia información movimientos Kronos
				urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo"

				procesoExterno := int64(res["Id"].(float64))

				var formatoArka map[string]interface{}
				var idMovArka int

				if jsonString, err := json.Marshal(res["FormatoTipoMovimientoId"]); err == nil {
					if err := json.Unmarshal(jsonString, &formatoArka); err != nil {
						panic(err.Error())
					}
					idMovArka = int(formatoArka["Id"].(float64))
				} else {
					panic(err.Error())
				}

				tipo := models.TipoMovimiento{Id: data.IdTipoMovimiento}
				movimientosKronos := models.MovimientoProcesoExterno{
					TipoMovimientoId:         &tipo,
					ProcesoExterno:           procesoExterno,
					Activo:                   true,
					MovimientoProcesoExterno: idMovArka,
				}

				if err = request.SendJson(urlcrud, "POST", &resM, &movimientosKronos); err == nil {
					// Cambia estado acta
					urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/" + strconv.Itoa(int(actaRecibidoId))
					actaRecibido[0].UltimoEstado.EstadoActaId.Id = 6
					actaRecibido[0].UltimoEstado.Id = 0

					if err = request.SendJson(urlcrud, "PUT", &resA, &actaRecibido[0]); err == nil {
						body := res
						body["Acta"] = resA
						resultado = body
					} else {
						panic(err.Error())
					}
				} else {
					panic(err.Error())
				}

			} else {
				panic(err.Error())
			}

		} else {
			panic(err.Error())
		}
	}

	return resultado
}

// GetEntrada ...
func GetEntrada(entradaId int) (consultaEntrada map[string]interface{}, outputError map[string]interface{}) {
	var (
		urlcrud        string
		tipoMovimiento map[string]interface{}
		movimientoArka map[string]interface{}
	)

	if entradaId != 0 { // (1) error parametro
		// Solicita información Movimientos KRONOS
		urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo?query=ProcesoExterno:" + strconv.Itoa(entradaId) + ",TipoMovimientoId.Acronimo:e_arka,Activo:true"

		if err := request.GetJson(urlcrud, &tipoMovimiento); err == nil {
			// Solicita información movimientos ARKA de acuedo a la información consultada por movimientos kronos
			var data []map[string]interface{}
			var movimientoId int

			if jsonString, err := json.Marshal(tipoMovimiento["Body"]); err == nil {
				if err := json.Unmarshal(jsonString, &data); err == nil {
					for _, movimiento := range data {
						movimientoId = int(movimiento["ProcesoExterno"].(float64))
					}
				}
			}

			urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento/" + strconv.Itoa(movimientoId)

			if response, err := request.GetJsonTest(urlcrud, &movimientoArka); err == nil { // (2) error servicio caido

				if response.StatusCode == 200 { // (3) error estado de la solicitud
					consultaEntrada = make(map[string]interface{})
					consultaEntrada = map[string]interface{}{"TipoMovimiento": data[0], "Movimiento": movimientoArka}
					return consultaEntrada, nil

				} else {
					logs.Info("Error (3) estado de la solicitud")
					outputError = map[string]interface{}{"Function": "GetEntrada:GetEntrada", "Error": response.Status}
					return nil, outputError
				}
			} else {
				logs.Info("Error (2) servicio caido")
				logs.Debug(err)
				outputError = map[string]interface{}{"Function": "GetEntrada", "Error": err}
				return nil, outputError
			}
		} else {
			return nil, outputError
		}
	} else {
		return nil, outputError
	}
}

// GetEntradas
func GetEntradas() (consultaEntradas []map[string]interface{}, outputError map[string]interface{}) {
	var (
		urlcrud                  string
		tipoMovimiento           map[string]interface{}
		movimientosId            []int
		tipoMovimientoEspecifico []interface{}
	)

	// Solicita información Movimientos KRONOS
	urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo?query=TipoMovimientoId.Acronimo:e_arka,Activo:true&limit=-1"
	if err := request.GetJson(urlcrud, &tipoMovimiento); err == nil {
		// Solicita información movimientos ARKA de acuedo a la información consultada por movimientos kronos
		var data []map[string]interface{}

		if jsonString, err := json.Marshal(tipoMovimiento["Body"]); err == nil {
			if err := json.Unmarshal(jsonString, &data); err == nil {
				for _, movimiento := range data {
					movimientosId = append(movimientosId, int(movimiento["ProcesoExterno"].(float64)))
					tipoMovimientoEspecifico = append(tipoMovimientoEspecifico, movimiento["TipoMovimientoId"])
				}
			}
		}

		// Solicita información a movimientos ARKA
		var contador = 0
		for _, i := range movimientosId {

			urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento/" + strconv.Itoa(i)

			var movimientoArka map[string]interface{}
			var aux map[string]interface{}

			if response, err := request.GetJsonTest(urlcrud, &movimientoArka); err == nil { // (2) error servicio caido

				if response.StatusCode == 200 { // (3) error estado de la solicitud
					aux = make(map[string]interface{})
					aux = map[string]interface{}{"TipoMovimiento": tipoMovimientoEspecifico[contador], "Movimiento": movimientoArka}
					consultaEntradas = append(consultaEntradas, aux)
					//logs.Info(movimientoArka)
				} else {
					logs.Info("Error (3) estado de la solicitud")
					outputError = map[string]interface{}{"Function": "GetEntrada:GetEntrada", "Error": response.Status}
				}
			} else {
				logs.Info("Error (2) servicio caido")
				logs.Debug(err)
				outputError = map[string]interface{}{"Function": "GetEntrada", "Error": err}
			}
			contador++
		}
	}
	if consultaEntradas != nil {
		return consultaEntradas, nil
	} else {
		return nil, outputError
	}
}

//GetEncargadoElemento busca al encargado de un elemento
func GetEncargadoElemento(placa string) (idElemento map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetEncargadoElemento - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var urlelemento string
	var detalle []map[string]interface{}

	if placa == "" {
		err := fmt.Errorf("La placa no puede estar en blanco")
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "GetEncargadoElemento - placa == ''", "status": "400", "err": err}
		return nil, outputError
	}

	if id, err := actaRecibido.GetIdElementoPlaca(placa); err == nil {
		if id == "" {
			err = fmt.Errorf("La placa '%s' no ha sido asignada a una salida", placa)
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "GetEncargadoElemento - id == ''", "status": "404", "err": err}
			return nil, outputError
		}
		urlelemento = "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento/?query=ElementoActaId:" + id
		fmt.Println(urlelemento)
		if resp, err := request.GetJsonTest(urlelemento, &detalle); err == nil && resp.StatusCode == 200 {
			logs.Info(detalle)
			cadena := detalle[0]["MovimientoId"].(map[string]interface{})["Detalle"]
			if resultado, err := utilsHelper.ConvertirStringJson(cadena); err == nil {
				idtercero := fmt.Sprintf("%v", resultado["funcionario"])
				if response, err := tercerosHelper.GetNombreTerceroById(idtercero); err == nil {
					if len(response) == 0 { // posible validación adicional:  || response["Id"].(string) == "0"
						err := fmt.Errorf("Respuesta inesperada en la respuesta de la función GetNombreTerceroById")
						logs.Error(err)
						outputError = map[string]interface{}{"funcion": "GetEncargadoElemento - tercerosHelper.GetNombreTerceroById(idtercero)", "status": "404", "err": err}
						return nil, outputError
					}
					var nombrecompleto = response["NombreCompleto"].(string)
					var aux = make(map[string]interface{})
					aux = map[string]interface{}{"Id": idtercero, "NombreCompleto": nombrecompleto}
					return aux, nil
				} else {
					return nil, err
				}
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "GetEncargadoElemento - actaRecibido.GetIdElementoPlaca(placa);", "status": "500", "err": err}
				return nil, outputError
			}
		} else {
			if err == nil {
				err = fmt.Errorf("Undesired Status Code: %d", resp.StatusCode)
			}
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "GetEncargadoElemento - request.GetJsonTest(urlelemento, &detalle) ", "status": "500", "err": err}
			return nil, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "GetEncargadoElemento", "status": "502", "err": err}
		return nil, outputError
	}
}

// AnularEntrada Anula una entrada y los movimientos posteriores a esta, el acta asociada queda en estado aceptada
func AnularEntrada(movimientoId int) (response map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetEncargadoElemento - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var (
		urlcrud        string
		res            map[string]interface{}
		resArka        map[string]interface{}
		resA           map[string]interface{}
		resM           map[string]interface{}
		movimientoArka models.Movimiento
		resP           []models.Movimiento
		actaRecibido   []models.TransaccionActaRecibido
	)

	res = make(map[string]interface{})
	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento/" + strconv.Itoa(int(movimientoId))
	if _, err := request.GetJsonTest(urlcrud, &movimientoArka); err == nil { // Get movimiento de api movimientos_arka_crud
		fmt.Println(movimientoArka.Id)
		if movimientoArka.Id > 0 {

			movimientoArka.EstadoMovimientoId.Id = 1
			if err := request.SendJson(urlcrud, "PUT", &resArka, &movimientoArka); err == nil { // Put en el api movimientos_arka_crud
				res["Movimiento"] = resArka

				var data0 map[string]interface{}
				urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo?query=ProcesoExterno:" + strconv.Itoa(int(movimientoId))
				if _, err := request.GetJsonTest(urlcrud, &data0); err == nil { // Get movimiento de api movimientos_crud
					var data1 []models.MovimientoProcesoExterno
					if jsonString1, err := json.Marshal(data0["Body"]); err == nil {
						if err := json.Unmarshal(jsonString1, &data1); err == nil {

							tipo := models.TipoMovimiento{Id: 34}
							movimientosKronos := models.MovimientoProcesoExterno{
								Id:                       data1[0].Id,
								TipoMovimientoId:         &tipo,
								ProcesoExterno:           data1[0].ProcesoExterno,
								Activo:                   true,
								MovimientoProcesoExterno: data1[0].MovimientoProcesoExterno,
							}
							urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo/" + strconv.Itoa(int(data1[0].Id))
							if err = request.SendJson(urlcrud, "PUT", &resM, &movimientosKronos); err == nil { // Put api movimientos_crud
								var detalleMovimiento map[string]interface{}
								if err := json.Unmarshal([]byte(movimientoArka.Detalle), &detalleMovimiento); err == nil {

									urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/" + fmt.Sprint(detalleMovimiento["acta_recibido_id"])
									if err := request.GetJson(urlcrud, &actaRecibido); err == nil { // Get informacion acta de api acta_recibido_crud
										actaRecibido[0].UltimoEstado.EstadoActaId.Id = 5
										actaRecibido[0].UltimoEstado.Id = 0
										if err = request.SendJson(urlcrud, "PUT", &resA, &actaRecibido[0]); err == nil { // Puesto que se anula la entrada, el acta debe quedar disponible para volver ser asociada a una entrada
											res["Acta Recibido"] = resA

											urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento?query=EstadoMovimientoId__Id:3,MovimientoPadreId__Id:" + strconv.Itoa(int(movimientoId))
											if _, err := request.GetJsonTest(urlcrud, &resP); err == nil { // Se determina si hay salidas que deban ser anuladas
												if resP[0].Id > 0 {
													// Cuando esté hecho el método para anular las salidas, agregarlo
													res["Salidas"] = resP
												}
											} else {
												logs.Error(err)
												outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
												return nil, outputError
											}
										} else {
											logs.Error(err)
											outputError = map[string]interface{}{"funcion": "AnularEntrada - actaRecibido.TransaccionActaRecibido(acta);", "status": "502", "err": err}
											return nil, outputError
										}
									} else {
										logs.Error(err)
										outputError = map[string]interface{}{"funcion": "AnularEntrada - actaRecibido.TransaccionActaRecibido(acta);", "status": "502", "err": err}
										return nil, outputError
									}
								} else {
									logs.Error(err)
									outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
									return nil, outputError
								}
							} else {
								logs.Error(err)
								outputError = map[string]interface{}{"funcion": "AnularEntrada - movimientos.MovimientoProcesoExterno(id);", "status": "502", "err": err}
								return nil, outputError
							}
						} else {
							logs.Error(err)
							outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
							return nil, outputError
						}
					} else {
						logs.Error(err)
						outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
						return nil, outputError
					}
				} else {
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "AnularEntrada - movimientos.MovimientoProcesoExterno(id);", "status": "502", "err": err}
					return nil, outputError
				}
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "AnularEntrada - movimientosArka.Movimiento(id);", "status": "502", "err": err}
				return nil, outputError
			}
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "AnularEntrada - movimientosArka.Movimiento(id); El movimiento no existe", "status": "502", "err": err}
			return nil, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "AnularEntrada - movimientosArka.Movimiento(id);", "status": "502", "err": err}
		return nil, outputError
	}
	return res, nil
}
