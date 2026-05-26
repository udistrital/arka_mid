package salidaHelper

import (
	"fmt"
	"strings"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

const (
	estadoSalidaAprobada    = "Salida Aprobada"
	estadoSalidaAnulada     = "Salida anulada"
	estadoSalidaRechazada   = "Salida Rechazada"
	estadoEntradaAprobada   = "Entrada Aprobada"
	estadoEntradaConSalida  = "Entrada Con Salida"
	formatoAjusteAutomatico = "AAT"
	estadoAjusteAprobado    = "Ajuste Aprobado"
)

// AnularSalida anula una salida aprobada, genera la reversa contable y restablece la entrada padre cuando corresponde.
func AnularSalida(salidaID int, request *models.AnulacionSalidaRequest, resultado *models.ResultadoAnulacionSalida) (outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("AnularSalida - Unhandled Error!", "500")

	if resultado == nil {
		return map[string]interface{}{
			"funcion": "AnularSalida - resultado",
			"err":     "resultado nil",
			"status":  "500",
		}
	}

	if request == nil {
		request = &models.AnulacionSalidaRequest{}
	}

	trSalida, outputError := movimientosArka.GetTrSalida(salidaID)
	if outputError != nil {
		return outputError
	}
	if trSalida == nil || trSalida.Salida == nil {
		return map[string]interface{}{
			"funcion": "AnularSalida - movimientosArka.GetTrSalida",
			"err":     "tr_salida nil",
			"status":  "404",
		}
	}

	if !validarSalidaAnulable(trSalida, resultado) {
		return nil
	}

	transaccionOriginal, outputError := movimientosContables.GetTransaccion(*trSalida.Salida.ConsecutivoId, "consecutivo", true)
	if outputError != nil {
		return outputError
	}
	if transaccionOriginal == nil || len(transaccionOriginal.Movimientos) == 0 {
		resultado.Error = "No se encontró la transacción contable original de la salida."
		return nil
	}

	for _, elemento := range trSalida.Elementos {
		historial, err := movimientosArka.GetHistorialElemento(elemento.Id, true)
		if err != nil {
			return err
		}

		if ok, msg := historialPermiteAnularSalida(historial, salidaID); !ok {
			resultado.Error = msg
			return nil
		}
	}

	salidaOriginal := cloneMovimientoSalida(trSalida.Salida)
	var entradaOriginal *models.Movimiento
	var entradaActualizada *models.Movimiento
	if trSalida.Salida.MovimientoPadreId != nil && trSalida.Salida.MovimientoPadreId.Id > 0 {
		if entradaPadre, err := movimientosArka.GetMovimientoById(trSalida.Salida.MovimientoPadreId.Id); err != nil {
			return err
		} else if entradaPadre != nil {
			entradaOriginal = cloneMovimientoSalida(entradaPadre)
		}
	}

	movimientoReversion, transaccionReversion, outputError := construirReversionSalida(trSalida, transaccionOriginal, request.Observacion)
	if outputError != nil {
		return outputError
	}
	if movimientoReversion == nil || transaccionReversion == nil {
		resultado.Error = "No se pudo construir la reversa contable de la salida."
		return nil
	}

	resultado.TransaccionContable.Concepto = transaccionReversion.Descripcion
	resultado.TransaccionContable.Fecha = transaccionReversion.FechaTransaccion
	resultado.TransaccionContable.Movimientos, outputError = asientoContable.GetDetalleContable(transaccionReversion.Movimientos, nil)
	if outputError != nil {
		return outputError
	}

	if outputError = aplicarEstadoSalidaAnulada(trSalida.Salida, request.Observacion); outputError != nil {
		return outputError
	}
	if outputError = movimientosArka.PutMovimiento(trSalida.Salida, trSalida.Salida.Id); outputError != nil {
		return outputError
	}

	if entradaOriginal != nil {
		entradaActualizada, outputError = actualizarEstadoEntradaPadre(entradaOriginal.Id, salidaID)
		if outputError != nil {
			rollbackSalidaAnulada(salidaOriginal, resultado)
			return outputError
		}
	}

	if outputError = movimientosArka.PostMovimiento(movimientoReversion); outputError != nil {
		rollbackEntradaPadreAnulada(entradaOriginal, resultado)
		rollbackSalidaAnulada(salidaOriginal, resultado)
		return outputError
	}

	if _, outputError = movimientosContables.PostTrContable(transaccionReversion); outputError != nil {
		desactivarMovimientoReversionSalida(movimientoReversion, resultado)
		rollbackEntradaPadreAnulada(entradaOriginal, resultado)
		rollbackSalidaAnulada(salidaOriginal, resultado)
		return outputError
	}

	resultado.Salida = *trSalida.Salida
	if entradaActualizada != nil {
		resultado.Entrada = *entradaActualizada
	} else if entradaOriginal != nil {
		resultado.Entrada = *entradaOriginal
	}
	resultado.MovimientoReversion = *movimientoReversion

	return nil
}

func validarSalidaAnulable(trSalida *models.TrSalida, resultado *models.ResultadoAnulacionSalida) bool {
	if trSalida == nil || trSalida.Salida == nil || trSalida.Salida.EstadoMovimientoId == nil {
		resultado.Error = "No se pudo determinar el estado actual de la salida."
		return false
	}

	if trSalida.Salida.FormatoTipoMovimientoId == nil {
		resultado.Error = "No se pudo determinar el tipo de movimiento de la salida."
		return false
	}

	if trSalida.Salida.FormatoTipoMovimientoId.CodigoAbreviacion != "SAL" && trSalida.Salida.FormatoTipoMovimientoId.CodigoAbreviacion != "SAL_CONS" {
		resultado.Error = "El movimiento indicado no corresponde a una salida anulable."
		return false
	}

	if trSalida.Salida.EstadoMovimientoId.Nombre != estadoSalidaAprobada {
		resultado.Error = "Solo se pueden anular salidas en estado aprobada."
		return false
	}

	if trSalida.Salida.ConsecutivoId == nil || *trSalida.Salida.ConsecutivoId <= 0 {
		resultado.Error = "La salida no tiene consecutivo contable y no se puede reversar."
		return false
	}

	if len(trSalida.Elementos) == 0 {
		resultado.Error = "La salida no tiene elementos asociados y no se puede anular."
		return false
	}

	return true
}

func historialPermiteAnularSalida(historial *models.Historial, salidaID int) (bool, string) {
	if historial == nil {
		return false, "No se pudo consultar el historial actual de uno de los elementos de la salida."
	}

	if historial.Salida == nil || historial.Salida.Id != salidaID {
		return false, "La salida ya no es el último movimiento efectivo de uno de los elementos y no se puede anular."
	}

	if len(historial.Traslados) > 0 {
		return false, "La salida tiene elementos con traslados posteriores y no se puede anular."
	}

	if historial.Baja != nil {
		return false, "La salida tiene elementos con bajas posteriores y no se puede anular."
	}

	if len(historial.Novedades) > 0 {
		return false, "La salida tiene elementos con movimientos posteriores y no se puede anular."
	}

	return true, ""
}

func construirReversionSalida(
	trSalida *models.TrSalida,
	transaccionOriginal *models.TransaccionMovimientos,
	observacion string,
) (movimiento *models.Movimiento, transaccion *models.TransaccionMovimientos, outputError map[string]interface{}) {
	var (
		consecutivo         models.Consecutivo
		formatoMovimientoID int
		estadoMovimientoID  int
		elementosAjuste     []int
	)

	outputError = movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&formatoMovimientoID, formatoAjusteAutomatico)
	if outputError != nil {
		return nil, nil, outputError
	}

	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&estadoMovimientoID, estadoAjusteAprobado)
	if outputError != nil {
		return nil, nil, outputError
	}

	outputError = consecutivos.Get("contxtAjusteCons", "Anulación salida Arka", &consecutivo)
	if outputError != nil {
		return nil, nil, outputError
	}

	for _, elemento := range trSalida.Elementos {
		if elemento != nil && elemento.ElementoActaId != nil && *elemento.ElementoActaId > 0 {
			elementosAjuste = append(elementosAjuste, *elemento.ElementoActaId)
		}
	}

	detalle := models.FormatoAjusteAutomatico{Elementos: elementosAjuste}
	detalleJSON := ""
	outputError = utilsHelper.Marshal(detalle, &detalleJSON)
	if outputError != nil {
		return nil, nil, outputError
	}

	movimiento = &models.Movimiento{
		Observacion: descripcionReversionSalida(trSalida.Salida, observacion),
		Detalle:     detalleJSON,
		Activo:      true,
		MovimientoPadreId: &models.Movimiento{
			Id: trSalida.Salida.Id,
		},
		FormatoTipoMovimientoId: &models.FormatoTipoMovimiento{Id: formatoMovimientoID},
		EstadoMovimientoId:      &models.EstadoMovimiento{Id: estadoMovimientoID},
	}
	movimiento.Consecutivo = utilsHelper.String(consecutivos.Format("%05d", getTipoComprobanteAnulacionSalida(), &consecutivo))
	movimiento.ConsecutivoId = &consecutivo.Id

	transaccion = &models.TransaccionMovimientos{
		ConsecutivoId: consecutivo.Id,
		Activo:        true,
	}
	transaccion.Movimientos, outputError = invertirMovimientosContablesSalida(transaccionOriginal, descripcionReversionSalida(trSalida.Salida, observacion))
	if outputError != nil {
		return nil, nil, outputError
	}

	msg, outputError := asientoContable.CreateTransaccionContable(getTipoComprobanteAnulacionSalida(), "Reversa contable anulación salida", transaccion)
	if outputError != nil {
		return nil, nil, outputError
	}
	if msg != "" {
		return nil, nil, map[string]interface{}{
			"funcion": "construirReversionSalida - asientoContable.CreateTransaccionContable",
			"err":     msg,
			"status":  "400",
		}
	}

	return movimiento, transaccion, nil
}

