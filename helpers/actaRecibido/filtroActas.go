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

func filtrarActasPorEstados(actas []map[string]interface{}, states []string) []map[string]interface{} {
	fin := len(actas)
	for i := 0; i < fin; {
		dejar := false
		for _, reqState := range states {
			if actas[i]["Estado"] == reqState {
				dejar = true
				break
			}
		}
		if dejar {
			i++
		} else {
			actas[i] = actas[fin-1]
			fin--
		}
	}
	actas = actas[:fin]

	return actas
}
