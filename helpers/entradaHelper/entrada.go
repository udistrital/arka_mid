package entradaHelper

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/asientoContable"
	crud_actas "github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	crudMovimientosArka "github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	crudTerceros "github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/salidaHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

// AprobarEntrada Actualiza una entrada a estado aprobada y hace los respectivos registros en kronos y transacciones contables
func AprobarEntrada(entradaId int) (result map[string]interface{}, outputError map[string]interface{}) {

	funcion := "AprobarEntrada"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		detalleMovimiento map[string]interface{}
		historico         *models.HistoricoActa
		movimiento        *models.Movimiento
		elementos         []*models.Elemento
	)

	resultado := make(map[string]interface{})

	if mov, err := crudMovimientosArka.GetMovimientoById(entradaId); err != nil {
		return nil, err
	} else {
		movimiento = mov
	}

	if err := json.Unmarshal([]byte(movimiento.Detalle), &detalleMovimiento); err != nil {
		logs.Error(err)
		eval := " - json.Unmarshal([]byte(movimiento.Detalle), &detalleMovimiento)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	if sm, err := crudMovimientosArka.GetAllEstadoMovimiento(url.QueryEscape("Entrada Aprobada")); err != nil {
		return nil, err
	} else {
		movimiento.EstadoMovimientoId = sm[0]
	}

	query := "Activo:true,ActaRecibidoId__Id:" + fmt.Sprint(detalleMovimiento["acta_recibido_id"])
	if ha, err := crud_actas.GetAllHistoricoActa(query, "", "FechaCreacion", "desc", "", "-1"); err != nil {
		return nil, err
	} else {
		historico = ha[0]
	}

	if el_, err := crud_actas.GetAllElemento(query, "", "FechaCreacion", "desc", "", "-1"); err != nil {
		return nil, err
	} else {
		elementos = el_
	}

	detalle := ""
	for k, v := range detalleMovimiento {
		if k != "consecutivo" {
			detalle = detalle + k + ": " + fmt.Sprintf("%v", v) + " "
		}
	}

	var groups = make(map[int]float64)
	for _, elemento := range elementos {
		x := float64(0)
		if val, ok := groups[elemento.SubgrupoCatalogoId]; ok {
			x = val + elemento.ValorTotal
		} else {
			x = elemento.ValorTotal
		}
		groups[elemento.SubgrupoCatalogoId] = x
	}

	var trContable map[string]interface{}
	if tr_, err := asientoContable.AsientoContable(groups, strconv.Itoa(movimiento.FormatoTipoMovimientoId.Id), detalle, "Entrada de almacen", historico.ProveedorId, true); tr_ == nil || err != nil {
		return nil, err
	} else {
		trContable = tr_
		if tr_["errorTransaccion"].(string) != "" {
			return tr_, nil
		}
	}

	t := trContable["resultadoTransaccion"]
	detalleMovimiento["ConsecutivoContableId"] = t.(*models.DetalleTrContable).ConsecutivoId
	if jsonString, err := json.Marshal(detalleMovimiento); err != nil {
		logs.Error(err)
		eval := " - json.Marshal(detalleMovimiento)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	} else {
		movimiento.Detalle = string(jsonString[:])
	}

	if movimiento_, err := crudMovimientosArka.PutMovimiento(movimiento, movimiento.Id); err != nil {
		return nil, err
	} else {
		movimiento = movimiento_
	}

	resultado["movimientoArka"] = movimiento
	resultado["transaccionContable"] = trContable["resultadoTransaccion"]
	resultado["tercero"] = trContable["tercero"]
	resultado["errorTransaccion"] = ""

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
		return nil, err
	} else {
		for _, elemento := range detalleElementos {
			placa := ""
			if elemento.SubgrupoCatalogoId.TipoBienId.NecesitaPlaca {
				if placa_, _, err := consecutivos.Get("%05.0f", ctxPlaca, "Registro Placa Arka"); err != nil {
					return nil, err
				} else {
					year, month, day := time.Now().Date()
					placa = consecutivos.Format(fmt.Sprintf("%04d%02d%02d", year, month, day), placa_, "")
				}
			}
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
				Placa:              placa,
				SubgrupoCatalogoId: elemento.SubgrupoCatalogoId.SubgrupoId.Id,
				EstadoElementoId:   &models.EstadoElemento{Id: elemento.EstadoElementoId.Id},
				ActaRecibidoId:     &models.ActaRecibido{Id: elemento.ActaRecibidoId.Id},
				Activo:             true,
			}
			elementos = append(elementos, &elemento_)
		}
		return elementos, nil
	}

}

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
