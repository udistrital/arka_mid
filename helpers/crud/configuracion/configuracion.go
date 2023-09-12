package configuracion

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	"github.com/udistrital/arka_mid/utils_oas/request"
)

var basePath, _ = beego.AppConfig.String("configuracionService")

func GetAllPerfilXMenuOpcion(query string, opciones *[]*models.PerfilXMenuOpcion) (outputError map[string]interface{}) {

	funcion := "GetAllPerfilXMenuOpcion - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := basePath + "perfil_x_menu_opcion?" + query
	if err := request.GetJson(urlcrud, &opciones); err != nil {
		eval := "request.GetJson(urlcrud, &opciones)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

func GetAllParametro(query string, parametros *[]models.ParametroConfiguracion) (outputError map[string]interface{}) {

	funcion := "GetAllParametro - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := basePath + "parametro?query=Aplicacion__Nombre:arka_ii_main," + query
	if err := request.GetJson(urlcrud, &parametros); err != nil {
		logs.Error(err, urlcrud)
		eval := "request.GetJson(urlcrud, &opciones)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

func PutParametro(id int, parametro *models.ParametroConfiguracion) (outputError map[string]interface{}) {

	funcion := "PutParametro - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := basePath + "parametro/" + strconv.Itoa(id)
	if err := request.SendJson(urlcrud, "PUT", &parametro, &parametro); err != nil {
		logs.Error(err, urlcrud)
		eval := `request.SendJson(urlcrud, "PUT", &parametro, &parametro)`
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}
