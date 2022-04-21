package asientoContable

import (
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
)

// GetDetalleContable Consulta los detalles de una transacción contable para ser mostrada en el cliente
func GetDetalleContable(movimientos []*models.MovimientoTransaccion) (movimientos_ []*models.DetalleMovimientoContable, outputError map[string]interface{}) {

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
			mov_.Debito = 0
		} else if mov.TipoMovimientoId == dbId {
			mov_.Debito = mov.Valor
			mov_.Credito = 0
		}

		movs = append(movs, mov_)
	}

	for _, mov := range movs {
		mov_ := new(models.DetalleMovimientoContable)
		var cta *models.DetalleCuenta

		if ctaCr_, err := cuentasContables.GetCuentaContable(mov.Cuenta); err != nil {
			return nil, err
		} else {
			if err := formatdata.FillStruct(ctaCr_, &cta); err != nil {
				logs.Error(err)
				eval := " - formatdata.FillStruct(ctaCr_, &ctaCr)"
				return nil, errorctrl.Error(funcion+eval, err, "500")
			}
			mov_.Cuenta = cta
		}

		if mov.TerceroId > 0 {
			if tercero_, err := terceros.GetTerceroById(mov.TerceroId); err != nil {
				return nil, err
			} else {
				mov_.TerceroId = tercero_
			}
		}
		mov_.Credito = mov.Credito
		mov_.Debito = mov.Debito
		mov_.Descripcion = mov.Descripcion
		movimientos_ = append(movimientos_, mov_)
	}

	return movimientos_, nil

}
