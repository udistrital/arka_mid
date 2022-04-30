package utilsHelper

import (
	"encoding/json"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/errorctrl"
)

// Unmarshal Hace el decode de un string a el tipo de variable indicado.
func Unmarshal(in string, out interface{}) (outputError map[string]interface{}) {

	funcion := "Unmarshal"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	if err := json.Unmarshal([]byte(in), &out); err != nil {
		logs.Error(err)
		eval := " - json.Unmarshal([]byte(int), &out)"
		return errorctrl.Error(funcion+eval, err, "500")
	}

	return

}
