package bodegaConsumoHelper

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// GetSolicitudById trae el nombre de un encargado por su id
func GetSolicitudById(id int) (Solicitud map[string]interface{}, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GetSolicitudById - Unhandled Error", "500")

	var solicitud_ = make([]map[string]interface{}, 1)
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

	outputError = utilsHelper.FillStruct(mov[0], &solicitud_[0])
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

	solicitud_[0]["Funcionario"] = tercero
	Solicitud = map[string]interface{}{
		"Solicitud": solicitud_,
		"Elementos": elementos___,
	}

	return

}
