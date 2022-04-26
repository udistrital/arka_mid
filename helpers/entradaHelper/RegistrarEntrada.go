package entradaHelper

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

// RegistrarEntrada Crea registro de entrada en estado en trámite
func RegistrarEntrada(data *models.TransaccionEntrada) (result map[string]interface{}, outputError map[string]interface{}) {

	funcion := "RegistrarEntrada"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		urlcrud          string
		res              map[string]interface{}
		actaRecibido     models.TransaccionActaRecibido
		tipoMovimiento   int
		estadoMovimiento int
	)

	resultado := make(map[string]interface{})

	detalleJSON := map[string]interface{}{}
	if err := json.Unmarshal([]byte(data.Detalle), &detalleJSON); err != nil {
		panic(err.Error())
	}

	ctxConsecutivo, _ := beego.AppConfig.Int("contxtEntradaCons")
	if consecutivo, consecutivoId, err := consecutivos.Get("%05.0f", ctxConsecutivo, "Entradas Arka"); err != nil {
		return nil, outputError
	} else {
		consecutivo = consecutivos.Format(getTipoComprobanteEntradas()+"-", consecutivo, fmt.Sprintf("%s%04d", "-", time.Now().Year()))
		detalleJSON["consecutivo"] = consecutivo
		detalleJSON["ConsecutivoId"] = consecutivoId
		resultado["Consecutivo"] = detalleJSON["consecutivo"]
	}

	if jsonData, err := json.Marshal(detalleJSON); err != nil {
		logs.Error(err)
		eval := " - json.Marshal(detalleJSON)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	} else {
		data.Detalle = string(jsonData[:])
	}

	// Consulta el acta
	actaRecibidoId := int(detalleJSON["acta_recibido_id"].(float64))
	urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/" + strconv.Itoa(int(actaRecibidoId)) + "?elementos=false"
	if err := request.GetJson(urlcrud, &actaRecibido); err != nil {
		logs.Error(err)
		eval := " - request.GetJson(urlcrud, &actaRecibido)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	// Crea registro en api movimientos_arka_crud
	if err := movimientosArka.GetEstadoMovimientoIdByNombre(&estadoMovimiento, "Entrada En Trámite"); err != nil {
		return nil, err
	}

	if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&tipoMovimiento, data.FormatoTipoMovimientoId); err != nil {
		return nil, err
	}

	movimiento := models.Movimiento{
		Observacion:             data.Observacion,
		Detalle:                 data.Detalle,
		Activo:                  true,
		FormatoTipoMovimientoId: &models.FormatoTipoMovimiento{Id: tipoMovimiento},
		EstadoMovimientoId:      &models.EstadoMovimiento{Id: estadoMovimiento},
	}

	if mov, err := movimientosArka.PostMovimiento(&movimiento); err != nil {
		return nil, err
	} else {
		movimiento = *mov
	}

	// Crea registro en table soporte_movimiento si es necesario
	if data.SoporteMovimientoId != 0 {
		soporte := models.SoporteMovimiento{
			DocumentoId:  data.SoporteMovimientoId,
			Activo:       true,
			MovimientoId: &models.Movimiento{Id: movimiento.Id},
		}

		if _, err := movimientosArka.PostSoporteMovimiento(&soporte); err != nil {
			return nil, err
		}

	}

	if elementos, err := asignarPlacaActa(actaRecibidoId); err != nil {
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
