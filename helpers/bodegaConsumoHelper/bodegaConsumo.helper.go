package bodegaConsumoHelper

import (
	"encoding/json"
	"fmt"

	// "strconv"
	"strings"
	// "reflect"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/tercerosHelper"
	"github.com/udistrital/utils_oas/request"
)

//GetTerceroById trae el nombre de un encargado por su id
func GetSolicitudById(id int) (Solicitud map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetSolicitudById", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var solicitud_ []map[string]interface{}
	var elementos___ []map[string]interface{}

	// url := "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento/" + fmt.Sprintf("%v", id) + ""
	url := "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento?query=Id:" + fmt.Sprintf("%v", id) + ""
	logs.Debug(url)
	if res, err := request.GetJsonTest(url, &solicitud_); err == nil && res.StatusCode == 200 {

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
		logs.Debug(fmt.Sprintf("str: %s", str))
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(str), &data); err == nil {

			fmt.Println(data)
			if tercero, err := tercerosHelper.GetNombreTerceroById(fmt.Sprintf("%v", data["Funcionario"])); err == nil {
				solicitud_[0]["Funcionario"] = tercero
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "/GetSolicitudById - tercerosHelper.GetNombreTerceroById(fmt.Sprintf(\"%v\", data[\"Funcionario\"]))",
					"err":     err,
					"status":  "502",
				}
				return nil, outputError
			}
			var data_ []map[string]interface{}
			if jsonString, err := json.Marshal(data["Elementos"]); err == nil {
				if err2 := json.Unmarshal(jsonString, &data_); err2 == nil {

					for _, elementos := range data_ {

						if Elemento__, err := TraerElementoSolicitud(elementos); err == nil {
							Elemento__["Cantidad"] = elementos["Cantidad"]
							fmt.Println(elementos["CantidadAprobada"])
							if elementos["CantidadAprobada"] != nil {
								Elemento__["CantidadAprobada"] = elementos["CantidadAprobada"]
							} else {
								Elemento__["CantidadAprobada"] = 0
							}

							elementos___ = append(elementos___, Elemento__)
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

func TraerElementoSolicitud(Elemento map[string]interface{}) (Elemento_ map[string]interface{}, err error) {

	urlcrud3 := "http://" + beego.AppConfig.String("oikos2Service") + "asignacion_espacio_fisico_dependencia?query=Id:" + fmt.Sprintf("%v", Elemento["Ubicacion"])

	var ubicacion []map[string]interface{}
	var sede []map[string]interface{}

	fmt.Println("elemento asdasdadasdfasd: ", Elemento)

	if _, err := request.GetJsonTest(urlcrud3, &ubicacion); err == nil {

		ubicacion2 := ubicacion[0]["EspacioFisicoId"].(map[string]interface{})

		z := strings.Split(fmt.Sprintf("%v", ubicacion2["CodigoAbreviacion"]), "")

		urlcrud4 := "http://" + beego.AppConfig.String("oikos2Service") + "espacio_fisico?query=CodigoAbreviacion:" + z[0] + z[1] + z[2] + z[3]

		if _, err := request.GetJsonTest(urlcrud4, &sede); err == nil {
			fmt.Println("Sede: ", sede)

		} else {
			panic(err.Error())
			return nil, err
		}
		if Elemento___, err := UltimoMovimientoKardex(fmt.Sprintf("%v", Elemento["ElementoActa"])); err == nil {

			Elemento___["Sede"] = sede[0]
			Elemento___["Dependencia"] = ubicacion[0]["DependenciaId"]
			Elemento___["Ubicacion"] = ubicacion[0]["EspacioFisicoId"]

			return Elemento___, nil

		} else {
			panic(err.Error())
			return nil, err
		}

	} else {
		panic(err.Error())
		return nil, err
	}

}

func GetElementosSinAsignar() (Elementos []map[string]interface{}, err error) {

	fmt.Println("aaaaaaaaaaaaaaaaaaaaa")
	url := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento?query=MovimientoId.FormatoTipoMovimientoId.Id:9,Activo:true&limit=-1"
	if _, err := request.GetJsonTest(url, &Elementos); err == nil {

		if keys := len(Elementos[0]); keys != 0 {
			for i, elemento := range Elementos {
				var detalle []map[string]interface{}
				url2 := "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento?query=Id:" + fmt.Sprintf("%v", elemento["ElementoActaId"]) + "&fields=Id,Nombre,Marca,Serie,SubgrupoCatalogoId"
				if _, err := request.GetJsonTest(url2, &detalle); err == nil {

					var subgrupo []map[string]interface{}
					url3 := "http://" + beego.AppConfig.String("catalogoElementosService") + "subgrupo?query=Id:" + fmt.Sprintf("%v", detalle[0]["SubgrupoCatalogoId"])
					if _, err := request.GetJsonTest(url3, &subgrupo); err == nil {

						Elementos[i]["Nombre"] = detalle[0]["Nombre"]
						Elementos[i]["Marca"] = detalle[0]["Marca"]
						Elementos[i]["Serie"] = detalle[0]["Serie"]
						Elementos[i]["SubgrupoCatalogoId"] = subgrupo[0]

					} else {
						panic(err.Error())
						return nil, err
					}
				} else {
					panic(err.Error())
					return nil, err
				}
			}

			return Elementos, nil
		} else {
			return Elementos, nil
		}

	} else {
		panic(err.Error())
		return nil, err
	}
}

func GetAperturasKardex() (Elementos []map[string]interface{}, err error) {

	var Elementos___ []map[string]interface{}
	url := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento?query=MovimientoId.FormatoTipoMovimientoId.CodigoAbreviacion:AP_KDX&limit=-1"
	if _, err := request.GetJsonTest(url, &Elementos___); err == nil {
		fmt.Println("Elementos___", Elementos___)

		if keys := len(Elementos___[0]); keys != 0 {
			for _, elemento := range Elementos___ {

				fmt.Println("Elemento", elemento)
				var data map[string]interface{}
				if jsonString, err := json.Marshal(elemento["MovimientoId"]); err == nil {
					fmt.Println("Movimiento", jsonString)

					if err2 := json.Unmarshal(jsonString, &data); err2 == nil {
						fmt.Println("DetalleMovimiento", data)

						str := fmt.Sprintf("%v", data["Detalle"])
						var data2 map[string]interface{}
						if err := json.Unmarshal([]byte(str), &data2); err == nil {
							fmt.Println("Detalle", data2)

							var elemento_catalogo []map[string]interface{}
							url3 := "http://" + beego.AppConfig.String("catalogoElementosService") + "elemento?query=Id:" + fmt.Sprintf("%v", elemento["ElementoCatalogoId"])
							if _, err := request.GetJsonTest(url3, &elemento_catalogo); err == nil {

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
								panic(err.Error())
								return nil, err
							}

						} else {
							panic(err.Error())
							return nil, err
						}
					} else {
						panic(err.Error())
						return nil, err
					}
				} else {
					panic(err.Error())
					return nil, err
				}

			}

			return Elementos, nil
		} else {
			return Elementos___, nil
		}

	} else {
		panic(err.Error())
		return nil, err
	}
}

func GetExistenciasKardex() (Elementos []map[string]interface{}, err error) {

	var Elementos___ []map[string]interface{}
	url := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento?query=MovimientoId.FormatoTipoMovimientoId.CodigoAbreviacion:AP_KDX,Activo:true&limit=-1&fields=ElementoCatalogoId"
	if _, err := request.GetJsonTest(url, &Elementos___); err == nil {
		fmt.Println("Elementos", Elementos___[0])

		if keys := len(Elementos___[0]); keys != 0 {

			for _, elemento := range Elementos___ {

				if Elemento, err := UltimoMovimientoKardex(fmt.Sprintf("%v", elemento["ElementoCatalogoId"])); err == nil {
					Elementos = append(Elementos, Elemento)
				} else {
					panic(err.Error())
					return nil, err
				}
			}

			return Elementos, nil
		} else {

			return Elementos___, nil
		}

	} else {
		panic(err.Error())
		return nil, err
	}
}

func UltimoMovimientoKardex(id_catalogo string) (Elemento_Movimiento map[string]interface{}, err error) {

	var elemento_catalogo []map[string]interface{}

	fmt.Println("id asdasdadasdfasd: ", id_catalogo)
	url3 := "http://" + beego.AppConfig.String("catalogoElementosService") + "elemento?query=Id:" + id_catalogo
	if _, err := request.GetJsonTest(url3, &elemento_catalogo); err == nil {

		fmt.Println(elemento_catalogo)
		var ultimo_movimiento_kardex []map[string]interface{}
		url4 := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento?query=ElementoCatalogoId:" +
			id_catalogo + ",Activo:true&limit=1&sortby=FechaCreacion&order=desc&fields=ElementoCatalogoId,Id,SaldoCantidad,SaldoValor"
		if _, err := request.GetJsonTest(url4, &ultimo_movimiento_kardex); err == nil {

			Elemento := ultimo_movimiento_kardex[0]
			Elemento["ElementoCatalogoId"] = elemento_catalogo[0]
			Elemento["Nombre"] = elemento_catalogo[0]["Nombre"]

			return Elemento, nil

		} else {
			panic(err.Error())
			return nil, err
		}

	} else {
		panic(err.Error())
		return nil, err
	}
}
