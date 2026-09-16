package configuracion

import (
	"context"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/models"
	errorCtrl "github.com/udistrital/utils_oas/v2/errorctrl"
	requestV2 "github.com/udistrital/utils_oas/v2/request"
)

var basePath, _ = beego.AppConfig.String("configuracionService")

func GetAllPerfilXMenuOpcion(ctx context.Context, query string, opciones *[]*models.PerfilXMenuOpcion) (outputError map[string]interface{}) {

	funcion := "GetAllPerfilXMenuOpcion - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := basePath + "perfil_x_menu_opcion?" + query
	if _, err := requestV2.GetWithContext(ctx, urlcrud, opciones); err != nil {
		eval := "requestV2.GetWithContext(ctx, urlcrud, opciones)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

func GetAllParametro(ctx context.Context, query string, parametros *[]models.ParametroConfiguracion) (outputError map[string]interface{}) {

	funcion := "GetAllParametro - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := basePath + "parametro?query=Aplicacion__Nombre:arka_ii_main," + query
	if _, err := requestV2.GetWithContext(ctx, urlcrud, parametros); err != nil {
		logs.Error(err, urlcrud)
		eval := "requestV2.GetWithContext(ctx, urlcrud, parametros)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

func PutParametro(ctx context.Context, id int, parametro *models.ParametroConfiguracion) (outputError map[string]interface{}) {

	funcion := "PutParametro - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := basePath + "parametro/" + strconv.Itoa(id)
	if _, err := requestV2.PutWithContext(ctx, urlcrud, parametro, parametro); err != nil {
		logs.Error(err, urlcrud)
		eval := "requestV2.PutWithContext(ctx, urlcrud, parametro, parametro)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}
