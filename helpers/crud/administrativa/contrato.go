package administrativa

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/utils_oas/request"
)

// GetContrato ...
func GetContrato(contratoId int, vigencia string) (contrato map[string]interface{}, outputError map[string]interface{}) {
	if contratoId != 0 { // (1) error parametro
		request.GetJsonWSO2("http://"+beego.AppConfig.String("administrativaJbpmService")+"informacion_contrato/"+strconv.Itoa(int(contratoId))+"/"+vigencia, &contrato)
		return contrato, nil

	} else {
		logs.Info("Error (1) Parametro")
		outputError = map[string]interface{}{"Function": "FuncionalidadMidController:GetContrato", "Error": "null parameter"}
		return nil, outputError
	}
}
