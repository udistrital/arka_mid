package utilsHelper

import (
	"encoding/json"
)

// ConvertirInterfaceMap
func ConvertirInterfaceMap(Objeto interface{}) (Salida map[string]interface{}, err error) {

	if jsonString, err := json.Marshal(Objeto); err == nil {
		if err2 := json.Unmarshal(jsonString, &Salida); err2 != nil {
			panic(err.Error())
			return nil, err2
		}
	} else {
		panic(err.Error())
		return nil, err
	}
	return Salida, nil
}
// ConvertirInterfaceArrayMap
func ConvertirInterfaceArrayMap(Objeto_ interface{}) (Salida []map[string]interface{}, err error) {

	if jsonString, err := json.Marshal(Objeto_); err == nil {
		if err2 := json.Unmarshal(jsonString, &Salida); err2 != nil {
			panic(err.Error())
			return nil, err2
		}
	} else {
		panic(err.Error())
		return nil, err
	}
	return Salida, nil
}


// ArrayFind
func ArrayFind(Objeto__ []map[string]interface{}, campo string, valor string) ( Busqueda map[string]interface{}, err error) {

	Busqueda_ := make(map[string]interface{}, 0)
	if keys := len(Objeto__[0]); keys != 0 {
		
		for _, value := range Objeto__ {
			if value[campo] == valor {
				Busqueda_ = value;
				return Busqueda_, nil
			}
		}

	} else {
		panic(err.Error())
		return nil, err
	}
	
	return Busqueda_, nil
}

