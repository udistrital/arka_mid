package tercerosMidHelper

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetDetalle Consulta El nombre, número de identificación, correo y cargo asociado a un funcionario
func GetDetalleFuncionario(id int) (DetalleFuncionario *models.DetalleFuncionario, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetDetalleFuncionario", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var (
		urlcrud  string
		response []*models.DetalleTercero
		cargo    []*models.Parametro
		correo   []*models.InfoComplementariaTercero
	)

	DetalleFuncionario = new(models.DetalleFuncionario)

	// Consulta información general y documento de identidad
	urlcrud = "http://" + beego.AppConfig.String("tercerosMidService") + "tipo/funcionarios/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &response); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetDetalleFuncionario - request.GetJson(urlcrud, &response1)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	DetalleFuncionario.Tercero = response

	// Consulta correo
	urlcrud = "http://" + beego.AppConfig.String("tercerosService") + "info_complementaria_tercero?limit=1&fields=Dato&sortby=Id&order=desc"
	urlcrud += "&query=Activo%3Atrue,InfoComplementariaId__Nombre__icontains%3Acorreo,TerceroId__Id%3A" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &correo); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetDetalleFuncionario - request.GetJson(urlcrud, &response2)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	DetalleFuncionario.Correo = correo

	// Consulta cargo
	urlcrud = "http://" + beego.AppConfig.String("tercerosMidService") + "propiedad/cargo/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &cargo); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetDetalleFuncionario - request.GetJson(urlcrud, &response3)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	DetalleFuncionario.Cargo = cargo

	return DetalleFuncionario, nil
}
