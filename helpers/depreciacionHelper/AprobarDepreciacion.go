package depreciacionHelper

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/asientoContable"
	crudMovimientosArka "github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	crudTerceros "github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// AprobarDepreciacion Registra las novedades para los elementos depreciados y realiza la transaccion contable
func AprobarDepreciacion(id int) (detalleD map[string]interface{}, outputError map[string]interface{}) {

	funcion := "AprobarDepreciacion"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		movimiento *models.Movimiento
		detalle    *models.FormatoDepreciacion
		infoCorte  []*models.DetalleCorteDepreciacion
		fechaCorte time.Time
		terceroUD  int
	)

	detalleD = make(map[string]interface{})

	if mov, err := crudMovimientosArka.GetAllMovimiento("query=Id:" + strconv.Itoa(id)); err != nil {
		return nil, err
	} else {
		movimiento = mov[0]
	}

	if err := json.Unmarshal([]byte(movimiento.Detalle), &detalle); err != nil {
		logs.Error(err)
		eval := " - json.Unmarshal([]byte(movimiento.Detalle), &detalle)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	// Consulta los valores para generar el corte de depreciación
	if t, err := time.Parse("2006-01-02", detalle.FechaCorte); err != nil {
		logs.Error(err)
		eval := ` - time.Parse("2006-01-02", detalle.FechaCorte)`
		return nil, errorctrl.Error(funcion+eval, err, "500")
	} else {
		fechaCorte = t
	}

	if elementos, err := crudMovimientosArka.GetCorteDepreciacion(fechaCorte.AddDate(0, 0, 1).Format("2006-01-02")); err != nil {
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
			if movimiento.FormatoTipoMovimientoId.Nombre == "Depreciación" && el.SubgrupoCatalogoId.Depreciacion {
				subgrupoBien[el.Id] = el.SubgrupoCatalogoId.SubgrupoId.Id
			} else if movimiento.FormatoTipoMovimientoId.Nombre == "Amortizacion" && el.SubgrupoCatalogoId.Amortizacion {
				subgrupoBien[el.Id] = el.SubgrupoCatalogoId.SubgrupoId.Id
			}
		}
	}

	totales := make(map[int]float64)
	novedades := []*models.NovedadElemento{}
	for _, dt := range infoCorte {
		// Determina si al elemento se le debe aplicar depreciación
		if _, ok := subgrupoBien[dt.ElementoActaId]; ok {

			// Calcula el valor de la depreciación
			if dt.NovedadElementoId > 0 {
				dt.FechaRef = dt.FechaRef.AddDate(0, 0, 1)
			}
			deltaT := GetDeltaTiempo(dt.FechaRef, fechaCorte.AddDate(0, 0, 1))
			if deltaT > dt.VidaUtil {
				deltaT = dt.VidaUtil
			}

			depreciacion := (dt.ValorPresente - dt.ValorResidual) * deltaT / dt.VidaUtil
			if dt.ValorPresente-depreciacion < 0 {
				depreciacion = dt.ValorPresente
			}

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

	query := "query=TipoDocumentoId__Nombre:NIT,Numero:" + crudTerceros.GetDocUD()
	if terceroUD_, err := crudTerceros.GetAllDatosIdentificacion(query); err != nil {
		return nil, err
	} else {
		terceroUD = terceroUD_[0].TerceroId.Id
	}

	// Registra la transacción contable
	if trContable, err := asientoContable.AsientoContable(totales, strconv.Itoa(movimiento.FormatoTipoMovimientoId.Id), "", descAsiento(), terceroUD, true); err != nil {
		return nil, err
	} else {
		detalleD["trContable"] = trContable
		detalle.TrContable = trContable["resultadoTransaccion"].(*models.DetalleTrContable).ConsecutivoId
		if trContable["errorTransaccion"].(string) != "" {
			return detalleD, nil
		}
	}

	for _, nov := range novedades {
		if _, err := crudMovimientosArka.PostTrNovedadElemento(nov); err != nil {
			return nil, err
		}
	}

	if sm, err := crudMovimientosArka.GetAllEstadoMovimiento(url.QueryEscape("Depr Aprobada")); err != nil {
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

	if movimiento_, err := crudMovimientosArka.PutMovimiento(movimiento, movimiento.Id); err != nil {
		return nil, err
	} else {
		detalleD["Movimiento"] = movimiento_
	}

	return detalleD, nil
}
