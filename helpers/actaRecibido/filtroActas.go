package actaRecibido

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

func filtrarActas(actas []map[string]interface{}, usuarioWSO2 string) (filtradas []map[string]interface{}, outputError map[string]interface{}) {

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

	url := "http://" + beego.AppConfig.String("autenticacionService") + "token/userRol"

	var dataUsuario models.UsuarioAutenticacion
	req := models.UsuarioDataRequest{User: usuarioWSO2}
	fmt.Print("\nREQUEST_USUARIO: ")
	fmt.Println(req)

	if err := request.SendJson(url, "POST", &dataUsuario, &req); err == nil {
		fmt.Print("\nDATA_USUARIO: ")
		fmt.Println(dataUsuario)
		return actas, nil
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/filtrarActas",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}
