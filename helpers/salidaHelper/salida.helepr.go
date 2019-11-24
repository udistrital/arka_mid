package salidaHelper

import (
	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// AddEntrada Transacción para registrar la información de una salida
func AddSalida(data *models.TrSalida) map[string]interface{} {
	var (
		urlcrud   string
		res       map[string]interface{}
		resM      map[string]interface{}
		resultado map[string]interface{}
	)

	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "tr_salida/"

	// Inserta salida en Movimientos ARKA
	if err := request.SendJson(urlcrud, "POST", &res, &data); err == nil {
		// Inserta salida en Movimientos KRONOS
		urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo"

		var salidaId int

		dataSalida := res["Salida"].(map[string]interface{})
		salidaId = int(dataSalida["Id"].(float64))

		procesoExterno := int64(salidaId)
		logs.Debug(procesoExterno)
		tipo := models.TipoMovimiento{Id: 16}
		movimientosKronos := models.MovimientoProcesoExterno{
			TipoMovimientoId: &tipo,
			ProcesoExterno:   procesoExterno,
			Activo:           true,
		}

		if err = request.SendJson(urlcrud, "POST", &resM, &movimientosKronos); err == nil {
			body := res
			body["MovimientosKronos"] = resM["Body"]
			resultado = body
		} else {
			panic(err.Error())
		}
	} else {
		panic(err.Error())
	}

	return resultado
}
