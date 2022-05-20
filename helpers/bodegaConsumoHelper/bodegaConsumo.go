package bodegaConsumoHelper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	// "github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

func TraerElementoSolicitud(Elemento map[string]interface{}) (Elemento_ map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "TraerElementoSolicitud - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var idStr int
	if id, err := strconv.Atoi(fmt.Sprintf("%v", Elemento["Ubicacion"])); err == nil {
		idStr = id
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "TraerElementoSolicitud - strconv.Atoi(fmt.Sprintf(\"%v\", Elemento[\"Ubicacion\"]))",
			"err":     err,
			"status":  "400",
		}
		return nil, outputError
	}

	urlcrud3 := "http://" + beego.AppConfig.String("oikosService") + "asignacion_espacio_fisico_dependencia?query=Id:" + strconv.Itoa(idStr)
	// logs.Debug("urlcrud3:", urlcrud3)

	var ubicacion []map[string]interface{}
	var sede []map[string]interface{}

	// fmt.Println("elemento asdasdadasdfasd: ", Elemento)

	if res, err := request.GetJsonTest(urlcrud3, &ubicacion); err == nil && res.StatusCode == 200 {

		ubicacion2 := ubicacion[0]["EspacioFisicoId"].(map[string]interface{})

		z := strings.Split(fmt.Sprintf("%v", ubicacion2["CodigoAbreviacion"]), "")

		urlcrud4 := "http://" + beego.AppConfig.String("oikosService") + "espacio_fisico?query=CodigoAbreviacion:" + z[0] + z[1] + z[2] + z[3]

		if res, err := request.GetJsonTest(urlcrud4, &sede); err != nil || res.StatusCode != 200 {
			if err == nil {
				err = fmt.Errorf("undesired Status Code: %d", res.StatusCode)
			}
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "TraerElementoSolicitud - request.GetJsonTest(urlcrud4, &sede)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}

		var idElemento int
		if id, err := strconv.Atoi(fmt.Sprintf("%v", Elemento["ElementoCatalogoId"])); err == nil {
			idElemento = id
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "TraerElementoSolicitud - strconv.Atoi(fmt.Sprintf(\"%v\", Elemento[\"ElementoCatalogoId\"]))",
				"err":     err,
				"status":  "400",
			}
			return nil, outputError
		}
		// logs.Debug("elemActa:", elemActa)
		if Elemento___, err := UltimoMovimientoKardex(idElemento); err == nil {

			Elemento___["Sede"] = sede[0]
			Elemento___["Dependencia"] = ubicacion[0]["DependenciaId"]
			Elemento___["Ubicacion"] = ubicacion[0]["EspacioFisicoId"]

			return Elemento___, nil

		} else {
			return nil, err
		}

	} else {
		if err == nil {
			err = fmt.Errorf("undesired Status Code: %d", res.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "TraerElementoSolicitud - request.GetJsonTest(urlcrud3, &ubicacion)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

}

func GetExistenciasKardex() (Elementos []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetExistenciasKardex - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var Elementos___ []map[string]interface{}
	url := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento?"
	url += "query=MovimientoId__FormatoTipoMovimientoId__CodigoAbreviacion__in:AP_KDX|SAL_KDX,Activo:true&limit=-1&fields=ElementoCatalogoId"
	if res, err := request.GetJsonTest(url, &Elementos___); err == nil && res.StatusCode == 200 {
		// fmt.Println("Elementos", Elementos___[0])

		if keys := len(Elementos___[0]); keys != 0 {

			for _, elemento := range Elementos___ {

				var idCatalogo int
				if id, err := strconv.Atoi(fmt.Sprintf("%v", elemento["ElementoCatalogoId"])); err == nil {
					idCatalogo = id
				} else {
					logs.Warn(err)
					continue
				}

				if Elemento, err := UltimoMovimientoKardex(idCatalogo); err == nil {
					if s, ok := Elemento["SaldoCantidad"]; ok {
						if v, ok := s.(float64); ok && v > 0 {
							Elementos = append(Elementos, Elemento)
						}
					}
				}
			}

			return Elementos, nil
		} else {

			return Elementos___, nil
		}

	} else {
		if err == nil {
			err = fmt.Errorf("undesired Status Code: %d", res.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetExistenciasKardex - request.GetJsonTest(url, &Elementos___)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}

func UltimoMovimientoKardex(id_catalogo int) (Elemento_Movimiento map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "UltimoMovimientoKardex - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	if id_catalogo <= 0 {
		err := fmt.Errorf("id MUST be > 0")
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "UltimoMovimientoKardex - id_catalogo <= 0",
			"err":     err,
			"status":  "400",
		}
		return nil, outputError
	}

	idStr := strconv.Itoa(id_catalogo)

	var elemento_catalogo []map[string]interface{}

	// fmt.Println("id asdasdadasdfasd: ", id_catalogo)
	url3 := "http://" + beego.AppConfig.String("catalogoElementosService") + "elemento?query=Id:" + idStr
	// logs.Debug("url3:", url3)
	if res, err := request.GetJsonTest(url3, &elemento_catalogo); err == nil && res.StatusCode == 200 {

		if len(elemento_catalogo) != 1 || len(elemento_catalogo[0]) == 0 {
			err = fmt.Errorf("no hay un elemento del Catalogo de Elementos con id:%s", idStr)
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "UltimoMovimientoKardex - len(elemento_catalogo) != 1 || len(elemento_catalogo[0]) == 0",
				"err":     err,
				"status":  "404",
			}
			return nil, outputError
		}

		// fmt.Println(elemento_catalogo)
		var ultimo_movimiento_kardex []map[string]interface{}
		url4 := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento?query=ElementoCatalogoId:" +
			idStr + ",Activo:true&limit=1&sortby=FechaCreacion&order=desc&fields=ElementoCatalogoId,Id,SaldoCantidad,SaldoValor"
		// logs.Debug("url4:", url4)
		if res, err := request.GetJsonTest(url4, &ultimo_movimiento_kardex); err == nil && res.StatusCode == 200 {

			Elemento := ultimo_movimiento_kardex[0]
			Elemento["ElementoCatalogoId"] = elemento_catalogo[0]
			Elemento["SubgrupoCatalogoId"] = elemento_catalogo[0]["SubgrupoId"]

			return Elemento, nil

		} else {
			if err == nil {
				err = fmt.Errorf("undesired Status Code: %d", res.StatusCode)
			}
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "UltimoMovimientoKardex - request.GetJsonTest(url4, &ultimo_movimiento_kardex)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}

	} else {
		if err == nil {
			err = fmt.Errorf("undesired Status Code: %d", res.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "UltimoMovimientoKardex - request.GetJsonTest(url3, &elemento_catalogo)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}
