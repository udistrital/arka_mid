package terceros

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	"github.com/udistrital/arka_mid/utils_oas/request"
)

var tercerosMID, _ = beego.AppConfig.String("tercerosMidService")

func GetFuncionario(id int) (funcionario []*models.DetalleTercero, outputError map[string]interface{}) {

	funcion := "GetFuncionario"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	// Consulta información general y documento de identidad
	urlcrud := "http://" + tercerosMID + "tipo/funcionarios/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &funcionario); err != nil {
		logs.Error(urlcrud + ", " + err.Error())
		eval := " - request.GetJson(urlcrud, &funcionario)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}
