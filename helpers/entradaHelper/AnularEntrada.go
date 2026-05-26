package entradaHelper

import (
	"fmt"
	"strings"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	crudActaRecibido "github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

const (
	estadoEntradaAprobada     = "Entrada Aprobada"
	estadoEntradaConSalida    = "Entrada Con Salida"
	estadoEntradaAnulada      = "Entrada anulada"
	estadoSalidaAnulada       = "Salida anulada"
	estadoActaAsociadaEntrada = "AsociadoEntrada"
	estadoActaEnVerificacion  = "EnVerificacion"
	formatoAjusteAutomatico   = "AAT"
	estadoAjusteAprobado      = "Ajuste Aprobado"
)

// AnularEntrada anula una entrada, devuelve el acta a verificación y registra la reversa contable.
func AnularEntrada(entradaID int, request *models.AnulacionEntradaRequest, resultado *models.ResultadoAnulacionEntrada) (outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("AnularEntrada - Unhandled Error!", "500")

	if resultado == nil {
		return map[string]interface{}{
			"funcion": "AnularEntrada - resultado",
			"err":     "resultado nil",
			"status":  "500",
		}
	}

	if request == nil {
		request = &models.AnulacionEntradaRequest{}
	}

	entrada, outputError := movimientosArka.GetMovimientoById(entradaID)
	if outputError != nil {
		return outputError
	}
	if entrada == nil {
		return map[string]interface{}{
			"funcion": "AnularEntrada - movimientosArka.GetMovimientoById",
			"err":     "movimiento nil",
			"status":  "404",
		}
	}

	salidas, outputError := consultarSalidasAsociadasEntrada(entradaID)
	if outputError != nil {
		return outputError
	}

	if !validarEntradaAnulable(entrada, salidas, resultado) {
		return nil
	}

	var formato models.FormatoBaseEntrada
	if outputError = utilsHelper.Unmarshal(entrada.Detalle, &formato); outputError != nil {
		return outputError
	}
	if formato.ActaRecibidoId <= 0 {
		resultado.Error = "La entrada no tiene un acta asociada y no se puede anular con este flujo."
		return nil
	}

	actaTransaccion := new(models.TransaccionActaRecibido)
	outputError = crudActaRecibido.GetTransaccionActaRecibidoById(formato.ActaRecibidoId, false, actaTransaccion)
	if outputError != nil {
		return outputError
	}
	if actaTransaccion == nil || actaTransaccion.UltimoEstado == nil || actaTransaccion.UltimoEstado.EstadoActaId == nil {
		resultado.Error = "No se pudo consultar el estado actual del acta asociada."
		return nil
	}
	if actaTransaccion.UltimoEstado.EstadoActaId.CodigoAbreviacion != estadoActaAsociadaEntrada {
		resultado.Error = "El acta asociada no está en estado 'Asociada a Entrada' y no se puede anular la entrada."
		return nil
	}

	if entrada.ConsecutivoId == nil || *entrada.ConsecutivoId <= 0 {
		resultado.Error = "La entrada no tiene consecutivo contable y no se puede reversar."
		return nil
	}

	transaccionOriginal, outputError := movimientosContables.GetTransaccion(*entrada.ConsecutivoId, "consecutivo", true)
	if outputError != nil {
		return outputError
	}
	if transaccionOriginal == nil || len(transaccionOriginal.Movimientos) == 0 {
		resultado.Error = "No se encontró la transacción contable original de la entrada."
		return nil
	}

	entradaOriginal := cloneMovimiento(entrada)
	actaOriginal := cloneTransaccionActa(actaTransaccion)

	movimientoReversion, transaccionReversion, outputError := construirReversionEntrada(entrada, formato, transaccionOriginal, request.Observacion)
	if outputError != nil {
		return outputError
	}
	if movimientoReversion == nil || transaccionReversion == nil {
		resultado.Error = "No se pudo construir la reversa contable de la entrada."
		return nil
	}

	resultado.TransaccionContable.Concepto = transaccionReversion.Descripcion
	resultado.TransaccionContable.Fecha = transaccionReversion.FechaTransaccion
	resultado.TransaccionContable.Movimientos, outputError = asientoContable.GetDetalleContable(transaccionReversion.Movimientos, nil)
	if outputError != nil {
		return outputError
	}

	if outputError = aplicarEstadoEntradaAnulada(entrada, request.Observacion); outputError != nil {
		return outputError
	}
	if outputError = movimientosArka.PutMovimiento(entrada, entradaID); outputError != nil {
		return outputError
	}

	if outputError = aplicarEstadoActaEnVerificacion(actaTransaccion); outputError != nil {
		rollbackEntradaAnulada(entradaOriginal, resultado)
		return outputError
	}
	if outputError = crudActaRecibido.PutTransaccionActaRecibido(formato.ActaRecibidoId, actaTransaccion); outputError != nil {
		rollbackEntradaAnulada(entradaOriginal, resultado)
		return outputError
	}

	if outputError = movimientosArka.PostMovimiento(movimientoReversion); outputError != nil {
		rollbackActaAnulada(formato.ActaRecibidoId, actaOriginal, resultado)
		rollbackEntradaAnulada(entradaOriginal, resultado)
		return outputError
	}

	if _, outputError = movimientosContables.PostTrContable(transaccionReversion); outputError != nil {
		desactivarMovimientoReversion(movimientoReversion, resultado)
		rollbackActaAnulada(formato.ActaRecibidoId, actaOriginal, resultado)
		rollbackEntradaAnulada(entradaOriginal, resultado)
		return outputError
	}

	resultado.Entrada = *entrada
	resultado.MovimientoReversion = *movimientoReversion
	return nil
}

func consultarSalidasAsociadasEntrada(entradaID int) (salidas []*models.Movimiento, outputError map[string]interface{}) {
	payload := "limit=-1&query=Activo:true,MovimientoPadreId__Id:" +
		fmt.Sprint(entradaID) +
		",FormatoTipoMovimientoId__CodigoAbreviacion__in:SAL|SAL_CONS"

	salidas, _, outputError = movimientosArka.GetAllMovimiento(payload)
	return
}

func validarEntradaAnulable(entrada *models.Movimiento, salidas []*models.Movimiento, resultado *models.ResultadoAnulacionEntrada) bool {
	if entrada == nil || entrada.EstadoMovimientoId == nil {
		resultado.Error = "No se pudo determinar el estado actual de la entrada."
		return false
	}

	if entrada.EstadoMovimientoId.Nombre != estadoEntradaAprobada && entrada.EstadoMovimientoId.Nombre != estadoEntradaConSalida {
		resultado.Error = "Solo se pueden anular entradas en estado aprobada o con salidas previamente anuladas."
		return false
	}

	for _, salida := range salidas {
		if salida == nil || salida.EstadoMovimientoId == nil {
			resultado.Error = "No se pudo validar el estado de las salidas asociadas."
			return false
		}
		if salida.EstadoMovimientoId.Nombre != estadoSalidaAnulada {
			resultado.Error = "La entrada tiene salidas asociadas que no están anuladas y no se puede anular."
			return false
		}
	}

	return true
}

func construirReversionEntrada(
	entrada *models.Movimiento,
	formato models.FormatoBaseEntrada,
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

	outputError = consecutivos.Get("contxtAjusteCons", "Anulación entrada Arka", &consecutivo)
	if outputError != nil {
		return nil, nil, outputError
	}

	if formato.ActaRecibidoId > 0 {
		query := "Activo:true,ActaRecibidoId__Id:" + fmt.Sprint(formato.ActaRecibidoId)
		if elementos, err := crudActaRecibido.GetAllElemento(query, "Id", "Id", "asc", "", "-1"); err != nil {
			return nil, nil, err
		} else {
			for _, elemento := range elementos {
				if elemento != nil && elemento.Id > 0 {
					elementosAjuste = append(elementosAjuste, elemento.Id)
				}
			}
		}
	}

	detalle := models.FormatoAjusteAutomatico{Elementos: elementosAjuste}
	detalleJSON := ""
	outputError = utilsHelper.Marshal(detalle, &detalleJSON)
	if outputError != nil {
		return nil, nil, outputError
	}

	movimiento = &models.Movimiento{
		Observacion: descripcionReversionEntrada(entrada, observacion),
		Detalle:     detalleJSON,
		Activo:      true,
		MovimientoPadreId: &models.Movimiento{
			Id: entrada.Id,
		},
		FormatoTipoMovimientoId: &models.FormatoTipoMovimiento{Id: formatoMovimientoID},
		EstadoMovimientoId:      &models.EstadoMovimiento{Id: estadoMovimientoID},
	}
	movimiento.Consecutivo = utilsHelper.String(consecutivos.Format("%05d", getTipoComprobanteAnulacionEntrada(), &consecutivo))
	movimiento.ConsecutivoId = &consecutivo.Id

	transaccion = &models.TransaccionMovimientos{
		ConsecutivoId: consecutivo.Id,
		Activo:        true,
	}
	transaccion.Movimientos, outputError = invertirMovimientosContables(transaccionOriginal, descripcionReversionEntrada(entrada, observacion))
	if outputError != nil {
		return nil, nil, outputError
	}

	msg, outputError := asientoContable.CreateTransaccionContable(getTipoComprobanteAnulacionEntrada(), "Reversa contable anulación entrada", transaccion)
	if outputError != nil {
		return nil, nil, outputError
	}
	if msg != "" {
		return nil, nil, map[string]interface{}{
			"funcion": "construirReversionEntrada - asientoContable.CreateTransaccionContable",
			"err":     msg,
			"status":  "400",
		}
	}

	return movimiento, transaccion, nil
}

func invertirMovimientosContables(original *models.TransaccionMovimientos, descripcion string) (movimientos []*models.MovimientoTransaccion, outputError map[string]interface{}) {
	if original == nil || len(original.Movimientos) == 0 {
		return nil, map[string]interface{}{
			"funcion": "invertirMovimientosContables - original",
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
				"funcion": "invertirMovimientosContables - TipoMovimientoId",
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
		}
		if movimientoOriginal.TerceroId != nil {
			terceroID := *movimientoOriginal.TerceroId
			movimiento.TerceroId = &terceroID
		}

		movimientos = append(movimientos, movimiento)
	}

	return movimientos, nil
}

func aplicarEstadoEntradaAnulada(entrada *models.Movimiento, observacion string) (outputError map[string]interface{}) {
	if entrada == nil || entrada.EstadoMovimientoId == nil {
		return map[string]interface{}{
			"funcion": "aplicarEstadoEntradaAnulada - entrada.EstadoMovimientoId",
			"err":     "estado movimiento nil",
			"status":  "500",
		}
	}

	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&entrada.EstadoMovimientoId.Id, estadoEntradaAnulada)
	if outputError != nil {
		return outputError
	}

	motivo := strings.TrimSpace(observacion)
	if motivo == "" {
		motivo = "Sin observación adicional."
	}
	entrada.Observacion = strings.TrimSpace(entrada.Observacion + "\nANULACION DE ENTRADA: " + motivo)
	return nil
}

