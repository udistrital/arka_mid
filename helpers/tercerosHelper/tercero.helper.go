package tercerosHelper

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

//GetNombreTerceroById trae el nombre de un encargado por su id
func GetNombreTerceroById(idTercero string) (tercero map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetNombreTerceroById - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	if v, err := strconv.Atoi(idTercero); err != nil || v <= 0 {
		err := fmt.Errorf("ID MUST be an integer > 0 - Got:%s", idTercero)
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetNombreTerceroById - len(personas) != 1 || len(personas[0])",
			"err":     err,
			"status":  "400",
		}
		return nil, outputError
	}

	var personas []map[string]interface{}

	urltercero := "http://" + beego.AppConfig.String("tercerosService") + "datos_identificacion"
	urltercero += "?query=Activo:true,TerceroId__Id:" + idTercero
	// logs.Debug("urltercero:", urltercero)
	if resp, err := request.GetJsonTest(urltercero, &personas); err == nil && resp.StatusCode == 200 {

		if len(personas) != 1 || len(personas[0]) == 0 {
			var status string
			if len(personas) > 1 {
				err = fmt.Errorf("Hay más de un documento para Tercero.Id=%s", idTercero)
				status = "409"
			} else {
				err = fmt.Errorf("No se encontró el Tercero.Id=%s y/o un documento asociado", idTercero)
				status = "404"
			}
			logs.Warn(err)
			outputError = map[string]interface{}{
				"funcion": "GetNombreTerceroById - len(personas) != 1 || len(personas[0])",
				"err":     err,
				"status":  status,
			}
			return nil, outputError
		}

		return map[string]interface{}{
			"Id":             personas[0]["TerceroId"].(map[string]interface{})["Id"],
			"Numero":         personas[0]["Numero"],
			"NombreCompleto": personas[0]["TerceroId"].(map[string]interface{})["NombreCompleto"],
		}, nil
	} else {
		if err == nil {
			err = fmt.Errorf("Undesired Status Code: %d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetNombreTerceroById - request.GetJsonTest(urltercero, &personas)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
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
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetTerceroByDoc - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()
	urltercero := "http://" + beego.AppConfig.String("tercerosService") + "datos_identificacion/?query=Activo:true,"
	urltercero += "Numero:" + doc
	var terceros []*models.DatosIdentificacion

	if resp, err := request.GetJsonTest(urltercero, &terceros); err == nil && resp.StatusCode == 200 {
		if len(terceros) == 1 {
			return terceros[0], nil
		} else if len(terceros) == 0 {
			err := fmt.Errorf("El documento '%s' aún no está asignado a un registro en Terceros", doc)
			outputError = map[string]interface{}{
				"funcion": "GetTerceroByDoc - len(datosTerceros) == 1 ",
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
			err := fmt.Errorf("El Documento '%s' tiene más de un registro en Terceros (%d registros%s)", doc, q, s)
			logs.Warn(err)
			outputError = map[string]interface{}{
				"funcion": "GetTerceroByDoc - len(datosTerceros) == 1 && datosTerceros[0].TerceroId != nil",
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

}

// GetTerceroById get controlador tercero/{id} del api terceros_crud
func GetTerceroById(id int) (tercero *models.Tercero, outputError map[string]interface{}) {

	funcion := "GetTerceroById"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("tercerosService") + "tercero/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &tercero); err != nil {
		eval := " - request.GetJson(urlcrud, &tercero)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}
	return tercero, nil
}
