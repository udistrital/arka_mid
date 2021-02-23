package tercerosHelper

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/models"
	// "github.com/udistrital/arka_mid/helpers/autenticacion"
	// "github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/utils_oas/request"
)

func GetContratista(idTercero int) (terceros []map[string]interface{}, outputError map[string]interface{}) {

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
	}
	logs.Debug("parametroContratistaID:", parametroContratistaID)

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

			for _, vincul := range vinculaciones {
				add := true
				for _, tercero := range terceros {
					if vincul.Id == tercero["Id"] {
						add = false
						break
					}
				}
				if add {
					terceros = append(terceros, map[string]interface{}{
						"Tercero": vincul.TerceroPrincipalId,
					})
				}
			}

		}
	}
	logs.Debug("terceros:", terceros)

	return terceros, nil
}
