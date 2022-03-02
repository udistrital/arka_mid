package parametrosHelper

import (
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/utils_oas/errorctrl"
)

// GetParametrosDebitoCredito consulta los parametros de movimientos credito y debito
func GetParametrosDebitoCredito() (dbId, crId int, outputError map[string]interface{}) {

	funcion := "GetParametrosDebitoCredito"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var query string

	query = "query=CodigoAbreviacion:MCD"
	if par_, err := parametros.GetAllParametro(query); err != nil {
		return 0, 0, err
	} else {
		dbId = par_[0].Id
	}

	query = "query=CodigoAbreviacion:MCC"
	if par_, err := parametros.GetAllParametro(query); err != nil {
		return 0, 0, err
	} else {
		crId = par_[0].Id
	}

	return dbId, crId, nil
}
