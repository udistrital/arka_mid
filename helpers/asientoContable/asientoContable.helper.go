package asientoContable

import (
	"net/url"
	"time"

	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	crudTerceros "github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

func CreaMovimiento(valor float64, descripcionMovto string, idTercero int, cuenta *models.CuentaContable, tipo int) (movimiento *models.MovimientoTransaccion) {
	movimiento = new(models.MovimientoTransaccion)

	if cuenta.RequiereTercero {
		movimiento.TerceroId = &idTercero
	} else {
		movimiento.TerceroId = nil
	}

	movimiento.CuentaId = cuenta.Id
	movimiento.NombreCuenta = cuenta.Nombre
	movimiento.TipoMovimientoId = tipo
	movimiento.Valor = valor
	movimiento.Descripcion = descripcionMovto
	movimiento.Activo = true

	return movimiento
}

// AsientoContable realiza el asiento contable. totales tiene los valores por clase, tipomvto el tipo de mvto
func AsientoContable(totales map[int]float64, comprobante, tipomvto, descripcionMovto, descripcionAsiento string, idTercero, consecutivoId int, submit bool) (response map[string]interface{}, outputError map[string]interface{}) {

	funcion := "AsientoContable"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		res                  map[string]interface{}
		transaccion          models.TransaccionMovimientos
		parametroTipoDebito  int
		parametroTipoCredito int
		cuentasSubgrupo      []*models.CuentasSubgrupo
		comprobanteID        string
	)

	res = make(map[string]interface{})
	res["errorTransaccion"] = ""

	if db_, cr_, err := parametros.GetParametrosDebitoCredito(); err != nil {
		return nil, err
	} else {
		parametroTipoDebito = db_
		parametroTipoCredito = cr_
	}

	if comprobante != "" {
		if err := cuentasContables.GetComprobante(comprobante, &comprobanteID); err != nil {
			return nil, err
		}
	}

	if comprobanteID != "" {
		etiquetas := *new(models.Etiquetas)
		etiquetas.ComprobanteId = comprobanteID
		outputError = utilsHelper.Marshal(etiquetas, &transaccion.Etiquetas)
		if outputError != nil {
			return
		}
	} else {
		transaccion.Etiquetas = ""
	}

	transaccion.ConsecutivoId = consecutivoId
	transaccion.Activo = true
	transaccion.Descripcion = descripcionAsiento
	transaccion.FechaTransaccion = time.Now()

	idsSubgrupos := make([]int, len(totales))

	i := 0
	for k := range totales {
		idsSubgrupos[i] = k
		i++
	}

	query := "limit=-1&fields=CuentaDebitoId,CuentaCreditoId,SubgrupoId&sortby=Id&order=desc&"
	query += "query=SubtipoMovimientoId:" + tipomvto + ",Activo:true,SubgrupoId__Id__in:"
	query += url.QueryEscape(utilsHelper.ArrayToString(idsSubgrupos, "|"))
	if elementos_, err := catalogoElementos.GetAllCuentasSubgrupo(query); err != nil {
		return nil, err
	} else {
		cuentasSubgrupo = elementos_
	}

	infoCuentas := make(map[string]*models.CuentaContable)
	for id := range totales {
		if idx := FindInArray(cuentasSubgrupo, id); idx > -1 {

			if ctaCr_, err := cuentasContables.GetCuentaContable(cuentasSubgrupo[idx].CuentaCreditoId); err != nil {
				return nil, err
			} else if ctaCr_ == nil {
				subgrupo, err := catalogoElementos.GetSubgrupoById(id)
				if err != nil {
					return nil, err
				} else {
					res["errorTransaccion"] = "Debe parametrizar las cuentas del subgrupo " + subgrupo.Codigo + " " + subgrupo.Nombre
					return res, nil
				}
			} else {
				infoCuentas[cuentasSubgrupo[idx].CuentaCreditoId] = ctaCr_
			}

			if ctaDb_, err := cuentasContables.GetCuentaContable(cuentasSubgrupo[idx].CuentaDebitoId); err != nil {
				return nil, err
			} else if ctaDb_ == nil {
				subgrupo, err := catalogoElementos.GetSubgrupoById(id)
				if err != nil {
					return nil, err
				} else {
					res["errorTransaccion"] = "Debe parametrizar las cuentas del subgrupo " + subgrupo.Codigo + " " + subgrupo.Nombre
					return res, nil
				}
			} else {
				infoCuentas[cuentasSubgrupo[idx].CuentaDebitoId] = ctaDb_
			}

			movimientoCredito := CreaMovimiento(totales[id], descripcionMovto, idTercero, infoCuentas[cuentasSubgrupo[idx].CuentaCreditoId], parametroTipoCredito)
			movimientoDebito := CreaMovimiento(totales[id], descripcionMovto, idTercero, infoCuentas[cuentasSubgrupo[idx].CuentaDebitoId], parametroTipoDebito)
			transaccion.Movimientos = append(transaccion.Movimientos, movimientoDebito)
			transaccion.Movimientos = append(transaccion.Movimientos, movimientoCredito)

		} else {
			subgrupo, err := catalogoElementos.GetSubgrupoById(id)
			if err != nil {
				return nil, err
			} else {
				res["errorTransaccion"] = "Debe parametrizar las cuentas del subgrupo " + subgrupo.Codigo + " " + subgrupo.Nombre
				return res, nil
			}
		}
	}

	if submit {
		if tr, err := movimientosContables.PostTrContable(&transaccion); err != nil {
			return nil, err
		} else {
			if tercero, err := crudTerceros.GetNombreTerceroById(idTercero); err != nil {
				return nil, err
			} else {
				res["resultadoTransaccion"] = fillDetalle(infoCuentas, tr, tercero.Numero)
			}
		}
	} else {
		if tercero, err := crudTerceros.GetNombreTerceroById(idTercero); err != nil {
			return nil, err
		} else {
			res["simulacro"] = fillDetalle(infoCuentas, &transaccion, tercero.Numero)
		}
	}

	return res, nil
}

func fillDetalle(cuentas map[string]*models.CuentaContable, transaccion *models.TransaccionMovimientos, tercero string) (detalle *models.DetalleTrContable) {

	detalle = new(models.DetalleTrContable)
	movimientos := make([]*models.DetalleMovimientoTransaccion, 0)
	for _, mov := range transaccion.Movimientos {
		movimiento_ := new(models.DetalleMovimientoTransaccion)
		movimiento_.CuentaId = cuentas[mov.CuentaId]
		movimiento_.TipoMovimientoId = mov.TipoMovimientoId
		movimiento_.Valor = mov.Valor
		movimiento_.Descripcion = mov.Descripcion

		if movimiento_.CuentaId.RequiereTercero {
			movimiento_.TerceroId = tercero
		}

		movimientos = append(movimientos, movimiento_)
	}

	detalle.ConsecutivoId = transaccion.ConsecutivoId
	detalle.Descripcion = transaccion.Descripcion
	detalle.Etiquetas = transaccion.Etiquetas
	detalle.FechaTransaccion = transaccion.FechaTransaccion
	detalle.Movimientos = movimientos

	return detalle
}

// findIdInArray Retorna la posicion en que se encuentra el id específicado
func FindInArray(cuentasSg []*models.CuentasSubgrupo, subgrupoId int) (i int) {
	for i, cuentaSg := range cuentasSg {
		if int(cuentaSg.SubgrupoId.Id) == subgrupoId {
			return i
		}
	}
	return -1
}
