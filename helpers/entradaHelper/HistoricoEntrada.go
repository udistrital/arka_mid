package entradaHelper

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

var getConsecutivoByIDEntradaHistorica = consecutivos.GetById
var getTransaccionActaRecibidoEntradaHistorica = actaRecibido.GetTransaccionActaRecibidoById

// RegistrarEntradaHistorica crea y aprueba una entrada histórica usando un consecutivo y año específicos.
func RegistrarEntradaHistorica(data *models.TransaccionEntradaHistorica, resultado *models.ResultadoMovimiento) (outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("RegistrarEntradaHistorica - Unhandled Error!", "500")

	if data == nil {
		return buildHistoricoPasoError("validar payload", "RegistrarEntradaHistorica - data", "payload nil", "400", nil)
	}
	if err := validarAprobacionHistorica(data); err != nil {
		return buildHistoricoPasoError("validar payload", "RegistrarEntradaHistorica - validarAprobacionHistorica", err.Error(), "400", nil)
	}
	if resultado == nil {
		return buildHistoricoPasoError("inicializar resultado", "RegistrarEntradaHistorica - resultado", "resultado nil", "500", nil)
	}

	resultado.Movimiento = models.Movimiento{
		Observacion:             data.Observacion,
		Activo:                  true,
		FormatoTipoMovimientoId: &models.FormatoTipoMovimiento{},
		EstadoMovimientoId:      &models.EstadoMovimiento{},
		FechaCreacion:           data.FechaCreacion,
		FechaModificacion:       normalizarFechaModificacionHistorica(data),
	}

	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&resultado.Movimiento.EstadoMovimientoId.Id, "Entrada En Trámite")
	if outputError != nil {
		return wrapHistoricoDependencyError(
			"resolver estado inicial",
			"RegistrarEntradaHistorica - GetEstadoMovimientoIdByNombre",
			"no se pudo resolver el estado inicial 'Entrada En Trámite'",
			outputError,
			"500",
		)
	}

	outputError = movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&resultado.Movimiento.FormatoTipoMovimientoId.Id, data.FormatoTipoMovimientoId)
	if outputError != nil {
		return buildHistoricoPasoError(
			"resolver tipo de movimiento",
			"RegistrarEntradaHistorica - GetFormatoTipoMovimientoIdByCodigoAbreviacion",
			fmt.Sprintf("FormatoTipoMovimientoId %q inválido o no parametrizado", data.FormatoTipoMovimientoId),
			"400",
			outputError,
		)
	}

	outputError = crearDetalleEntrada(data.Detalle, &resultado.Movimiento.Detalle)
	if outputError != nil {
		return wrapHistoricoDependencyError(
			"construir detalle del movimiento",
			"RegistrarEntradaHistorica - crearDetalleEntrada",
			"no se pudo construir el detalle de la entrada histórica",
			outputError,
			"400",
		)
	}

	var acta models.TransaccionActaRecibido
	if data.Detalle.ActaRecibidoId > 0 {
		outputError = getTransaccionActaRecibidoEntradaHistorica(data.Detalle.ActaRecibidoId, false, &acta)
		if outputError != nil {
			return wrapHistoricoDependencyError(
				"consultar acta asociada",
				fmt.Sprintf("RegistrarEntradaHistorica - GetTransaccionActaRecibidoById(acta_recibido_id=%d)", data.Detalle.ActaRecibidoId),
				fmt.Sprintf("no se pudo consultar el acta %d", data.Detalle.ActaRecibidoId),
				outputError,
				"404",
			)
		}
		if acta.UltimoEstado == nil || acta.UltimoEstado.EstadoActaId == nil || acta.UltimoEstado.EstadoActaId.CodigoAbreviacion != "Aceptada" {
			return buildHistoricoPasoError(
				"validar estado del acta",
				fmt.Sprintf("RegistrarEntradaHistorica - validarEstadoActa(acta_recibido_id=%d)", data.Detalle.ActaRecibidoId),
				mensajeEstadoActaHistoricaInvalido(data.Detalle.ActaRecibidoId, acta.UltimoEstado),
				"400",
				nil,
			)
		}
	}

	if outputError = aplicarConsecutivoHistoricoEntrada(&resultado.Movimiento, data.ConsecutivoId, data.Year); outputError != nil {
		return
	}

	if data.Detalle.ActaRecibidoId > 0 {
		resultado.Error, outputError = asignarPlacas(data.Detalle.ActaRecibidoId, &acta.Elementos)
		if outputError != nil {
			return wrapHistoricoDependencyError(
				"asignar placas a elementos del acta",
				fmt.Sprintf("RegistrarEntradaHistorica - asignarPlacas(acta_recibido_id=%d)", data.Detalle.ActaRecibidoId),
				fmt.Sprintf("no se pudieron asignar placas para el acta %d", data.Detalle.ActaRecibidoId),
				outputError,
				"400",
			)
		}
		if resultado.Error != "" {
			return buildHistoricoPasoError(
				"asignar placas a elementos del acta",
				fmt.Sprintf("RegistrarEntradaHistorica - asignarPlacas(acta_recibido_id=%d)", data.Detalle.ActaRecibidoId),
				resultado.Error,
				"400",
				nil,
			)
		}
		if len(acta.Elementos) == 0 {
			return buildHistoricoPasoError(
				"validar elementos del acta",
				fmt.Sprintf("RegistrarEntradaHistorica - validarElementosActa(acta_recibido_id=%d)", data.Detalle.ActaRecibidoId),
				fmt.Sprintf("El acta %d no tiene elementos asociados y no puede usarse para una entrada histórica.", data.Detalle.ActaRecibidoId),
				"400",
				nil,
			)
		}
	}

	outputError = movimientosArka.PostMovimiento(&resultado.Movimiento)
	if outputError != nil {
		return wrapHistoricoDependencyError(
			"registrar movimiento en arka",
			"RegistrarEntradaHistorica - PostMovimiento",
			"no se pudo registrar el movimiento histórico en movimientos_arka_crud",
			outputError,
			"502",
		)
	}

	resultado.Movimiento.FechaCreacion = data.FechaCreacion
	resultado.Movimiento.FechaModificacion = normalizarFechaModificacionHistorica(data)
	outputError = movimientosArka.PutMovimiento(&resultado.Movimiento, resultado.Movimiento.Id)
	if outputError != nil {
		return wrapHistoricoDependencyError(
			"actualizar fechas históricas del movimiento",
			fmt.Sprintf("RegistrarEntradaHistorica - PutMovimiento(movimiento_id=%d)", resultado.Movimiento.Id),
			fmt.Sprintf("no se pudieron persistir las fechas históricas del movimiento %d", resultado.Movimiento.Id),
			outputError,
			"502",
		)
	}

	if data.SoporteMovimientoId > 0 {
		soporte := models.SoporteMovimiento{
			DocumentoId:  data.SoporteMovimientoId,
			Activo:       true,
			MovimientoId: &models.Movimiento{Id: resultado.Movimiento.Id},
		}

		outputError = movimientosArka.PostSoporteMovimiento(&soporte)
		if outputError != nil {
			return wrapHistoricoDependencyError(
				"asociar soporte al movimiento",
				fmt.Sprintf("RegistrarEntradaHistorica - PostSoporteMovimiento(movimiento_id=%d)", resultado.Movimiento.Id),
				fmt.Sprintf("no se pudo asociar el soporte %d al movimiento %d", data.SoporteMovimientoId, resultado.Movimiento.Id),
				outputError,
				"502",
			)
		}
	}

	if data.Detalle.ActaRecibidoId > 0 {
		acta.UltimoEstado.EstadoActaId.Id = 6
		acta.UltimoEstado.Id = 0
		outputError = actaRecibido.PutTransaccionActaRecibido(data.Detalle.ActaRecibidoId, &acta)
		if outputError != nil {
			return wrapHistoricoDependencyError(
				"actualizar estado del acta después del registro",
				fmt.Sprintf("RegistrarEntradaHistorica - PutTransaccionActaRecibido(acta_recibido_id=%d)", data.Detalle.ActaRecibidoId),
				fmt.Sprintf("no se pudo actualizar el estado del acta %d después de registrar la entrada histórica", data.Detalle.ActaRecibidoId),
				outputError,
				"502",
			)
		}
	}

	outputError = aprobarEntradaHistorica(resultado.Movimiento.Id, data, resultado)
	return
}

