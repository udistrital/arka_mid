package entradaHelper

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	crudTerceros "github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/salidaHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

//GetEncargadoElemento busca al encargado de un elemento
func GetEncargadoElemento(placa string) (idElemento *models.Tercero, outputError map[string]interface{}) {

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
		err := fmt.Errorf("la placa no puede estar en blanco")
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "GetEncargadoElemento - placa == ''", "status": "400", "err": err}
		return nil, outputError
	}

	if id, err := actaRecibido.GetIdElementoPlaca(placa); err == nil {
		if id == "" {
			err := fmt.Errorf("la placa '%s' no ha sido asignada a una salida", placa)
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "GetEncargadoElemento - id == ''", "status": "404", "err": err}
			return nil, outputError
		}
		urlelemento = "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento/?query=ElementoActaId:" + id
		if resp, err := request.GetJsonTest(urlelemento, &detalle); err == nil && resp.StatusCode == 200 {
			cadena := detalle[0]["MovimientoId"].(map[string]interface{})["Detalle"]
			if resultado, err := utilsHelper.ConvertirStringJson(cadena); err == nil {
				idtercero := int(resultado["funcionario"].(float64))
				if tercero, err := crudTerceros.GetTerceroById(idtercero); err == nil {
					return tercero, nil
				} else {
					return nil, err
				}
			} else {
				return nil, err
			}

		} else {
			if err == nil {
				err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
			}
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "GetEncargadoElemento - request.GetJsonTest(urlelemento, &detalle) ", "status": "500", "err": err}
			return nil, outputError
		}
	} else {
		return nil, err
	}
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

						for i := range salidas {
							if salidaCompleta, err := salidaHelper.TraerDetalle(nil, models.FormatoSalida{}, nil, nil, nil); err == nil {
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

// GetConsecutivoEntrada Retorna el consecutivo de una entrada a partir del detalle del movimiento.
func GetConsecutivoEntrada(detalle string) (consecutivo string, outputError map[string]interface{}) {

	funcion := "GetConsecutivoEntrada"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")
	var (
		detalle_ map[string]interface{}
	)

	if err := json.Unmarshal([]byte(detalle), &detalle_); err != nil {
		logs.Error(err)
		eval := " - json.Unmarshal([]byte(detalle), &detalle_)"
		return "", errorctrl.Error(funcion+eval, err, "500")
	}

	return detalle_["consecutivo"].(string), nil
}

func getTipoComprobanteEntradas() string {
	return "P8"
}
