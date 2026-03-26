package actaRecibido

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	"github.com/udistrital/arka_mid/utils_oas/request"
)

var path, _ = beego.AppConfig.String("actaRecibidoService")

// GetElementoById consulta controlador elemento/{id} del api acta_recibido_crud
func GetElementoById(id int, elemento *models.Elemento) (outputError map[string]interface{}) {

	funcion := "GetElementoById - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + path + "elemento/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &elemento); err != nil {
		logs.Error(urlcrud+", ", err)
		eval := "request.GetJson(urlcrud, &elemento)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetAllElemento query controlador elemento del api acta_recibido_crud
func GetAllElemento(query string, fields string, sortby string, order string, offset string, limit string) (elementos []*models.Elemento, outputError map[string]interface{}) {

	funcion := "GetAllElemento - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + path + "elemento?" + utilsHelper.EncodeUrl(query, fields, sortby, order, offset, limit)
	if err := request.GetJson(urlcrud, &elementos); err != nil {
		logs.Error(urlcrud+", ", err)
		eval := "request.GetJson(urlcrud, &elementos)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return elementos, nil
}

// GetAllHistoricoActa query controlador historico_acta del api acta_recibido_crud
func GetAllHistoricoActa(query string, fields string, sortby string, order string, offset string, limit string) (historicos []models.HistoricoActa, outputError map[string]interface{}) {

	funcion := "GetAllHistoricoActa - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + path + "historico_acta?" + utilsHelper.EncodeUrl(query, fields, sortby, order, offset, limit)
	if err := request.GetJson(urlcrud, &historicos); err != nil {
		logs.Error(urlcrud+", ", err)
		eval := "request.GetJson(urlcrud, &historicos)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return historicos, nil
}

// GetAllHistoricoActas query controlador historico_acta del api acta_recibido_crud teniendo el cuenta el número de registros totales
func GetAllHistoricoActas(query string, fields string, sortby string, order string, offset string, limit string) (historicos []*models.HistoricoActa, count string, outputError map[string]interface{}) {

	funcion := "GetAllHistoricoActas - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + path + "historico_acta?" + utilsHelper.EncodeUrl(query, fields, sortby, order, offset, limit)
	response, err := request.GetJsonTest(urlcrud, &historicos)
	if err != nil {
		logs.Error(urlcrud+", ", err)
		eval := "request.GetJsonTest(urlcrud, &historicos)"
		return nil, "", errorCtrl.Error(funcion+eval, err, "502")
	}

	count = response.Header.Get("total-count")
	return
}

// GetAllActaRecibido query controlador acta_recibido del api acta_recibido_crud
func GetAllActaRecibido(payload string) (actas []models.ActaRecibido, outputError map[string]interface{}) {

	funcion := "GetAllActaRecibido - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + path + "acta_recibido?" + payload
	err := request.GetJson(urlcrud, &actas)
	if err != nil {
		logs.Error(urlcrud, err)
		eval := "request.GetJson(urlcrud, &actas)"
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetAllCampo query controlador acta_recibido del api acta_recibido_crud
func GetAllCampo(payload string) (campos []models.Campo, outputError map[string]interface{}) {

	funcion := "GetAllCampo - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + path + "campo?" + payload
	err := request.GetJson(urlcrud, &campos)
	if err != nil {
		logs.Error(urlcrud, err)
		eval := "request.GetJson(urlcrud, &campos)"
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// PutElemento put controlador elemento del api acta_recibido_crud
func PutElemento(elemento *models.Elemento, elementoId int) (outputError map[string]interface{}) {

	funcion := "PutElemento - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + path + "elemento/" + strconv.Itoa(elementoId)
	err := request.SendJson(urlcrud, "PUT", &elemento, &elemento)
	if err != nil {
		logs.Error(urlcrud, err)
		eval := `request.SendJson(urlcrud, "PUT", &elemento, &elemento)`
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return

}

// PutElementoCampo put controlador elemento del api acta_recibido_crud
func PutElementoCampo(elemento *models.ElementoCampo, elementoId int) (outputError map[string]interface{}) {

	funcion := "PutElementoCampo - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + path + "elemento_campo/" + strconv.Itoa(elementoId)
	err := request.SendJson(urlcrud, "PUT", &elemento, &elemento)
	if err != nil {
		logs.Error(urlcrud, err)
		eval := `request.SendJson(urlcrud, "PUT", &elemento, &elemento)`
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetSoporteById query controlador soporte_acta del api acta_recibido_crud
func GetSoporteById(id int, soporte *models.SoporteActa) (outputError map[string]interface{}) {

	funcion := "GetSoporteById"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + path + "soporte_acta/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &soporte); err != nil {
		logs.Error(urlcrud+", ", err)
		eval := " - request.GetJson(urlcrud, &soporte)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetTransaccionActaRecibidoById consulta controlador transaccion_acta_recibido/{id} del api acta_recibido_crud
func GetTransaccionActaRecibidoById(id int, elementos bool, transaccion *models.TransaccionActaRecibido) (outputError map[string]interface{}) {

	funcion := "GetTransaccionActaRecibidoById - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + path + "transaccion_acta_recibido/" + strconv.Itoa(id) + "?elementos=" + strconv.FormatBool(elementos)
	if err := request.GetJson(urlcrud, &transaccion); err != nil {
		logs.Error(err)
		eval := `request.GetJson(urlcrud, &transaccion)`
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// PutTransaccionActaRecibido put controlador transaccion_acta_recibido/{id} del api acta_recibido_crud
func PutTransaccionActaRecibido(id int, transaccion *models.TransaccionActaRecibido) (outputError map[string]interface{}) {

	funcion := "PutTransaccionActaRecibido - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + path + "transaccion_acta_recibido/" + strconv.Itoa(id)
	if err := request.SendJson(urlcrud, "PUT", &transaccion, &transaccion); err != nil {
		logs.Error(err)
		eval := `request.SendJson(urlcrud, "PUT", &transaccion, &transaccion)`
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetAllElementoCampo query controlador historico_acta del api acta_recibido_crud teniendo el cuenta el número de registros totales
func GetAllElementoCampo(payload string) (elementosCampo []models.ElementoCampo, outputError map[string]interface{}) {

	funcion := "GetAllElementoCampo - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + path + "elemento_campo?" + payload
	if err := request.GetJson(urlcrud, &elementosCampo); err != nil {
		logs.Error(err)
		eval := `request.GetJson(urlcrud, &elementosCampo)`
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// PostActaRecibido post controlador acta_recibido del api acta_recibido_crud
func PostActaRecibido(acta *models.ActaRecibido) (outputError map[string]interface{}) {

	funcion := "PostActaRecibido - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := "http://" + path + "acta_recibido/"
	err := request.SendJson(urlcrud, "POST", &acta, &acta)
	if err != nil {
		logs.Error(err, urlcrud)
		eval := `request.SendJson(urlcrud, "POST", &acta, &acta)`
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// PostElemento post controlador elemento del api acta_recibido_crud
func PostElemento(elemento *models.Elemento) (outputError map[string]interface{}) {

	funcion := "PostElemento - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := "http://" + path + "elemento/"
	err := request.SendJson(urlcrud, "POST", &elemento, &elemento)
	if err != nil {
		logs.Error(err, urlcrud)
		eval := `request.SendJson(urlcrud, "POST", &elemento, &elemento)`
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// PostElementoCampo post controlador elemento_campo del api acta_recibido_crud
func PostElementoCampo(elemento *models.ElementoCampo) (outputError map[string]interface{}) {

	funcion := "PostElementoCampo - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := "http://" + path + "elemento_campo/"
	err := request.SendJson(urlcrud, "POST", &elemento, &elemento)
	if err != nil {
		logs.Error(err, urlcrud)
		eval := `request.SendJson(urlcrud, "POST", &elemento, &elemento)`
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}
