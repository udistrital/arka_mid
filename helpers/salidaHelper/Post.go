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

	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&estadoMovimientoId, "Salida En Trámite")
	if outputError != nil {
		return
	}

	for _, salida := range m.Salidas {

		salida.Salida.EstadoMovimientoId = &models.EstadoMovimiento{Id: estadoMovimientoId}
		if !etl {
			outputError = setConsecutivoSalida(salida.Salida)
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
