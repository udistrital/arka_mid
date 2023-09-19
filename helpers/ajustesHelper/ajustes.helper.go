package ajustesHelper

import (
	"time"

	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

func PostAjuste(trContable *models.PreTrAjuste) (movimiento *models.Movimiento, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("PostAjuste - Unhandled Error!", "500")

	movimiento = &models.Movimiento{
		FormatoTipoMovimientoId: &models.FormatoTipoMovimiento{},
		EstadoMovimientoId:      &models.EstadoMovimiento{},
	}

	detalle := &models.FormatoAjuste{PreTrAjuste: trContable}
	outputError = utilsHelper.Marshal(detalle, &movimiento.Detalle)
	if outputError != nil {
		return
	}

	outputError = movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&movimiento.FormatoTipoMovimientoId.Id, "AJ_CBE")
	if outputError != nil {
		return
	}

	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&movimiento.EstadoMovimientoId.Id, "Ajuste En Trámite")
	if outputError != nil {
		return
	}

	var consecutivo models.Consecutivo
	outputError = consecutivos.Get("contxtAjusteCons", "Ajuste Contable Arka", &consecutivo)
	if outputError != nil {
		return
	}

	movimiento.Consecutivo = utilsHelper.String(consecutivos.Format("%05d", getTipoComprobanteAjustes(), &consecutivo))
	movimiento.ConsecutivoId = &consecutivo.Id
	movimiento.Activo = true

	outputError = movimientosArka.PostMovimiento(movimiento)

	return
}

// GetDetalleAjuste Consulta los detalles de un ajuste contable
func GetDetalleAjuste(id int) (Ajuste *models.DetalleAjuste, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("GetDetalleAjuste - Unhandled Error!", "500")

	var (
		detalle     models.FormatoAjuste
		movimientos []*models.PreMovAjuste
	)

	Ajuste = new(models.DetalleAjuste)

	movimiento, outputError := movimientosArka.GetMovimientoById(id)
	if outputError != nil {
		return
	}

	outputError = utilsHelper.Unmarshal(movimiento.Detalle, &detalle)
	if outputError != nil {
		return
	}

	parametroDebitoId, parametroCreditoId, outputError := parametros.GetParametrosDebitoCredito()
	if outputError != nil {
		return
	}

	if detalle.PreTrAjuste != nil && movimiento.EstadoMovimientoId.Nombre != "Ajuste Aprobado" {
		movimientos = detalle.PreTrAjuste.Movimientos
	} else if movimiento.EstadoMovimientoId.Nombre == "Ajuste Aprobado" && movimiento.ConsecutivoId != nil && *movimiento.ConsecutivoId > 0 {
		if tr, err := movimientosContables.GetTransaccion(*movimiento.ConsecutivoId, "consecutivo", true); err != nil {
			return nil, err
		} else {
			for _, mov := range tr.Movimientos {
				mov_ := new(models.PreMovAjuste)
				mov_.Cuenta = mov.CuentaId
				if mov.TerceroId != nil {
					mov_.TerceroId = *mov.TerceroId
				}
				mov_.Descripcion = mov.Descripcion

				if mov.TipoMovimientoId == parametroCreditoId {
					mov_.Credito = mov.Valor
				} else if mov.TipoMovimientoId == parametroDebitoId {
					mov_.Debito = mov.Valor
				}

				movimientos = append(movimientos, mov_)
			}
		}
	}

	movs := make([]*models.DetalleMovimientoContable, 0)
	for _, mov := range movimientos {
		mov_ := new(models.DetalleMovimientoContable)

		if cta, err := cuentasContables.GetCuentaContable(mov.Cuenta); err != nil {
			return nil, err
		} else {
			mov_.Cuenta = &models.DetalleCuenta{
				Id:              cta.Id,
				Codigo:          cta.Codigo,
				Nombre:          cta.Nombre,
				RequiereTercero: cta.RequiereTercero,
			}
		}

		if mov.TerceroId > 0 {
			if tercero, err := terceros.GetNombreTerceroById(mov.TerceroId); err != nil {
				return nil, err
			} else {
				mov_.TerceroId = tercero
			}
		}
		mov_.Credito = mov.Credito
		mov_.Debito = mov.Debito
		mov_.Descripcion = mov.Descripcion
		movs = append(movs, mov_)
	}

	Ajuste.TrContable = movs
	Ajuste.Movimiento = movimiento

	return Ajuste, nil

}

// AprobarAjuste Realiza la transacción contable correspondiente
func AprobarAjuste(id int) (movimiento *models.Movimiento, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("AprobarAjuste - Unhandled Error!", "500")

	movimiento, outputError = movimientosArka.GetMovimientoById(id)
	if outputError != nil {
		return
	}

	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&movimiento.EstadoMovimientoId.Id, "Ajuste Aprobado")
	if outputError != nil {
		return
	}

	var detalle models.FormatoAjuste
	if err := utilsHelper.Unmarshal(movimiento.Detalle, &detalle); err != nil {
		return nil, err
	}

	parametroDebitoId, parametroCreditoId, outputError := parametros.GetParametrosDebitoCredito()
	if outputError != nil {
		return
	}

	movs := make([]*models.MovimientoTransaccion, 0)
	for _, mov := range detalle.PreTrAjuste.Movimientos {
		mov_ := new(models.MovimientoTransaccion)
		var cta *models.DetalleCuenta

		if ctaCr_, err := cuentasContables.GetCuentaContable(mov.Cuenta); err != nil {
			return nil, err
		} else {
			outputError = utilsHelper.FillStruct(ctaCr_, &cta)
			if outputError != nil {
				return
			}
			mov_.CuentaId = cta.Id
			mov_.NombreCuenta = cta.Nombre
		}

		if mov.TerceroId > 0 {
			mov_.TerceroId = &mov.TerceroId
		}

		if mov.Credito > 0 {
			mov_.TipoMovimientoId = parametroCreditoId
			mov_.Valor = mov.Credito
		} else if mov.Debito > 0 {
			mov_.TipoMovimientoId = parametroDebitoId
			mov_.Valor = mov.Debito
		}

		mov_.Activo = true
		movs = append(movs, mov_)
	}

	transaccion := new(models.TransaccionMovimientos)

	transaccion.ConsecutivoId = *movimiento.ConsecutivoId
	transaccion.Movimientos = movs
	transaccion.FechaTransaccion = time.Now()
	transaccion.Activo = true
	transaccion.Etiquetas = ""
	transaccion.Descripcion = ""

	if _, err := movimientosContables.PostTrContable(transaccion); err != nil {
		return nil, err
	}

	movimiento.Detalle = "{}"
	outputError = movimientosArka.PutMovimiento(movimiento, movimiento.Id)

	return
}
