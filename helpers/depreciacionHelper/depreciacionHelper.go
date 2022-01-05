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

	// Consulta los valores para generar el corte de depreciación
	if elementos, err := movimientosArkaHelper.GetCorteDepreciacion(info.FechaCorte.AddDate(0, 0, 1).Format("2006-01-02")); err != nil {
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

	if len(totales) == 0 {
		return detalleD, nil
	}

	// Simula la transacción contable en caso de aprobarse
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
		TrContable:   0,
		FechaCorte:   info.FechaCorte.Format("2006-01-02"),
		Totales:      totales,
		RazonRechazo: info.RazonRechazo,
	}

	if detalle_, err := json.Marshal(detalle); err != nil {
		logs.Error(err)
		eval := " - json.Marshal(detalle)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	} else {
		movimiento.Detalle = string(detalle_[:])
	}

	movimiento.Observacion = info.Observaciones
	movimiento.Activo = true

	if info.Id > 0 {
		// Actualiza el registro el registro del movimiento correspondiente a la solicitud
		movimiento.Id = info.Id
		if movimiento_, err := movimientosArkaHelper.PutMovimiento(movimiento, info.Id); err != nil {
			return nil, err
		} else {
			detalleD["Movimiento"] = movimiento_
		}
	} else {
		// Crea el registro del movimiento correspondiente a la solicitud
		if movimiento_, err := movimientosArkaHelper.PostMovimiento(movimiento); err != nil {
			return nil, err
		} else {
			detalleD["Movimiento"] = movimiento_
		}
	}

	return detalleD, nil
}

func GetDepreciacion(id int) (detalleD map[string]interface{}, outputError map[string]interface{}) {

	funcion := "GetDepreciacion"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		detalle    *models.FormatoDepreciacion
		movimiento *models.Movimiento
		// query      string
	)
	detalleD = make(map[string]interface{})

	if mov_, err := movimientosArkaHelper.GetMovimientoById(id); err != nil {
		return nil, err
	} else {
		movimiento = mov_
	}

	if err := json.Unmarshal([]byte(movimiento.Detalle), &detalle); err != nil {
		logs.Error(err)
		eval := " - json.Unmarshal([]byte(movimiento.Detalle), &detalle)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	if trSimulada, err := cuentasContablesHelper.AsientoContable(detalle.Totales, "17", "Depreciación almacén", "", 9445, false); err != nil {
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

// AprobarDepreciacion Registra las novedades para los elementos depreciados y realiza la transaccion contable
func AprobarDepreciacion(id int) (detalleD map[string]interface{}, outputError map[string]interface{}) {

	funcion := "AprobarDepreciacion"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		movimiento *models.Movimiento
		detalle    *models.FormatoDepreciacion
		infoCorte  []*models.DetalleCorteDepreciacion
		fechaCorte time.Time
	)

	detalleD = make(map[string]interface{})
	if mov, err := movimientosArkaHelper.GetMovimientoById(id); err != nil {
		return nil, err
	} else {
		movimiento = mov
	}

	if err := json.Unmarshal([]byte(movimiento.Detalle), &detalle); err != nil {
		logs.Error(err)
		eval := " - json.Unmarshal([]byte(movimiento.Detalle), &detalle)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	// Consulta los valores para generar el corte de depreciación
	if elementos, err := movimientosArkaHelper.GetCorteDepreciacion(detalle.FechaCorte); err != nil {
		return nil, err
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

	if t, err := time.Parse("2006-01-02", detalle.FechaCorte); err != nil {
		logs.Error(err)
		eval := ` - time.Parse("2006-01-02", detalle.FechaCorte)`
		return nil, errorctrl.Error(funcion+eval, err, "500")
	} else {
		fechaCorte = t
	}

	totales := make(map[int]float64)
	novedades := []*models.NovedadElemento{}
	for _, dt := range infoCorte {
		// Determina si al elemento se le debe aplicar depreciación
		if _, ok := subgrupoBien[dt.ElementoActaId]; ok {

			// calcula el valor de la depreciación
			deltaT := GetDeltaTiempo(dt.FechaRef, fechaCorte)
			if deltaT > dt.VidaUtil {
				deltaT = dt.VidaUtil
			}
			depreciacion := (dt.ValorPresente - dt.ValorResidual) * deltaT / dt.VidaUtil

			nov := &models.NovedadElemento{
				VidaUtil:             dt.VidaUtil - deltaT,
				ValorLibros:          dt.ValorPresente - depreciacion,
				ValorResidual:        dt.ValorResidual,
				MovimientoId:         movimiento,
				ElementoMovimientoId: &models.ElementosMovimiento{Id: dt.ElementoMovimientoId},
				Activo:               true,
			}

			novedades = append(novedades, nov)

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

	if len(totales) == 0 {
		return detalleD, nil
	}

	// Registra la transacción contable
	if trContable, err := cuentasContablesHelper.AsientoContable(totales, "17", "Depreciación almacén", "", 9445, true); err != nil {
		return nil, err
	} else {
		detalleD["trContable"] = trContable
		detalle.TrContable = trContable["resultadoTransaccion"].(models.TransaccionMovimientos).ConsecutivoId
		if trContable["errorTransaccion"].(string) != "" {
			return detalleD, nil
		}
	}

	for _, nov := range novedades {
		if _, err := movimientosArkaHelper.PostTrNovedadElemento(nov); err != nil {
			return nil, err
		}
	}

	if sm, err := movimientosArkaHelper.GetAllEstadoMovimiento(url.QueryEscape("Depr Aprobada")); err != nil {
		return nil, err
	} else {
		movimiento.EstadoMovimientoId = sm[0]
	}

	detalle.RazonRechazo = ""
	detalle.Totales = totales

	if detalle_, err := json.Marshal(detalle); err != nil {
		logs.Error(err)
		eval := " - json.Marshal(detalle)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	} else {
		movimiento.Detalle = string(detalle_[:])
	}

	if movimiento_, err := movimientosArkaHelper.PutMovimiento(movimiento, movimiento.Id); err != nil {
		return nil, err
	} else {
		detalleD["Movimiento"] = movimiento_
	}

	return detalleD, nil
}

// GetDeltaTiempo retorna el tiempo en años entre dos fechas
func GetDeltaTiempo(ref, fin time.Time) float64 {

	y1, m1, d1 := fin.Date()
	y0, m0, d0 := ref.Date()

	years := y1 - y0
	months := int(m1 - m0)
	days := d1 - d0

	if days < 0 {
		days += DaysIn(y1, m1-1)
		months--
	}

	if days < 0 {
		days += DaysIn(y1, m1)
		months--
	}

	if months < 0 {
		months += 12
		years--
	}

	prct := float64(years)
	prct += float64(months) / 12
	if days > 0 {
		var firstOfMonth time.Time
		if months > 0 {
			firstOfMonth = time.Date(fin.Year(), fin.Month(), 1, 0, 0, 0, 0, time.UTC)
		} else {
			firstOfMonth = time.Date(ref.Year(), ref.Month(), 1, 0, 0, 0, 0, time.UTC)
		}
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
		prct += (float64(days) / float64(lastOfMonth.Day())) / 12
	}

	return prct
}

func DaysIn(year int, month time.Month) int {
	return time.Date(year, month, 0, 0, 0, 0, 0, time.UTC).Day()
}
