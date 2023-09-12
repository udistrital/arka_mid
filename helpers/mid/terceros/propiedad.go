package terceros

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	"github.com/udistrital/arka_mid/utils_oas/request"
)

var basePath, _ = beego.AppConfig.String("tercerosMidService")

func GetCargoFuncionario(id int) (cargo []*models.Parametro, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetCargoFuncionario", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	// Consulta cargo
	urlcrud := "http://" + basePath + "propiedad/cargo/" + strconv.Itoa(id)
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
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	// Consulta documento
	urlcrud := "http://" + basePath + "propiedad/documento/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &documento); err != nil {
		eval := " - request.GetJson(urlcrud, &documento)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}
