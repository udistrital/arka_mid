package terceros

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	requestV2 "github.com/udistrital/utils_oas/v2/request"
)

var path, _ = beego.AppConfig.String("tercerosService")

// GetNombreTerceroById trae el nombre de un encargado por su id
func GetNombreTerceroById(ctx context.Context, idTercero int) (tercero *models.IdentificacionTercero, outputError map[string]interface{}) {

	funcion := "GetNombreTerceroById"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	if idTercero <= 0 {
		err := errors.New("el idTercero debe ser mayor a 0")
		logs.Error(err)
		eval := " - strconv.Atoi(idTercero)"
		return nil, errorCtrl.Error(funcion+eval, err, "500")
	}

	urlcrud := "limit=1&sortby=TipoDocumentoId&order=desc&query=Activo:true,TerceroId__Id:" + strconv.Itoa(idTercero) + ",TipoDocumentoId__Id__in:3|6|7"
	if datosId, err := GetAllDatosIdentificacion(ctx, urlcrud); err != nil {
		return nil, err
	} else {
		tercero = new(models.IdentificacionTercero)
		if len(datosId) == 0 || datosId[0].Id == 0 {
			urltercero := basePath + "tercero/" + strconv.Itoa(idTercero)
			tercero_ := new(models.Tercero)
			if _, err := requestV2.GetWithContext(ctx, urltercero, tercero_); err != nil {
				eval := " - requestV2.GetWithContext(ctx, urltercero, tercero_)"
				return nil, errorCtrl.Error(funcion+eval, err, "502")
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
func GetTerceroByUsuarioWSO2(ctx context.Context, usuario string) (tercero map[string]interface{}, outputError map[string]interface{}) {

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
	urltercero := path + "tercero"
	urltercero += "?fields=Id,NombreCompleto,TipoContribuyenteId"
	urltercero += "&query=Activo:true,UsuarioWSO2:" + usuario
	// logs.Info(urltercero)
	if status, err := requestV2.GetWithContext(ctx, urltercero, &terceros); err == nil && status == 200 {
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
			err = fmt.Errorf("undesired status code: %d", status)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetTerceroByUsuarioWSO2 - requestV2.GetWithContext(ctx, urltercero, &terceros)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	return tercero, nil
}

func GetTerceroUD(ctx context.Context) (int, map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("GetTerceroUD - Unhandled Error!", "500")

	payload := "query=TipoDocumentoId__Nombre:NIT,Numero:" + GetDocUD()
	if tercero, err := GetAllDatosIdentificacion(ctx, payload); err != nil {
		return 0, err
	} else if len(tercero) > 0 && tercero[0].TerceroId.Id > 0 {
		return tercero[0].TerceroId.Id, nil
	} else {
		return 0, nil
	}
}
