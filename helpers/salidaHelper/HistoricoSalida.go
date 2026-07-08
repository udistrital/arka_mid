package salidaHelper

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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
		return buildHistoricoSalidaPasoError("validar payload", "RegistrarSalidaHistorica - data", "payload nil", "400", nil)
	}
	if err := validarAprobacionHistoricaSalida(data); err != nil {
		return buildHistoricoSalidaPasoError("validar payload", "RegistrarSalidaHistorica - validarAprobacionHistoricaSalida", err.Error(), "400", nil)
	}
	if resultado == nil {
		return buildHistoricoSalidaPasoError("inicializar resultado", "RegistrarSalidaHistorica - resultado", "resultado nil", "500", nil)
	}

	if len(data.Salidas) != 1 || data.Salidas[0].Salida == nil {
		return buildHistoricoSalidaPasoError("validar payload", "RegistrarSalidaHistorica - salidas", "debe especificar exactamente una salida histórica", "400", nil)
	}

	var estadoMovimientoId int
	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&estadoMovimientoId, "Salida En Trámite")
	if outputError != nil {
		return wrapHistoricoSalidaDependencyError(
			"resolver estado inicial",
			"RegistrarSalidaHistorica - GetEstadoMovimientoIdByNombre",
			"no se pudo resolver el estado inicial 'Salida En Trámite'",
			outputError,
			"500",
		)
	}

	payload := &models.SalidaGeneral{
		Salidas: make([]models.TrSalida, len(data.Salidas)),
	}
	copy(payload.Salidas, data.Salidas)

	if outputError = normalizarNuevaSalida(&payload.Salidas[0], estadoMovimientoId); outputError != nil {
		return wrapHistoricoSalidaDependencyError(
			"normalizar salida histórica",
			"RegistrarSalidaHistorica - normalizarNuevaSalida",
			"no se pudo normalizar la salida histórica",
			outputError,
			"400",
		)
	}
	payload.Salidas[0].Salida.FechaCreacion = data.FechaCreacion
	payload.Salidas[0].Salida.FechaModificacion = normalizarFechaModificacionHistoricaSalida(data)

	if outputError = aplicarConsecutivoHistoricoSalida(payload.Salidas[0].Salida, data.ConsecutivoId, data.Year); outputError != nil {
		return wrapHistoricoSalidaDependencyError(
			"aplicar consecutivo histórico",
			"RegistrarSalidaHistorica - aplicarConsecutivoHistoricoSalida",
			"no se pudo aplicar el consecutivo histórico de la salida",
			outputError,
			"404",
		)
	}

	if outputError = movimientosArka.PostTrSalida(payload); outputError != nil {
		return wrapPostTrSalidaHistoricaError(payload, outputError)
	}
	if outputError = completarSalidasPersistidas(payload); outputError != nil {
		return wrapHistoricoSalidaDependencyError(
			"confirmar persistencia de la salida",
			"RegistrarSalidaHistorica - completarSalidasPersistidas",
			"no se confirmó la salida histórica creada",
			outputError,
			"502",
		)
	}
	if len(payload.Salidas) == 0 || payload.Salidas[0].Salida == nil || payload.Salidas[0].Salida.Id <= 0 {
		return buildHistoricoSalidaPasoError("confirmar persistencia de la salida", "RegistrarSalidaHistorica - completarSalidasPersistidas", "no se confirmó la salida histórica creada", "502", nil)
	}

	trSalida, outputError := movimientosArka.GetTrSalida(payload.Salidas[0].Salida.Id)
	if outputError != nil {
		return wrapHistoricoSalidaDependencyError(
			"consultar salida recién creada",
			fmt.Sprintf("RegistrarSalidaHistorica - GetTrSalida(salida_id=%d)", payload.Salidas[0].Salida.Id),
			fmt.Sprintf("no se pudo consultar la salida histórica %d recién creada", payload.Salidas[0].Salida.Id),
			outputError,
			"404",
		)
	}
	if trSalida == nil || trSalida.Salida == nil {
		return buildHistoricoSalidaPasoError(
			"consultar salida recién creada",
			fmt.Sprintf("RegistrarSalidaHistorica - GetTrSalida(salida_id=%d)", payload.Salidas[0].Salida.Id),
			"tr_salida nil",
			"404",
			nil,
		)
	}

	trSalida.Salida.FechaCreacion = data.FechaCreacion
	trSalida.Salida.FechaModificacion = normalizarFechaModificacionHistoricaSalida(data)
	trSalida.Salida.FechaCorte = &data.FechaCorte
	outputError = movimientosArka.PutMovimiento(trSalida.Salida, trSalida.Salida.Id)
	if outputError != nil {
		payloadMovimiento, marshalErr := serializarMovimientoHistoricoSalida(*trSalida.Salida)
		if marshalErr != nil {
			payloadMovimiento = "no se pudo serializar el payload de movimiento histórico: " + marshalErr.Error()
		}
		return wrapHistoricoSalidaDependencyError(
			"actualizar fechas históricas del movimiento",
			fmt.Sprintf("RegistrarSalidaHistorica - PutMovimiento(salida_id=%d)", trSalida.Salida.Id),
			fmt.Sprintf("no se pudieron persistir las fechas históricas de la salida %d", trSalida.Salida.Id),
			mergeHistoricoSalidaErrorCause(outputError, map[string]interface{}{
				"payload_movimiento": payloadMovimiento,
				"salida_id":          trSalida.Salida.Id,
			}),
			"502",
		)
	}

	for _, elemento := range trSalida.Elementos {
		if elemento == nil || elemento.Id <= 0 {
			continue
		}
		elementoHistorico := normalizarElementoMovimientoHistoricoSalida(elemento, data)
		if elementoHistorico != nil && (elementoHistorico.MovimientoId == nil || elementoHistorico.MovimientoId.Id <= 0) {
			elementoHistorico.MovimientoId = &models.Movimiento{Id: trSalida.Salida.Id}
		}
		if _, outputError = movimientosArka.PutElementosMovimiento(elementoHistorico, elementoHistorico.Id); outputError != nil {
			payloadElemento, marshalErr := serializarElementoMovimientoHistoricoSalida(*elementoHistorico)
			if marshalErr != nil {
				payloadElemento = "no se pudo serializar el payload de elemento_movimiento: " + marshalErr.Error()
			}
			return wrapHistoricoSalidaDependencyError(
				"actualizar fechas históricas de los elementos de la salida",
				fmt.Sprintf("RegistrarSalidaHistorica - PutElementosMovimiento(elemento_movimiento_id=%d)", elementoHistorico.Id),
				fmt.Sprintf("no se pudieron persistir las fechas históricas del elemento_movimiento %d", elementoHistorico.Id),
				mergeHistoricoSalidaErrorCause(outputError, map[string]interface{}{
					"salida_id":               trSalida.Salida.Id,
					"elemento_movimiento_id":  elementoHistorico.Id,
					"elemento_acta_id":        valorElementoActaHistoricoSalida(elementoHistorico.ElementoActaId),
					"payload_elemento_movimiento": payloadElemento,
				}),
				"502",
			)
		}
	}

	outputError = AprobarSalida(payload.Salidas[0].Salida.Id, resultado)
	if outputError != nil {
		return wrapHistoricoSalidaDependencyError(
			"aprobar salida histórica",
			fmt.Sprintf("RegistrarSalidaHistorica - AprobarSalida(salida_id=%d)", payload.Salidas[0].Salida.Id),
			fmt.Sprintf("no se pudo aprobar la salida histórica %d", payload.Salidas[0].Salida.Id),
			outputError,
			"502",
		)
	}

	return nil
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

