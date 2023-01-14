package bodegaConsumoHelper

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/oikos"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
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

func UltimoMovimientoKardex(elementoId int) (detalleElemento map[string]interface{}, outputError map[string]interface{}) {

	funcion := "UltimoMovimientoKardex - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	if elementoId <= 0 {
		err := fmt.Errorf("id MUST be > 0")
		logs.Error(err)
		eval := "id_catalogo <= 0"
		return nil, errorctrl.Error(funcion+eval, err, "400")
	}

	detalleElemento = make(map[string]interface{})
	idStr := strconv.Itoa(elementoId)

	var elemento_ []models.ElementoCatalogo
	outputError = catalogoElementos.GetAllElemento("query=Id:"+idStr, &elemento_)
	if outputError != nil || len(elemento_) != 1 {
		return
	}

	payload := "limit=1&sortby=FechaCreacion&order=desc&fields=ElementoCatalogoId,Id,SaldoCantidad,SaldoValor&query=ElementoCatalogoId:" + idStr
	elemento, err := movimientosArka.GetAllElementosMovimiento(payload)
	if err != nil || len(elemento) != 1 {
		return nil, err
	}

	outputError = utilsHelper.FillStruct(elemento[0], &detalleElemento)
	if outputError != nil {
		return
	}

	detalleElemento["ElementoCatalogoId"] = elemento_[0]
	detalleElemento["SubgrupoCatalogoId"] = elemento_[0].SubgrupoId

	return
}
