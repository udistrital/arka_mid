package tercerosMidHelper

import (
	"github.com/udistrital/arka_mid/helpers/tercerosHelper"
	"github.com/udistrital/arka_mid/models"
)

// GetDetalle Consulta El nombre, número de identificación, correo y cargo asociado a un funcionario
func GetDetalleFuncionario(id int) (DetalleFuncionario *models.DetalleFuncionario, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetDetalleFuncionario", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	DetalleFuncionario = new(models.DetalleFuncionario)

	// Consulta información general y documento de identidad
	if tercero_, err := GetFuncionario(id); err != nil {
		return nil, err
	} else {
		DetalleFuncionario.Tercero = tercero_
	}

	// Consulta correo
	if correo_, err := tercerosHelper.GetCorreo(id); err != nil {
		return nil, err
	} else {
		DetalleFuncionario.Correo = correo_
	}

	// Consulta cargo
	if cargo_, err := GetCargoFuncionario(id); err != nil {
		return nil, err
	} else {
		DetalleFuncionario.Cargo = cargo_
	}

	return DetalleFuncionario, nil
}
