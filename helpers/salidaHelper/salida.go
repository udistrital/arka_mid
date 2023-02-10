package salidaHelper

import (
	"net/url"

	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// PostTrSalidas Completa los detalles de las salidas y hace el respectivo registro en api movimientos_arka_crud
func PostTrSalidas(m *models.SalidaGeneral, etl bool) (resultado map[string]interface{}, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("PostTrSalidas - Unhandled Error!", "500")

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

func GetSalidas(tramiteOnly bool) (Salidas []map[string]interface{}, outputError map[string]interface{}) {

	funcion := "GetSalidas - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	asignaciones := make(map[int]models.AsignacionEspacioFisicoDependencia)
	sedes := make(map[string]models.EspacioFisico)
	funcionarios := make(map[int]models.Tercero)

	query := "limit=-1&sortby=Id&order=desc&query=Activo:true,FormatoTipoMovimientoId__CodigoAbreviacion__in:SAL|SAL_CONS,EstadoMovimientoId__Nombre"
	if tramiteOnly {
		query += url.QueryEscape(":Salida En Trámite")
	} else {
		query += url.QueryEscape("__startswith:Salida")
	}

	salidas_, outputError := movimientosArka.GetAllMovimiento(query)
	if outputError != nil {
		return
	}

	for _, salida := range salidas_ {

		var formato models.FormatoSalida
		outputError = utilsHelper.Unmarshal(salida.Detalle, &formato)
		if outputError != nil {
			return
		}

		salida__, err := TraerDetalle(salida, formato, asignaciones, sedes, funcionarios)
		if err != nil {
			return nil, err
		}

		Salidas = append(Salidas, salida__)
	}

	return
}

// GetInfoSalida Retorna el funcionario de una salida a partir del detalle del movimiento
func GetInfoSalida(detalle string) (funcionarioId int, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GetInfoSalida - Unhandled Error!", "500")

	var detalle_ models.FormatoSalida
	if err := utilsHelper.Unmarshal(detalle, &detalle_); err != nil {
		return 0, err
	}

	return detalle_.Funcionario, nil
}

func getTipoComprobanteSalidas() string {
	return "H21"
}
