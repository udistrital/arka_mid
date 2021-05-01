package utilsHelper

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego/logs"

	"github.com/udistrital/utils_oas/request"
)

// BufferGeneric actúa como proxy para evitar consultas repetidas a un helper
// o callback, que podría mediar entre el helper
//
// Para usarlo, hay que haber creado previamente el diccionario en que se irán
// almacenando los resultados:
//
// diccionario := make(map[int]interface{})
//
// También requiere punteros a un par de enteros (previamente inicializados) que
// registrarán las consultas realizadas y evitadas
func BufferGeneric(id int, diccionario map[int]interface{}, callback func() (interface{}, map[string]interface{}),
	consultasNecesarias *int, consultasEvitadas *int) (elemento interface{}, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			panic(map[string]interface{}{
				"funcion": "BufferGeneric - Unhandled Error!",
				"err":     err,
				"status":  "500",
			})
		}
	}()

	if elemento, ok := diccionario[id]; ok {
		if consultasEvitadas != nil {
			*consultasEvitadas++
		}
		return elemento, nil
	}

	if consultasNecesarias != nil {
		*consultasNecesarias++
	}

	if res, err := callback(); err == nil {
		if res != nil {
			diccionario[id] = res
		}
		return res, nil
	} else {
		return nil, err
	}
}

// BufferGetStat actúa como proxy para evitar consultas repetidas a una URL determinada
// siempre y cuando la URL y id especificados tengan una relación directa y bidireccional
//
// Es un caso particular de BufferGeneric
func BufferGetStat(id int, diccionario map[int]interface{}, url string,
	consultasNecesarias *int, consultasEvitadas *int) (elemento map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			panic(map[string]interface{}{
				"funcion": "BufferStat - Unhandled Error!",
				"err":     err,
				"status":  "500",
			})
		}
	}()

	makeRequest := func() (interface{}, map[string]interface{}) {
		var data []map[string]interface{}
		if res, err := request.GetJsonTest(url, &data); err == nil && res.StatusCode == 200 {
			if len(data) == 1 && len(data[0]) > 0 {
				diccionario[id] = data[0]
				return data[0], nil
			}
			return nil, nil
		} else {
			if err == nil {
				err = fmt.Errorf("undesired Status Code: %d != 200", res.StatusCode)
			}
			logs.Error(err)
			return nil, map[string]interface{}{
				"funcion": "BufferGet/makeRequest - request.GetJsonTest(url, &data)",
				"err":     err,
				"status":  "502",
			}
		}
	}

	if v, err := BufferGeneric(id, diccionario, makeRequest, consultasNecesarias, consultasEvitadas); err == nil {
		if v2, ok := v.(map[string]interface{}); ok {
			return v2, nil
		} else {
			err := errors.New("no se pudo convertir interface{} a map[string]interface{}")
			logs.Error(err)
			return nil, map[string]interface{}{
				"funcion": "BufferGet - v.(map[string]interface{})",
				"err":     err,
				"status":  "502",
			}
		}
	} else {
		return nil, err
	}
}

// BufferGet es igual que BufferGetStat pero sin estadísticas
func BufferGet(id int, diccionario map[int]interface{}, url string) (elemento map[string]interface{}, outputError map[string]interface{}) {
	return BufferGetStat(id, diccionario, url, nil, nil)
}
