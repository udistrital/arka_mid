package actaRecibido

import (
	"fmt"

	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/autenticacion"
)

func filtrarActasSegunRoles(actas *[]map[string]interface{}, usuarioWSO2 string) (outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/filtrarActas",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	if data, err := autenticacion.DataUsuario(usuarioWSO2); err == nil {
		fmt.Print("\nDATA_USUARIO: ")
		fmt.Println(data)
		return nil
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/filtrarActas",
			"err":     err,
			"status":  "502",
		}
		return outputError
	}
}
