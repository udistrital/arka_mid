package salidaHelper

import (
	"net/url"

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

		if outputError = normalizarNuevaSalida(&m.Salidas[idx], estadoMovimientoId); outputError != nil {
			return nil, outputError
		}

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
	if outputError == nil {
		outputError = completarSalidasPersistidas(m)
	}
	resultado["trSalida"] = m

	return
}

func normalizarNuevaSalida(trSalida *models.TrSalida, estadoMovimientoId int) (outputError map[string]interface{}) {
	if trSalida == nil || trSalida.Salida == nil {
		return errorCtrl.Error("normalizarNuevaSalida - salida", "salida nil", "400")
	}

	salida := trSalida.Salida
	salida.Id = 0
	salida.FechaCorte = nil
	salida.EstadoMovimientoId = &models.EstadoMovimiento{Id: estadoMovimientoId, Nombre: "Salida En Trámite"}

	if salida.MovimientoPadreId == nil || salida.MovimientoPadreId.Id <= 0 {
		return errorCtrl.Error("normalizarNuevaSalida - MovimientoPadreId", "la salida debe tener una entrada padre válida", "400")
	}
	salida.MovimientoPadreId = &models.Movimiento{Id: salida.MovimientoPadreId.Id}

	if salida.FormatoTipoMovimientoId == nil || salida.FormatoTipoMovimientoId.Id <= 0 {
		return errorCtrl.Error("normalizarNuevaSalida - FormatoTipoMovimientoId", "la salida debe tener un formato de movimiento válido", "400")
	}
	salida.FormatoTipoMovimientoId = &models.FormatoTipoMovimiento{Id: salida.FormatoTipoMovimientoId.Id}

	for idx := range trSalida.Elementos {
		elemento := trSalida.Elementos[idx]
		if elemento == nil {
			return errorCtrl.Error("normalizarNuevaSalida - Elementos", "uno de los elementos de la salida es nil", "400")
		}
		if elemento.ElementoActaId == nil || *elemento.ElementoActaId <= 0 {
			return errorCtrl.Error("normalizarNuevaSalida - ElementoActaId", "uno de los elementos no tiene ElementoActaId válido", "400")
		}

		elemento.Id = 0
		elemento.MovimientoId = nil
	}

	return nil
}

func completarSalidasPersistidas(m *models.SalidaGeneral) (outputError map[string]interface{}) {
	for idx := range m.Salidas {
		salida := m.Salidas[idx].Salida
		if salida == nil || salida.Id > 0 || salida.Consecutivo == nil || *salida.Consecutivo == "" {
			continue
		}

		query := "limit=1&sortby=Id&order=desc&query=Consecutivo:" + url.QueryEscape(*salida.Consecutivo)
		movimientos, _, err := movimientosArka.GetAllMovimiento(query)
		if err != nil {
			return err
		}
		if len(movimientos) == 0 || movimientos[0] == nil || movimientos[0].Id <= 0 {
			return errorCtrl.Error("completarSalidasPersistidas - GetAllMovimiento", "no se confirmó la persistencia de la salida creada", "502")
		}

		m.Salidas[idx].Salida = movimientos[0]
	}

	return nil
}

func getTipoComprobanteSalidas() string {
	return "H21"
}