func aplicarEstadoActaEnVerificacion(transaccion *models.TransaccionActaRecibido) (outputError map[string]interface{}) {
	if transaccion == nil || transaccion.UltimoEstado == nil || transaccion.UltimoEstado.EstadoActaId == nil {
		return map[string]interface{}{
			"funcion": "aplicarEstadoActaEnVerificacion - transaccion.UltimoEstado.EstadoActaId",
			"err":     "estado acta nil",
			"status":  "500",
		}
	}

	estadoID := 0
	outputError = crudActaRecibido.GetEstadoActaIdByCodigoAbreviacion(&estadoID, estadoActaEnVerificacion)
	if outputError != nil {
		return outputError
	}

	transaccion.UltimoEstado.Id = 0
	transaccion.UltimoEstado.EstadoActaId.Id = estadoID
	return nil
}

func rollbackEntradaAnulada(original *models.Movimiento, resultado *models.ResultadoAnulacionEntrada) {
	if original == nil {
		return
	}
	if err := movimientosArka.PutMovimiento(original, original.Id); err != nil {
		appendAnulacionError(resultado, "No se pudo revertir el estado original de la entrada.")
	}
}

func rollbackActaAnulada(actaID int, original *models.TransaccionActaRecibido, resultado *models.ResultadoAnulacionEntrada) {
	if original == nil || actaID <= 0 {
		return
	}
	if err := crudActaRecibido.PutTransaccionActaRecibido(actaID, original); err != nil {
		appendAnulacionError(resultado, "No se pudo revertir el estado original del acta.")
	}
}

