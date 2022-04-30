package salidaHelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/asientoContable"
	crudActas "github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	crudMovimientosArka "github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

// PostTrSalidas Completa los detalles de las salidas y hace el respectivo registro en api movimientos_arka_crud
func PostTrSalidas(m *models.SalidaGeneral) (resultado map[string]interface{}, outputError map[string]interface{}) {

	funcion := "PostTrSalidas"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

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

		var consecutivo models.Consecutivo
		if err := consecutivos.Get(ctxSalida, "Registro Salida Arka", &consecutivo); err != nil {
			return nil, err
		}

		detalle["consecutivo"] = consecutivos.Format("%05d", getTipoComprobanteSalidas(), &consecutivo)
		detalle["ConsecutivoId"] = consecutivo.Id

		if detalleJSON, err := json.Marshal(detalle); err != nil {
			logs.Error(err)
			eval := " - json.Marshal(detalle)"
			return nil, errorctrl.Error(funcion+eval, err, "500")
		} else {
			salida.Salida.Detalle = string(detalleJSON)
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
		consecutivoId     int
	)

	resultado := make(map[string]interface{})

	if tr_, err := crudMovimientosArka.GetTrSalida(salidaId); err != nil {
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
	if elementos_, err := crudActas.GetAllElemento(query, fields, "", "", "", "-1"); err != nil {
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

	if val, ok := detallePrincipal["funcionario"]; ok && val != nil {
		funcionarioId = int(val.(float64))
	}

	detalle := ""
	for k, v := range detalleMovimiento {
		if k == "consecutivo" {
			detalle = "Entrada: " + fmt.Sprintf("%v", v) + " "
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
	if fm, err := crudMovimientosArka.GetAllFormatoTipoMovimiento(query); err != nil {
		return nil, err
	} else {
		tipoMovimiento = fm[0].Id
	}

	if val, ok := detallePrincipal["ConsecutivoId"]; ok && val != nil {
		consecutivoId = int(val.(float64))
	}

	var trContable map[string]interface{}
	if len(groups) > 0 && funcionarioId > 0 {
		if tr_, err := asientoContable.AsientoContable(groups, getTipoComprobanteSalidas(), strconv.Itoa(tipoMovimiento), detalle, "Salida de almacén", funcionarioId, consecutivoId, true); err != nil {
			return nil, err
		} else {
			trContable = tr_
			if tr_["errorTransaccion"].(string) != "" {
				return tr_, nil
			}
		}
	}

	if sm, err := crudMovimientosArka.GetAllEstadoMovimiento("query=Nombre:" + url.QueryEscape("Salida Aprobada")); err != nil {
		return nil, err
	} else {
		trSalida.Salida.EstadoMovimientoId = sm[0]
	}

	if movimiento_, err := crudMovimientosArka.PutMovimiento(trSalida.Salida, trSalida.Salida.Id); err != nil {
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

	funcion := "GetSalida"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		trSalida           *models.TrSalida
		detalle            map[string]interface{}
		ids                []int
		elementosActa      []*models.DetalleElemento
		elementosCompletos []*models.DetalleElemento__
	)

	if tr_, err := crudMovimientosArka.GetTrSalida(id); err != nil {
		return nil, err
	} else {
		trSalida = tr_
	}

	for _, el := range trSalida.Elementos {
		ids = append(ids, el.ElementoActaId)
	}

	if len(ids) > 0 {
		if elementosActa, outputError = actaRecibido.GetElementos(0, ids); outputError != nil {
			return nil, outputError
		}
	}

	for _, el := range elementosActa {
		var idx int
		var elemento_ *models.DetalleElemento__
		detalle := new(models.ElementosMovimiento)

		if idx = utilsHelper.FindElementoInArrayElementosMovimiento(trSalida.Elementos, el.Id); idx > -1 {
			detalle = trSalida.Elementos[idx]
			detalle.ValorResidual = detalle.ValorResidual * 100 / detalle.ValorTotal
		} else {
			detalle.ValorResidual = 0
			detalle.VidaUtil = 0
		}

		if elemento_, outputError = utilsHelper.FillElemento(el, detalle); outputError != nil {
			return nil, outputError
		}

		elementosCompletos = append(elementosCompletos, elemento_)
	}

	if salida__, err := TraerDetalle(trSalida.Salida); err != nil {
		return nil, err
	} else {
		detalle = salida__
	}

	Salida_final := map[string]interface{}{
		"Elementos": elementosCompletos,
		"Salida":    detalle,
	}

	if trSalida.Salida.EstadoMovimientoId.Nombre == "Salida Aprobada" {
		if val, ok := detalle["ConsecutivoId"]; ok && val != nil {
			if tr, err := movimientosContables.GetTransaccion(int(val.(float64)), "consecutivo", true); err != nil {
				return nil, err
			} else if len(tr.Movimientos) > 0 {
				if detalleContable, err := asientoContable.GetDetalleContable(tr.Movimientos); err != nil {
					return nil, err
				} else {
					trContable := map[string]interface{}{
						"movimientos": detalleContable,
						"concepto":    tr.Descripcion,
						"fecha":       tr.FechaTransaccion,
					}
					Salida_final["trContable"] = trContable
				}
			}
		}
	}

	return Salida_final, nil

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

	query := "limit=-1&sortby=Id&order=desc&query=Activo:true,FormatoTipoMovimientoId__CodigoAbreviacion__in:SAL|SAL_CONS,EstadoMovimientoId__Nombre"
	if tramiteOnly {
		query += url.QueryEscape(":Salida En Trámite")
	} else {
		query += url.QueryEscape("__startswith:Salida")
	}

	if salidas_, err := crudMovimientosArka.GetAllMovimiento(query); err != nil {
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
