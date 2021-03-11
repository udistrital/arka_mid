package tercerosHelper

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

//GetNombreTerceroById trae el nombre de un encargado por su id
func GetNombreTerceroById(idTercero string) (tercero map[string]interface{}, err error) {
	var urltercero string
	var personas []map[string]interface{}

	urltercero = "http://" + beego.AppConfig.String("tercerosService") + "datos_identificacion/?query=TerceroId__Id:" + idTercero + ",Activo:true"
	if resp, err := request.GetJsonTest(urltercero, &personas); err == nil && resp.StatusCode == 200 {
		for _, element := range personas {
			if len(element) == 0 {
				return nil, errors.New("No se encontro registro")
			} else {
				// fmt.Println("encargado: ", element)
				return map[string]interface{}{
					"Id":             element["TerceroId"].(map[string]interface{})["Id"],
					"Numero":         element["Numero"],
					"NombreCompleto": element["TerceroId"].(map[string]interface{})["NombreCompleto"],
				}, nil
			}
		}
	} else {
		if err == nil {
			err = fmt.Errorf("Undesired Status Code: %d", resp.StatusCode)
		}
		logs.Error(err)
		return nil, err
	}
	return

}

// GetTerceroByUsuarioWSO2 trae la información de un tercero a partir de su UsuarioWSO2
func GetTerceroByUsuarioWSO2(usuario string) (tercero map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetTerceroByUsuarioWSO2 - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var terceros []*models.Tercero
	urltercero := "http://" + beego.AppConfig.String("tercerosService") + "tercero"
	urltercero += "?fields=Id,NombreCompleto,TipoContribuyenteId"
	urltercero += "&query=Activo:true,UsuarioWSO2:" + usuario
	// logs.Info(urltercero)
	if resp, err := request.GetJsonTest(urltercero, &terceros); err == nil && resp.StatusCode == 200 {
		if len(terceros) == 1 && terceros[0].TipoContribuyenteId != nil {
			data := terceros[0]
			tercero = map[string]interface{}{
				"Id":             data.Id,
				"Numero":         "",
				"NombreCompleto": data.NombreCompleto,
			}
		} else if len(terceros) == 0 || terceros[0].TipoContribuyenteId == nil {
			err := fmt.Errorf("El usuario '%s' aún no está asignado a un registro en Terceros", usuario)
			outputError = map[string]interface{}{
				"funcion": "GetTerceroByUsuarioWSO2 - len(datosTerceros) == 1 && datosTerceros[0].TerceroId != nil",
				"err":     err,
				"status":  "404",
			}
			return nil, outputError
		} else { // len(terceros) > 1
			q := len(terceros)
			s := ""
			if q >= 10 {
				s = " - o más"
			}
			err := fmt.Errorf("El usuario '%s' tiene más de un registro en Terceros (%d registros%s)", usuario, q, s)
			logs.Warn(err)
			outputError = map[string]interface{}{
				"funcion": "GetTerceroByUsuarioWSO2 - len(datosTerceros) == 1 && datosTerceros[0].TerceroId != nil",
				"err":     err,
				"status":  "409",
			}
			return nil, outputError
		}
	} else {
		if err == nil {
			err = fmt.Errorf("Undesired Status Code: %d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetTerceroByUsuarioWSO2 - request.GetJsonTest(urltercero, &datosTerceros)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	return tercero, nil
}

func GetTerceroByDoc(doc string) (tercero *models.DatosIdentificacion, outputError map[string]interface{}) {
	urltercero := "http://" + beego.AppConfig.String("tercerosService") + "datos_identificacion/?query=Activo:true"
	urltercero += "Numero:" + doc
	var terceros []*models.DatosIdentificacion

	if resp, err := request.GetJsonTest(urltercero, &terceros); err == nil && resp.StatusCode == 200 {
		return terceros[0], nil
	}

	var vacio models.DatosIdentificacion
	return &vacio, nil
}
