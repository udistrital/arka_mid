package utilsHelper

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/logs"
)

// ConvertirInterfaceMap
func ConvertirInterfaceMap(Objeto interface{}) (Salida map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/ConvertirInterfaceMap - Unhandled Error!",
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
				"funcion": "/ConvertirInterfaceMap - json.Unmarshal(jsonString, &Salida)",
				"err":     err,
				"status":  "500",
			}
			return nil, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/ConvertirInterfaceMap - json.Marshal(Objeto)",
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
func ConvertirStringJson(Objeto_ interface{}) (Salida map[string]interface{}, err error) {

	str := fmt.Sprintf("%v", Objeto_)
	if err := json.Unmarshal([]byte(str), &Salida); err != nil {
		panic(err.Error())
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
