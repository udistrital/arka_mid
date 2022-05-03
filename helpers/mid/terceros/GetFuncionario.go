package terceros

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

var tercerosMID = beego.AppConfig.String("tercerosMidService")

func GetFuncionario(id int) (funcionario []*models.DetalleTercero, outputError map[string]interface{}) {

	funcion := "GetFuncionario"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	// Consulta información general y documento de identidad
	urlcrud := "http://" + tercerosMID + "tipo/funcionarios/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &funcionario); err != nil {
		logs.Error(urlcrud + ", " + err.Error())
		eval := " - request.GetJson(urlcrud, &funcionario)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}

	return
}
