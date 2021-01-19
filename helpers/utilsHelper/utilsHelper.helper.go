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
				"funcion": "/ConvertirInterfaceMap",
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
				"funcion": "/ConvertirInterfaceMap",
				"err":     err,
				"status":  "500",
			}
			return nil, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/ConvertirInterfaceMap",
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
