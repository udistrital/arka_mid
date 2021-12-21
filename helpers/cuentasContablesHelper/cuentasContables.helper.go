package cuentasContablesHelper

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/tercerosHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

const ID_SALIDA_PRUEBAS = "16"
const ID_SALIDA_CONSUMO_PRUEBAS = "22"

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

// AsientoContable realiza el asiento contable. totales tiene los valores por clase, tipomvto el tipo de mvto
func AsientoContable(totales map[int]float64, tipomvto string, descripcionMovto string, descripcionAsiento string, idTercero int) (response map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "AsientoContable - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var (
		res                     map[string]interface{}
		resMap                  map[string]interface{}
		elemento                []map[string]interface{}
		transaccion             models.TransaccionMovimientos
		respuesta_peticion      map[string]interface{}
		parametroTipoDebito     models.Parametro
		parametroTipoCredito    models.Parametro
		tipoComprobanteContable models.TipoComprobanteContable
	)

	if tipomvto == ID_SALIDA_CONSUMO_PRUEBAS {
		tipomvto = ID_SALIDA_PRUEBAS
	}
	idconsecutivo := ""
	if idconsecutivo1, err := utilsHelper.GetConsecutivo("%05.0f", 1, "CNTB"); err != nil {

		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "RegistrarEntrada - utilsHelper.GetConsecutivo()",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	} else {
		idconsecutivo = idconsecutivo1
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

	intcons, _ := strconv.Atoi(idconsecutivo)
	transaccion.ConsecutivoId = intcons
	transaccion.Activo = true
	transaccion.Descripcion = descripcionMovto
	transaccion.Etiquetas = string(jsonString)
	transaccion.FechaTransaccion = time_bogota.Tiempo_bogota()

	tercerodebito := true
	tercerocredito := true
	for clave, _ := range totales {
		urlcuentas := "http://" + beego.AppConfig.String("catalogoElementosService") + "cuentas_subgrupo/?sortby=Id&order=desc&query=SubgrupoId.Id:" + strconv.Itoa(clave) + ",Activo:true,SubtipoMovimientoId:" + tipomvto
		logs.Debug("******* La url de las cuentas", urlcuentas)
		if respuesta, err := request.GetJsonTest(urlcuentas, &elemento); err == nil && respuesta.StatusCode == 200 {
			//			for _, element := range elemento { //deberia existir un solo para de cuentas para cada tipo de movimiento, pero esto hay que discutirlo
			element := elemento[0]
			if len(element) == 0 {
				outputError = map[string]interface{}{"funcion": "asientoContable - if len(element) == 0 ", "status": "500", "err": err}
				return nil, outputError
			} else {
				nombrecuentadebito := ""
				nombrecuentacredito := ""
				urlcuenta := "http://" + beego.AppConfig.String("cuentasContablesService") + "nodo_cuenta_contable/" + element["CuentaDebitoId"].(string)
				if respuesta, err := request.GetJsonTest(urlcuenta, &respuesta_peticion); err == nil && respuesta.StatusCode == 200 {
					nombrecuentadebito = respuesta_peticion["Body"].(interface{}).(map[string]interface{})["Nombre"].(string)
					tercerodebito = respuesta_peticion["Body"].(interface{}).(map[string]interface{})["RequiereTercero"].(bool)
				} else {
					if err == nil {
						err = fmt.Errorf("Undesired Status Code: %d", respuesta.StatusCode)
					}
					logs.Error(err)
					outputError = map[string]interface{}{
						"funcion": "GetCuentaContable -  if respuesta, err := request.GetJsonTest(urlcuenta, &respuesta_peticion);)",
						"err":     err,
						"status":  "502",
					}
					return nil, outputError
				}

				urlcuenta = "http://" + beego.AppConfig.String("cuentasContablesService") + "nodo_cuenta_contable/" + element["CuentaCreditoId"].(string)
				if respuesta, err := request.GetJsonTest(urlcuenta, &respuesta_peticion); err == nil && respuesta.StatusCode == 200 {
					nombrecuentacredito = respuesta_peticion["Body"].(interface{}).(map[string]interface{})["Nombre"].(string)
					tercerocredito = respuesta_peticion["Body"].(interface{}).(map[string]interface{})["RequiereTercero"].(bool)
				} else {
					if err == nil {
						err = fmt.Errorf("Undesired Status Code: %d", respuesta.StatusCode)
					}
					logs.Error(err)
					outputError = map[string]interface{}{
						"funcion": "GetCuentaContable -  if respuesta, err := request.GetJsonTest(urlcuenta, &respuesta_peticion);)",
						"err":     err,
						"status":  "502",
					}
					return nil, outputError
				}

				var movimientoDebito models.MovimientoTransaccion
				var movimientoCredito models.MovimientoTransaccion

				if tercerodebito {
					movimientoDebito.TerceroId = idTercero
				} else {
					movimientoDebito.TerceroId = 0
				}
				movimientoDebito.CuentaId = element["CuentaDebitoId"].(string)
				movimientoDebito.NombreCuenta = nombrecuentadebito
				movimientoDebito.TipoMovimientoId = parametroTipoDebito.Id
				movimientoDebito.Valor = totales[clave]
				movimientoDebito.Descripcion = descripcionAsiento
				transaccion.Movimientos = append(transaccion.Movimientos, movimientoDebito)

				if tercerocredito {
					movimientoCredito.TerceroId = idTercero
				} else {
					movimientoCredito.TerceroId = 0
				}
				movimientoCredito.CuentaId = element["CuentaCreditoId"].(string)
				movimientoCredito.NombreCuenta = nombrecuentacredito
				movimientoCredito.TipoMovimientoId = parametroTipoCredito.Id
				movimientoCredito.Valor = totales[clave]
				movimientoCredito.Descripcion = descripcionAsiento
				transaccion.Movimientos = append(transaccion.Movimientos, movimientoCredito)
			}
			//		}
		} else {

			if err == nil {
				err = fmt.Errorf("Undesired Status Code: %d", respuesta.StatusCode)
			}
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "AsientoContable - if respuesta, err := request.GetJsonTest(urlcuentas, &elemento);",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
	}

	apiMvtoContables := "http://" + beego.AppConfig.String("movimientosContablesmidService") + "transaccion_movimientos/transaccion_movimientos/"
	if err := request.SendJson(apiMvtoContables, "POST", &res, &transaccion); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "AsientoContable -if err := request.SendJson(apiMvtoContables, \"POST\", &res, &transaccion);",
			"err":     err,
			"status":  "502",
		}
	} else {
		res["resultadoTransaccion"] = transaccion
		if tercero, err := tercerosHelper.GetTerceroById(idTercero); err == nil {
			logs.Debug(tercero)
			res["tercero"] = tercero
		} else {
			return nil, err
		}
		return res, nil
	}
	return res, nil
}
