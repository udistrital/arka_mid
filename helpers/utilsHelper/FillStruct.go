package utilsHelper

import (
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
)

func FillStruct(in, out interface{}) (outputError map[string]interface{}) {

	funcion := "FillStruct - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	if err := formatdata.FillStruct(in, &out); err != nil {
		logs.Error(err)
		eval := "formatdata.FillStruct(in, &out)"
		return errorctrl.Error(funcion+eval, err, "500")
	}

	return

}
