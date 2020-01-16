package salidaHelper

import (
	"encoding/json"
	"strconv"
	"fmt"
	// "reflect"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// AddEntrada Transacción para registrar la información de una salida
func AddSalida(data *models.TrSalida) map[string]interface{} {
	var (
		urlcrud   string
		res       map[string]interface{}
		resM      map[string]interface{}
		resultado map[string]interface{}
	)

	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "tr_salida/"

	// Inserta salida en Movimientos ARKA
	if err := request.SendJson(urlcrud, "POST", &res, &data); err == nil {
		// Inserta salida en Movimientos KRONOS
		urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo"

		var salidaId int

		dataSalida := res["Salida"].(map[string]interface{})
		salidaId = int(dataSalida["Id"].(float64))

		procesoExterno := int64(salidaId)
		logs.Debug(procesoExterno)
		tipo := models.TipoMovimiento{Id: 16}
		movimientosKronos := models.MovimientoProcesoExterno{
			TipoMovimientoId: &tipo,
			ProcesoExterno:   procesoExterno,
			Activo:           true,
		}

		if err = request.SendJson(urlcrud, "POST", &resM, &movimientosKronos); err == nil {
			body := res
			body["MovimientosKronos"] = resM["Body"]
			resultado = body
		} else {
			panic(err.Error())
		}
	} else {
		panic(err.Error())
	}

	return resultado
}

func GetSalida(id int) (Salida map[string]interface{}, err error) {

	
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
							panic(err.Error())
							return nil, err
						}
					} else {
						panic(err.Error())
						return nil, err
					}

					if _, err := request.GetJsonTest(urlcrud, &salida_); err == nil {
					
					} else {
						panic(err.Error())
						return nil, err
					}
					

				} 


			} else {
				panic(err.Error())
				return nil, err
			}
		} else {
			panic(err.Error())
			return nil, err
		}

		if salida__, err := TraerDetalle(salida_["Salida"]); err == nil {

			Salida_final := map[string]interface{}{
				"Elementos":	data_,
				"Salida":		salida__,
			}
			return Salida_final, nil

		} else {
			panic(err.Error())
			return nil, err
		}
		
	} else {
		panic(err.Error())
		return nil, err
	}
}

func GetSalidas() (Salidas []map[string]interface{}, err error) {


	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento?query=FormatoTipoMovimientoId.CodigoAbreviacion__contains:SAL,Activo:true"

	

	var salidas_ []map[string]interface{}
	if _, err := request.GetJsonTest(urlcrud, &salidas_); err == nil {
	
		for _, salida := range salidas_ {
			
			if salida__, err := TraerDetalle(salida); err == nil {

				Salidas = append(Salidas, salida__)
				
	
			} else {
				panic(err.Error())
				return nil, err
			}
		}
		return Salidas, nil
	
	} else {
		panic(err.Error())
		return nil, err
	}
}

func TraerDetalle(salida interface{}) (salida_ map[string]interface{}, err error) {

	var data map[string]interface{}
	if jsonString, err := json.Marshal(salida); err == nil {
		if err2 := json.Unmarshal(jsonString, &data); err2 == nil {

			str := fmt.Sprintf("%v", data["Detalle"])

			var data2 map[string]interface{}

			if err := json.Unmarshal([]byte(str), &data2); err == nil {

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

				if data2["funcionario"] != nil {

					urlcrud2 := "http://" + beego.AppConfig.String("tercerosService") + "tercero/?query=Id:" + fmt.Sprintf("%v", data2["funcionario"]) + "&fields=Id,NombreCompleto" 
					if _, err := request.GetJsonTest(urlcrud2, &tercero); err != nil {
						panic(err.Error())
						return nil, err
					}

					Salida2 := map[string]interface{}{
						"Id":							data["Id"],
						"Observacion":					data["Observacion"],
						"Funcionario":					tercero[0],
						"Sede":							sede[0],
						"Dependencia":					ubicacion[0]["DependenciaId"],
						"Ubicacion":					ubicacion[0]["EspacioFisicoId"],
						"FechaCreacion":				data["FechaCreacion"],
						"FechaModificacion":			data["FechaModificacion"],
						"Activo":						data["Activo"],
						"MovimientoPadreId":			data["MovimientoPadreId"],
						"FormatoTipoMovimientoId":		data["FormatoTipoMovimientoId"],
						"EstadoMovimientoId":			data["EstadoMovimientoId"],

					}
					
					return Salida2, nil
					
				} else {

					Salida2 := map[string]interface{}{
						"Id":							data["Id"],
						"Observacion":					data["Observacion"],
						"Sede":							sede[0],
						"Dependencia":					ubicacion[0]["DependenciaId"],
						"Ubicacion":					ubicacion[0]["EspacioFisicoId"],
						"FechaCreacion":				data["FechaCreacion"],
						"FechaModificacion":			data["FechaModificacion"],
						"Activo":						data["Activo"],
						"MovimientoPadreId":			data["MovimientoPadreId"],
						"FormatoTipoMovimientoId":		data["FormatoTipoMovimientoId"],
						"EstadoMovimientoId":			data["EstadoMovimientoId"],

					}
					return Salida2, nil
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