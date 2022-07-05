package administrativa

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/utils_oas/errorctrl"
	e "github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

// GetContrato ...
func GetContrato(contratoId int, vigencia string) (contrato map[string]interface{}, outputError map[string]interface{}) {
	const funcion = "GetContrato"
	defer e.ErrorControlFunction(funcion+" - Unhandled Error!", fmt.Sprint(http.StatusInternalServerError))

	if contratoId != 0 { // (1) error parametro
		urlContrato := "http://" + beego.AppConfig.String("administrativaJbpmService")
		urlContrato += fmt.Sprintf("informacion_contrato/%d/%s", contratoId, vigencia)
		if err := request.GetJsonWSO2(urlContrato, &contrato); err != nil {
			logs.Error(err)
			context := " - request.GetJsonWSO2(urlContrato, &contrato)"
			return nil, e.Error(funcion+context, err, fmt.Sprint(http.StatusBadGateway))
		}
		return contrato, nil

	} else {
		err := errors.New("id del contrato debe ser distinto de cero")
		context := " - contratoId != 0"
		logs.Error(err)
		return nil, e.Error(funcion+context, err, fmt.Sprint(http.StatusBadRequest))
	}
}

// GetTipoContratoById Consulta endpoint tipo_contrato/:id del api administrativa_amazon_api
func GetTipoContratoById(tipoContratoId int, tipoContrato interface{}) (outputError map[string]interface{}) {

	const funcion = "GetTipoContratoById - "
	defer e.ErrorControlFunction(funcion+"Unhandled Error!", fmt.Sprint(http.StatusInternalServerError))

	urlcrud := "http://" + beego.AppConfig.String("administrativaService") + "tipo_contrato/" + strconv.Itoa(tipoContratoId)
	if err := request.GetJson(urlcrud, &tipoContrato); err != nil {
		logs.Error(err, urlcrud)
		eval := "request.GetJson(urlcrud, &tipoContrato)"
		return errorctrl.Error(funcion+eval, err, "502")
	}

	return
}
