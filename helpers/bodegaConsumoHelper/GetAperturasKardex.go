package bodegaConsumoHelper

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

func GetAperturasKardex() (Elementos []map[string]interface{}, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GetAperturasKardex - Unhandled Error", "500")

	Elementos = make([]map[string]interface{}, 0)
	payload := "limit=-1&query=MovimientoId__FormatoTipoMovimientoId__CodigoAbreviacion:AP_KDX&sortby=FechaCreacion&order=desc"

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
			"MetodoValoracion":  detalle["Metodo_Valoracion"],
			"CantidadMinima":    detalle["Cantidad_Minima"],
			"CantidadMaxima":    detalle["Cantidad_Maxima"],
			"FechaCreacion":     elemento.FechaCreacion,
			"Observaciones":     elemento.MovimientoId.Observacion,
			"Id":                elemento.MovimientoId.Id,
			"MovimientoPadreId": elemento.MovimientoId.MovimientoPadreId,
		}

		var elemento_ []models.ElementoCatalogo
		outputError = catalogoElementos.GetAllElemento("query=Id:"+strconv.Itoa(elemento.ElementoCatalogoId), &elemento_)
		if outputError != nil {
			return
		} else if len(elemento_) == 1 {
			Elemento["ElementoCatalogoId"] = elemento_[0]
		}

		Elementos = append(Elementos, Elemento)
	}

	return

}
