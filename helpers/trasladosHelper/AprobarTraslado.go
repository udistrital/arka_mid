package trasladoshelper

import (
	"time"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/inventarioHelper"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// AprobarTraslado Actualiza el estado del traslado y genera la transaccion contable correspondiente
func AprobarTraslado(id int, response *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("AprobarTraslado - Unhandled Error!", "500")

	var (
		detalle     models.FormatoTraslado
		tipoSalida  int
		transaccion models.TransaccionMovimientos
	)

	movimiento_, outputError := movimientosArka.GetMovimientoById(id)
	if outputError != nil || movimiento_.EstadoMovimientoId.Nombre != "Traslado Confirmado" {
		return
	}

	response.Movimiento = *movimiento_
	if err := utilsHelper.Unmarshal(response.Movimiento.Detalle, &detalle); err != nil {
		return err
	}

	if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&tipoSalida, "SAL"); err != nil {
		return err
	}

	if err := movimientosArka.GetEstadoMovimientoIdByNombre(&response.Movimiento.EstadoMovimientoId.Id, "Traslado Aprobado"); err != nil {
		return err
	}

	bufferCuentas := make(map[string]models.CuentaContable)
	bufferSubgrupos := make(map[int]models.DetalleSubgrupo)
	for _, el := range detalle.Elementos {

		historial, err := movimientosArka.GetHistorialElemento(el, true)
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
		outputError = actaRecibido.GetElementoById(historial.Elemento.ElementoActaId, &elementoActa)
		if outputError != nil {
			return
		}

		elementoActa.ValorTotal = valor
		elementosActa := []*models.Elemento{&elementoActa}
		tipoEntrada := historial.Salida.MovimientoPadreId.FormatoTipoMovimientoId.Id

		response.Error, outputError = asientoContable.CalcularMovimientosContables(elementosActa, descMovDestino(), tipoEntrada, tipoSalida, detalle.FuncionarioDestino, detalle.FuncionarioOrigen, bufferCuentas, bufferSubgrupos, &transaccion.Movimientos)
		if outputError != nil || response.Error != "" {
			return
		}
	}

	transaccion.ConsecutivoId = *response.Movimiento.ConsecutivoId
	response.Error, outputError = asientoContable.CreateTransaccionContable(getTipoComprobanteTraslados(), "Traslado de elementos", &transaccion)
	if outputError != nil || response.Error != "" {
		return
	}

	response.TransaccionContable.Concepto = transaccion.Descripcion
	response.TransaccionContable.Fecha = transaccion.FechaTransaccion
	response.TransaccionContable.Movimientos, outputError = asientoContable.GetDetalleContable(transaccion.Movimientos, bufferCuentas)
	if outputError != nil {
		return
	}

	_, outputError = movimientosContables.PostTrContable(&transaccion)
	if outputError != nil {
		return
	}

	response.Movimiento.FechaCorte = utilsHelper.Time(time.Now())
	_, outputError = movimientosArka.PutMovimiento(&response.Movimiento, response.Movimiento.Id)

	return
}

func descMovDestino() string {
	return "Traslado de elementos"
}
