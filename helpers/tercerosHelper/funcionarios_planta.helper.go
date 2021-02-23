package tercerosHelper

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetFuncionariosPlanta trae los funcionarios de planta
func GetFuncionariosPlanta(idTercero int) (terceros []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/GetFuncionariosPlanta - Uncaught Error!",
				"err":     err,
				"status":  "500", // Error no manejado!
			}
			panic(outputError)
		}
	}()

	// PARTE 1. Traer los ID de los parámetros asociados a funcionarios de planta

	// Los siguientes son los códigos de los registros de la tabla "parametro" de la API
	// de parámetros, cuyo tipo_parámetro_id sea el asociado a Tipo de Vinculacion.
	// Específicamente los códigos de parámetros que se asignen a administrativos o docentes de planta:
	codigosParametroFuncionarioPlanta := []string{"DP", "AP"}
	CodigoTipoParamVinculacion := "TV"
	parametroPlantaID := make(map[string]int)

	var respBody models.RespuestaAPI1Arr
	urlParametros := "http://" + beego.AppConfig.String("parametrosService") + "parametro?limit=-1"
	urlParametros += "&fields=Id,CodigoAbreviacion"
	urlParametros += "&query=Activo:true,TipoParametroId__Activo:true,TipoParametroId__CodigoAbreviacion:" + CodigoTipoParamVinculacion
	// fmt.Println(urlParametros)
	if resp, err := request.GetJsonTest(urlParametros, &respBody); err == nil && resp.StatusCode == 200 {
		// fmt.Printf("Data: %v\n", respBody.Data)

		for _, paramVinculacion := range respBody.Data {
			// fmt.Printf("Param #%d: %#v\n", k, paramVinculacion)
			codParam := paramVinculacion["CodigoAbreviacion"]
			// fmt.Printf("codParam (%T): %v\n", codParam, codParam)
			for _, codigoFuncPlanta := range codigosParametroFuncionarioPlanta {
				if codigoFuncPlanta == codParam {
					// fmt.Printf("P=P %v - T(id):%T - v:%f\n", paramVinculacion, paramVinculacion["Id"], paramVinculacion["Id"])
					parametroPlantaID[codigoFuncPlanta] = int(paramVinculacion["Id"].(float64))
				}
			}
		}
		// fmt.Printf("ids: %#v\n", parametroPlantaID)
	} else {
		if err == nil {
			err = fmt.Errorf("Undesired status code - Got:%d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/GetFuncionariosPlanta - request.GetJsonTest(urlParametros, &respBody)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	// PARTE 2. Traer los terceros que tengan estos IDs en la tabla vinculacion

	// NOTA: Esta parte se podría mejorar aplicando concurrencia. Vease:
	// https://gobyexample.com/goroutines
	// https://gobyexample.com/waitgroups
	// https://mayurwadekar2.medium.com/concurrency-and-parallelism-in-golang-c8327701fd94
	for _, paramID := range parametroPlantaID {

		var vinculaciones []models.Vinculacion
		urlTerceros := "http://" + beego.AppConfig.String("tercerosService") + "vinculacion?limit=-1"
		urlTerceros += "&fields=Id,TerceroPrincipalId,TipoVinculacionId,DependenciaId"
		urlTerceros += "&query=Activo:true,TipoVinculacionId:" + fmt.Sprint(paramID)
		if idTercero > 0 {
			urlTerceros += ",TerceroPrincipalId__Id:" + fmt.Sprint(idTercero)
		}
		// fmt.Println(urlTerceros)
		if resp, err := request.GetJsonTest(urlTerceros, &vinculaciones); err == nil && resp.StatusCode == 200 {

			if len(vinculaciones) == 0 || vinculaciones[0].Id == 0 {
				continue
			}
			// fmt.Println("paramId:", paramId, "#vinculaciones: ", len(vinculaciones))

			// Lo siguiente es para que no se vuelva a agregar un tercero
			// cuando el tercero tenga más de una vinculación
			for _, vincul := range vinculaciones {
				add := true
				for _, tercero := range terceros {
					if mTercero := tercero["TerceroPrincipal"].(models.Tercero); vincul.TerceroPrincipalId.Id == mTercero.Id {
						add = false
						break
					}
				}
				if add {
					terceros = append(terceros, map[string]interface{}{
						"TerceroPrincipal": vincul.TerceroPrincipalId,
						"TipoVinculacion":  vincul.TipoVinculacionId,
						"DependenciaId":    vincul.DependenciaId,
					})
				}
			}
		} else {
			if err == nil {
				err = fmt.Errorf("Undesired status code - Got:%d", resp.StatusCode)
			}
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "/GetFuncionariosPlanta - request.GetJsonTest(urlTerceros, &vinculaciones)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
	}
	// fmt.Println("#terceros:", len(terceros))

	// PARTE 3 - Agregar Información complementaria de Sede y Dependencia (si la hay)

	var sedesDependencias []models.AsignacionEspacioFisicoDependencia
	for _, tercero := range terceros {
		// fmt.Println("k:", k, "tercero:", tercero)

		// 3.1 traer los registros necesarios/disponibles
		consultar := true
		for _, seDependencia := range sedesDependencias {
			if seDependencia.DependenciaId.Id == tercero["Dependencia"] {
				// fmt.Println("seDepend:", seDependencia.Id)
				consultar = false
			}
		}
		if consultar {

			var resBody []models.AsignacionEspacioFisicoDependencia
			urlOikos := "http://" + beego.AppConfig.String("oikos2Service") + "asignacion_espacio_fisico_dependencia?limit=-1"
			urlOikos += "&fields=Id,EspacioFisicoId,DependenciaId&query=Activo:true"
			// urlOikos += ",EspacioFisicoId__TipoEspacioFisicoId__Nombre:SEDE"
			urlOikos += ",EspacioFisicoId__TipoEspacioFisicoId__CodigoAbreviacion:Tipo_1"
			urlOikos += ",DependenciaId__Id:" + fmt.Sprint(tercero["DependenciaId"])
			// fmt.Println(urlOikos)
			if resp, err := request.GetJsonTest(urlOikos, &resBody); err == nil && resp.StatusCode == 200 {

				if len(resBody) == 0 || resBody[0].Id == 0 {
					// No se encontró relación sede-dependencia para el tercero actual
					continue
				}

				for _, v := range resBody {
					sedesDependencias = append(sedesDependencias, v)
				}

			} else {
				if err == nil {
					err = fmt.Errorf("Undesired status code - Got:%d", resp.StatusCode)
				}
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "/GetFuncionariosPlanta - request.GetJsonTest(urlOikos, &resBody)",
					"err":     err,
					"status":  "502",
				}
				return nil, outputError
			}
		}

		// 3.2 asignar la información disponible
		for _, seDep := range sedesDependencias {
			if tercero["DependenciaId"] == seDep.DependenciaId.Id {
				tercero["Sede"] = seDep.EspacioFisicoId
				tercero["Dependencia"] = seDep.DependenciaId
				break
			}
		}
	}
	// fmt.Printf("%v\n", sedesDependencias)

	return terceros, nil
}
