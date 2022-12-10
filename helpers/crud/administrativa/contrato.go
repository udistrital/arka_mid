package administrativa

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	administrativa "github.com/udistrital/administrativa_mid_api/models"
	"github.com/udistrital/utils_oas/errorctrl"
	e "github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

// GetContrato ...
func GetContrato(contratoId int, vigencia string, contrato *administrativa.InformacionContrato) (outputError map[string]interface{}) {

	const funcion = "GetContrato - "
	defer e.ErrorControlFunction(funcion+"Unhandled Error!", fmt.Sprint(http.StatusInternalServerError))

	if contratoId <= 0 {
		err := errors.New("id del contrato debe ser distinto de cero")
		context := "contratoId <= 0"
		logs.Error(err)
		return e.Error(funcion+context, err, fmt.Sprint(http.StatusBadRequest))
	}

	urlCrud := "http://" + beego.AppConfig.String("administrativaJbpmService") +
		fmt.Sprintf("informacion_contrato/%d/%s", contratoId, vigencia)
	if err := request.GetJsonWSO2(urlCrud, &contrato); err != nil {
		logs.Error(err)
		context := "request.GetJsonWSO2(urlCrud, &contrato)"
		return e.Error(funcion+context, err, fmt.Sprint(http.StatusBadGateway))
	}

	return
}

// GetTipoContratoById Consulta endpoint tipo_contrato/:id del api administrativa_amazon_api
func GetTipoContratoById(tipoContratoId string, tipoContrato interface{}) (outputError map[string]interface{}) {

	const funcion = "GetTipoContratoById - "
	defer e.ErrorControlFunction(funcion+"Unhandled Error!", fmt.Sprint(http.StatusInternalServerError))

	urlcrud := "http://" + beego.AppConfig.String("administrativaService") + "tipo_contrato/" + tipoContratoId
	if err := request.GetJson(urlcrud, &tipoContrato); err != nil {
		logs.Error(err, urlcrud)
		eval := "request.GetJson(urlcrud, &tipoContrato)"
		return errorctrl.Error(funcion+eval, err, "502")
	}

	return
}
