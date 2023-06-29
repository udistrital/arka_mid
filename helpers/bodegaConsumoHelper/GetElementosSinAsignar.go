package bodegaConsumoHelper

import (
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

func GetElementosSinAsignar() (Elementos []map[string]interface{}, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("GetElementosSinAsignar - Unhandled Error", "500")

	payload := "limit=-1&query=Activo:true,MovimientoId__FormatoTipoMovimientoId__CodigoAbreviacion:SAL_CONS" +
		",MovimientoId__EstadoMovimientoId__Nombre:Salida%20Aprobada"
	elementos, err := movimientosArka.GetAllElementosMovimiento(payload)
	if err != nil {
		return nil, err
	}

	Elementos = make([]map[string]interface{}, 0)
	subgruposBuffer := make(map[int]models.Subgrupo)

	for _, el := range elementos {

		var el_ models.Elemento
		outputError = actaRecibido.GetElementoById(*el.ElementoActaId, &el_)
		if outputError != nil {
			return
		}

		_, ok := subgruposBuffer[el_.SubgrupoCatalogoId]
		if !ok {
			sg, err := catalogoElementos.GetSubgrupoById(el_.SubgrupoCatalogoId)
			if err != nil {
				return nil, err
			}
			subgruposBuffer[el_.SubgrupoCatalogoId] = sg
		}

		detalle := make(map[string]interface{})
		outputError = utilsHelper.FillStruct(el, &detalle)
		if outputError != nil {
			return
		}

		detalle["Nombre"] = el_.Nombre
		detalle["Marca"] = el_.Marca
		detalle["Serie"] = el_.Serie
		detalle["SubgrupoCatalogoId"] = subgruposBuffer[el_.SubgrupoCatalogoId]

		Elementos = append(Elementos, detalle)
	}

	return
}
