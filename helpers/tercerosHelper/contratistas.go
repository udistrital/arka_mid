package tercerosHelper

import (
	"fmt"

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

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/GetContratista - Uncaught Error!",
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
	codigosParametroContratista := []string{"CPS"} // , "OPS", "PS"
	codigoTipoParamVinculacion := "TV"
	parametroContratistaID := make(map[string]int)

	var respBody models.RespuestaAPI1Arr
	urlParametros := "http://" + beego.AppConfig.String("parametrosService") + "parametro?limit=-1"
	urlParametros += "&fields=Id,CodigoAbreviacion"
	urlParametros += "&query=Activo:true,TipoParametroId__Activo:true,TipoParametroId__CodigoAbreviacion:" + codigoTipoParamVinculacion
	// fmt.Println(urlParametros)
	if resp, err := request.GetJsonTest(urlParametros, &respBody); err == nil && resp.StatusCode == 200 {
		for _, paramVinculacion := range respBody.Data {
			// fmt.Printf("Param #%d: %#v\n", k, paramVinculacion)
			codParam := paramVinculacion["CodigoAbreviacion"]
			// fmt.Printf("codParam (%T): %v\n", codParam, codParam)
			for _, codigoContratista := range codigosParametroContratista {
				if codigoContratista == codParam {
					// fmt.Printf("P=P %v - T(id):%T - v:%f\n", paramVinculacion, paramVinculacion["Id"], paramVinculacion["Id"])
					parametroContratistaID[codigoContratista] = int(paramVinculacion["Id"].(float64))
				}
			}
		}
	} else {
		if err == nil {
			err = fmt.Errorf("Undesired status code - Got:%d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/GetContratista - request.GetJsonTest(urlParametros, &respBody)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	// logs.Debug("parametroContratistaID:", parametroContratistaID)

	// PARTE 2. Traer los terceros que tengan los ID anteriores en la tabla vinculacion
	for _, paramID := range parametroContratistaID {
		var vinculaciones []models.Vinculacion
		urlTerceros := "http://" + beego.AppConfig.String("tercerosService") + "vinculacion?limit=-1"
		urlTerceros += "&fields=Id,TerceroPrincipalId"
		urlTerceros += "&query=Activo:true,TipoVinculacionId:" + fmt.Sprint(paramID)
		if idTercero > 0 {
			urlTerceros += ",TerceroPrincipalId__Id:" + fmt.Sprint(idTercero)
		}
		// fmt.Println(urlTerceros)
		if resp, err := request.GetJsonTest(urlTerceros, &vinculaciones); err == nil && resp.StatusCode == 200 {
			if len(vinculaciones) == 0 || vinculaciones[0].TerceroPrincipalId == nil {
				continue
			}
			// fmt.Println("paramId:", paramId, "#vinculaciones: ", len(vinculaciones))

			// Lo siguiente es para que no se vuelva a agregar un tercero
			// cuando el tercero tenga más de una vinculación
			for _, vincul := range vinculaciones {
				add := true
				for _, tercero := range terceros {
					if mTercero := tercero["Tercero"].(*models.Tercero); vincul.TerceroPrincipalId.Id == mTercero.Id {
						add = false
						break
					}
				}
				if add {
					terceroRecortado := map[string]interface{}{
						"Id":             vincul.TerceroPrincipalId.Id,
						"NombreCompleto": vincul.TerceroPrincipalId.NombreCompleto,
						// "UsuarioWSO2":    vincul.TerceroPrincipalId.UsuarioWSO2,
					}
					terceros = append(terceros, map[string]interface{}{
						// "Tercero": vincul.TerceroPrincipalId,
						"Tercero": terceroRecortado,
						// "TipoVinculacion":  vincul.TipoVinculacionId,
					})
				}
			}

		} else {
			if err == nil {
				err = fmt.Errorf("Undesired status code - Got:%d", resp.StatusCode)
			}
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "/GetContratista - request.GetJsonTest(urlTerceros, &vinculaciones)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
	}
	// logs.Debug("terceros:", terceros)

	// PARTE 3 Traer identificación disponible...
	for _, tercero := range terceros {

		var terceroModelo models.Tercero
		if err := mapstructure.Decode(tercero["Tercero"], &terceroModelo); err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "/GetContratista - mapstructure.Decode(tercero[\"Tercero\"], &terceroModelo)",
				"err":     err,
				"status":  "500",
			}
			return nil, outputError
		}
		// logs.Debug("terceroModelo:", terceroModelo)

		// 3.1 ... de terceros?
		var dataTerceros []map[string]interface{} // models.DatosIdentificacion
		urlDocTercero := "http://" + beego.AppConfig.String("tercerosService") + "datos_identificacion"
		urlDocTercero += "?fields=TipoDocumentoId,Numero"
		urlDocTercero += "&query=Activo:true,TerceroId__Id:" + fmt.Sprint(terceroModelo.Id)
		// logs.Debug("urlDocTercero: ", urlDocTercero)
		if resp, err := request.GetJsonTest(urlDocTercero, &dataTerceros); err == nil && resp.StatusCode == 200 {
			// tercero["DataTercerosDocumento"] = dataTerceros[0]
			if len(dataTerceros) == 1 && dataTerceros[0]["Id"] != 0 {
				var dataModel models.DatosIdentificacion
				if err := mapstructure.Decode(dataTerceros[0], &dataModel); err != nil {
					logs.Error(err)
					outputError = map[string]interface{}{
						"funcion": "/GetContratista - mapstructure.Decode(dataTerceros[0], &dataModel)",
						"err":     err,
						"status":  "500",
					}
					return nil, outputError
				}
				dataRecortada := map[string]interface{}{
					// "TipoDocumentoId": dataModel.TipoDocumentoId,
					"TipoDocumentoId": map[string]interface{}{
						"Id":     dataModel.TipoDocumentoId.Id,
						"Nombre": dataModel.TipoDocumentoId.Nombre,
					},
					"Numero": dataModel.Numero,
				}
				tercero["Identificacion"] = dataRecortada
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
				"funcion": "/GetContratista - request.GetJsonTest(urlDocTercero, &dataTerceros)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}

		// 3.2 ... de Autenticacion MID?
		/*
			if data, err := autenticacion.DataUsuario(terceroModelo.UsuarioWSO2); err == nil {
				tercero["DataAutenticacion"] = data
				// logs.Debug("dataAutenticacion:", data)
			} else {
				return nil, err
			}
		*/

	}

	return terceros, nil
}
