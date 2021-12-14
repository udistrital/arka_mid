package entradaHelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/cuentasContablesHelper"
	"github.com/udistrital/arka_mid/helpers/salidaHelper"
	"github.com/udistrital/arka_mid/helpers/tercerosHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

// RegistrarEntrada Crea registro de entrada en estado en trámite
func RegistrarEntrada(data models.Movimiento) (result map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "RegistrarEntrada - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var (
		urlcrud             string
		res                 map[string]interface{}
		actaRecibido        models.TransaccionActaRecibido
		resEstadoMovimiento []models.EstadoMovimiento
	)
	resultado := make(map[string]interface{})

	detalleJSON := map[string]interface{}{}
	if err := json.Unmarshal([]byte(data.Detalle), &detalleJSON); err != nil {
		panic(err.Error())
	}

	ctxConsecutivo, _ := beego.AppConfig.Int("contxtEntradaCons")
	if consecutivo, err := utilsHelper.GetConsecutivo("%05.0f", ctxConsecutivo, "Registro Entrada Arka"); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "RegistrarEntrada - utilsHelper.GetConsecutivo(\"%05.0f\", ctxConsecutivo, \"Registro Entrada Arka\")",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	} else {
		consecutivo = utilsHelper.FormatConsecutivo(getTipoComprobanteEntradas()+"-", consecutivo, fmt.Sprintf("%s%04d", "-", time.Now().Year()))
		detalleJSON["consecutivo"] = consecutivo
		resultado["Consecutivo"] = detalleJSON["consecutivo"]
	}

	if jsonData, err := json.Marshal(detalleJSON); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "RegistrarEntrada - json.Marshal(detalleJSON)",
			"err":     err,
			"status":  "500",
		}
		return nil, outputError
	} else {
		data.Detalle = string(jsonData[:])
	}

	// Consulta el acta
	actaRecibidoId := int(detalleJSON["acta_recibido_id"].(float64))
	urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/" + strconv.Itoa(int(actaRecibidoId)) + "?elementos=false"
	if err := request.GetJson(urlcrud, &actaRecibido); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "RegistrarEntrada - request.GetJson(urlcrud, &actaRecibido)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	// Crea registro en api movimientos_arka_crud
	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "estado_movimiento?query=Nombre:Entrada%20En%20Trámite"
	if err := request.GetJson(urlcrud, &resEstadoMovimiento); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "RegistrarEntrada - request.GetJson(urlcrud, &resEstadoMovimiento)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	} else if len(resEstadoMovimiento) == 0 {
		err = errors.New("len(resEstadoMovimiento) == 0")
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "RegistrarEntrada - request.GetJson(urlcrud, &resEstadoMovimiento)",
			"err":     err,
			"status":  "404",
		}
		return nil, outputError
	}
	data.EstadoMovimientoId.Id = resEstadoMovimiento[0].Id

	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento"
	if err := request.SendJson(urlcrud, "POST", &res, &data); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "RegistrarEntrada - request.SendJson(urlcrud, \"POST\", &res, &data)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	resultado["MovimientoId"] = res["Id"]

	// Crea registro en table soporte_movimiento si es necesario
	if data.SoporteMovimientoId != 0 {
		urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "soporte_movimiento"

		idEntrada := int(res["Id"].(float64))
		movimientoEntrada := models.Movimiento{Id: idEntrada}
		soporteMovimiento := models.SoporteMovimiento{
			DocumentoId:  data.SoporteMovimientoId,
			Activo:       true,
			MovimientoId: &movimientoEntrada,
		}

		if err := request.SendJson(urlcrud, "POST", &res, &soporteMovimiento); err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "RegistrarEntrada - request.SendJson(urlcrud, \"POST\", &resS, &soporteMovimiento)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
	}

	if elementos, err := asignarPlacaActa(actaRecibidoId); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "RegistrarEntrada - asignarPlacaActa(actaRecibidoId)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	} else {
		actaRecibido.Elementos = elementos
	}

	// Actualiza el estado del acta
	urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/" + strconv.Itoa(int(actaRecibidoId))
	actaRecibido.UltimoEstado.EstadoActaId.Id = 6
	actaRecibido.UltimoEstado.Id = 0

	if err := request.SendJson(urlcrud, "PUT", &res, &actaRecibido); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "RegistrarEntrada - request.SendJson(urlcrud, \"PUT\", &res, &actaRecibido)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	return resultado, nil

}

