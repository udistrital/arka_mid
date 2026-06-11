package trasladoshelper

import (
	"errors"

	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// Post Crea registro de traslado en estado en trámite
func Post(traslado *models.Movimiento) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("Post - Unhandled Error!", "500")

	var consecutivo models.Consecutivo
	outputError = consecutivos.Get("contxtAjusteCons", "Registro Traslado Arka", &consecutivo)
	if outputError != nil {
		return
	}

	traslado.Consecutivo = utilsHelper.String(consecutivos.Format("%05d", getTipoComprobanteTraslados(), &consecutivo))
	traslado.ConsecutivoId = &consecutivo.Id

	outputError = movimientosArka.PostMovimiento(traslado)

	return
}

// PostInterno crea registro de traslado interno en estado confirmado
func PostInterno(traslado *models.Movimiento) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("PostInterno - Unhandled Error!", "500")

	if err := validarTrasladoInterno(traslado); err != nil {
		return errorCtrl.Error("PostInterno - validarTrasladoInterno(traslado)", err, "400")
	}

	var consecutivo models.Consecutivo
	outputError = consecutivos.Get("contxtAjusteCons", "Registro Traslado Arka", &consecutivo)
	if outputError != nil {
		return
	}

	if traslado.EstadoMovimientoId == nil {
		traslado.EstadoMovimientoId = &models.EstadoMovimiento{}
	}
	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&traslado.EstadoMovimientoId.Id, "Traslado Confirmado")
	if outputError != nil {
		return
	}

	traslado.Consecutivo = utilsHelper.String(consecutivos.Format("%05d", getTipoComprobanteTraslados(), &consecutivo))
	traslado.ConsecutivoId = &consecutivo.Id

	outputError = movimientosArka.PostMovimiento(traslado)

	return
}

func validarTrasladoInterno(traslado *models.Movimiento) error {
	if traslado == nil {
		return errors.New("se debe especificar un traslado válido")
	}

	var detalle models.FormatoTraslado
	if err := utilsHelper.Unmarshal(traslado.Detalle, &detalle); err != nil {
		if e, ok := err["err"].(error); ok {
			return e
		}
		return errors.New("detalle de traslado inválido")
	}

	if detalle.FuncionarioOrigen <= 0 {
		return errors.New("se debe especificar un funcionario origen válido")
	}

	if detalle.FuncionarioDestino <= 0 {
		return errors.New("se debe especificar un funcionario destino válido")
	}

	if detalle.FuncionarioOrigen == detalle.FuncionarioDestino {
		return errors.New("no se puede realizar un traslado al mismo funcionario")
	}

	return nil
}
