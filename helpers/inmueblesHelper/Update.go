package inmuebleshelper

import (
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/models"
)

func Update(inmueble *models.Inmueble) (resultado models.ResultadoMovimiento, outputError map[string]interface{}) {

	if inmueble.ElementoMovimiento.ValorTotal <= 0 {
		resultado.Error = "No se indicó un valor no nulo para el inmueble."
	} else if inmueble.SubgrupoId.Id <= 0 {
		resultado.Error = "No se indicó la clase del inmueble."
	} else if inmueble.Cuentas.CuentaCreditoId.Id == "" || inmueble.Cuentas.CuentaDebitoId.Id == "" {
		resultado.Error = "No se indicaron las cuentas para el inmueble."
	} else if inmueble.ElementoMovimiento.VidaUtil > 0 {
		if inmueble.CuentasMediciones.CuentaCreditoId.Id == "" || inmueble.CuentasMediciones.CuentaDebitoId.Id == "" {
			resultado.Error = "No se indicaron las cuentas para la depreciación del inmueble."
		}
	}

	if resultado.Error != "" {
		return
	}

	var elemento models.Elemento
	outputError = actaRecibido.GetElementoById(inmueble.Elemento.Id, &elemento)
	if outputError != nil {
		return
	}

	elemento.EspacioFisicoId = inmueble.EspacioFisico.Id
	elemento.Nombre = inmueble.Elemento.Nombre
	elemento.SubgrupoCatalogoId = inmueble.SubgrupoId.Id

	outputError = actaRecibido.PutElemento(&elemento, elemento.Id)
	if outputError != nil {
		return
	}

	resultado.Error, outputError = registrarCuentas(*inmueble)

	return
}
