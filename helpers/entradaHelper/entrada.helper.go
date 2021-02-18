package entradaHelper

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/tercerosHelper"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// AddEntrada Transacción para registrar la información de una entrada
func AddEntrada(data models.Movimiento) map[string]interface{} {
	var (
		urlcrud      string
		res          map[string]interface{}
		resA         map[string]interface{}
		resM         map[string]interface{}
		resS         map[string]interface{}
		resH         []map[string]interface{}
		actaRecibido []models.TransaccionActaRecibido
		resultado    map[string]interface{}
	)

	// fmt.Println("Activo", data)
	// fmt.Println("Detalle", data.Detalle)
	// fmt.Println("estado", data.EstadoMovimientoId)
	// fmt.Println("estado", data.EstadoMovimientoId)
	// fmt.Println("estado", data.EstadoMovimientoId)
	detalleJSON := map[string]interface{}{}
	if err := json.Unmarshal([]byte(data.Detalle), &detalleJSON); err != nil {
		panic(err.Error())
	}
	actaRecibidoId := int(detalleJSON["acta_recibido_id"].(float64))

	urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "historico_acta?query=Activo:true,ActaRecibidoId__Id:" + strconv.Itoa(int(actaRecibidoId))

	if _, err := request.GetJsonTest(urlcrud, &resH); err == nil { // (2) error servicio caido
		// Determina si el acta recibido está asociada a una entrada, de ser asi, hace el put : post
		var data1 map[string]interface{}
		if jsonString, err := json.Marshal(resH[0]); err == nil {
			if err := json.Unmarshal(jsonString, &data1); err == nil {
				var ss = data1["EstadoActaId"]
				var data2 map[string]interface{}
				if jsonString1, err := json.Marshal(ss); err == nil {
					if err := json.Unmarshal(jsonString1, &data2); err == nil {
						var ent = data2["CodigoAbreviacion"]
						if data.Id > 0 {
							fmt.Println("Editar Entrada")
							resultado = data2

							urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento/" + strconv.Itoa(int(data.Id))

							if err = request.SendJson(urlcrud, "PUT", &res, &data); err == nil {

								urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo?query=ProcesoExterno:" + strconv.Itoa(int(data.Id))

								var sss map[string]interface{}
								if _, err := request.GetJsonTest(urlcrud, &sss); err == nil {
									var data3 map[string]interface{}
									if jsonString2, err := json.Marshal(sss); err == nil {
										if err := json.Unmarshal(jsonString2, &data3); err == nil {
											var sssd = data3["Body"]
											fmt.Println("sssd:", sssd)
											var data4 []map[string]interface{}
											if jsonString3, err := json.Marshal(sssd); err == nil {
												if err := json.Unmarshal(jsonString3, &data4); err == nil {

													urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo/" + strconv.Itoa(int(data.Id))

													procesoExterno := data.Id
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
														Id:                       int(data4[0]["Id"].(float64)),
														TipoMovimientoId:         &tipo,
														ProcesoExterno:           int64(procesoExterno),
														Activo:                   true,
														MovimientoProcesoExterno: idMovArka,
													}
													if err = request.SendJson(urlcrud, "PUT", &resM, &movimientosKronos); err == nil {

													} else {
														fmt.Println("failed")
													}
												} else {
													fmt.Println(err)
												}
											}
										}
									}

									fmt.Println("hier")
								} else {
									fmt.Println(err)
								}

								// urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo"

								// 	procesoExterno := int64(res["Id"].(float64))

								// 	var formato_arka map[string]interface{}
								// 	var id_mov_arka int

								// 	if jsonString, err := json.Marshal(res["FormatoTipoMovimientoId"]); err == nil {
								// 		if err := json.Unmarshal(jsonString, &formato_arka); err != nil {
								// 			panic(err.Error())
								// 		}
								// 		id_mov_arka = int(formato_arka["Id"].(float64))
								// 	} else {
								// 		panic(err.Error())
								// 	}

								// 	tipo := models.TipoMovimiento{Id: data.IdTipoMovimiento}
								// 	movimientosKronos := models.MovimientoProcesoExterno{
								// 		TipoMovimientoId:         &tipo,
								// 		ProcesoExterno:           procesoExterno,
								// 		Activo:                   true,
								// 		MovimientoProcesoExterno: id_mov_arka,
								// 	}

								// 	if err = request.SendJson(urlcrud, "POST", &resM, &movimientosKronos); err == nil {
								// 		resultado = res
								// 	} else {
								// 		panic(err.Error())
								// 	}

							} else {
								fmt.Println("hiersss")
								panic(err.Error())
							}

						} else if ent == "AsociadoEntrad" {
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

									var formato_arka map[string]interface{}
									var id_mov_arka int

									if jsonString, err := json.Marshal(res["FormatoTipoMovimientoId"]); err == nil {
										if err := json.Unmarshal(jsonString, &formato_arka); err != nil {
											panic(err.Error())
										}
										id_mov_arka = int(formato_arka["Id"].(float64))
									} else {
										panic(err.Error())
									}

									tipo := models.TipoMovimiento{Id: data.IdTipoMovimiento}
									movimientosKronos := models.MovimientoProcesoExterno{
										TipoMovimientoId:         &tipo,
										ProcesoExterno:           procesoExterno,
										Activo:                   true,
										MovimientoProcesoExterno: id_mov_arka,
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
						} else {
							fmt.Println(res)
						}
					}
				}

			}
		}
	} else {
		panic(err.Error())
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

// GetEncargado busca al encargado de un elemento
func GetEncargadoElemento(placa string) (idElemento map[string]interface{}, outputError map[string]interface{}) {
	var urlelemento string
	var elemento map[string]interface{}
	if id, err := actaRecibido.GetIdElementoPlaca(placa); err == nil {
		urlelemento = "http://" + beego.AppConfig.String("movimientosArkaService") + "tr_encargado_elemento/" + id
		if response, err := request.GetJsonTest(urlelemento, &elemento); err == nil {
			fmt.Println("status: ", elemento)
			if response.StatusCode == 200 {
				if response, err := tercerosHelper.GetNombreTerceroById("elemento"); err == nil {
					elemento["funcionario"] = response["NombreCompleto"].(string)
					idElemento = elemento
					return
				} else {
					outputError = map[string]interface{}{"Function": "GetEncargadoElemento", "Error": err}
					return
				}

			} else {
				outputError = map[string]interface{}{"Function": "GetEncargadoElemento", "status": int(response.StatusCode), "Error": response.Status}
				return
			}
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{"Function": "GetEncargadoElemento", "status": int(response.StatusCode), "Error": err}
			return nil, outputError
		}
	} else {
		outputError = map[string]interface{}{"Function": "GetEncargadoElemento", "status": "404", "Error": err}
		return nil, outputError
	}
	return
}
