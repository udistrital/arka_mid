package salidaHelper

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

func PutTrSalidas(m *models.SalidaGeneral, salidaId int) (resultado map[string]interface{}, outputError map[string]interface{}) {

	funcion := "PutTrSalidas - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	var (
		salidaOriginal     *models.Movimiento
		estadoMovimientoId int
		query              string
	)

	resultado = make(map[string]interface{})

	// El objetivo es generar los respectivos consecutivos en caso de generarse más de una salida a partir de la original

	// Se consulta la salida original
	query = "limit=1&query=Id:" + strconv.Itoa(salidaId)
	if salida_, err := movimientosArka.GetAllMovimiento(query); err != nil {
		return nil, err
	} else if len(salida_) == 1 && salida_[0].EstadoMovimientoId.Nombre == "Salida Rechazada" {
		salidaOriginal = salida_[0]
	} else {
		return
	}

	if err := movimientosArka.GetEstadoMovimientoIdByNombre(&estadoMovimientoId, "Salida En Trámite"); err != nil {
		return nil, err
	}

	if len(m.Salidas) == 1 {
		// Si no se generan nuevas salidas, tan solo se debe actualizar el funcionario y la ubicación de la salida original y la vida útil y valor residual de los elementos

		m.Salidas[0].Salida.EstadoMovimientoId.Id = estadoMovimientoId
		m.Salidas[0].Salida.Id = salidaId
		if trRes, err := movimientosArka.PutTrSalida(m); err != nil {
			return nil, err
		} else {
			resultado["trSalida"] = trRes
		}

	} else {
		// Si se generaron salidas a partir de la original, se debe asignar un consecutivo a cada una y una de ellas debe tener el original

		// Se consulta la salida original
		ctxSalida, _ := beego.AppConfig.Int("contxtSalidaCons")

		detalleOriginal := map[string]interface{}{}
		if err := utilsHelper.Unmarshal(salidaOriginal.Detalle, &detalleOriginal); err != nil {
			return nil, err
		}

		// Se debe decidir a cuál de las nuevas asignarle el id y el consecutivo original
		index := -1
		detalleNueva := map[string]interface{}{}
		for idx, l := range m.Salidas {
			if err := utilsHelper.Unmarshal(l.Salida.Detalle, &detalleNueva); err != nil {
				return nil, err
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

		if index == -1 {
			index = 0
		}

		for idx, salida := range m.Salidas {
			salida.Salida.EstadoMovimientoId.Id = estadoMovimientoId
			detalle := map[string]interface{}{}
			if err := utilsHelper.Unmarshal(salida.Salida.Detalle, &detalle); err != nil {
				return nil, err
			}

			if idx != index {
				var consecutivo models.Consecutivo
				if err := consecutivos.Get(ctxSalida, "Registro Salida Arka", &consecutivo); err != nil {
					return nil, err
				}

				detalle["consecutivo"] = consecutivos.Format("%05d", getTipoComprobanteSalidas(), &consecutivo)
				detalle["ConsecutivoId"] = consecutivo.Id
				salida.Salida.Id = 0
				if err := utilsHelper.Marshal(detalle, &salida.Salida.Detalle); err != nil {
					return nil, err
				}

			} else {
				detalle["consecutivo"] = detalleOriginal["consecutivo"]
				detalle["ConsecutivoId"] = detalleOriginal["ConsecutivoId"]
				salida.Salida.Id = salidaId
				if err := utilsHelper.Marshal(detalle, &salida.Salida.Detalle); err != nil {
					return nil, err
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
