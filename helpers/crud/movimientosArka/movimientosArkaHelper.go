package movimientosArka

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	"github.com/udistrital/arka_mid/utils_oas/request"
)

var basePath, _ = beego.AppConfig.String("movimientosArkaService")

// GetAllEstadoMovimiento query controlador estado_movimiento del api movimientos_arka_crud
func GetAllEstadoMovimiento(query string) (estados []*models.EstadoMovimiento, outputError map[string]interface{}) {

	funcion := "GetAllEstadoMovimiento"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error", "500")

	urlcrud := "http://" + basePath + "estado_movimiento?" + query
	if err := request.GetJson(urlcrud, &estados); err != nil {
		eval := " - request.GetJson(urlcrud, &estados)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetAllFormatoTipoMovimiento query controlador formato_tipo_movimiento del api movimientos_arka_crud
func GetAllFormatoTipoMovimiento(query string) (formatos []*models.FormatoTipoMovimiento, outputError map[string]interface{}) {

	funcion := "GetAllFormatoTipoMovimiento"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "formato_tipo_movimiento?" + query
	if err := request.GetJson(urlcrud, &formatos); err != nil {
		eval := " - request.GetJson(urlcrud, &formatos)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}
	return formatos, nil
}

// GetAllElementosMovimiento query controlador elementos_movimiento del api movimientos_arka_crud
func GetAllElementosMovimiento(query string) (elementos []*models.ElementosMovimiento, outputError map[string]interface{}) {

	funcion := "GetAllElementosMovimiento"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "elementos_movimiento?" + query
	if err := request.GetJson(urlcrud, &elementos); err != nil {
		eval := " - request.GetJson(urlcrud, &elementos)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}
	return elementos, nil
}

// GetAllSoporteMovimiento query controlador soporte_movimiento del api movimientos_arka_crud
func GetAllSoporteMovimiento(query string) (soportes []models.SoporteMovimiento, outputError map[string]interface{}) {

	funcion := "GetAllSoporteMovimiento"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "soporte_movimiento?" + query
	if err := request.GetJson(urlcrud, &soportes); err != nil {
		eval := " - request.GetJson(urlcrud, &soportes)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}
	return soportes, nil
}

// GetAllMovimiento query controlador movimiento del api movimientos_arka_crud
func GetAllMovimiento(query string) (movimientos []*models.Movimiento, outputError map[string]interface{}) {

	funcion := "GetAllMovimiento"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "movimiento?" + query
	if err := request.GetJson(urlcrud, &movimientos); err != nil {
		eval := " - request.GetJson(urlcrud, &movimientos)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}
	return movimientos, nil
}

// GetAllNovedadElemento query controlador novedad_elemento del api movimientos_arka_crud
func GetAllNovedadElemento(query string) (novedades []*models.NovedadElemento, outputError map[string]interface{}) {

	funcion := "GetAllNovedadElemento"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "novedad_elemento?" + query
	if err := request.GetJson(urlcrud, &novedades); err != nil {
		eval := " - request.GetJson(urlcrud, &novedades)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}
	return
}

// GetMovimientoById consulta controlador movimiento/{id} del api movimientos_arka_crud
func GetMovimientoById(id int) (movimiento *models.Movimiento, outputError map[string]interface{}) {

	funcion := "GetMovimientoById"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error", "500")

	// Se consulta el movimiento
	urlcrud := "http://" + basePath + "movimiento/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &movimiento); err != nil {
		eval := " - request.GetJson(urlcrud, &movimiento)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return movimiento, nil
}

// GetElementosMovimientoById consulta controlador elementos_movimiento/{id} del api movimientos_arka_crud
func GetElementosMovimientoById(id int, elemento *models.ElementosMovimiento) (outputError map[string]interface{}) {

	funcion := "GetElementosMovimientoById - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := "http://" + basePath + "elementos_movimiento/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &elemento); err != nil {
		logs.Error(err, urlcrud)
		eval := "request.GetJson(urlcrud, &elemento)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetTrSalida consulta controlador tr_salida/{id} del api movimientos_arka_crud
func GetTrSalida(id int) (trSalida *models.TrSalida, outputError map[string]interface{}) {

	funcion := "GetTrSalida"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error", "500")

	// Se consulta el movimiento
	urlcrud := "http://" + basePath + "tr_salida/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &trSalida); err != nil {
		eval := " - request.GetJson(urlcrud, &trSalida)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return trSalida, nil
}

// PostMovimiento post controlador movimiento del api movimientos_arka_crud
func PostMovimiento(movimiento *models.Movimiento) (outputError map[string]interface{}) {

	funcion := "PostMovimiento - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := "http://" + basePath + "movimiento"
	if err := request.SendJson(urlcrud, "POST", &movimiento, &movimiento); err != nil {
		logs.Error(err)
		eval := `request.SendJson(urlcrud, "POST", &movimiento, &movimiento)`
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// PostSoporteMovimiento post controlador soporte_movimiento del api movimientos_arka_crud
func PostSoporteMovimiento(soporte *models.SoporteMovimiento) (outputError map[string]interface{}) {

	funcion := "PostSoporteMovimiento - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "soporte_movimiento"
	if err := request.SendJson(urlcrud, "POST", &soporte, &soporte); err != nil {
		logs.Error(err)
		eval := `request.SendJson(urlcrud, "POST", &soporte, &soporte)`
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// PostElementosMovimiento post controlador elementos_movimiento del api movimientos_arka_crud
func PostElementosMovimiento(elemento *models.ElementosMovimiento) (outputError map[string]interface{}) {

	funcion := "PostElementosMovimiento - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := "http://" + basePath + "elementos_movimiento"
	err := request.SendJson(urlcrud, "POST", &elemento, &elemento)
	if err != nil {
		logs.Error(err, urlcrud)
		eval := `request.SendJson(urlcrud, "POST", &elemento, &elemento)`
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// PutTrSalida put controlador tr_salida del api movimientos_arka_crud
func PutTrSalida(trSalida *models.SalidaGeneral) (trResultado *models.SalidaGeneral, outputError map[string]interface{}) {

	funcion := "PutTrSalida"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "tr_salida/"
	if err := request.SendJson(urlcrud, "PUT", &trResultado, &trSalida); err != nil {
		eval := " - request.SendJson(urlcrud, \"PUT\", &trResultado, &trSalida)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return trResultado, nil

}

// PostTrSalida post controlador tr_salida del api movimientos_arka_crud
func PostTrSalida(trSalida *models.SalidaGeneral) (outputError map[string]interface{}) {

	funcion := "PostTrSalida - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "tr_salida"
	err := request.SendJson(urlcrud, "POST", &trSalida, &trSalida)
	if err != nil {
		logs.Error(err)
		eval := `request.SendJson(urlcrud, "POST", &trSalida, &trSalida)`
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// PutMovimiento put controlador movimiento del api movimientos_arka_crud
func PutMovimiento(movimiento *models.Movimiento, movimientoId int) (outputError map[string]interface{}) {

	funcion := "PutMovimiento"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "movimiento/" + strconv.Itoa(movimientoId)
	err := request.SendJson(urlcrud, "PUT", &movimiento, &movimiento)
	if err != nil {
		eval := `request.SendJson(urlcrud, "PUT", &movimientoRes, &movimiento)`
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// PutRevision put controlador bajas/ del api movimientos_arka_crud
func PutRevision(revision *models.TrRevisionBaja) (ids []int, outputError map[string]interface{}) {

	funcion := "PutRevision"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "bajas/"
	if err := request.SendJson(urlcrud, "PUT", &ids, &revision); err != nil {
		eval := " - request.SendJson(urlcrud, \"PUT\", &ids, &revision)"
		return nil, errorCtrl.Error(funcion+eval, err, "500")
	}

	return ids, nil

}

// PutSoporteMovimiento put controlador soporte_movimiento del api movimientos_arka_crud
func PutSoporteMovimiento(soporte *models.SoporteMovimiento, soporteId int) (soporteR *models.SoporteMovimiento, outputError map[string]interface{}) {

	funcion := "PutSoporteMovimiento"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "soporte_movimiento/" + strconv.Itoa(soporteId)
	if err := request.SendJson(urlcrud, "PUT", &soporteR, &soporte); err != nil {
		eval := " - request.SendJson(urlcrud, \"PUT\", &soporteR, &soporte)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return soporteR, nil

}

// PutElementosMovimiento put controlador elementos_movimiento del api movimientos_arka_crud
func PutElementosMovimiento(elementoM *models.ElementosMovimiento, elementoId int) (elementoM_ *models.ElementosMovimiento, outputError map[string]interface{}) {

	funcion := "PutElementosMovimiento"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "elementos_movimiento/" + strconv.Itoa(elementoId)
	if err := request.SendJson(urlcrud, "PUT", &elementoM_, &elementoM); err != nil {
		eval := ` - request.SendJson(urlcrud, "PUT", &soporteR, &soporte)`
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return elementoM_, nil

}

// PutNovedadElemento put controlador novedad_elemento del api movimientos_arka_crud
func PutNovedadElemento(novedad *models.NovedadElemento, novedadId int) (novedad_ *models.NovedadElemento, outputError map[string]interface{}) {

	funcion := "PutNovedadElemento"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "novedad_elemento/" + strconv.Itoa(novedadId)
	if err := request.SendJson(urlcrud, "PUT", &novedad_, &novedad); err != nil {
		eval := ` - request.SendJson(urlcrud, "PUT", &novedad_, &novedad)`
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return novedad_, nil

}

// PostNovedadElemento post controlador novedad_elemento del api movimientos_arka_crud
func PostNovedadElemento(novedad *models.NovedadElemento) (outputError map[string]interface{}) {

	funcion := "PostNovedadElemento - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := "http://" + basePath + "novedad_elemento"
	err := request.SendJson(urlcrud, "POST", &novedad, &novedad)
	if err != nil {
		logs.Error(err, urlcrud)
		eval := `request.SendJson(urlcrud, "POST", &novedad, &novedad)`
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetElementosFuncionario query controlador elementos_movimiento/funcionario/{funcionarioId} del api movimientos_arka_crud
func GetElementosFuncionario(funcionarioId int) (movimientos []int, outputError map[string]interface{}) {

	funcion := "GetElementosFuncionario"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "elementos_movimiento/funcionario/" + strconv.Itoa(funcionarioId)
	if err := request.GetJson(urlcrud, &movimientos); err != nil {
		eval := " - request.GetJson(urlcrud, &movimientos)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}
	return movimientos, nil
}

// GetHistorialElemento query controlador elementos_movimiento/historial/{elementoId} del api movimientos_arka_crud
func GetHistorialElemento(elementoId int, final bool) (historial *models.Historial, outputError map[string]interface{}) {

	funcion := "GetHistorialElemento"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "elementos_movimiento/historial/" + strconv.Itoa(elementoId)
	urlcrud += "?final=" + strconv.FormatBool(final)
	if err := request.GetJson(urlcrud, &historial); err != nil {
		eval := " - request.GetJson(urlcrud, &historial)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}
	return historial, nil
}

// GetCorteDepreciacion query controlador cierre/?fechaCorte={fechaCorte} del api movimientos_arka_crud
func GetCorteDepreciacion(fechaCorte string) (corte []models.DepreciacionElemento, outputError map[string]interface{}) {

	funcion := "GetCorteDepreciacion - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "cierre/?fechaCorte=" + fechaCorte
	if err := request.GetJson(urlcrud, &corte); err != nil {
		logs.Error(err, urlcrud)
		eval := "request.GetJson(urlcrud, &corte)"
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// AprobarCierre post controlador cierre del api movimientos_arka_crud
func AprobarCierre(cierre *models.Movimiento) (outputError map[string]interface{}) {

	funcion := "AprobarCierre - "
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error", "500")

	urlcrud := "http://" + basePath + "cierre/"
	if err := request.SendJson(urlcrud, "POST", &cierre, &cierre); err != nil {
		logs.Error(err, urlcrud)
		eval := `request.SendJson(urlcrud, "POST", &cierre, &data)`
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetEntradaByActa consulta controlador movimiento/entrada/{acta_recibido_id} del api movimientos_arka_crud
func GetEntradaByActa(acta_recibido_id int) (entrada *models.Movimiento, outputError map[string]interface{}) {

	funcion := "GetEntradaByActa"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "movimiento/entrada/" + strconv.Itoa(acta_recibido_id)
	if err := request.GetJson(urlcrud, &entrada); err != nil {
		eval := " - request.GetJson(urlcrud, &entrada)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}
	return entrada, nil
}

// GetTrasladosByTerceroId consulta controlador movimiento/traslado/{tercero_id} del api movimientos_arka_crud
func GetTrasladosByTerceroId(terceroId int, confirmar bool, traslados *[]*models.Movimiento) (outputError map[string]interface{}) {

	funcion := "GetTrasladosByTerceroId - "
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "movimiento/traslado/" + strconv.Itoa(terceroId)
	if confirmar {
		urlcrud += "?confirmar=true"
	}
	if err := request.GetJson(urlcrud, &traslados); err != nil {
		eval := "request.GetJson(urlcrud, &traslados)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}
	return
}

// GetBajasByTerceroId consulta controlador movimiento/baja/{tercero_id} del api movimientos_arka_crud
func GetBajasByTerceroId(terceroId int, bajas *[]*models.Movimiento) (outputError map[string]interface{}) {

	funcion := "GetBajasByTerceroId - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "movimiento/baja/" + strconv.Itoa(terceroId)
	if err := request.GetJson(urlcrud, &bajas); err != nil {
		eval := "request.GetJson(urlcrud, &bajas)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetBodegaByTerceroId consulta controlador movimiento/bodega/{tercero_id} del api movimientos_arka_crud
func GetBodegaByTerceroId(terceroId int, solicitudes *[]*models.Movimiento) (outputError map[string]interface{}) {

	funcion := "GetBodegaByTerceroId - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "movimiento/bodega/" + strconv.Itoa(terceroId)
	if err := request.GetJson(urlcrud, &solicitudes); err != nil {
		eval := "request.GetJson(urlcrud, &solicitudes)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

// GetAperturas consulta controlador tr_kardex/aperturas del api movimientos_arka_crud
func GetAperturas(conSaldo bool, aperturas *[]models.Apertura) (outputError map[string]interface{}) {

	funcion := "GetAperturas - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "tr_kardex/aperturas?ConSaldo=" + strconv.FormatBool(conSaldo)
	if err := request.GetJson(urlcrud, &aperturas); err != nil {
		eval := "request.GetJson(urlcrud, &aperturas)"
		return errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}

func GetAllCentroCostos(payload string) (centroCostos []models.CentroCostos, outputError map[string]interface{}) {

	funcion := "GetAllCentroCostos - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error", "500")

	urlcrud := "http://" + basePath + "centro_costos?" + payload
	err := request.GetJson(urlcrud, &centroCostos)
	if err != nil {
		logs.Error(err)
		eval := "request.GetJson(urlcrud, &centroCostos)"
		outputError = errorCtrl.Error(funcion+eval, err, "502")
	}

	return
}
