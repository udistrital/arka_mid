package depreciacionHelper

import (
	"github.com/udistrital/arka_mid/helpers/crud/configuracion"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// AprobarDepreciacion Registra las novedades para los elementos depreciados y realiza la transaccion contable
func AprobarDepreciacion(id int, resultado *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("AprobarDepreciacion - Unhandled Error!", "500")

	var (
		parametros []models.ParametroConfiguracion
		detalle    models.FormatoDepreciacion
	)

	outputError = configuracion.GetAllParametro("Nombre:cierreEnCurso", &parametros)
	if outputError != nil {
		return
	}
	if len(parametros) != 1 || parametros[0].Valor != "true" {
		resultado.Error = "No hay un cierre en curso listo para aprobación."
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
	if err := utilsHelper.Unmarshal(resultado.Movimiento.Detalle, &detalle); err != nil {
		return err
	}
	if detalle.CalculoError != "" {
		resultado.Error = "El cierre no está listo para aprobación: " + detalle.CalculoError
		return
	}
	if !detalle.CalculoListo || detalle.Transaccion == nil || len(detalle.Transaccion.Movimientos) == 0 {
		resultado.Error = "El cálculo del cierre aún está en proceso o no produjo una vista previa aprobable."
		return
	}

	transaccion := *detalle.Transaccion
	transaccion.ConsecutivoId = *resultado.Movimiento.ConsecutivoId
	if transaccion.FechaTransaccion.IsZero() && resultado.Movimiento.FechaCorte != nil {
		transaccion.FechaTransaccion = *resultado.Movimiento.FechaCorte
	}

	if detalle.PreviewContable != nil {
		resultado.TransaccionContable = *detalle.PreviewContable
	}
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
