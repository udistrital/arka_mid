package bajasHelper

import (
	"encoding/json"
	"fmt"
	"strconv"

	// "strings"
	// "reflect"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/salidaHelper"
	"github.com/udistrital/arka_mid/helpers/tercerosHelper"
	"github.com/udistrital/arka_mid/helpers/ubicacionHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/utils_oas/request"
)

func TraerDatosElemento(id int) (Elemento map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/TraerDatosElemento", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var elemento_movimiento_ []map[string]interface{}
	// var movimiento_ map[string]interface{}

	url := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento/?query=ElementoActaId:" + fmt.Sprintf("%v", id) + ",Activo:true"
	if _, err := request.GetJsonTest(url, &elemento_movimiento_); err == nil {

		if v, err := salidaHelper.TraerDetalle(elemento_movimiento_[0]["MovimientoId"]); err == nil {

			fmt.Println("Elemento Movimiento: ", elemento_movimiento_)

			var movimiento_ map[string]interface{}
			if jsonString3, err := json.Marshal(elemento_movimiento_[0]["MovimientoId"]); err == nil {
				if err2 := json.Unmarshal(jsonString3, &movimiento_); err2 == nil {
					movimiento_["MovimientoPadreId"] = nil
				}
			}

			var elemento_ []map[string]interface{}

			urlcrud_ := "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento?query=Id:" + fmt.Sprintf("%v", elemento_movimiento_[0]["ElementoActaId"]) + "&fields=Id,Nombre,TipoBienId,Marca,Serie,Placa,SubgrupoCatalogoId"
			if _, err := request.GetJsonTest(urlcrud_, &elemento_); err == nil {

				fmt.Println("Elemento: ", elemento_)

				var subgrupo_ map[string]interface{}
				urlcrud_2 := "http://" + beego.AppConfig.String("catalogoElementosService") + "subgrupo/" + fmt.Sprintf("%v", elemento_[0]["SubgrupoCatalogoId"])
				if _, err := request.GetJsonTest(urlcrud_2, &subgrupo_); err == nil {
					Elemento := map[string]interface{}{
						"Id":                 elemento_[0]["Id"],
						"Placa":              elemento_[0]["Placa"],
						"Nombre":             elemento_[0]["Nombre"],
						"TipoBienId":         elemento_[0]["TipoBienId"],
						"Entrada":            v["MovimientoPadreId"],
						"Salida":             movimiento_,
						"SubgrupoCatalogoId": subgrupo_,
						"Marca":              elemento_[0]["Marca"],
						"Serie":              elemento_[0]["Serie"],
						"Funcionario":        v["Funcionario"],
						"Sede":               v["Sede"],
						"Dependencia":        v["Dependencia"],
						"Ubicacion":          v["Ubicacion"],
					}

					// elemento_[0]["SubgrupoCatalogoId"] = subgrupo_
					// elemento_movimiento_[0]["ElementoActaId"] = elemento_[0]
					// Elemento = elemento_movimiento_[0]
					return Elemento, nil

				} else {
					logs.Error(err)
					outputError = map[string]interface{}{
						"funcion": "/TraerDatosElemento - request.GetJsonTest(urlcrud_2, &subgrupo_)",
						"err":     err,
						"status":  "502",
					}
					return nil, outputError
				}
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "/TraerDatosElemento - request.GetJsonTest(urlcrud_, &elemento_)",
					"err":     err,
					"status":  "502",
				}
				return nil, outputError
			}
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "/TraerDatosElemento - salidaHelper.TraerDetalle(elemento_movimiento_[0][\"MovimientoId\"])",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/TraerDatosElemento - request.GetJsonTest(url, &elemento_movimiento_)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}

func GetAllSolicitudes() (historicoActa []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetAllSolicitudes - Uncaught Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var Solicitudes []map[string]interface{}

	url := "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento?query=FormatotipoMovimientoId.CodigoAbreviacion:SOL_BAJA,Activo:true&limit=-1"

	if _, err := request.GetJsonTest(url, &Solicitudes); err == nil { // (2) error servicio caido
		fmt.Println("solicitudes: ", Solicitudes)

		if len(Solicitudes) == 0 || len(Solicitudes[0]) == 0 {
			logs.Warn(err)
			outputError = map[string]interface{}{
				"funcion": "GetAllSolicitudes - len(Solicitudes) == 0 || len(Solicitudes[0]) == 0",
				"err":     "sin",
				"status":  "200", // TODO: Debería ser un 204 pero el cliente (Angular) se ofende... (hay que hacer varios ajustes)
			}
			return nil, outputError
		}

		tercerosBuffer := make(map[int]interface{})
		ubicacionesBuffer := make(map[int]interface{})

		for _, solicitud := range Solicitudes {

			var data_ map[string]interface{}
			var data2_ map[string]interface{}
			var data3_ map[string]interface{}
			var Tercero_ map[string]interface{}
			var Revisor_ map[string]interface{}
			var Ubicacion_ map[string]interface{}

			Ubicacion_ = nil

			if data, err := utilsHelper.ConvertirStringJson(solicitud["Detalle"]); err == nil {
				data_ = data
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "GetAllSolicitudes - utilsHelper.ConvertirStringJson(solicitud[\"Detalle\"])",
					"err":     err,
					"status":  "500",
				}
				return nil, outputError
			}
			if data, err := utilsHelper.ConvertirInterfaceMap(solicitud["EstadoMovimientoId"]); err == nil {
				data2_ = data
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "GetAllSolicitudes - utilsHelper.ConvertirInterfaceMap(solicitud[\"EstadoMovimientoId\"])",
					"err":     err,
					"status":  "500",
				}
				return nil, outputError
			}

			requestTercero := func(id string) func() (interface{}, map[string]interface{}) {
				return func() (interface{}, map[string]interface{}) {
					if Tercero, err := tercerosHelper.GetNombreTerceroById(id); err == nil {
						return Tercero, nil
					}
					return nil, nil
				}
			}

			funcionarioIDstr := fmt.Sprintf("%v", data_["Funcionario"])
			if funcionarioID, err := strconv.Atoi(funcionarioIDstr); err == nil {
				if v, err := utilsHelper.BufferGeneric(funcionarioID, tercerosBuffer, requestTercero(funcionarioIDstr), nil, nil); err == nil {
					if v2, ok := v.(map[string]interface{}); ok {
						Tercero_ = v2
					}
				}
			}

			revisorIDstr := fmt.Sprintf("%v", data_["Revisor"])
			if revisorID, err := strconv.Atoi(revisorIDstr); err == nil {
				if v, err := utilsHelper.BufferGeneric(revisorID, tercerosBuffer, requestTercero(revisorIDstr), nil, nil); err == nil {
					if v2, ok := v.(map[string]interface{}); ok {
						Revisor_ = v2
					}
				}
			}

			ubicacionIdStr := fmt.Sprintf("%v", data_["Ubicacion"])
			requestUbicacion := func() (interface{}, map[string]interface{}) {
				if ubicacion, err := ubicacionHelper.GetAsignacionSedeDependencia(ubicacionIdStr); err == nil {
					return ubicacion, nil
				}
				return nil, nil
			}
			if ubicacionId, err := strconv.Atoi(ubicacionIdStr); err == nil {
				if v, err := utilsHelper.BufferGeneric(ubicacionId, ubicacionesBuffer, requestUbicacion, nil, nil); err == nil {
					if v2, ok := v.(map[string]interface{}); ok {
						Ubicacion_ = v2
					}
				}
			}

			if Ubicacion_ != nil {
				if jsonString2, err := json.Marshal(Ubicacion_["EspacioFisicoId"]); err == nil {
					if err2 := json.Unmarshal(jsonString2, &data3_); err2 != nil {
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "/GetAllSolicitudes",
							"err":     err,
							"status":  "500",
						}
						return nil, outputError
					}
				}
			} else {
				data3_ = map[string]interface{}{
					"Nombre": "Ubicacion No Especificada",
				}
			}

			fmt.Println(data3_)
			Acta := map[string]interface{}{
				"Ubicacion":         data3_["Nombre"],
				"Activo":            solicitud["Activo"],
				"FechaCreacion":     solicitud["FechaCreacion"],
				"FechaVistoBueno":   data_["FechaVistoBueno"],
				"FechaModificacion": solicitud["FechaModificacion"],
				"Id":                solicitud["Id"],
				"Observaciones":     solicitud["Observacion"],
				"Funcionario":       Tercero_["NombreCompleto"],
				"Revisor":           Revisor_["NombreCompleto"],
				"Estado":            data2_["Nombre"],
			}
			historicoActa = append(historicoActa, Acta)
		}
		return historicoActa, nil

	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/GetAllSolicitudes",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}

