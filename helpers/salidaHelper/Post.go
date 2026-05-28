package salidaHelper

import (
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// Post Completa los detalles de las salidas y hace el respectivo registro en api movimientos_arka_crud
func Post(m *models.SalidaGeneral, etl bool) (resultado map[string]interface{}, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("Post - Unhandled Error!", "500")

	var estadoMovimientoId int
	resultado = make(map[string]interface{})

	if m == nil || len(m.Salidas) == 0 {
		return nil, errorCtrl.Error("Post - validacion salidas", "debe especificar al menos una salida", "400")
	}

	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&estadoMovimientoId, "Salida En Trámite")
	if outputError != nil {
		return
	}

	for idx := range m.Salidas {
		salida := m.Salidas[idx].Salida
		if salida == nil {
			return nil, errorCtrl.Error("Post - validacion salida", "una de las salidas no tiene movimiento asociado", "400")
		}

		salida.Id = 0
		salida.EstadoMovimientoId = &models.EstadoMovimiento{Id: estadoMovimientoId, Nombre: "Salida En Trámite"}
		if !etl {
			salida.Consecutivo = nil
			salida.ConsecutivoId = nil
			outputError = setConsecutivoSalida(salida)
			if outputError != nil {
				return
			}
		}
	}

	outputError = movimientosArka.PostTrSalida(m)
	resultado["trSalida"] = m

	return
}

func getTipoComprobanteSalidas() string {
	return "H21"
}