// AprobarEntrada Actualiza una entrada a estado aprobada y hace los respectivos registros en kronos y transacciones contables
func AprobarEntrada(entradaId int) (result map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "AprobarEntrada - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var (
		urlcrud                 string
		res                     map[string]interface{}
		movArka                 []models.Movimiento
		resEstadoMovimiento     []models.EstadoMovimiento
		detalleMovimiento       map[string]interface{}
		transaccionActaRecibido models.TransaccionActaRecibido
	)
	resultado := make(map[string]interface{})

	// Se cambia el estado del movimiento en movimientos_arka_crud
	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento?query=Id:" + strconv.Itoa(int(entradaId))
	if err := request.GetJson(urlcrud, &movArka); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "AprobarEntrada - request.GetJson(urlcrud, &movArka)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "estado_movimiento?query=Nombre:Entrada%20Aprobada"
	if err := request.GetJson(urlcrud, &resEstadoMovimiento); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "AprobarEntrada - request.GetJson(urlcrud, &resEstadoMovimiento)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	} else if len(resEstadoMovimiento) == 0 {
		err = errors.New("len(resEstadoMovimiento) == 0")
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "AprobarEntrada - request.GetJson(urlcrud, &resEstadoMovimiento)",
			"err":     err,
			"status":  "404",
		}
		return nil, outputError
	}

	urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "tipo_movimiento?query=Nombre:" + movArka[0].FormatoTipoMovimientoId.Nombre
	urlcrud = strings.ReplaceAll(urlcrud, " ", "%20")
	if err := request.GetJson(urlcrud, &res); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "AprobarEntrada - request.GetJson(urlcrud, &res)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	} else if reflect.TypeOf(res["Body"]).Kind() != reflect.Slice {
		err = errors.New("no se encuentra tipo_movimiento en api movimientos_crud")
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "AprobarEntrada - reflect.TypeOf(res[\"Body\"]).Kind() != reflect.Slice",
			"err":     err,
			"status":  "404",
		}
		return nil, outputError
	}
	tipomvto := strconv.Itoa(int(res["Body"].([]interface{})[0].(map[string]interface{})["Id"].(float64)))
	movArka[0].EstadoMovimientoId.Id = resEstadoMovimiento[0].Id

	if err := json.Unmarshal([]byte(movArka[0].Detalle), &detalleMovimiento); err == nil {
		var resTrActa models.TransaccionActaRecibido

		urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/" + fmt.Sprint(detalleMovimiento["acta_recibido_id"])
		if err := request.GetJson(urlcrud, &resTrActa); err == nil { // Get informacion acta de api acta_recibido_crud
			transaccionActaRecibido = resTrActa
		} else {
			err = errors.New("no se encuentra el acta de recibido")
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "AprobarEntrada - request.GetJson(urlcrud, &resTrActa); err == nil",
				"err":     err,
				"status":  "404",
			}
			return nil, outputError
		}
	}

	detalle := ""
	for k, v := range detalleMovimiento {
		if k != "consecutivo" {
			detalle = detalle + k + ": " + fmt.Sprintf("%v", v) + " "
		}
	}

	var groups = make(map[int]float64)
	i := 0
	for _, elemento := range transaccionActaRecibido.Elementos {
		logs.Debug("entra:", elemento.SubgrupoCatalogoId)
		x := float64(0)
		if val, ok := groups[elemento.SubgrupoCatalogoId]; ok {
			x = val + elemento.ValorFinal
		} else {
			x = elemento.ValorFinal
		}
		groups[elemento.SubgrupoCatalogoId] = x
		i++
	}

	var trContable map[string]interface{}
	if resA, outputError := cuentasContablesHelper.AsientoContable(groups, tipomvto, "Entrada de almacen", detalle, transaccionActaRecibido.UltimoEstado.ProveedorId); res == nil || outputError != nil {
		if outputError == nil {
			outputError = map[string]interface{}{
				"funcion": "AddEntrada -cuentasContablesHelper.AsientoContable(groups, tipomvto, \"asiento contable\");",
				"err":     res,
				"status":  "502",
			}
		}
		return nil, outputError
	} else {
		trContable = resA
	}

	t := trContable["resultadoTransaccion"]
	detalleMovimiento["ConsecutivoContableId"] = t.(models.TransaccionMovimientos).ConsecutivoId
	if jsonString, err := json.Marshal(detalleMovimiento); err == nil {
		movArka[0].Detalle = string(jsonString[:])
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "AprobarEntrada - request.SendJson(urlcrud, \"PUT\", &res, &movArka)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento/" + strconv.Itoa(int(entradaId))
	if err := request.SendJson(urlcrud, "PUT", &res, &movArka[0]); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "AprobarEntrada - request.SendJson(urlcrud, \"PUT\", &res, &movArka)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo"
	procesoExterno := int64(entradaId)
	idMovArka := int(movArka[0].FormatoTipoMovimientoId.Id)
	tipoMovimientoId := models.TipoMovimiento{Id: int(res["Body"].([]interface{})[0].(map[string]interface{})["Id"].(float64))}
	movimientosKronos := models.MovimientoProcesoExterno{
		TipoMovimientoId:         &tipoMovimientoId,
		ProcesoExterno:           procesoExterno,
		Activo:                   true,
		MovimientoProcesoExterno: idMovArka,
	}

	if err := request.SendJson(urlcrud, "POST", &res, &movimientosKronos); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "AprobarEntrada - request.SendJson(urlcrud, \"POST\", &resM, &movimientosKronos)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	resultado["movimientoArka"] = movArka[0]
	resultado["transaccionContable"] = trContable["resultadoTransaccion"]
	resultado["tercero"] = trContable["tercero"]
	return resultado, nil
}

