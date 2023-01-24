package bodegaConsumoHelper

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

func GetExistenciasKardex() (Elementos []map[string]interface{}, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GetExistenciasKardex - Unhandled Error!", "500")

	var aperturas []models.Apertura
	outputError = movimientosArka.GetAperturas(true, &aperturas)
	if outputError != nil {
		return nil, outputError
	}

	for _, apertura := range aperturas {

		catalogo, err := detalleElementoCatalogo(apertura.ElementoCatalogoId)
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

func detalleElementoCatalogo(elementoId int) (elemento models.ElementoCatalogo, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("detalleElementoCatalogo - Unhandled Error!", "500")

	var elemento_ []models.ElementoCatalogo
	outputError = catalogoElementos.GetAllElemento("query=Id:"+strconv.Itoa(elementoId), &elemento_)
	if outputError != nil || len(elemento_) != 1 {
		return
	}

	elemento = elemento_[0]
	return
}
