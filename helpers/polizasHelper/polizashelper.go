package polizashelper

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	utilsHelper "github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/request"
)

// GetElementosPoliza Obtiene todos los elementos que necesiten poliza
func GetElementosPoliza(offset int, limit int, fields []string, order []string, query map[string]string, sortby []string) (ElementosPoliza *[]models.Elemento, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetElementosPoliza - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var (
		ActasIdsEntradas []int
		SubgrupoPoliza   []int
	)

	// Trae las ActasId que son entradas en movimientoservice.movimiento
	if Elementos, err := GetIdsActasEntrada(); err != nil {
		return nil, err
	} else {
		ActasIdsEntradas = Elementos
	}

	// Obtiene IdDetalleSubgrupo de los tipos de bien en catalogoelemento.detallesubgrupo
	if idsubgrupo, err := GetSubgruposPoliza(); err != nil {
		return nil, err
	} else {
		SubgrupoPoliza = idsubgrupo
	}

	// Respuesta de los Elementos que necesiten poliza y pudiendo filtrar
	if RElementosPoliza, err := GetElementosPolizas(ActasIdsEntradas, SubgrupoPoliza, limit, offset, fields, order, query, sortby); err != nil {
		return nil, err
	} else {
		ElementosPoliza = RElementosPoliza
	}

	return
}

// GetElementosEntrada obtiene todos las ActasId que sean entradas (P8)
func GetIdsActasEntrada() (IdsActasEntradas []int, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetIdsActasEntrada - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var (
		RespuestaAPI     interface{}
		Movimientos      []models.Movimiento
		ActasIdsEntradas []int
	)

	//Traer el detalle de las entradas (json)
	params := url.Values{}
	params.Add("query", "Activo:true,Detalle__contains:consecutivo\": \"P")
	params.Add("fields", "Detalle")
	params.Add("limit", "-1")
	path, _ := beego.AppConfig.String("movimientosArkaService")
	urlRespuestaAPI := "http://" + path + "movimiento?" + params.Encode()
	if _, err := request.GetJsonTest(urlRespuestaAPI, &RespuestaAPI); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetIdsActasEntrada - request.GetJsonTest(urlRespuestaAPI, &RespuestaAPI)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	// Asignar a Movimientos RespuestaAPI
	if err := utilsHelper.FillStruct(&RespuestaAPI, &Movimientos); err != nil {
		logs.Error(err)
	}

	//Itera los Movimientos para poder guardar solo las actas de recibido
	for _, movimiento := range Movimientos {
		var detalle map[string]interface{}
		if err := json.Unmarshal([]byte(movimiento.Detalle), &detalle); err != nil {
			logs.Warn(err)
		} else {
			ActasIdsEntradas = append(ActasIdsEntradas, int(detalle["acta_recibido_id"].(float64)))
		}
	}

	//Remueve los Ids replicados
	IdsActasEntradas = utilsHelper.RemoveDuplicateInt(ActasIdsEntradas)

	return
}

// Obtiene todos los idSubgrupos que necesitan poliza
func GetSubgruposPoliza() (IdSubgruposPoliza []int, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetSubgruposPoliza - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var (
		Respuesta         interface{}
		DetalleSubgrupo   []models.DetalleSubgrupo
		IdDetalleSubgrupo []int
	)

	//Traer el detalle de las entradas (json)
	params := url.Values{}
	params.Add("query", "TipoBienId__NecesitaPoliza:True,Activo:True")
	params.Add("limit", "-1")
	path, _ := beego.AppConfig.String("catalogoElementosService")
	urlSubgrupos := "http://" + path + "detalle_subgrupo?" + params.Encode()
	if _, err := request.GetJsonTest(urlSubgrupos, &Respuesta); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetSubgruposPoliza - request.GetJsonTest(urlSubgrupos, &Respuesta)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	// Asignar a DetalleSubgrupo RespuestaAPI
	if err := utilsHelper.FillStruct(&Respuesta, &DetalleSubgrupo); err != nil {
		logs.Warn(err)
	}

	// Itera los DetalleSubgrupo para guardar solo los SubgrupoId
	for _, detallesub := range DetalleSubgrupo {
		IdDetalleSubgrupo = append(IdDetalleSubgrupo, detallesub.SubgrupoId.Id)
	}

	// Remueve los subgrupos replicados
	IdSubgruposPoliza = utilsHelper.RemoveDuplicateInt(IdDetalleSubgrupo)

	return
}

// Consulta los elementos por ActaId y retorna algunos parametros
func GetElementosPolizas(ActasIdsEntradas []int, SubgrupoPoliza []int, limit int, offset int, fields []string, order []string,
	query map[string]string, sortby []string) (ElementosEntradas *[]models.Elemento, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetElementosPolizas - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var (
		Respuest interface{}
		Query    []string
		QUERY    string
	)

	// Se convierten los filtros en terminos de string para poder hacer la consulta por url
	lim := strconv.Itoa(limit)
	off := strconv.Itoa(offset)
	//field := strings.Join(fields, ",")
	sort := strings.Join(sortby, ",")
	orde := strings.Join(order, ",")

	params := url.Values{}

	for key, value := range query {
		Query = append(Query, key+":"+value)
		QUER := strings.Trim(strings.Replace(fmt.Sprint(Query), " ", ",", -1), "[]")
		QUERY = QUER
	}

	if QUERY != "" {
		params.Add("query", "Activo:True,"+QUERY+",ActaRecibidoId__Id__in:"+utilsHelper.ArrayToString(ActasIdsEntradas, "|")+
			",SubgrupoCatalogoId__in:"+utilsHelper.ArrayToString(SubgrupoPoliza, "|"))
	} else {
		params.Add("query", "Activo:True,ActaRecibidoId__Id__in:"+utilsHelper.ArrayToString(ActasIdsEntradas, "|")+
			",SubgrupoCatalogoId__in:"+utilsHelper.ArrayToString(SubgrupoPoliza, "|"))
	}
	if lim != "10" {
		params.Add("limit", lim)
	}
	if off != "0" {
		params.Add("offset", off)
	}
	// if field != "" {
	// 	params.Add("fields", field)
	// }
	if sort != "" {
		params.Add("sortby", sort)
	}
	if orde != "" {
		params.Add("order", orde)
	}

	path, _ := beego.AppConfig.String("actaRecibidoService")
	urlElementosSubgrupos := "http://" + path + "elemento?" + params.Encode()
	//logs.Debug(urlElementosSubgrupos)
	if _, err := request.GetJsonTest(urlElementosSubgrupos, &Respuest); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetElementosPolizas - request.GetJsonTest(urlElementosSubgrupos, &Respuest)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	if err := utilsHelper.FillStruct(&Respuest, &ElementosEntradas); err != nil {
		logs.Error(err)
	}

	return
}
