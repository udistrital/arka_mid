package salidaHelper

import (
	"net/url"

	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

func GetAll(tramiteOnly bool) (Salidas []map[string]interface{}, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("GetAll - Unhandled Error!", "500")

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
