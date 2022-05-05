package entradaHelper

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
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
		acta             models.TransaccionActaRecibido
		tipoMovimiento   int
		estadoMovimiento int
		consecutivo      models.Consecutivo
	)

	resultado := make(map[string]interface{})

	detalleJSON := map[string]interface{}{}
	if err := utilsHelper.Unmarshal(data.Detalle, &detalleJSON); err != nil {
		return nil, err
	}

	ctxConsecutivo, _ := beego.AppConfig.Int("contxtEntradaCons")
	if err := consecutivos.Get(ctxConsecutivo, "Entradas Arka", &consecutivo); err != nil {
		return nil, err
	}

	detalleJSON["consecutivo"] = consecutivos.Format("%05d", getTipoComprobanteEntradas(), &consecutivo)
	detalleJSON["ConsecutivoId"] = consecutivo.Id
	resultado["Consecutivo"] = detalleJSON["consecutivo"]

	if err := utilsHelper.Marshal(detalleJSON, &data.Detalle); err != nil {
		return nil, err
	}

	// Consulta el acta
	actaRecibidoId := int(detalleJSON["acta_recibido_id"].(float64))
	if err := actaRecibido.GetTransaccionActaRecibidoById(actaRecibidoId, &acta); err != nil {
		return nil, err
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

	if err := movimientosArka.PostMovimiento(&movimiento); err != nil {
		return nil, err
	}

	// Crea registro en table soporte_movimiento si es necesario
	if data.SoporteMovimientoId != 0 {
		soporte := models.SoporteMovimiento{
			DocumentoId:  data.SoporteMovimientoId,
			Activo:       true,
			MovimientoId: &models.Movimiento{Id: movimiento.Id},
		}

		if err := movimientosArka.PostSoporteMovimiento(&soporte); err != nil {
			return nil, err
		}

	}

	if elementos, err := asignarPlacaActa(actaRecibidoId); err != nil {
		return nil, outputError
	} else {
		acta.Elementos = elementos
	}

	// Actualiza el estado del acta
	urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/" + strconv.Itoa(int(actaRecibidoId))
	acta.UltimoEstado.EstadoActaId.Id = 6
	acta.UltimoEstado.Id = 0

	if err := request.SendJson(urlcrud, "PUT", &res, &acta); err != nil {
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
