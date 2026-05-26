package actaRecibido

import (
	"net/url"

	"github.com/beego/beego/v2/core/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	"github.com/udistrital/arka_mid/utils_oas/request"
)

// GetAllEstadoActa consulta controlador estado_acta del api acta_recibido_crud
func GetAllEstadoActa(query string) (estados []models.EstadoActa, outputError map[string]interface{}) {
	funcion := "GetAllEstadoActa - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := path + "estado_acta?" + query
	if err := request.GetJson(urlcrud, &estados); err != nil {
		logs.Error(urlcrud+", ", err)
		eval := "request.GetJson(urlcrud, &estados)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetEstadoActaIdByCodigoAbreviacion consulta el Id de un EstadoActa según su código abreviación
func GetEstadoActaIdByCodigoAbreviacion(id *int, codigo string) (outputError map[string]interface{}) {
	funcion := "GetEstadoActaIdByCodigoAbreviacion - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	query := "query=CodigoAbreviacion:" + url.QueryEscape(codigo)
	if estados, err := GetAllEstadoActa(query); err != nil {
		return err
	} else if len(estados) == 0 {
		errMsg := "No se encuentra el estado acta: " + codigo
		logs.Error(errMsg)
		eval := "GetAllEstadoActa(query)"
		return errorCtrl.Error(funcion+eval, errMsg, "500")
	} else {
		*id = estados[0].Id
	}

	return
}
