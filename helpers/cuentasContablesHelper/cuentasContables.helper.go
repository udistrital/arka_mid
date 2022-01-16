package cuentasContablesHelper

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/tercerosHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
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

// GetCuentaContable ...
func GetCuentaContable(cuentaContableId string) (cuentaContable map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetCuentaContable - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	urlcrud := "http://" + beego.AppConfig.String("cuentasContablesService") + "nodo_cuenta_contable/" + cuentaContableId

	var data models.RespuestaAPI2obj
	if resp, err := request.GetJsonTest(urlcrud, &data); err == nil && resp.StatusCode == 200 && data.Code == 200 {
		return data.Body, nil
	} else {
		if err == nil {
			if resp.StatusCode != 200 {
				err = fmt.Errorf("Undesired Status Code: %d", resp.StatusCode)
			} else {
				err = fmt.Errorf("Undesired Status Code: %d - in Body: %d", resp.StatusCode, data.Code)
			}
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetCuentaContable - request.GetJsonTest(urlcrud, &cuentaContable)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}

func creaMovimiento(valor float64, descripcionMovto string, idTercero int, cuentaId string, cuenta *models.CuentaContable, tipo int) (movimientoDebito models.MovimientoTransaccion) {
	var movimiento models.MovimientoTransaccion

	if cuenta.RequiereTercero {
		movimiento.TerceroId = &idTercero
	} else {
		movimiento.TerceroId = nil
	}
	movimiento.CuentaId = cuentaId
	movimiento.NombreCuenta = cuenta.Nombre
	movimiento.TipoMovimientoId = tipo
	movimiento.Valor = valor
	movimiento.Descripcion = descripcionMovto
	return movimiento
}

func GetInfoSubgrupo(subgrupoId int) (detalleSubgrupo map[string]interface{}, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetInfoSubgrupo - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()
	if subgrupoId <= 0 {
		err := fmt.Errorf("subgrupoId MUST be > 0 - Got: %d", subgrupoId)
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetInfoSubgrupo - subgrupoId <= 0",
			"err":     err,
			"status":  "400",
		}
		panic(outputError)
	}

	var detalles []map[string]interface{}

	urlcrud := "http://" + beego.AppConfig.String("catalogoElementosService") + "subgrupo?limit=-1"
	urlcrud += "&query=Activo:True,Id:" + strconv.Itoa(int(subgrupoId))
	if response, err := request.GetJsonTest(urlcrud, &detalles); err == nil && response.StatusCode == 200 { // (2) error servicio caido
		// fmt.Println(cuentasSubgrupo[0])
		if detalles[0]["Id"].(float64) != 0 {
			fmt.Println("el detalle", detalles[0])
			return detalles[0], nil
		} else {
			err = fmt.Errorf("Cuenta no existe")
			outputError = map[string]interface{}{
				"funcion": "GetInfoSubgrupo - request.GetJsonTest(urlcrud, &cuentasSubgrupo)",
				"err":     err,
				"status:": "502",
			}
			return nil, outputError
		}
	} else {
		if err == nil {
			err = fmt.Errorf("Undesired Status Code: %d", response.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetInfoSubgrupo - request.GetJsonTest(urlcrud, &cuentasSubgrupo)",
			"err":     err,
			"status:": "502",
		}
		return nil, outputError
	}
}

// AsientoContable realiza el asiento contable. totales tiene los valores por clase, tipomvto el tipo de mvto
func AsientoContable(totales map[int]float64, tipomvto string, descripcionMovto string, descripcionAsiento string, idTercero int, submit bool) (response map[string]interface{}, outputError map[string]interface{}) {

	funcion := "AsientoContable"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		res    map[string]interface{}
		resMap map[string]interface{}
		//	elemento                []map[string]interface{}
		transaccion models.TransaccionMovimientos
		//	respuesta_peticion      map[string]interface{}
		parametroTipoDebito     models.Parametro
		parametroTipoCredito    models.Parametro
		tipoComprobanteContable models.TipoComprobanteContable
		cuentasSubgrupo         []*models.CuentaSubgrupo

		ctaCr *models.CuentaContable
		ctaDb *models.CuentaContable
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
	var jsonString []byte
	var err1 error

	//captura el id del movimiento credito
	urlcrud := "http://" + beego.AppConfig.String("parametrosService") + "parametro?query=CodigoAbreviacion:MCD"
	if err := request.GetJson(urlcrud, &resMap); err != nil { // Get parámetro tipo movimiento contable débito
		outputError = map[string]interface{}{"funcion": "asientoContable - request.GetJson(urlcrud, &resMap)", "status": "502", "err": err}
		return nil, outputError
	}

	if jsonString, err1 = json.Marshal(resMap["Data"]); err1 != nil {
		outputError = map[string]interface{}{"funcion": "asientoContable - json.Marshal(resMap[\"Data\"])", "status": "500", "err": err1}
		return nil, outputError
	}
	var parametro []models.Parametro
	if err1 = json.Unmarshal(jsonString, &parametro); err1 != nil {
		outputError = map[string]interface{}{"funcion": "asientoContable - json.Unmarshal(jsonString, &parametro)", "status": "500", "err": err1}
		return nil, outputError
	}
	parametroTipoDebito = parametro[0]

	//captura el id del movimiento debito
	urlcrud = "http://" + beego.AppConfig.String("parametrosService") + "parametro?query=CodigoAbreviacion:MCC"
	if err1 = request.GetJson(urlcrud, &resMap); err1 != nil { // Get parámetro tipo movimiento contable débito
		outputError = map[string]interface{}{"funcion": "asientoContable - request.GetJson(urlcrud, &resMap)", "status": "502", "err": err1}
		return nil, outputError
	}

	if jsonString, err1 = json.Marshal(resMap["Data"]); err1 != nil {
		outputError = map[string]interface{}{"funcion": "asientoContable - json.Marshal(resMap[\"Data\"])", "status": "500", "err": err1}
		return nil, outputError
	}
	if err1 = json.Unmarshal(jsonString, &parametro); err1 != nil {
		outputError = map[string]interface{}{"funcion": "asientoContable - json.Unmarshal(jsonString, &parametro)", "status": "500", "err": err1}
		return nil, outputError
	}

	//	resMap = make(map[string]interface{})
	parametroTipoCredito = parametro[0]

	urlcrud = "http://" + beego.AppConfig.String("cuentasContablesService") + "tipo_comprobante"
	if err := request.GetJson(urlcrud, &resMap); err == nil { // Para obtener código del tipo de comprobante
		for _, sliceTipoComprobante := range resMap["Body"].([]interface{}) {

			if valor, ok := sliceTipoComprobante.(map[string]interface{})["TipoDocumento"]; ok && valor == "E" {
				if jsonString, err = json.Marshal(sliceTipoComprobante); err == nil {
					if err = json.Unmarshal(jsonString, &tipoComprobanteContable); err == nil {
						resMap = make(map[string]interface{})
					} else {
						logs.Error(err)
						outputError = map[string]interface{}{"funcion": "asientoContable -json.Unmarshal(jsonString, &tipoComprobanteContable)", "status": "500", "err": err}
						return nil, outputError
					}
				} else {
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "asientoContable - if jsonString, err = json.Marshal(sliceTipoComprobante);", "status": "500", "err": err}
					return nil, outputError
				}
			}
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "asientoContable - if err := request.GetJson(urlcrud, &resMap);", "status": "502", "err": err}
		return nil, outputError
	}

	etiquetas := make(map[string]interface{})
	etiquetas["TipoComprobanteId"] = tipoComprobanteContable.Codigo

	if jsonString, err1 = json.Marshal(etiquetas); err1 != nil {
		outputError = map[string]interface{}{"funcion": "asientoContable - if jsonString, err1 = json.Marshal(etiquetas);", "status": "500", "err": err1}
		return nil, outputError
	}

	transaccion.ConsecutivoId = consecutivoId
	transaccion.Activo = true
	transaccion.Descripcion = descripcionMovto
	transaccion.Etiquetas = string(jsonString)
	transaccion.FechaTransaccion = time_bogota.Tiempo_bogota()

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
	for clave, _ := range totales {
		if idx := FindInArray(cuentasSubgrupo, clave); idx > -1 {
			if ctaCr_, err := GetCuentaContable(cuentasSubgrupo[idx].CuentaCreditoId); err != nil {
				return nil, err
			} else {
				if err := formatdata.FillStruct(ctaCr_, &ctaCr); err != nil {
					logs.Error(err)
					eval := " - formatdata.FillStruct(ctaCr_, &ctaCr)"
					return nil, errorctrl.Error(funcion+eval, err, "500")
				}
			}

			if ctaDb_, err := GetCuentaContable(cuentasSubgrupo[idx].CuentaDebitoId); err != nil {
				return nil, err
			} else {
				if err := formatdata.FillStruct(ctaDb_, &ctaDb); err != nil {
					logs.Error(err)
					eval := " - formatdata.FillStruct(ctaDb_, &ctaDb)"
					return nil, errorctrl.Error(funcion+eval, err, "500")
				}
			}

			infoCuentas[clave] = new(InfoCuentasSubgrupos)
			infoCuentas[clave].CuentaCredito = ctaCr
			infoCuentas[clave].CuentaDebito = ctaDb
			//			var movimientoDebito models.MovimientoTransaccion
			//		var movimientoCredito models.MovimientoTransaccion
			movimientoCredito := creaMovimiento(totales[clave], descripcionAsiento, idTercero, cuentasSubgrupo[idx].CuentaCreditoId, ctaCr, parametroTipoCredito.Id)
			movimientoDebito := creaMovimiento(totales[clave], descripcionAsiento, idTercero, cuentasSubgrupo[idx].CuentaDebitoId, ctaDb, parametroTipoDebito.Id)
			transaccion.Movimientos = append(transaccion.Movimientos, movimientoDebito)
			transaccion.Movimientos = append(transaccion.Movimientos, movimientoCredito)

		} else {
			subgrupo, error := GetInfoSubgrupo(clave)
			if error == nil {
				res["errorTransaccion"] = fmt.Sprintf("Debe parametrizar las cuentas del subgrupo %v", subgrupo["Nombre"].(string))
				return res, nil
			} else {
				return nil, error
			}
		}
	}

	if submit {
		if resp, err := PostTrContable(&transaccion); err != nil || !resp.Success {
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

// PostTrContable post controlador transaccion_movimientos/transaccion_movimientos/ del api movimientos_contables_mid
func PostTrContable(tr *models.TransaccionMovimientos) (resp *models.RespuestaAPI1Str, outputError map[string]interface{}) {

	funcion := "PostTrContable"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error", "500")

	urlcrud := "http://" + beego.AppConfig.String("movimientosContablesmidService") + "transaccion_movimientos/transaccion_movimientos/"
	if err := request.SendJson(urlcrud, "POST", &resp, &tr); err != nil {
		eval := ` - request.SendJson(urlcrud, "POST", &novedadR, &novedad)`
		return nil, errorctrl.Error(funcion+eval, err, "502")
	} else if strings.Contains(resp.Data, "invalid character") {
		logs.Error(resp.Data)
		resp, outputError = PostTrContable(tr)
	}

	return resp, nil
}