func asignarPlacaActa(actaRecibidoId int) (elementos []*models.Elemento, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "asignarPlacaActa - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	ctxPlaca, _ := beego.AppConfig.Int("contxtPlaca")
	if detalleElementos, err := actaRecibido.GetElementos(actaRecibidoId, nil); err != nil {
		outputError = map[string]interface{}{
			"funcion": "asignarPlacaActa - actaRecibido.GetElementos(actaRecibidoId)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	} else {
		for _, elemento := range detalleElementos {
			if elemento.SubgrupoCatalogoId.TipoBienId.NecesitaPlaca == true {
				if placa, err := utilsHelper.GetConsecutivo("%05.0f", ctxPlaca, "Registro Placa Arka"); err != nil {
					logs.Error(err)
					outputError = map[string]interface{}{
						"funcion": "asignarPlacaActa - utilsHelper.GetConsecutivo(\"%05.0f\", ctxPlaca, \"Registro Placa Arka\")",
						"err":     err,
						"status":  "502",
					}
					return nil, outputError
				} else {
					year, month, day := time.Now().Date()
					elemento.Placa = utilsHelper.FormatConsecutivo(fmt.Sprintf("%04d%02d%02d", year, month, day), placa, "")
					elemento_ := models.Elemento{
						Id:                 elemento.Id,
						Nombre:             elemento.Nombre,
						Cantidad:           elemento.Cantidad,
						Marca:              elemento.Marca,
						Serie:              elemento.Serie,
						UnidadMedida:       elemento.UnidadMedida,
						ValorUnitario:      elemento.ValorUnitario,
						Subtotal:           elemento.Subtotal,
						Descuento:          elemento.Descuento,
						ValorTotal:         elemento.ValorTotal,
						PorcentajeIvaId:    elemento.PorcentajeIvaId,
						ValorIva:           elemento.ValorIva,
						ValorFinal:         elemento.ValorFinal,
						Placa:              elemento.Placa,
						SubgrupoCatalogoId: elemento.SubgrupoCatalogoId.SubgrupoId.Id,
						EstadoElementoId:   &models.EstadoElemento{Id: elemento.EstadoElementoId.Id},
						ActaRecibidoId:     &models.ActaRecibido{Id: elemento.ActaRecibidoId.Id},
						Activo:             true,
					}
					elementos = append(elementos, &elemento_)
				}
			}
		}
		return elementos, nil
	}

}

//GetEncargadoElemento busca al encargado de un elemento
func GetEncargadoElemento(placa string) (idElemento map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetEncargadoElemento - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var urlelemento string
	var detalle []map[string]interface{}

	if placa == "" {
		err := fmt.Errorf("La placa no puede estar en blanco")
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "GetEncargadoElemento - placa == ''", "status": "400", "err": err}
		return nil, outputError
	}

	if id, err := actaRecibido.GetIdElementoPlaca(placa); err == nil {
		if id == "" {
			err := fmt.Errorf("La placa '%s' no ha sido asignada a una salida", placa)
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "GetEncargadoElemento - id == ''", "status": "404", "err": err}
			return nil, outputError
		}
		urlelemento = "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento/?query=ElementoActaId:" + id
		if resp, err := request.GetJsonTest(urlelemento, &detalle); err == nil && resp.StatusCode == 200 {
			cadena := detalle[0]["MovimientoId"].(map[string]interface{})["Detalle"]
			if resultado, err := utilsHelper.ConvertirStringJson(cadena); err == nil {
				idtercero := fmt.Sprintf("%v", resultado["funcionario"])
				if response, err := tercerosHelper.GetNombreTerceroById(idtercero); err == nil {
					if len(response) == 0 { // posible validación adicional:  || response["Id"].(string) == "0"
						err := fmt.Errorf("Respuesta inesperada en la respuesta de la función GetNombreTerceroById")
						logs.Error(err)
						outputError = map[string]interface{}{"funcion": "GetEncargadoElemento - len(response) == 0", "status": "404", "err": err}
						return nil, outputError
					}
					var nombrecompleto = response["NombreCompleto"].(string)
					var aux = make(map[string]interface{})
					aux = map[string]interface{}{"Id": idtercero, "NombreCompleto": nombrecompleto}
					return aux, nil
				} else {
					return nil, err
				}
			} else {
				return nil, err
			}

		} else {
			if err == nil {
				err = fmt.Errorf("Undesired Status Code: %d", resp.StatusCode)
			}
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "GetEncargadoElemento - request.GetJsonTest(urlelemento, &detalle) ", "status": "500", "err": err}
			return nil, outputError
		}
	} else {
		return nil, err
	}
}

