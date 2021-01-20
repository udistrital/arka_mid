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
				"funcion": "/GetAllSolicitudes",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var Solicitudes []map[string]interface{}
	var Terceros []map[string]interface{}
	var Ubicaciones []map[string]interface{}

	url := "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento?query=FormatotipoMovimientoId.CodigoAbreviacion:SOL_BAJA,Activo:true&limit=-1"

	if _, err := request.GetJsonTest(url, &Solicitudes); err == nil { // (2) error servicio caido
		fmt.Println("solicitudes: ", Solicitudes)

		if len(Solicitudes) == 0 || len(Solicitudes[0]) == 0 {
			logs.Warn(err)
			outputError = map[string]interface{}{
				"funcion": "/GetAllSolicitudes",
				"err":     "sin",
				"status":  "200", // TODO: Debería ser un 204 pero el cliente (Angular) se ofende... (hay que hacer varios ajustes)
			}
			return nil, outputError
		}
	




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
					"funcion": "/GetAllSolicitudes",
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
					"funcion": "/GetAllSolicitudes",
					"err":     err,
					"status":  "500",
				}
				return nil, outputError
			}

			if Terceros == nil {
				if Tercero, err := tercerosHelper.GetNombreTerceroById(fmt.Sprintf("%v", data_["Funcionario"])); err == nil {
					Tercero_ = Tercero
					Terceros = append(Terceros, Tercero)
				} else {
					logs.Error(err)
					outputError = map[string]interface{}{
						"funcion": "/GetAllSolicitudes",
						"err":     err,
						"status":  "502",
					}
					return nil, outputError
				}
			} else {
				if keys := len(Terceros[0]); keys != 0 {
					if Tercero, err := utilsHelper.ArrayFind(Terceros, "Id", fmt.Sprintf("%v", data_["Funcionario"])); err == nil {
						if keys := len(Tercero); keys == 0 {
							if Tercero, err := tercerosHelper.GetNombreTerceroById2(fmt.Sprintf("%v", data_["Funcionario"])); err == nil {
								Tercero_ = Tercero
								Terceros = append(Terceros, Tercero)
							} else {
								logs.Error(err)
								outputError = map[string]interface{}{
									"funcion": "/GetAllSolicitudes",
									"err":     err,
									"status":  "502",
								}
								return nil, outputError
							}
						} else {
							Tercero_ = Tercero
						}
					} else {
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "/GetAllSolicitudes",
							"err":     err,
							"status":  "500",
						}
						return nil, outputError
					}
				} else {
					if Tercero, err := tercerosHelper.GetNombreTerceroById2(fmt.Sprintf("%v", data_["Revisor"])); err == nil {
						Tercero_ = Tercero
						Terceros = append(Terceros, Tercero)
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
			}

			if Terceros == nil {
				if Tercero, err := tercerosHelper.GetNombreTerceroById2(fmt.Sprintf("%v", data_["Revisor"])); err == nil {
					Revisor_ = Tercero
					Terceros = append(Terceros, Tercero)
				} else {
					logs.Error(err)
					outputError = map[string]interface{}{
						"funcion": "/GetAllSolicitudes",
						"err":     err,
						"status":  "502",
					}
					return nil, outputError
				}
			} else {
				if keys := len(Terceros[0]); keys != 0 {
					if Tercero, err := utilsHelper.ArrayFind(Terceros, "Id", fmt.Sprintf("%v", data_["Revisor"])); err == nil {
						if keys := len(Tercero); keys == 0 {
							if Tercero, err := tercerosHelper.GetNombreTerceroById2(fmt.Sprintf("%v", data_["Revisor"])); err == nil {
								Revisor_ = Tercero
								Terceros = append(Terceros, Tercero)
							} else {
								logs.Error(err)
								outputError = map[string]interface{}{
									"funcion": "/GetAllSolicitudes",
									"err":     err,
									"status":  "502",
								}
								return nil, outputError
							}
						} else {
							Revisor_ = Tercero
						}
					} else {
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "/GetAllSolicitudes",
							"err":     err,
							"status":  "500",
						}
						return nil, outputError
					}
				} else {
					if Tercero, err := tercerosHelper.GetNombreTerceroById2(fmt.Sprintf("%v", data_["Revisor"])); err == nil {
						Revisor_ = Tercero
						Terceros = append(Terceros, Tercero)
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
			}

			if Ubicaciones == nil {
				if ubicacion, err := ubicacionHelper.GetAsignacionSedeDependencia(fmt.Sprintf("%v", data_["Ubicacion"])); err == nil {
					fmt.Println(ubicacion)
					if keys := len(ubicacion); keys != 0 {
						Ubicacion_ = ubicacion
						Ubicaciones = append(Ubicaciones, ubicacion)
					}

				} else {
					logs.Error(err)
					outputError = map[string]interface{}{
						"funcion": "/GetAllSolicitudes",
						"err":     err,
						"status":  "502",
					}
					return nil, outputError
				}
			} else {
				if keys := len(Ubicaciones[0]); keys != 0 {
					if ubicacion, err := utilsHelper.ArrayFind(Ubicaciones, "Id", fmt.Sprintf("%v", data_["Ubicacion"])); err == nil {
						if keys := len(ubicacion); keys == 0 {
							if ubicacion, err := ubicacionHelper.GetAsignacionSedeDependencia(fmt.Sprintf("%v", data_["Ubicacion"])); err == nil {
								fmt.Println(ubicacion)
								if keys := len(ubicacion); keys != 0 {
									Ubicacion_ = ubicacion
									Ubicaciones = append(Ubicaciones, ubicacion)
								}
							} else {
								logs.Error(err)
								outputError = map[string]interface{}{
									"funcion": "/GetAllSolicitudes",
									"err":     err,
									"status":  "502",
								}
								return nil, outputError
							}
						} else {
							Ubicacion_ = ubicacion
						}
					} else {
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "/GetAllSolicitudes",
							"err":     err,
							"status":  "500",
						}
						return nil, outputError
					}
				} else {
					if ubicacion, err := ubicacionHelper.GetAsignacionSedeDependencia(fmt.Sprintf("%v", data_["Ubicacion"])); err == nil {
						fmt.Println(ubicacion)
						if keys := len(ubicacion); keys != 0 {
							Ubicacion_ = ubicacion
							Ubicaciones = append(Ubicaciones, ubicacion)
						}
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
				if Funcionario, err := tercerosHelper.GetNombreTerceroById2(fmt.Sprintf("%v", data_["Funcionario"])); err == nil {
					if Revisor, err := tercerosHelper.GetNombreTerceroById2(fmt.Sprintf("%v", data_["Revisor"])); err == nil {
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
						"status":  "502",
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
