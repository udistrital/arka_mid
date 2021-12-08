package salidaHelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	// "reflect"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/cuentasContablesHelper"
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

// AsignarPlaca Transacción para asignar las placas
func AsignarPlaca(m *models.Elemento) (resultado map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "AsignarPlaca - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	// fmt.Printf("entro a asignar")
	// fmt.Printf("%+v\n", m)
	year, month, day := time.Now().Date()
	//	fecstring := fmt.Sprintf("%4d", year) + fmt.Sprintf("%02d", int(month)) + fmt.Sprintf("%02d", day)

	consec := Consecutivo{0, 0, year, 0, "Placas", true}
	var (
		res map[string]interface{} // models.SalidaGeneral
	)

	apiCons := "http://" + beego.AppConfig.String("consecutivosService") + "consecutivo"
	putElemento := "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento/" + fmt.Sprintf("%d", m.Id)

	// Inserta salida en Movimientos ARKA
	// AsignarPlaca Transacción para asignar las placas
	if err := request.SendJson(apiCons, "POST", &res, &consec); err == nil {
		resultado, _ := res["Data"].(map[string]interface{})
		// fmt.Printf("%+v\n", &resultado)
		// fmt.Printf("%05.0f\n", resultado["Consecutivo"])
		fecstring := fmt.Sprintf("%4d", year) + fmt.Sprintf("%02d", int(month)) + fmt.Sprintf("%02d", day) + fmt.Sprintf("%05.0f", resultado["Consecutivo"])
		// fmt.Println(fecstring)
		m.Placa = fecstring
		if err := request.SendJson(putElemento, "PUT", &resultado, &m); err == nil {
			return resultado, nil
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "AsignarPlaca - request.SendJson(putElemento, \"PUT\", &resultado, &m)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "AsignarPlaca - request.SendJson(apiCons, \"POST\", &res, &consec)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}

// PostTrSalidas Completa los detalles de las salidas y hace el respectivo registro en api movimientos_arka_crud
func PostTrSalidas(m *models.SalidaGeneral) (resultado map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "PostTrSalidas - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var (
		res                 map[string][](map[string]interface{})
		resEstadoMovimiento []models.EstadoMovimiento
	)

	resultado = make(map[string]interface{})

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "estado_movimiento?query=Nombre:Salida%20En%20Trámite"
	if err := request.GetJson(urlcrud, &resEstadoMovimiento); err != nil || len(resEstadoMovimiento) == 0 {
		status := "502"
		if err == nil {
			err = errors.New("len(resEstadoMovimiento) == 0")
			status = "404"
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "PostTrSalidas - request.GetJson(urlcrud, &resEstadoMovimiento)",
			"err":     err,
			"status":  status,
		}
		return nil, outputError
	}

	for _, salida := range m.Salidas {

		detalle := map[string]interface{}{}
		if err := json.Unmarshal([]byte(salida.Salida.Detalle), &detalle); err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "PostTrSalidas - json.Unmarshal([]byte(salida.Salida.Detalle), &detalle)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}

		if consecutivo, err := utilsHelper.GetConsecutivo("H21", 230, "Registro Salida Arka"); err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "PostTrSalidas - utilsHelper.GetConsecutivo(\"H21\", 230, \"Registro Salida Arka\")",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		} else {
			detalle["consecutivo"] = consecutivo
			if detalleJSON, err := json.Marshal(detalle); err != nil {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "PostTrSalidas - json.Marshal(detalle)",
					"err":     err,
					"status":  "500",
				}
				return nil, outputError
			} else {
				salida.Salida.Detalle = string(detalleJSON)
			}
		}

		salida.Salida.EstadoMovimientoId.Id = resEstadoMovimiento[0].Id
	}

	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "tr_salida"

	// Crea registros en api movimientos_arka_crud
	if err := request.SendJson(urlcrud, "POST", &res, &m); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "PostTrSalidas - request.SendJson(movArka, \"POST\", &res, &m)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	resultado["trSalida"] = res

	return resultado, nil
}

// AprobarSalida Aprobacion de una salida
func AprobarSalida(salidaId int) (result map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "AprobarSalida - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var (
		urlcrud                 string
		res                     map[string]interface{}
		resEstadoMovimiento     []models.EstadoMovimiento
		movArka                 []models.Movimiento
		detalleMovimiento       map[string]interface{}
		transaccionActaRecibido models.TransaccionActaRecibido
	)
	resultado := make(map[string]interface{})

	logs.Debug("********* PASO *********************")

	var detalle []map[string]interface{}

	urlelemento := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento/?query=MovimientoId:" + strconv.Itoa(salidaId)
	if resp, err := request.GetJsonTest(urlelemento, &detalle); err == nil && resp.StatusCode == 200 {
	} else {
		if err == nil {
			err = fmt.Errorf("Undesired Status Code: %d", resp.StatusCode)
		}
		outputError = map[string]interface{}{"funcion": "AprobarSalida - request.GetJsonTest(urlelemento, &detalle)", "status": "500", "err": err}
		return nil, outputError
	}

	/*
		var groupscatalogo = make(map[int]float64)
		i := 0
		for _, elemento := range transaccionActaRecibido.Elementos {
			fmt.Println("entra:")
			x := float64(0)
			if val, ok := groups[elemento.SubgrupoCatalogoId]; ok {
				x = val
			}
			groups[elemento.SubgrupoCatalogoId] = groups[elemento.SubgrupoCatalogoId] + x + elemento.ValorFinal
			i++
		}
	*/

	//	var groupscatalogo = make(map[int]float64)
	for i := 0; i < len(detalle); i++ {
		logs.Debug("****** el valor *******", detalle[i]["ElementoCatalogoId"].(float64))
	}

	// Se cambia el estado del movimiento en movimientos_arka_crud
	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento?query=Id:" + strconv.Itoa(int(salidaId))
	if err := request.GetJson(urlcrud, &movArka); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "AprobarSalida - request.GetJson(urlcrud, &movArka)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "estado_movimiento?query=Nombre:Salida%20Aprobada"
	if err := request.GetJson(urlcrud, &resEstadoMovimiento); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "AprobarSalida - request.GetJson(urlcrud, &resEstadoMovimiento)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	} else if len(resEstadoMovimiento) == 0 {
		err = errors.New("len(resEstadoMovimiento) == 0")
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "AprobarSalida - request.GetJson(urlcrud, &resEstadoMovimiento)",
			"err":     err,
			"status":  "404",
		}
		return nil, outputError
	}

	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento/" + strconv.Itoa(int(salidaId))
	movArka[0].EstadoMovimientoId.Id = resEstadoMovimiento[0].Id
	if err := request.SendJson(urlcrud, "PUT", &res, &movArka[0]); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "AprobarSalida - request.SendJson(urlcrud, \"PUT\", &res, &movArka[0])",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	resultado["movimientoArka"] = movArka[0]

	// Crea registro en movimientos_crud
	urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "tipo_movimiento?query=Nombre:" + movArka[0].FormatoTipoMovimientoId.Nombre
	urlcrud = strings.ReplaceAll(urlcrud, " ", "%20")
	if err := request.GetJson(urlcrud, &res); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "AprobarSalida - request.GetJson(urlcrud, &res)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	} else if reflect.TypeOf(res["Body"]).Kind() != reflect.Slice {
		err = errors.New("no se encuentra tipo_movimiento en api movimientos_crud")
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "AprobarSalida - reflect.TypeOf(res[\"Body\"]).Kind() != reflect.Slice",
			"err":     err,
			"status":  "404",
		}
		return nil, outputError
	}

	urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo"
	procesoExterno := int64(salidaId)
	idMovArka := int(movArka[0].FormatoTipoMovimientoId.Id)
	tipomvto := strconv.Itoa(int(res["Body"].([]interface{})[0].(map[string]interface{})["Id"].(float64)))
	tipoMovimientoId := models.TipoMovimiento{Id: int(res["Body"].([]interface{})[0].(map[string]interface{})["Id"].(float64))}
	movimientosKronos := models.MovimientoProcesoExterno{
		TipoMovimientoId:         &tipoMovimientoId,
		ProcesoExterno:           procesoExterno,
		Activo:                   true,
		MovimientoProcesoExterno: idMovArka,
	}

	if err := request.SendJson(urlcrud, "POST", &res, &movimientosKronos); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "AprobarSalida - request.SendJson(urlcrud, \"POST\", &resM, &movimientosKronos)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	if err := json.Unmarshal([]byte(movArka[0].MovimientoPadreId.Detalle), &detalleMovimiento); err == nil {
		var resTrActa models.TransaccionActaRecibido

		logs.Debug("********* 1 *********************", fmt.Sprint(detalleMovimiento["acta_recibido_id"]))
		urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/" + fmt.Sprint(detalleMovimiento["acta_recibido_id"])
		if err := request.GetJson(urlcrud, &resTrActa); err == nil { // Get informacion acta de api acta_recibido_crud
			transaccionActaRecibido = resTrActa
		} else {

			transaccionActaRecibido = resTrActa
		}
		logs.Debug("********* 2 *********************", resTrActa)
	}

	resultado["movimientoArka"] = movArka[0]

	// Transaccion contable

	logs.Debug("********* Llega *********************", tipomvto)
	var groups = make(map[int]float64)
	i := 0
	for _, elemento := range transaccionActaRecibido.Elementos {
		fmt.Println("entra:", elemento)
		x := float64(0)
		if val, ok := groups[elemento.SubgrupoCatalogoId]; ok {
			x = val
		}
		groups[elemento.SubgrupoCatalogoId] = groups[elemento.SubgrupoCatalogoId] + x + elemento.ValorFinal
		i++
	}
	logs.Debug("los grupos", groups)
	if i == 0 {
		return resultado, nil
	}

	if resA, outputError := cuentasContablesHelper.AsientoContable(groups, tipomvto, "Salida de almacen"); res == nil || outputError != nil {
		if outputError == nil {
			outputError = map[string]interface{}{
				"funcion": "AddEntrada -cuentasContablesHelper.AsientoContable(groups, tipomvto, \"asiento contable\");",
				"err":     res,
				"status":  "502",
			}
		}
		logs.Debug("Fue error")
		return nil, outputError
	} else {
		logs.Debug("Fue exito")
		resultado["transaccionContable"] = resA["resultadoTransaccion"]
	}
	return resultado, nil
}

