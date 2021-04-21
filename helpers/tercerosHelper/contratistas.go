package tercerosHelper

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/mitchellh/mapstructure"

	// "github.com/udistrital/arka_mid/helpers/utilsHelper"
	// "github.com/udistrital/arka_mid/helpers/autenticacion"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetContratista trae la lista de contratistas registrados en Terceros, con opción de filtrar por ID
func GetContratista(idTercero int) (terceros []map[string]interface{}, outputError map[string]interface{}) {

	logs.Warn("prueba GetContratista actualizado")

	step := "0"

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetContratista - Uncaught Error! - after step:" + step,
				"err":     err,
				"status":  "500", // Uncaught error!
			}
			panic(outputError)
		}
	}()

	// PARTE 1. Traer los ID de los parámetros asociados a contratistas

	// Los siguientes son los códigos de los registros de la tabla "parametro" de la API
	// de parámetros, cuyo tipo_parámetro_id sea el asociado a Cargos.
	// Específicamente, los códigos de parámetros que se asignen a contratistas
	logs.Warn("prueba GetContratista actualizado - AAA")
	codigosParametroContratista := []string{"CPS"} // , "OPS", "PS"
	codigoTipoParamVinculacion := "TV"
	parametroContratistaID := make(map[string]int)
	logs.Warn("prueba GetContratista actualizado - BBB")

	panic(fmt.Errorf("paila"))

	logs.Warn("prueba GetContratista actualizado - CCC")

	var respBody models.RespuestaAPI1Arr
	urlParametros := "http://" + beego.AppConfig.String("parametrosService") + "parametro?limit=-1"
	urlParametros += "&fields=Id,CodigoAbreviacion"
	urlParametros += "&query=Activo:true,TipoParametroId__Activo:true,TipoParametroId__CodigoAbreviacion:" + codigoTipoParamVinculacion
	// logs.Debug("urlParametros:", urlParametros)
	if resp, err := request.GetJsonTest(urlParametros, &respBody); err == nil && resp.StatusCode == 200 {
		step = "1.1"
		if respBody.Data == nil || len(respBody.Data) == 0 || len(respBody.Data[0]) == 0 {
			err := fmt.Errorf("No están registrados los parámetros asociados a contratistas")
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetContratista - respBody.Data == nil || len(respBody.Data) == 0 || len(respBody.Data[0]) == 0",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
		step = "1.2"
		for _, paramVinculacion := range respBody.Data {
			// fmt.Printf("Param #%d: %#v\n", k, paramVinculacion)
			var codParam string
			if v, ok := paramVinculacion["CodigoAbreviacion"].(string); ok {
				codParam = v
			} else {
				continue
			}
			// fmt.Printf("codParam (%T): %v\n", codParam, codParam)
			for _, codigoContratista := range codigosParametroContratista {
				if id, ok := paramVinculacion["Id"].(float64); ok && codigoContratista == codParam {
					// fmt.Printf("P=P %v - T(id):%T - v:%f\n", paramVinculacion, paramVinculacion["Id"], paramVinculacion["Id"])
					parametroContratistaID[codigoContratista] = int(id)
					break
				}
			}
		}
		step = "1.3"
	} else {
		if err == nil {
			err = fmt.Errorf("Undesired status code - Got:%d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetContratista - request.GetJsonTest(urlParametros, &respBody)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	// logs.Debug("parametroContratistaID:", parametroContratistaID)

	// PARTE 2. Traer los terceros que tengan los ID anteriores en la tabla vinculacion
	tercerosMap := make(map[int](map[string]interface{}))
	for k, paramID := range parametroContratistaID {
		step = "2.1_" + k
		var vinculaciones []models.Vinculacion
		urlTerceros := "http://" + beego.AppConfig.String("tercerosService") + "vinculacion?limit=-1"
		urlTerceros += "&fields=Id,TerceroPrincipalId"
		urlTerceros += "&query=Activo:true,TipoVinculacionId:" + fmt.Sprint(paramID)
		if idTercero > 0 {
			urlTerceros += ",TerceroPrincipalId__Id:" + fmt.Sprint(idTercero)
		}
		// logs.Debug("urlTerceros:", urlTerceros)
		if resp, err := request.GetJsonTest(urlTerceros, &vinculaciones); err == nil && resp.StatusCode == 200 {
			step = "2.2_" + k
			if len(vinculaciones) == 0 || vinculaciones[0].TerceroPrincipalId == nil {
				continue
			}
			// fmt.Println("paramId:", paramId, "#vinculaciones: ", len(vinculaciones))

			// Lo siguiente es para que no se vuelva a agregar un tercero
			// cuando el tercero tenga más de una vinculación
			step = "2.3_" + k
			for k2, vincul := range vinculaciones {
				if _, found := tercerosMap[vincul.TerceroPrincipalId.Id]; found {
					continue
				}
				step = "2.4.1_" + k + "_" + strconv.Itoa(k2)
				terceroRecortado := map[string]interface{}{
					"Id":             vincul.TerceroPrincipalId.Id,
					"NombreCompleto": vincul.TerceroPrincipalId.NombreCompleto,
					// "UsuarioWSO2":    vincul.TerceroPrincipalId.UsuarioWSO2,
				}
				tercerosMap[vincul.TerceroPrincipalId.Id] = map[string]interface{}{
					// "Tercero": vincul.TerceroPrincipalId,
					"Tercero": terceroRecortado,
					// "TipoVinculacion":  vincul.TipoVinculacionId,
				}
				step = "2.4.2_" + k + "_" + strconv.Itoa(k2)
			}
			step = "2.5_" + k

		} else {
			if err == nil {
				err = fmt.Errorf("Undesired status code - Got:%d", resp.StatusCode)
			}
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetContratista - request.GetJsonTest(urlTerceros, &vinculaciones)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}

	}
	step = "2.6"
	for _, tercero := range tercerosMap {
		terceros = append(terceros, tercero)
	}
	step = "2.7"
	// logs.Debug("terceros:", terceros)

	// PARTE 3 Traer identificación disponible...
	for k, tercero := range terceros {
		kStr := strconv.Itoa(k)

		step = "3.1_" + kStr
		var terceroModelo models.Tercero
		if err := mapstructure.Decode(tercero["Tercero"], &terceroModelo); err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetContratista - mapstructure.Decode(tercero[\"Tercero\"], &terceroModelo)",
				"err":     err,
				"status":  "500",
			}
			return nil, outputError
		}
		// logs.Debug("terceroModelo:", terceroModelo)
		step = "3.2_" + kStr
		// 3.1 ... de terceros?
		var dataTerceros []map[string]interface{} // models.DatosIdentificacion
		urlDocTercero := "http://" + beego.AppConfig.String("tercerosService") + "datos_identificacion"
		urlDocTercero += "?fields=TipoDocumentoId,Numero"
		urlDocTercero += "&query=Activo:true,TerceroId__Id:" + fmt.Sprint(terceroModelo.Id)
		// logs.Debug("urlDocTercero: ", urlDocTercero)
		if resp, err := request.GetJsonTest(urlDocTercero, &dataTerceros); err == nil && resp.StatusCode == 200 {
			step = "3.3_" + kStr
			// tercero["DataTercerosDocumento"] = dataTerceros[0]
			if len(dataTerceros) == 1 && dataTerceros[0]["Id"] != 0 {
				step = "3.4_" + kStr
				var dataModel models.DatosIdentificacion
				if err := mapstructure.Decode(dataTerceros[0], &dataModel); err != nil {
					logs.Error(err)
					outputError = map[string]interface{}{
						"funcion": "GetContratista - mapstructure.Decode(dataTerceros[0], &dataModel)",
						"err":     err,
						"status":  "500",
					}
					return nil, outputError
				}
				step = "3.5_" + kStr
				dataRecortada := map[string]interface{}{
					// "TipoDocumentoId": dataModel.TipoDocumentoId,
					"TipoDocumentoId": map[string]interface{}{
						"Id":     dataModel.TipoDocumentoId.Id,
						"Nombre": dataModel.TipoDocumentoId.Nombre,
					},
					"Numero": dataModel.Numero,
				}
				tercero["Identificacion"] = dataRecortada
				step = "3.6_" + kStr
			} else {
				err := fmt.Errorf("Hay +/- un documento registrado como Activo para el Tercero con ID: %d", terceroModelo.Id)
				logs.Warn(err)
			}
		} else {
			if err == nil {
				err = fmt.Errorf("Undesired status code - Got:%d", resp.StatusCode)
			}
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetContratista - request.GetJsonTest(urlDocTercero, &dataTerceros)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}

	}

	return terceros, nil
}
