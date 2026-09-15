package cuentasContables

import (
	"context"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	requestV2 "github.com/udistrital/utils_oas/v2/request"
)

var basePath, _ = beego.AppConfig.String("cuentasContablesService")

// GetCuentaContable Consulta controlador nodo_cuenta_contable/{UUID}
func GetCuentaContable(ctx context.Context, cuentaContableId string) (cuentaContable *models.CuentaContable, outputError map[string]interface{}) {

	funcion := "GetCuentaContable"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "nodo_cuenta_contable/" + cuentaContableId
	var data models.RespuestaAPI2obj
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &data); err != nil || data.Code != 200 {
		if data.Message == "document-no-found" {
			return nil, nil
		}
		eval := " - requestV2.GetWithContext(ctx, urlcrud, &data)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	} else {
		outputError = utilsHelper.FillStruct(data.Body, &cuentaContable)
		return
	}
}

// GetTipoComprobante Consulta controlador tipo_comprobante/{UUID}
func GetTipoComprobante(ctx context.Context, tipoDocumento string) (tipoComprobante *models.TipoComprobante, outputError map[string]interface{}) {

	funcion := "GetTipoComprobante"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var data models.RespuestaAPI2arr
	urlcrud := basePath + "tipo_comprobante"
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &data); err != nil || data.Code != 200 {
		eval := " - requestV2.GetWithContext(ctx, urlcrud, &data)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	var tiposComprobante []*models.TipoComprobante
	outputError = utilsHelper.FillStruct(data.Body, &tiposComprobante)
	if outputError != nil {
		return
	}

	for _, tipoComprobante := range tiposComprobante {
		if tipoComprobante.TipoDocumento == tipoDocumento {
			return tipoComprobante, nil
		}

	}

	return nil, nil
}

// GetComprobante Retorna el id para un tipo de comprobante dato
func GetComprobante(ctx context.Context, tipoDocumento string, id *string) (outputError map[string]interface{}) {

	funcion := "GetComprobante"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		data         models.RespuestaAPI2arr
		comprobantes []*models.Comprobante
	)

	urlcrud := basePath + "comprobante"
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &data); err != nil || data.Code != 200 {
		eval := " - requestV2.GetWithContext(ctx, urlcrud, &data)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	outputError = utilsHelper.FillStruct(data.Body, &comprobantes)
	if outputError != nil {
		return
	}

	for _, comprobante := range comprobantes {
		if comprobante.TipoComprobante.TipoDocumento+strconv.Itoa(comprobante.Numero) == tipoDocumento {
			*id = comprobante.Id
			return
		}
	}

	return

}
