package bodegaConsumoHelper

import (
	"context"
	"strconv"

	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/oikos"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	errorCtrl "github.com/udistrital/utils_oas/v2/errorctrl"
)

// GetSolicitudById trae el nombre de un encargado por su id
func GetSolicitudById(ctx context.Context, id int) (Solicitud map[string]interface{}, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("GetSolicitudById - Unhandled Error", "500")

	var solicitud_ = make(map[string]interface{})
	var elementos___ []map[string]interface{}

	mov, outputError := movimientosArka.GetMovimientoById(ctx, id)
	if outputError != nil {
		return
	}

	var detalle models.FormatoSolicitudBodega
	outputError = utilsHelper.Unmarshal(mov.Detalle, &detalle)
	if outputError != nil {
		return
	}

	outputError = utilsHelper.FillStruct(mov, &solicitud_)
	if outputError != nil {
		return
	}

	tercero, outputError := terceros.GetNombreTerceroById(ctx, detalle.Funcionario)
	if outputError != nil {
		return
	}

	for _, elementos := range detalle.Elementos {
		Elemento__, err := traerElementoSolicitud(ctx, elementos)
		if err != nil {
			return nil, err
		}

		Elemento__["Cantidad"] = elementos.Cantidad
		Elemento__["CantidadAprobada"] = elementos.CantidadAprobada
		elementos___ = append(elementos___, Elemento__)
	}

	solicitud_["Funcionario"] = tercero
	Solicitud = map[string]interface{}{
		"Solicitud": solicitud_,
		"Elementos": elementos___,
	}

	return

}

func traerElementoSolicitud(ctx context.Context, Elemento models.ElementoSolicitud_) (Elemento_ map[string]interface{}, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("traerElementoSolicitud - Unhandled Error", "500")

	ubicacionInfo, outputError := oikos.GetSedeDependenciaUbicacion(ctx, Elemento.Ubicacion)
	if outputError != nil {
		return
	}

	ultimo, outputError := ultimoMovimientoKardex(ctx, Elemento.ElementoCatalogoId)
	if outputError != nil {
		return
	}

	outputError = utilsHelper.FillStruct(ultimo, &Elemento_)
	if outputError != nil {
		return
	}

	catalogo, outputError := detalleElementoCatalogo(ctx, Elemento.ElementoCatalogoId)
	if outputError != nil {
		return
	}

	Elemento_["ElementoCatalogoId"] = catalogo

	if ubicacionInfo == nil || ubicacionInfo.Ubicacion == nil {
		return
	}

	Elemento_["Sede"] = ubicacionInfo.Sede
	Elemento_["Dependencia"] = ubicacionInfo.Dependencia
	Elemento_["Ubicacion"] = ubicacionInfo.Ubicacion.EspacioFisicoId

	return
}

func ultimoMovimientoKardex(ctx context.Context, elementoId int) (ultimo models.ElementosMovimiento, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("ultimoMovimientoKardex - Unhandled Error!", "500")

	payload := "limit=1&sortby=FechaCreacion&order=desc&fields=ElementoCatalogoId,Id,SaldoCantidad,SaldoValor&query=ElementoCatalogoId:"
	elemento, err := movimientosArka.GetAllElementosMovimiento(ctx, payload+strconv.Itoa(elementoId))
	if err != nil || len(elemento) != 1 {
		return ultimo, err
	}

	ultimo = *elemento[0]
	return
}
