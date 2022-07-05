package salidaHelper

import (
	"regexp"
	"strconv"

	"github.com/udistrital/arka_mid/helpers/crud/oikos"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

func TraerDetalle(salida interface{}) (salida_ map[string]interface{}, outputError map[string]interface{}) {

	funcion := "TraerDetalle - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	var (
		data      models.Movimiento
		data2     models.FormatoSalida
		ubicacion models.AsignacionEspacioFisicoDependencia
		sede      models.EspacioFisico
		query     string
	)

	if err := utilsHelper.FillStruct(salida, &data); err != nil {
		return nil, err
	}

	if err := utilsHelper.Unmarshal(data.Detalle, &data2); err != nil {
		return nil, err
	}

	if data2.Ubicacion > 0 {

		query = "?query=Id:" + strconv.Itoa(data2.Ubicacion)
		if asignacion_, err := oikos.GetAllAsignacion(query); err != nil {
			return nil, err
		} else if len(asignacion_) > 0 {
			ubicacion = *asignacion_[0]
			ubicacion.EspacioFisicoId.Id = ubicacion.Id
		}
	}

	if ubicacion.Id > 0 && ubicacion.EspacioFisicoId.CodigoAbreviacion != "" {
		rgxp := regexp.MustCompile("[0-9]")
		str := rgxp.ReplaceAllString(ubicacion.EspacioFisicoId.CodigoAbreviacion, "")

		query = "?query=CodigoAbreviacion:" + str
		if sede_, err := oikos.GetAllEspacioFisico(query); err != nil {
			return nil, err
		} else if len(sede_) > 0 {
			sede = *sede_[0]
		}
	}

	Salida2 := map[string]interface{}{
		"Id":                      data.Id,
		"Observacion":             data.Observacion,
		"Sede":                    sede,
		"Dependencia":             ubicacion.DependenciaId,
		"Ubicacion":               ubicacion.EspacioFisicoId,
		"FechaCreacion":           data.FechaCreacion,
		"FechaModificacion":       data.FechaModificacion,
		"Activo":                  data.Activo,
		"MovimientoPadreId":       data.MovimientoPadreId,
		"FormatoTipoMovimientoId": data.FormatoTipoMovimientoId,
		"EstadoMovimientoId":      data.EstadoMovimientoId.Id,
		"Consecutivo":             data2.Consecutivo,
		"ConsecutivoId":           data2.ConsecutivoId,
	}

	if data2.Funcionario > 0 {
		if funcionario_, err := terceros.GetTerceroById(data2.Funcionario); err != nil {
			return nil, err
		} else {
			Salida2["Funcionario"] = funcionario_
		}
	}

	return Salida2, nil

}
