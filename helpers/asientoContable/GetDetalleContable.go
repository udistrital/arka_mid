package asientoContable

import (
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// GetDetalleContable Consulta los detalles de una transacción contable para ser mostrada en el cliente
func GetDetalleContable(movimientos []*models.MovimientoTransaccion, detalleCuentas map[string]models.CuentaContable) (movimientos_ []*models.DetalleMovimientoContable, outputError map[string]interface{}) {

	funcion := "GetDetalleContable"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		dbId int
		crId int
		movs []*models.PreMovAjuste
	)

	movimientos_ = make([]*models.DetalleMovimientoContable, 0)

	if dbId, crId, outputError = parametros.GetParametrosDebitoCredito(); outputError != nil {
		return nil, outputError
	}

	for _, mov := range movimientos {
		mov_ := new(models.PreMovAjuste)
		mov_.Cuenta = mov.CuentaId
		if mov.TerceroId != nil {
			mov_.TerceroId = *mov.TerceroId
		}
		mov_.Descripcion = mov.Descripcion

		if mov.TipoMovimientoId == crId {
			mov_.Credito = mov.Valor
		} else if mov.TipoMovimientoId == dbId {
			mov_.Debito = mov.Valor
		}

		movs = append(movs, mov_)
	}

	if detalleCuentas == nil {
		detalleCuentas = make(map[string]models.CuentaContable)
	}

	for _, mov := range movs {
		mov_ := new(models.DetalleMovimientoContable)
		if cta, ok := detalleCuentas[mov.Cuenta]; !ok {
			if cta_, err := cuentasContables.GetCuentaContable(mov.Cuenta); err != nil {
				return nil, err
			} else {
				if cta_ != nil {
					detalleCuentas[mov.Cuenta] = *cta_
				} else {
					cta_ = new(models.CuentaContable)
					detalleCuentas[mov.Cuenta] = *cta_
				}
				mov_.Cuenta = &models.DetalleCuenta{
					Id:              cta_.Id,
					Codigo:          cta_.Codigo,
					Nombre:          cta_.Nombre,
					RequiereTercero: cta_.RequiereTercero,
				}
			}
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
		movimientos_ = append(movimientos_, mov_)
	}

	return movimientos_, nil

}

func GetFullDetalleContable(consecutivoId int) (trContable models.InfoTransaccionContable, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GetFullDetalleContable - Unhandled Error!", "500")

	transaccion, outputError := movimientosContables.GetTransaccion(consecutivoId, "consecutivo", true)
	if outputError != nil {
		return
	}

	trContable = models.InfoTransaccionContable{
		Concepto: transaccion.Descripcion,
		Fecha:    transaccion.FechaTransaccion,
	}

	if len(transaccion.Movimientos) > 0 {
		trContable.Movimientos, outputError = GetDetalleContable(transaccion.Movimientos, nil)
	}

	return
}
