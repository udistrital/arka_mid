package trasladoshelper

import (
	"context"
	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/inventarioHelper"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	errorCtrl "github.com/udistrital/utils_oas/v2/errorctrl"
	timebogota "github.com/udistrital/utils_oas/v2/time_bogota"
)

// AprobarTraslado Actualiza el estado del traslado y genera la transaccion contable correspondiente
func AprobarTraslado(ctx context.Context, id int, response *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("AprobarTraslado - Unhandled Error!", "500")

	var (
		detalle     models.FormatoTraslado
		tipoSalida  int
		transaccion models.TransaccionMovimientos
	)

	movimiento_, outputError := movimientosArka.GetMovimientoById(ctx, id)
	if outputError != nil || movimiento_.EstadoMovimientoId.Nombre != "Traslado Confirmado" {
		return
	}

	response.Movimiento = *movimiento_
	if err := utilsHelper.Unmarshal(response.Movimiento.Detalle, &detalle); err != nil {
		return err
	}

	if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(ctx, &tipoSalida, "SAL"); err != nil {
		return err
	}

	if err := movimientosArka.GetEstadoMovimientoIdByNombre(ctx, &response.Movimiento.EstadoMovimientoId.Id, "Traslado Aprobado"); err != nil {
		return err
	}

	bufferCuentas := make(map[string]models.CuentaContable)
	bufferSubgrupos := make(map[int]models.DetalleSubgrupo)
	for _, el := range detalle.Elementos {

		historial, err := movimientosArka.GetHistorialElemento(ctx, el, true)
		if err != nil {
			return err
		} else if historial == nil {
			response.Error = "No se pudo la parametrización de los elementos. Contacte soporte."
			return
		}

		valor, _, _, _, err := inventarioHelper.GetUltimoValor(*historial)
		if err != nil {
			return err
		}

		if valor <= 0 {
			continue
		}

		var elementoActa models.Elemento
		outputError = actaRecibido.GetElementoById(ctx, *historial.Elemento.ElementoActaId, &elementoActa)
		if outputError != nil {
			return
		}

		elementoActa.ValorTotal = valor
		elementosActa := []*models.Elemento{&elementoActa}
		tipoEntrada := historial.Salida.MovimientoPadreId.FormatoTipoMovimientoId.Id

		response.Error, outputError = asientoContable.CalcularMovimientosContables(ctx, elementosActa, descMovDestino(), tipoEntrada, tipoSalida, detalle.FuncionarioDestino, detalle.FuncionarioOrigen, bufferCuentas, bufferSubgrupos, &transaccion.Movimientos)
		if outputError != nil || response.Error != "" {
			return
		}
	}

	transaccion.ConsecutivoId = *response.Movimiento.ConsecutivoId
	response.Error, outputError = asientoContable.CreateTransaccionContable(ctx, getTipoComprobanteTraslados(), "Traslado de elementos", &transaccion)
	if outputError != nil || response.Error != "" {
		return
	}

	response.TransaccionContable.Concepto = transaccion.Descripcion
	response.TransaccionContable.Fecha = transaccion.FechaTransaccion
	response.TransaccionContable.Movimientos, outputError = asientoContable.GetDetalleContable(ctx, transaccion.Movimientos, bufferCuentas)
	if outputError != nil {
		return
	}

	_, outputError = movimientosContables.PostTrContable(ctx, &transaccion)
	if outputError != nil {
		return
	}

	response.Movimiento.FechaCorte = utilsHelper.Time(timebogota.TiempoBogota())
	outputError = movimientosArka.PutMovimiento(ctx, &response.Movimiento, response.Movimiento.Id)

	return
}

func descMovDestino() string {
	return "Traslado de elementos"
}
