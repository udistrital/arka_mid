package utilsHelper

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
)

// ConvertirInterfaceMap
func ConvertirInterfaceMap(Objeto interface{}) (Salida map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "ConvertirInterfaceMap - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	if jsonString, err := json.Marshal(Objeto); err == nil {
		if err2 := json.Unmarshal(jsonString, &Salida); err2 != nil {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "ConvertirInterfaceMap - json.Unmarshal(jsonString, &Salida)",
				"err":     err,
				"status":  "500",
			}
			return nil, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "ConvertirInterfaceMap - json.Marshal(Objeto)",
			"err":     err,
			"status":  "500",
		}
		return nil, outputError
	}
	return Salida, nil
}

// ConvertirInterfaceArrayMap
func ConvertirInterfaceArrayMap(Objeto_ interface{}) (Salida []map[string]interface{}, err error) {
	fmt.Println(Objeto_)
	if jsonString, err := json.Marshal(Objeto_); err == nil {
		if err2 := json.Unmarshal(jsonString, &Salida); err2 != nil {
			panic(err.Error())
		}
	} else {
		panic(err.Error())
	}
	return Salida, nil
}

// ConvertirStringJson
func ConvertirStringJson(Objeto_ interface{}) (Salida map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "ConvertirStringJson - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	str := fmt.Sprintf("%v", Objeto_)
	if err := json.Unmarshal([]byte(str), &Salida); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "ConvertirStringJson - json.Unmarshal([]byte(str), &Salida)",
			"err":     err,
			"status":  "500",
		}
		return nil, outputError
	}
	return Salida, nil

}

// ArrayFind
func ArrayFind(Objeto__ []map[string]interface{}, campo string, valor string) (Busqueda map[string]interface{}, err error) {

	if len(Objeto__) == 0 {
		return nil, nil
	}

	Busqueda_ := make(map[string]interface{}, 0)
	if keys := len(Objeto__[0]); keys != 0 {

		for _, value := range Objeto__ {
			if value[campo] == valor {
				Busqueda_ = value
				return Busqueda_, nil
			}
		}

	} else {
		panic(err.Error())
	}

	return Busqueda_, nil
}

// KeysValuesMap descompone un mapeo en dos arreglos con sus claves y valores
func KeysValuesMap(m map[interface{}]interface{}) (keys []interface{}, vals []interface{}) {

	defer func() {
		if err := recover(); err != nil {
			panic(map[string]interface{}{
				"funcion": "KeysValuesMap - Unhandled Error!",
				"err":     err,
				"status":  "500",
			})
		}
	}()

	for k, v := range m {
		keys = append(keys, k)
		vals = append(vals, v)
	}
	return
}

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

// removeDuplicateIds Remueve de un vector los enteros duplicados
func RemoveDuplicateIds(addrs []int) []int {
	result := make([]int, 0, len(addrs))
	temp := map[int]struct{}{}
	for _, item := range addrs {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
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
