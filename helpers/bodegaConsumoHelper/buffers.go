package bodegaConsumoHelper

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/utils_oas/request"
)

func bufferElementoActa(elementoID int, elementos map[int](map[string]interface{})) (elemento map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "bufferElementoActa - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	if elemento, ok := elementos[elementoID]; ok {
		return elemento, nil
	}

	idStr := strconv.Itoa(elementoID)

	var detalle []map[string]interface{}
	url2 := "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento"
	url2 += "?query=Id:" + idStr
	url2 += "&fields=Id,Nombre,Marca,Serie,SubgrupoCatalogoId"
	// logs.Debug("url2:", url2)
	if res, err := request.GetJsonTest(url2, &detalle); err == nil && res.StatusCode == 200 {
		if len(detalle) == 1 && len(detalle[0]) > 0 {
			elementos[elementoID] = detalle[0]
			return detalle[0], nil
		}
	} else {
		if err == nil {
			err = fmt.Errorf("Undesired Status Code: %d", res.StatusCode)
		}
		logs.Error(err)
		return nil, map[string]interface{}{
			"funcion": "bufferElementoActa - request.GetJsonTest(url2, &detalle)",
			"err":     err,
			"status":  "502",
		}
	}
	return nil, nil
}

func bufferSubgrupoCatalogo(subgrupoID int, subgrupos map[int](map[string]interface{})) (subgrupo map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "bufferSubgrupoCatalogo - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	if subgrupo, ok := subgrupos[subgrupoID]; ok {
		return subgrupo, nil
	}

	idStr := strconv.Itoa(subgrupoID)

	var detalle []map[string]interface{}

	url3 := "http://" + beego.AppConfig.String("catalogoElementosService") + "subgrupo"
	url3 += "?query=Id:" + idStr
	// logs.Debug("url3:", url3)
	if res, err := request.GetJsonTest(url3, &detalle); err == nil && res.StatusCode == 200 {
		if len(detalle) == 1 && len(detalle[0]) > 0 {
			subgrupos[subgrupoID] = detalle[0]
			return detalle[0], nil
		}
	} else {
		if err == nil {
			err = fmt.Errorf("Undesired Status Code: %d", res.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "bufferSubgrupoCatalogo - request.GetJsonTest(url3, &detalle)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	return nil, nil
}
