package salidaHelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/movimientosArkaHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

type Consecutivo struct {
	Id          int
	ContextoId  int
	Year        int
	Consecutivo int
	Descripcion string
	Activo      bool
}

// AsignarPlaca Transacción para asignar las placas
func AsignarPlaca(m *models.Elemento) (resultado map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "AsignarPlaca - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	year, month, day := time.Now().Date()

	consec := Consecutivo{0, 0, year, 0, "Placas", true}
	var (
		res map[string]interface{} // models.SalidaGeneral
	)

	apiCons := "http://" + beego.AppConfig.String("consecutivosService") + "consecutivo"
	putElemento := "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento/" + fmt.Sprintf("%d", m.Id)

	// Inserta salida en Movimientos ARKA
	// AsignarPlaca Transacción para asignar las placas
	if err := request.SendJson(apiCons, "POST", &res, &consec); err == nil {
		resultado, _ := res["Data"].(map[string]interface{})
		fecstring := fmt.Sprintf("%4d", year) + fmt.Sprintf("%02d", int(month)) + fmt.Sprintf("%02d", day) + fmt.Sprintf("%05.0f", resultado["Consecutivo"])
		m.Placa = fecstring
		if err := request.SendJson(putElemento, "PUT", &resultado, &m); err == nil {
			return resultado, nil
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "AsignarPlaca - request.SendJson(putElemento, \"PUT\", &resultado, &m)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "AsignarPlaca - request.SendJson(apiCons, \"POST\", &res, &consec)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}

// PostTrSalidas Completa los detalles de las salidas y hace el respectivo registro en api movimientos_arka_crud
func PostTrSalidas(m *models.SalidaGeneral) (resultado map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "PostTrSalidas - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var (
		res                 map[string][](map[string]interface{})
		resEstadoMovimiento []models.EstadoMovimiento
	)

	resultado = make(map[string]interface{})

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "estado_movimiento?query=Nombre:Salida%20En%20Trámite"
	if err := request.GetJson(urlcrud, &resEstadoMovimiento); err != nil || len(resEstadoMovimiento) == 0 {
		status := "502"
		if err == nil {
			err = errors.New("len(resEstadoMovimiento) == 0")
			status = "404"
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "PostTrSalidas - request.GetJson(urlcrud, &resEstadoMovimiento)",
			"err":     err,
			"status":  status,
		}
		return nil, outputError
	}

	ctxSalida, _ := beego.AppConfig.Int("contxtSalidaCons")
	for _, salida := range m.Salidas {

		detalle := map[string]interface{}{}
		if err := json.Unmarshal([]byte(salida.Salida.Detalle), &detalle); err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "PostTrSalidas - json.Unmarshal([]byte(salida.Salida.Detalle), &detalle)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}

		if consecutivo, _, err := utilsHelper.GetConsecutivo("%05.0f", ctxSalida, "Registro Salida Arka"); err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "PostTrSalidas - utilsHelper.GetConsecutivo(\"%05.0f\", ctxSalida, \"Registro Salida Arka\")",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		} else {
			consecutivo = utilsHelper.FormatConsecutivo(getTipoComprobanteSalidas()+"-", consecutivo, fmt.Sprintf("%s%04d", "-", time.Now().Year()))
			detalle["consecutivo"] = consecutivo
			if detalleJSON, err := json.Marshal(detalle); err != nil {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "PostTrSalidas - json.Marshal(detalle)",
					"err":     err,
					"status":  "500",
				}
				return nil, outputError
			} else {
				salida.Salida.Detalle = string(detalleJSON)
			}
		}

		salida.Salida.EstadoMovimientoId.Id = resEstadoMovimiento[0].Id
	}

	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "tr_salida"

	// Crea registros en api movimientos_arka_crud
	if err := request.SendJson(urlcrud, "POST", &res, &m); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "PostTrSalidas - request.SendJson(movArka, \"POST\", &res, &m)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	resultado["trSalida"] = res

	return resultado, nil
}

// AprobarSalida Aprobacion de una salida
func AprobarSalida(salidaId int) (result map[string]interface{}, outputError map[string]interface{}) {

	funcion := "AprobarSalida"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		detalleMovimiento map[string]interface{}
		detallePrincipal  map[string]interface{}
		trSalida          *models.TrSalida
		elementosActa     []*models.Elemento
		funcionarioId     int
		tipoMovimiento    int
	)

	resultado := make(map[string]interface{})

	if tr_, err := movimientosArkaHelper.GetTrSalida(salidaId); err != nil {
		return nil, err
	} else {
		trSalida = tr_
	}

	var idsElementos []int
	for _, el := range trSalida.Elementos {
		idsElementos = append(idsElementos, el.ElementoActaId)
	}

	fields := "SubgrupoCatalogoId,ValorTotal"
	query := "Id__in:" + utilsHelper.ArrayToString(idsElementos, "|")
	if elementos_, err := actaRecibido.GetAllElemento(query, fields, "", "", "", "-1"); err != nil {
		return nil, err
	} else {
		if len(elementos_) == 0 {
			return resultado, nil
		}
		elementosActa = elementos_
	}

	if err := json.Unmarshal([]byte(trSalida.Salida.MovimientoPadreId.Detalle), &detalleMovimiento); err != nil {
		logs.Error(err)
		eval := " - json.Unmarshal([]byte(trSalida.Salida.MovimientoPadreId.Detalle), &detalleMovimiento)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	if err := json.Unmarshal([]byte(trSalida.Salida.Detalle), &detallePrincipal); err != nil {
		logs.Error(err)
		eval := " - json.Unmarshal([]byte(trSalida.Salida.Detalle), &detallePrincipal)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	funcionario := ""
	for k, v := range detallePrincipal {
		if k == "funcionario" {
			funcionario = fmt.Sprintf("%v", v)
		}
	}

	if func_, err := strconv.Atoi(funcionario); err != nil {
		logs.Error(err)
		eval := " - strconv.Atoi(funcionario)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	} else {
		funcionarioId = func_
	}

	detalle := ""
	for k, v := range detalleMovimiento {
		if k == "consecutivo" {
			detalle = detalle + k + ": " + fmt.Sprintf("%v", v) + " "
		}
	}

	var groups = make(map[int]float64)
	for _, elemento := range elementosActa {
		x := float64(0)
		if val, ok := groups[elemento.SubgrupoCatalogoId]; ok {
			x = val + elemento.ValorTotal
		} else {
			x = elemento.ValorTotal
		}
		groups[elemento.SubgrupoCatalogoId] = x
	}

	query = "query=CodigoAbreviacion:SAL"
	if fm, err := movimientosArkaHelper.GetAllFormatoTipoMovimiento(query); err != nil {
		return nil, err
	} else {
		tipoMovimiento = fm[0].Id
	}

	var trContable map[string]interface{}
	if tr_, err := asientoContable.AsientoContable(groups, strconv.Itoa(tipoMovimiento), "Salida de almacen", detalle, funcionarioId, true); err != nil {
		return nil, err
	} else {
		trContable = tr_
		if tr_["errorTransaccion"].(string) != "" {
			return tr_, nil
		}
	}

	t := trContable["resultadoTransaccion"]
	detallePrincipal["ConsecutivoContableId"] = t.(*models.DetalleTrContable).ConsecutivoId

	if jsonString, err := json.Marshal(detallePrincipal); err != nil {
		logs.Error(err)
		eval := " - json.Marshal(detallePrincipal)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	} else {
		trSalida.Salida.Detalle = string(jsonString[:])
	}

	if sm, err := movimientosArkaHelper.GetAllEstadoMovimiento(url.QueryEscape("Salida Aprobada")); err != nil {
		return nil, err
	} else {
		trSalida.Salida.EstadoMovimientoId = sm[0]
	}

	if movimiento_, err := movimientosArkaHelper.PutMovimiento(trSalida.Salida, trSalida.Salida.Id); err != nil {
		return nil, err
	} else {
		trSalida.Salida = movimiento_
	}

	resultado["movimientoArka"] = trSalida.Salida
	resultado["transaccionContable"] = trContable["resultadoTransaccion"]
	resultado["tercero"] = trContable["tercero"]
	resultado["errorTransaccion"] = ""

	return resultado, nil
}

func GetSalida(id int) (Salida map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetSalida - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "tr_salida/" + strconv.Itoa(id)
	var salida_ map[string]interface{}
	if _, err := request.GetJsonTest(urlcrud, &salida_); err == nil {

		var data_ []map[string]interface{}
		if jsonString, err := json.Marshal(salida_["Elementos"]); err == nil {

			if err2 := json.Unmarshal(jsonString, &data_); err2 == nil {

				for i, elemento := range data_ {

					var elemento_ []map[string]interface{}
					urlcrud_ := "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento?query=Id:" + fmt.Sprintf("%v", elemento["ElementoActaId"]) + "&fields=Id,Nombre,Marca,Serie,Placa,SubgrupoCatalogoId"
					if _, err := request.GetJsonTest(urlcrud_, &elemento_); err == nil {
						var subgrupo_ []map[string]interface{}

						urlcrud_2 := "http://" + beego.AppConfig.String("catalogoElementosService") + "detalle_subgrupo?query=SubgrupoId__Id:" + fmt.Sprintf("%v", elemento_[0]["SubgrupoCatalogoId"])
						if _, err := request.GetJsonTest(urlcrud_2, &subgrupo_); err == nil {
							data_[i]["Nombre"] = elemento_[0]["Nombre"]
							data_[i]["TipoBienId"] = subgrupo_[0]["TipoBienId"]
							data_[i]["SubgrupoCatalogoId"] = subgrupo_[0]["SubgrupoId"]
							data_[i]["Marca"] = elemento_[0]["Marca"]
							data_[i]["Serie"] = elemento_[0]["Serie"]
							data_[i]["Placa"] = elemento_[0]["Placa"]

						} else {
							logs.Error(err)
							outputError = map[string]interface{}{
								"funcion": "GetSalida - request.GetJsonTest(urlcrud_2, &subgrupo_)",
								"err":     err,
								"status":  "502",
							}
							return nil, outputError
						}
					} else {
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "GetSalida - request.GetJsonTest(urlcrud_, &elemento_)",
							"err":     err,
							"status":  "502",
						}
						return nil, outputError
					}

					if _, err := request.GetJsonTest(urlcrud, &salida_); err != nil {
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "GetSalida - request.GetJsonTest(urlcrud, &salida_) (BIS)",
							"err":     err,
							"status":  "502",
						}
						return nil, outputError
					}

				}

			} else {
				logs.Error(err2)
				outputError = map[string]interface{}{
					"funcion": "GetSalida - json.Unmarshal(jsonString, &data_)",
					"err":     err2,
					"status":  "500",
				}
				return nil, outputError
			}
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetSalida - json.Marshal(salida_[\"Elementos\"])",
				"err":     err,
				"status":  "500",
			}
			return nil, outputError
		}

		if salida__, err := TraerDetalle(salida_["Salida"]); err == nil {

			Salida_final := map[string]interface{}{
				"Elementos": data_,
				"Salida":    salida__,
			}
			return Salida_final, nil

		} else {
			return nil, err
		}

	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetSalida - request.GetJsonTest(urlcrud, &salida_)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}

func GetSalidas(tramiteOnly bool) (Salidas []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetSalidas - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	query := "limit=-1&sortby=Id&order=desc&query=Activo:true,EstadoMovimientoId__Nombre"
	if tramiteOnly {
		query += url.QueryEscape(":Salida En Trámite")
	} else {
		query += url.QueryEscape("__startswith:Salida")
	}

	if salidas_, err := movimientosArkaHelper.GetAllMovimiento(query); err != nil {
		return nil, err
	} else {
		if len(salidas_) == 0 {
			return nil, nil
		}

		for _, salida := range salidas_ {
			if salida__, err := TraerDetalle(salida); err == nil {
				Salidas = append(Salidas, salida__)
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "GetSalidas - TraerDetalle(salida)",
					"err":     err,
					"status":  "502",
				}
				return nil, err
			}
		}
	}
	return Salidas, nil
}

// GetInfoSalida Retorna el funcionario y el consecutivo de una salida a partir del detalle del movimiento
func GetInfoSalida(detalle string) (funcionarioId int, consecutivo string, outputError map[string]interface{}) {

	funcion := "GetInfoSalida"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var detalle_ map[string]interface{}

	if err := json.Unmarshal([]byte(detalle), &detalle_); err != nil {
		logs.Error(err)
		eval := " - json.Unmarshal([]byte(detalle), &detalle_)"
		return 0, "", errorctrl.Error(funcion+eval, err, "500")
	}

	return int(detalle_["funcionario"].(float64)), detalle_["consecutivo"].(string), nil

}

func getTipoComprobanteSalidas() string {
	return "H21"
}