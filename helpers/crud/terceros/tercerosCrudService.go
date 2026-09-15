package terceros

import (
	"context"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	requestV2 "github.com/udistrital/utils_oas/v2/request"
)

var basePath, _ = beego.AppConfig.String("tercerosService")

// GetCorreo Consulta el correo de un tercero
func GetCorreo(ctx context.Context, id int) (DetalleFuncionario []*models.InfoComplementariaTercero, outputError map[string]interface{}) {

	funcion := "GetCorreo"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		urlcrud string
		correo  []*models.InfoComplementariaTercero
	)

	// Consulta correo
	urlcrud = basePath + "info_complementaria_tercero?limit=1&fields=Dato&sortby=Id&order=desc"
	urlcrud += "&query=Activo%3Atrue,InfoComplementariaId__Nombre__icontains%3Acorreo,TerceroId__Id%3A" + strconv.Itoa(id)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &correo); err != nil {
		logs.Error(err)
		eval := " - requestV2.GetWithContext(ctx, urlcrud, &correo)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return correo, nil
}

// GetAllDatosIdentificacion get controlador datos_identificacion de api terceros_crud
func GetAllDatosIdentificacion(ctx context.Context, query string) (datosId []models.DatosIdentificacion, outputError map[string]interface{}) {

	funcion := "GetAllDatosIdentificacion"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	// Consulta correo
	query = strings.TrimPrefix(query, "?")
	urlcrud := basePath + "datos_identificacion?" + query
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &datosId); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllDatosIdentificacion - requestV2.GetWithContext(ctx, urlcrud, &datosId)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	return datosId, nil
}

// GetTerceroById get controlador tercero/{id} del api terceros_crud
func GetTerceroById(ctx context.Context, id int) (tercero *models.Tercero, outputError map[string]interface{}) {

	funcion := "GetTerceroById"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "tercero/" + strconv.Itoa(id)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &tercero); err != nil {
		eval := " - requestV2.GetWithContext(ctx, urlcrud, &tercero)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}
	return tercero, nil
}

// GetTrTerceroIdentificacionById get controlador tercero/{id} del api terceros_crud
func GetTrTerceroIdentificacionById(ctx context.Context, id int) (tercero models.DetalleTercero, outputError map[string]interface{}) {

	funcion := "GetTrTerceroIdentificacionById"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "tercero/identificacion/" + strconv.Itoa(id)
	_, err := requestV2.GetWithContext(ctx, urlcrud, &tercero)
	if err != nil {
		logs.Error(err, urlcrud)
		eval := " - requestV2.GetWithContext(ctx, urlcrud, &tercero)"
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetAllTrTerceroIdentificacion get controlador tercero/identificacion del api terceros_crud
func GetAllTrTerceroIdentificacion(ctx context.Context, payload string) (terceros []models.DetalleTercero, outputError map[string]interface{}) {

	funcion := "GetAllTrTerceroIdentificacion - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := basePath + "tercero/identificacion?" + payload
	_, err := requestV2.GetWithContext(ctx, urlcrud, &terceros)
	if err != nil {
		logs.Error(err, urlcrud)
		eval := "requestV2.GetWithContext(ctx, urlcrud, &terceros)"
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetDocUD Get documento de identificación UD
func GetDocUD() string {
	return "899999230"
}
