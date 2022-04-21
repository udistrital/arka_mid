package terceros

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

type IdentificacionTercero struct {
	Id             int
	Numero         string
	NombreCompleto string
}

//GetNombreTerceroById trae el nombre de un encargado por su id
func GetNombreTerceroById(idTercero string) (tercero *IdentificacionTercero, outputError map[string]interface{}) {

	funcion := "GetNombreTerceroById"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var terceroId int
	if v, err := strconv.Atoi(idTercero); err != nil || v <= 0 {
		if err == nil {
			err = errors.New("el idTercero debe ser mayor a 0")
		}
		logs.Error(err)
		eval := " - strconv.Atoi(idTercero)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	} else {
		terceroId = v
	}

	urlcrud := "?limit=1&sortby=TipoDocumentoId&order=desc&query=Activo:true,TerceroId__Id:" + idTercero
	if datosId, err := GetAllDatosIdentificacion(urlcrud); err != nil {
		return nil, err
	} else {
		tercero = new(IdentificacionTercero)
		if len(datosId) == 0 || datosId[0].Id == 0 {
			if tercero_, err := GetTerceroById(terceroId); err != nil {
				return nil, err
			} else {
				tercero.Id = tercero_.Id
				tercero.NombreCompleto = tercero_.NombreCompleto
			}
			return tercero, nil
		}

		tercero.Id = datosId[0].TerceroId.Id
		tercero.Numero = datosId[0].Numero
		tercero.NombreCompleto = datosId[0].TerceroId.NombreCompleto
		return tercero, nil
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
			err := fmt.Errorf("el usuario '%s' aún no está asignado a un registro en Terceros", usuario)
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
			err := fmt.Errorf("el usuario '%s' tiene más de un registro en Terceros (%d registros%s)", usuario, q, s)
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
			err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
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
			err := fmt.Errorf("el documento '%s' aún no está asignado a un registro en Terceros", doc)
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
			err := fmt.Errorf("el Documento '%s' tiene más de un registro en Terceros (%d registros%s)", doc, q, s)
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
			err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
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
