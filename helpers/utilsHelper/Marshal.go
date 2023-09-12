package utilsHelper

import (
	"encoding/json"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// Marshal Hace el encode de cualquier estructura y la retorna en un string.
func Marshal(in interface{}, out *string) (outputError map[string]interface{}) {

	funcion := "Marshal - "
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	if out_, err := json.Marshal(in); err != nil {
		logs.Error(err)
		eval := "json.Marshal(in)"
		outputError = errorCtrl.Error(funcion+eval, err, "500")
	} else {
		*out = string(out_[:])
	}

	return

}

func String(v string) *string     { return &v }
func Time(v time.Time) *time.Time { return &v }
func Int(v int) *int              { return &v }
