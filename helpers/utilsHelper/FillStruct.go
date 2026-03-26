package utilsHelper

import (
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
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
