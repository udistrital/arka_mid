package terceros

import (
	crudTerceros "github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// GetDetalleFuncionario Consulta El nombre, número de identificación, correo y cargo asociado a un funcionario
func GetDetalleFuncionario(id int) (DetalleFuncionario *models.DetalleFuncionario, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("GetDetalleFuncionario - Unhandled Error!", "500")

	DetalleFuncionario = new(models.DetalleFuncionario)

	// Consulta información general y documento de identidad
	tercero_, outputError := crudTerceros.GetTrTerceroIdentificacionById(id)
	if outputError != nil {
		return
	}

	// Consulta correo
	correo_, outputError := crudTerceros.GetCorreo(id)
	if outputError != nil {
		return
	}

	// Consulta cargo
	cargo_, outputError := GetCargoFuncionario(id)
	if outputError != nil {
		return
	}

	DetalleFuncionario.Tercero = []models.DetalleTercero{tercero_}
	DetalleFuncionario.Correo = correo_
	DetalleFuncionario.Cargo = cargo_
	return DetalleFuncionario, nil
}

// GetInfoTerceroById Consulta El nombre y  número de identificación de cualquier tercero
func GetInfoTerceroById(id int) (InfoTercero *models.InfoTercero, outputError map[string]interface{}) {

	funcion := "GetInfoTerceroById"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	InfoTercero = new(models.InfoTercero)

	// Consulta nombre
	if tercero_, err := crudTerceros.GetTerceroById(id); err != nil {
		return nil, err
	} else {
		InfoTercero.Tercero = tercero_
	}

	// Consulta documento
	if documento_, err := GetDocumentoTercero(id); err != nil {
		return nil, err
	} else {
		if len(documento_) != 0 {
			InfoTercero.Identificacion = documento_[0]
		}
	}

	return InfoTercero, nil
}
