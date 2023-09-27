package autenticacion

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	"github.com/udistrital/arka_mid/utils_oas/request"
)

var basePath, _ = beego.AppConfig.String("autenticacionService")

// DataUsuario Consulta datos asociados a un usuario de la MID API de Autenticación
func DataUsuario(usuarioWSO2 string) (dataUsuario models.UsuarioAutenticacion, outputError map[string]interface{}) {

	funcion := "DataUsuario - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	url := "http://" + basePath + "token/userRol"
	req := models.UsuarioDataRequest{User: usuarioWSO2}
	// logs.Debug("url:", url, "- req:", req)
	if err := request.SendJson(url, "POST", &dataUsuario, &req); err == nil {
		return dataUsuario, nil
	} else {
		var empty models.UsuarioAutenticacion
		logs.Error(err)
		eval := `request.SendJson(url, "POST", &dataUsuario, &req)`
		return empty, errorCtrl.Error(funcion+eval, err, "500")
	}

}

// GetInfoUser Consulta los roles y el TerceroId asociado a un usuario determinado
func GetInfoUser(usr string, terceroId *int, roles *[]string) (outputError map[string]interface{}) {

	funcion := "GetInfoUser - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	user, err := DataUsuario(usr)
	if err != nil {
		return err
	}

	*roles = user.Role

	return GetTerceroUser(user, terceroId)
}

// GetTerceroUser Consulta los roles y el TerceroId asociado a un usuario determinado
func GetTerceroUser(user models.UsuarioAutenticacion, terceroId *int) (outputError map[string]interface{}) {

	funcion := "GetTerceroUser - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	if user.Documento == "" {
		return
	}

	payload := "documento=" + user.Documento

	tercero, outputError := terceros.GetAllTrTerceroIdentificacion(payload)
	if outputError != nil {
		return
	}

	if len(tercero) > 0 {
		*terceroId = tercero[0].Tercero.Id
	}

	return

}
