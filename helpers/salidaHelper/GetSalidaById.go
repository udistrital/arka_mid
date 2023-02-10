package salidaHelper

import (
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

func GetSalidaById(id int) (Salida map[string]interface{}, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GetSalidaById - Unhandled Error!", "500")

	var (
		formato       models.FormatoSalida
		ids           []int
		elementosActa []*models.DetalleElemento
	)

	trSalida, outputError := movimientosArka.GetTrSalida(id)
	if outputError != nil || (trSalida.Salida.FormatoTipoMovimientoId.CodigoAbreviacion != "SAL" && trSalida.Salida.FormatoTipoMovimientoId.CodigoAbreviacion != "SAL_CONS") {
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

	outputError = utilsHelper.Unmarshal(trSalida.Salida.Detalle, &formato)
	if outputError != nil {
		return
	}

	detalle, outputError := TraerDetalle(trSalida.Salida, formato, nil, nil, nil)
	if outputError != nil {
		return
	}

	Salida = map[string]interface{}{
		"Elementos": elementosCompletos,
		"Salida":    detalle,
	}

	if trSalida.Salida.EstadoMovimientoId.Nombre == "Salida Aprobada" && trSalida.Salida.ConsecutivoId != nil && *trSalida.Salida.ConsecutivoId > 0 {
		Salida["TransaccionContable"] = models.InfoTransaccionContable{}
		Salida["TransaccionContable"], outputError = asientoContable.GetFullDetalleContable(*trSalida.Salida.ConsecutivoId)
	}

	return
}
