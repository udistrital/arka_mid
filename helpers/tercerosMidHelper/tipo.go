package tercerosMidHelper

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

func GetFuncionario(id int) (funcionario []*models.DetalleTercero, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetFuncionario", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	// Consulta información general y documento de identidad
	urlcrud := "http://" + beego.AppConfig.String("tercerosMidService") + "tipo/funcionarios/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &funcionario); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetFuncionario - request.GetJson(urlcrud, &funcionario)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	return
}
