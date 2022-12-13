package parametros

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

func GetUVTByVigencia(vigencia int) (uvt float64, outputError map[string]interface{}) {

	funcion := "GetUVTByVigencia - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	var parametros__ []models.ParametroPeriodo
	payload := "fields=Valor&limit=1&sortby=Id&order=desc&query=Activo:true,ParametroId__CodigoAbreviacion:UVT," +
		"PeriodoId__Nombre:" + strconv.Itoa(vigencia)
	if err := GetAllParametroPeriodo(payload, &parametros__); err != nil {
		return 0, err
	}

	if len(parametros__) == 1 {
		var valor map[string]interface{}
		if err := utilsHelper.Unmarshal(parametros__[0].Valor, &valor); err != nil {
			return 0, err
		}
		uvt = valor["Valor"].(float64)
	}

	return
}
