package administrativa

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/arka_mid/models"
	errorCtrl "github.com/udistrital/utils_oas/v2/errorctrl"
	requestV2 "github.com/udistrital/utils_oas/v2/request"
)

var basePath, _ = beego.AppConfig.String("administrativaJbpmService")

// GetContrato ...
func GetContrato(ctx context.Context, contratoId int, vigencia string, contrato *models.InformacionContrato) (outputError map[string]interface{}) {

	const funcion = "GetContrato - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", fmt.Sprint(http.StatusInternalServerError))

	if contratoId <= 0 {
		err := errors.New("id del contrato debe ser distinto de cero")
		context := "contratoId <= 0"
		logs.Error(err)
		return errorCtrl.Error(funcion+context, err, fmt.Sprint(http.StatusBadRequest))
	}

	urlCrud := basePath + fmt.Sprintf("informacion_contrato/%d/%s", contratoId, vigencia)
	_, err := requestV2.GetWithContext(ctx, urlCrud, contrato)
	if err != nil {
		logs.Error(err, urlCrud)
		eval := "requestV2.GetWithContext(ctx, urlCrud, contrato)"
		outputError = errorCtrl.Error(funcion+eval, err, fmt.Sprint(http.StatusBadGateway))
	}

	return
}

// GetTipoContratoById Consulta endpoint tipo_contrato/:id del api administrativa_amazon_api
func GetTipoContratoById(ctx context.Context, tipoContratoId string, tipoContrato interface{}) (outputError map[string]interface{}) {

	const funcion = "GetTipoContratoById - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", fmt.Sprint(http.StatusInternalServerError))

	urlcrud := basePath + "tipo_contrato/" + tipoContratoId
	_, err := requestV2.GetWithContext(ctx, urlcrud, tipoContrato)
	if err != nil {
		logs.Error(err, urlcrud)
		eval := "requestV2.GetWithContext(ctx, urlcrud, tipoContrato)"
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}
