package tercerosHelper

import (
	"fmt"

	"github.com/astaxie/beego/logs"
)

var diccionarioTipoHelper = map[string](func(int) ([]map[string]interface{}, map[string]interface{})){
	"funcionarioPlanta": GetFuncionariosPlanta,
	"ordenadoresGasto":  GetOrdenadores,
	"contratista":       GetContratista,
	"proveedor":         GetProveedor,
}

// GetTipos retorna la lista de tipos que pueden ser usados con GetHelperTipo
func GetTipos() (tercero []string, outputError map[string]interface{}) {

	// Puede que ni sea necesario en este helper, pero se coloca por lineamiento...
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetTipos - Unhandled Error!",
				"err":     err,
				"status":  "500", // Error no manejado!
			}
			panic(outputError)
		}
	}()

	for k := range diccionarioTipoHelper {
		tercero = append(tercero, k)
	}
	return tercero, nil
}

// GetHelperTipo trae los terceros con el criterio especificado.
// El criterio debe ser alguno de los valores retornados por GetTipos
func GetHelperTipo(tipo string) (helper func(int) ([]map[string]interface{}, map[string]interface{}), outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetHelperTipo - Unhandled Error!",
				"err":     err,
				"status":  "500", // Error no manejado!
			}
			panic(outputError)
		}
	}()

	if helper, found := diccionarioTipoHelper[tipo]; found {
		return helper, nil
	}

	err := fmt.Errorf("\"%s\" not implemented", tipo)
	logs.Error(err)

	return nil, map[string]interface{}{
		"funcion": "GetHelperTipo - found := diccionarioTipoHelper[tipo]",
		"err":     err,
		"status":  "404",
	}
}
