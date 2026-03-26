package depreciacionHelper

import (
	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/configuracion"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// AprobarDepreciacion Registra las novedades para los elementos depreciados y realiza la transaccion contable
func AprobarDepreciacion(id int, resultado *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("AprobarDepreciacion - Unhandled Error!", "500")

	var (
		parametros  []models.ParametroConfiguracion
		transaccion models.TransaccionMovimientos
		cuentas     map[string]models.CuentaContable
	)

	outputError = configuracion.GetAllParametro("Nombre:cierreEnCurso", &parametros)
	if outputError != nil || len(parametros) != 1 || parametros[0].Valor != "true" {
		return
	}

	mov_, outputError := movimientosArka.GetMovimientoById(id)
	if outputError != nil {
		return
	} else if mov_.EstadoMovimientoId.Nombre != "Cierre En Curso" {
		resultado.Error = "El cierre no está en curso por lo que no puede ser aprobado."
		return
	}

	resultado.Movimiento = *mov_
	outputError = calcularCierre(resultado.Movimiento.FechaCorte.UTC().Format("2006-01-02"), cuentas, &transaccion, resultado)
	if outputError != nil || resultado.Error != "" || len(transaccion.Movimientos) == 0 {
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

	outputError = movimientosArka.AprobarCierre(&resultado.Movimiento)
	if outputError != nil {
		resultado.Error = "Se registró la transacción contable pero no se pudo aprobar el cierre correctamente. Contacte soporte"
		return
	}

	desbloquearSistema(parametros[0], *resultado)

	return
}
