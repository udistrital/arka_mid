package consecutivos

import (
	"encoding/json"
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	"github.com/udistrital/arka_mid/utils_oas/request"
)

var ConsecutivosCRUD, _ = beego.AppConfig.String("consecutivosService")

// Post post controlador consecutivo del api consecutivos_crud
func Post(consecutivo interface{}) (outputError map[string]interface{}) {

	funcion := "Post"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := ConsecutivosCRUD + "consecutivo"
	response := new(models.RespuestaAPI1Interface)

	if err := request.SendJson(urlcrud, "POST", response, consecutivo); err != nil {
		logs.Error(urlcrud + ", " + err.Error())
		eval := ` - request.SendJson(urlcrud, "POST", response, consecutivo)`
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	if !response.Success {
		err := fmt.Errorf("%v", response.Message)
		logs.Error(err)
		eval := ` - request.SendJson(urlcrud, "POST", response, consecutivo)`
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	// helper local para mapear cualquier interface{} al struct destino
	fill := func(in interface{}, out interface{}) map[string]interface{} {
		raw, err := json.Marshal(in)
		if err != nil {
			return errorCtrl.Error(funcion+" - json.Marshal(in)", err, "502")
		}
		if err := json.Unmarshal(raw, out); err != nil {
			return errorCtrl.Error(funcion+" - json.Unmarshal(raw, out)", err, "502")
		}
		return nil
	}

	switch data := response.Data.(type) {
	case map[string]interface{}:
		return fill(data, consecutivo)

	case []interface{}:
		if len(data) == 0 {
			err := fmt.Errorf("response.Data llegó como arreglo vacío")
			return errorCtrl.Error(funcion+" - response.Data vacío", err, "502")
		}

		// intentar encontrar el consecutivo correcto si el payload es *models.Consecutivo
		if cons, ok := consecutivo.(*models.Consecutivo); ok {
			for _, item := range data {
				tmp := models.Consecutivo{}
				raw, err := json.Marshal(item)
				if err != nil {
					continue
				}
				if err := json.Unmarshal(raw, &tmp); err != nil {
					continue
				}

				if tmp.ContextoId == cons.ContextoId &&
					tmp.Year == cons.Year &&
					tmp.Descripcion == cons.Descripcion {
					return fill(item, consecutivo)
				}
			}
		}

		// fallback: usar el último elemento
		return fill(data[len(data)-1], consecutivo)

	default:
		err := fmt.Errorf("tipo inesperado en response.Data: %T", response.Data)
		return errorCtrl.Error(funcion+" - tipo inesperado", err, "502")
	}
}
