package bodegaConsumoHelper

import (
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

func GetAperturasKardex() (Elementos []map[string]interface{}, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("GetAperturasKardex - Unhandled Error", "500")

	Elementos = make([]map[string]interface{}, 0)

	var aperturas []models.Apertura
	outputError = movimientosArka.GetAperturas(false, &aperturas)
	if outputError != nil {
		return nil, outputError
	}

	for _, elemento := range aperturas {

		catalogo, err := detalleElementoCatalogo(elemento.ElementoCatalogoId)
		if err != nil {
			return nil, err
		}

		Elemento := map[string]interface{}{
			"FechaCreacion":      elemento.FechaCreacion,
			"SaldoCantidad":      elemento.SaldoCantidad,
			"MetodoValoracion":   elemento.MetodoValoracion,
			"CantidadMinima":     elemento.CantidadMinima,
			"CantidadMaxima":     elemento.CantidadMaxima,
			"ElementoCatalogoId": catalogo,
		}

		Elementos = append(Elementos, Elemento)
	}

	return

}
