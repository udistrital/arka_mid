package actaRecibido

import (
	"context"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	errorCtrl "github.com/udistrital/utils_oas/v2/errorctrl"
	requestV2 "github.com/udistrital/utils_oas/v2/request"
)

var path, _ = beego.AppConfig.String("actaRecibidoService")

// GetElementoById consulta controlador elemento/{id} del api acta_recibido_crud
func GetElementoById(ctx context.Context, id int, elemento *models.Elemento) (outputError map[string]interface{}) {

	funcion := "GetElementoById - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := path + "elemento/" + strconv.Itoa(id)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, elemento); err != nil {
		logs.Error(urlcrud+", ", err)
		eval := "requestV2.GetWithContext(ctx, urlcrud, elemento)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetAllElemento query controlador elemento del api acta_recibido_crud
func GetAllElemento(ctx context.Context, query string, fields string, sortby string, order string, offset string, limit string) (elementos []*models.Elemento, outputError map[string]interface{}) {

	funcion := "GetAllElemento - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := path + "elemento?" + utilsHelper.EncodeUrl(query, fields, sortby, order, offset, limit)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &elementos); err != nil {
		logs.Error(urlcrud+", ", err)
		eval := "requestV2.GetWithContext(ctx, urlcrud, &elementos)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return elementos, nil
}

// GetAllHistoricoActa query controlador historico_acta del api acta_recibido_crud
func GetAllHistoricoActa(ctx context.Context, query string, fields string, sortby string, order string, offset string, limit string) (historicos []models.HistoricoActa, outputError map[string]interface{}) {

	funcion := "GetAllHistoricoActa - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := path + "historico_acta?" + utilsHelper.EncodeUrl(query, fields, sortby, order, offset, limit)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &historicos); err != nil {
		logs.Error(urlcrud+", ", err)
		eval := "requestV2.GetWithContext(ctx, urlcrud, &historicos)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return historicos, nil
}

// GetAllHistoricoActas query controlador historico_acta del api acta_recibido_crud teniendo el cuenta el número de registros totales
func GetAllHistoricoActas(ctx context.Context, query string, fields string, sortby string, order string, offset string, limit string) (historicos []*models.HistoricoActa, count string, outputError map[string]interface{}) {

	funcion := "GetAllHistoricoActas - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := path + "historico_acta?" + utilsHelper.EncodeUrl(query, fields, sortby, order, offset, limit)
	_, total, err := requestV2.GetWithTotalCount(ctx, urlcrud, &historicos)
	if err != nil {
		logs.Error(urlcrud+", ", err)
		eval := "requestV2.GetWithTotalCount(ctx, urlcrud, &historicos)"
		return nil, "", errorCtrl.Error(funcion+eval, err, "502")
	}

	// GetWithTotalCount returns zero when Total-Count is absent or invalid.
	if total != 0 {
		count = strconv.Itoa(total)
	}
	return
}

// GetAllActaRecibido query controlador acta_recibido del api acta_recibido_crud
func GetAllActaRecibido(ctx context.Context, payload string) (actas []models.ActaRecibido, outputError map[string]interface{}) {

	funcion := "GetAllActaRecibido - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := path + "acta_recibido?" + payload
	_, err := requestV2.GetWithContext(ctx, urlcrud, &actas)
	if err != nil {
		logs.Error(urlcrud, err)
		eval := "requestV2.GetWithContext(ctx, urlcrud, &actas)"
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetAllCampo query controlador acta_recibido del api acta_recibido_crud
func GetAllCampo(ctx context.Context, payload string) (campos []models.Campo, outputError map[string]interface{}) {

	funcion := "GetAllCampo - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := path + "campo?" + payload
	_, err := requestV2.GetWithContext(ctx, urlcrud, &campos)
	if err != nil {
		logs.Error(urlcrud, err)
		eval := "requestV2.GetWithContext(ctx, urlcrud, &campos)"
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// PutElemento put controlador elemento del api acta_recibido_crud
func PutElemento(ctx context.Context, elemento *models.Elemento, elementoId int) (outputError map[string]interface{}) {

	funcion := "PutElemento - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := path + "elemento/" + strconv.Itoa(elementoId)
	_, err := requestV2.PutWithContext(ctx, urlcrud, elemento, elemento)
	if err != nil {
		logs.Error(urlcrud, err)
		eval := `requestV2.PutWithContext(ctx, urlcrud, elemento, elemento)`
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return

}

// PutElementoCampo put controlador elemento del api acta_recibido_crud
func PutElementoCampo(ctx context.Context, elemento *models.ElementoCampo, elementoId int) (outputError map[string]interface{}) {

	funcion := "PutElementoCampo - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := path + "elemento_campo/" + strconv.Itoa(elementoId)
	_, err := requestV2.PutWithContext(ctx, urlcrud, elemento, elemento)
	if err != nil {
		logs.Error(urlcrud, err)
		eval := `requestV2.PutWithContext(ctx, urlcrud, elemento, elemento)`
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetSoporteById query controlador soporte_acta del api acta_recibido_crud
func GetSoporteById(ctx context.Context, id int, soporte *models.SoporteActa) (outputError map[string]interface{}) {

	funcion := "GetSoporteById"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := path + "soporte_acta/" + strconv.Itoa(id)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, soporte); err != nil {
		logs.Error(urlcrud+", ", err)
		eval := " - requestV2.GetWithContext(ctx, urlcrud, soporte)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetTransaccionActaRecibidoById consulta controlador transaccion_acta_recibido/{id} del api acta_recibido_crud
func GetTransaccionActaRecibidoById(ctx context.Context, id int, elementos bool, transaccion *models.TransaccionActaRecibido) (outputError map[string]interface{}) {

	funcion := "GetTransaccionActaRecibidoById - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := path + "transaccion_acta_recibido/" + strconv.Itoa(id) + "?elementos=" + strconv.FormatBool(elementos)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, transaccion); err != nil {
		logs.Error(err)
		eval := `requestV2.GetWithContext(ctx, urlcrud, transaccion)`
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// PutTransaccionActaRecibido put controlador transaccion_acta_recibido/{id} del api acta_recibido_crud
func PutTransaccionActaRecibido(ctx context.Context, id int, transaccion *models.TransaccionActaRecibido) (outputError map[string]interface{}) {

	funcion := "PutTransaccionActaRecibido - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := path + "transaccion_acta_recibido/" + strconv.Itoa(id)
	if _, err := requestV2.PutWithContext(ctx, urlcrud, transaccion, transaccion); err != nil {
		logs.Error(err)
		eval := `requestV2.PutWithContext(ctx, urlcrud, transaccion, transaccion)`
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetAllElementoCampo query controlador historico_acta del api acta_recibido_crud teniendo el cuenta el número de registros totales
func GetAllElementoCampo(ctx context.Context, payload string) (elementosCampo []models.ElementoCampo, outputError map[string]interface{}) {

	funcion := "GetAllElementoCampo - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := path + "elemento_campo?" + payload
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &elementosCampo); err != nil {
		logs.Error(err)
		eval := `requestV2.GetWithContext(ctx, urlcrud, &elementosCampo)`
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// PostActaRecibido post controlador acta_recibido del api acta_recibido_crud
func PostActaRecibido(ctx context.Context, acta *models.ActaRecibido) (outputError map[string]interface{}) {

	funcion := "PostActaRecibido - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := path + "acta_recibido/"
	_, err := requestV2.PostWithContext(ctx, urlcrud, acta, acta)
	if err != nil {
		logs.Error(err, urlcrud)
		eval := `requestV2.PostWithContext(ctx, urlcrud, acta, acta)`
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// PostElemento post controlador elemento del api acta_recibido_crud
func PostElemento(ctx context.Context, elemento *models.Elemento) (outputError map[string]interface{}) {

	funcion := "PostElemento - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := path + "elemento/"
	_, err := requestV2.PostWithContext(ctx, urlcrud, elemento, elemento)
	if err != nil {
		logs.Error(err, urlcrud)
		eval := `requestV2.PostWithContext(ctx, urlcrud, elemento, elemento)`
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// PostElementoCampo post controlador elemento_campo del api acta_recibido_crud
func PostElementoCampo(ctx context.Context, elemento *models.ElementoCampo) (outputError map[string]interface{}) {

	funcion := "PostElementoCampo - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := path + "elemento_campo/"
	_, err := requestV2.PostWithContext(ctx, urlcrud, elemento, elemento)
	if err != nil {
		logs.Error(err, urlcrud)
		eval := `requestV2.PostWithContext(ctx, urlcrud, elemento, elemento)`
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}
