package terceros

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	"github.com/udistrital/arka_mid/utils_oas/request"
)

var basePath, _ = beego.AppConfig.String("tercerosMidService")

func GetCargoFuncionario(id int) (cargo []*models.Parametro, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetCargoFuncionario", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	// Consulta cargo
	urlcrud := basePath + "propiedad/cargo/" + strconv.Itoa(id)
	req, err := http.NewRequest(http.MethodGet, urlcrud, nil)
	if err != nil {
		logs.Error(err)
		return nil, map[string]interface{}{
			"funcion": "GetCargoFuncionario - http.NewRequest(http.MethodGet, urlcrud, nil)",
			"err":     err,
			"status":  "502",
		}
	}

	req.Header.Set("Authorization", request.GetHeader())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logs.Error(err)
		return nil, map[string]interface{}{
			"funcion": "GetCargoFuncionario - http.DefaultClient.Do(req)",
			"err":     err,
			"status":  "502",
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err)
		return nil, map[string]interface{}{
			"funcion": "GetCargoFuncionario - io.ReadAll(resp.Body)",
			"err":     err,
			"status":  "502",
		}
	}

	if resp.StatusCode == http.StatusNotFound {
		return []*models.Parametro{}, nil
	}

	if resp.StatusCode >= http.StatusBadRequest {
		serviceError := map[string]interface{}{
			"status": resp.StatusCode,
			"body":   string(body),
		}
		logs.Error(serviceError)
		return nil, map[string]interface{}{
			"funcion": "GetCargoFuncionario - servicio terceros_mid",
			"err":     serviceError,
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
func GetDocumentoTercero(id int) (documento []*models.DatosIdentificacion, outputError map[string]interface{}) {

	funcion := "GetDocumentoTercero"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	// Consulta documento
	urlcrud := basePath + "propiedad/documento/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &documento); err != nil {
		eval := " - request.GetJson(urlcrud, &documento)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}