func aprobarEntradaHistorica(entradaId int, data *models.TransaccionEntradaHistorica, resultado *models.ResultadoMovimiento) (outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("aprobarEntradaHistorica - Unhandled Error!", "500")

	formato, outputError := getFormato(entradaId, resultado)
	if outputError != nil {
		return wrapHistoricoDependencyError(
			"consultar movimiento recién registrado",
			fmt.Sprintf("aprobarEntradaHistorica - getFormato(entrada_id=%d)", entradaId),
			fmt.Sprintf("no se pudo consultar el movimiento %d para aprobar la entrada histórica", entradaId),
			outputError,
			"404",
		)
	}
	if resultado.Error != "" {
		return buildHistoricoPasoError(
			"consultar movimiento recién registrado",
			fmt.Sprintf("aprobarEntradaHistorica - getFormato(entrada_id=%d)", entradaId),
			resultado.Error,
			"400",
			nil,
		)
	}

	resultado.Movimiento.FechaCreacion = data.FechaCreacion
	resultado.Movimiento.FechaModificacion = normalizarFechaModificacionHistorica(data)
	resultado.Movimiento.FechaCorte = &data.FechaCorte

	terceroId, outputError := getTerceroEntrada(formato, resultado)
	if outputError != nil {
		return wrapHistoricoDependencyError(
			"resolver tercero contable de la entrada",
			fmt.Sprintf("aprobarEntradaHistorica - getTerceroEntrada(entrada_id=%d)", entradaId),
			fmt.Sprintf("no se pudo resolver el tercero contable para la entrada histórica %d", entradaId),
			outputError,
			"502",
		)
	}
	if resultado.Error != "" {
		return buildHistoricoPasoError(
			"resolver tercero contable de la entrada",
			fmt.Sprintf("aprobarEntradaHistorica - getTerceroEntrada(entrada_id=%d)", entradaId),
			resultado.Error,
			"400",
			nil,
		)
	}

	elementos, novedades, outputError := getElementosEntrada(formato, entradaId, resultado)
	if outputError != nil {
		return wrapHistoricoDependencyError(
			"resolver elementos y novedades de la entrada",
			fmt.Sprintf("aprobarEntradaHistorica - getElementosEntrada(entrada_id=%d)", entradaId),
			fmt.Sprintf("no se pudieron resolver los elementos o novedades de la entrada histórica %d", entradaId),
			outputError,
			"502",
		)
	}
	if resultado.Error != "" {
		return buildHistoricoPasoError(
			"resolver elementos y novedades de la entrada",
			fmt.Sprintf("aprobarEntradaHistorica - getElementosEntrada(entrada_id=%d)", entradaId),
			resultado.Error,
			"400",
			nil,
		)
	}

	outputError = contabilidadEntrada(resultado, formato, elementos, terceroId)
	if outputError != nil {
		return wrapHistoricoDependencyError(
			"generar transacción contable de la entrada",
			fmt.Sprintf("aprobarEntradaHistorica - contabilidadEntrada(entrada_id=%d)", entradaId),
			fmt.Sprintf("no se pudo generar la transacción contable de la entrada histórica %d", entradaId),
			outputError,
			"502",
		)
	}
	if resultado.Error != "" {
		return buildHistoricoPasoError(
			"generar transacción contable de la entrada",
			fmt.Sprintf("aprobarEntradaHistorica - contabilidadEntrada(entrada_id=%d)", entradaId),
			resultado.Error,
			"400",
			nil,
		)
	}

	for _, nov := range novedades {
		outputError = movimientosArka.PostNovedadElemento(&nov)
		if outputError != nil {
			return wrapHistoricoDependencyError(
				"registrar novedades de los elementos",
				fmt.Sprintf("aprobarEntradaHistorica - PostNovedadElemento(entrada_id=%d)", entradaId),
				fmt.Sprintf("no se pudieron registrar las novedades de la entrada histórica %d", entradaId),
				outputError,
				"502",
			)
		}
	}

	outputError = movimientosArka.PutMovimiento(&resultado.Movimiento, resultado.Movimiento.Id)
	if outputError != nil {
		return wrapHistoricoDependencyError(
			"persistir aprobación histórica del movimiento",
			fmt.Sprintf("aprobarEntradaHistorica - PutMovimiento(movimiento_id=%d)", resultado.Movimiento.Id),
			fmt.Sprintf("no se pudo persistir la aprobación histórica del movimiento %d", resultado.Movimiento.Id),
			outputError,
			"502",
		)
	}

	return nil
}

