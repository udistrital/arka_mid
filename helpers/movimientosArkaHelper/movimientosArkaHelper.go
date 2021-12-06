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

func GetAllEstadoMovimiento(nombre string) (estado []*models.EstadoMovimiento, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetAllEstadoMovimiento", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

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

func GetAllElementosMovimiento(query string) (elementos []*models.ElementosMovimiento, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetAllElementosMovimiento", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento?" + query
	if err := request.GetJson(urlcrud, &elementos); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllElementosMovimiento - request.GetJson(urlcrud, &elementos)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	return elementos, nil
}

func GetAllSoporteMovimiento(query string) (soportes []*models.SoporteMovimiento, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetAllSoporteMovimiento", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "soporte_movimiento?" + query
	if err := request.GetJson(urlcrud, &soportes); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllSoporteMovimiento - request.GetJson(urlcrud, &soportes)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	return soportes, nil
}

func GetAllMovimiento(query string) (movimientos []*models.Movimiento, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetAllMovimiento", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento?" + query
	if err := request.GetJson(urlcrud, &movimientos); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllMovimiento - request.GetJson(urlcrud, &movimientos)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	return movimientos, nil
}

func GetMovimientoById(id int) (movimiento *models.Movimiento, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetMovimientoById", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	// Se consulta el movimiento
	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &movimiento); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetMovimientoById - request.GetJson(urlcrud, &movimiento)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	return movimiento, nil
}

func PostMovimiento(movimiento *models.Movimiento) (movimientoR *models.Movimiento, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/PostMovimiento", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	var (
		res *models.Movimiento
	)

	// Crea registro en api movimientos_arka_crud
	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento"
	if err := request.SendJson(urlcrud, "POST", &res, &movimiento); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "PostMovimiento - request.SendJson(urlcrud, \"POST\", &res, &movimiento)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	return res, nil
}

func PostSoporteMovimiento(soporte *models.SoporteMovimiento) (soporteR *models.SoporteMovimiento, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/PostSoporteMovimiento", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "soporte_movimiento"
	if err := request.SendJson(urlcrud, "POST", &soporteR, &soporte); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "PostSoporteMovimiento - request.SendJson(urlcrud, \"POST\", &soporteR, &soporte)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	return soporteR, nil
}

func PutTrSalida(trSalida *models.SalidaGeneral) (trResultado *models.SalidaGeneral, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/PutTrSalida", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "tr_salida/"
	if err := request.SendJson(urlcrud, "PUT", &trResultado, &trSalida); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "PutTrSalida - request.SendJson(movArka, \"PUT\", &trResultado, &trSalida)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	return trResultado, nil

}

func PutMovimiento(movimiento *models.Movimiento, movimientoId int) (movimientoRes *models.Movimiento, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/PutMovimiento", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento/" + strconv.Itoa(movimientoId)
	if err := request.SendJson(urlcrud, "PUT", &movimientoRes, &movimiento); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "PutMovimiento - request.SendJson(urlcrud, \"PUT\", &res, &m.Salidas[0].Salida)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	return movimientoRes, nil

}

func PutRevision(revision *models.TrRevisionBaja) (ids []int, outputError map[string]interface{}) {

	funcion := "PutRevision"
	defer errorctrl.ErrorControlFunction(funcion, "500")

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "bajas/"
	if err := request.SendJson(urlcrud, "PUT", &ids, &revision); err != nil {
		logs.Error(err)
		funcion += " - request.SendJson(urlcrud, \"PUT\", &ids, &revision)"
		return nil, errorctrl.Error(funcion, err, "500")
	}

	return ids, nil

}

func PutSoporteMovimiento(soporte *models.SoporteMovimiento, soporteId int) (soporteR *models.SoporteMovimiento, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/PutSoporteMovimiento", "err": err, "status": "500"}
			panic(outputError)
		}
	}()

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "soporte_movimiento/" + strconv.Itoa(soporteId)
	if err := request.SendJson(urlcrud, "PUT", &soporteR, &soporte); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "PutSoporteMovimiento - request.SendJson(urlcrud, \"PUT\", &soporteR, &soporte)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	return soporteR, nil

}
