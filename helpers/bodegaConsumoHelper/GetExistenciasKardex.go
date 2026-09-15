package bodegaConsumoHelper

import (
	"context"
	"strconv"

	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

func GetExistenciasKardex(ctx context.Context) (Elementos []map[string]interface{}, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("GetExistenciasKardex - Unhandled Error!", "500")

	var aperturas []models.Apertura
	outputError = movimientosArka.GetAperturas(ctx, true, &aperturas)
	if outputError != nil {
		return nil, outputError
	}

	for _, apertura := range aperturas {

		catalogo, err := detalleElementoCatalogo(ctx, apertura.ElementoCatalogoId)
		if err != nil {
			return nil, err
		}

		var detalle map[string]interface{}
		outputError = utilsHelper.FillStruct(apertura, &detalle)
		if outputError != nil {
			return
		}

		detalle["ElementoCatalogoId"] = catalogo
		Elementos = append(Elementos, detalle)

	}

	return Elementos, nil
}

func detalleElementoCatalogo(ctx context.Context, elementoId int) (elemento models.ElementoCatalogo, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("detalleElementoCatalogo - Unhandled Error!", "500")

	var elemento_ []models.ElementoCatalogo
	outputError = catalogoElementos.GetAllElemento(ctx, "query=Id:"+strconv.Itoa(elementoId), &elemento_)
	if outputError != nil || len(elemento_) != 1 {
		return
	}

	elemento = elemento_[0]
	return
}
