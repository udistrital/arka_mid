package bodegaConsumoHelper

import (
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/utils_oas/errorctrl"
)

func GetAperturasKardex() (Elementos []map[string]interface{}, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GetAperturasKardex - Unhandled Error", "500")

	Elementos = make([]map[string]interface{}, 0)
	payload := "limit=-1&query=MovimientoId__FormatoTipoMovimientoId__CodigoAbreviacion:AP_KDX" +
		"&sortby=FechaCreacion&order=desc&fields=FechaCreacion,MovimientoId,ElementoCatalogoId"

	aperturas, err := movimientosArka.GetAllElementosMovimiento(payload)
	if err != nil {
		return nil, err
	}

	for _, elemento := range aperturas {

		var detalle map[string]interface{}
		outputError = utilsHelper.Unmarshal(elemento.MovimientoId.Detalle, &detalle)
		if outputError != nil {
			return
		}

		Elemento := map[string]interface{}{
			"Id":               elemento.MovimientoId.Id,
			"FechaCreacion":    elemento.MovimientoId.FechaCreacion,
			"MetodoValoracion": detalle["Metodo_Valoracion"],
			"CantidadMinima":   detalle["Cantidad_Minima"],
			"CantidadMaxima":   detalle["Cantidad_Maxima"],
		}

		ultimo, err := ultimoMovimientoKardex(elemento.ElementoCatalogoId)
		if err != nil {
			return nil, err
		}

		catalogo, err := detalleElementoCatalogo(elemento.ElementoCatalogoId)
		if err != nil {
			return nil, err
		}

		Elemento["ElementoCatalogoId"] = catalogo
		Elemento["SaldoCantidad"] = ultimo.SaldoCantidad
		Elementos = append(Elementos, Elemento)
	}

	return

}
