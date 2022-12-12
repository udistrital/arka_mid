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

//GetNombreTerceroById trae el nombre de un encargado por su id
func GetNombreTerceroById(idTercero int) (tercero *models.IdentificacionTercero, outputError map[string]interface{}) {

	funcion := "GetNombreTerceroById"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	if idTercero <= 0 {
		err := errors.New("el idTercero debe ser mayor a 0")
		logs.Error(err)
		eval := " - strconv.Atoi(idTercero)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	urlcrud := "?limit=1&sortby=TipoDocumentoId&order=desc&query=Activo:true,TerceroId__Id:" + strconv.Itoa(idTercero)
	if datosId, err := GetAllDatosIdentificacion(urlcrud); err != nil {
		return nil, err
	} else {
		tercero = new(models.IdentificacionTercero)
		if len(datosId) == 0 || datosId[0].Id == 0 {
			if tercero_, err := GetTerceroById(idTercero); err != nil {
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
	urltercero := "http://" + beego.AppConfig.String("tercerosService") + "datos_identificacion"
	urltercero += "?query=Activo:true,TerceroId__Activo:true,Numero:" + doc
	// TODO: Alternativamente, se podría usar una de las siguientes:
	// urltercero += "?limit=1&sortby=TipoDocumentoId&order=desc&query=Activo:true,TerceroId__Id:" + strconv.Itoa(idTercero)
	// urltercero += "?limit=1&sortby=TipoDocumentoId&order=desc&query=Activo:true,Numero:" + doc
	// PERO
	// Depende de que en terceros_crud se agregue a la tabla datos_identificacion un CONSTRAINT UNIQUE
	// entre tipo_documento_id y numero. Una vez agregado dicho constraint, quizás las validaciones
	// a continuación se podrían prescindir
	var terceros []models.DatosIdentificacion
	// logs.Debug("urltercero:", urltercero)
	if resp, err := request.GetJsonTest(urltercero, &terceros); err == nil && resp.StatusCode == 200 {
		if len(terceros) == 1 {
			return &terceros[0], nil
		} else if len(terceros) == 0 {
			err := fmt.Errorf("el documento '%s' aún no está asignado a un registro en Terceros", doc)
			logs.Notice(err)
			outputError = map[string]interface{}{
				"funcion": "GetTerceroByDoc - len(datosTerceros) == 1 ",
				"err":     err,
				"status":  "404",
			}
			return nil, outputError
		} else { // len(terceros) > 1
			q := len(terceros)
			if q > 1 && DocumentosValidos(terceros, false, true) {
				return &terceros[0], nil
			}
			s := ""
			if q >= 10 {
				s = " - o más"
			}
			err := fmt.Errorf("el Documento '%s' tiene más de un registro activo en Terceros (%d registros%s)", doc, q, s)
			logs.Notice(err)
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
			"funcion": "GetTerceroByDoc - request.GetJsonTest(urltercero, &datosTerceros)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

}

func GetTerceroUD() (int, map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GetTerceroUD - Unhandled Error!", "500")

	payload := "query=TipoDocumentoId__Nombre:NIT,Numero:" + GetDocUD()
	if tercero, err := GetAllDatosIdentificacion(payload); err != nil {
		return 0, err
	} else if len(tercero) > 0 && tercero[0].TerceroId.Id > 0 {
		return tercero[0].TerceroId.Id, nil
	} else {
		return 0, nil
	}
}
