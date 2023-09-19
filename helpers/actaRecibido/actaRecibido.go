package actaRecibido

import (
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	"github.com/udistrital/arka_mid/utils_oas/request"
)

// GetAllParametrosActa Consulta diferentes valores paramétricos
func GetAllParametrosActa() (parametros_ map[string]interface{}, outputError map[string]interface{}) {

	funcion := "GetAllParametrosActa - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	var (
		EstadoActa     interface{}
		EstadoElemento interface{}
		Ivas           = make([]models.Iva, 0)
	)

	var path, _ = beego.AppConfig.String("actaRecibidoService")
	path = "http://" + path
	urlActasEstadoActa := path + "estado_acta?limit=-1"
	if _, err := request.GetJsonTest(urlActasEstadoActa, &EstadoActa); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllParametrosActa - request.GetJsonTest(urlActasEstadoActa, &EstadoActa)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	urlACtasEstadoElem := path + "estado_elemento?limit=-1"
	if _, err := request.GetJsonTest(urlACtasEstadoElem, &EstadoElemento); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllParametrosActa - request.GetJsonTest(urlACtasEstadoElem, &EstadoElemento)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	if err := parametros.GetAllIVAByPeriodo(strconv.Itoa(time.Now().Year()), &Ivas); err != nil {
		return nil, err
	}

	parametros_ = map[string]interface{}{
		"EstadoActa":     EstadoActa,
		"EstadoElemento": EstadoElemento,
		"IVA":            Ivas,
	}

	return parametros_, nil
}
