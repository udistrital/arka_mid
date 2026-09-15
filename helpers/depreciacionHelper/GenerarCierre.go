package depreciacionHelper

import (
	"context"
	"time"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/configuracion"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

var launchCierreAsync = func(fn func()) {
	go fn()
}

// GenerarCierre Crear el movimiento y transacción contable correspondientes al cierre a una fecha determinada
func GenerarCierre(ctx context.Context, info *models.InfoDepreciacion, resultado *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("GenerarCierre - Unhandled Error!", "500")

	var (
		detalle          models.FormatoDepreciacion
		parametros       []models.ParametroConfiguracion
		formatoCierre    int
		estadoMovimiento int
	)

	if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(ctx, &formatoCierre, "CRR"); err != nil {
		return err
	}

	if err := movimientosArka.GetEstadoMovimientoIdByNombre(ctx, &estadoMovimiento, "Cierre En Curso"); err != nil {
		return err
	}

	if err := configuracion.GetAllParametro(ctx, "Nombre__in:modificandoCuentas|cierreEnCurso&sortby=Nombre&order=desc&limit=2", &parametros); err != nil {
		return err
	} else if len(parametros) != 2 {
		resultado.Error = "No se pudo bloquear el sistema para iniciar el proceso de cierre. Contacte soporte."
		return
	} else if parametros[0].Valor == "true" {
		resultado.Error = "Cuentas en modificación. No se puede iniciar el proceso de cierre. Intente más tarde."
		return
	}

	parametros[1].Valor = "true"
	if err := configuracion.PutParametro(ctx, parametros[1].Id, &parametros[1]); err != nil {
		resultado.Error = "No se pudo bloquear el sistema para iniciar el proceso de cierre. Contacte soporte."
		return err
	}

	if info.Id == 0 {
		var consecutivo_ models.Consecutivo
		outputError = consecutivos.Get(ctx, "contxtMedicionesCons", "Registro cierre Arka", &consecutivo_)
		if outputError != nil {
			desbloquearSistema(ctx, parametros[1], *resultado)
			return
		}

		resultado.Movimiento.ConsecutivoId = &consecutivo_.Id
		resultado.Movimiento.Consecutivo = utilsHelper.String(consecutivos.Format("%02d", getTipoComprobanteCierre(), &consecutivo_))
	} else {
		if mov_, err := movimientosArka.GetMovimientoById(ctx, info.Id); err != nil {
			desbloquearSistema(ctx, parametros[1], *resultado)
			return err
		} else {
			resultado.Movimiento = *mov_
		}

		if err := utilsHelper.Unmarshal(resultado.Movimiento.Detalle, &detalle); err != nil {
			desbloquearSistema(ctx, parametros[1], *resultado)
			return err
		}
	}

	resultado.Movimiento.FormatoTipoMovimientoId = &models.FormatoTipoMovimiento{Id: formatoCierre}
	resultado.Movimiento.EstadoMovimientoId = &models.EstadoMovimiento{Id: estadoMovimiento}
	resultado.Movimiento.FechaCorte = &info.FechaCorte

	detalle.RazonRechazo = info.RazonRechazo
	detalle.CalculoListo = false
	detalle.CalculoError = ""
	detalle.ElementosCalculados = 0
	detalle.MovimientosContables = 0
	detalle.FechaCalculo = nil
	detalle.Transaccion = nil
	detalle.PreviewContable = nil

	if err := utilsHelper.Marshal(detalle, &resultado.Movimiento.Detalle); err != nil {
		desbloquearSistema(ctx, parametros[1], *resultado)
		return err
	}

	resultado.Movimiento.Observacion = info.Observaciones
	resultado.Movimiento.Activo = true

	if resultado.Movimiento.Id > 0 {
		outputError = movimientosArka.PutMovimiento(ctx, &resultado.Movimiento, resultado.Movimiento.Id)
		if outputError != nil {
			desbloquearSistema(ctx, parametros[1], *resultado)
			return
		}
	} else {
		if err := movimientosArka.PostMovimiento(ctx, &resultado.Movimiento); err != nil {
			desbloquearSistema(ctx, parametros[1], *resultado)
			return err
		}
	}

	resultado.TransaccionContable = models.InfoTransaccionContable{}
	resultado.Error = ""

	movimientoID := resultado.Movimiento.Id
	fechaCorte := info.FechaCorte
	launchCierreAsync(func() {
		procesarCierreDepreciacionAsync(ctx, movimientoID, fechaCorte)
	})

	return
}

func desbloquearSistema(ctx context.Context, parametro models.ParametroConfiguracion, resultado models.ResultadoMovimiento) {
	parametro.Valor = "false"
	if err := configuracion.PutParametro(ctx, parametro.Id, &parametro); err != nil {
		resultado.Error += " No se pudo desbloquear el sistema. Contacte soporte."
		return
	}
}

func procesarCierreDepreciacionAsync(ctx context.Context, movimientoID int, fechaCorte time.Time) {
	movimiento, outputError := movimientosArka.GetMovimientoById(ctx, movimientoID)
	if outputError != nil || movimiento == nil {
		return
	}

	var (
		detalle     models.FormatoDepreciacion
		resultado   models.ResultadoMovimiento
		transaccion models.TransaccionMovimientos
		cuentas     map[string]models.CuentaContable
	)

	if err := utilsHelper.Unmarshal(movimiento.Detalle, &detalle); err != nil {
		persistirFalloCalculoCierreAsync(ctx, movimiento, "No se pudo leer el detalle del cierre para procesar la depreciación.")
		return
	}

	transaccion.ConsecutivoId = 0
	if movimiento.ConsecutivoId != nil {
		transaccion.ConsecutivoId = *movimiento.ConsecutivoId
	}
	if movimiento.FechaCorte != nil && !movimiento.FechaCorte.IsZero() {
		transaccion.FechaTransaccion = *movimiento.FechaCorte
	} else {
		transaccion.FechaTransaccion = fechaCorte
	}

	outputError = calcularCierre(ctx, fechaCorte.Format("2006-01-02"), cuentas, &transaccion, &resultado)
	if outputError != nil {
		persistirFalloCalculoCierreAsync(ctx, movimiento, "No se pudo calcular el cierre de depreciación.")
		return
	}
	if resultado.Error != "" {
		persistirFalloCalculoCierreAsync(ctx, movimiento, resultado.Error)
		return
	}
	if len(transaccion.Movimientos) == 0 {
		persistirFalloCalculoCierreAsync(ctx, movimiento, "No se generaron movimientos contables para la fecha de corte indicada.")
		return
	}

	if msg, err := asientoContable.CreateTransaccionContable(ctx, getTipoComprobanteCierre(), dscTransaccionCierre(), &transaccion); err != nil || msg != "" {
		if msg == "" {
			msg = "No se pudo construir la transacción contable del cierre."
		}
		persistirFalloCalculoCierreAsync(ctx, movimiento, msg)
		return
	}

	detalleContable, outputError := asientoContable.GetDetalleContable(ctx, transaccion.Movimientos, cuentas)
	if outputError != nil {
		persistirFalloCalculoCierreAsync(ctx, movimiento, "No se pudo construir la vista previa contable del cierre.")
		return
	}

	now := time.Now().UTC()
	detalle.CalculoListo = true
	detalle.CalculoError = ""
	detalle.ElementosCalculados = contarMovimientosConValor(transaccion.Movimientos)
	detalle.MovimientosContables = len(transaccion.Movimientos)
	detalle.FechaCalculo = &now
	detalle.Transaccion = &transaccion
	detalle.PreviewContable = &models.InfoTransaccionContable{
		Movimientos: detalleContable,
		Concepto:    transaccion.Descripcion,
		Fecha:       transaccion.FechaTransaccion,
	}

	if err := utilsHelper.Marshal(detalle, &movimiento.Detalle); err != nil {
		persistirFalloCalculoCierreAsync(ctx, movimiento, "No se pudo persistir el detalle calculado del cierre.")
		return
	}

	_ = movimientosArka.PutMovimiento(ctx, movimiento, movimiento.Id)
}

func persistirFalloCalculoCierreAsync(ctx context.Context, movimiento *models.Movimiento, mensaje string) {
	if movimiento == nil {
		return
	}

	var detalle models.FormatoDepreciacion
	_ = utilsHelper.Unmarshal(movimiento.Detalle, &detalle)

	detalle.CalculoListo = false
	detalle.CalculoError = mensaje
	detalle.Transaccion = nil
	detalle.PreviewContable = nil
	detalle.ElementosCalculados = 0
	detalle.MovimientosContables = 0

	if err := utilsHelper.Marshal(detalle, &movimiento.Detalle); err != nil {
		return
	}

	if movimiento.EstadoMovimientoId == nil {
		movimiento.EstadoMovimientoId = &models.EstadoMovimiento{}
	}
	if err := movimientosArka.GetEstadoMovimientoIdByNombre(ctx, &movimiento.EstadoMovimientoId.Id, "Cierre Rechazado"); err == nil {
		movimiento.EstadoMovimientoId.Nombre = "Cierre Rechazado"
	}
	_ = movimientosArka.PutMovimiento(ctx, movimiento, movimiento.Id)

	var parametros []models.ParametroConfiguracion
	if err := configuracion.GetAllParametro(ctx, "Nombre:cierreEnCurso", &parametros); err == nil && len(parametros) == 1 {
		desbloquearSistema(ctx, parametros[0], models.ResultadoMovimiento{})
	}
}

func contarMovimientosConValor(movimientos []*models.MovimientoTransaccion) int {
	total := 0
	for _, movimiento := range movimientos {
		if movimiento == nil || movimiento.Valor == 0 {
			continue
		}
		total++
	}
	return total
}