func GetSalida(id int) (Salida map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetSalida - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
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
					urlcrud_ := "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento?query=Id:" + fmt.Sprintf("%v", elemento["ElementoActaId"]) + "&fields=Id,Nombre,Marca,Serie,Placa,SubgrupoCatalogoId"
					if _, err := request.GetJsonTest(urlcrud_, &elemento_); err == nil {
						var subgrupo_ []map[string]interface{}

						urlcrud_2 := "http://" + beego.AppConfig.String("catalogoElementosService") + "detalle_subgrupo?query=SubgrupoId__Id:" + fmt.Sprintf("%v", elemento_[0]["SubgrupoCatalogoId"])
						if _, err := request.GetJsonTest(urlcrud_2, &subgrupo_); err == nil {
							data_[i]["Nombre"] = elemento_[0]["Nombre"]
							data_[i]["TipoBienId"] = subgrupo_[0]["TipoBienId"]
							data_[i]["SubgrupoCatalogoId"] = subgrupo_[0]["SubgrupoId"]
							data_[i]["Marca"] = elemento_[0]["Marca"]
							data_[i]["Serie"] = elemento_[0]["Serie"]
							data_[i]["Placa"] = elemento_[0]["Placa"]

						} else {
							logs.Error(err)
							outputError = map[string]interface{}{
								"funcion": "GetSalida - request.GetJsonTest(urlcrud_2, &subgrupo_)",
								"err":     err,
								"status":  "502",
							}
							return nil, outputError
						}
					} else {
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "GetSalida - request.GetJsonTest(urlcrud_, &elemento_)",
							"err":     err,
							"status":  "502",
						}
						return nil, outputError
					}

					if _, err := request.GetJsonTest(urlcrud, &salida_); err != nil {
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "GetSalida - request.GetJsonTest(urlcrud, &salida_) (BIS)",
							"err":     err,
							"status":  "502",
						}
						return nil, outputError
					}

				}

			} else {
				logs.Error(err2)
				outputError = map[string]interface{}{
					"funcion": "GetSalida - json.Unmarshal(jsonString, &data_)",
					"err":     err2,
					"status":  "500",
				}
				return nil, outputError
			}
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetSalida - json.Marshal(salida_[\"Elementos\"])",
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
			return nil, err
		}

	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetSalida - request.GetJsonTest(urlcrud, &salida_)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}

