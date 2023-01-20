package bodegaConsumoHelper

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/oikos"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

func traerElementoSolicitud(Elemento models.ElementoSolicitud_) (Elemento_ map[string]interface{}, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("traerElementoSolicitud - Unhandled Error", "500")

	ubicacionInfo, err := oikos.GetSedeDependenciaUbicacion(Elemento.Ubicacion)
	if err != nil {
		return nil, err
	}

	ultimo, err := ultimoMovimientoKardex(Elemento.ElementoCatalogoId)
	if err != nil {
		return nil, err
	}

	outputError = utilsHelper.FillStruct(ultimo, &Elemento_)
	if outputError != nil {
		return
	}

	catalogo, err := detalleElementoCatalogo(Elemento.ElementoCatalogoId)
	if err != nil {
		return nil, err
	}

	Elemento_["ElementoCatalogoId"] = catalogo
	Elemento_["Sede"] = ubicacionInfo.Sede
	Elemento_["Dependencia"] = ubicacionInfo.Dependencia
	Elemento_["Ubicacion"] = ubicacionInfo.Ubicacion.EspacioFisicoId

	return
}

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

func ultimoMovimientoKardex(elementoId int) (ultimo models.ElementosMovimiento, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("ultimoMovimientoKardex - Unhandled Error!", "500")

	payload := "limit=1&sortby=FechaCreacion&order=desc&fields=ElementoCatalogoId,Id,SaldoCantidad,SaldoValor&query=ElementoCatalogoId:"
	elemento, err := movimientosArka.GetAllElementosMovimiento(payload + strconv.Itoa(elementoId))
	if err != nil || len(elemento) != 1 {
		return ultimo, err
	}

	ultimo = *elemento[0]
	return
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