func desactivarMovimientoReversion(movimiento *models.Movimiento, resultado *models.ResultadoAnulacionEntrada) {
	if movimiento == nil || movimiento.Id <= 0 {
		return
	}
	movimiento.Activo = false
	if err := movimientosArka.PutMovimiento(movimiento, movimiento.Id); err != nil {
		appendAnulacionError(resultado, "No se pudo desactivar el movimiento de reversión generado.")
	}
}

func appendAnulacionError(resultado *models.ResultadoAnulacionEntrada, msg string) {
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

func descripcionReversionEntrada(entrada *models.Movimiento, observacion string) string {
	consecutivo := fmt.Sprint(entrada.Id)
	if entrada != nil && entrada.Consecutivo != nil && *entrada.Consecutivo != "" {
		consecutivo = *entrada.Consecutivo
	}

	descripcion := "Reversa contable por anulación de entrada " + consecutivo
	if strings.TrimSpace(observacion) != "" {
		descripcion += ". " + strings.TrimSpace(observacion)
	}
	return descripcion
}

func getTipoComprobanteAnulacionEntrada() string {
	return "N39"
}

func cloneMovimiento(movimiento *models.Movimiento) *models.Movimiento {
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

func cloneTransaccionActa(transaccion *models.TransaccionActaRecibido) *models.TransaccionActaRecibido {
	if transaccion == nil {
		return nil
	}

	clone := *transaccion
	if transaccion.ActaRecibido != nil {
		acta := *transaccion.ActaRecibido
		clone.ActaRecibido = &acta
	}
	if transaccion.UltimoEstado != nil {
		historico := *transaccion.UltimoEstado
		clone.UltimoEstado = &historico
		if transaccion.UltimoEstado.EstadoActaId != nil {
			estado := *transaccion.UltimoEstado.EstadoActaId
			clone.UltimoEstado.EstadoActaId = &estado
		}
		if transaccion.UltimoEstado.ActaRecibidoId != nil {
			acta := *transaccion.UltimoEstado.ActaRecibidoId
			clone.UltimoEstado.ActaRecibidoId = &acta
		}
	}

	return &clone
}
