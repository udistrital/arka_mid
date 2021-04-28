package bodegaConsumoHelper

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	// "github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/arka_mid/helpers/tercerosHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
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
			if tercero, err := tercerosHelper.GetNombreTerceroById(fmt.Sprintf("%v", data["Funcionario"])); err == nil {
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

func TraerElementoSolicitud(Elemento map[string]interface{}) (Elemento_ map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "TraerElementoSolicitud - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var idStr int
	if id, err := strconv.Atoi(fmt.Sprintf("%v", Elemento["Ubicacion"])); err == nil {
		idStr = id
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "TraerElementoSolicitud - strconv.Atoi(fmt.Sprintf(\"%v\", Elemento[\"Ubicacion\"]))",
			"err":     err,
			"status":  "400",
		}
		return nil, outputError
	}

	urlcrud3 := "http://" + beego.AppConfig.String("oikos2Service") + "asignacion_espacio_fisico_dependencia?query=Id:" + strconv.Itoa(idStr)
	// logs.Debug("urlcrud3:", urlcrud3)

	var ubicacion []map[string]interface{}
	var sede []map[string]interface{}

	// fmt.Println("elemento asdasdadasdfasd: ", Elemento)

	if res, err := request.GetJsonTest(urlcrud3, &ubicacion); err == nil && res.StatusCode == 200 {

		ubicacion2 := ubicacion[0]["EspacioFisicoId"].(map[string]interface{})

		z := strings.Split(fmt.Sprintf("%v", ubicacion2["CodigoAbreviacion"]), "")

		urlcrud4 := "http://" + beego.AppConfig.String("oikos2Service") + "espacio_fisico?query=CodigoAbreviacion:" + z[0] + z[1] + z[2] + z[3]

		if res, err := request.GetJsonTest(urlcrud4, &sede); err != nil || res.StatusCode != 200 {
			if err == nil {
				err = fmt.Errorf("Undesired Status Code: %d", res.StatusCode)
			}
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "TraerElementoSolicitud - request.GetJsonTest(urlcrud4, &sede)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}

		var idElemento int
		if id, err := strconv.Atoi(fmt.Sprintf("%v", Elemento["ElementoActa"])); err == nil {
			idElemento = id
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "TraerElementoSolicitud - strconv.Atoi(fmt.Sprintf(\"%v\", Elemento[\"ElementoActa\"]))",
				"err":     err,
				"status":  "400",
			}
			return nil, outputError
		}
		// logs.Debug("elemActa:", elemActa)
		if Elemento___, err := UltimoMovimientoKardex(idElemento); err == nil {

			Elemento___["Sede"] = sede[0]
			Elemento___["Dependencia"] = ubicacion[0]["DependenciaId"]
			Elemento___["Ubicacion"] = ubicacion[0]["EspacioFisicoId"]

			return Elemento___, nil

		} else {
			return nil, err
		}

	} else {
		if err == nil {
			err = fmt.Errorf("Undesired Status Code: %d", res.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "TraerElementoSolicitud - request.GetJsonTest(urlcrud3, &ubicacion)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

}

func GetElementosSinAsignar() (Elementos []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetElementosSinAsignar - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	// fmt.Println("aaaaaaaaaaaaaaaaaaaaa")
	url := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento?limit=-1"
	url += "&query=Activo:true,MovimientoId.FormatoTipoMovimientoId.Id:9"
	// logs.Debug("url:", url)
	if res, err := request.GetJsonTest(url, &Elementos); err == nil && res.StatusCode == 200 {

		if keys := len(Elementos[0]); keys != 0 {

			elementosActaBuffer := make(map[int]interface{})
			subgruposCatalogoBuffer := make(map[int]interface{})

			for i, elemento := range Elementos {
				void := true

				var elementoActaId int
				if v, err := strconv.Atoi(fmt.Sprintf("%v", elemento["ElementoActaId"])); err == nil && v > 0 {
					elementoActaId = v
				} else {
					err = fmt.Errorf("ElementoActaId='%v', erroneo para 'elementos_movimiento.Id=%v'", elemento["ElementoActaId"], elemento["Id"])
					logs.Warn(err)
					// TODO: revisar si esto es suficiente
					continue
				}

				urlElemento := "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento"
				urlElemento += "?query=Id:" + strconv.Itoa(elementoActaId)
				urlElemento += "&fields=Id,Nombre,Marca,Serie,SubgrupoCatalogoId"
				if detalle, err := utilsHelper.BufferGet(elementoActaId, elementosActaBuffer, urlElemento); err == nil && detalle != nil {

					var subgrupoCatalogoId int
					if v, err := strconv.Atoi(fmt.Sprintf("%v", detalle["SubgrupoCatalogoId"])); err == nil && v > 0 {
						subgrupoCatalogoId = v
					} else {
						err = fmt.Errorf("SubgrupoCatalogoId='%v', erroneo para 'elemento(Acta).Id=%d'", detalle["SubgrupoCatalogoId"], elementoActaId)
						logs.Warn(err)
						// TODO: revisar si esto es suficiente
						continue
					}

					urlSubgrupo := "http://" + beego.AppConfig.String("catalogoElementosService") + "subgrupo"
					urlSubgrupo += "?query=Id:" + strconv.Itoa(subgrupoCatalogoId)
					if subgrupo, err := utilsHelper.BufferGet(subgrupoCatalogoId, subgruposCatalogoBuffer, urlSubgrupo); err == nil && subgrupo != nil {
						Elementos[i]["Nombre"] = detalle["Nombre"]
						Elementos[i]["Marca"] = detalle["Marca"]
						Elementos[i]["Serie"] = detalle["Serie"]
						Elementos[i]["SubgrupoCatalogoId"] = subgrupo

						void = false
					} else {
						if err == nil {
							logs.Warn(fmt.Errorf("no hay subgrupo_catalogo.Id=%d (CRUD catalogo) asociado al elemento.Id=%d (CRUD Actas)", subgrupoCatalogoId, elementoActaId))
						} else {
							logs.Warn(err)
						}
					}

				} else {
					if err == nil {
						logs.Warn(fmt.Errorf("no hay elemento.Id=%d (CRUD Actas) asociado al elemento.Id=%v (CRUD movimientos_arka)", elementoActaId, elemento["Id"]))
					} else {
						logs.Warn(err)
					}
				}

				if void {
					Elementos[i] = nil
				}
			}

			// Quitar elementos inconsistentes
			fin := len(Elementos)
			// logs.Debug("fin(antes):", fin)
			for i := 0; i < fin; {
				if Elementos[i] != nil {
					i++
				} else {
					Elementos[i] = Elementos[fin-1]
					fin--
				}
			}
			// logs.Debug("fin(despues):", fin)

		}
		return Elementos, nil

	} else {
		if err == nil {
			err = fmt.Errorf("Undesired Status Code: %d", res.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetElementosSinAsignar - request.GetJsonTest(url, &Elementos)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}

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

func GetExistenciasKardex() (Elementos []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetExistenciasKardex - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var Elementos___ []map[string]interface{}
	url := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento?query=MovimientoId.FormatoTipoMovimientoId.CodigoAbreviacion:AP_KDX,Activo:true&limit=-1&fields=ElementoCatalogoId"
	if res, err := request.GetJsonTest(url, &Elementos___); err == nil && res.StatusCode == 200 {
		// fmt.Println("Elementos", Elementos___[0])

		if keys := len(Elementos___[0]); keys != 0 {

			for _, elemento := range Elementos___ {

				var idCatalogo int
				if id, err := strconv.Atoi(fmt.Sprintf("%v", elemento["ElementoCatalogoId"])); err == nil {
					idCatalogo = id
				} else {
					logs.Warn(err)
					continue
				}

				if Elemento, err := UltimoMovimientoKardex(idCatalogo); err == nil {
					if s, ok := Elemento["SaldoCantidad"]; ok {
						if v, ok := s.(float64); ok && v > 0 {
							Elementos = append(Elementos, Elemento)
						}
					}
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
			"funcion": "GetExistenciasKardex - request.GetJsonTest(url, &Elementos___)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}

func UltimoMovimientoKardex(id_catalogo int) (Elemento_Movimiento map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "UltimoMovimientoKardex - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	if id_catalogo <= 0 {
		err := fmt.Errorf("id MUST be > 0")
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "UltimoMovimientoKardex - id_catalogo <= 0",
			"err":     err,
			"status":  "400",
		}
		return nil, outputError
	}

	idStr := strconv.Itoa(id_catalogo)

	var elemento_catalogo []map[string]interface{}

	// fmt.Println("id asdasdadasdfasd: ", id_catalogo)
	url3 := "http://" + beego.AppConfig.String("catalogoElementosService") + "elemento?query=Id:" + idStr
	// logs.Debug("url3:", url3)
	if res, err := request.GetJsonTest(url3, &elemento_catalogo); err == nil && res.StatusCode == 200 {

		if len(elemento_catalogo) != 1 || len(elemento_catalogo[0]) == 0 {
			err = fmt.Errorf("No hay un elemento del Catalogo de Elementos con id:%s", idStr)
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "UltimoMovimientoKardex - len(elemento_catalogo) != 1 || len(elemento_catalogo[0]) == 0",
				"err":     err,
				"status":  "404",
			}
			return nil, outputError
		}

		// fmt.Println(elemento_catalogo)
		var ultimo_movimiento_kardex []map[string]interface{}
		url4 := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento?query=ElementoCatalogoId:" +
			idStr + ",Activo:true&limit=1&sortby=FechaCreacion&order=desc&fields=ElementoCatalogoId,Id,SaldoCantidad,SaldoValor"
		// logs.Debug("url4:", url4)
		if res, err := request.GetJsonTest(url4, &ultimo_movimiento_kardex); err == nil && res.StatusCode == 200 {

			Elemento := ultimo_movimiento_kardex[0]
			Elemento["ElementoCatalogoId"] = elemento_catalogo[0]
			Elemento["Nombre"] = elemento_catalogo[0]["Nombre"]

			return Elemento, nil

		} else {
			if err == nil {
				err = fmt.Errorf("Undesired Status Code: %d", res.StatusCode)
			}
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "UltimoMovimientoKardex - request.GetJsonTest(url4, &ultimo_movimiento_kardex)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}

	} else {
		if err == nil {
			err = fmt.Errorf("Undesired Status Code: %d", res.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "UltimoMovimientoKardex - request.GetJsonTest(url3, &elemento_catalogo)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}
