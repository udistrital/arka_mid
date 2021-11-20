package trasladoshelper

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetDetalle Consulta los funcionarios, ubicación y elementos asociados a un traslado
func GetDetalleTraslado(id int) (Traslado map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetDetalleTraslado", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var (
		urlcrud    string
		movimiento models.Movimiento
		detalle    models.DetalleTraslado
	)
	Traslado = make(map[string]interface{})

	// Se consulta el movimiento
	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &movimiento); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetDetalleTraslado - request.GetJson(urlcrud, &movimiento)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
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
	if origen, err := GetDetalleFuncionario(detalle.FuncionarioOrigen); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetDetalleTraslado - GetDetalleFuncionario(detalle.FuncionarioOrigen)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	} else {
		Traslado["FuncionarioOrigen"] = origen
	}

	// Se consulta el detalle del funcionario destino
	if destino, err := GetDetalleFuncionario(detalle.FuncionarioDestino); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetDetalleTraslado - GetDetalleFuncionario(detalle.FuncionarioDestino)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	} else {
		Traslado["FuncionarioDestino"] = destino
	}

	// Se consulta la sede, dependencia correspondiente a la ubicacion
	if ubicacionDetalle, err := utilsHelper.GetSedeDependenciaUbicacion(detalle.Ubicacion); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetDetalleTraslado - GetSedeDependenciaUbicacion(detalle.Ubicacion)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	} else {
		Traslado["Ubicacion"] = ubicacionDetalle
	}

	// Se consultan los detalles de los elementos del traslado
	if elementos, err := GetElementosTraslado(detalle.Elementos); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetDetalleTraslado - GetElementosTraslado(detalle.Elementos)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	} else {
		Traslado["Elementos"] = elementos
	}

	Traslado["Observaciones"] = movimiento.Observacion

	return Traslado, nil

}

// GetDetalle Consulta los funcionarios, ubicación y elementos asociados a un traslado
func GetDetalleFuncionario(id int) (Tercero map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetDetalleFuncionario", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var (
		urlcrud  string
		response []map[string]interface{}
	)

	Tercero = make(map[string]interface{})

	// Consulta información general y documento de identidad
	urlcrud = "http://" + beego.AppConfig.String("tercerosMidService") + "tipo/funcionarios/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &response); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetDetalleFuncionario - request.GetJson(urlcrud, &response1)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	Tercero["Tercero"] = response

	// Consulta correo
	urlcrud = "http://" + beego.AppConfig.String("tercerosService") + "info_complementaria_tercero?limit=1&fields=Dato&sortby=Id&order=desc"
	urlcrud += "&query=Activo%3Atrue,InfoComplementariaId__Nombre__icontains%3Acorreo,TerceroId__Id%3A" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &response); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetDetalleFuncionario - request.GetJson(urlcrud, &response2)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	Tercero["Correo"] = response

	// Consulta cargo
	urlcrud = "http://" + beego.AppConfig.String("tercerosMidService") + "propiedad/cargo/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &response); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetDetalleFuncionario - request.GetJson(urlcrud, &response3)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	Tercero["Cargo"] = response

	return Tercero, nil
}

func GetElementosTraslado(ids []int) (Elementos []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetDetalleFuncionario", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var (
		urlcrud   string
		response  []map[string]interface{}
		elementos []map[string]interface{}
	)

	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento?limit=-1&fields=Id,ElementoActaId&sortby=ElementoActaId&order=desc"
	urlcrud += "&query=Id__in:" + url.QueryEscape(arrayToString(ids, ";"))
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
		idsActa = append(idsActa, int(val["ElementoActaId"].(float64)))
	}

	urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento?limit=-1&fields=Id,Placa,Nombre,Marca&sortby=Id&order=desc"
	urlcrud += "&query=Id__in:" + url.QueryEscape(arrayToString(idsActa, ";"))
	if err := request.GetJson(urlcrud, &response); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetElementosTraslado - request.GetJson(urlcrud, &response)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	for _, elemento := range elementos {
		if i := findIdInArray(response, int(elemento["ElementoActaId"].(float64))); i > -1 {
			if len(response) > 1 {
				response = append(response[:i], response[i+1:]...)
			}
			elemento["Placa"] = response[i]["Placa"]
			elemento["Nombre"] = response[i]["Nombre"]
			elemento["Marca"] = response[i]["Marca"]
		}
	}

	Elementos = elementos
	return
}

func arrayToString(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

// findIdInArray Retorna la posicion en que se encuentra el id específicado
func findIdInArray(idsList []map[string]interface{}, id int) (i int) {
	for i, id_ := range idsList {
		if int(id_["Id"].(float64)) == id {
			return i
		}
	}
	return -1
}
