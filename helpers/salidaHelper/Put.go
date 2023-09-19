package salidaHelper

import (
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

func Put(m *models.SalidaGeneral, salidaId int) (resultado map[string]interface{}, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("Put - Unhandled Error!", "500")

	var (
		detalleOriginal    models.FormatoSalida
		estadoMovimientoId int
	)

	resultado = make(map[string]interface{})

	if len(m.Salidas) == 0 {
		return
	}

	// El objetivo es generar los respectivos consecutivos en caso de generarse más de una salida a partir de la original

	// Se consulta la salida original
	salidaOriginal, outputError := movimientosArka.GetMovimientoById(salidaId)
	if outputError != nil || salidaOriginal.EstadoMovimientoId.Nombre != "Salida Rechazada" {
		return
	}

	outputError = utilsHelper.Unmarshal(salidaOriginal.Detalle, &detalleOriginal)
	if outputError != nil {
		return
	}

	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&estadoMovimientoId, "Salida En Trámite")
	if outputError != nil {
		return
	}

	index := -1
	if len(m.Salidas) == 1 {
		// Si no se generan nuevas salidas, tan solo se debe actualizar el funcionario y ubicación de la salida original así como la vida útil y valor residual de los elementos
		index = 0
	}

	for idx, l := range m.Salidas {
		// Si se generaron salidas a partir de la original, se debe asignar un consecutivo a cada una y una de ellas debe tener el original
		// Se debe decidir a cuál de las nuevas asignarle el id y el consecutivo original
		if index > -1 {
			break
		}

		var detalleNuevo models.FormatoSalida
		outputError = utilsHelper.Unmarshal(l.Salida.Detalle, &detalleNuevo)
		if outputError != nil {
			return
		}

		if detalleNuevo.Funcionario == detalleOriginal.Funcionario && detalleNuevo.Ubicacion == detalleOriginal.Ubicacion {
			index = idx
		} else if detalleNuevo.Funcionario == detalleOriginal.Funcionario {
			index = idx
		} else if detalleNuevo.Ubicacion == detalleOriginal.Ubicacion {
			index = idx
		}
	}

	if index == -1 {
		index = 0
	}

	for idx, salida := range m.Salidas {

		var id int
		if idx == index {
			id = salidaId
			salida.Salida.Consecutivo = salidaOriginal.Consecutivo
			salida.Salida.ConsecutivoId = salidaOriginal.ConsecutivoId
		}

		salida.Salida.Id = id
		salida.Salida.EstadoMovimientoId.Id = estadoMovimientoId
		outputError = setConsecutivoSalida(salida.Salida)
		if outputError != nil {
			return
		}
	}

	trRes, outputError := movimientosArka.PutTrSalida(m)
	if outputError != nil {
		return
	}

	resultado["trSalida"] = trRes
	return
}

func setConsecutivoSalida(salida *models.Movimiento) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("setConsecutivoSalida - Unhandled Error!", "500")

	if salida.Consecutivo == nil || salida.ConsecutivoId == nil || *salida.Consecutivo == "" || *salida.ConsecutivoId <= 0 {

		var consecutivo models.Consecutivo
		outputError = consecutivos.Get("contxtSalidaCons", "Registro Salida Arka", &consecutivo)
		if outputError != nil {
			return
		}

		salida.Consecutivo = utilsHelper.String(consecutivos.Format("%05d", getTipoComprobanteSalidas(), &consecutivo))
		salida.ConsecutivoId = &consecutivo.Id
	}

	return
}