// AnularEntrada Anula una entrada y los movimientos posteriores a esta, el acta asociada queda en estado aceptada
func AnularEntrada(movimientoId int) (response map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "AnularEntrada - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var (
		urlcrud                 string
		res                     map[string]interface{}
		resMap                  map[string]interface{}
		movimientoArka          models.Movimiento
		transaccionActaRecibido models.TransaccionActaRecibido
		movimientosKronos       models.MovimientoProcesoExterno
		detalleMovimiento       map[string]interface{}
		tipoMovimiento          models.TipoMovimiento
		estadoActa              models.EstadoActa
		estadoMovimiento        models.EstadoMovimiento
		parametroTipoDebito     models.Parametro
		parametroTipoCredito    models.Parametro
		tipoComprobanteContable models.TipoComprobanteContable
		consecutivoId           int
		consecutivo             int
		transaccion             models.TransaccionMovimientos
		cuentasSubgrupo         []models.CuentaSubgrupo
		TipoEntradaKronos       models.TipoMovimiento
	)

	res = make(map[string]interface{})

	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento?query=Id:" + strconv.Itoa(int(movimientoId))
	var resMovArka []models.Movimiento
	if err := request.GetJson(urlcrud, &resMovArka); err == nil { // Get movimiento de api movimientos_arka_crud
		movimientoArka = resMovArka[0]
		if movimientoArka.Id > 0 {

			urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo?query=ProcesoExterno:" + strconv.Itoa(int(movimientoId))
			if err = request.GetJson(urlcrud, &resMap); err == nil { // Get movimiento de api movimientos_crud
				var resMovKronos []models.MovimientoProcesoExterno
				if jsonString, err := json.Marshal(resMap["Body"]); err == nil {
					if err = json.Unmarshal(jsonString, &resMovKronos); err == nil {
						resMap = make(map[string]interface{})
						movimientosKronos = resMovKronos[0]

						if err = json.Unmarshal([]byte(movimientoArka.Detalle), &detalleMovimiento); err == nil {
							var resTrActa []models.TransaccionActaRecibido

							urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/" + fmt.Sprint(detalleMovimiento["acta_recibido_id"])
							if err = request.GetJson(urlcrud, &resTrActa); err == nil { // Get informacion acta de api acta_recibido_crud
								transaccionActaRecibido = resTrActa[0]
								var resEstadoMovimiento []models.EstadoMovimiento

								urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "estado_movimiento?query=Nombre:Entrada%20Anulada"
								if err = request.GetJson(urlcrud, &resEstadoMovimiento); err == nil { // Get parametrización estado_movimiento de api movimientos_arka_crud
									estadoMovimiento = resEstadoMovimiento[0]

									urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "tipo_movimiento?query=Nombre:Entrada%20Anulada"
									if err = request.GetJson(urlcrud, &resMap); err == nil { // Get parametrización tipo_movimiento de api movimientos_crud
										var resTipoMovimiento []models.TipoMovimiento
										if jsonString, err = json.Marshal(resMap["Body"]); err == nil {
											if err = json.Unmarshal(jsonString, &resTipoMovimiento); err == nil {
												resMap = make(map[string]interface{})
												tipoMovimiento = resTipoMovimiento[0]

												urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "tipo_movimiento?query=Nombre__iexact:" + strings.ReplaceAll(movimientoArka.FormatoTipoMovimientoId.Nombre, " ", "%20")
												if err = request.GetJson(urlcrud, &resMap); err == nil { // Get parametrización tipo_movimiento de api movimientos_crud
													if jsonString, err = json.Marshal(resMap["Body"]); err == nil {
														if err = json.Unmarshal(jsonString, &resTipoMovimiento); err == nil {
															resMap = make(map[string]interface{})
															TipoEntradaKronos = resTipoMovimiento[0]
															var resEstadoActa []models.EstadoActa

															urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "estado_acta?query=Nombre:Aceptada"
															if err = request.GetJson(urlcrud, &resEstadoActa); err == nil { // Get parametrización acta de api acta_recibido_crud
																estadoActa = resEstadoActa[0]
																movimientoArka.EstadoMovimientoId.Id = estadoMovimiento.Id
																movimientosKronos.TipoMovimientoId.Id = tipoMovimiento.Id
																transaccionActaRecibido.UltimoEstado.EstadoActaId.Id = estadoActa.Id
																transaccionActaRecibido.UltimoEstado.Id = 0

																urlcrud = "http://" + beego.AppConfig.String("parametrosService") + "parametro?query=CodigoAbreviacion:MCC"
																if err = request.GetJson(urlcrud, &resMap); err == nil { // Get parámetro tipo movimiento contable crédito
																	if jsonString, err = json.Marshal(resMap["Data"]); err == nil {
																		var parametro []models.Parametro
																		if err = json.Unmarshal(jsonString, &parametro); err == nil {
																			resMap = make(map[string]interface{})
																			parametroTipoDebito = parametro[0]

																			urlcrud = "http://" + beego.AppConfig.String("parametrosService") + "parametro?query=CodigoAbreviacion:MCD"
																			if err = request.GetJson(urlcrud, &resMap); err == nil { // Get parámetro tipo movimiento contable débito
																				if jsonString, err = json.Marshal(resMap["Data"]); err == nil {
																					if err = json.Unmarshal(jsonString, &parametro); err == nil {
																						resMap = make(map[string]interface{})
																						parametroTipoCredito = parametro[0]

																						urlcrud = "http://" + beego.AppConfig.String("cuentasContablesService") + "tipo_comprobante"
																						if err = request.GetJson(urlcrud, &resMap); err == nil { // Para obtener código del tipo de comprobante
																							for _, sliceTipoComprobante := range resMap["Body"].([]interface{}) {
																								if sliceTipoComprobante.(map[string]interface{})["TipoDocumento"] == "E" {
																									if jsonString, err = json.Marshal(sliceTipoComprobante); err == nil {
																										if err = json.Unmarshal(jsonString, &tipoComprobanteContable); err == nil {
																											resMap = make(map[string]interface{})
																										} else {
																											logs.Error(err)
																											outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
																											return nil, outputError
																										}
																									} else {
																										logs.Error(err)
																										outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
																										return nil, outputError
																									}
																								}
																							}

																							year, _, _ := time.Now().Date()
																							postConsecutivo := models.Consecutivo{
																								Id:          0,
																								ContextoId:  199,
																								Year:        year,
																								Consecutivo: 0,
																								Descripcion: "Ajustes Arka",
																								Activo:      true,
																							}
																							urlcrud = "http://" + beego.AppConfig.String("consecutivosService") + "consecutivo"
																							if err = request.SendJson(urlcrud, "POST", &resMap, &postConsecutivo); err == nil {
																								if consecutivoId, err = strconv.Atoi(fmt.Sprint(resMap["Data"].(map[string]interface{})["Id"])); err == nil {
																									if consecutivo, err = strconv.Atoi(fmt.Sprint(resMap["Data"].(map[string]interface{})["Consecutivo"])); err == nil {
																										resMap = make(map[string]interface{})
																										transaccion.ConsecutivoId = consecutivoId

																										// Se crea map para agrupar los valores totales según el código del subgrupo
																										mapSubgruposTotales := map[int]float64{}
																										for _, elemento := range transaccionActaRecibido.Elementos { // Proceso para registrar el movimiento contable para cada elemento
																											if mapSubgruposTotales[elemento.SubgrupoCatalogoId] == 0 {
																												mapSubgruposTotales[elemento.SubgrupoCatalogoId] = elemento.ValorTotal
																											} else {
																												mapSubgruposTotales[elemento.SubgrupoCatalogoId] += elemento.ValorTotal
																											}
																										}

																										etiquetas := make(map[string]interface{})
																										etiquetas["TipoComprobanteId"] = tipoComprobanteContable.Codigo
																										if jsonString, err = json.Marshal(etiquetas); err == nil {
																											transaccion.Etiquetas = string(jsonString)
																											transaccion.Activo = true
																											transaccion.FechaTransaccion = time_bogota.Tiempo_bogota()
																											transaccion.Descripcion = "Transacción para registrar movimientos contables correspondientes a entrada del sistema arka"

																											for SubgrupoId, valor := range mapSubgruposTotales {
																												var cuentaDebito models.CuentaContable
																												var cuentaCredito models.CuentaContable
																												var movimientoDebito models.MovimientoTransaccion
																												var movimientoCredito models.MovimientoTransaccion

																												urlcrud = "http://" + beego.AppConfig.String("catalogoElementosService") + "cuentas_subgrupo?query=SubgrupoId__Id:" + strconv.Itoa(SubgrupoId) + ",SubtipoMovimientoId:" + strconv.Itoa(TipoEntradaKronos.Id) + ",Activo:true"
																												if err = request.GetJson(urlcrud, &cuentasSubgrupo); err == nil { // Obtiene cuentas que deben ser afectadas

																													urlcrud = "http://" + beego.AppConfig.String("cuentasContablesService") + "nodo_cuenta_contable/" + cuentasSubgrupo[0].CuentaDebitoId
																													if err = request.GetJson(urlcrud, &resMap); err == nil { // Se trae información de cuenta débito de api cuentas_contables_crud

																														if jsonString, err = json.Marshal(resMap["Body"]); err == nil {
																															if err := json.Unmarshal(jsonString, &cuentaDebito); err == nil {
																																resMap = make(map[string]interface{})

																																movimientoDebito.NombreCuenta = cuentaDebito.Nombre
																																movimientoDebito.CuentaId = cuentaDebito.Codigo
																																movimientoDebito.TipoMovimientoId = parametroTipoCredito.Id
																																movimientoDebito.Valor = valor
																																movimientoDebito.Descripcion = "Movimiento crédito registrado desde sistema arka"
																																movimientoDebito.Activo = true
																																movimientoDebito.TerceroId = 1 // Provisional
																																transaccion.Movimientos = append(transaccion.Movimientos, movimientoDebito)

																																urlcrud = "http://" + beego.AppConfig.String("cuentasContablesService") + "nodo_cuenta_contable/" + cuentasSubgrupo[0].CuentaCreditoId
																																if err = request.GetJson(urlcrud, &resMap); err == nil { // Se trae información de cuenta crédito de api cuentas_contables_crud

																																	if jsonString, err = json.Marshal(resMap["Body"]); err == nil {
																																		if err = json.Unmarshal(jsonString, &cuentaCredito); err == nil {
																																			movimientoCredito.NombreCuenta = cuentaCredito.Nombre
																																			movimientoCredito.CuentaId = cuentaCredito.Codigo
																																			movimientoCredito.TipoMovimientoId = parametroTipoDebito.Id
																																			movimientoCredito.Valor = valor
																																			movimientoCredito.Descripcion = "Movimiento débito registrado desde sistema arka"
																																			movimientoCredito.Activo = true
																																			movimientoCredito.TerceroId = 1 // Provisional
																																			transaccion.Movimientos = append(transaccion.Movimientos, movimientoCredito)
																																		} else {
																																			logs.Error(err)
																																			outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
																																			return nil, outputError
																																		}
																																	} else {
																																		logs.Error(err)
																																		outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
																																		return nil, outputError
																																	}
																																} else {
																																	logs.Error(err)
																																	outputError = map[string]interface{}{"funcion": "AnularEntrada - cuentasContables.nodo_cuenta_contable(cuenta);", "status": "502", "err": err}
																																	return nil, outputError
																																}
																															} else {
																																logs.Error(err)
																																outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
																																return nil, outputError
																															}
																														} else {
																															logs.Error(err)
																															outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
																															return nil, outputError
																														}
																													} else {
																														logs.Error(err)
																														outputError = map[string]interface{}{"funcion": "AnularEntrada - cuentasContables.nodo_cuenta_contable(cuenta);", "status": "502", "err": err}
																														return nil, outputError
																													}
																												} else {
																													logs.Error(err)
																													outputError = map[string]interface{}{"funcion": "AnularEntrada - catalogoElementos.cuentasSubgrupo(subgrupo);", "status": "502", "err": err}
																													return nil, outputError
																												}
																											}

																											res["transaccion"] = transaccion
																											var resMovmientoContable interface{}

																											urlcrud = "http://" + beego.AppConfig.String("movimientosContablesmidService") + "transaccion_movimientos/transaccion_movimientos"
																											if err = request.SendJson(urlcrud, "POST", &resMovmientoContable, &transaccion); err == nil {
																												if resMovmientoContable.(map[string]interface{})["Status"] == "201" {
																													res["Respuesta movimientos contables Entradas"] = resMovmientoContable

																													// Anulación de salidas asociadas
																													// Si el estado de movimientoArka es Entrada Asociada a una salida, continuar con la anulación de las salidas

																													consecutivoAjuste := "H20-" + fmt.Sprintf("%05d", consecutivo) + "-" + strconv.Itoa(year)
																													detalleMovimiento["consecutivo_ajuste"] = consecutivoAjuste
																													detalleMovimiento["mov_contable_ajuste_consecutivo_id"] = transaccion.ConsecutivoId

																													if jsonString, err = json.Marshal(detalleMovimiento); err == nil {
																														movimientoArka.Detalle = string(jsonString)
																														urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento/" + strconv.Itoa(int(movimientoId))
																														if err = request.SendJson(urlcrud, "PUT", &movimientoArka, &movimientoArka); err == nil { // Put en el api movimientos_arka_crud
																															res["arka"] = movimientoArka.Detalle
																															urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo/" + strconv.Itoa(movimientoArka.Id)
																															if err = request.SendJson(urlcrud, "PUT", &movimientosKronos, &movimientosKronos); err == nil { // Put api movimientos_crud

																																urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/" + fmt.Sprint(detalleMovimiento["acta_recibido_id"])
																																if err = request.SendJson(urlcrud, "PUT", &transaccionActaRecibido, &transaccionActaRecibido); err == nil { // Puesto que se anula la entrada, el acta debe quedar disponible para volver ser asociada a una entrada
																																	res["movArkaId"] = movimientoArka.EstadoMovimientoId.Id
																																	res["movKronosId"] = movimientosKronos.TipoMovimientoId.Id
																																	res["EstadoActaId"] = transaccionActaRecibido.UltimoEstado.EstadoActaId.Id
																																} else {
																																	logs.Error(err)
																																	outputError = map[string]interface{}{"funcion": "AnularEntrada - actaRecibido.TransaccionActaRecibido(acta);", "status": "502", "err": err}
																																	return nil, outputError
																																}
																															} else {
																																logs.Error(err)
																																outputError = map[string]interface{}{"funcion": "AnularEntrada - movimientos.MovimientoProcesoExterno(id);", "status": "502", "err": err}
																																return nil, outputError
																															}
																														} else {
																															logs.Error(err)
																															outputError = map[string]interface{}{"funcion": "AnularEntrada - movimientosArka.Movimiento(id);", "status": "502", "err": err}
																															return nil, outputError
																														}
																													} else {
																														logs.Error(err)
																														outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
																														return nil, outputError
																													}
																												} else {
																													res["Respuesta movimientos contables Entradas"] = resMovmientoContable.(map[string]interface{})["Data"]
																													outputError = map[string]interface{}{"funcion": "AnularEntrada - movimientosContablesMid.postTransaccion;", "status": "502", "err": resMovmientoContable.(map[string]interface{})["Data"]}
																													return res, outputError
																												}
																											} else {
																												logs.Error(err)
																												outputError = map[string]interface{}{"funcion": "AnularEntrada - movimientosContablesMid.postTransaccion(movimiento);", "status": "502", "err": err}
																												return nil, outputError
																											}

																										} else {
																											logs.Error(err)
																											outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
																											return nil, outputError
																										}
																									} else {
																										logs.Error(err)
																										outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
																										return nil, outputError
																									}
																								} else {
																									logs.Error(err)
																									outputError = map[string]interface{}{"funcion": "AnularEntrada - consecutivos.postConsecutivo; No se retornó un consecutivo válido", "status": "502", "err": err}
																									return nil, outputError
																								}
																							} else {
																								logs.Error(err)
																								outputError = map[string]interface{}{"funcion": "AnularEntrada - consecutivos.postConsecutivo;", "status": "502", "err": err}
																								return nil, outputError
																							}
																						} else {
																							logs.Error(err)
																							outputError = map[string]interface{}{"funcion": "AnularEntrada - cuentasContablesCrud.TipoComprobante(Codigo);", "status": "502", "err": err}
																							return nil, outputError
																						}
																					} else {
																						logs.Error(err)
																						outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
																						return nil, outputError
																					}
																				} else {
																					logs.Error(err)
																					outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
																					return nil, outputError
																				}
																			} else {
																				logs.Error(err)
																				outputError = map[string]interface{}{"funcion": "AnularEntrada - parametros.Parametro(CodigoAbreviación);", "status": "502", "err": err}
																				return nil, outputError
																			}
																		} else {
																			logs.Error(err)
																			outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
																			return nil, outputError
																		}
																	} else {
																		logs.Error(err)
																		outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
																		return nil, outputError
																	}
																} else {
																	logs.Error(err)
																	outputError = map[string]interface{}{"funcion": "AnularEntrada - parametros.Parametro(CodigoAbreviación);", "status": "502", "err": err}
																	return nil, outputError
																}

															} else {
																logs.Error(err)
																outputError = map[string]interface{}{"funcion": "AnularEntrada - actaRecibido.EstadoActa", "status": "502", "err": err}
																return nil, outputError
															}

														} else {
															logs.Error(err)
															outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
															return nil, outputError
														}
													} else {
														logs.Error(err)
														outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
														return nil, outputError
													}
												} else {
													logs.Error(err)
													outputError = map[string]interface{}{"funcion": "AnularEntrada - movimientos.TipoMovimiento", "status": "502", "err": err}
													return nil, outputError
												}

											} else {
												logs.Error(err)
												outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
												return nil, outputError
											}
										} else {
											logs.Error(err)
											outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
											return nil, outputError
										}
									} else {
										logs.Error(err)
										outputError = map[string]interface{}{"funcion": "AnularEntrada - movimientos.TipoMovimiento", "status": "502", "err": err}
										return nil, outputError
									}
								} else {
									logs.Error(err)
									outputError = map[string]interface{}{"funcion": "AnularEntrada - movimientosArka.EstadoMovimiento", "status": "502", "err": err}
									return nil, outputError
								}
							} else {
								logs.Error(err)
								outputError = map[string]interface{}{"funcion": "AnularEntrada - actaRecibido.TransaccionActaRecibido(acta);", "status": "502", "err": err}
								return nil, outputError
							}
						} else {
							logs.Error(err)
							outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
							return nil, outputError
						}
					} else {
						logs.Error(err)
						outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
						return nil, outputError
					}
				} else {
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
					return nil, outputError
				}
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "AnularEntrada - movimientos.MovimientoProcesoExterno(id);", "status": "502", "err": err}
				return nil, outputError
			}
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "AnularEntrada - movimientosArka.Movimiento(id); El movimiento no existe", "status": "502", "err": err}
			return nil, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "AnularEntrada - movimientosArka.Movimiento(id);", "status": "502", "err": err}
		return nil, outputError
	}
	return res, nil
}

