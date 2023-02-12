package depreciacionHelper

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/configuracion"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// AprobarDepreciacion Registra las novedades para los elementos depreciados y realiza la transaccion contable
func AprobarDepreciacion(id int, resultado *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("AprobarDepreciacion - Unhandled Error!", "500")

	var (
		parametros        []models.ParametroConfiguracion
		transaccionCierre models.TransaccionCierre
		transaccion       models.TransaccionMovimientos
		cuentas           map[string]models.CuentaContable
	)

	if err := configuracion.GetAllParametro("Nombre:cierreEnCurso", &parametros); err != nil {
		return err
	} else if len(parametros) != 1 || parametros[0].Valor != "true" {
		return
	}

	if mov_, err := movimientosArka.GetAllMovimiento("limit=1&query=Id:" + strconv.Itoa(id)); err != nil {
		return err
	} else if len(mov_) == 1 && mov_[0].EstadoMovimientoId.Nombre == "Cierre En Curso" {
		resultado.Movimiento = *mov_[0]
	} else {
		resultado.Error = "No se pudo consultar la información del cierre. Contacte soporte."
		return
	}

	if err := calcularCierre(resultado.Movimiento.FechaCorte.UTC().Format("2006-01-02"), &transaccionCierre.ElementoMovimientoId, cuentas, &transaccion, resultado); err != nil {
		return err
	} else if resultado.Error != "" || len(transaccionCierre.ElementoMovimientoId) == 0 || len(transaccion.Movimientos) == 0 {
		return
	}

	resultado.Error, outputError = asientoContable.CreateTransaccionContable(getTipoComprobanteCierre(), dscTransaccionCierre(), &transaccion)
	if outputError != nil || resultado.Error != "" {
		return
	}

	resultado.TransaccionContable.Concepto = transaccion.Descripcion
	resultado.TransaccionContable.Fecha = transaccion.FechaTransaccion
	resultado.TransaccionContable.Movimientos, outputError = asientoContable.GetDetalleContable(transaccion.Movimientos, cuentas)
	if outputError != nil {
		return
	}

	transaccion.ConsecutivoId = *resultado.Movimiento.ConsecutivoId
	_, outputError = movimientosContables.PostTrContable(&transaccion)
	if outputError != nil {
		resultado.Error = "Error al registrar la transacción contable. Contacte soporte"
		return
	}

	transaccionCierre.MovimientoId = id
	outputError = movimientosArka.AprobarCierre(&transaccionCierre, &resultado.Movimiento)
	if outputError != nil {
		resultado.Error = "Se registró la transacción contable pero no se pudo aprobar el cierre correctamente. Contacte soporte"
		return
	}

	desbloquearSistema(parametros[0], *resultado)

	return
}