func invertirMovimientosContablesSalida(original *models.TransaccionMovimientos, descripcion string) (movimientos []*models.MovimientoTransaccion, outputError map[string]interface{}) {
	if original == nil || len(original.Movimientos) == 0 {
		return nil, map[string]interface{}{
			"funcion": "invertirMovimientosContablesSalida - original",
			"err":     "transacción original vacía",
			"status":  "400",
		}
	}

	dbID, crID, outputError := parametros.GetParametrosDebitoCredito()
	if outputError != nil {
		return nil, outputError
	}

	movimientos = make([]*models.MovimientoTransaccion, 0, len(original.Movimientos))
	for _, movimientoOriginal := range original.Movimientos {
		if movimientoOriginal == nil {
			continue
		}

		tipoMovimiento := 0
		switch movimientoOriginal.TipoMovimientoId {
		case dbID:
			tipoMovimiento = crID
		case crID:
			tipoMovimiento = dbID
		default:
			return nil, map[string]interface{}{
				"funcion": "invertirMovimientosContablesSalida - TipoMovimientoId",
				"err":     "tipo de movimiento contable no reconocido",
				"status":  "400",
			}
		}

		movimiento := &models.MovimientoTransaccion{
			CuentaId:         movimientoOriginal.CuentaId,
			NombreCuenta:     movimientoOriginal.NombreCuenta,
			TipoMovimientoId: tipoMovimiento,
			Valor:            movimientoOriginal.Valor,
			Descripcion:      descripcion,
			Activo:           true,
			TerceroId:        normalizarTerceroIdSalida(movimientoOriginal.TerceroId),
		}

		movimientos = append(movimientos, movimiento)
	}

	return movimientos, nil
}

