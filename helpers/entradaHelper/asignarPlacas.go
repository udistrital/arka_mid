package entradaHelper

import (
	"fmt"
	"strconv"
	"time"

	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

func asignarPlacas(actaRecibidoId int, elementos *[]*models.Elemento) (errMsg string, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("asignarPlacas - Unhandled Error!", "500")

	var detalle_ []*models.DetalleElemento
	if detalleElementos, err := actaRecibido.GetElementos(actaRecibidoId, nil); err != nil {
		return "", err
	} else {
		detalle_ = detalleElementos
	}

	var uvt float64 = 1
	// if uvt_, err := parametros.GetUVTByVigencia(time.Now().Year()); err != nil {
	// 	return "", err
	// } else if uvt_ == 0 {
	// 	return "No se pudo consultar el valor del UVT. Intente más tarde o contacte soporte.", nil
	// } else {
	// 	uvt = uvt_
	// }

	var bufferTiposBien = make(map[int]*models.TipoBien)
	for _, el := range detalle_ {

		placa := ""
		tipoBien := 0
		if el.TipoBienId != nil {
			if el.TipoBienId.TipoBienPadreId.Id != el.SubgrupoCatalogoId.TipoBienId.Id {
				return "El tipo bien asignado manualmente no corresponde a la clase correspondiente.", nil
			}
			tipoBien = el.TipoBienId.Id
			bufferTiposBien[el.TipoBienId.Id] = el.TipoBienId
			if el.TipoBienId.NecesitaPlaca {
				if err := generarPlaca(&placa); err != nil {
					return "", err
				}
			}
		} else {
			if placa_, msj, err := checkPlacaElemento(el.SubgrupoCatalogoId.TipoBienId.Id, el.ValorUnitario/uvt, bufferTiposBien); err != nil || msj != "" {
				return msj, err
			} else if placa_ {
				if err := generarPlaca(&placa); err != nil {
					return "", err
				}
			}
		}

		elemento_ := models.Elemento{
			Id:                 el.Id,
			Nombre:             el.Nombre,
			Cantidad:           el.Cantidad,
			Marca:              el.Marca,
			Serie:              el.Serie,
			UnidadMedida:       el.UnidadMedida,
			ValorUnitario:      el.ValorUnitario,
			Subtotal:           el.Subtotal,
			Descuento:          el.Descuento,
			ValorTotal:         el.ValorTotal,
			PorcentajeIvaId:    el.PorcentajeIvaId,
			ValorIva:           el.ValorIva,
			ValorFinal:         el.ValorFinal,
			Placa:              placa,
			SubgrupoCatalogoId: el.SubgrupoCatalogoId.SubgrupoId.Id,
			TipoBienId:         tipoBien,
			EstadoElementoId:   &models.EstadoElemento{Id: el.EstadoElementoId.Id},
			ActaRecibidoId:     &models.ActaRecibido{Id: el.ActaRecibidoId.Id},
			Activo:             true,
		}
		*elementos = append(*elementos, &elemento_)
	}

	return

}

func checkPlacaElemento(tbPadreId int, normalizado float64, bufferTiposBien map[int]*models.TipoBien) (placa bool, errMsg string, outputError map[string]interface{}) {

	funcion := "checkPlacaElemento - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	if tbPadreId <= 0 {
		errMsg = "La asignación de la clase a los elementos no es correcta."
		return
	}

	for _, tb_ := range bufferTiposBien {
		if tb_.TipoBienPadreId.Id == tbPadreId && tb_.LimiteInferior <= normalizado && normalizado < tb_.LimiteSuperior {
			return tb_.NecesitaPlaca, "", nil
		}
	}

	var tb__ []models.TipoBien
	payload := "limit=1&query=Activo:true,TipoBienPadreId__Id:" + strconv.Itoa(tbPadreId) + ",LimiteInferior__lte:" + fmt.Sprintf("%f", normalizado) +
		",LimiteSuperior__gt:" + fmt.Sprintf("%f", normalizado)
	if err := catalogoElementos.GetAllTipoBien(payload, &tb__); err != nil {
		return false, "", err
	} else if len(tb__) != 1 {
		errMsg = "La asignación de la clase a los elementos no es correcta."
		return
	}

	bufferTiposBien[tb__[0].Id] = &tb__[0]
	return tb__[0].NecesitaPlaca, "", nil
}

func generarPlaca(placa *string) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("generarPlaca - Unhandled Error!", "500")

	var consecutivo models.Consecutivo
	if err := consecutivos.Get("contxtPlaca", "Registro Placa Arka", &consecutivo); err != nil {
		return err
	}

	year, month, day := time.Now().Date()
	*placa = fmt.Sprintf("%04d%02d%02d%05d", year, month, day, consecutivo.Consecutivo)

	return
}
