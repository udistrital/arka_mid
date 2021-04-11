package entradaHelper

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/tercerosHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
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
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "GetEncargadoElemento - tercerosHelper.GetNombreTerceroById(idtercero)", "status": "502", "err": err}
					return nil, outputError
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
		urlcrud                 string
		res                     map[string]interface{}
		resMap                  map[string]interface{}
		movimientoArka          models.Movimiento
		transaccionActaRecibido models.TransaccionActaRecibido
		movimientosKronos       models.MovimientoProcesoExterno
		detalleMovimiento       map[string]interface{}
		tipoMovimiento          models.TipoMovimiento
		estadoActa              models.EstadoActa
		estadoMovimiento        models.EstadoMovimiento
		parametroTipoDebito     models.Parametro
		parametroTipoCredito    models.Parametro
		tipoComprobanteContable models.TipoComprobanteContable
		consecutivoId           int
		consecutivo             int
		transaccion             models.TransaccionMovimientos
		cuentasSubgrupo         []models.CuentaSubgrupo
		TipoEntradaKronos       models.TipoMovimiento
	)

	res = make(map[string]interface{})

	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento?query=Id:" + strconv.Itoa(int(movimientoId))
	var resMovArka []models.Movimiento
	if err := request.GetJson(urlcrud, &resMovArka); err == nil { // Get movimiento de api movimientos_arka_crud
		movimientoArka = resMovArka[0]
		if movimientoArka.Id > 0 {

			urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo?query=ProcesoExterno:" + strconv.Itoa(int(movimientoId))
			if err = request.GetJson(urlcrud, &resMap); err == nil { // Get movimiento de api movimientos_crud
				var resMovKronos []models.MovimientoProcesoExterno
				if jsonString, err := json.Marshal(resMap["Body"]); err == nil {
					if err = json.Unmarshal(jsonString, &resMovKronos); err == nil {
						resMap = make(map[string]interface{})
						movimientosKronos = resMovKronos[0]

						if err = json.Unmarshal([]byte(movimientoArka.Detalle), &detalleMovimiento); err == nil {
							var resTrActa []models.TransaccionActaRecibido

							urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/" + fmt.Sprint(detalleMovimiento["acta_recibido_id"])
							if err = request.GetJson(urlcrud, &resTrActa); err == nil { // Get informacion acta de api acta_recibido_crud
								transaccionActaRecibido = resTrActa[0]
								var resEstadoMovimiento []models.EstadoMovimiento

								urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "estado_movimiento?query=Nombre:Entrada%20Anulada"
								if err = request.GetJson(urlcrud, &resEstadoMovimiento); err == nil { // Get parametrización estado_movimiento de api movimientos_arka_crud
									estadoMovimiento = resEstadoMovimiento[0]

									urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "tipo_movimiento?query=Nombre:Entrada%20Anulada"
									if err = request.GetJson(urlcrud, &resMap); err == nil { // Get parametrización tipo_movimiento de api movimientos_crud
										var resTipoMovimiento []models.TipoMovimiento
										if jsonString, err = json.Marshal(resMap["Body"]); err == nil {
											if err = json.Unmarshal(jsonString, &resTipoMovimiento); err == nil {
												resMap = make(map[string]interface{})
												tipoMovimiento = resTipoMovimiento[0]

												urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "tipo_movimiento?query=Nombre__iexact:" + strings.ReplaceAll(movimientoArka.FormatoTipoMovimientoId.Nombre, " ", "%20")
												if err = request.GetJson(urlcrud, &resMap); err == nil { // Get parametrización tipo_movimiento de api movimientos_crud
													if jsonString, err = json.Marshal(resMap["Body"]); err == nil {
														if err = json.Unmarshal(jsonString, &resTipoMovimiento); err == nil {
															resMap = make(map[string]interface{})
															TipoEntradaKronos = resTipoMovimiento[0]
															var resEstadoActa []models.EstadoActa

															urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "estado_acta?query=Nombre:Aceptada"
															if err = request.GetJson(urlcrud, &resEstadoActa); err == nil { // Get parametrización acta de api acta_recibido_crud
																estadoActa = resEstadoActa[0]
																movimientoArka.EstadoMovimientoId.Id = estadoMovimiento.Id
																movimientosKronos.TipoMovimientoId.Id = tipoMovimiento.Id
																transaccionActaRecibido.UltimoEstado.EstadoActaId.Id = estadoActa.Id
																transaccionActaRecibido.UltimoEstado.Id = 0

																urlcrud = "http://" + beego.AppConfig.String("parametrosService") + "parametro?query=CodigoAbreviacion:MCC"
																if err = request.GetJson(urlcrud, &resMap); err == nil { // Get parámetro tipo movimiento contable crédito
																	if jsonString, err = json.Marshal(resMap["Data"]); err == nil {
																		var parametro []models.Parametro
																		if err = json.Unmarshal(jsonString, &parametro); err == nil {
																			resMap = make(map[string]interface{})
																			parametroTipoDebito = parametro[0]

																			urlcrud = "http://" + beego.AppConfig.String("parametrosService") + "parametro?query=CodigoAbreviacion:MCD"
																			if err = request.GetJson(urlcrud, &resMap); err == nil { // Get parámetro tipo movimiento contable débito
																				if jsonString, err = json.Marshal(resMap["Data"]); err == nil {
																					if err = json.Unmarshal(jsonString, &parametro); err == nil {
																						resMap = make(map[string]interface{})
																						parametroTipoCredito = parametro[0]

																						urlcrud = "http://" + beego.AppConfig.String("cuentasContablesService") + "tipo_comprobante"
																						if err = request.GetJson(urlcrud, &resMap); err == nil { // Para obtener código del tipo de comprobante
																							for _, sliceTipoComprobante := range resMap["Body"].([]interface{}) {
																								if sliceTipoComprobante.(map[string]interface{})["TipoDocumento"] == "E" {
																									if jsonString, err = json.Marshal(sliceTipoComprobante); err == nil {
																										if err = json.Unmarshal(jsonString, &tipoComprobanteContable); err == nil {
																											resMap = make(map[string]interface{})
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
																								}
																							}

																							year, _, _ := time.Now().Date()
																							postConsecutivo := Consecutivo{0, 2, year, 0, "Ajustes", true}
																							urlcrud = "http://" + beego.AppConfig.String("consecutivosService") + "consecutivo"
																							if err = request.SendJson(urlcrud, "POST", &resMap, &postConsecutivo); err == nil {
																								if consecutivoId, err = strconv.Atoi(fmt.Sprint(resMap["Data"].(map[string]interface{})["Id"])); err == nil {
																									if consecutivo, err = strconv.Atoi(fmt.Sprint(resMap["Data"].(map[string]interface{})["Consecutivo"])); err == nil {
																										resMap = make(map[string]interface{})
																										transaccion.ConsecutivoId = consecutivoId

																										// Se crea map para agrupar los valores totales según el código del subgrupo
																										mapSubgruposTotales := map[int]float64{}
																										for _, elemento := range transaccionActaRecibido.SoportesActa[0].Elementos { // Proceso para registrar el movimiento contable para cada elemento
																											if mapSubgruposTotales[elemento.SubgrupoCatalogoId] == 0 {
																												mapSubgruposTotales[elemento.SubgrupoCatalogoId] = elemento.ValorTotal
																											} else {
																												mapSubgruposTotales[elemento.SubgrupoCatalogoId] += elemento.ValorTotal
																											}
																										}

																										etiquetas := make(map[string]interface{})
																										etiquetas["TipoComprobanteId"] = tipoComprobanteContable.Codigo
																										if jsonString, err = json.Marshal(etiquetas); err == nil {
																											transaccion.Etiquetas = string(jsonString)
																											transaccion.Activo = true
																											transaccion.FechaTransaccion = time_bogota.Tiempo_bogota()
																											transaccion.Descripcion = "Transacción para registrar movimientos contables correspondientes a entrada del sistema arka"

																											for SubgrupoId, valor := range mapSubgruposTotales {
																												var cuentaDebito models.CuentaContable
																												var cuentaCredito models.CuentaContable
																												var movimientoDebito models.MovimientoTransaccion
																												var movimientoCredito models.MovimientoTransaccion

																												urlcrud = "http://" + beego.AppConfig.String("catalogoElementosService") + "cuentas_subgrupo?query=SubgrupoId__Id:" + strconv.Itoa(SubgrupoId) + ",SubtipoMovimientoId:" + strconv.Itoa(TipoEntradaKronos.Id) + ",Activo:true"
																												if err = request.GetJson(urlcrud, &cuentasSubgrupo); err == nil { // Obtiene cuentas que deben ser afectadas

																													urlcrud = "http://" + beego.AppConfig.String("cuentasContablesService") + "nodo_cuenta_contable/" + cuentasSubgrupo[0].CuentaDebitoId
																													if err = request.GetJson(urlcrud, &resMap); err == nil { // Se trae información de cuenta débito de api cuentas_contables_crud

																														if jsonString, err = json.Marshal(resMap["Body"]); err == nil {
																															if err := json.Unmarshal(jsonString, &cuentaDebito); err == nil {
																																resMap = make(map[string]interface{})

																																movimientoDebito.NombreCuenta = cuentaDebito.Nombre
																																movimientoDebito.CuentaId = cuentaDebito.Codigo
																																movimientoDebito.TipoMovimientoId = parametroTipoCredito.Id
																																movimientoDebito.Valor = valor
																																movimientoDebito.Descripcion = "Movimiento crédito registrado desde sistema arka"
																																movimientoDebito.Activo = true
																																movimientoDebito.TerceroId = 1 // Provisional
																																transaccion.Movimientos = append(transaccion.Movimientos, movimientoDebito)

																																urlcrud = "http://" + beego.AppConfig.String("cuentasContablesService") + "nodo_cuenta_contable/" + cuentasSubgrupo[0].CuentaCreditoId
																																if err = request.GetJson(urlcrud, &resMap); err == nil { // Se trae información de cuenta crédito de api cuentas_contables_crud

																																	if jsonString, err = json.Marshal(resMap["Body"]); err == nil {
																																		if err = json.Unmarshal(jsonString, &cuentaCredito); err == nil {
																																			movimientoCredito.NombreCuenta = cuentaCredito.Nombre
																																			movimientoCredito.CuentaId = cuentaCredito.Codigo
																																			movimientoCredito.TipoMovimientoId = parametroTipoDebito.Id
																																			movimientoCredito.Valor = valor
																																			movimientoCredito.Descripcion = "Movimiento débito registrado desde sistema arka"
																																			movimientoCredito.Activo = true
																																			movimientoCredito.TerceroId = 1 // Provisional
																																			transaccion.Movimientos = append(transaccion.Movimientos, movimientoCredito)
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
																																	outputError = map[string]interface{}{"funcion": "AnularEntrada - cuentasContables.nodo_cuenta_contable(cuenta);", "status": "502", "err": err}
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
																														outputError = map[string]interface{}{"funcion": "AnularEntrada - cuentasContables.nodo_cuenta_contable(cuenta);", "status": "502", "err": err}
																														return nil, outputError
																													}
																												} else {
																													logs.Error(err)
																													outputError = map[string]interface{}{"funcion": "AnularEntrada - catalogoElementos.cuentasSubgrupo(subgrupo);", "status": "502", "err": err}
																													return nil, outputError
																												}
																											}

																											res["transaccion"] = transaccion
																											var resMovmientoContable interface{}

																											urlcrud = "http://" + beego.AppConfig.String("movimientosContablesmidService") + "transaccion_movimientos/transaccion_movimientos"
																											if err = request.SendJson(urlcrud, "POST", &resMovmientoContable, &transaccion); err == nil {
																												if resMovmientoContable.(map[string]interface{})["Status"] == "201" {
																													res["Respuesta movimientos contables Entradas"] = resMovmientoContable

																													// Anulación de salidas asociadas
																													// Si el estado de movimientoArka es Entrada Asociada a una salida, continuar con la anulación de las salidas

																													consecutivoAjuste := "H20-" + fmt.Sprintf("%05d", consecutivo) + "-" + strconv.Itoa(year)
																													detalleMovimiento["consecutivo_ajuste"] = consecutivoAjuste
																													detalleMovimiento["mov_contable_ajuste_consecutivo_id"] = transaccion.ConsecutivoId

																													if jsonString, err = json.Marshal(detalleMovimiento); err == nil {
																														movimientoArka.Detalle = string(jsonString)
																														urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento/" + strconv.Itoa(int(movimientoId))
																														if err = request.SendJson(urlcrud, "PUT", &movimientoArka, &movimientoArka); err == nil { // Put en el api movimientos_arka_crud
																															res["arka"] = movimientoArka.Detalle
																															urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo/" + strconv.Itoa(movimientoArka.Id)
																															if err = request.SendJson(urlcrud, "PUT", &movimientosKronos, &movimientosKronos); err == nil { // Put api movimientos_crud

																																urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/" + fmt.Sprint(detalleMovimiento["acta_recibido_id"])
																																if err = request.SendJson(urlcrud, "PUT", &transaccionActaRecibido, &transaccionActaRecibido); err == nil { // Puesto que se anula la entrada, el acta debe quedar disponible para volver ser asociada a una entrada
																																	res["movArkaId"] = movimientoArka.EstadoMovimientoId.Id
																																	res["movKronosId"] = movimientosKronos.TipoMovimientoId.Id
																																	res["EstadoActaId"] = transaccionActaRecibido.UltimoEstado.EstadoActaId.Id
																																} else {
																																	logs.Error(err)
																																	outputError = map[string]interface{}{"funcion": "AnularEntrada - actaRecibido.TransaccionActaRecibido(acta);", "status": "502", "err": err}
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
																														outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
																														return nil, outputError
																													}
																												} else {
																													res["Respuesta movimientos contables Entradas"] = resMovmientoContable.(map[string]interface{})["Data"]
																													outputError = map[string]interface{}{"funcion": "AnularEntrada - movimientosContablesMid.postTransaccion;", "status": "502", "err": resMovmientoContable.(map[string]interface{})["Data"]}
																													return res, outputError
																												}
																											} else {
																												logs.Error(err)
																												outputError = map[string]interface{}{"funcion": "AnularEntrada - movimientosContablesMid.postTransaccion(movimiento);", "status": "502", "err": err}
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
																									outputError = map[string]interface{}{"funcion": "AnularEntrada - consecutivos.postConsecutivo; No se retornó un consecutivo válido", "status": "502", "err": err}
																									return nil, outputError
																								}
																							} else {
																								logs.Error(err)
																								outputError = map[string]interface{}{"funcion": "AnularEntrada - consecutivos.postConsecutivo;", "status": "502", "err": err}
																								return nil, outputError
																							}
																						} else {
																							logs.Error(err)
																							outputError = map[string]interface{}{"funcion": "AnularEntrada - cuentasContablesCrud.TipoComprobante(Codigo);", "status": "502", "err": err}
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
																				outputError = map[string]interface{}{"funcion": "AnularEntrada - parametros.Parametro(CodigoAbreviación);", "status": "502", "err": err}
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
																	outputError = map[string]interface{}{"funcion": "AnularEntrada - parametros.Parametro(CodigoAbreviación);", "status": "502", "err": err}
																	return nil, outputError
																}

															} else {
																logs.Error(err)
																outputError = map[string]interface{}{"funcion": "AnularEntrada - actaRecibido.EstadoActa", "status": "502", "err": err}
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
													outputError = map[string]interface{}{"funcion": "AnularEntrada - movimientos.TipoMovimiento", "status": "502", "err": err}
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
										outputError = map[string]interface{}{"funcion": "AnularEntrada - movimientos.TipoMovimiento", "status": "502", "err": err}
										return nil, outputError
									}
								} else {
									logs.Error(err)
									outputError = map[string]interface{}{"funcion": "AnularEntrada - movimientosArka.EstadoMovimiento", "status": "502", "err": err}
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
