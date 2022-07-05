package depreciacionHelper

import (
	"github.com/astaxie/beego"

	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// GenerarCierre Crear el movimiento y transacción contable correspondientes al cierre a una fecha determinada
func GenerarCierre(info *models.InfoDepreciacion, resultado *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GenerarCierre - Unhandled Error!", "500")

	var (
		movimiento       models.Movimiento
		detalle          models.FormatoDepreciacion
		formatoCierre    int
		estadoMovimiento int
	)

	if err := calcularCierre(info.FechaCorte.Format("2006-01-02"), nil, nil, resultado); err != nil {
		return err
	}

	if resultado.Error != "" || len(resultado.TransaccionContable.Movimientos) == 0 {
		return
	}

	if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&formatoCierre, "CRR"); err != nil {
		return err
	}

	if err := movimientosArka.GetEstadoMovimientoIdByNombre(&estadoMovimiento, "Cierre En Curso"); err != nil {
		return err
	}

	movimiento.FormatoTipoMovimientoId = &models.FormatoTipoMovimiento{Id: formatoCierre}
	movimiento.EstadoMovimientoId = &models.EstadoMovimiento{Id: estadoMovimiento}

	if info.Id == 0 {
		var consecutivo_ models.Consecutivo
		ctxt, _ := beego.AppConfig.Int("contxtMedicionesCons")
		if err := consecutivos.Get(ctxt, "Registro cierre Arka", &consecutivo_); err != nil {
			return outputError
		}

		detalle.ConsecutivoId = consecutivo_.Id
		detalle.Consecutivo = consecutivos.Format("%02d", getTipoComprobanteCierre(), &consecutivo_)
	} else {
		var (
			movimiento_ models.Movimiento
		)

		if mov_, err := movimientosArka.GetMovimientoById(info.Id); err != nil {
			return err
		} else {
			movimiento_ = *mov_
		}

		if err := utilsHelper.Unmarshal(movimiento_.Detalle, &detalle); err != nil {
			return err
		}
	}

	detalle.FechaCorte = info.FechaCorte.Format("2006-01-02")
	detalle.RazonRechazo = info.RazonRechazo
	if err := utilsHelper.Marshal(detalle, &movimiento.Detalle); err != nil {
		return err
	}

	movimiento.Observacion = info.Observaciones
	movimiento.Activo = true

	if info.Id > 0 {
		movimiento.Id = info.Id
		if movimiento_, err := movimientosArka.PutMovimiento(&movimiento, info.Id); err != nil {
			return err
		} else {
			resultado.Movimiento = *movimiento_
		}
	} else {
		if err := movimientosArka.PostMovimiento(&movimiento); err != nil {
			return err
		}
		resultado.Movimiento = movimiento
	}

	return
}
