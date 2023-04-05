package utilsHelper

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
)

func ArrayToString(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

// findIdInArray Retorna la posicion en que se encuentra el id específicado
func FindIdInArray(idsList []*models.Elemento, id int) (i int) {
	for i, id_ := range idsList {
		if int(id_.Id) == id {
			return i
		}
	}
	return -1
}

// findElementoInArray Retorna la posicion en que se encuentra el id específicado
func FindElementoInArrayElementosMovimiento(elementos []*models.ElementosMovimiento, id int) (i int) {
	for i, el_ := range elementos {
		if int(*el_.ElementoActaId) == id {
			return i
		}
	}
	return -1
}

// fillElemento Agrega la vida útil y valor residual al elemento del acta
func FillElemento(elActa *models.DetalleElemento, elMov *models.ElementosMovimiento) (completo *models.DetalleElemento__, outputError map[string]interface{}) {

	funcion := "fillElemento"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	if err := formatdata.FillStruct(elActa, &completo); err != nil {
		logs.Error(err)
		eval := " - formatdata.FillStruct(elActa, &completo)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}
	completo.VidaUtil = elMov.VidaUtil
	completo.ValorResidual = elMov.ValorResidual

	return completo, nil

}

// RemoveDuplicateInt Remueve de un slice los int duplicados
func RemoveDuplicateInt(intSlice []int) []int {
	allKeys := make(map[int]bool)
	list := []int{}
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func EncodeUrl(query string, fields string, sortby string, order string, offset string, limit string) string {
	params := url.Values{}

	if len(query) > 0 {
		params.Add("query", query)
	}

	if len(fields) > 0 {
		params.Add("fields", fields)
	}

	if len(sortby) > 0 {
		params.Add("sortby", sortby)

	}

	if len(order) > 0 {
		params.Add("order", order)

	}

	if len(offset) > 0 {
		params.Add("offset", offset)

	}

	if len(limit) > 0 {
		params.Add("limit", limit)

	}

	return params.Encode()
}
