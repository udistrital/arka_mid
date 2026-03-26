package depreciacionHelper

import (
	"github.com/udistrital/arka_mid/helpers/crud/configuracion"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// RechazarCierre Verifica el estado de las cuentas contables y actualiza el estado del cierre.
func RechazarCierre(info *models.InfoDepreciacion, resultado *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("RechazarCierre - Unhandled Error!", "500")

	var (
		detalle    *models.FormatoDepreciacion
		parametros []models.ParametroConfiguracion
	)

	if err := configuracion.GetAllParametro("Nombre:cierreEnCurso", &parametros); err != nil {
		return err
	} else if len(parametros) != 1 || parametros[0].Valor != "true" {
		return
	}

	mov_, outputError := movimientosArka.GetMovimientoById(info.Id)
	if outputError != nil || mov_.EstadoMovimientoId.Nombre != "Cierre En Curso" {
		return
	}

	resultado.Movimiento = *mov_
	if err := movimientosArka.GetEstadoMovimientoIdByNombre(&resultado.Movimiento.EstadoMovimientoId.Id, "Cierre Rechazado"); err != nil {
		return err
	}

	if err := utilsHelper.Unmarshal(resultado.Movimiento.Detalle, &detalle); err != nil {
		return err
	}

	detalle.RazonRechazo = info.RazonRechazo
	if err := utilsHelper.Marshal(detalle, &resultado.Movimiento.Detalle); err != nil {
		return err
	}

	outputError = movimientosArka.PutMovimiento(&resultado.Movimiento, info.Id)
	if outputError != nil {
		return
	}

	parametros[0].Valor = "false"
	if err := configuracion.PutParametro(parametros[0].Id, &parametros[0]); err != nil {
		return err
	}

	return
}
