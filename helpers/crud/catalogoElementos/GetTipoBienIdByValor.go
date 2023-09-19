package catalogoElementos

import (
	"fmt"
	"strconv"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// GetTipoBienIdByValor Determina el tipo bien al que pertenece un elemento dado el valor en UVT
func GetTipoBienIdByValor(tbPadreId int, normalizado float64, bufferTiposBien map[int]models.TipoBien) (tipoBienId int, outputError map[string]interface{}) {

	funcion := "GetTipoBienIdByValor - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	if tbPadreId <= 0 {
		return
	}

	for _, tb_ := range bufferTiposBien {
		if tb_.TipoBienPadreId.Id == tbPadreId && tb_.LimiteInferior <= normalizado && normalizado < tb_.LimiteSuperior {
			return tb_.Id, nil
		}
	}

	var tb__ []models.TipoBien
	payload := "limit=1&query=Activo:true,TipoBienPadreId__Id:" + strconv.Itoa(tbPadreId) + ",LimiteInferior__lte:" +
		fmt.Sprintf("%f", normalizado) + ",LimiteSuperior__gt:" + fmt.Sprintf("%f", normalizado)
	if err := GetAllTipoBien(payload, &tb__); err != nil || len(tb__) != 1 {
		return 0, err
	}

	bufferTiposBien[tb__[0].Id] = tb__[0]
	return tb__[0].Id, nil
}
