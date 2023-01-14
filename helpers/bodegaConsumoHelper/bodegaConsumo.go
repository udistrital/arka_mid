package bodegaConsumoHelper

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/oikos"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

func TraerElementoSolicitud(Elemento models.ElementoSolicitud_) (Elemento_ map[string]interface{}, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("TraerElementoSolicitud - Unhandled Error", "500")

	ubicacionInfo, err := oikos.GetSedeDependenciaUbicacion(Elemento.Ubicacion)
	if err != nil {
		return nil, err
	}

	if Elemento___, err := UltimoMovimientoKardex(Elemento.ElementoCatalogoId); err == nil {

		Elemento___["Sede"] = ubicacionInfo.Sede
		Elemento___["Dependencia"] = ubicacionInfo.Dependencia
		Elemento___["Ubicacion"] = ubicacionInfo.Ubicacion.EspacioFisicoId

		return Elemento___, nil

	} else {
		return nil, err
	}

}

func GetExistenciasKardex() (Elementos []map[string]interface{}, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GetExistenciasKardex - Unhandled Error!", "500")

	var Elementos___ []*models.ElementosMovimiento
	url := "query=MovimientoId__FormatoTipoMovimientoId__CodigoAbreviacion__in:AP_KDX|SAL_KDX," +
		"Activo:true,ElementoCatalogoId__gt:0&limit=-1&fields=ElementoCatalogoId"

	if elementos_, err := movimientosArka.GetAllElementosMovimiento(url); err != nil {
		return nil, err
	} else {
		Elementos___ = elementos_
	}

	if len(Elementos___) > 0 {

		for _, elemento := range Elementos___ {
			if Elemento, err := UltimoMovimientoKardex(elemento.ElementoCatalogoId); err == nil {
				if s, ok := Elemento["SaldoCantidad"]; ok {
					if v, ok := s.(float64); ok && v > 0 {
						Elementos = append(Elementos, Elemento)
					}
				}
			}
		}

	}

	return Elementos, nil

}

func UltimoMovimientoKardex(id_catalogo int) (Elemento_Movimiento map[string]interface{}, outputError map[string]interface{}) {

	funcion := "UltimoMovimientoKardex - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	if id_catalogo <= 0 {
		err := fmt.Errorf("id MUST be > 0")
		logs.Error(err)
		eval := " - id_catalogo <= 0"
		return nil, errorctrl.Error(funcion+eval, err, "400")
	}

	idStr := strconv.Itoa(id_catalogo)

	var elemento_catalogo []map[string]interface{}

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
