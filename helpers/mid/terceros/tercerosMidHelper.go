package terceros

import (
	"context"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"

	crudTerceros "github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	requestV2 "github.com/udistrital/utils_oas/v2/request"
)

// GetDetalleFuncionario Consulta El nombre, número de identificación, correo y cargo asociado a un funcionario
func GetDetalleFuncionario(ctx context.Context, id int) (DetalleFuncionario *models.DetalleFuncionario, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("GetDetalleFuncionario - Unhandled Error!", "500")

	DetalleFuncionario = new(models.DetalleFuncionario)

	// Consulta información general y documento de identidad
	tercero_, outputError := crudTerceros.GetTrTerceroIdentificacionById(ctx, id)
	if outputError != nil {
		return
	}

	// Consulta correo
	correo_, outputError := crudTerceros.GetCorreo(ctx, id)
	if outputError != nil {
		return
	}

	// Consulta cargo
	cargo_, outputError := GetCargoFuncionario(ctx, id)
	if outputError != nil {
		return
	}

	DetalleFuncionario.Tercero = []models.DetalleTercero{tercero_}
	DetalleFuncionario.Correo = correo_
	DetalleFuncionario.Cargo = cargo_
	return DetalleFuncionario, nil
}

// GetInfoTerceroById Consulta El nombre y  número de identificación de cualquier tercero
func GetInfoTerceroById(ctx context.Context, id int) (InfoTercero *models.InfoTercero, outputError map[string]interface{}) {

	funcion := "GetInfoTerceroById"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	InfoTercero = new(models.InfoTercero)

	// Consulta nombre
	tercerosService, _ := beego.AppConfig.String("tercerosService")
	urltercero := tercerosService + "tercero/" + strconv.Itoa(id)
	tercero_ := new(models.Tercero)
	if _, err := requestV2.GetWithContext(ctx, urltercero, &tercero_); err != nil {
		eval := " - requestV2.GetWithContext(ctx, urltercero, &tercero_)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	} else {
		InfoTercero.Tercero = tercero_
	}

	// Consulta documento
	if documento_, err := GetDocumentoTercero(ctx, id); err != nil {
		return nil, err
	} else {
		if len(documento_) != 0 {
			InfoTercero.Identificacion = documento_[0]
		}
	}

	return InfoTercero, nil
}
