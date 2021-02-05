package tercerosHelper

import (
	"fmt"

	"github.com/astaxie/beego/logs"
)

func GetJefes(idTercero int) (terceros []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/GetJefes",
				"err":     err,
				"status":  "500", // Error no manejado!
			}
			panic(outputError)
		}
	}()

	// TODO: (Eliminar este comentario e) Implementar

	err := fmt.Errorf("No implementado (aún)")
	logs.Error(err)
	outputError = map[string]interface{}{
		"funcion": "/GetJefes",
		"err":     err,
		"status":  "501",
	}
	return nil, outputError
}
