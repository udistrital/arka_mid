package consecutivos

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	"github.com/udistrital/arka_mid/utils_oas/request"
)

// GetById consulta controlador consecutivo/{id} del api consecutivos_crud.
func GetById(id int, consecutivo *models.Consecutivo) (outputError map[string]interface{}) {
	funcion := "GetById - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := ConsecutivosCRUD + "consecutivo/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &consecutivo); err != nil {
		logs.Error(urlcrud, err)
		eval := "request.GetJson(urlcrud, &consecutivo)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return nil
}
