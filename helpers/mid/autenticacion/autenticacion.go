package autenticacion

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

var basePath = "http://" + beego.AppConfig.String("autenticacionService")

// DataUsuario Consulta datos asociados a un usuario de la MID API de Autenticación
func DataUsuario(usuarioWSO2 string) (dataUsuario models.UsuarioAutenticacion, outputError map[string]interface{}) {

	funcion := "DataUsuario - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	url := basePath + "token/userRol"
	req := models.UsuarioDataRequest{User: usuarioWSO2}
	// logs.Debug("url:", url, "- req:", req)
	if err := request.SendJson(url, "POST", &dataUsuario, &req); err == nil {
		return dataUsuario, nil
	} else {
		var empty models.UsuarioAutenticacion
		logs.Error(err)
		eval := `request.SendJson(url, "POST", &dataUsuario, &req)`
		return empty, errorctrl.Error(funcion+eval, err, "500")
	}

}
