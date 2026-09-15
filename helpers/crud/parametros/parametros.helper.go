package parametros

import (
	"context"

	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// GetParametrosDebitoCredito consulta los parametros de movimientos credito y debito
func GetParametrosDebitoCredito(ctx context.Context) (dbId, crId int, outputError map[string]interface{}) {

	funcion := "GetParametrosDebitoCredito"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var query string

	query = "query=CodigoAbreviacion:MCD,Activo:true"
	if par_, err := GetAllParametro(ctx, query); err != nil {
		return 0, 0, err
	} else {
		dbId = par_[0].Id
	}

	query = "query=CodigoAbreviacion:MCC,Activo:true"
	if par_, err := GetAllParametro(ctx, query); err != nil {
		return 0, 0, err
	} else {
		crId = par_[0].Id
	}

	return dbId, crId, nil
}
