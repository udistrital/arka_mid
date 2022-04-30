package salidaHelper

import (
	"encoding/json"
	"net/url"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

func PutTrSalidas(m *models.SalidaGeneral, salidaId int) (resultado map[string]interface{}, outputError map[string]interface{}) {

	funcion := "PutTrSalidas"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")
	var (
		estadoMovimiento *models.EstadoMovimiento
		salidaOriginal   *models.Movimiento
	)

	resultado = make(map[string]interface{})

	// El objetivo es generar los respectivos consecutivos en caso de generarse más de una salida a partir de la original

	if estadosMovimiento, err := movimientosArka.GetAllEstadoMovimiento("query=Nombre:" + url.QueryEscape("Salida En Trámite")); err != nil {
		return nil, err
	} else {
		estadoMovimiento = estadosMovimiento[0]
	}

	// En caso de generarse más de una salida, se debe actualizar

	if len(m.Salidas) == 1 {
		// Si no se generan nuevas salidas, simplemente se debe actualizar el funcionario y la ubicación del movimiento original

		m.Salidas[0].Salida.EstadoMovimientoId.Id = estadoMovimiento.Id
		if trRes, err := movimientosArka.PutTrSalida(m); err != nil {
			return nil, err
		} else {
			resultado["trSalida"] = trRes
		}

	} else {

		// Si se generaron salidas a partir de la original, se debe asignar un consecutivo a cada una y una de ellas debe tener el original

		// Se consulta la salida original
		ctxSalida, _ := beego.AppConfig.Int("contxtSalidaCons")

		// Se consulta el movimiento
		if movimientoA, err := movimientosArka.GetMovimientoById(salidaId); err != nil {
			return nil, err
		} else {
			salidaOriginal = movimientoA
		}

		detalleOriginal := map[string]interface{}{}
		if err := json.Unmarshal([]byte(salidaOriginal.Detalle), &detalleOriginal); err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "PutTrSalidas - json.Unmarshal([]byte(salidaOriginal.Detalle), &detalleOriginal)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}

		// Se debe decidir a cuál de las nuevas asignarle el id y el consecutivo original
		index := -1
		detalleNueva := map[string]interface{}{}
		for idx, l := range m.Salidas {
			if err := json.Unmarshal([]byte(l.Salida.Detalle), &detalleNueva); err != nil {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "PutTrSalidas - json.Unmarshal([]byte(l.Salida.Detalle), &detalleNueva)",
					"err":     err,
					"status":  "502",
				}
				return nil, outputError
			}
			funcNuevo := detalleNueva["funcionario"]
			funcOriginal := detalleOriginal["funcionario"]
			ubcNuevo := detalleNueva["ubicacion"]
			ubcOriginal := detalleOriginal["ubicacion"]
			if funcNuevo == funcOriginal && ubcNuevo == ubcOriginal {
				index = idx
				break
			} else if funcNuevo == funcOriginal {
				index = idx
				break
			} else if ubcNuevo == ubcOriginal {
				index = idx
				break
			}
		}

		for idx, salida := range m.Salidas {

			salida.Salida.EstadoMovimientoId.Id = estadoMovimiento.Id
			detalle := map[string]interface{}{}
			if err := json.Unmarshal([]byte(salida.Salida.Detalle), &detalle); err != nil {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "PutTrSalidas - json.Unmarshal([]byte(salida.Salida.Detalle), &detalle)",
					"err":     err,
					"status":  "502",
				}
				return nil, outputError
			}

			if idx != index {
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
					// Si ninguna salida tiene el mismo funcionario o ubicación que la original, se asigna el id de la original a la primera del arreglo
					if index == -1 && idx == 0 {
						salida.Salida.Id = salidaId
					}
				}

			} else {
				detalle["consecutivo"] = detalleOriginal["consecutivo"]
				detalle["ConsecutivoId"] = detalleOriginal["ConsecutivoId"]
				if detalleJSON, err := json.Marshal(detalle); err != nil {
					logs.Error(err)
					outputError = map[string]interface{}{
						"funcion": "PutTrSalidas - json.Marshal(detalle)",
						"err":     err,
						"status":  "500",
					}
					return nil, outputError
				} else {
					salida.Salida.Detalle = string(detalleJSON)
					salida.Salida.Id = salidaId
				}
			}
		}

		// Hace el put api movimientos_arka_crud
		if trRes, err := movimientosArka.PutTrSalida(m); err != nil {
			return nil, err
		} else {
			resultado["trSalida"] = trRes
		}
	}

	return resultado, nil
}
