package cuentasContablesHelper

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

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
	// logs.Debug("urlcrud:", urlcrud)

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

// realiza el asiento contable. totales tiene los valores por clase, tipomvto el tipo de mvto
func AsientoContable(totales map[int]float64, tipomvto string, descripcionMovto string) (response map[string]interface{}, outputError map[string]interface{}) {

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

	//consecutivo del asiento
	year, _, _ := time.Now().Date()
	consec := models.Consecutivo{0, 1, year, 0, "CNTB", true}
	apiCons := "http://" + beego.AppConfig.String("consecutivosService") + "consecutivo"
	idconsecutivo := float64(0)
	if err := request.SendJson(apiCons, "POST", &res, &consec); err == nil {
		resultado, _ := res["Data"].(map[string]interface{})
		idconsecutivo = resultado["Id"].(float64)
	} else {
		outputError = map[string]interface{}{"funcion": "asientoContable -response, err := request.SendJson(apiCons,", "status": "500", "err": err}
		return nil, outputError
	}

	var jsonString []byte
	var err1 error

	//captura el id del movimiento credito
	urlcrud := "http://" + beego.AppConfig.String("parametrosService") + "parametro?query=CodigoAbreviacion:MCD"
	if err := request.GetJson(urlcrud, &resMap); err != nil { // Get parámetro tipo movimiento contable débito
		outputError = map[string]interface{}{"funcion": "asientoContable - if err := request.GetJson(urlcrud, &resMap);", "status": "502", "err": err}
		return nil, outputError
	}

	if jsonString, err1 = json.Marshal(resMap["Data"]); err1 != nil {
		outputError = map[string]interface{}{"funcion": "asientoContable - if jsonString, err1 = json.Marshal(resMap[\"Data\"]);", "status": "500", "err": err1}
		return nil, outputError
	}
	var parametro []models.Parametro
	if err1 = json.Unmarshal(jsonString, &parametro); err1 != nil {
		outputError = map[string]interface{}{"funcion": "asientoContable - if err1 = json.Unmarshal(jsonString, &parametro);", "status": "500", "err": err1}
		return nil, outputError
	}
	//	resMap = make(map[string]interface{})
	parametroTipoDebito = parametro[0]

	//captura el id del movimiento debito
	urlcrud = "http://" + beego.AppConfig.String("parametrosService") + "parametro?query=CodigoAbreviacion:MCC"
	if err1 = request.GetJson(urlcrud, &resMap); err1 != nil { // Get parámetro tipo movimiento contable débito
		outputError = map[string]interface{}{"funcion": "asientoContable - if err1 = request.GetJson(urlcrud, &resMap);", "status": "502", "err": err1}
		return nil, outputError
	}

	if jsonString, err1 = json.Marshal(resMap["Data"]); err1 != nil {
		outputError = map[string]interface{}{"funcion": "asientoContable - if jsonString, err1 = json.Marshal(resMap[\"Data\"]);", "status": "500", "err": err1}
		return nil, outputError
	}
	if err1 = json.Unmarshal(jsonString, &parametro); err1 != nil {
		outputError = map[string]interface{}{"funcion": "asientoContable - if err1 = json.Unmarshal(jsonString, &parametro);", "status": "500", "err": err1}
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
						outputError = map[string]interface{}{"funcion": "asientoContable - if err = json.Unmarshal(jsonString, &tipoComprobanteContable);", "status": "500", "err": err}
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

	transaccion.Activo = true
	transaccion.ConsecutivoId = int(idconsecutivo)
	transaccion.Descripcion = descripcionMovto
	transaccion.Etiquetas = string(jsonString)
	transaccion.FechaTransaccion = time_bogota.Tiempo_bogota()

	for clave, _ := range totales {
		urlcuentas := "http://" + beego.AppConfig.String("catalogoElementosService") + "cuentas_subgrupo/?query=SubgrupoId.Id:" + strconv.Itoa(clave) + ",Activo:true,SubtipoMovimientoId:" + tipomvto
		logs.Info(urlcuentas)
		if respuesta, err := request.GetJsonTest(urlcuentas, &elemento); err == nil && respuesta.StatusCode == 200 {
			for _, element := range elemento {
				if len(element) == 0 {
					outputError = map[string]interface{}{"funcion": "asientoContable - if len(element) == 0 ", "status": "500", "err": err}
					return nil, outputError
				} else {
					nombrecuentadebito := ""
					nombrecuentacredito := ""
					urlcuenta := "http://" + beego.AppConfig.String("cuentasContablesService") + "nodo_cuenta_contable/" + element["CuentaDebitoId"].(string)
					if respuesta, err := request.GetJsonTest(urlcuenta, &respuesta_peticion); err == nil && respuesta.StatusCode == 200 {
						nombrecuentadebito = respuesta_peticion["Body"].(interface{}).(map[string]interface{})["Nombre"].(string)
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

					movimientoDebito.TerceroId = 1
					movimientoDebito.CuentaId = element["CuentaDebitoId"].(string)
					movimientoDebito.NombreCuenta = nombrecuentadebito
					movimientoDebito.TipoMovimientoId = parametroTipoDebito.Id
					movimientoDebito.Valor = totales[clave]
					movimientoDebito.Descripcion = "primer movimiento"
					transaccion.Movimientos = append(transaccion.Movimientos, movimientoDebito)

					movimientoCredito.TerceroId = 1
					movimientoCredito.CuentaId = element["CuentaCreditoId"].(string)
					movimientoCredito.NombreCuenta = nombrecuentacredito
					movimientoCredito.TipoMovimientoId = parametroTipoCredito.Id
					movimientoCredito.Valor = totales[clave]
					movimientoCredito.Descripcion = "segundo movimiento"
					transaccion.Movimientos = append(transaccion.Movimientos, movimientoCredito)
				}
			}
			apiMvtoContables := "http://" + beego.AppConfig.String("movimientosContablesmidService") + "transaccion_movimientos/transaccion_movimientos/"
			logs.Info(fmt.Sprintf("apiMvtoContables: %s - body: %v", apiMvtoContables, transaccion))
			outputError = map[string]interface{}{"funcion": "asientoContable - en desarrollo", "status": "500", "err": err}
			return nil, outputError

			/*	este es el post de la contabilidad if err := request.SendJson(apiMvtoContables, "POST", &res, &mvto); err == nil {
				logs.Info("Termino bien")
			}*/
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
	return res, nil
}
