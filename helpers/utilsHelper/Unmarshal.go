package utilsHelper

import (
	"encoding/json"

	"github.com/beego/beego/v2/core/logs"
	errorCtrl "github.com/udistrital/utils_oas/v2/errorctrl"
)

// Unmarshal Hace el decode de un string a el tipo de variable indicado.
func Unmarshal(in string, out interface{}) (outputError map[string]interface{}) {

	funcion := "Unmarshal - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	err := json.Unmarshal([]byte(in), &out)
	if err != nil {
		logs.Error(err)
		eval := "json.Unmarshal([]byte(int), &out)"
		outputError = errorCtrl.Error(funcion+eval, err, "500")
	}

	return

}
