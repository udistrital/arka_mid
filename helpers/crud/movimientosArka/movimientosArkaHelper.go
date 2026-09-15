package movimientosArka

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	requestV2 "github.com/udistrital/utils_oas/v2/request"
)

var basePath, _ = beego.AppConfig.String("movimientosArkaService")
var getAllMovimientoRequest = requestV2.GetWithTotalCount

// GetAllEstadoMovimiento query controlador estado_movimiento del api movimientos_arka_crud
func GetAllEstadoMovimiento(ctx context.Context, query string) (estados []*models.EstadoMovimiento, outputError map[string]interface{}) {

	funcion := "GetAllEstadoMovimiento"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error", "500")

	urlcrud := basePath + "estado_movimiento?" + query
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &estados); err != nil {
		eval := " - requestV2.GetWithContext(ctx, urlcrud, &estados)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetAllFormatoTipoMovimiento query controlador formato_tipo_movimiento del api movimientos_arka_crud
func GetAllFormatoTipoMovimiento(ctx context.Context, query string) (formatos []*models.FormatoTipoMovimiento, outputError map[string]interface{}) {

	funcion := "GetAllFormatoTipoMovimiento"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "formato_tipo_movimiento?" + query
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &formatos); err != nil {
		eval := " - requestV2.GetWithContext(ctx, urlcrud, &formatos)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}
	return formatos, nil
}

// GetAllElementosMovimiento query controlador elementos_movimiento del api movimientos_arka_crud
func GetAllElementosMovimiento(ctx context.Context, query string) (elementos []*models.ElementosMovimiento, outputError map[string]interface{}) {

	funcion := "GetAllElementosMovimiento"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "elementos_movimiento?" + query
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &elementos); err != nil {
		eval := " - requestV2.GetWithContext(ctx, urlcrud, &elementos)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}
	return elementos, nil
}

// GetAllSoporteMovimiento query controlador soporte_movimiento del api movimientos_arka_crud
func GetAllSoporteMovimiento(ctx context.Context, query string) (soportes []models.SoporteMovimiento, outputError map[string]interface{}) {

	funcion := "GetAllSoporteMovimiento"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "soporte_movimiento?" + query
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &soportes); err != nil {
		eval := " - requestV2.GetWithContext(ctx, urlcrud, &soportes)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}
	return soportes, nil
}

// GetAllMovimiento query controlador movimiento del api movimientos_arka_crud
func GetAllMovimiento(ctx context.Context, payload string) (movimientos []*models.Movimiento, count string, outputError map[string]interface{}) {

	funcion := "GetAllMovimiento - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := basePath + "movimiento?" + payload
	_, total, err := getAllMovimientoRequest(ctx, urlcrud, &movimientos)
	if err != nil {
		eval := "requestV2.GetWithTotalCount(ctx, urlcrud, &movimientos)"
		return nil, "", errorCtrl.Error(funcion+eval, err, "502")
	}
	// GetWithTotalCount returns zero when Total-Count is absent or invalid.
	if total != 0 {
		count = strconv.Itoa(total)
	}
	return movimientos, count, nil
}

// GetAllNovedadElemento query controlador novedad_elemento del api movimientos_arka_crud
func GetAllNovedadElemento(ctx context.Context, query string) (novedades []*models.NovedadElemento, outputError map[string]interface{}) {

	funcion := "GetAllNovedadElemento"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "novedad_elemento?" + query
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &novedades); err != nil {
		eval := " - requestV2.GetWithContext(ctx, urlcrud, &novedades)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}
	return
}

// GetMovimientoById consulta controlador movimiento/{id} del api movimientos_arka_crud
func GetMovimientoById(ctx context.Context, id int) (movimiento *models.Movimiento, outputError map[string]interface{}) {

	funcion := "GetMovimientoById"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error", "500")

	// Se consulta el movimiento
	urlcrud := basePath + "movimiento/" + strconv.Itoa(id)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &movimiento); err != nil {
		eval := " - requestV2.GetWithContext(ctx, urlcrud, &movimiento)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return movimiento, nil
}

// GetElementosMovimientoById consulta controlador elementos_movimiento/{id} del api movimientos_arka_crud
func GetElementosMovimientoById(ctx context.Context, id int, elemento *models.ElementosMovimiento) (outputError map[string]interface{}) {

	funcion := "GetElementosMovimientoById - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := basePath + "elementos_movimiento/" + strconv.Itoa(id)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &elemento); err != nil {
		logs.Error(err, urlcrud)
		eval := "requestV2.GetWithContext(ctx, urlcrud, &elemento)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetTrSalida consulta controlador tr_salida/{id} del api movimientos_arka_crud
func GetTrSalida(ctx context.Context, id int) (trSalida *models.TrSalida, outputError map[string]interface{}) {

	funcion := "GetTrSalida"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error", "500")

	// Se consulta el movimiento
	urlcrud := basePath + "tr_salida/" + strconv.Itoa(id)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &trSalida); err != nil {
		eval := " - requestV2.GetWithContext(ctx, urlcrud, &trSalida)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return trSalida, nil
}

// PostMovimiento post controlador movimiento del api movimientos_arka_crud
func PostMovimiento(ctx context.Context, movimiento *models.Movimiento) (outputError map[string]interface{}) {

	funcion := "PostMovimiento - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := basePath + "movimiento"

	var raw interface{}

	if _, err := requestV2.PostWithContext(ctx, urlcrud, movimiento, &raw); err != nil {
		logs.Error(err)
		eval := `requestV2.PostWithContext(ctx, urlcrud, movimiento, &raw)`
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	switch data := raw.(type) {
	case map[string]interface{}:
		b, err := json.Marshal(data)
		if err != nil {
			return errorCtrl.Error(funcion+"json.Marshal(data)", err, "502")
		}
		if err := json.Unmarshal(b, movimiento); err != nil {
			return errorCtrl.Error(funcion+"json.Unmarshal(b, movimiento)", err, "502")
		}

	case []interface{}:
		if len(data) == 0 {
			err := fmt.Errorf("response vacío: arreglo sin elementos")
			return errorCtrl.Error(funcion+"response vacío", err, "502")
		}

		// normalmente tomo el último elemento creado
		b, err := json.Marshal(data[len(data)-1])
		if err != nil {
			return errorCtrl.Error(funcion+"json.Marshal(data[len(data)-1])", err, "502")
		}
		if err := json.Unmarshal(b, movimiento); err != nil {
			return errorCtrl.Error(funcion+"json.Unmarshal(b, movimiento)", err, "502")
		}

	default:
		err := fmt.Errorf("tipo inesperado de respuesta: %T", raw)
		return errorCtrl.Error(funcion+"tipo inesperado", err, "502")
	}

	return
}

// PostSoporteMovimiento post controlador soporte_movimiento del api movimientos_arka_crud
func PostSoporteMovimiento(ctx context.Context, soporte *models.SoporteMovimiento) (outputError map[string]interface{}) {

	funcion := "PostSoporteMovimiento - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := basePath + "soporte_movimiento"
	if _, err := requestV2.PostWithContext(ctx, urlcrud, &soporte, &soporte); err != nil {
		logs.Error(err)
		eval := `requestV2.PostWithContext(ctx, urlcrud, &soporte, &soporte)`
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// PostElementosMovimiento post controlador elementos_movimiento del api movimientos_arka_crud
func PostElementosMovimiento(ctx context.Context, elemento *models.ElementosMovimiento) (outputError map[string]interface{}) {

	funcion := "PostElementosMovimiento - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := basePath + "elementos_movimiento"
	_, err := requestV2.PostWithContext(ctx, urlcrud, &elemento, &elemento)
	if err != nil {
		logs.Error(err, urlcrud)
		eval := `requestV2.PostWithContext(ctx, urlcrud, &elemento, &elemento)`
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// PutTrSalida put controlador tr_salida del api movimientos_arka_crud
func PutTrSalida(ctx context.Context, trSalida *models.SalidaGeneral) (trResultado *models.SalidaGeneral, outputError map[string]interface{}) {

	funcion := "PutTrSalida"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "tr_salida/"
	if _, err := requestV2.PutWithContext(ctx, urlcrud, trSalida, &trResultado); err != nil {
		eval := " - requestV2.PutWithContext(ctx, urlcrud, trSalida, &trResultado)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return trResultado, nil

}

// PostTrSalida post controlador tr_salida del api movimientos_arka_crud
func PostTrSalida(ctx context.Context, trSalida *models.SalidaGeneral) (outputError map[string]interface{}) {

	funcion := "PostTrSalida - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := basePath + "tr_salida"
	_, err := requestV2.PostWithContext(ctx, urlcrud, trSalida, trSalida)
	if err != nil {
		logs.Error(err)
		eval := `requestV2.PostWithContext(ctx, urlcrud, trSalida, trSalida)`
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// PutMovimiento put controlador movimiento del api movimientos_arka_crud
func PutMovimiento(ctx context.Context, movimiento *models.Movimiento, movimientoId int) (outputError map[string]interface{}) {

	funcion := "PutMovimiento"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "movimiento/" + strconv.Itoa(movimientoId)
	_, err := requestV2.PutWithContext(ctx, urlcrud, &movimiento, &movimiento)
	if err != nil {
		eval := `requestV2.PutWithContext(ctx, urlcrud, &movimiento, &movimiento)`
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// PutRevision put controlador bajas/ del api movimientos_arka_crud
func PutRevision(ctx context.Context, revision *models.TrRevisionBaja) (ids []int, outputError map[string]interface{}) {

	funcion := "PutRevision"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "bajas/"
	if _, err := requestV2.PutWithContext(ctx, urlcrud, &revision, &ids); err != nil {
		eval := " - requestV2.PutWithContext(ctx, urlcrud, &revision, &ids)"
		return nil, errorCtrl.Error(funcion+eval, err, "500")
	}

	return ids, nil

}

// PutSoporteMovimiento put controlador soporte_movimiento del api movimientos_arka_crud
func PutSoporteMovimiento(ctx context.Context, soporte *models.SoporteMovimiento, soporteId int) (soporteR *models.SoporteMovimiento, outputError map[string]interface{}) {

	funcion := "PutSoporteMovimiento"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "soporte_movimiento/" + strconv.Itoa(soporteId)
	if _, err := requestV2.PutWithContext(ctx, urlcrud, &soporte, &soporteR); err != nil {
		eval := " - requestV2.PutWithContext(ctx, urlcrud, &soporte, &soporteR)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return soporteR, nil

}

// PutElementosMovimiento put controlador elementos_movimiento del api movimientos_arka_crud
func PutElementosMovimiento(ctx context.Context, elementoM *models.ElementosMovimiento, elementoId int) (elementoM_ *models.ElementosMovimiento, outputError map[string]interface{}) {

	funcion := "PutElementosMovimiento"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "elementos_movimiento/" + strconv.Itoa(elementoId)
	if _, err := requestV2.PutWithContext(ctx, urlcrud, &elementoM, &elementoM_); err != nil {
		eval := ` - requestV2.PutWithContext(ctx, urlcrud, &elementoM, &elementoM_)`
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return elementoM_, nil

}

// PutNovedadElemento put controlador novedad_elemento del api movimientos_arka_crud
func PutNovedadElemento(ctx context.Context, novedad *models.NovedadElemento, novedadId int) (novedad_ *models.NovedadElemento, outputError map[string]interface{}) {

	funcion := "PutNovedadElemento"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "novedad_elemento/" + strconv.Itoa(novedadId)
	if _, err := requestV2.PutWithContext(ctx, urlcrud, &novedad, &novedad_); err != nil {
		eval := ` - requestV2.PutWithContext(ctx, urlcrud, &novedad, &novedad_)`
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return novedad_, nil

}

// PostNovedadElemento post controlador novedad_elemento del api movimientos_arka_crud
func PostNovedadElemento(ctx context.Context, novedad *models.NovedadElemento) (outputError map[string]interface{}) {

	funcion := "PostNovedadElemento - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := basePath + "novedad_elemento"
	_, err := requestV2.PostWithContext(ctx, urlcrud, &novedad, &novedad)
	if err != nil {
		logs.Error(err, urlcrud)
		eval := `requestV2.PostWithContext(ctx, urlcrud, &novedad, &novedad)`
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetElementosFuncionario query controlador elementos_movimiento/funcionario/{funcionarioId} del api movimientos_arka_crud
func GetElementosFuncionario(ctx context.Context, funcionarioId int) (movimientos []int, outputError map[string]interface{}) {

	funcion := "GetElementosFuncionario"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "elementos_movimiento/funcionario/" + strconv.Itoa(funcionarioId)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &movimientos); err != nil {
		eval := " - requestV2.GetWithContext(ctx, urlcrud, &movimientos)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}
	return movimientos, nil
}

// GetHistorialElemento query controlador elementos_movimiento/historial/{elementoId} del api movimientos_arka_crud
func GetHistorialElemento(ctx context.Context, elementoId int, final bool) (historial *models.Historial, outputError map[string]interface{}) {

	funcion := "GetHistorialElemento"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "elementos_movimiento/historial/" + strconv.Itoa(elementoId)
	urlcrud += "?final=" + strconv.FormatBool(final)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &historial); err != nil {
		eval := " - requestV2.GetWithContext(ctx, urlcrud, &historial)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}
	return historial, nil
}

// GetCorteDepreciacion query controlador cierre/?fechaCorte={fechaCorte} del api movimientos_arka_crud
func GetCorteDepreciacion(ctx context.Context, fechaCorte string) (corte []models.DepreciacionElemento, outputError map[string]interface{}) {

	funcion := "GetCorteDepreciacion - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := buildCorteDepreciacionURL(fechaCorte)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &corte); err != nil {
		logs.Error(err, urlcrud)
		eval := "requestV2.GetWithContext(ctx, urlcrud, &corte)"
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

func buildCorteDepreciacionURL(fechaCorte string) string {
	return basePath + "cierre?fechaCorte=" + fechaCorte
}

// AprobarCierre post controlador cierre del api movimientos_arka_crud
func AprobarCierre(ctx context.Context, cierre *models.Movimiento) (outputError map[string]interface{}) {

	funcion := "AprobarCierre - "
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error", "500")

	urlcrud := basePath + "cierre/"
	if _, err := requestV2.PostWithContext(ctx, urlcrud, &cierre, &cierre); err != nil {
		logs.Error(err, urlcrud)
		eval := `requestV2.PostWithContext(ctx, urlcrud, &cierre, &cierre)`
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetEntradaByActa consulta controlador movimiento/entrada/{acta_recibido_id} del api movimientos_arka_crud
func GetEntradaByActa(ctx context.Context, acta_recibido_id int) (entrada *models.Movimiento, outputError map[string]interface{}) {

	funcion := "GetEntradaByActa"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "movimiento/entrada/" + strconv.Itoa(acta_recibido_id)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &entrada); err != nil {
		eval := " - requestV2.GetWithContext(ctx, urlcrud, &entrada)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}
	return entrada, nil
}

// GetTrasladosByTerceroId consulta controlador movimiento/traslado/{tercero_id} del api movimientos_arka_crud
func GetTrasladosByTerceroId(ctx context.Context, terceroId int, confirmar bool, traslados *[]*models.Movimiento) (outputError map[string]interface{}) {

	funcion := "GetTrasladosByTerceroId - "
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := basePath + "movimiento/traslado/" + strconv.Itoa(terceroId)
	if confirmar {
		urlcrud += "?confirmar=true"
	}
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &traslados); err != nil {
		eval := "requestV2.GetWithContext(ctx, urlcrud, &traslados)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}
	return
}

// GetBajasByTerceroId consulta controlador movimiento/baja/{tercero_id} del api movimientos_arka_crud
func GetBajasByTerceroId(ctx context.Context, terceroId int, bajas *[]*models.Movimiento) (outputError map[string]interface{}) {

	funcion := "GetBajasByTerceroId - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := basePath + "movimiento/baja/" + strconv.Itoa(terceroId)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &bajas); err != nil {
		eval := "requestV2.GetWithContext(ctx, urlcrud, &bajas)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetBodegaByTerceroId consulta controlador movimiento/bodega/{tercero_id} del api movimientos_arka_crud
func GetBodegaByTerceroId(ctx context.Context, terceroId int, solicitudes *[]*models.Movimiento) (outputError map[string]interface{}) {

	funcion := "GetBodegaByTerceroId - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := basePath + "movimiento/bodega/" + strconv.Itoa(terceroId)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &solicitudes); err != nil {
		eval := "requestV2.GetWithContext(ctx, urlcrud, &solicitudes)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetAperturas consulta controlador tr_kardex/aperturas del api movimientos_arka_crud
func GetAperturas(ctx context.Context, conSaldo bool, aperturas *[]models.Apertura) (outputError map[string]interface{}) {

	funcion := "GetAperturas - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := basePath + "tr_kardex/aperturas?ConSaldo=" + strconv.FormatBool(conSaldo)
	if _, err := requestV2.GetWithContext(ctx, urlcrud, &aperturas); err != nil {
		eval := "requestV2.GetWithContext(ctx, urlcrud, &aperturas)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

func GetAllCentroCostos(ctx context.Context, payload string) (centroCostos []models.CentroCostos, outputError map[string]interface{}) {

	funcion := "GetAllCentroCostos - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := basePath + "centro_costos?" + normalizarConsultaCentroCostos(payload)

	_, err := requestV2.GetWithContext(ctx, urlcrud, &centroCostos)
	if err != nil {
		logs.Error(err)
		eval := "requestV2.GetWithContext(ctx, urlcrud, &centroCostos)"
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// normalizarConsultaCentroCostos evita el filtro exacto por Id del CRUD, que
// conserva un plan de PostgreSQL incompatible desde que Codigo pasó a ser
// texto. Id__in con un único valor tiene la misma semántica y retorna el
// arreglo esperado por este MID.
func normalizarConsultaCentroCostos(payload string) string {
	const filtroID = "query=Id:"
	if strings.HasPrefix(payload, filtroID) {
		return "query=Id__in:" + strings.TrimPrefix(payload, filtroID)
	}

	return payload
}
