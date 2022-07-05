package configuracion

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

var basePath = beego.AppConfig.String("configuracionService")

func GetAllPerfilXMenuOpcion(query string, opciones *[]*models.PerfilXMenuOpcion) (outputError map[string]interface{}) {

	funcion := "GetAllPerfilXMenuOpcion - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := basePath + "perfil_x_menu_opcion?" + query
	if err := request.GetJson(urlcrud, &opciones); err != nil {
		eval := "request.GetJson(urlcrud, &opciones)"
		return errorctrl.Error(funcion+eval, err, "502")
	}

	return
}
