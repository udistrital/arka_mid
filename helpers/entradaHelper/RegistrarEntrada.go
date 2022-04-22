package entradaHelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// RegistrarEntrada Crea registro de entrada en estado en trámite
func RegistrarEntrada(data models.Movimiento) (result map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "RegistrarEntrada - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var (
		urlcrud             string
		res                 map[string]interface{}
		actaRecibido        models.TransaccionActaRecibido
		resEstadoMovimiento []models.EstadoMovimiento
	)
	resultado := make(map[string]interface{})

	detalleJSON := map[string]interface{}{}
	if err := json.Unmarshal([]byte(data.Detalle), &detalleJSON); err != nil {
		panic(err.Error())
	}

	ctxConsecutivo, _ := beego.AppConfig.Int("contxtEntradaCons")
	if consecutivo, consecutivoId, err := consecutivos.Get("%05.0f", ctxConsecutivo, "Registro Entrada Arka"); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "RegistrarEntrada - utilsHelper.GetConsecutivo(\"%05.0f\", ctxConsecutivo, \"Registro Entrada Arka\")",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	} else {
		consecutivo = consecutivos.Format(getTipoComprobanteEntradas()+"-", consecutivo, fmt.Sprintf("%s%04d", "-", time.Now().Year()))
		detalleJSON["consecutivo"] = consecutivo
		detalleJSON["ConsecutivoId"] = consecutivoId
		resultado["Consecutivo"] = detalleJSON["consecutivo"]
	}

	if jsonData, err := json.Marshal(detalleJSON); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "RegistrarEntrada - json.Marshal(detalleJSON)",
			"err":     err,
			"status":  "500",
		}
		return nil, outputError
	} else {
		data.Detalle = string(jsonData[:])
	}

	// Consulta el acta
	actaRecibidoId := int(detalleJSON["acta_recibido_id"].(float64))
	urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/" + strconv.Itoa(int(actaRecibidoId)) + "?elementos=false"
	if err := request.GetJson(urlcrud, &actaRecibido); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "RegistrarEntrada - request.GetJson(urlcrud, &actaRecibido)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	// Crea registro en api movimientos_arka_crud
	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "estado_movimiento?query=Nombre:Entrada%20En%20Trámite"
	if err := request.GetJson(urlcrud, &resEstadoMovimiento); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "RegistrarEntrada - request.GetJson(urlcrud, &resEstadoMovimiento)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	} else if len(resEstadoMovimiento) == 0 {
		err = errors.New("len(resEstadoMovimiento) == 0")
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "RegistrarEntrada - request.GetJson(urlcrud, &resEstadoMovimiento)",
			"err":     err,
			"status":  "404",
		}
		return nil, outputError
	}
	data.EstadoMovimientoId.Id = resEstadoMovimiento[0].Id

	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento"
	if err := request.SendJson(urlcrud, "POST", &res, &data); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "RegistrarEntrada - request.SendJson(urlcrud, \"POST\", &res, &data)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	resultado["MovimientoId"] = res["Id"]

	// Crea registro en table soporte_movimiento si es necesario
	if data.SoporteMovimientoId != 0 {
		urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "soporte_movimiento"

		idEntrada := int(res["Id"].(float64))
		movimientoEntrada := models.Movimiento{Id: idEntrada}
		soporteMovimiento := models.SoporteMovimiento{
			DocumentoId:  data.SoporteMovimientoId,
			Activo:       true,
			MovimientoId: &movimientoEntrada,
		}

		if err := request.SendJson(urlcrud, "POST", &res, &soporteMovimiento); err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "RegistrarEntrada - request.SendJson(urlcrud, \"POST\", &resS, &soporteMovimiento)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
	}

	if elementos, err := asignarPlacaActa(actaRecibidoId); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "RegistrarEntrada - asignarPlacaActa(actaRecibidoId)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	} else {
		actaRecibido.Elementos = elementos
	}

	// Actualiza el estado del acta
	urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/" + strconv.Itoa(int(actaRecibidoId))
	actaRecibido.UltimoEstado.EstadoActaId.Id = 6
	actaRecibido.UltimoEstado.Id = 0

	if err := request.SendJson(urlcrud, "PUT", &res, &actaRecibido); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "RegistrarEntrada - request.SendJson(urlcrud, \"PUT\", &res, &actaRecibido)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	return resultado, nil

}
