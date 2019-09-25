package entradaHelper

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
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

// GetEntrada ...
func GetEntrada(entradaId int) (consultaEntrada *models.ConsultaEntrada, outputError map[string]interface{}) {
	var (
		urlcrud  string
		entrada  models.EntradaElemento
		contrato models.Contrato
	)

	if entradaId != 0 { // (1) error parametro
		// Solicita información elementos acta
		urlcrud = "http://" + beego.AppConfig.String("entradaService") + "entrada_elemento/" + strconv.Itoa(entradaId)

		if err := request.GetJson(urlcrud+strconv.Itoa(int(entradaId)), &entrada); err == nil {

			logs.Debug(urlcrud)

			logs.Debug(entrada)

			urlcrud = "http://" + beego.AppConfig.String("administrativaService") + "informacion_contrato/" + strconv.Itoa(entrada.ContratoId) + "/" + entrada.Vigencia

			logs.Debug(urlcrud)

			if response, err := request.GetJsonTest(urlcrud, &contrato); err == nil { // (2) error servicio caido

				if response.StatusCode == 200 { // (3) error estado de la solicitud

					consultaEntrada.Id = entrada.Id
					consultaEntrada.Solicitante = entrada.Solicitante
					consultaEntrada.Observacion = entrada.Observacion
					consultaEntrada.Importacion = entrada.Importacion
					consultaEntrada.FechaCreacion = entrada.FechaCreacion
					consultaEntrada.FechaModificacion = entrada.FechaModificacion
					consultaEntrada.Activo = entrada.Activo
					consultaEntrada.TipoEntradaId = entrada.TipoEntradaId
					// CONTRATO
					consultaEntrada.ContratoId.NumeroContratoSuscrito = contrato.NumeroContratoSuscrito
					consultaEntrada.ContratoId.OrdenadorGasto = contrato.OrdenadorGasto
					consultaEntrada.ContratoId.Supervisor = contrato.Supervisor

					consultaEntrada.ElementoId = entrada.ElementoId
					consultaEntrada.DocumentoContableId = entrada.DocumentoContableId
					consultaEntrada.Consecutivo = entrada.Consecutivo
					consultaEntrada.Vigencia = entrada.Vigencia

					return consultaEntrada, nil

				} else {
					logs.Info("Error (3) estado de la solicitud")
					outputError = map[string]interface{}{"Function": "GetEntrada:GetEntrada", "Error": response.Status}
					return nil, outputError
				}
			} else {
				logs.Info("Error (2) servicio caido")
				logs.Debug(err)
				outputError = map[string]interface{}{"Function": "GetEntrada", "Error": err}
				return nil, outputError
			}
		} else {
			return nil, outputError
		}
	} else {
		return nil, outputError
	}
}
