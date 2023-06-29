package depreciacionHelper

import (
	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/configuracion"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// GenerarCierre Crear el movimiento y transacción contable correspondientes al cierre a una fecha determinada
func GenerarCierre(info *models.InfoDepreciacion, resultado *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("GenerarCierre - Unhandled Error!", "500")

	var (
		detalle          models.FormatoDepreciacion
		parametros       []models.ParametroConfiguracion
		formatoCierre    int
		estadoMovimiento int
		transaccion      models.TransaccionMovimientos
		cuentas          map[string]models.CuentaContable
	)

	if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&formatoCierre, "CRR"); err != nil {
		return err
	}

	if err := movimientosArka.GetEstadoMovimientoIdByNombre(&estadoMovimiento, "Cierre En Curso"); err != nil {
		return err
	}

	if err := configuracion.GetAllParametro("Nombre__in:modificandoCuentas|cierreEnCurso&sortby=Nombre&order=desc&limit=2", &parametros); err != nil {
		return err
	} else if len(parametros) != 2 {
		resultado.Error = "No se pudo bloquear el sistema para iniciar el proceso de cierre. Contacte soporte."
		return
	} else if parametros[0].Valor == "true" {
		resultado.Error = "Cuentas en modificación. No se puede iniciar el proceso de cierre. Intente más tarde."
		return
	}

	parametros[1].Valor = "true"
	if err := configuracion.PutParametro(parametros[1].Id, &parametros[1]); err != nil {
		resultado.Error = "No se pudo bloquear el sistema para iniciar el proceso de cierre. Contacte soporte."
		return err
	}

	if err := calcularCierre(info.FechaCorte.Format("2006-01-02"), cuentas, &transaccion, resultado); err != nil || resultado.Error != "" || len(transaccion.Movimientos) == 0 {
		desbloquearSistema(parametros[1], *resultado)
		return err
	}

	if msg, err := asientoContable.CreateTransaccionContable(getTipoComprobanteCierre(), dscTransaccionCierre(), &transaccion); err != nil || msg != "" {
		resultado.Error = msg
		desbloquearSistema(parametros[1], *resultado)
		return err
	}

	if detalleContable, err := asientoContable.GetDetalleContable(transaccion.Movimientos, cuentas); err != nil {
		desbloquearSistema(parametros[1], *resultado)
		return err
	} else if len(detalleContable) > 0 {
		trContable := models.InfoTransaccionContable{
			Movimientos: detalleContable,
			Concepto:    transaccion.Descripcion,
		}
		resultado.TransaccionContable = trContable
	}

	if info.Id == 0 {
		var consecutivo_ models.Consecutivo
		outputError = consecutivos.Get("contxtMedicionesCons", "Registro cierre Arka", &consecutivo_)
		if outputError != nil {
			desbloquearSistema(parametros[1], *resultado)
			return
		}

		resultado.Movimiento.ConsecutivoId = &consecutivo_.Id
		resultado.Movimiento.Consecutivo = utilsHelper.String(consecutivos.Format("%02d", getTipoComprobanteCierre(), &consecutivo_))
	} else {
		if mov_, err := movimientosArka.GetMovimientoById(info.Id); err != nil {
			desbloquearSistema(parametros[1], *resultado)
			return err
		} else {
			resultado.Movimiento = *mov_
		}

		if err := utilsHelper.Unmarshal(resultado.Movimiento.Detalle, &detalle); err != nil {
			desbloquearSistema(parametros[1], *resultado)
			return err
		}
	}

	resultado.Movimiento.FormatoTipoMovimientoId = &models.FormatoTipoMovimiento{Id: formatoCierre}
	resultado.Movimiento.EstadoMovimientoId = &models.EstadoMovimiento{Id: estadoMovimiento}
	resultado.Movimiento.FechaCorte = &info.FechaCorte

	detalle.RazonRechazo = info.RazonRechazo

	if err := utilsHelper.Marshal(detalle, &resultado.Movimiento.Detalle); err != nil {
		desbloquearSistema(parametros[1], *resultado)
		return err
	}

	resultado.Movimiento.Observacion = info.Observaciones
	resultado.Movimiento.Activo = true

	if resultado.Movimiento.Id > 0 {
		if movimiento_, err := movimientosArka.PutMovimiento(&resultado.Movimiento, resultado.Movimiento.Id); err != nil {
			desbloquearSistema(parametros[1], *resultado)
			return err
		} else {
			resultado.Movimiento = *movimiento_
		}
	} else {
		if err := movimientosArka.PostMovimiento(&resultado.Movimiento); err != nil {
			desbloquearSistema(parametros[1], *resultado)
			return err
		}
	}

	return
}

func desbloquearSistema(parametro models.ParametroConfiguracion, resultado models.ResultadoMovimiento) {
	parametro.Valor = "false"
	if err := configuracion.PutParametro(parametro.Id, &parametro); err != nil {
		resultado.Error += " No se pudo desbloquear el sistema. Contacte soporte."
		return
	}
}
