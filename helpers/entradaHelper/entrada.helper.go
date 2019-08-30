package entradaHelper

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// AddEntrada Transacción para registrar la información de una entrada
func AddEntrada(data models.EntradaElemento) map[string]interface{} {
	var (
		urlcrud      string
		res          interface{}
		resA         interface{}
		actaRecibido []models.TransaccionActaRecibido
		resultado    map[string]interface{}
	)

	urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/"

	// Solicita información acta
	if err := request.GetJson(urlcrud+strconv.Itoa(int(data.ActaRecibidoId)), &actaRecibido); err == nil {
		//Envia información entrada
		urlcrud = "http://" + beego.AppConfig.String("entradaService") + "entrada_elemento"

		if err = request.SendJson(urlcrud, "POST", &res, &data); err == nil {
			// Cambia estado acta
			switch res.(type) {
			case map[string]interface{}:
				urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/" + strconv.Itoa(int(data.ActaRecibidoId))
				actaRecibido[0].UltimoEstado.EstadoActaId.Id = 6
				actaRecibido[0].UltimoEstado.Id = 0

				if err = request.SendJson(urlcrud, "PUT", &resA, &actaRecibido[0]); err == nil {

					switch resA.(type) {
					case map[string]interface{}:
						body := res.(map[string]interface{})
						body["Acta"] = resA
						resultado = body
					default:
						beego.Error("res acta", resA)
						panic(helpers.ExternalAPIErrorMessage())
					}
				} else {
					panic(err.Error())
				}
			default:
				beego.Error("res entrada", res)
				panic(helpers.ExternalAPIErrorMessage())
			}

		} else {
			panic(err.Error())
		}

	} else {
		panic(err.Error())
	}

	return resultado
}
