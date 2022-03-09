package actaRecibido

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

// GetElementoById consulta controlador elemento/{id} del api acta_recibido_crud
func GetElementoById(id int) (elemento *models.Elemento, outputError map[string]interface{}) {

	funcion := "GetElementoById"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &elemento); err != nil {
		eval := " - request.GetJson(urlcrud, &elemento)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}

	return elemento, nil
}

// GetAllElemento query controlador elemento del api acta_recibido_crud
func GetAllElemento(query string, fields string, sortby string, order string, offset string, limit string) (elementos []*models.Elemento, outputError map[string]interface{}) {

	funcion := "GetAllElemento"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento?" + utilsHelper.EncodeUrl(query, fields, sortby, order, offset, limit)
	if err := request.GetJson(urlcrud, &elementos); err != nil {
		eval := " - request.GetJson(urlcrud, &elementos)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}

	return elementos, nil
}

// GetAllHistoricoActa query controlador historico_acta del api acta_recibido_crud
func GetAllHistoricoActa(query string, fields string, sortby string, order string, offset string, limit string) (historicos []*models.HistoricoActa, outputError map[string]interface{}) {

	funcion := "GetAllHistoricoActa"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("actaRecibidoService") + "historico_acta?" + utilsHelper.EncodeUrl(query, fields, sortby, order, offset, limit)
	if err := request.GetJson(urlcrud, &historicos); err != nil {
		eval := " - request.GetJson(urlcrud, &historicos)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}

	return historicos, nil
}

// PutElemento put controlador elemento del api acta_recibido_crud
func PutElemento(elemento *models.Elemento, elementoId int) (elemento_ *models.Elemento, outputError map[string]interface{}) {

	funcion := "PutElemento"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento/" + strconv.Itoa(elementoId)
	if err := request.SendJson(urlcrud, "PUT", &elemento_, &elemento); err != nil {
		eval := ` - request.SendJson(urlcrud, "PUT", &elemento_, &elemento)`
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}

	return elemento_, nil

}
