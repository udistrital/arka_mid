package bodegaConsumoHelper

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"

	"github.com/udistrital/arka_mid/helpers/actaRecibidoHelper"
	"github.com/udistrital/utils_oas/request"
)

//GetTerceroById trae el nombre de un encargado por su id
func GetSolicitudById(id int) (solicitudRespuesta map[string]interface{}, outputError map[string]interface{}) {
	var url string
	var solicitud []map[string]interface{}
	var detalleMovimiento map[string]interface{}
	var elementos []map[string]interface{}

	fmt.Println("prueba  ", strconv.Itoa(id))
	url = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento/?query=Id:" + strconv.Itoa(id) + ",Activo:true"
	// se trae un movimiento solicitud bodega
	if response, err := request.GetJsonTest(url, &solicitud); err == nil {
		if response.StatusCode == 200 {
			for _, movSolicitud := range solicitud {
				if len(movSolicitud) == 0 {
					outputError = map[string]interface{}{"Function": "GetSolicitudById", "Cause": "No se encontro registro", "Error": err}
					return nil, outputError

				} else {
					// se abre el json interno del movimiento
					if err := json.Unmarshal([]byte(movSolicitud["Detalle"].(string)), &detalleMovimiento); err == nil {
						fmt.Println(detalleMovimiento, "detalle2: ", detalleMovimiento["Funcionario"])
						fmt.Println("elementossss: ", detalleMovimiento["Elementos"])
						// se recorren los elementos existentes en el detalle
						for _, elementoMovimiento := range detalleMovimiento["Elementos"].([]interface{}) {
							fmt.Println("elemento id: ", elementoMovimiento.(map[string]interface{})["Cantidad"])
							cant := elementoMovimiento.(map[string]interface{})["ElementoActa"]
							// se llama función para traer informacion de elemento especifico de esquema movimientos
							if elementoMovimiento, outputError := GetElementoMovimientoById(strconv.Itoa(int(elementoMovimiento.(map[string]interface{})["ElementoActa"].(float64)))); outputError == nil {
								//fmt.Println("Movimiento elem: ", elementoMovimiento)
								// se llama función para traer informacion de elemento especifico de esquema actas
								if elementoActa, outputError := actaRecibidoHelper.GetElementoById(strconv.Itoa(int(elementoMovimiento["Id"].(float64)))); outputError == nil {
									//	fmt.Println("acta elem: ", elementoActa)

									elemento := map[string]interface{}{
										"Id":                 elementoMovimiento["Id"],
										"Nombre":             elementoActa["Nombre"],
										"Marca":              elementoActa["Marca"],
										"Serie":              elementoActa["Serie"],
										"CantidadDisponible": elementoMovimiento["SaldoCantidad"],
										"CantidadSolicitada": cant,
										"ValorUnitario":      elementoActa["ValorUnitario"],
									}
									fmt.Println("prue: ", elemento)
									elementos = append(elementos, elemento)
								}
							}
						}
						return map[string]interface{}{
							"Id":            movSolicitud["Id"],
							"FechaCreacion": movSolicitud["FechaCreacion"],
							"Observacion":   movSolicitud["Observacion"],
							"Elementos":     elementos,
						}, nil
					} else {
						outputError = map[string]interface{}{"Function": "GetSolicitudById", "Error": err}
						return nil, outputError
					}

				}

			}
		} else if response.StatusCode == 400 {
			outputError = map[string]interface{}{"Function": "GetSolicitudById", "Error": err}
			return nil, outputError
		}
	} else {
		fmt.Println("error: ", err)
		outputError = map[string]interface{}{"Function": "GetSolicitudById", "Error": err}
		return nil, outputError
	}

	return

}

func GetElementoMovimientoById(Id string) (Elemento map[string]interface{}, outputError map[string]interface{}) {
	var urlelemento string
	var elemento []map[string]interface{}
	urlelemento = "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento/?query=Id:" + Id + "&limit=1"
	if response, err := request.GetJsonTest(urlelemento, &elemento); err == nil {

		if response.StatusCode == 200 {
			for _, element := range elemento {
				if len(element) == 0 {
					return nil, map[string]interface{}{"Function": "GetElementoById", "Error": "Sin Elementos"}
				} else {
					return element, nil
				}

			}
		}
	} else {
		return nil, map[string]interface{}{"Function": "GetElementoById", "Error": err}
	}
	return
}

func GetElementosSinAsignar() (Elementos []map[string]interface{}, err error) {

	url := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento?query=MovimientoId.FormatoTipoMovimientoId.Id:9,Activo:true&limit=-1"
	if _, err := request.GetJsonTest(url, &Elementos); err == nil {

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
		panic(err.Error())
		return nil, err
	}
}

func GetAperturasKardex() (Elementos []map[string]interface{}, err error) {

	var Elementos___ []map[string]interface{}
	url := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento?query=MovimientoId.FormatoTipoMovimientoId.CodigoAbreviacion:AP_KDX,Activo:true&limit=-1"
	if _, err := request.GetJsonTest(url, &Elementos___); err == nil {
		fmt.Println("Elementos",Elementos___)

		for _, elemento := range Elementos___ {

			fmt.Println("Elemento", elemento)
			var data map[string]interface{}
			if jsonString, err := json.Marshal(elemento["MovimientoId"]); err == nil {
				fmt.Println("Movimiento",jsonString)

				if err2 := json.Unmarshal(jsonString, &data); err2 == nil {
					fmt.Println("DetalleMovimiento", data)

					str := fmt.Sprintf("%v", data["Detalle"])
					var data2 map[string]interface{}
					if err := json.Unmarshal([]byte(str), &data2); err == nil {
						fmt.Println("Detalle",data2)

						
						var elemento_catalogo []map[string]interface{}
						url3 := "http://" + beego.AppConfig.String("catalogoElementosService") + "elemento?query=Id:" + fmt.Sprintf("%v", elemento["ElementoCatalogoId"])
						if _, err := request.GetJsonTest(url3, &elemento_catalogo); err == nil {

							Elemento := map[string]interface{}{
								"MetodoValoracion":		data2["Metodo_Valoracion"],
								"CantidadMinima":		data2["Cantidad_Minima"],
								"CantidadMaxima": 		data2["Cantidad_Maxima"],
								"FechaCreacion":		elemento["FechaCreacion"],
								"Observaciones":		data["Observacion"],
								"Id":					data["Id"],
								"MovimientoPadreId":	data["MovimientoPadreId"],
								"ElementoCatalogoId":	elemento_catalogo[0],
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
		panic(err.Error())
		return nil, err
	}
}
