package autenticacion

import (
	"context"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	requestV2 "github.com/udistrital/utils_oas/v2/request"
)

var basePath, _ = beego.AppConfig.String("autenticacionService")

// DataUsuario Consulta datos asociados a un usuario de la MID API de Autenticación
func DataUsuario(ctx context.Context, usuarioWSO2 string) (dataUsuario models.UsuarioAutenticacion, outputError map[string]interface{}) {

	funcion := "DataUsuario - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	url := basePath + "token/userRol"
	req := models.UsuarioDataRequest{User: usuarioWSO2}
	// logs.Debug("url:", url, "- req:", req)
	if _, err := requestV2.PostWithContext(ctx, url, &req, &dataUsuario); err == nil {
		return dataUsuario, nil
	} else {
		var empty models.UsuarioAutenticacion
		logs.Error(err)
		eval := `requestV2.PostWithContext(ctx, url, &req, &dataUsuario)`
		return empty, errorCtrl.Error(funcion+eval, err, "500")
	}

}

// GetInfoUser Consulta los roles y el TerceroId asociado a un usuario determinado
func GetInfoUser(ctx context.Context, usr string, terceroId *int, roles *[]string) (outputError map[string]interface{}) {

	funcion := "GetInfoUser - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	user, err := DataUsuario(ctx, usr)
	if err != nil {
		return err
	}

	*roles = user.Role

	return GetTerceroUser(ctx, user, terceroId)
}

// GetTerceroUser Consulta los roles y el TerceroId asociado a un usuario determinado
func GetTerceroUser(ctx context.Context, user models.UsuarioAutenticacion, terceroId *int) (outputError map[string]interface{}) {

	funcion := "GetTerceroUser - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	if user.Documento == "" {
		return
	}

	payload := "documento=" + user.Documento

	tercero, outputError := terceros.GetAllTrTerceroIdentificacion(ctx, payload)
	if outputError != nil {
		return
	}

	if len(tercero) > 0 {
		*terceroId = tercero[0].Tercero.Id
	}

	return

}
