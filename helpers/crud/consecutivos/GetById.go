package consecutivos

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	"github.com/udistrital/arka_mid/models"
	errorCtrl "github.com/udistrital/utils_oas/v2/errorctrl"
	requestV2 "github.com/udistrital/utils_oas/v2/request"
)

// GetById consulta controlador consecutivo/{id} del api consecutivos_crud.
func GetById(ctx context.Context, id int, consecutivo *models.Consecutivo) (outputError map[string]interface{}) {
	funcion := "GetById - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := ConsecutivosCRUD + "consecutivo/" + strconv.Itoa(id)
	response := new(models.RespuestaAPI1Interface)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, response); err != nil {
		logs.Error(urlcrud, err)
		eval := "requestV2.GetWithContext(ctx, urlcrud, response)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	if !response.Success {
		err := fmt.Errorf("%v", response.Message)
		eval := "requestV2.GetWithContext(ctx, urlcrud, response)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	fill := func(in interface{}, out interface{}) map[string]interface{} {
		raw, err := json.Marshal(in)
		if err != nil {
			return errorCtrl.Error(funcion+"json.Marshal(in)", err, "502")
		}
		if err := json.Unmarshal(raw, out); err != nil {
			return errorCtrl.Error(funcion+"json.Unmarshal(raw, out)", err, "502")
		}
		return nil
	}

	switch data := response.Data.(type) {
	case map[string]interface{}:
		return fill(data, consecutivo)
	case []interface{}:
		if len(data) == 0 {
			err := fmt.Errorf("response.Data llegó como arreglo vacío")
			return errorCtrl.Error(funcion+"response.Data vacío", err, "502")
		}
		return fill(data[len(data)-1], consecutivo)
	default:
		err := fmt.Errorf("tipo inesperado en response.Data: %T", response.Data)
		return errorCtrl.Error(funcion+"tipo inesperado", err, "502")
	}
}
