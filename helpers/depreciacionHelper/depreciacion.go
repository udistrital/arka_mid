package depreciacionHelper

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/asientoContable"
	crudMovimientosArka "github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	crudTerceros "github.com/udistrital/arka_mid/helpers/crud/terceros"

	// "github.com/udistrital/arka_mid/helpers/asientoContable"
	// "github.com/udistrital/arka_mid/helpers/movimientosArkaHelper"
	// "github.com/udistrital/arka_mid/helpers/tercerosHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

func GetDepreciacion(id int) (detalleD map[string]interface{}, outputError map[string]interface{}) {

	funcion := "GetDepreciacion"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		detalle    *models.FormatoDepreciacion
		movimiento *models.Movimiento
		terceroUD  int
	)
	detalleD = make(map[string]interface{})

	if mov_, err := crudMovimientosArka.GetMovimientoById(id); err != nil {
		return nil, err
	} else {
		movimiento = mov_
	}

	if err := json.Unmarshal([]byte(movimiento.Detalle), &detalle); err != nil {
		logs.Error(err)
		eval := " - json.Unmarshal([]byte(movimiento.Detalle), &detalle)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	query := "query=TipoDocumentoId__Nombre:NIT,Numero:" + crudTerceros.GetDocUD()
	if terceroUD_, err := crudTerceros.GetAllDatosIdentificacion(query); err != nil {
		return nil, err
	} else {
		terceroUD = terceroUD_[0].TerceroId.Id
	}

	if trSimulada, err := asientoContable.AsientoContable(detalle.Totales, strconv.Itoa(movimiento.FormatoTipoMovimientoId.Id), "", descAsiento(), terceroUD, false); err != nil {
		return nil, outputError
	} else {
		detalleD["TrContable"] = trSimulada
		if trSimulada["errorTransaccion"].(string) != "" {
			return detalleD, nil
		}
	}

	detalleD["Movimiento"] = movimiento
	return detalleD, nil
}

// GetDeltaTiempo retorna el tiempo en años entre dos fechas
func GetDeltaTiempo(ref, fin time.Time) (prct float64) {

	ref = time.Date(ref.Year(), ref.Month(), ref.Day(), 0, 0, 0, 0, time.UTC)
	fin = time.Date(fin.Year(), fin.Month(), fin.Day(), 0, 0, 0, 0, time.UTC)

	prct = fin.Sub(ref).Hours() / (24 * 365)

	return prct
}

// CalculaDp Genera el valor y el tiempo en años a depreciar
// ref: Fecha de referencia para determinar el tiempo por el cual se correra la depreciacion.
// dp: Valor calculado correspondiente a la depreciacion.
func CalculaDp(presente, residual, vUtil float64, ref, fCorte time.Time) (dp, deltaT float64) {

	if vUtil == 0 {
		return 0, 0
	}

	deltaT = GetDeltaTiempo(ref, fCorte.AddDate(0, 0, 1))
	if deltaT > vUtil {
		deltaT = vUtil
	}

	if residual > presente {
		presente = residual
	}

	dp = (presente - residual) * deltaT / vUtil
	return dp, deltaT

}

// GetDetalleDepreciacion Consulta el detalle de una medición determinada
func GetDetalleDepreciacion(detalle string) (detalle_ *models.FormatoDepreciacion, outputError map[string]interface{}) {

	funcion := "GetDetalleDepreciacion"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	if err := json.Unmarshal([]byte(detalle), &detalle_); err != nil {
		logs.Error(err)
		eval := " - json.Unmarshal([]byte(detalle), &detalle_)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	return detalle_, nil
}

func descAsiento() string {
	return "Depreciación almacén"
}
