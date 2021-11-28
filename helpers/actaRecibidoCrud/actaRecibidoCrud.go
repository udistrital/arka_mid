package actarecibidocrud

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

func GetElementoById(id int) (elemento *models.Elemento, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetElementoById", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	// Se consulta el elemento
	urlcrud := "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &elemento); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetElementoById - request.GetJson(urlcrud, &elemento)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	return elemento, nil
}
