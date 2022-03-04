package movimientosArkaHelper

import (
	"errors"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

// GetAllEstadoMovimiento query controlador estado_movimiento del api movimientos_arka_crud
func GetAllEstadoMovimiento(nombre string) (estado []*models.EstadoMovimiento, outputError map[string]interface{}) {

	funcion := "GetAllEstadoMovimiento"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error", "500")

	var (
		resEstadoMovimiento []*models.EstadoMovimiento
		urlcrud             string
	)

	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "estado_movimiento?query=Nombre:" + nombre
	if err := request.GetJson(urlcrud, &resEstadoMovimiento); err != nil || len(resEstadoMovimiento) == 0 {
		status := "502"
		if err == nil {
			err = errors.New("len(resEstadoMovimiento) == 0")
			status = "404"
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllEstadoMovimiento - request.GetJson(urlcrud, &resEstadoMovimiento)",
			"err":     err,
			"status":  status,
		}
		return nil, outputError
	}
	return resEstadoMovimiento, nil
}

// GetAllFormatoTipoMovimiento query controlador formato_tipo_movimiento del api movimientos_arka_crud
func GetAllFormatoTipoMovimiento(query string) (formatos []*models.FormatoTipoMovimiento, outputError map[string]interface{}) {

	funcion := "GetAllFormatoTipoMovimiento"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "formato_tipo_movimiento?" + query
	if err := request.GetJson(urlcrud, &formatos); err != nil {
		eval := " - request.GetJson(urlcrud, &formatos)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}
	return formatos, nil
}

// GetAllElementosMovimiento query controlador elementos_movimiento del api movimientos_arka_crud
func GetAllElementosMovimiento(query string) (elementos []*models.ElementosMovimiento, outputError map[string]interface{}) {

	funcion := "GetAllElementosMovimiento"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento?" + query
	if err := request.GetJson(urlcrud, &elementos); err != nil {
		eval := " - request.GetJson(urlcrud, &elementos)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}
	return elementos, nil
}

// GetAllSoporteMovimiento query controlador soporte_movimiento del api movimientos_arka_crud
func GetAllSoporteMovimiento(query string) (soportes []*models.SoporteMovimiento, outputError map[string]interface{}) {

	funcion := "GetAllSoporteMovimiento"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "soporte_movimiento?" + query
	if err := request.GetJson(urlcrud, &soportes); err != nil {
		eval := " - request.GetJson(urlcrud, &soportes)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}
	return soportes, nil
}

// GetAllMovimiento query controlador movimiento del api movimientos_arka_crud
func GetAllMovimiento(query string) (movimientos []*models.Movimiento, outputError map[string]interface{}) {

	funcion := "GetAllMovimiento"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento?" + query
	if err := request.GetJson(urlcrud, &movimientos); err != nil {
		eval := " - request.GetJson(urlcrud, &movimientos)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}
	return movimientos, nil
}

// GetMovimientoById consulta controlador movimiento/{id} del api movimientos_arka_crud
func GetMovimientoById(id int) (movimiento *models.Movimiento, outputError map[string]interface{}) {

	funcion := "GetMovimientoById"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error", "500")

	// Se consulta el movimiento
	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &movimiento); err != nil {
		eval := " - request.GetJson(urlcrud, &movimiento)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}

	return movimiento, nil
}

// GetTrSalida consulta controlador tr_salida/{id} del api movimientos_arka_crud
func GetTrSalida(id int) (trSalida *models.TrSalida, outputError map[string]interface{}) {

	funcion := "GetTrSalida"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error", "500")

	// Se consulta el movimiento
	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "tr_salida/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &trSalida); err != nil {
		eval := " - request.GetJson(urlcrud, &trSalida)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}

	return trSalida, nil
}

// PostMovimiento post controlador movimiento del api movimientos_arka_crud
func PostMovimiento(movimiento *models.Movimiento) (movimientoR *models.Movimiento, outputError map[string]interface{}) {

	funcion := "PostMovimiento"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error", "500")

	var (
		res *models.Movimiento
	)

	// Crea registro en api movimientos_arka_crud
	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento"
	if err := request.SendJson(urlcrud, "POST", &res, &movimiento); err != nil {
		eval := " - request.SendJson(urlcrud, \"POST\", &res, &movimiento)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}

	return res, nil
}

// PostSoporteMovimiento post controlador soporte_movimiento del api movimientos_arka_crud
func PostSoporteMovimiento(soporte *models.SoporteMovimiento) (soporteR *models.SoporteMovimiento, outputError map[string]interface{}) {

	funcion := "PostSoporteMovimiento"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "soporte_movimiento"
	if err := request.SendJson(urlcrud, "POST", &soporteR, &soporte); err != nil {
		eval := " - request.SendJson(urlcrud, \"POST\", &soporteR, &soporte)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}

	return soporteR, nil
}

// PutTrSalida put controlador tr_salida del api movimientos_arka_crud
func PutTrSalida(trSalida *models.SalidaGeneral) (trResultado *models.SalidaGeneral, outputError map[string]interface{}) {

	funcion := "PutTrSalida"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "tr_salida/"
	if err := request.SendJson(urlcrud, "PUT", &trResultado, &trSalida); err != nil {
		eval := " - request.SendJson(urlcrud, \"PUT\", &trResultado, &trSalida)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}

	return trResultado, nil

}

// PutMovimiento put controlador movimiento del api movimientos_arka_crud
func PutMovimiento(movimiento *models.Movimiento, movimientoId int) (movimientoRes *models.Movimiento, outputError map[string]interface{}) {

	funcion := "PutMovimiento"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento/" + strconv.Itoa(movimientoId)
	if err := request.SendJson(urlcrud, "PUT", &movimientoRes, &movimiento); err != nil {
		eval := " - request.SendJson(urlcrud, \"PUT\", &movimientoRes, &movimiento)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}

	return movimientoRes, nil

}

// PutRevision put controlador bajas/ del api movimientos_arka_crud
func PutRevision(revision *models.TrRevisionBaja) (ids []int, outputError map[string]interface{}) {

	funcion := "PutRevision"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "bajas/"
	if err := request.SendJson(urlcrud, "PUT", &ids, &revision); err != nil {
		eval := " - request.SendJson(urlcrud, \"PUT\", &ids, &revision)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	return ids, nil

}

// PutSoporteMovimiento put controlador soporte_movimiento del api movimientos_arka_crud
func PutSoporteMovimiento(soporte *models.SoporteMovimiento, soporteId int) (soporteR *models.SoporteMovimiento, outputError map[string]interface{}) {

	funcion := "PutSoporteMovimiento"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "soporte_movimiento/" + strconv.Itoa(soporteId)
	if err := request.SendJson(urlcrud, "PUT", &soporteR, &soporte); err != nil {
		eval := " - request.SendJson(urlcrud, \"PUT\", &soporteR, &soporte)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}

	return soporteR, nil

}

// GetElementosFuncionario query controlador elementos_movimiento/funcionario/{funcionarioId} del api movimientos_arka_crud
func GetElementosFuncionario(funcionarioId int) (movimientos []int, outputError map[string]interface{}) {

	funcion := "GetElementosFuncionario"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento/funcionario/" + strconv.Itoa(funcionarioId)
	if err := request.GetJson(urlcrud, &movimientos); err != nil {
		eval := " - request.GetJson(urlcrud, &movimientos)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}
	return movimientos, nil
}

// GetHistorialElemento query controlador elementos_movimiento/historial/{elementoId} del api movimientos_arka_crud
func GetHistorialElemento(elementoId int, final bool) (historial *models.Historial, outputError map[string]interface{}) {

	funcion := "GetHistorialElemento"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento/historial/" + strconv.Itoa(elementoId)
	urlcrud += "?final=" + strconv.FormatBool(final)
	if err := request.GetJson(urlcrud, &historial); err != nil {
		eval := " - request.GetJson(urlcrud, &historial)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}
	return historial, nil
}

// GetCorteDepreciacion query controlador depreciacion/?fechaCorte={fechaCorte} del api movimientos_arka_crud
func GetCorteDepreciacion(fechaCorte string) (corte []*models.DetalleCorteDepreciacion, outputError map[string]interface{}) {

	funcion := "GetCorteDepreciacion"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "depreciacion/"
	urlcrud += "?fechaCorte=" + fechaCorte
	if err := request.GetJson(urlcrud, &corte); err != nil {
		eval := " - request.GetJson(urlcrud, &corte)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}
	return corte, nil
}

// PostTrNovedadElemento post controlador depreciacion del api movimientos_arka_crud
func PostTrNovedadElemento(novedad *models.NovedadElemento) (novedadR *models.NovedadElemento, outputError map[string]interface{}) {

	funcion := "PostMovimiento"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error", "500")

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "depreciacion/"
	if err := request.SendJson(urlcrud, "POST", &novedadR, &novedad); err != nil {
		eval := ` - request.SendJson(urlcrud, "POST", &novedadR, &novedad)`
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}

	return novedadR, nil
}

// GetEntradaByActa consulta controlador movimiento/entrada/{acta_recibido_id} del api movimientos_arka_crud
func GetEntradaByActa(acta_recibido_id int) (entrada []*models.Movimiento, outputError map[string]interface{}) {

	funcion := "GetEntradaByActa"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento/entrada/" + strconv.Itoa(acta_recibido_id)
	if err := request.GetJson(urlcrud, &entrada); err != nil {
		eval := " - request.GetJson(urlcrud, &entrada)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}
	return entrada, nil
}
