package bajasHelper

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/salidaHelper"
	"github.com/udistrital/arka_mid/models"
	e "github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

func TraerDatosElemento(id int) (Elemento models.ElementoBajaDetalle, outputError map[string]interface{}) {
	const funcion = "TraerDatosElemento - "
	defer e.ErrorControlFunction(funcion+"unhandled error!", fmt.Sprint(http.StatusInternalServerError))

	var elemento_movimiento_ []map[string]interface{}
	// var movimiento_ map[string]interface{}

	url := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento/?query=ElementoActaId:" + fmt.Sprintf("%v", id) + ",Activo:true"
	if _, err := request.GetJsonTest(url, &elemento_movimiento_); err == nil {
		logs.Debug("len(elemento_movimiento_):", len(elemento_movimiento_))
		if len(elemento_movimiento_) == 0 || len(elemento_movimiento_[0]) == 0 {
			return Elemento, e.Error(funcion+"len(elemento_movimiento_) == 0 || len(elemento_movimiento_[0]) == 0",
				fmt.Errorf("no se encontraron elementos_movimiento con ElementoActaId:%d", id), fmt.Sprint(http.StatusNotFound))
		}

		if v, err := salidaHelper.TraerDetalle(elemento_movimiento_[0]["MovimientoId"]); err == nil {

			fmt.Println("Elemento Movimiento: ", elemento_movimiento_)

			var movimiento_ map[string]interface{}
			if jsonString3, err := json.Marshal(elemento_movimiento_[0]["MovimientoId"]); err == nil {
				if err2 := json.Unmarshal(jsonString3, &movimiento_); err2 == nil {
					movimiento_["MovimientoPadreId"] = nil
				}
			}

			var elemento_ []map[string]interface{}

			urlcrud_ := "http://" + beego.AppConfig.String("actaRecibidoService") +
				"elemento?query=Id:" + fmt.Sprintf("%v", elemento_movimiento_[0]["ElementoActaId"]) +
				"&fields=Id,Nombre,Marca,Serie,Placa,SubgrupoCatalogoId"
				// "&fields=Id,Nombre,TipoBienId,Marca,Serie,Placa,SubgrupoCatalogoId"
			logs.Debug("urlcrud_:", urlcrud_)
			if _, err := request.GetJsonTest(urlcrud_, &elemento_); err == nil {
				logs.Debug("len(elemento_):", len(elemento_))

				fmt.Println("Elemento: ", elemento_)

				var subgrupo_ map[string]interface{}
				urlcrud_2 := "http://" + beego.AppConfig.String("catalogoElementosService") + "subgrupo/" + fmt.Sprintf("%v", elemento_[0]["SubgrupoCatalogoId"])
				if _, err := request.GetJsonTest(urlcrud_2, &subgrupo_); err == nil {
					Elemento = models.ElementoBajaDetalle{
						Id:                 elemento_[0]["Id"],
						Placa:              elemento_[0]["Placa"],
						Nombre:             elemento_[0]["Nombre"],
						Entrada:            v["MovimientoPadreId"],
						Salida:             movimiento_,
						SubgrupoCatalogoId: subgrupo_,
						Marca:              elemento_[0]["Marca"],
						Serie:              elemento_[0]["Serie"],
						Funcionario:        v["Funcionario"],
						Sede:               v["Sede"],
						Dependencia:        v["Dependencia"],
						Ubicacion:          v["Ubicacion"],
						// TipoBienId:         elemento_[0]["TipoBienId"],
					}
				} else {
					logs.Error(err)
					outputError = map[string]interface{}{
						"funcion": "/TraerDatosElemento - request.GetJsonTest(urlcrud_2, &subgrupo_)",
						"err":     err,
						"status":  "502",
					}
				}
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "/TraerDatosElemento - request.GetJsonTest(urlcrud_, &elemento_)",
					"err":     err,
					"status":  "502",
				}
			}
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "/TraerDatosElemento - salidaHelper.TraerDetalle(elemento_movimiento_[0][\"MovimientoId\"])",
				"err":     err,
				"status":  "502",
			}
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/TraerDatosElemento - request.GetJsonTest(url, &elemento_movimiento_)",
			"err":     err,
			"status":  "502",
		}
	}
	return
}
