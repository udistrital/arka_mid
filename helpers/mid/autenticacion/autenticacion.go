package autenticacion

import (
	"fmt"
	"regexp"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
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

// GetInfoUser Consulta los roles y el TerceroId asociado a un usuario determinado
func GetInfoUser(usr string, terceroId *int, roles *[]string) (outputError map[string]interface{}) {

	funcion := "GetInfoUser - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

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
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	if user.Documento == "" {
		return
	}

	rgxp := regexp.MustCompile("\\d.*")
	tipo := rgxp.ReplaceAllString(user.DocumentoCompuesto, "")

	if tipo != "" {
		payload := "query=Activo:true,TerceroId__Activo:true,TipoDocumentoId__CodigoAbreviacion:" + tipo + ",Numero:" + user.Documento
		datosId, err := terceros.GetAllDatosIdentificacion(payload)

		if err != nil {
			return err
		}

		if len(datosId) == 1 && datosId[0].TerceroId != nil {
			*terceroId = datosId[0].TerceroId.Id
			return
		} else if len(datosId) > 1 {
			if terceros.DocumentosValidos(datosId, false, true) {
				*terceroId = datosId[0].TerceroId.Id
				return
			}

			err := fmt.Errorf("el Documento '%s' tiene más de un registro activo en Terceros (%d registros).", user.DocumentoCompuesto, len(datosId))
			logs.Notice(err)
			outputError = errorctrl.Error(funcion, err, "409")

			return outputError
		}
	}

	tercero, err := terceros.GetTerceroByDoc(user.Documento)
	if err != nil {
		return err
	}

	if tercero.TerceroId != nil && tercero.TerceroId.Id > 0 {
		*terceroId = tercero.TerceroId.Id
	}

	return

}
