package salidaHelper

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

var getMovimientoByIDSalidaPost = movimientosArka.GetMovimientoById
var getAllElementosMovimientoSalidaPost = movimientosArka.GetAllElementosMovimiento
var putElementosMovimientoSalidaPost = movimientosArka.PutElementosMovimiento

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

		if outputError = liberarElementosDeSalidasAnuladas(&m.Salidas[idx]); outputError != nil {
			return nil, outputError
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
	entradaPadre, err := getMovimientoByIDSalidaPost(salida.MovimientoPadreId.Id)
	if err != nil {
		return err
	}
	if entradaPadre == nil || entradaPadre.Id <= 0 {
		return errorCtrl.Error("normalizarNuevaSalida - MovimientoPadreId", "la entrada padre indicada no existe", "400")
	}
	if entradaPadre.EstadoMovimientoId == nil || (entradaPadre.EstadoMovimientoId.Nombre != "Entrada Aprobada" && entradaPadre.EstadoMovimientoId.Nombre != "Entrada Con Salida") {
		return errorCtrl.Error("normalizarNuevaSalida - EstadoMovimientoId", "la entrada padre no está en un estado válido para registrar salidas", "400")
	}
	if entradaPadre.FormatoTipoMovimientoId == nil || !strings.HasPrefix(entradaPadre.FormatoTipoMovimientoId.CodigoAbreviacion, "ENT_") {
		return errorCtrl.Error("normalizarNuevaSalida - FormatoTipoMovimientoId", "el movimiento padre indicado no corresponde a una entrada válida", "400")
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

func liberarElementosDeSalidasAnuladas(trSalida *models.TrSalida) (outputError map[string]interface{}) {
	if trSalida == nil || trSalida.Salida == nil || len(trSalida.Elementos) == 0 {
		return nil
	}

	desactivados := make(map[int]struct{})
	parentID := 0
	if trSalida.Salida.MovimientoPadreId != nil {
		parentID = trSalida.Salida.MovimientoPadreId.Id
	}

	for _, elemento := range trSalida.Elementos {
		if elemento == nil || elemento.ElementoActaId == nil || *elemento.ElementoActaId <= 0 {
			continue
		}

		query := "limit=-1&query=Activo:true,ElementoActaId:" + strconv.Itoa(*elemento.ElementoActaId) +
			",MovimientoId__EstadoMovimientoId__Nombre:" + url.QueryEscape(estadoSalidaAnulada)
		if parentID > 0 {
			query += ",MovimientoId__MovimientoPadreId__Id:" + strconv.Itoa(parentID)
		}

		elementosMovimiento, err := getAllElementosMovimientoSalidaPost(query)
		if err != nil {
			return err
		}

		for _, elementoMovimiento := range elementosMovimiento {
			if elementoMovimiento == nil || elementoMovimiento.Id <= 0 {
				continue
			}
			if _, ok := desactivados[elementoMovimiento.Id]; ok {
				continue
			}

			elementoMovimiento.Activo = false
			if _, err = putElementosMovimientoSalidaPost(elementoMovimiento, elementoMovimiento.Id); err != nil {
				return err
			}

			desactivados[elementoMovimiento.Id] = struct{}{}
		}
	}

	return nil
}

func getTipoComprobanteSalidas() string {
	return "H21"
}
