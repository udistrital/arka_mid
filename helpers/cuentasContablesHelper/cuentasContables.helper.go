package cuentasContablesHelper

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/movimientosContablesMidHelper"
	"github.com/udistrital/arka_mid/helpers/parametrosHelper"
	"github.com/udistrital/arka_mid/helpers/tercerosHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

const ID_SALIDA_PRUEBAS = "16"
const ID_SALIDA_CONSUMO_PRUEBAS = "22"

type InfoCuentasSubgrupos struct {
	CuentaDebito  *models.CuentaContable
	CuentaCredito *models.CuentaContable
}

// GetAllCuentasSubgrupo query controlador cuentas_subgrupo del api catalogo_elementos_crud
func GetAllCuentasSubgrupo(query string) (elementos []*models.CuentaSubgrupo, outputError map[string]interface{}) {

	funcion := "GetAllCuentasSubgrupo"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("catalogoElementosService") + "cuentas_subgrupo?" + query
	if err := request.GetJson(urlcrud, &elementos); err != nil {
		logs.Error(err)
		eval := " - request.GetJson(urlcrud, &elementos)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	return elementos, nil
}

func creaMovimiento(valor float64, descripcionMovto string, idTercero int, cuenta *models.CuentaContable, tipo int) (movimiento *models.MovimientoTransaccion) {
	movimiento = new(models.MovimientoTransaccion)

	if cuenta.RequiereTercero {
		movimiento.TerceroId = &idTercero
	} else {
		movimiento.TerceroId = nil
	}

	movimiento.CuentaId = cuenta.Codigo
	movimiento.NombreCuenta = cuenta.Nombre
	movimiento.TipoMovimientoId = tipo
	movimiento.Valor = valor
	movimiento.Descripcion = descripcionMovto

	return movimiento
}

// GetSubgrupoById Consulta controlador subgrupo/{id} del api catalogo_elementos_crud
func GetSubgrupoById(id int) (subgrupo *models.Subgrupo, outputError map[string]interface{}) {

	funcion := "GetSubgrupoById"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("catalogoElementosService") + "subgrupo/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &subgrupo); err != nil {
		logs.Error(err)
		eval := " - request.GetJson(urlcrud, &subgrupo)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	return subgrupo, nil
}

// AsientoContable realiza el asiento contable. totales tiene los valores por clase, tipomvto el tipo de mvto
func AsientoContable(totales map[int]float64, tipomvto string, descripcionMovto string, descripcionAsiento string, idTercero int, submit bool) (response map[string]interface{}, outputError map[string]interface{}) {

	funcion := "AsientoContable"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		res map[string]interface{}
		//	elemento                []map[string]interface{}
		transaccion models.TransaccionMovimientos
		//	respuesta_peticion      map[string]interface{}
		parametroTipoDebito  int
		parametroTipoCredito int
		cuentasSubgrupo      []*models.CuentaSubgrupo
	)

	res = make(map[string]interface{})
	res["errorTransaccion"] = ""
	if tipomvto == ID_SALIDA_CONSUMO_PRUEBAS {
		tipomvto = ID_SALIDA_PRUEBAS
	}

	consecutivoId := 0
	if submit {
		if _, consecutivoId_, err := utilsHelper.GetConsecutivo("%05.0f", 1, "CNTB"); err != nil {
			return nil, outputError
		} else {
			consecutivoId = consecutivoId_
		}
	}

	if db_, cr_, err := parametrosHelper.GetParametrosDebitoCredito(); err != nil {
		return nil, err
	} else {
		parametroTipoDebito = db_
		parametroTipoCredito = cr_
	}

	etiquetas := make(map[string]interface{})

	if tipoComprobante_, err := GetTipoComprobante("E"); err != nil {
		return nil, err
	} else if tipoComprobante_ != nil {
		etiquetas["TipoComprobanteId"] = tipoComprobante_.Codigo
		if jsonData, err := json.Marshal(etiquetas); err != nil {
			logs.Error(err)
			eval := " - json.Marshal(etiquetas)"
			return nil, errorctrl.Error(funcion+eval, err, "500")
		} else {
			transaccion.Etiquetas = string(jsonData[:])
		}
	} else {
		transaccion.Etiquetas = ""
	}

	transaccion.ConsecutivoId = consecutivoId
	transaccion.Activo = true
	transaccion.Descripcion = descripcionMovto
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
	if elementos_, err := GetAllCuentasSubgrupo(query); err != nil {
		return nil, err
	} else {
		cuentasSubgrupo = elementos_
	}

	infoCuentas := make(map[int]*InfoCuentasSubgrupos)
	for id := range totales {
		if idx := FindInArray(cuentasSubgrupo, id); idx > -1 {

			infoCuentas[id] = new(InfoCuentasSubgrupos)
			if ctaCr_, err := GetCuentaContable(cuentasSubgrupo[idx].CuentaCreditoId); err != nil {
				return nil, err
			} else {
				infoCuentas[id].CuentaCredito = ctaCr_
			}

			if ctaDb_, err := GetCuentaContable(cuentasSubgrupo[idx].CuentaDebitoId); err != nil {
				return nil, err
			} else {
				infoCuentas[id].CuentaDebito = ctaDb_
			}

			movimientoCredito := creaMovimiento(totales[id], descripcionAsiento, idTercero, infoCuentas[id].CuentaCredito, parametroTipoCredito)
			movimientoDebito := creaMovimiento(totales[id], descripcionAsiento, idTercero, infoCuentas[id].CuentaDebito, parametroTipoDebito)
			transaccion.Movimientos = append(transaccion.Movimientos, movimientoDebito)
			transaccion.Movimientos = append(transaccion.Movimientos, movimientoCredito)

		} else {
			subgrupo, err := GetSubgrupoById(id)
			if err != nil {
				return nil, err
			} else {
				res["errorTransaccion"] = fmt.Sprintf("Debe parametrizar las cuentas del subgrupo ") + subgrupo.Nombre
				return res, nil
			}
		}
	}

	if submit {
		if resp, err := movimientosContablesMidHelper.PostTrContable(&transaccion); err != nil || !resp.Success {
			if err == nil {
				eval := " - PostTrContable(&transaccion)"
				return nil, errorctrl.Error(funcion+eval, resp.Data, "502")
			}
			return nil, err
		} else {
			res["resultadoTransaccion"] = transaccion
			if tercero, err := tercerosHelper.GetNombreTerceroById(strconv.Itoa(idTercero)); err != nil {
				return nil, err
			} else {
				res["tercero"] = tercero
			}
			return res, nil
		}
	} else {
		if tercero, err := tercerosHelper.GetNombreTerceroById(strconv.Itoa(idTercero)); err != nil {
			return nil, err
		} else {
			res["tercero"] = tercero
		}
		res["simulacro"] = transaccion
		return res, nil
	}
}

// findIdInArray Retorna la posicion en que se encuentra el id específicado
func FindInArray(cuentasSg []*models.CuentaSubgrupo, subgrupoId int) (i int) {
	for i, cuentaSg := range cuentasSg {
		if int(cuentaSg.SubgrupoId.Id) == subgrupoId {
			return i
		}
	}
	return -1
}
