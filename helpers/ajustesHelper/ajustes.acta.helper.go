package ajustesHelper

import (
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
)

// determinarDeltaActa Determina el tipo de ajuste del elemento.
// org: Elementos antes del ajuste.
// nvo: Elementos editados.
// msc: Cambios miscelaneos, elementos a los que unicamente se les debe ajustar nombre, marca, serie, unidad.
// vls: Cambios a valores, elementos a los que se les debe cambiar el valor total.
// sg: Cambia el subgrupo del elemento. Se ajusta la placa de acuerdo al nuevo subgrupo.
func determinarDeltaActa(org *models.Elemento, nvo *models.DetalleElemento_) (msc, vls, sg bool, outputError map[string]interface{}) {

	funcion := "determinarDeltaActa"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	if org.SubgrupoCatalogoId != nvo.SubgrupoCatalogoId {

		urlcrud := "fields=TipoBienId&sortby=Id&order=desc&query=Activo:true,SubgrupoId__Id:" + strconv.Itoa(nvo.SubgrupoCatalogoId)
		if detalleSubgrupo_, err := catalogoElementos.GetAllDetalleSubgrupo(urlcrud); err != nil {
			return false, false, false, err
		} else if len(detalleSubgrupo_) == 0 {
			err := "len(detalleSubgrupo_) = 0"
			eval := " - catalogoElementosHelper.GetAllDetalleSubgrupo(urlcrud)"
			return false, false, false, errorctrl.Error(funcion+eval, err, "500")
		} else {
			if detalleSubgrupo_[0].TipoBienId.NecesitaPlaca && nvo.Placa == "" {
				var consecutivo models.Consecutivo
				ctxPlaca, _ := beego.AppConfig.Int("contxtPlaca")
				if err := consecutivos.Get(ctxPlaca, "Registro Placa Arka", &consecutivo); err != nil {
					return false, false, false, err
				} else {
					year, month, day := time.Now().Date()
					nvo.Placa = fmt.Sprintf("%04d%02d%02d%05d", year, month, day, consecutivo.Consecutivo)
				}
			} else if !detalleSubgrupo_[0].TipoBienId.NecesitaPlaca && nvo.Placa != "" {
				nvo.Placa = ""
			}

		}

		nvo.Activo = true
		sg = true

	} else if org.ValorTotal != nvo.ValorTotal {
		nvo.Activo = true
		vls = true

	} else if org.Nombre != nvo.Nombre || org.Marca != nvo.Marca ||
		org.Serie != nvo.Serie || org.UnidadMedida != nvo.UnidadMedida ||
		org.Cantidad != nvo.Cantidad || org.ValorUnitario != nvo.ValorUnitario ||
		org.Subtotal != nvo.Subtotal || org.Descuento != nvo.Descuento ||
		org.PorcentajeIvaId != nvo.PorcentajeIvaId ||
		org.ValorIva != nvo.ValorIva || org.ValorFinal != nvo.ValorFinal {

		nvo.Activo = true
		msc = true

	}

	return msc, vls, sg, nil

}

// fillElementos Consulta el detalle de los subgrupos
func fillElementos(elsOrg []*models.DetalleElemento_) (completos []*models.DetalleElemento__, outputError map[string]interface{}) {

	funcion := "fillElementos"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		ids       []int
		query     string
		subgrupos map[int]*models.DetalleSubgrupo
	)

	for _, el := range elsOrg {
		ids = append(ids, el.SubgrupoCatalogoId)
	}

	query = "fields=SubgrupoId,TipoBienId,Depreciacion,Amortizacion,ValorResidual,VidaUtil&sortby=Id&order=desc"
	query += "&query=Activo:true,SubgrupoId__Id__in:" + utilsHelper.ArrayToString(ids, "|")
	if sg, err := catalogoElementos.GetAllDetalleSubgrupo(query); err != nil {
		return nil, err
	} else {
		subgrupos = make(map[int]*models.DetalleSubgrupo)
		for _, sg := range sg {
			if subgrupos[sg.SubgrupoId.Id] == nil {
				subgrupos[sg.SubgrupoId.Id] = sg
			}
		}
	}

	for _, el := range elsOrg {
		elC := new(models.DetalleElemento__)
		elC.Id = 0
		elC.Nombre = el.Nombre
		elC.Cantidad = el.Cantidad
		elC.Marca = el.Marca
		elC.Serie = el.Serie
		elC.UnidadMedida = el.UnidadMedida
		elC.ValorUnitario = el.ValorUnitario
		elC.Subtotal = el.Subtotal
		elC.Descuento = el.Descuento
		elC.ValorTotal = el.ValorTotal
		elC.PorcentajeIvaId = el.PorcentajeIvaId
		elC.ValorIva = el.ValorIva
		elC.ValorFinal = el.ValorFinal
		elC.SubgrupoCatalogoId = subgrupos[el.SubgrupoCatalogoId]
		elC.EstadoElementoId = el.EstadoElementoId
		elC.ActaRecibidoId = el.ActaRecibidoId
		elC.Placa = el.Placa
		elC.ValorResidual = el.ValorResidual
		elC.VidaUtil = el.VidaUtil

		completos = append(completos, elC)

	}

	return completos, nil

}

func generarNuevosActa(nuevos []*models.DetalleElemento_) (actualizados []*models.Elemento, outputError map[string]interface{}) {

	funcion := "generarNuevosActa"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	if err := formatdata.FillStruct(nuevos, &actualizados); err != nil {
		logs.Error(err)
		eval := " - formatdata.FillStruct(nuevos, &actualizados)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	return actualizados, nil

}

func generarNuevos(nuevos []*models.DetalleElemento) (actualizados []*models.DetalleElemento__, outputError map[string]interface{}) {

	funcion := "generarNuevosActa"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	if err := formatdata.FillStruct(nuevos, &actualizados); err != nil {
		logs.Error(err)
		eval := " - formatdata.FillStruct(nuevos, &actualizados)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	return actualizados, nil

}
