package depreciacionHelper

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/cuentasContablesHelper"
	"github.com/udistrital/arka_mid/helpers/movimientosArkaHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// GenerarTrDepreciacion Calcula la transacción contable que se generará
func GenerarTrDepreciacion(info *models.InfoDepreciacion) (detalleD map[string]interface{}, outputError map[string]interface{}) {

	funcion := "GenerarTrDepreciacion"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")
	var (
		infoCorte  []*models.DetalleCorteDepreciacion
		movimiento *models.Movimiento
		query      string
	)

	detalleD = make(map[string]interface{})
	info.FechaCorte = info.FechaCorte.AddDate(0, 0, 1)
	stringFecha := info.FechaCorte.Format("2006-01-02")

	// Consulta los valores para generar el corte de depreciación
	if elementos, err := movimientosArkaHelper.GetCorteDepreciacion(stringFecha); err != nil {
		return nil, err
	} else if elementos == nil {
		return detalleD, nil
	} else {
		infoCorte = elementos
	}

	// Consulta los subgrupos de cada elemento o novedad
	idsActa := []int{}
	for _, val := range infoCorte {
		idsActa = append(idsActa, int(val.ElementoActaId))
	}

	subgrupoBien := make(map[int]int)
	if elemento_, err := actaRecibido.GetElementos(0, idsActa); err != nil {
		return nil, err
	} else {
		for _, el := range elemento_ {
			// Determina qué elementos se deben depreciar de acuerdo a la parametrización del tipo de bien
			// Asociar subgrupo a los elementos que requieren depreciacion
			if el.SubgrupoCatalogoId.Depreciacion {
				subgrupoBien[el.Id] = el.SubgrupoCatalogoId.SubgrupoId.Id
			}
		}
	}
	totales := make(map[int]float64)
	for _, dt := range infoCorte {
		// Determina si al elemento se le debe aplicar depreciación
		if _, ok := subgrupoBien[dt.ElementoActaId]; ok {
			// calcula el valor de la depreciación
			deltaT := GetDeltaTiempo(dt.FechaRef, info.FechaCorte)
			if deltaT > dt.VidaUtil {
				deltaT = dt.VidaUtil
			}
			depreciacion := (dt.ValorPresente - dt.ValorResidual) * deltaT / dt.VidaUtil

			// Agrupa las cantidades por subgrupo
			x := float64(0)
			if val, ok := totales[subgrupoBien[dt.ElementoActaId]]; ok {
				x = val + depreciacion
			} else {
				x = depreciacion
			}
			totales[subgrupoBien[dt.ElementoActaId]] = x
		}
	}

	// Genera la transacción simulada de acuerdo a los valores obtenidos
	if len(totales) == 0 {
		return detalleD, nil
	}
	if trSimulada, err := cuentasContablesHelper.AsientoContable(totales, "17", "Depreciación almacén", "", 9445, false); err != nil {
		return nil, outputError
	} else {
		detalleD["trContable"] = trSimulada
		if trSimulada["errorTransaccion"].(string) != "" {
			return detalleD, nil
		}
	}

	movimiento = new(models.Movimiento)
	if sm, err := movimientosArkaHelper.GetAllEstadoMovimiento(url.QueryEscape("Depr Generada")); err != nil {
		return nil, err
	} else {
		movimiento.EstadoMovimientoId = sm[0]
	}

	query = "query=Nombre:Depreciación"
	if fm, err := movimientosArkaHelper.GetAllFormatoTipoMovimiento(query); err != nil {
		return nil, err
	} else {
		movimiento.FormatoTipoMovimientoId = fm[0]
	}

	detalle := models.FormatoDepreciacion{
		TrContable: 0,
		FechaCorte: stringFecha,
		Totales:    totales,
	}

	if detalle_, err := json.Marshal(detalle); err != nil {
		logs.Error(err)
		eval := " - json.Marshal(detalle)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	} else {
		movimiento.Detalle = string(detalle_[:])
	}

	movimiento.Observacion = info.Observacion
	movimiento.Activo = true

	// Crea el registro del movimiento correspondiente a la solicitud
	if movimiento_, err := movimientosArkaHelper.PostMovimiento(movimiento); err != nil {
		return nil, err
	} else {
		detalleD["Movimiento"] = movimiento_
	}

	return detalleD, nil
}

// GetDeltaTiempo retorna el tiempo en años entre dos fechas
func GetDeltaTiempo(ref time.Time, fin time.Time) float64 {

	prct := 0.0
	y1, m1, d1 := fin.Date()
	y0, m0, d0 := ref.Date()
	year := int(y1 - y0)
	month := int(m1 - m0)
	day := int(d1 - d0)
	prct += float64(year)

	if day < 0 {
		// days in month:
		t := time.Date(y1, m1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}

	if month < 0 {
		month += 12
		year--
	}

	prct += float64(month) / 12
	if day > 0 {
		firstOfMonth := time.Date(fin.Year(), fin.Month(), 1, 0, 0, 0, 0, time.UTC)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
		prct += (float64(day) / float64(lastOfMonth.Day())) / 12
	}
	return prct
}
