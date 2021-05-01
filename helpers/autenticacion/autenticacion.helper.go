package autenticacion

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// DataUsuario Consulta datos asociados a un usuario de la MID API de Autenticación
func DataUsuario(usuarioWSO2 string) (dataUsuario models.UsuarioAutenticacion, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "DataUsuario - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	url := "http://" + beego.AppConfig.String("autenticacionService") + "token/userRol"

	req := models.UsuarioDataRequest{User: usuarioWSO2}

	if err := request.SendJson(url, "POST", &dataUsuario, &req); err == nil {
		return dataUsuario, nil
	} else {
		var empty models.UsuarioAutenticacion
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "DataUsuario - request.SendJson(url, \"POST\", &dataUsuario, &req)",
			"err":     err,
			"status":  "502",
		}
		return empty, outputError
	}

}
