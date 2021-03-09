package tercerosHelper

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/mitchellh/mapstructure"

	// "github.com/udistrital/arka_mid/helpers/utilsHelper"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

// GetProveedor trae la lista de proveedores registrados en Terceros, con opcion de filtrar por ID
func GetProveedor(idProveedor int) (terceros []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/GetProveedor - Uncaught Error!",
				"err":     err,
				"status":  "500", // Uncaught error!
			}
			panic(outputError)
		}
	}()

	// PARTE 1. Traer los ID de los tipo_tercero asociados a proveedores
	codigosTipoTerceroProveedor := []string{"ENTIDAD_PUBLICA", "ENTIDAD_PRIVADA", "ENTIDAD_MIXTA"} // , "OPS", "PS"
	tipoTerceroIDs := make(map[string]int)

	var data []map[string]interface{}
	urlTipos := "http://" + beego.AppConfig.String("tercerosService") + "tipo_tercero?limit=-1"
	urlTipos += "&fields=Id,CodigoAbreviacion"
	urlTipos += "&query=Activo:true"
	// fmt.Println(urlParametros)
	if resp, err := request.GetJsonTest(urlTipos, &data); err == nil && resp.StatusCode == 200 {
		if len(data) == 0 || len(data[0]) == 0 {
			err := fmt.Errorf("No hay tipo_tercero registrados")
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "/GetProveedor - request.GetJsonTest(urlTipos, &data)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
		for _, tipoDisponible := range data {
			for _, tipoSuficiente := range codigosTipoTerceroProveedor {
				if tipoDisponible["CodigoAbreviacion"] == tipoSuficiente {
					// fmt.Printf("P=P %v - T(id):%T - v:%f\n", paramVinculacion, paramVinculacion["Id"], paramVinculacion["Id"])
					tipoTerceroIDs[tipoSuficiente] = int(tipoDisponible["Id"].(float64))
				}
			}
		}
	} else {
		if err == nil {
			err = fmt.Errorf("Undesired status code - Got:%d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/GetProveedor - request.GetJsonTest(urlTipos, &data)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	// logs.Debug(tipoTerceroIDs)

	// PARTE 2 - Traer los terceros con los tipo_tercero requeridos
	tercerosMap := make(map[int](map[string]interface{}))
	for _, id := range tipoTerceroIDs {
		// logs.Debug("param:", param, "- id:", id)
		data = make([]map[string]interface{}, 0)
		urlTerceros := "http://" + beego.AppConfig.String("tercerosService") + "tercero_tipo_tercero?limit=-1"
		urlTerceros += "&fields=TerceroId"
		urlTerceros += "&query=Activo:true,TipoTerceroId__Id:" + fmt.Sprint(id)
		if idProveedor > 0 {
			// logs.Debug("idProveedor:", idProveedor)
			urlTerceros += ",TerceroPrincipalId__Id:" + fmt.Sprint(idProveedor)
		}
		// logs.Debug("urlTerceros:", urlTerceros)
		// fmt.Println(urlTerceros)
		if resp, err := request.GetJsonTest(urlTerceros, &data); err == nil && resp.StatusCode == 200 {
			if len(data) == 0 || len(data[0]) == 0 {
				logs.Debug("No se encontraron terceros. Saltando al siguiente parametro")
				continue
			}

			for _, terceroTipo := range data {
				// logs.Debug("terceroTipo:", terceroTipo)

				var terData models.Tercero
				if err := mapstructure.Decode(terceroTipo["TerceroId"], &terData); err != nil {
					logs.Error(err)
					outputError = map[string]interface{}{
						"funcion": "/GetProveedor - mapstructure.Decode(terceroTipo[\"TerceroId\"], &terData)",
						"err":     err,
						"status":  "500",
					}
					return nil, outputError
				}

				// logs.Debug("terData:", terData)
				terceroClean := map[string]interface{}{
					"Id":             terData.Id,
					"NombreCompleto": terData.NombreCompleto,
				}
				tercerosMap[terData.Id] = terceroClean
				// logs.Debug("terceroClean:", terceroClean)
			}
		}
	}
	formatdata.JsonPrint(tercerosMap)

	return
}
