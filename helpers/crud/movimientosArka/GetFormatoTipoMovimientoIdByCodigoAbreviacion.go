package movimientosArka

import (
	"net/url"

	"github.com/beego/beego/v2/core/logs"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// GetFormatoTipoMovimientoIdByCodigoAbreviacion Consulta el Id de un FormatoTipoMovimiento según el Codigo de abreviación del mismo
func GetFormatoTipoMovimientoIdByCodigoAbreviacion(id *int, codigoAbreviacion string) (outputError map[string]interface{}) {

	funcion := "GetFormatoTipoMovimientoIdByCodigoAbreviacion"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	query := "query=CodigoAbreviacion:" + url.QueryEscape(codigoAbreviacion)
	if fm, err := GetAllFormatoTipoMovimiento(query); err != nil {
		return err
	} else if len(fm) == 0 {
		err := "No se encuentra el formato tipo movimiento: " + codigoAbreviacion
		logs.Error(err)
		eval := " - GetAllFormatoTipoMovimiento(query)"
		return errorCtrl.Error(funcion+eval, err, "500")
	} else {
		*id = fm[0].Id
	}

	return
}

// GetEstadoMovimientoIdByNombre Consulta el Id de un EstadoMovimiento según el nombre del mismo
func GetEstadoMovimientoIdByNombre(id *int, nombre string) (outputError map[string]interface{}) {

	funcion := "GetEstadoMovimientoIdByNombre"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	query := "query=Nombre:" + url.QueryEscape(nombre)
	if em, err := GetAllEstadoMovimiento(query); err != nil {
		return err
	} else if len(em) == 0 {
		err := "No se encuentra el estado movimiento: " + nombre
		logs.Error(err)
		eval := " - GetAllEstadoMovimiento(query)"
		return errorCtrl.Error(funcion+eval, err, "500")
	} else {
		*id = em[0].Id
	}

	return
}
