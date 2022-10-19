package parametros

import (
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

func GetAllIVAByPeriodo(vigencia string, ivas *[]models.Iva) (outputError map[string]interface{}) {

	funcion := "GetAllIVAByPeriodo - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	var parametros__ []models.ParametroPeriodo

	payload := "query=ParametroId__TipoParametroId__CodigoAbreviacion:I,PeriodoId__Nombre:" + vigencia
	if err := GetAllParametroPeriodo(payload, &parametros__); err != nil {
		return err
	}

	for _, par := range parametros__ {
		var iva_ models.Iva
		if err := utilsHelper.Unmarshal(par.Valor, &iva_); err != nil {
			return err
		}

		iva_.Id = par.Id
		iva_.CodigoAbreviacion = par.ParametroId.CodigoAbreviacion
		*ivas = append(*ivas, iva_)
	}

	return
}
