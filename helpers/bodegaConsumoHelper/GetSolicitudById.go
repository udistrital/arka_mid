package bodegaConsumoHelper

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/oikos"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// GetSolicitudById trae el nombre de un encargado por su id
func GetSolicitudById(id int) (Solicitud map[string]interface{}, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("GetSolicitudById - Unhandled Error", "500")

	var solicitud_ = make(map[string]interface{})
	var elementos___ []map[string]interface{}

	mov, err := movimientosArka.GetAllMovimiento("query=Id:" + strconv.Itoa(id))
	if err != nil || len(mov) != 1 {
		return nil, err
	}

	var detalle models.FormatoSolicitudBodega
	outputError = utilsHelper.Unmarshal(mov[0].Detalle, &detalle)
	if outputError != nil {
		return
	}

	outputError = utilsHelper.FillStruct(mov[0], &solicitud_)
	if outputError != nil {
		return
	}

	tercero, err := terceros.GetNombreTerceroById(detalle.Funcionario)
	if err != nil {
		return nil, err
	}

	for _, elementos := range detalle.Elementos {
		Elemento__, err := traerElementoSolicitud(elementos)
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

func traerElementoSolicitud(Elemento models.ElementoSolicitud_) (Elemento_ map[string]interface{}, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("traerElementoSolicitud - Unhandled Error", "500")

	ubicacionInfo, outputError := oikos.GetSedeDependenciaUbicacion(Elemento.Ubicacion)
	if outputError != nil {
		return
	}

	ultimo, outputError := ultimoMovimientoKardex(Elemento.ElementoCatalogoId)
	if outputError != nil {
		return
	}

	outputError = utilsHelper.FillStruct(ultimo, &Elemento_)
	if outputError != nil {
		return
	}

	catalogo, outputError := detalleElementoCatalogo(Elemento.ElementoCatalogoId)
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

func ultimoMovimientoKardex(elementoId int) (ultimo models.ElementosMovimiento, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("ultimoMovimientoKardex - Unhandled Error!", "500")

	payload := "limit=1&sortby=FechaCreacion&order=desc&fields=ElementoCatalogoId,Id,SaldoCantidad,SaldoValor&query=ElementoCatalogoId:"
	elemento, err := movimientosArka.GetAllElementosMovimiento(payload + strconv.Itoa(elementoId))
	if err != nil || len(elemento) != 1 {
		return ultimo, err
	}

	ultimo = *elemento[0]
	return
}
