package movimientosArkaHelper

import (
	"errors"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
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