func aplicarConsecutivoHistoricoEntrada(entrada *models.Movimiento, consecutivoID, year int) (outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("aplicarConsecutivoHistoricoEntrada - Unhandled Error!", "500")

	if entrada == nil {
		return buildHistoricoPasoError("aplicar consecutivo histórico", "aplicarConsecutivoHistoricoEntrada - entrada", "entrada nil", "500", nil)
	}

	var consecutivo models.Consecutivo
	outputError = getConsecutivoByIDEntradaHistorica(consecutivoID, &consecutivo)
	if outputError != nil {
		return wrapHistoricoDependencyError(
			"consultar consecutivo histórico",
			fmt.Sprintf("aplicarConsecutivoHistoricoEntrada - GetById(consecutivo_id=%d)", consecutivoID),
			fmt.Sprintf("no se pudo consultar el consecutivo %d", consecutivoID),
			outputError,
			"404",
		)
	}
	if consecutivo.Id <= 0 {
		return buildHistoricoPasoError(
			"validar consecutivo histórico",
			"aplicarConsecutivoHistoricoEntrada - getConsecutivoByIDEntradaHistorica",
			fmt.Sprintf("ConsecutivoId %d no existe o no está disponible para construir el consecutivo histórico", consecutivoID),
			"404",
			nil,
		)
	}

	consecutivo.Year = year
	entrada.ConsecutivoId = &consecutivo.Id
	entrada.Consecutivo = stringHistoricoConsecutivoEntrada(&consecutivo)
	return nil
}

