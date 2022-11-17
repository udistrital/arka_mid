package actaRecibido

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

var path = "http://" + beego.AppConfig.String("actaRecibidoService")

// GetElementoById consulta controlador elemento/{id} del api acta_recibido_crud
func GetElementoById(id int, elemento *models.Elemento) (outputError map[string]interface{}) {

	funcion := "GetElementoById - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := path + "elemento/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &elemento); err != nil {
		logs.Error(urlcrud+", ", err)
		eval := "request.GetJson(urlcrud, &elemento)"
		return errorctrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetAllElemento query controlador elemento del api acta_recibido_crud
func GetAllElemento(query string, fields string, sortby string, order string, offset string, limit string) (elementos []*models.Elemento, outputError map[string]interface{}) {

	funcion := "GetAllElemento - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := path + "elemento?" + utilsHelper.EncodeUrl(query, fields, sortby, order, offset, limit)
	if err := request.GetJson(urlcrud, &elementos); err != nil {
		logs.Error(urlcrud+", ", err)
		eval := "request.GetJson(urlcrud, &elementos)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}

	return elementos, nil
}

// GetAllHistoricoActa query controlador historico_acta del api acta_recibido_crud
func GetAllHistoricoActa(query string, fields string, sortby string, order string, offset string, limit string) (historicos []*models.HistoricoActa, outputError map[string]interface{}) {

	funcion := "GetAllHistoricoActa - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := path + "historico_acta?" + utilsHelper.EncodeUrl(query, fields, sortby, order, offset, limit)
	if err := request.GetJson(urlcrud, &historicos); err != nil {
		logs.Error(urlcrud+", ", err)
		eval := "request.GetJson(urlcrud, &historicos)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}

	return historicos, nil
}

// PutElemento put controlador elemento del api acta_recibido_crud
func PutElemento(elemento *models.Elemento, elementoId int) (elemento_ *models.Elemento, outputError map[string]interface{}) {

	funcion := "PutElemento - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := path + "elemento/" + strconv.Itoa(elementoId)
	if err := request.SendJson(urlcrud, "PUT", &elemento_, &elemento); err != nil {
		logs.Error(urlcrud+", ", err)
		eval := `request.SendJson(urlcrud, "PUT", &elemento_, &elemento)`
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}

	return

}

// GetSoporteById query controlador soporte_acta del api acta_recibido_crud
func GetSoporteById(id int, soporte *models.SoporteActa) (outputError map[string]interface{}) {

	funcion := "GetSoporteById"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := path + "soporte_acta/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &soporte); err != nil {
		logs.Error(urlcrud+", ", err)
		eval := " - request.GetJson(urlcrud, &soporte)"
		return errorctrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetTransaccionActaRecibidoById consulta controlador transaccion_acta_recibido/{id} del api acta_recibido_crud
func GetTransaccionActaRecibidoById(id int, elementos bool, transaccion *models.TransaccionActaRecibido) (outputError map[string]interface{}) {

	funcion := "GetTransaccionActaRecibidoById - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := path + "transaccion_acta_recibido/" + strconv.Itoa(id) + "?elementos=" + strconv.FormatBool(elementos)
	if err := request.GetJson(urlcrud, &transaccion); err != nil {
		logs.Error(err)
		eval := `request.GetJson(urlcrud, &transaccion)`
		return errorctrl.Error(funcion+eval, err, "502")
	}

	return
}

// PutTransaccionActaRecibido put controlador transaccion_acta_recibido/{id} del api acta_recibido_crud
func PutTransaccionActaRecibido(id int, transaccion *models.TransaccionActaRecibido) (outputError map[string]interface{}) {

	funcion := "PutTransaccionActaRecibido - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := path + "transaccion_acta_recibido/" + strconv.Itoa(id)
	if err := request.SendJson(urlcrud, "PUT", &transaccion, &transaccion); err != nil {
		logs.Error(err)
		eval := `request.SendJson(urlcrud, "PUT", &transaccion, &transaccion)`
		return errorctrl.Error(funcion+eval, err, "502")
	}

	return
}
