package salidaHelper

import (
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

func GetSalidaById(id int) (Salida map[string]interface{}, outputError map[string]interface{}) {

	funcion := "GetSalidaById - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	var (
		trSalida      *models.TrSalida
		detalle       map[string]interface{}
		formato       models.FormatoSalida
		ids           []int
		elementosActa []*models.DetalleElemento
	)

	if tr_, err := movimientosArka.GetTrSalida(id); err != nil {
		return nil, err
	} else if tr_.Salida.FormatoTipoMovimientoId.CodigoAbreviacion == "SAL" ||
		tr_.Salida.FormatoTipoMovimientoId.CodigoAbreviacion == "SAL_CONS" {
		trSalida = tr_
	} else {
		return
	}

	for _, el := range trSalida.Elementos {
		ids = append(ids, el.ElementoActaId)
	}

	if len(ids) > 0 {
		if elementosActa, outputError = actaRecibido.GetElementos(0, ids); outputError != nil {
			return nil, outputError
		}
	}

	var elementosCompletos = make([]models.DetalleElementoSalida, 0)
	for _, el := range elementosActa {

		if idx := utilsHelper.FindElementoInArrayElementosMovimiento(trSalida.Elementos, el.Id); idx > -1 {

			detalle := models.DetalleElementoSalida{
				Cantidad:           el.Cantidad,
				ElementoActaId:     el.Id,
				Id:                 trSalida.Elementos[idx].Id,
				Marca:              el.Marca,
				Nombre:             el.Nombre,
				Placa:              el.Placa,
				Serie:              el.Serie,
				SubgrupoCatalogoId: el.SubgrupoCatalogoId,
				TipoBienId:         el.TipoBienId,
				ValorResidual:      (trSalida.Elementos[idx].ValorResidual * 10000) / (trSalida.Elementos[idx].ValorTotal * 100),
				ValorUnitario:      el.ValorUnitario,
				ValorTotal:         trSalida.Elementos[idx].ValorTotal,
				VidaUtil:           trSalida.Elementos[idx].VidaUtil,
			}

			elementosCompletos = append(elementosCompletos, detalle)
		}

	}

	if err := utilsHelper.Unmarshal(trSalida.Salida.Detalle, &formato); err != nil {
		return nil, err
	}

	if salida__, err := TraerDetalle(trSalida.Salida, formato, nil, nil, nil); err != nil {
		return nil, err
	} else {
		detalle = salida__
	}

	Salida_final := map[string]interface{}{
		"Elementos": elementosCompletos,
		"Salida":    detalle,
	}

	if trSalida.Salida.EstadoMovimientoId.Nombre == "Salida Aprobada" && formato.ConsecutivoId > 0 {
		if tr, err := movimientosContables.GetTransaccion(formato.ConsecutivoId, "consecutivo", true); err != nil {
			return nil, err
		} else if len(tr.Movimientos) > 0 {
			if detalleContable, err := asientoContable.GetDetalleContable(tr.Movimientos, nil); err != nil {
				return nil, err
			} else {
				trContable := models.InfoTransaccionContable{
					Movimientos: detalleContable,
					Concepto:    tr.Descripcion,
					Fecha:       tr.FechaTransaccion,
				}
				Salida_final["TransaccionContable"] = trContable
			}
		}
	}

	return Salida_final, nil

}
