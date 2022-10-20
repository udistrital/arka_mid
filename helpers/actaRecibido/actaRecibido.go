package actaRecibido

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
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

// GetAllParametrosSoporte ...
func GetAllParametrosSoporte() (Parametros []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetAllParametrosSoporte - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var Dependencias interface{}
	var Sedes interface{}
	var Ubicaciones interface{}
	parametros := make([]map[string]interface{}, 0)

	urlOikosDependencia := "http://" + beego.AppConfig.String("oikosService") + "dependencia?limit=-1"
	if _, err := request.GetJsonTest(urlOikosDependencia, &Dependencias); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllParametrosSoporte - request.GetJsonTest(urlOikosDependencia, &Dependencias)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	urlOikosAsignacion := "http://" + beego.AppConfig.String("oikosService") + "asignacion_espacio_fisico_dependencia?limit=-1"
	if _, err := request.GetJsonTest(urlOikosAsignacion, &Ubicaciones); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllParametrosSoporte - request.GetJsonTest(urlOikosAsignacion, &Ubicaciones)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	urlOikosEspFis := "http://" + beego.AppConfig.String("oikosService") + "espacio_fisico?query=TipoEspacioFisicoId__Nombre:SEDE&limit=-1"
	if _, err := request.GetJsonTest(urlOikosEspFis, &Sedes); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllParametrosSoporte - request.GetJsonTest(urlOikosEspFis, &Sedes)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	parametros = append(parametros, map[string]interface{}{
		"Dependencias": Dependencias,
		"Ubicaciones":  Ubicaciones,
		"Sedes":        Sedes,
	})

	return parametros, nil
}

// GetAsignacionSedeDependencia ...
func GetAsignacionSedeDependencia(Datos models.GetSedeDependencia) (Parametros []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetAsignacionSedeDependencia - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	if Datos.Sede == nil {
		err := fmt.Errorf("sede no especificada")
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAsignacionSedeDependencia - Datos.Sede == nil",
			"err":     err,
			"status":  "400",
		}
		return nil, outputError
	}

	var Ubicaciones []map[string]interface{}
	var Parametros2 []map[string]interface{}
	// logs.Debug("Datos:")
	// formatdata.JsonPrint(Datos)
	// fmt.Println("")
	oikosUrl := "http://" + beego.AppConfig.String("oikosService") + "asignacion_espacio_fisico_dependencia?limit=-1"
	oikosUrl += "&query=DependenciaId.Id:" + strconv.Itoa(Datos.Dependencia.Id)
	// logs.Debug("oikosUrl:", oikosUrl)
	if resp, err := request.GetJsonTest(oikosUrl, &Ubicaciones); err == nil && resp.StatusCode == 200 { // (2) error servicio caido
		for _, relacion := range Ubicaciones {
			var data map[string]interface{}
			if jsonString, err := json.Marshal(relacion["EspacioFisicoId"]); err == nil {
				if err2 := json.Unmarshal(jsonString, &data); err2 == nil {
					if number := strings.Index(fmt.Sprintf("%v", data["CodigoAbreviacion"]), Datos.Sede.CodigoAbreviacion); number != -1 {
						Parametros2 = append(Parametros2, map[string]interface{}{
							"Id":              relacion["Id"],
							"DependenciaId":   relacion["DependenciaId"],
							"EspacioFisicoId": relacion["EspacioFisicoId"],
							"Estado":          relacion["Estado"],
							"FechaFin":        relacion["FechaFin"],
							"FechaInicio":     relacion["FechaInicio"],
							"Nombre":          data["Nombre"],
						})
					}
					Parametros = append(Parametros, map[string]interface{}{
						"Relaciones": Parametros2,
					})

				} else {
					logs.Error(err2)
					outputError = map[string]interface{}{
						"funcion": "GetAsignacionSedeDependencia - json.Unmarshal(jsonString, &data)",
						"err":     err2,
						"status":  "500",
					}
					return nil, outputError
				}
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "GetAsignacionSedeDependencia - json.Marshal(relacion[\"EspacioFisicoId\"])",
					"err":     err,
					"status":  "500",
				}
				return nil, outputError
			}
		}

		return Parametros, nil

	} else {
		if err == nil {
			err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAsignacionSedeDependencia - request.GetJsonTest(oikosUrl, &Ubicaciones)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

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

// GetAllElementosConsumo obtiene todos los elementos de consumo
func GetAllElementosConsumo() (elementos []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetAllElementosConsumo - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	url := "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento?query=TipoBienId:1,Activo:true"
	if response, err := request.GetJsonTest(url, &elementos); err == nil && response.StatusCode == 200 {
		if len(elementos) == 0 {
			err := errors.New("no hay elementos")
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetAllElementosConsumo - len(elementos) == 0",
				"err":     err,
				"status":  "404",
			}
			return nil, outputError
		} else {
			return elementos, nil
		}

	} else {
		if err == nil {
			err = fmt.Errorf("undesired Status Code: %d", response.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllElementosConsumo - request.GetJsonTest(url, &elementos)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

}
