package salidaHelper

import (
	"context"
	"strconv"

	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

var consultarCentroCostosSalida = movimientosArka.GetAllCentroCostos

const mensajeCentroCostosNoEncontrado = "Error en la búsqueda, consultar a soporte"

func traerDetalle(
	ctx context.Context,
	movimiento *models.Movimiento,
	salida models.FormatoSalidaCostos,
	centrosCostos map[string]models.CentroCostos,
	funcionarios map[int]models.Tercero,
) (salida_ map[string]interface{}, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("TraerDetalle - Unhandled Error!", "500")

	var (
		sede         models.EspacioFisico
		centroCostos models.CentroCostos
		funcionario  models.Tercero
	)

	if funcionarios == nil {
		funcionarios = make(map[int]models.Tercero)
	}

	if centrosCostos == nil {
		centrosCostos = make(map[string]models.CentroCostos)
	}

	buscaCentroCostos := salida.Ubicacion > 0 || salida.CentroCostos != ""
	if salida.Ubicacion > 0 {
		key := "id:" + strconv.Itoa(salida.Ubicacion)
		if val, ok := centrosCostos[key]; !ok {
			payload := "query=Id:" + strconv.Itoa(salida.Ubicacion)
			if centrosCostos_, err := consultarCentroCostosSalida(ctx, payload); err != nil {
				return nil, err
			} else if len(centrosCostos_) == 1 {
				centroCostos = centrosCostos_[0]
				centrosCostos[key] = centroCostos
			}
		} else {
			centroCostos = val
		}
	} else if salida.CentroCostos != "" {
		key := "codigo:" + salida.CentroCostos
		if val, ok := centrosCostos[key]; !ok {
			payload := "query=Codigo:" + salida.CentroCostos
			centroCostos_, err := consultarCentroCostosSalida(ctx, payload)
			if err != nil {
				return nil, err
			} else if len(centroCostos_) == 1 {
				centroCostos = centroCostos_[0]
				centrosCostos[key] = centroCostos
			}
		} else {
			centroCostos = val
		}
	}

	if buscaCentroCostos && centroCostos.Id == 0 {
		centroCostos = models.CentroCostos{
			Codigo: "0",
			Nombre: mensajeCentroCostosNoEncontrado,
		}
	}

	var dependencia *models.Dependencia
	if centroCostos.Id > 0 {
		if centroCostos.Sede == "" && centroCostos.Dependencia == "" {
			dependencia = &models.Dependencia{Nombre: centroCostos.Nombre}
		} else {
			sede = models.EspacioFisico{Nombre: centroCostos.Sede}
			dependencia = &models.Dependencia{Nombre: centroCostos.Dependencia}
		}
	} else if buscaCentroCostos {
		dependencia = &models.Dependencia{Nombre: mensajeCentroCostosNoEncontrado}
	}

	if salida.Funcionario > 0 {

		if val, ok := funcionarios[salida.Funcionario]; !ok {
			if funcionario_, err := terceros.GetTerceroById(ctx, salida.Funcionario); err != nil {
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
		"Dependencia":             dependencia,
		"Ubicacion":               centroCostos,
		"FechaCreacion":           movimiento.FechaCreacion,
		"FechaCorte":              movimiento.FechaCorte,
		"Activo":                  movimiento.Activo,
		"MovimientoPadreId":       movimiento.MovimientoPadreId,
		"FormatoTipoMovimientoId": movimiento.FormatoTipoMovimientoId,
		"EstadoMovimientoId":      movimiento.EstadoMovimientoId,
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
