package salidaHelper

import (
	"fmt"
	"strconv"
	"time"

	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

var getConsecutivoByIDSalidaHistorica = consecutivos.GetById

// RegistrarSalidaHistorica crea y aprueba una salida histórica usando un consecutivo y año específicos.
func RegistrarSalidaHistorica(data *models.TransaccionSalidaHistorica, resultado *models.ResultadoMovimiento) (outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("RegistrarSalidaHistorica - Unhandled Error!", "500")

	if data == nil {
		return errorCtrl.Error("RegistrarSalidaHistorica - data", "payload nil", "400")
	}
	if err := validarAprobacionHistoricaSalida(data); err != nil {
		return errorCtrl.Error("RegistrarSalidaHistorica - validarAprobacionHistoricaSalida", err, "400")
	}
	if resultado == nil {
		return errorCtrl.Error("RegistrarSalidaHistorica - resultado", "resultado nil", "500")
	}

	if len(data.Salidas) != 1 || data.Salidas[0].Salida == nil {
		return errorCtrl.Error("RegistrarSalidaHistorica - salidas", "debe especificar exactamente una salida histórica", "400")
	}

	var estadoMovimientoId int
	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&estadoMovimientoId, "Salida En Trámite")
	if outputError != nil {
		return
	}

	payload := &models.SalidaGeneral{
		Salidas: make([]models.TrSalida, len(data.Salidas)),
	}
	copy(payload.Salidas, data.Salidas)

	if outputError = normalizarNuevaSalida(&payload.Salidas[0], estadoMovimientoId); outputError != nil {
		return
	}
	payload.Salidas[0].Salida.FechaCreacion = data.FechaCreacion
	payload.Salidas[0].Salida.FechaModificacion = normalizarFechaModificacionHistoricaSalida(data)

	if outputError = aplicarConsecutivoHistoricoSalida(payload.Salidas[0].Salida, data.ConsecutivoId, data.Year); outputError != nil {
		return
	}

	if outputError = movimientosArka.PostTrSalida(payload); outputError != nil {
		return
	}
	if outputError = completarSalidasPersistidas(payload); outputError != nil {
		return
	}
	if len(payload.Salidas) == 0 || payload.Salidas[0].Salida == nil || payload.Salidas[0].Salida.Id <= 0 {
		return errorCtrl.Error("RegistrarSalidaHistorica - completarSalidasPersistidas", "no se confirmó la salida histórica creada", "502")
	}

	trSalida, outputError := movimientosArka.GetTrSalida(payload.Salidas[0].Salida.Id)
	if outputError != nil {
		return
	}
	if trSalida == nil || trSalida.Salida == nil {
		return errorCtrl.Error("RegistrarSalidaHistorica - movimientosArka.GetTrSalida", "tr_salida nil", "404")
	}

	trSalida.Salida.FechaCreacion = data.FechaCreacion
	trSalida.Salida.FechaModificacion = normalizarFechaModificacionHistoricaSalida(data)
	trSalida.Salida.FechaCorte = &data.FechaCorte
	outputError = movimientosArka.PutMovimiento(trSalida.Salida, trSalida.Salida.Id)
	if outputError != nil {
		return
	}

	for _, elemento := range trSalida.Elementos {
		if elemento == nil || elemento.Id <= 0 {
			continue
		}
		elemento.FechaCreacion = data.FechaCreacion
		elemento.FechaModificacion = normalizarFechaModificacionHistoricaSalida(data)
		if _, outputError = movimientosArka.PutElementosMovimiento(elemento, elemento.Id); outputError != nil {
			return
		}
	}

	return AprobarSalida(payload.Salidas[0].Salida.Id, resultado)
}

func aplicarConsecutivoHistoricoSalida(salida *models.Movimiento, consecutivoID, year int) (outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("aplicarConsecutivoHistoricoSalida - Unhandled Error!", "500")

	if salida == nil {
		return errorCtrl.Error("aplicarConsecutivoHistoricoSalida - salida", "salida nil", "500")
	}

	var consecutivo models.Consecutivo
	outputError = getConsecutivoByIDSalidaHistorica(consecutivoID, &consecutivo)
	if outputError != nil {
		return
	}
	if consecutivo.Id <= 0 {
		return errorCtrl.Error("aplicarConsecutivoHistoricoSalida - getConsecutivoByIDSalidaHistorica", "no se encontró el consecutivo indicado", "404")
	}

	consecutivo.Year = year
	salida.ConsecutivoId = &consecutivo.Id
	salida.Consecutivo = stringHistoricoConsecutivoSalida(&consecutivo)
	return nil
}

func stringHistoricoConsecutivoSalida(consecutivo *models.Consecutivo) *string {
	if consecutivo == nil {
		return nil
	}
	formatted := consecutivos.Format("%05d", getTipoComprobanteSalidas(), consecutivo)
	return &formatted
}

func normalizarFechaModificacionHistoricaSalida(data *models.TransaccionSalidaHistorica) time.Time {
	if data == nil {
		return time.Time{}
	}
	if !data.FechaModificacion.IsZero() {
		return data.FechaModificacion
	}
	if !data.FechaCorte.IsZero() {
		return data.FechaCorte
	}
	return data.FechaCreacion
}

func validarAprobacionHistoricaSalida(data *models.TransaccionSalidaHistorica) error {
	if data == nil {
		return fmt.Errorf("payload nil")
	}
	if data.ConsecutivoId <= 0 {
		return fmt.Errorf("consecutivo_id inválido")
	}
	if data.Year <= 0 {
		return fmt.Errorf("year inválido")
	}
	if data.FechaCreacion.IsZero() {
		return fmt.Errorf("fecha_creacion vacía")
	}
	if data.FechaCorte.IsZero() {
		return fmt.Errorf("fecha_corte vacía")
	}
	if len(data.Salidas) != 1 || data.Salidas[0].Salida == nil {
		return fmt.Errorf("debe especificar exactamente una salida histórica")
	}
	return nil
}

func debugHistoricoSalida(data *models.TransaccionSalidaHistorica) string {
	if data == nil {
		return "<nil>"
	}
	return "consecutivo_id=" + strconv.Itoa(data.ConsecutivoId) +
		" year=" + strconv.Itoa(data.Year) +
		" fecha_creacion=" + data.FechaCreacion.Format(time.RFC3339) +
		" fecha_corte=" + data.FechaCorte.Format(time.RFC3339)
}