func aplicarEstadoSalidaAnulada(salida *models.Movimiento, observacion string) (outputError map[string]interface{}) {
	if salida == nil || salida.EstadoMovimientoId == nil {
		return map[string]interface{}{
			"funcion": "aplicarEstadoSalidaAnulada - salida.EstadoMovimientoId",
			"err":     "estado movimiento nil",
			"status":  "500",
		}
	}

	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&salida.EstadoMovimientoId.Id, estadoSalidaAnulada)
	if outputError != nil {
		return outputError
	}

	motivo := strings.TrimSpace(observacion)
	if motivo == "" {
		motivo = "Sin observación adicional."
	}
	salida.Observacion = strings.TrimSpace(salida.Observacion + "\nANULACION DE SALIDA: " + motivo)
	return nil
}

func actualizarEstadoEntradaPadre(entradaID, salidaID int) (entrada *models.Movimiento, outputError map[string]interface{}) {
	entrada, outputError = movimientosArka.GetMovimientoById(entradaID)
	if outputError != nil || entrada == nil {
		return
	}

	salidas, outputError := consultarSalidasAsociadasEntradaPadre(entradaID)
	if outputError != nil {
		return nil, outputError
	}

	if !debeRestaurarEntradaAprobada(salidas, salidaID) {
		return entrada, nil
	}

	if entrada.EstadoMovimientoId == nil {
		return nil, map[string]interface{}{
			"funcion": "actualizarEstadoEntradaPadre - entrada.EstadoMovimientoId",
			"err":     "estado movimiento nil",
			"status":  "500",
		}
	}

	if entrada.EstadoMovimientoId.Nombre != estadoEntradaConSalida && entrada.EstadoMovimientoId.Nombre != estadoEntradaAprobada {
		return entrada, nil
	}

	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&entrada.EstadoMovimientoId.Id, estadoEntradaAprobada)
	if outputError != nil {
		return nil, outputError
	}

	if outputError = movimientosArka.PutMovimiento(entrada, entrada.Id); outputError != nil {
		return nil, outputError
	}

	return entrada, nil
}

