package inventarioHelper

import (
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/oikos"
	"github.com/udistrital/arka_mid/helpers/mid/terceros"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// GetDetalleElemento Consulta historial de un elemento dado el id del elemento en el api acta_recibido_crud
func GetDetalleElemento(id int, Elemento *models.DetalleElementoBaja) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GetDetalleElemento - Unhandled Error!", "500")

	var (
		elemento           models.DetalleElemento
		elementoMovimiento models.ElementosMovimiento
	)

	outputError = movimientosArka.GetElementosMovimientoById(id, &elementoMovimiento)
	if outputError != nil || elementoMovimiento.Id == 0 {
		return
	}

	Elemento.Historial, outputError = movimientosArka.GetHistorialElemento(elementoMovimiento.Id, true)
	if outputError != nil {
		return
	}

	// Consulta de Marca, Nombre, Serie y Subgrupo se hace mediante el actaRecibidoHelper
	ids := []int{elementoMovimiento.ElementoActaId}
	if elementos, err := actaRecibido.GetElementos(0, ids); err != nil || len(elementos) != 1 {
		return err
	} else {
		elemento = *elementos[0]
	}

	fc, ub, outputError := GetEncargado(Elemento.Historial)
	if outputError != nil {
		return
	}

	Elemento.Ubicacion, outputError = oikos.GetSedeDependenciaUbicacion(ub)
	if outputError != nil {
		return
	}

	Elemento.Funcionario, outputError = terceros.GetInfoTerceroById(fc)
	if outputError != nil {
		return
	}

	Elemento.Id = elementoMovimiento.Id
	Elemento.Placa = elemento.Placa
	Elemento.Nombre = elemento.Nombre
	Elemento.Marca = elemento.Marca
	Elemento.Serie = elemento.Serie
	Elemento.SubgrupoCatalogoId = elemento.SubgrupoCatalogoId

	return
}

// GetEncargado Retorna el funcionario y ubicacion actual de un elemento de acuerdo a su historial
func GetEncargado(historial *models.Historial) (funcionarioId int, ubicacionId int, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GetEncargado - Unhandled Error!", "500")

	if historial.Traslados != nil {
		var detalleTr models.DetalleTraslado
		outputError = utilsHelper.Unmarshal(historial.Traslados[0].Detalle, &detalleTr)
		if outputError != nil {
			return
		}

		funcionarioId, ubicacionId = detalleTr.FuncionarioDestino, detalleTr.Ubicacion
	} else if historial.Salida != nil {
		var detalleS models.FormatoSalida
		outputError = utilsHelper.Unmarshal(historial.Salida.Detalle, &detalleS)
		if outputError != nil {
			return
		}

		funcionarioId, ubicacionId = detalleS.Funcionario, detalleS.Ubicacion
	}

	return
}