func GetSalidas(tramiteOnly bool) (Salidas []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetSalidas - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()
	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento?limit=-1"
	urlcrud += "&query=Activo:true,EstadoMovimientoId__Nombre"

	if tramiteOnly {
		urlcrud += ":Salida%20En%20Trámite"
	} else {
		urlcrud += "__startswith:Salida"
	}

	var salidas_ []map[string]interface{}
	if resp, err := request.GetJsonTest(urlcrud, &salidas_); err == nil && resp.StatusCode == 200 {
		logs.Info(fmt.Sprintf("#Salidas %d:  %v", len(salidas_), salidas_))

		if len(salidas_) == 0 || len(salidas_[0]) == 0 {
			return nil, nil
		}

		for _, salida := range salidas_ {
			// fmt.Println("Salidas: ", salida)
			if salida__, err := TraerDetalle(salida); err == nil {
				Salidas = append(Salidas, salida__)
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "GetSalidas - TraerDetalle(salida)",
					"err":     err,
					"status":  "502",
				}
				return nil, err
			}
		}
		return Salidas, nil

	} else {
		if err == nil {
			err = fmt.Errorf("Undesired Status Code: %d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetSalidas - request.GetJsonTest(urlcrud, &salidas_)",
			"err":     err,
			"status":  "502", // (2) error servicio caido
		}
		return nil, outputError
	}
}

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
			// fmt.Println("Salida: ", data)
			str := fmt.Sprintf("%v", data["Detalle"])

			var data2 map[string]interface{}
			if err := json.Unmarshal([]byte(str), &data2); err == nil {
				// fmt.Println("Detalle Salida: ", data2)

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

								z := strings.Split(str2, "")

								urlcrud4 := "http://" + beego.AppConfig.String("oikos2Service") + "espacio_fisico?query=CodigoAbreviacion:" + z[0] + z[1] + z[2] + z[3]

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
