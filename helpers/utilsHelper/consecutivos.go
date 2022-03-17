package utilsHelper

import (
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

func GetConsecutivo(format string, contextoId int, descripcion string) (consecutivo string, consecutivoId int, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetConsecutivo - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var res map[string]interface{}

	year := time.Now().Year()
	data := models.Consecutivo{
		Id:          0,
		ContextoId:  contextoId,
		Year:        year,
		Consecutivo: 0,
		Descripcion: descripcion,
		Activo:      true,
	}
	url := "http://" + beego.AppConfig.String("consecutivosService") + "consecutivo"

	if err := request.SendJson(url, "POST", &res, &data); err == nil {
		consecutivo = fmt.Sprintf(format, res["Data"].(map[string]interface{})["Consecutivo"])
		consecutivoId = int(res["Data"].(map[string]interface{})["Id"].(float64))
	} else if strings.Contains(err.Error(), "invalid character") {
		logs.Error(err)
		consecutivo, consecutivoId, outputError = GetConsecutivo(format, contextoId, descripcion)
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetConsecutivo - request.SendJson(url, \"POST\", &res, &data)",
			"err":     err,
			"status":  "502",
		}
	}
	return consecutivo, consecutivoId, outputError
}

func FormatConsecutivo(prefix string, consecutivo string, suffix string) (consFormat string) {
	return prefix + consecutivo + suffix
}
