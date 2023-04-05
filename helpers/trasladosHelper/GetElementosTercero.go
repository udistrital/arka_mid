package trasladoshelper

import (
	"net/url"

	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/mid/terceros"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

func GetElementosTercero(terceroId int, inventario *models.InventarioTercero) (outputError map[string]interface{}) {

	funcion := "GetElementosTercero - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	var (
		elementosF    []int
		elementosM    []*models.ElementosMovimiento
		elementosActa []*models.Elemento
	)

	inventario.Elementos = make([]models.DetalleElementoPlaca, 0)

	if tercero, err := terceros.GetDetalleFuncionario(terceroId); err != nil {
		return err
	} else {
		inventario.Tercero = *tercero
	}

	// Consulta lista de elementos asignados al tercero
	if elemento_, err := movimientosArka.GetElementosFuncionario(terceroId); err != nil {
		return err
	} else {
		elementosF = elemento_
	}

	// Consulta id del elemento en el api acta_recibido_crud
	if len(elementosF) > 0 {
		query := "limit=-1&sortby=ElementoActaId&order=desc&query=Id__in:"
		query += url.QueryEscape(utilsHelper.ArrayToString(elementosF, "|"))
		if elementoMovimiento_, err := movimientosArka.GetAllElementosMovimiento(query); err != nil {
			return err
		} else {
			elementosM = elementoMovimiento_
		}
	} else {
		return
	}

	ids := []int{}
	for _, el := range elementosM {
		ids = append(ids, *el.ElementoActaId)
	}

	// Consulta de Nombre, Placa, Marca, Serie se hace al api acta_recibido_crud
	query := "Id__in:" + utilsHelper.ArrayToString(ids, "|")
	if elemento_, err := actaRecibido.GetAllElemento(query, "", "Id", "desc", "", "-1"); err != nil {
		return err
	} else {
		elementosActa = elemento_
	}

	if len(elementosActa) == len(elementosM) {
		for i := 0; i < len(elementosActa); i++ {
			if elementosActa[i].Placa != "" {

				var elemento = models.DetalleElementoPlaca{
					Id:             elementosM[i].Id,
					ElementoActaId: elementosActa[i].Id,
					Placa:          elementosActa[i].Placa,
					Nombre:         elementosActa[i].Nombre,
					Marca:          elementosActa[i].Marca,
					Serie:          elementosActa[i].Serie,
					Valor:          elementosActa[i].ValorTotal,
				}

				inventario.Elementos = append(inventario.Elementos, elemento)
			}
		}
	}

	return
}