func normalizarElementoMovimientoHistoricoSalida(
	elemento *models.ElementosMovimiento,
	data *models.TransaccionSalidaHistorica,
) *models.ElementosMovimiento {
	if elemento == nil {
		return nil
	}

	elementoHistorico := &models.ElementosMovimiento{
		Id:                elemento.Id,
		ElementoCatalogoId: elemento.ElementoCatalogoId,
		Unidad:            elemento.Unidad,
		ValorUnitario:     elemento.ValorUnitario,
		ValorTotal:        elemento.ValorTotal,
		SaldoCantidad:     elemento.SaldoCantidad,
		SaldoValor:        elemento.SaldoValor,
		VidaUtil:          elemento.VidaUtil,
		ValorResidual:     elemento.ValorResidual,
		Activo:            elemento.Activo,
	}

	if elemento.ElementoActaId != nil {
		elementoActaID := *elemento.ElementoActaId
		elementoHistorico.ElementoActaId = &elementoActaID
	}

	if elemento.MovimientoId != nil && elemento.MovimientoId.Id > 0 {
		elementoHistorico.MovimientoId = &models.Movimiento{Id: elemento.MovimientoId.Id}
	}

	if data != nil {
		elementoHistorico.FechaCreacion = data.FechaCreacion
		elementoHistorico.FechaModificacion = normalizarFechaModificacionHistoricaSalida(data)
	}

	return elementoHistorico
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

func wrapPostTrSalidaHistoricaError(payload *models.SalidaGeneral, outputError map[string]interface{}) map[string]interface{} {
	if outputError == nil {
		return nil
	}

	errMsg, _ := outputError["err"].(string)
	if !strings.Contains(errMsg, "uq_elemento_acta_id") {
		payloadSalida, marshalErr := serializarSalidaGeneralHistorica(payload)
		if marshalErr != nil {
			payloadSalida = "no se pudo serializar el payload enviado a tr_salida: " + marshalErr.Error()
		}
		return wrapHistoricoSalidaDependencyError(
			"registrar salida en arka",
			"RegistrarSalidaHistorica - PostTrSalida",
			"no se pudo registrar la salida histórica en movimientos_arka_crud",
			mergeHistoricoSalidaErrorCause(outputError, map[string]interface{}{
				"payload_tr_salida": payloadSalida,
			}),
			"502",
		)
	}

	var elementoActaIDs []string
	parentID := 0
	if payload != nil && len(payload.Salidas) > 0 {
		if payload.Salidas[0].Salida != nil && payload.Salidas[0].Salida.MovimientoPadreId != nil {
			parentID = payload.Salidas[0].Salida.MovimientoPadreId.Id
		}
		for _, elemento := range payload.Salidas[0].Elementos {
			if elemento == nil || elemento.ElementoActaId == nil || *elemento.ElementoActaId <= 0 {
				continue
			}
			elementoActaIDs = append(elementoActaIDs, strconv.Itoa(*elemento.ElementoActaId))
		}
	}

	detalle := "uno o más elementos ya están asociados a otra salida en movimientos_arka_crud"
	if len(elementoActaIDs) > 0 {
		detalle += ". elemento_acta_id=" + strings.Join(elementoActaIDs, ",")
	}
	if parentID > 0 {
		detalle += ", movimiento_padre_id=" + strconv.Itoa(parentID)
	}
	detalle += ". Es probable que exista un registro parcial previo y deba limpiarse en movimientos_arka_crud antes de reintentar."

	return buildHistoricoSalidaPasoError("registrar salida en arka", "RegistrarSalidaHistorica - PostTrSalida", detalle, "409", nil)
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

func buildHistoricoSalidaPasoError(paso, funcion, detalle, status string, causa interface{}) map[string]interface{} {
	errPayload := map[string]interface{}{
		"paso":    paso,
		"detalle": detalle,
	}
	if causa != nil {
		errPayload["causa"] = causa
	}

	return errorCtrl.Error(funcion, errPayload, status)
}

func wrapHistoricoSalidaDependencyError(paso, funcion, mensaje string, outputError map[string]interface{}, fallbackStatus string) map[string]interface{} {
	if outputError == nil {
		return nil
	}

	var causa interface{}
	if errValue, ok := outputError["err"]; ok && errValue != nil {
		causa = errValue
	}

	return buildHistoricoSalidaPasoError(paso, funcion, mensaje, outputErrorStatusHistoricoSalida(outputError, fallbackStatus), causa)
}

func mergeHistoricoSalidaErrorCause(outputError map[string]interface{}, extra map[string]interface{}) map[string]interface{} {
	if outputError == nil || len(extra) == 0 {
		return outputError
	}

	merged := make(map[string]interface{}, len(outputError))
	for key, value := range outputError {
		merged[key] = value
	}

	if errValue, ok := outputError["err"].(map[string]interface{}); ok {
		errMerged := make(map[string]interface{}, len(errValue)+len(extra))
		for key, value := range errValue {
			errMerged[key] = value
		}
		for key, value := range extra {
			errMerged[key] = value
		}
		merged["err"] = errMerged
		return merged
	}

	if errValue, ok := outputError["err"]; ok && errValue != nil {
		merged["err"] = map[string]interface{}{
			"causa_original": errValue,
		}
		if errMerged, ok := merged["err"].(map[string]interface{}); ok {
			for key, value := range extra {
				errMerged[key] = value
			}
		}
		return merged
	}

	merged["err"] = extra
	return merged
}

func outputErrorStatusHistoricoSalida(err map[string]interface{}, fallback string) string {
	if err == nil {
		return fallback
	}

	if status, ok := err["status"].(string); ok && status != "" {
		return status
	}

	return fallback
}

func serializarSalidaGeneralHistorica(payload *models.SalidaGeneral) (string, error) {
	rawBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(rawBytes), nil
}

func serializarMovimientoHistoricoSalida(movimiento models.Movimiento) (string, error) {
	rawBytes, err := json.Marshal(movimiento)
	if err != nil {
		return "", err
	}
	return string(rawBytes), nil
}

func serializarElementoMovimientoHistoricoSalida(elemento models.ElementosMovimiento) (string, error) {
	rawBytes, err := json.Marshal(elemento)
	if err != nil {
		return "", err
	}
	return string(rawBytes), nil
}

func valorElementoActaHistoricoSalida(elementoActaID *int) int {
	if elementoActaID == nil {
		return 0
	}
	return *elementoActaID
}
