package depreciacionHelper

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	crudMovimientosArka "github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	crudTerceros "github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// GenerarTrDepreciacion Calcula la transacción contable que se generará
func GenerarTrDepreciacion(info *models.InfoDepreciacion) (detalleD map[string]interface{}, outputError map[string]interface{}) {

	funcion := "GenerarTrDepreciacion"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		infoCorte     []*models.DetalleCorteDepreciacion
		movimiento    *models.Movimiento
		consecutivoId int
		terceroUD     int
		consecutivo   string
		query         string
	)

	detalleD = make(map[string]interface{})

	// Consulta los valores para generar el corte de depreciación
	if elementos, err := crudMovimientosArka.GetCorteDepreciacion(info.FechaCorte.AddDate(0, 0, 1).Format("2006-01-02")); err != nil {
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
			if info.Tipo == "Depreciación" && el.SubgrupoCatalogoId.Depreciacion {
				subgrupoBien[el.Id] = el.SubgrupoCatalogoId.SubgrupoId.Id
			} else if info.Tipo == "Amortizacion" && el.SubgrupoCatalogoId.Amortizacion {
				subgrupoBien[el.Id] = el.SubgrupoCatalogoId.SubgrupoId.Id
			}
		}
	}

	totales := make(map[int]float64)
	for _, dt := range infoCorte {
		// Determina si al elemento se le debe aplicar depreciación
		if _, ok := subgrupoBien[dt.ElementoActaId]; ok {

			// Calcula el valor de la depreciación
			if dt.NovedadElementoId > 0 {
				dt.FechaRef = dt.FechaRef.AddDate(0, 0, 1)
			}
			deltaT := GetDeltaTiempo(dt.FechaRef, info.FechaCorte.AddDate(0, 0, 1))
			if deltaT > dt.VidaUtil {
				deltaT = dt.VidaUtil
			}

			depreciacion := (dt.ValorPresente - dt.ValorResidual) * deltaT / dt.VidaUtil
			if dt.ValorPresente-depreciacion < 0 {
				depreciacion = dt.ValorPresente
			}

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

	query = "query=TipoDocumentoId__Nombre:NIT,Numero:" + crudTerceros.GetDocUD()
	if terceroUD_, err := crudTerceros.GetAllDatosIdentificacion(query); err != nil {
		return nil, err
	} else {
		terceroUD = terceroUD_[0].TerceroId.Id
	}

	movimiento = new(models.Movimiento)

	query = "query=Nombre:" + info.Tipo
	if fm, err := crudMovimientosArka.GetAllFormatoTipoMovimiento(query); err != nil {
		return nil, err
	} else {
		movimiento.FormatoTipoMovimientoId = fm[0]
	}

	// Simula la transacción contable en caso de aprobarse
	if trSimulada, err := asientoContable.AsientoContable(totales, "",
		strconv.Itoa(movimiento.FormatoTipoMovimientoId.Id), getDescripcionMovmientoCierre(), descAsiento(), terceroUD, 0, false); err != nil {
		return nil, outputError
	} else {
		detalleD["trContable"] = trSimulada
		if trSimulada["errorTransaccion"].(string) != "" {
			return detalleD, nil
		}
	}

	if sm, err := crudMovimientosArka.GetAllEstadoMovimiento("query=Nombre:" + url.QueryEscape("Depr Generada")); err != nil {
		return nil, err
	} else {
		movimiento.EstadoMovimientoId = sm[0]
	}

	if info.Id == 0 {
		ctxt, _ := beego.AppConfig.Int("contxtMedicionesCons")
		if consecutivo_, consecutivoId_, err := consecutivos.Get("%02.0f", ctxt, "Registro "+info.Tipo+" Arka"); err != nil {
			return nil, outputError
		} else {
			consecutivoId = consecutivoId_
			consecutivo = consecutivos.Format(getTipoComprobanteCierre()+"-", consecutivo_, fmt.Sprintf("%s%02d", "-", time.Now().Year()))
		}
	} else {
		var (
			detalle_    models.FormatoDepreciacion
			movimiento_ models.Movimiento
		)

		if mov_, err := crudMovimientosArka.GetMovimientoById(info.Id); err != nil {
			return nil, err
		} else {
			movimiento_ = *mov_
		}

		if err := json.Unmarshal([]byte(movimiento_.Detalle), &detalle_); err != nil {
			logs.Error(err)
			eval := " - json.Unmarshal([]byte(movimiento_.Detalle), &detalle_)"
			return nil, errorctrl.Error(funcion+eval, err, "500")
		}

		consecutivo, consecutivoId = detalle_.Consecutivo, detalle_.ConsecutivoId
	}

	detalle := models.FormatoDepreciacion{
		FechaCorte:    info.FechaCorte.Format("2006-01-02"),
		Totales:       totales,
		RazonRechazo:  info.RazonRechazo,
		ConsecutivoId: consecutivoId,
		Consecutivo:   consecutivo,
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
		if movimiento_, err := crudMovimientosArka.PutMovimiento(movimiento, info.Id); err != nil {
			return nil, err
		} else {
			detalleD["Movimiento"] = movimiento_
		}
	} else {
		// Crea el registro del movimiento correspondiente a la solicitud
		if movimiento_, err := crudMovimientosArka.PostMovimiento(movimiento); err != nil {
			return nil, err
		} else {
			detalleD["Movimiento"] = movimiento_
		}
	}

	return detalleD, nil
}

func getTipoComprobanteCierre() string {
	return "H22"
}

func getDescripcionMovmientoCierre() string {
	return "Mediciones posteriores"
}
