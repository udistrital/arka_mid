package actaRecibido

import (
	"context"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/arka_mid/models"
	errorCtrl "github.com/udistrital/utils_oas/v2/errorctrl"
	requestV2 "github.com/udistrital/utils_oas/v2/request"
)

// GetAllParametrosActa Consulta diferentes valores paramétricos
func GetAllParametrosActa(ctx context.Context) (parametros_ map[string]interface{}, outputError map[string]interface{}) {

	funcion := "GetAllParametrosActa - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	var (
		EstadoActa     interface{}
		EstadoElemento interface{}
		Ivas           = make([]models.Iva, 0)
	)

	var path, _ = beego.AppConfig.String("actaRecibidoService")
	urlActasEstadoActa := path + "estado_acta?limit=-1"
	if _, err := requestV2.GetWithContext(ctx, urlActasEstadoActa, &EstadoActa); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllParametrosActa - requestV2.GetWithContext(ctx, urlActasEstadoActa, &EstadoActa)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	urlACtasEstadoElem := path + "estado_elemento?limit=-1"
	if _, err := requestV2.GetWithContext(ctx, urlACtasEstadoElem, &EstadoElemento); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllParametrosActa - requestV2.GetWithContext(ctx, urlACtasEstadoElem, &EstadoElemento)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	if err := parametros.GetAllIVAByPeriodo(ctx, strconv.Itoa(time.Now().Year()), &Ivas); err != nil {
		return nil, err
	}

	parametros_ = map[string]interface{}{
		"EstadoActa":     EstadoActa,
		"EstadoElemento": EstadoElemento,
		"IVA":            Ivas,
	}

	return parametros_, nil
}
