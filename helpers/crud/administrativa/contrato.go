package administrativa

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	administrativa "github.com/udistrital/administrativa_mid_api/models"
	e "github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

var basePath = beego.AppConfig.String("administrativaJbpmService")

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

	urlCrud := "http://" + basePath + fmt.Sprintf("informacion_contrato/%d/%s", contratoId, vigencia)
	err := request.GetJsonWSO2(urlCrud, &contrato)
	if err != nil {
		logs.Error(err, urlCrud)
		context := "request.GetJsonWSO2(urlCrud, &contrato)"
		outputError = e.Error(funcion+context, err, fmt.Sprint(http.StatusBadGateway))
	}

	return
}

// GetTipoContratoById Consulta endpoint tipo_contrato/:id del api administrativa_amazon_api
func GetTipoContratoById(tipoContratoId string, tipoContrato interface{}) (outputError map[string]interface{}) {

	const funcion = "GetTipoContratoById - "
	defer e.ErrorControlFunction(funcion+"Unhandled Error!", fmt.Sprint(http.StatusInternalServerError))

	urlcrud := "http://" + basePath + "tipo_contrato/" + tipoContratoId
	err := request.GetJsonWSO2(urlcrud, &tipoContrato)
	if err != nil {
		logs.Error(err, urlcrud)
		eval := "request.GetJson(urlcrud, &tipoContrato)"
		outputError = e.Error(funcion+eval, err, "502")
	}

	return
}
