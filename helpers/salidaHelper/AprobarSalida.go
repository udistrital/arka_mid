package salidaHelper

import (
	"time"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// AprobarSalida Aprobacion de una salida
func AprobarSalida(salidaId int, res *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("AprobarSalida - Unhandled Error!", "500")

	var (
		salida         models.FormatoSalida
		tipoMovimiento int
	)

	trSalida, outputError := movimientosArka.GetTrSalida(salidaId)
	if outputError != nil || trSalida.Salida.EstadoMovimientoId.Nombre != "Salida En Trámite" {
		return
	} else if len(trSalida.Elementos) == 0 || trSalida.Salida.ConsecutivoId == nil || *trSalida.Salida.ConsecutivoId == 0 {
		res.Error = "No se pudo continuar con la transacción contable. Contacte soporte."
		return
	}

	res.Movimiento = *trSalida.Salida
	outputError = utilsHelper.Unmarshal(trSalida.Salida.Detalle, &salida)
	if outputError != nil {
		return
	}

	if salida.Funcionario == 0 || salida.Ubicacion == 0 {
		res.Error = "No se pudo continuar con la transacción contable. Contacte soporte."
		return
	}

	outputError = movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&tipoMovimiento, "SAL")
	if outputError != nil {
		return
	}

	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&trSalida.Salida.EstadoMovimientoId.Id, "Salida Aprobada")
	if outputError != nil {
		return
	}

	var idsElementos []int
	for _, el := range trSalida.Elementos {
		idsElementos = append(idsElementos, *el.ElementoActaId)
	}

	query := "Id__in:" + utilsHelper.ArrayToString(idsElementos, "|")
	elementosActa, outputError := actaRecibido.GetAllElemento(query, "ValorUnitario,ValorTotal,SubgrupoCatalogoId,TipoBienId", "SubgrupoCatalogoId", "desc", "", "-1")
	if outputError != nil {
		return
	}

	dsc := ""
	if trSalida.Salida.MovimientoPadreId != nil && trSalida.Salida.MovimientoPadreId.Consecutivo != nil {
		dsc = "Entrada: " + *trSalida.Salida.MovimientoPadreId.Consecutivo
	}

	bufferCuentas := make(map[string]models.CuentaContable)
	transaccion := models.TransaccionMovimientos{ConsecutivoId: *trSalida.Salida.ConsecutivoId}
	res.Error, outputError = asientoContable.CalcularMovimientosContables(elementosActa, dsc, res.Movimiento.MovimientoPadreId.FormatoTipoMovimientoId.Id, tipoMovimiento, salida.Funcionario, salida.Funcionario, bufferCuentas, nil, &transaccion.Movimientos)
	if outputError != nil || res.Error != "" {
		return
	}

	res.Error, outputError = asientoContable.CreateTransaccionContable(getTipoComprobanteSalidas(), "Salida de Almacén", &transaccion)
	if outputError != nil || res.Error != "" {
		return
	}

	res.TransaccionContable.Movimientos, outputError = asientoContable.GetDetalleContable(transaccion.Movimientos, bufferCuentas)
	if outputError != nil {
		return
	}

	_, outputError = movimientosContables.PostTrContable(&transaccion)
	if outputError != nil {
		return
	}

	res.TransaccionContable.Concepto = transaccion.Descripcion
	res.TransaccionContable.Fecha = transaccion.FechaTransaccion
	trSalida.Salida.FechaCorte = utilsHelper.Time(time.Now())
	trSalida.Salida, outputError = movimientosArka.PutMovimiento(trSalida.Salida, trSalida.Salida.Id)

	return
}
