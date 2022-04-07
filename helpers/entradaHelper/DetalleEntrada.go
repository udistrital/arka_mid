package entradaHelper

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/contratoHelper"
	"github.com/udistrital/arka_mid/helpers/movimientosArkaHelper"
	"github.com/udistrital/arka_mid/helpers/movimientosContablesMidHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// DetalleEntrada Consulta el detalle de una entrada incluyendo la transaccion contable (si aplica)
func DetalleEntrada(entradaId int) (result map[string]interface{}, outputError map[string]interface{}) {

	funcion := "DetalleEntrada"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		detalle    map[string]interface{}
		movimiento *models.Movimiento
		query      string
	)

	resultado := make(map[string]interface{})

	query = "query=Id:" + strconv.Itoa(entradaId)
	if mov, err := movimientosArkaHelper.GetAllMovimiento(query); err != nil {
		return nil, err
	} else if len(mov) > 0 {
		movimiento = mov[0]
	}

	if err := json.Unmarshal([]byte(movimiento.Detalle), &detalle); err != nil {
		logs.Error(err)
		eval := " - json.Unmarshal([]byte(movimiento.Detalle), &detalleMovimiento)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	if val, ok := detalle["contrato_id"]; ok {
		if val_, ok := detalle["vigencia_contrato"]; ok {
			if contrato, err := contratoHelper.GetContrato(int(val.(float64)), val_.(string)); err != nil {
				return nil, err
			} else {
				resultado["contrato"] = contrato["contrato"]
			}
		}
	}

	if movimiento.EstadoMovimientoId.Nombre == "Entrada Aprobada" || movimiento.EstadoMovimientoId.Nombre == "Entrada Con Salida" {
		if val, ok := detalle["ConsecutivoContableId"]; ok {
			if tr, err := movimientosContablesMidHelper.GetTransaccion(int(val.(float64)), "consecutivo", true); err != nil {
				return nil, err
			} else {
				if detalleContable, err := asientoContable.GetDetalleContable(tr.Movimientos); err != nil {
					return nil, err
				} else {
					resultado["trContable"] = detalleContable
				}
			}
		}
	}

	resultado["movimiento"] = movimiento

	return resultado, nil
}
