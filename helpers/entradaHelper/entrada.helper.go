package entradaHelper

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// AddEntrada Transacción para registrar la información de una entrada
func AddEntrada(data models.Movimiento) map[string]interface{} {
	var (
		urlcrud      string
		res          map[string]interface{}
		resA         map[string]interface{}
		resM         map[string]interface{}
		resS         map[string]interface{}
		actaRecibido []models.TransaccionActaRecibido
		resultado    map[string]interface{}
	)

	urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/"

	detalleJSON := map[string]interface{}{}
	if err := json.Unmarshal([]byte(data.Detalle), &detalleJSON); err != nil {
		panic(err.Error())
	}

	// Solicita información acta
	actaRecibidoId, err := strconv.Atoi(detalleJSON["acta_recibido_id"].(string))

	if err != nil {
		panic(err.Error())
	}

	if err := request.GetJson(urlcrud+strconv.Itoa(int(actaRecibidoId)), &actaRecibido); err == nil {
		// Envia información entrada
		urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento"

		if err = request.SendJson(urlcrud, "POST", &res, &data); err == nil {
			// Si la entrada tiene soportes
			if data.SoporteMovimientoId != 0 {
				// Envia información soporte (Si tiene)
				urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "soporte_movimiento"

				idEntrada := int(res["Id"].(float64))

				movimientoEntrada := models.Movimiento{Id: idEntrada}
				soporteMovimiento := models.SoporteMovimiento{
					DocumentoId:  data.SoporteMovimientoId,
					Activo:       true,
					MovimientoId: &movimientoEntrada,
				}

				if err = request.SendJson(urlcrud, "POST", &resS, &soporteMovimiento); err != nil {
					panic(err.Error())
				}
			}

			// Envia información movimientos Kronos
			urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo"

			procesoExterno := int64(res["Id"].(float64))
			tipo := models.TipoMovimiento{Id: data.IdTipoMovimiento}
			movimientosKronos := models.MovimientoProcesoExterno{
				TipoMovimientoId: &tipo,
				ProcesoExterno:   procesoExterno,
				Activo:           true,
			}

			if err = request.SendJson(urlcrud, "POST", &resM, &movimientosKronos); err == nil {
				// Cambia estado acta
				urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/" + strconv.Itoa(int(actaRecibidoId))
				actaRecibido[0].UltimoEstado.EstadoActaId.Id = 6
				actaRecibido[0].UltimoEstado.Id = 0

				if err = request.SendJson(urlcrud, "PUT", &resA, &actaRecibido[0]); err == nil {
					body := res
					body["Acta"] = resA
					resultado = body
				} else {
					panic(err.Error())
				}
			} else {
				panic(err.Error())
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

			urlcrud = "http://" + beego.AppConfig.String("administrativaService") + "informacion_contrato/" + strconv.Itoa(entrada.ContratoId) + "/" + entrada.Vigencia

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
