package trasladoshelper

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/movimientosArkaHelper"
	"github.com/udistrital/arka_mid/helpers/tercerosMidHelper"
	"github.com/udistrital/arka_mid/helpers/ubicacionHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetDetalle Consulta los funcionarios, ubicación y elementos asociados a un traslado
func GetDetalleTraslado(id int) (Traslado *models.TrTraslado, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetDetalleTraslado", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var (
		movimiento *models.Movimiento
		detalle    models.DetalleTraslado
	)
	Traslado = new(models.TrTraslado)

	// Se consulta el movimiento
	if movimientoA, err := movimientosArkaHelper.GetMovimientoById(id); err != nil {
		return nil, err
	} else {
		movimiento = movimientoA
	}

	if err := json.Unmarshal([]byte(movimiento.Detalle), &detalle); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetDetalleTraslado - json.Unmarshal([]byte(movimiento.Detalle), &detalle)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	// Se consulta el detalle del funcionario origen
	if origen, err := tercerosMidHelper.GetDetalleFuncionario(detalle.FuncionarioOrigen); err != nil {
		return nil, err
	} else {
		Traslado.FuncionarioOrigen = origen
	}

	// Se consulta el detalle del funcionario destino
	if destino, err := tercerosMidHelper.GetDetalleFuncionario(detalle.FuncionarioDestino); err != nil {
		return nil, err
	} else {
		Traslado.FuncionarioDestino = destino
	}

	// Se consulta la sede, dependencia correspondiente a la ubicacion
	if ubicacionDetalle, err := ubicacionHelper.GetSedeDependenciaUbicacion(detalle.Ubicacion); err != nil {
		return nil, err
	} else {
		Traslado.Ubicacion = ubicacionDetalle
	}

	// Se consultan los detalles de los elementos del traslado
	if elementos, err := GetElementosTraslado(detalle.Elementos); err != nil {
		return nil, err
	} else {
		Traslado.Elementos = elementos
	}
	Traslado.Detalle = movimiento.Detalle
	Traslado.Observaciones = movimiento.Observacion

	return Traslado, nil

}

func GetElementosTraslado(ids []int) (Elementos []*models.DetalleElementoPlaca, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetDetalleFuncionario", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var (
		urlcrud   string
		elementos []*models.DetalleElementoPlaca
	)

	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento?limit=-1&fields=Id,ElementoActaId&sortby=ElementoActaId&order=desc"
	urlcrud += "&query=Id__in:" + url.QueryEscape(utilsHelper.ArrayToString(ids, ";"))
	if err := request.GetJson(urlcrud, &elementos); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetElementosTraslado - request.GetJson(urlcrud, &elementos)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	idsActa := []int{}
	for _, val := range elementos {
		idsActa = append(idsActa, int(val.ElementoActaId))
	}

	if response, err := actaRecibido.GetElementosByIds(idsActa); err != nil {
		return nil, err
	} else {
		for _, elemento := range elementos {
			if i := utilsHelper.FindIdInArray(response, int(elemento.ElementoActaId)); i > -1 {
				if len(response) > 1 {
					response = append(response[:i], response[i+1:]...)
				}
				elemento.Placa = response[i].Placa
				elemento.Nombre = response[i].Nombre
				elemento.Marca = response[i].Marca
			}
		}
	}

	Elementos = elementos
	return
}

// RegistrarEntrada Crea registro de traslado en estado en trámite
func RegistrarTraslado(data *models.Movimiento) (result *models.Movimiento, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "RegistrarTraslado - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	result = new(models.Movimiento)

	detalleJSON := map[string]interface{}{}
	if err := json.Unmarshal([]byte(data.Detalle), &detalleJSON); err != nil {
		panic(err.Error())
	}

	ctxConsecutivo, _ := beego.AppConfig.Int("contxtTrasladoCons")
	if consecutivo, err := utilsHelper.GetConsecutivo("%05.0f", ctxConsecutivo, "Registro Traslado Arka"); err != nil {
		return nil, err
	} else {
		consecutivo = utilsHelper.FormatConsecutivo(getTipoComprobanteTraslados()+"-", consecutivo, fmt.Sprintf("%s%04d", "-", time.Now().Year()))
		detalleJSON["Consecutivo"] = consecutivo
	}

	if jsonData, err := json.Marshal(detalleJSON); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "RegistrarTraslado - json.Marshal(detalleJSON)",
			"err":     err,
			"status":  "500",
		}
		return nil, outputError
	} else {
		data.Detalle = string(jsonData[:])
	}

	// Crea registro en api movimientos_arka_crud
	if res, err := movimientosArkaHelper.PostMovimiento(data); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func getTipoComprobanteTraslados() string {
	return "T"
}
