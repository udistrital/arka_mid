package utilsHelper

import (
	"encoding/json"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/errorctrl"
)

// Marshal Hace el encode de cualquier estructura y la retorna en un string.
func Marshal(in interface{}, out *string) (outputError map[string]interface{}) {

	funcion := "Marshal"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	if out_, err := json.Marshal(in); err != nil {
		logs.Error(err)
		eval := " - json.Marshal(in)"
		return errorctrl.Error(funcion+eval, err, "500")
	} else {
		*out = string(out_[:])
	}

	return

}
