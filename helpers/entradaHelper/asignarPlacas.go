package entradaHelper

import (
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

func asignarPlacas(actaRecibidoId int) (elementos []*models.Elemento, outputError map[string]interface{}) {

	funcion := "asignarPlacas - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	var detalle_ []*models.DetalleElemento
	if detalleElementos, err := actaRecibido.GetElementos(actaRecibidoId, nil); err != nil {
		return nil, err
	} else {
		detalle_ = detalleElementos
	}

	var uvt float64
	if uvt_, err := parametros.GetUVTByVigencia(time.Now().Year()); err != nil {
		return nil, err
	} else if uvt_ == 0 {
		return
	} else {
		uvt = uvt_
	}

	var bufferTiposBien = make(map[int]*models.TipoBien, 0)
	for _, el := range detalle_ {

		placa := ""
		if el.TipoBienId != nil {
			bufferTiposBien[el.TipoBienId.Id] = el.TipoBienId
			if el.TipoBienId.NecesitaPlaca {
				if err := generarPlaca(&placa); err != nil {
					return nil, err
				}
			}
		} else {
			if placa_, _, err := checkPlacaElemento(el.SubgrupoCatalogoId.TipoBienId.Id, int(el.ValorUnitario/uvt), bufferTiposBien); err != nil {
				return nil, err
			} else if placa_ {
				if err := generarPlaca(&placa); err != nil {
					return nil, err
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
			EstadoElementoId:   &models.EstadoElemento{Id: el.EstadoElementoId.Id},
			ActaRecibidoId:     &models.ActaRecibido{Id: el.ActaRecibidoId.Id},
			Activo:             true,
		}
		elementos = append(elementos, &elemento_)
	}

	return elementos, nil

}

func checkPlacaElemento(tbPadreId int, normalizado int, bufferTiposBien map[int]*models.TipoBien) (placa bool, err_ string, outputError map[string]interface{}) {

	funcion := "checkPlacaElemento - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	if tbPadreId <= 0 {
		err_ = "La asignación de la clase a los elementos no es correcta."
		return
	}

	for _, tb_ := range bufferTiposBien {
		if tb_.TipoBienPadreId.Id == tbPadreId && tb_.LimiteInferior <= normalizado && normalizado < tb_.LimiteSuperior {
			return tb_.NecesitaPlaca, "", nil
		}
	}

	var tb__ []models.TipoBien
	payload := "limit=1&query=Activo:true,TipoBienPadreId__Id:" + strconv.Itoa(tbPadreId) + ",LimiteInferior__lte:" + strconv.Itoa(normalizado) +
		",LimiteSuperior__gt:" + strconv.Itoa(normalizado)
	if err := catalogoElementos.GetAllTipoBien(payload, &tb__); err != nil {
		return false, "", err
	} else if len(tb__) != 1 {
		err_ = "La asignación de la clase a los elementos no es correcta."
		return
	}

	bufferTiposBien[tb__[0].Id] = &tb__[0]
	return tb__[0].NecesitaPlaca, "", nil
}

func generarPlaca(placa *string) (outputError map[string]interface{}) {

	funcion := "generarPlaca - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	ctxPlaca, _ := beego.AppConfig.Int("contxtPlaca")
	var consecutivo models.Consecutivo

	if err := consecutivos.Get(ctxPlaca, "Registro Placa Arka", &consecutivo); err != nil {
		return err
	}

	year, month, day := time.Now().Date()
	*placa = fmt.Sprintf("%04d%02d%02d%05d", year, month, day, consecutivo.Consecutivo)

	return
}