func consultarSalidasAsociadasEntradaPadre(entradaID int) (salidas []*models.Movimiento, outputError map[string]interface{}) {
	payload := "limit=-1&query=Activo:true,MovimientoPadreId__Id:" +
		fmt.Sprint(entradaID) +
		",FormatoTipoMovimientoId__CodigoAbreviacion__in:SAL|SAL_CONS"

	salidas, _, outputError = movimientosArka.GetAllMovimiento(payload)
	return
}

func debeRestaurarEntradaAprobada(salidas []*models.Movimiento, salidaID int) bool {
	for _, salida := range salidas {
		if salida == nil || salida.EstadoMovimientoId == nil {
			return false
		}

		if salida.Id == salidaID {
			continue
		}

		if salida.EstadoMovimientoId.Nombre != estadoSalidaAnulada && salida.EstadoMovimientoId.Nombre != estadoSalidaRechazada {
			return false
		}
	}

	return true
}

func rollbackSalidaAnulada(original *models.Movimiento, resultado *models.ResultadoAnulacionSalida) {
	if original == nil {
		return
	}

	if err := movimientosArka.PutMovimiento(original, original.Id); err != nil {
		appendAnulacionSalidaError(resultado, "No se pudo revertir el estado original de la salida.")
	}
}

func rollbackEntradaPadreAnulada(original *models.Movimiento, resultado *models.ResultadoAnulacionSalida) {
	if original == nil {
		return
	}

	if err := movimientosArka.PutMovimiento(original, original.Id); err != nil {
		appendAnulacionSalidaError(resultado, "No se pudo revertir el estado original de la entrada padre.")
	}
}

func desactivarMovimientoReversionSalida(movimiento *models.Movimiento, resultado *models.ResultadoAnulacionSalida) {
	if movimiento == nil || movimiento.Id <= 0 {
		return
	}

	movimiento.Activo = false
	if err := movimientosArka.PutMovimiento(movimiento, movimiento.Id); err != nil {
		appendAnulacionSalidaError(resultado, "No se pudo desactivar el movimiento de reversión generado.")
	}
}

func appendAnulacionSalidaError(resultado *models.ResultadoAnulacionSalida, msg string) {
	if resultado == nil || msg == "" {
		return
	}

	if resultado.Error == "" {
		resultado.Error = msg
		return
	}

	if !strings.Contains(resultado.Error, msg) {
		resultado.Error += " " + msg
	}
}

func descripcionReversionSalida(salida *models.Movimiento, observacion string) string {
	consecutivo := fmt.Sprint(salida.Id)
	if salida != nil && salida.Consecutivo != nil && *salida.Consecutivo != "" {
		consecutivo = *salida.Consecutivo
	}

	descripcion := "Reversa contable por anulación de salida " + consecutivo
	if strings.TrimSpace(observacion) != "" {
		descripcion += ". " + strings.TrimSpace(observacion)
	}
	return descripcion
}

func normalizarTerceroIdSalida(terceroId *int) *int {
	if terceroId == nil || *terceroId <= 0 {
		return nil
	}

	tercero := *terceroId
	return &tercero
}

func cloneMovimientoSalida(movimiento *models.Movimiento) *models.Movimiento {
	if movimiento == nil {
		return nil
	}

	clone := *movimiento
	if movimiento.ConsecutivoId != nil {
		consecutivoID := *movimiento.ConsecutivoId
		clone.ConsecutivoId = &consecutivoID
	}
	if movimiento.Consecutivo != nil {
		consecutivo := *movimiento.Consecutivo
		clone.Consecutivo = &consecutivo
	}
	if movimiento.FechaCorte != nil {
		fechaCorte := *movimiento.FechaCorte
		clone.FechaCorte = &fechaCorte
	}
	if movimiento.FormatoTipoMovimientoId != nil {
		formato := *movimiento.FormatoTipoMovimientoId
		clone.FormatoTipoMovimientoId = &formato
	}
	if movimiento.EstadoMovimientoId != nil {
		estado := *movimiento.EstadoMovimientoId
		clone.EstadoMovimientoId = &estado
	}
	if movimiento.MovimientoPadreId != nil {
		padre := *movimiento.MovimientoPadreId
		clone.MovimientoPadreId = &padre
	}

	return &clone
}

func getTipoComprobanteAnulacionSalida() string {
	return "N39"
}
