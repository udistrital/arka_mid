package utilsHelper

import (
	errorCtrl "github.com/udistrital/utils_oas/v2/errorctrl"
)

func FillStruct(in, out interface{}) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("FillStruct - Unhandled Error!", "500")

	var str string
	outputError = Marshal(in, &str)
	if outputError != nil {
		return
	}

	outputError = Unmarshal(str, &out)

	return

}
