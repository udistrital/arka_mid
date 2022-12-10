package actaRecibido

import (
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/crud/administrativa"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

func RemoveIndex(s []byte, index int) []byte {
	return append(s[:index], s[index+1:]...)
}

// GetAllParametrosActa Consulta diferentes valores paramétricos
func GetAllParametrosActa() (parametros_ []map[string]interface{}, outputError map[string]interface{}) {

	funcion := "GetAllParametrosActa - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	var (
		Unidades       interface{}
		EstadoActa     interface{}
		EstadoElemento interface{}
		Ivas           []models.Iva
	)

	urlActasEstadoActa := "http://" + beego.AppConfig.String("actaRecibidoService") + "estado_acta?limit=-1"
	if _, err := request.GetJsonTest(urlActasEstadoActa, &EstadoActa); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllParametrosActa - request.GetJsonTest(urlActasEstadoActa, &EstadoActa)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	urlACtasEstadoElem := "http://" + beego.AppConfig.String("actaRecibidoService") + "estado_elemento?limit=-1"
	if _, err := request.GetJsonTest(urlACtasEstadoElem, &EstadoElemento); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllParametrosActa - request.GetJsonTest(urlACtasEstadoElem, &EstadoElemento)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	if err := parametros.GetAllIVAByPeriodo(strconv.Itoa(time.Now().Year()-1), &Ivas); err != nil {
		return nil, err
	}

	if outputError = administrativa.GetUnidades(&Unidades); outputError != nil {
		return
	}

	parametros_ = append(parametros_, map[string]interface{}{
		"Unidades":       Unidades,
		"EstadoActa":     EstadoActa,
		"EstadoElemento": EstadoElemento,
		"IVA":            Ivas,
	})

	return parametros_, nil
}

// GetIdElementoPlaca Busca el id de un elemento a partir de su placa
func GetIdElementoPlaca(placa string) (idElemento string, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetIdElementoPlaca - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var urlelemento string
	var elemento []map[string]interface{}
	urlelemento = "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento/?query=Placa:" + placa + "&fields=Id&limit=1"
	if response, err := request.GetJsonTest(urlelemento, &elemento); err == nil {

		if response.StatusCode == 200 {
			for _, element := range elemento {
				if len(element) == 0 {
					return "", nil
				} else {
					return strconv.Itoa(int((element["Id"]).(float64))), nil
				}

			}
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetSoportes - request.GetJsonTest(urlelemento, &elemento)",
			"err":     err,
			"status":  "502",
		}
		return "", outputError
	}
	return
}
