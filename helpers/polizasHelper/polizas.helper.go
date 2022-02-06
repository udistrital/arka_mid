package polizashelper

import (
	"encoding/json"

	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

// GetElementosPoliza Obtiene todos los elementos que necesiten poliza
//
func GetElementosPoliza(limit int, offset int, fields []string, order []string, query map[string]string) (ElementosPoliza []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetAllPolizasEntradas - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	type M map[string]interface{}

	var (
		ActasIdsEntradas  []int
		ElementosEntradas []interface{}
		SubgruposPoliza   []int
		//ElementosPorId    interface{}
	)

	// Trae las ActasId que son entradas
	if Elementos, err := GetIdsActasEntrada(); err != nil {
		return nil, err
	} else {
		ActasIdsEntradas = Elementos
	}
	logs.Debug(ActasIdsEntradas)

	if IdSubgrupo, err := SubgrupoParaPoliza(); err != nil {
		return nil, err
	} else {
		SubgruposPoliza = IdSubgrupo
	}
	logs.Debug(SubgruposPoliza)

	for i := 0; i < len(ActasIdsEntradas); i++ {

		// Con las ActasId trae los elementos de esa acta
		if Elementos, err := GetElementosPorActaIds(ActasIdsEntradas[i]); err != nil {
			return nil, err
		} else {
			ElementosEntradas = Elementos
		}

		m1 := M{"Id Acta": ActasIdsEntradas[i], "Elementos": ElementosEntradas}
		ElementosPoliza = append(ElementosPoliza, m1)
	}

	// for i := 0; i < len(ActasIdsEntradas); i++ {
	// 	//m1 := M{"IdActa": ElementosEntrada[i], "FechaEntrada": i, "NPoliza": i, "Descripcion": i}
	// 	m1 := M{"IdActa": ActasIdsEntradas[i]}
	// 	ElementosPoliza = append(ElementosPoliza, m1)
	// }
	return ElementosPoliza, nil
}

// GetElementosEntrada obtiene todos las ActasId que sean entradas (P8)
//
func GetIdsActasEntrada() (IdsActasEntradas []int, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetActaIdEntradas - Unhandled Error!",
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
	urlRespuestaAPI := "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento?query=Activo%3Atrue%2CDetalle__contains%3Aconsecutivo%5C%22%3A%20%5C%22P&fields=Detalle" //+ "&limit=-1"
	if _, err := request.GetJsonTest(urlRespuestaAPI, &RespuestaAPI); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetIdsActasEntradas - request.GetJsonTest(urlRespuestaAPI, &RespuestaAPI)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	// Asignar a Movimientos RespuestaAPI
	if err := formatdata.FillStruct(&RespuestaAPI, &Movimientos); err != nil {
		logs.Error(err)
	}

	for _, movimiento := range Movimientos {
		var detalle map[string]interface{}
		if err := json.Unmarshal([]byte(movimiento.Detalle), &detalle); err != nil {
			logs.Error(err)
		}
		IdsActasEntradas = append(IdsActasEntradas, int(detalle["acta_recibido_id"].(float64)))
	}

	ActasIdsEntradas = removeDuplicateElement(IdsActasEntradas)

	return ActasIdsEntradas, nil
}

// Consulta los elementos por ActaId y retorna algunos parametros
//
func GetElementosPorActaIds(ActaId int) (ElementosEntradas []interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetIdElementoPlaca - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	ActaStr := strconv.Itoa(ActaId)
	urlElementosEntradas := "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento?query=Activo%3ATrue%2CActaRecibidoId:" + ActaStr + "&fields=Id%2CNombre%2CSubgrupoCatalogoId%2CCantidad"
	if _, err := request.GetJsonTest(urlElementosEntradas, &ElementosEntradas); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetElementosPoliza - request.GetJsonTest(urlElementosEntradas, &ElementosEntradas)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	return
}

// Obtiene el Id de los subgrupos que necesitan poliza
//
func SubgrupoParaPoliza() (IdSubgruposPoliza []int, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetIdElementoPlaca - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var IdTipoBien interface{}

	urlIdTipoBien := "http://" + beego.AppConfig.String("catalogoElementosService") + "tipo_bien?query=Activo%3ATrue%2CNecesitaPoliza%3ATrue&fields=Id"
	if _, err := request.GetJsonTest(urlIdTipoBien, &IdTipoBien); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "SubgrupoParaPoliza - request.GetJsonTest(urlIdTipoBien, &IdTipoBien)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	logs.Debug(IdTipoBien)
	return
}

// Compara los ids de las actas y remueve las actas duplicadas
//
func removeDuplicateElement(addrs []int) []int {
	result := make([]int, 0, len(addrs))
	temp := map[int]struct{}{}
	for _, item := range addrs {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
