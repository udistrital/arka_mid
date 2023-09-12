package salidaHelper

import (
	"regexp"
	"strconv"

	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/oikos"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

func traerDetalle(movimiento *models.Movimiento, salida models.FormatoSalidaCostos,
	asignaciones map[int]models.AsignacionEspacioFisicoDependencia,
	sedes map[string]models.EspacioFisico,
	funcionarios map[int]models.Tercero) (salida_ map[string]interface{}, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("TraerDetalle - Unhandled Error!", "500")

	var (
		query       string
		sede        models.EspacioFisico
		ubicacion   models.AsignacionEspacioFisicoDependencia
		funcionario models.Tercero
	)

	if asignaciones == nil {
		asignaciones = make(map[int]models.AsignacionEspacioFisicoDependencia)
	}

	if sedes == nil {
		sedes = make(map[string]models.EspacioFisico)
	}

	if funcionarios == nil {
		funcionarios = make(map[int]models.Tercero)
	}

	if salida.Ubicacion > 0 {
		if val, ok := asignaciones[salida.Ubicacion]; !ok {
			query = "query=Id:" + strconv.Itoa(salida.Ubicacion)
			if asignacion_, err := oikos.GetAllAsignacion(query); err != nil {
				return nil, err
			} else if len(asignacion_) == 1 {
				ubicacion = asignacion_[0]
				asignaciones[salida.Ubicacion] = ubicacion
			}
		} else {
			ubicacion = val
		}
	} else if salida.CentroCostos != "" {
		payload := "query=Codigo:" + salida.CentroCostos
		centroCostos, outputError := movimientosArka.GetAllCentroCostos(payload)
		if outputError != nil {
			return nil, outputError
		} else if len(centroCostos) == 1 {
			if centroCostos[0].Sede == "" && centroCostos[0].Dependencia == "" {
				ubicacion = models.AsignacionEspacioFisicoDependencia{
					DependenciaId: &models.Dependencia{Nombre: centroCostos[0].Nombre},
				}
			} else {
				if centroCostos[0].Sede != "" {
					sede = models.EspacioFisico{Nombre: centroCostos[0].Sede}
				}

				if centroCostos[0].Dependencia != "" {
					ubicacion.DependenciaId = &models.Dependencia{Nombre: centroCostos[0].Dependencia}
				}
			}
		}

	}

	if ubicacion.Id > 0 && ubicacion.EspacioFisicoId.CodigoAbreviacion != "" {
		rgxp := regexp.MustCompile(`\d.*`)
		str := ubicacion.EspacioFisicoId.CodigoAbreviacion
		str = str[0:2] + rgxp.ReplaceAllString(str[2:], "")

		if val, ok := sedes[str]; !ok {
			sede_, err := oikos.GetSedeEspacioFisico(*ubicacion.EspacioFisicoId)
			if err != nil {
				return nil, err
			} else if sede_.Id > 0 {
				sede = sede_
				sedes[str] = sede
			}
		} else {
			sede = val
		}
	}

	if salida.Funcionario > 0 {

		if val, ok := funcionarios[salida.Funcionario]; !ok {
			if funcionario_, err := terceros.GetTerceroById(salida.Funcionario); err != nil {
				return nil, err
			} else {
				funcionario = *funcionario_
				funcionarios[salida.Funcionario] = *funcionario_
			}
		} else {
			funcionario = val
		}
	}

	Salida2 := map[string]interface{}{
		"Id":                      movimiento.Id,
		"Observacion":             movimiento.Observacion,
		"Sede":                    sede,
		"Dependencia":             ubicacion.DependenciaId,
		"Ubicacion":               ubicacion,
		"FechaCreacion":           movimiento.FechaCreacion,
		"FechaModificacion":       movimiento.FechaModificacion,
		"Activo":                  movimiento.Activo,
		"MovimientoPadreId":       movimiento.MovimientoPadreId,
		"FormatoTipoMovimientoId": movimiento.FormatoTipoMovimientoId,
		"EstadoMovimientoId":      movimiento.EstadoMovimientoId.Id,
		"Consecutivo":             movimiento.Consecutivo,
		"ConsecutivoId":           movimiento.ConsecutivoId,
		"Funcionario":             funcionario,
	}

	return Salida2, nil

}

// GetInfoSalida Retorna el funcionario de una salida a partir del detalle del movimiento
func GetInfoSalida(detalle string) (funcionarioId int, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("GetInfoSalida - Unhandled Error!", "500")

	var detalle_ models.FormatoSalida
	outputError = utilsHelper.Unmarshal(detalle, &detalle_)
	if outputError != nil {
		return
	}

	funcionarioId = detalle_.Funcionario
	return
}
