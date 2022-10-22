package salidaHelper

import (
	"net/url"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

// PostTrSalidas Completa los detalles de las salidas y hace el respectivo registro en api movimientos_arka_crud
func PostTrSalidas(m *models.SalidaGeneral, etl bool) (resultado map[string]interface{}, outputError map[string]interface{}) {

	funcion := "PostTrSalidas"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		res                map[string][](map[string]interface{})
		estadoMovimientoId int
	)

	resultado = make(map[string]interface{})

	if err := movimientosArka.GetEstadoMovimientoIdByNombre(&estadoMovimientoId, "Salida En Trámite"); err != nil {
		return nil, err
	}

	for _, salida := range m.Salidas {

		if !etl {
			if err := setDetalleSalida("", 0, &salida.Salida.Detalle); err != nil {
				return nil, err
			}
		}

		salida.Salida.EstadoMovimientoId.Id = estadoMovimientoId
	}

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "tr_salida"

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

func GetSalida(id int) (Salida map[string]interface{}, outputError map[string]interface{}) {

	funcion := "GetSalida"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		trSalida      *models.TrSalida
		detalle       map[string]interface{}
		ids           []int
		elementosActa []*models.DetalleElemento
	)

	if tr_, err := movimientosArka.GetTrSalida(id); err != nil {
		return nil, err
	} else if tr_.Salida.FormatoTipoMovimientoId.CodigoAbreviacion == "SAL" ||
		tr_.Salida.FormatoTipoMovimientoId.CodigoAbreviacion == "SAL_BOD" {
		trSalida = tr_
	} else {
		return
	}

	for _, el := range trSalida.Elementos {
		ids = append(ids, el.ElementoActaId)
	}

	if len(ids) > 0 {
		if elementosActa, outputError = actaRecibido.GetElementos(0, ids); outputError != nil {
			return nil, outputError
		}
	}

	var elementosCompletos = make([]models.DetalleElementoSalida, 0)
	for _, el := range elementosActa {

		if idx := utilsHelper.FindElementoInArrayElementosMovimiento(trSalida.Elementos, el.Id); idx > -1 {

			detalle := models.DetalleElementoSalida{
				Cantidad:           el.Cantidad,
				ElementoActaId:     el.Id,
				Id:                 trSalida.Elementos[idx].Id,
				Marca:              el.Marca,
				Nombre:             el.Nombre,
				Placa:              el.Placa,
				Serie:              el.Serie,
				SubgrupoCatalogoId: el.SubgrupoCatalogoId,
				ValorResidual:      (trSalida.Elementos[idx].ValorResidual * 10000) / (trSalida.Elementos[idx].ValorTotal * 100),
				ValorTotal:         trSalida.Elementos[idx].ValorTotal,
				VidaUtil:           trSalida.Elementos[idx].VidaUtil,
			}

			elementosCompletos = append(elementosCompletos, detalle)
		}

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
		if val, ok := detalle["ConsecutivoId"].(int); ok && val > 0 {
			if tr, err := movimientosContables.GetTransaccion(val, "consecutivo", true); err != nil {
				return nil, err
			} else if len(tr.Movimientos) > 0 {
				if detalleContable, err := asientoContable.GetDetalleContable(tr.Movimientos, nil); err != nil {
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

	query := "limit=20&sortby=Id&order=desc&query=Activo:true,FormatoTipoMovimientoId__CodigoAbreviacion__in:SAL|SAL_CONS,EstadoMovimientoId__Nombre"
	if tramiteOnly {
		query += url.QueryEscape(":Salida En Trámite")
	} else {
		query += url.QueryEscape("__startswith:Salida")
	}

	if salidas_, err := movimientosArka.GetAllMovimiento(query); err != nil {
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

	defer errorctrl.ErrorControlFunction("GetInfoSalida - Unhandled Error!", "500")

	var detalle_ models.FormatoSalida
	if err := utilsHelper.Unmarshal(detalle, &detalle_); err != nil {
		return 0, "", err
	}

	return detalle_.Funcionario, detalle_.Consecutivo, nil
}

func getTipoComprobanteSalidas() string {
	return "H21"
}
