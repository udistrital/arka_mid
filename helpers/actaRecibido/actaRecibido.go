package actaRecibido

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	// "github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/arka_mid/helpers/crud/administrativa"
	"github.com/udistrital/arka_mid/helpers/crud/oikos"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

func RemoveIndex(s []byte, index int) []byte {
	return append(s[:index], s[index+1:]...)
}

// GetAllParametrosActa ...
func GetAllParametrosActa() (Parametros []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetAllParametrosActa - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var (
		Unidades       interface{}
		TipoBien       interface{}
		EstadoActa     interface{}
		EstadoElemento interface{}
		ss             map[string]interface{}
		Parametro      []interface{}
		Valor          []interface{}
		IvaTest        []Imp
		Ivas           []Imp
	)

	parametros := make([]map[string]interface{}, 0)

	urlActasTipoBien := "http://" + beego.AppConfig.String("actaRecibidoService") + "tipo_bien?limit=-1"
	if _, err := request.GetJsonTest(urlActasTipoBien, &TipoBien); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllParametrosActa - request.GetJsonTest(urlActasTipoBien, &TipoBien)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

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

	urlParametros := "http://" + beego.AppConfig.String("parametrosService") + "parametro_periodo?query=PeriodoId__Nombre:2021,ParametroId__TipoParametroId__Id:12"
	if _, err := request.GetJsonTest(urlParametros, &ss); err == nil {

		var data []map[string]interface{}
		if jsonString, err := json.Marshal(ss["Data"]); err == nil {
			if err := json.Unmarshal(jsonString, &data); err == nil {
				for _, valores := range data {
					Parametro = append(Parametro, valores["ParametroId"])
					v := []byte(fmt.Sprintf("%v", valores["Valor"]))
					var valorUnm interface{}
					if err := json.Unmarshal(v, &valorUnm); err == nil {
						Valor = append(Valor, valorUnm)
					}
				}
			}
		}

		if jsonbody1, err := json.Marshal(Parametro); err == nil {
			if err := json.Unmarshal(jsonbody1, &Ivas); err != nil {
				fmt.Println(err)
				return
			}
		}

		if jsonbody1, err := json.Marshal(Valor); err == nil {
			if err := json.Unmarshal(jsonbody1, &IvaTest); err != nil {
				fmt.Println(err)
				return
			}
		}

		for i, valores := range IvaTest {
			IvaTest[i].CodigoAbreviacion = valores.CodigoAbreviacion
		}
		for i, valores := range Ivas {
			IvaTest[i].BasePesos = valores.BasePesos
			IvaTest[i].BaseUvt = valores.BaseUvt
			IvaTest[i].PorcentajeAplicacion = valores.PorcentajeAplicacion
			IvaTest[i].CodigoAbreviacion = valores.CodigoAbreviacion
		}

	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllParametrosActa - request.GetJsonTest(urlParametros, &ss)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	if outputError = administrativa.GetUnidades(&Unidades); outputError != nil {
		return
	}

	parametros = append(parametros, map[string]interface{}{
		"Unidades":       Unidades,
		"TipoBien":       TipoBien,
		"EstadoActa":     EstadoActa,
		"EstadoElemento": EstadoElemento,
		"IVA":            IvaTest,
	})

	return parametros, nil
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

	var (
		Dependencias []models.Dependencia
		Sedes        interface{}
		Ubicaciones  interface{}
	)
	parametros := make([]map[string]interface{}, 0)

	if err := oikos.GetDependencia("", "", "", "", -1, 0, &Dependencias); err != nil {
		logs.Warning(err)
	}

	if Ubicaciones, outputError = oikos.GetAllAsignacion("?limit=-1"); outputError != nil {
		return
	}

	if Sedes, outputError = oikos.GetAllEspacioFisico("?query=TipoEspacioFisicoId.Id:1&limit=-1"); outputError != nil {
		return
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