func stringHistoricoConsecutivoEntrada(consecutivo *models.Consecutivo) *string {
	if consecutivo == nil {
		return nil
	}
	formatted := consecutivos.Format("%05d", getTipoComprobanteEntradas(), consecutivo)
	return &formatted
}

func normalizarFechaModificacionHistorica(data *models.TransaccionEntradaHistorica) time.Time {
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

func validarAprobacionHistorica(data *models.TransaccionEntradaHistorica) error {
	if data == nil {
		return fmt.Errorf("payload nil")
	}
	var errores []string

	if strings.TrimSpace(data.FormatoTipoMovimientoId) == "" {
		errores = append(errores, "FormatoTipoMovimientoId es obligatorio")
	}
	if data.ConsecutivoId <= 0 {
		errores = append(errores, "ConsecutivoId es obligatorio y debe ser mayor a 0")
	}
	if data.Year <= 0 {
		errores = append(errores, "Year es obligatorio y debe ser mayor a 0")
	}
	if data.FechaCreacion.IsZero() {
		errores = append(errores, "FechaCreacion es obligatoria y debe venir en formato RFC3339")
	}
	if data.FechaCorte.IsZero() {
		errores = append(errores, "FechaCorte es obligatoria y debe venir en formato RFC3339")
	}
	if data.Detalle.ActaRecibidoId <= 0 && len(data.Detalle.Elementos) == 0 {
		errores = append(errores, "Detalle.acta_recibido_id o Detalle.elementos es obligatorio")
	}

	if len(errores) > 0 {
		return fmt.Errorf("payload histórico inválido: %s", strings.Join(errores, "; "))
	}

	return nil
}

func debugHistoricoEntrada(data *models.TransaccionEntradaHistorica) string {
	if data == nil {
		return "<nil>"
	}
	return "consecutivo_id=" + strconv.Itoa(data.ConsecutivoId) +
		" year=" + strconv.Itoa(data.Year) +
		" fecha_creacion=" + data.FechaCreacion.Format(time.RFC3339) +
		" fecha_corte=" + data.FechaCorte.Format(time.RFC3339)
}

func buildHistoricoPasoError(paso, funcion, detalle, status string, causa interface{}) map[string]interface{} {
	errPayload := map[string]interface{}{
		"paso":    paso,
		"detalle": detalle,
	}
	if causa != nil {
		errPayload["causa"] = causa
	}

	return errorCtrl.Error(funcion, errPayload, status)
}

func wrapHistoricoDependencyError(paso, funcion, mensaje string, outputError map[string]interface{}, fallbackStatus string) map[string]interface{} {
	if outputError == nil {
		return nil
	}

	var causa interface{}
	if errValue, ok := outputError["err"]; ok && errValue != nil {
		causa = errValue
	}

	return buildHistoricoPasoError(paso, funcion, mensaje, outputErrorStatusHistorico(outputError, fallbackStatus), causa)
}

func outputErrorStatusHistorico(err map[string]interface{}, fallback string) string {
	if err == nil {
		return fallback
	}

	if status, ok := err["status"].(string); ok && status != "" {
		return status
	}

	return fallback
}

func mensajeEstadoActaHistoricaInvalido(actaID int, historico *models.HistoricoActa) string {
	if historico == nil || historico.EstadoActaId == nil {
		return fmt.Sprintf("El acta %d no tiene estado actual identificable. Para registrar una entrada histórica debe estar en estado 'Aceptada'.", actaID)
	}

	estado := historico.EstadoActaId
	if strings.TrimSpace(estado.Nombre) != "" {
		return fmt.Sprintf(
			"El acta %d está en estado %q (%s). Para registrar una entrada histórica debe estar en estado 'Aceptada'.",
			actaID,
			estado.Nombre,
			estado.CodigoAbreviacion,
		)
	}

	return fmt.Sprintf(
		"El acta %d está en estado %q. Para registrar una entrada histórica debe estar en estado 'Aceptada'.",
		actaID,
		estado.CodigoAbreviacion,
	)
}
