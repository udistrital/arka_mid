package bajasHelper

import (
	"encoding/json"
	"fmt"
	// "strconv"
	// "strings"
	// "reflect"

	"github.com/astaxie/beego"

	// "github.com/udistrital/arka_mid/helpers/tercerosHelper"
	"github.com/udistrital/arka_mid/helpers/salidaHelper"
	"github.com/udistrital/utils_oas/request"
)

func TraerDatosElemento(id int) (Elemento map[string]interface{}, err error) {

	var elemento_movimiento_ []map[string]interface{}
	// var movimiento_ map[string]interface{}

	url := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento/?query=ElementoActaId:"+ fmt.Sprintf("%v", id) + ",Activo:true"
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
						"Id":						elemento_[0]["Id"],
						"Placa":					elemento_[0]["Placa"],
						"Nombre":					elemento_[0]["Nombre"],
						"TipoBienId":				elemento_[0]["TipoBienId"],
						"Entrada":					v["MovimientoPadreId"],
						"Salida":					movimiento_,
						"SubgrupoCatalogoId":		subgrupo_,
						"Marca":					elemento_[0]["Marca"],
						"Serie":					elemento_[0]["Serie"],
						"Funcionario":				v["Funcionario"],
						"Sede":						v["Sede"],
						"Dependencia":				v["Dependencia"],
						"Ubicacion":				v["Ubicacion"],
					}


					// elemento_[0]["SubgrupoCatalogoId"] = subgrupo_
					// elemento_movimiento_[0]["ElementoActaId"] = elemento_[0]
					// Elemento = elemento_movimiento_[0]
					return Elemento, nil

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