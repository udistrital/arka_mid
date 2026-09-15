package movimientosContables

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	requestV2 "github.com/udistrital/utils_oas/v2/request"
)

var getWithContext = requestV2.GetWithContext
var postWithContext = requestV2.PostWithContext

const maxPostTrContableRetries = 3

func getBasePath() string {
	basePath, _ := beego.AppConfig.String("movimientosContablesmidService")
	basePath = normalizeMovimientosContablesBasePath(basePath)

	return basePath
}

func normalizeMovimientosContablesBasePath(basePath string) string {
	basePath = strings.TrimSpace(basePath)

	if basePath != "" && !strings.HasPrefix(basePath, "http://") && !strings.HasPrefix(basePath, "https://") {
		basePath = "http://" + basePath
	}

	// movimientos_contables_mid responde 301 desde http hacia https y el POST termina degradándose.
	if strings.HasPrefix(basePath, "http://pruebasapi.intranetoas.udistrital.edu.co/") {
		basePath = "https://" + strings.TrimPrefix(basePath, "http://")
	}

	if basePath != "" && !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}

	return basePath
}

// GetTransaccion query controlador transaccion_movimientos del api movimientos_contables_mid
func GetTransaccion(ctx context.Context, id int, criteria string, detail bool) (transaccion *models.TransaccionMovimientos, outputError map[string]interface{}) {
	funcion := "GetTransaccion"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	basePath := getBasePath()
	urlcrud := basePath + "transaccion_movimientos/" + criteria + "/" + strconv.Itoa(id)
	if detail {
		urlcrud += "?detailed=true"
	}

	if _, err := getWithContext(ctx, urlcrud, &transaccion); err != nil {
		logs.Error("%s -> error requestV2.GetWithContext: %v, url=%s", funcion, err, urlcrud)
		eval := " - requestV2.GetWithContext(ctx, urlcrud, &transaccion)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return transaccion, nil
}

// PostTrContable post controlador transaccion_movimientos/transaccion_movimientos/ del api movimientos_contables_mid
func PostTrContable(ctx context.Context, tr *models.TransaccionMovimientos) (tr_ *models.TransaccionMovimientos, outputError map[string]interface{}) {
	return postTrContable(ctx, tr, 0)
}

func postTrContable(ctx context.Context, tr *models.TransaccionMovimientos, attempt int) (tr_ *models.TransaccionMovimientos, outputError map[string]interface{}) {
	funcion := "PostTrContable"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error", "500")

	var resp map[string]interface{}
	basePath := getBasePath()
	urlcrud := basePath + "transaccion_movimientos"

	logs.Info("==== INICIO %s ====", funcion)
	logs.Info("%s -> url=%s", funcion, urlcrud)
	logs.Info("%s -> payload=%+v", funcion, tr)

	if status, err := postWithContext(ctx, urlcrud, tr, &resp); err != nil {
		if shouldRetryPostTrContable(status, err, attempt) {
			time.Sleep(time.Duration(attempt+1) * 200 * time.Millisecond)
			logs.Warn("%s -> reintentando POST tras error transitorio: %v", funcion, err)
			return postTrContable(ctx, tr, attempt+1)
		}
		logs.Error("%s -> error requestV2.PostWithContext: %v", funcion, err)
		eval := ` - requestV2.PostWithContext(ctx, urlcrud, tr, &resp)`
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	logs.Info("%s -> respuesta cruda=%#v", funcion, resp)
	logs.Info("%s -> resp[Success]=%#v tipo=%T", funcion, resp["Success"], resp["Success"])
	logs.Info("%s -> resp[Data]=%#v tipo=%T", funcion, resp["Data"], resp["Data"])
	logs.Info("%s -> resp[Message]=%#v tipo=%T", funcion, resp["Message"], resp["Message"])
	logs.Info("%s -> resp[Status]=%#v tipo=%T", funcion, resp["Status"], resp["Status"])

	success, ok := resp["Success"].(bool)
	if !ok {
		logs.Error("%s -> resp[Success] no es bool o no vino. valor=%#v tipo=%T", funcion, resp["Success"], resp["Success"])
		eval := ` - resp["Success"] inválido`
		return nil, errorCtrl.Error(funcion+eval, fmt.Sprintf("campo Success inválido: %T", resp["Success"]), "502")
	}

	if !success {
		data := resp["Data"]

		if data == nil {
			logs.Error("%s -> resp[Data] vino nil", funcion)
			eval := ` - response Data nil`
			return nil, errorCtrl.Error(funcion+eval, "response.Data llegó nil", "502")
		}

		switch d := data.(type) {
		case string:
			logs.Error("%s -> resp[Data] string=%q", funcion, d)

			if strings.Contains(d, "invalid character") {
				logs.Error("%s -> se detectó 'invalid character', reintentando PostTrContable", funcion)
				return postTrContable(ctx, tr, attempt+1)
			}

			eval := ` - requestV2.PostWithContext(ctx, urlcrud, tr, &resp)`
			return nil, errorCtrl.Error(funcion+eval, d, "502")

		case map[string]interface{}:
			logs.Error("%s -> resp[Data] map=%#v", funcion, d)

			if errVal, exists := d["err"]; exists && errVal != nil {
				eval := ` - requestV2.PostWithContext(ctx, urlcrud, tr, &resp)`
				return nil, errorCtrl.Error(funcion+eval, errVal, "502")
			}

			// fallback: serializar mapa completo para verlo en logs/error
			raw, _ := json.Marshal(d)
			eval := ` - requestV2.PostWithContext(ctx, urlcrud, tr, &resp)`
			return nil, errorCtrl.Error(funcion+eval, string(raw), "502")

		default:
			logs.Error("%s -> resp[Data] tipo inesperado=%T valor=%#v", funcion, d, d)
			eval := ` - requestV2.PostWithContext(ctx, urlcrud, tr, &resp)`
			return nil, errorCtrl.Error(funcion+eval, fmt.Sprintf("tipo inesperado en response.Data: %T", d), "502")
		}
	}

	logs.Info("%s -> POST exitoso, retornando transacción", funcion)
	logs.Info("==== FIN %s ====", funcion)
	return tr, nil
}

func shouldRetryPostTrContable(status int, err error, attempt int) bool {
	if err == nil || attempt >= maxPostTrContableRetries {
		return false
	}

	return status == 404 || strings.Contains(err.Error(), "http 404:")
}