// GetMovimientosByActa ...
func GetMovimientosByActa(actaRecibidoId int) (movimientos map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetMovimientosByActa - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var (
		res      map[string]interface{}
		urlcrud  string
		entradas []models.Movimiento
		salidas  []map[string]interface{}
	)

	res = make(map[string]interface{})

	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento/entrada/" + strconv.Itoa(actaRecibidoId)
	if err := request.GetJson(urlcrud, &entradas); err == nil { // Se consulta la entrada asociada al acta

		for i, entrada := range entradas {
			var entradaCompleta []models.Movimiento

			urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento?query=Id:" + strconv.Itoa(entrada.Id)
			if err = request.GetJson(urlcrud, &entradaCompleta); err == nil { // Hace la consulta para obtener el detalle completo de la entrada
				entradas[i] = entradaCompleta[0]

				urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento?query=MovimientoPadreId__Id:" + strconv.Itoa(entrada.Id)
				if err = request.GetJson(urlcrud, &salidas); err == nil { // Se consultan las salidas asociada al acta

					if salidas[0]["Id"] != nil {

						for i, salida := range salidas {
							if salidaCompleta, err := salidaHelper.TraerDetalle(salida); err == nil {
								salidas[i] = salidaCompleta
							} else {
								return nil, err
							}

						}
						res["Salidas"] = salidas

						// Cuando esté completo el flujo de las bajas, incluir consulta de bajas de elementos asociadas a la entrada

					} else {
						res["Salidas"] = nil
					}
				} else {
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "GetMovimiento - movimientosArka.Movimiento(movimiento_padre_id);", "status": "502", "err": err}
					return nil, outputError
				}
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "GetMovimiento - movimientosArka.Movimiento(id);", "status": "502", "err": err}
				return nil, outputError
			}
		}
		res["Entradas"] = entradas
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "GetMovimientosByActa - movimientosArka.Movimiento(acta_id);", "status": "502", "err": err}
		return nil, outputError
	}
	return res, nil
}

func getTipoComprobanteEntradas() string {
	return "P8"
}