func TraerDetalle(id int) (Solicitud map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/TraerDetalle",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var Elementos__ []map[string]interface{}

	var data map[string]interface{}

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento/" + fmt.Sprintf("%v", id)

	if _, err := request.GetJsonTest(urlcrud, &data); err == nil {

		if data_, err := utilsHelper.ConvertirStringJson(data["Detalle"]); err == nil {

			if Sede, Dependencia, Ubicacion, err := ubicacionHelper.GetSedeDependenciaUbicacion(fmt.Sprintf("%v", data_["Ubicacion"])); err == nil {
				if Funcionario, err := tercerosHelper.GetNombreTerceroById(fmt.Sprintf("%v", data_["Funcionario"])); err == nil {
					if Revisor, err := tercerosHelper.GetNombreTerceroById(fmt.Sprintf("%v", data_["Revisor"])); err == nil {
						if Elementos, err := utilsHelper.ConvertirInterfaceArrayMap(data_["Elementos"]); err == nil {
							for _, elemento := range Elementos {
								id_, _ := strconv.Atoi(fmt.Sprintf("%v", elemento["Id"]))

								if Elemento_, err := TraerDatosElemento(id_); err == nil {

									Elemento_["Observaciones"] = elemento["Observaciones"]
									Elemento_["Soporte"] = elemento["Soporte"]
									Elemento_["TipoBaja"] = elemento["TipoBaja"]
									Elementos__ = append(Elementos__, Elemento_)
								}

							}
							Solicitud = map[string]interface{}{
								"Id":                data["Id"],
								"Sede":              Sede,
								"Dependencia":       Dependencia,
								"Ubicacion":         Ubicacion,
								"Funcionario":       Funcionario,
								"Revisor":           Revisor,
								"FechaCreacion":     data["FechaCreacion"],
								"FechaModificacion": data["FechaModificacion"],
								"FechaVistoBueno":   data_["FechaVistoBueno"],
								"Estado":            data["FechaCreacion"],
								"Activo":            data["FechaCreacion"],
								"Elementos":         Elementos__,
							}
							return Solicitud, nil

						} else {
							logs.Error(err)
							outputError = map[string]interface{}{
								"funcion": "/TraerDetalle",
								"err":     err,
								"status":  "500",
							}
							return nil, outputError
						}
					} else {
						return nil, err
					}
				} else {
					return nil, err
				}
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "/TraerDetalle",
					"err":     err,
					"status":  "502",
				}
				return nil, outputError
			}
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "/TraerDetalle",
				"err":     err,
				"status":  "500",
			}
			return nil, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/TraerDetalle",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}
