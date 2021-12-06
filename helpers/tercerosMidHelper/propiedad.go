package tercerosMidHelper

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

func GetCargoFuncionario(id int) (cargo []*models.Parametro, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetCargoFuncionario", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	// Consulta cargo
	urlcrud := "http://" + beego.AppConfig.String("tercerosMidService") + "propiedad/cargo/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &cargo); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetCargoFuncionario - request.GetJson(urlcrud, &cargo)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	return
}

// GetDocumentoTercero get controlador propiedad/documento/{id} del api terceros_mid
func GetDocumentoTercero(id int) (documento []*models.DatosIdentificacion, outputError map[string]interface{}) {

	funcion := "GetDocumentoTercero"
	defer errorctrl.ErrorControlFunction(funcion, "500")

	// Consulta documento
	urlcrud := "http://" + beego.AppConfig.String("tercerosMidService") + "propiedad/documento/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &documento); err != nil {
		funcion += " - request.GetJson(urlcrud, &documento)"
		return nil, errorctrl.Error(funcion, err, "502")
	}

	return
}
