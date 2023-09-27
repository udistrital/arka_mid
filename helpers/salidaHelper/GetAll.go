package salidaHelper

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

func GetAll(estados []string, fechaCreacion, fechaAprobacion, consecutivo, entrada,
	sortby, order string, limit, page int) (Salidas []map[string]interface{}, total string, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("GetAll - Unhandled Error!", "500")

	asignaciones := make(map[int]models.AsignacionEspacioFisicoDependencia)
	sedes := make(map[string]models.EspacioFisico)
	centrosCostos := make(map[string]models.CentroCostos)
	funcionarios := make(map[int]models.Tercero)
	Salidas = make([]map[string]interface{}, 0)

	if order != "" && (sortby == "Consecutivo" || sortby == "FechaCreacion" || sortby == "FechaCorte" || sortby == "MovimientoPadreId" || sortby == "EstadoMovimientoId") {
		order = strings.ToLower(order)
		if sortby == "MovimientoPadreId" {
			sortby = "MovimientoPadreId__FechaCreacion"
		} else if sortby == "EstadoMovimientoId" {
			sortby = "EstadoMovimientoId__Nombre"
		}
	} else {
		sortby = "FechaCreacion"
		order = "desc"
	}

	payload := "limit=" + fmt.Sprint(limit) +
		"&offset=" + fmt.Sprint(limit*(page-1)) +
		"&sortby=" + sortby +
		"&order=" + order +
		"&query=Activo:true,FormatoTipoMovimientoId__CodigoAbreviacion__in:SAL|SAL_CONS"

	if len(estados) > 0 {
		payload += ",EstadoMovimientoId__Id__Nombre:" + url.QueryEscape(strings.Join(estados, "|"))
	}

	if fechaCreacion != "" {
		payload += ",FechaCreacion__contains:" + strings.ReplaceAll(fechaCreacion, "/", "-")
	}

	if fechaAprobacion != "" {
		payload += ",FechaCorte__contains:" + strings.ReplaceAll(fechaAprobacion, "/", "-")
	}

	if consecutivo != "" {
		payload += ",Consecutivo__icontains:" + consecutivo
	}

	if entrada != "" {
		payload += ",MovimientoPadreId__Consecutivo__icontains:" + entrada
	}

	salidas_, total, outputError := movimientosArka.GetAllMovimiento(payload)
	if outputError != nil {
		return
	}

	for _, salida := range salidas_ {

		var formato models.FormatoSalidaCostos
		outputError = utilsHelper.Unmarshal(salida.Detalle, &formato)
		if outputError != nil {
			return
		}

		salida_, err := traerDetalle(salida, formato, asignaciones, sedes, centrosCostos, funcionarios)
		if err != nil {
			outputError = err
			return
		}

		Salidas = append(Salidas, salida_)
	}

	return
}
