package tercerosHelper

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetCorreo Consulta el correo de un tercero
func GetCorreo(id int) (DetalleFuncionario []*models.InfoComplementariaTercero, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetCorreo", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var (
		urlcrud string
		correo  []*models.InfoComplementariaTercero
	)

	// Consulta correo
	urlcrud = "http://" + beego.AppConfig.String("tercerosService") + "info_complementaria_tercero?limit=1&fields=Dato&sortby=Id&order=desc"
	urlcrud += "&query=Activo%3Atrue,InfoComplementariaId__Nombre__icontains%3Acorreo,TerceroId__Id%3A" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &correo); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetCorreo - request.GetJson(urlcrud, &response2)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	return correo, nil
}
