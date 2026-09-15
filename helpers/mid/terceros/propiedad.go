package terceros

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	requestV2 "github.com/udistrital/utils_oas/v2/request"
)

var basePath, _ = beego.AppConfig.String("tercerosMidService")

func GetCargoFuncionario(ctx context.Context, id int) (cargo []*models.Parametro, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetCargoFuncionario", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	// Consulta cargo
	urlcrud := basePath + "propiedad/cargo/" + strconv.Itoa(id)
	var body json.RawMessage
	status, err := requestV2.GetWithContext(ctx, urlcrud, &body)
	if status == 404 {
		return []*models.Parametro{}, nil
	}
	if err != nil && errors.Is(err, io.EOF) {
		return []*models.Parametro{}, nil
	}
	if err != nil {
		logs.Error(err)
		return nil, map[string]interface{}{
			"funcion": "GetCargoFuncionario - servicio terceros_mid",
			"err":     err,
			"status":  "502",
		}
	}

	if len(body) == 0 {
		return []*models.Parametro{}, nil
	}

	if err := json.Unmarshal(body, &cargo); err == nil {
		return cargo, nil
	}

	var parametro models.Parametro
	if err := json.Unmarshal(body, &parametro); err == nil {
		return []*models.Parametro{&parametro}, nil
	}

	logs.Error(err)
	outputError = map[string]interface{}{
		"funcion": "GetCargoFuncionario - json.Unmarshal(body, &cargo)",
		"err":     err,
		"status":  "502",
	}
	return
}

// GetDocumentoTercero get controlador propiedad/documento/{id} del api terceros_mid
func GetDocumentoTercero(ctx context.Context, id int) (documento []*models.DatosIdentificacion, outputError map[string]interface{}) {

	funcion := "GetDocumentoTercero"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	// Consulta documento
	urlcrud := basePath + "propiedad/documento/" + strconv.Itoa(id)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &documento); err != nil {
		eval := " - requestV2.GetWithContext(ctx, urlcrud, &documento)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}
