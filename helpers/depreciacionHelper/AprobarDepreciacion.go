package depreciacionHelper

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/configuracion"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// AprobarDepreciacion Registra las novedades para los elementos depreciados y realiza la transaccion contable
func AprobarDepreciacion(id int, resultado *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("AprobarDepreciacion - Unhandled Error!", "500")

	var (
		parametros        []models.ParametroConfiguracion
		detalle           models.FormatoDepreciacion
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

	if err := utilsHelper.Unmarshal(resultado.Movimiento.Detalle, &detalle); err != nil {
		return err
	}

	if err := calcularCierre(detalle.FechaCorte, &transaccionCierre.ElementoMovimientoId, cuentas, &transaccion, resultado); err != nil {
		return err
	} else if resultado.Error != "" || len(transaccionCierre.ElementoMovimientoId) == 0 || len(transaccion.Movimientos) == 0 {
		return
	}

	if msg, err := asientoContable.CreateTransaccionContable(getTipoComprobanteCierre(), dscTransaccionCierre(), &transaccion); err != nil || msg != "" {
		resultado.Error = msg
		return err
	}

	transaccion.ConsecutivoId = detalle.ConsecutivoId
	if _, err := movimientosContables.PostTrContable(&transaccion); err != nil {
		resultado.Error = "Error al registrar la transacción contable. Contacte soporte"
		return err
	}

	transaccionCierre.MovimientoId = id
	if err := movimientosArka.AprobarCierre(&transaccionCierre, &resultado.Movimiento); err != nil {
		resultado.Error = "Se registró la transacción contable pero no se pudo aprobar el cierre correctamente. Contacte soporte"
		return err
	}

	if detalleContable, err := asientoContable.GetDetalleContable(transaccion.Movimientos, cuentas); err != nil {
		return err
	} else if len(detalleContable) > 0 {
		trContable := models.InfoTransaccionContable{
			Movimientos: detalleContable,
			Concepto:    transaccion.Descripcion,
			Fecha:       transaccion.FechaTransaccion,
		}
		resultado.TransaccionContable = trContable
	}

	desbloquearSistema(parametros[0], *resultado)

	return
}
