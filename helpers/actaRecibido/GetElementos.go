package actaRecibido

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"

	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// GetElementos Consulta una lista de elementos así como el tipo de bien y el subgrupo
func GetElementos(actaId int, ids []int) (elementosActa []*models.DetalleElemento, outputError map[string]interface{}) {

	funcion := "GetElementos - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	if actaId <= 0 && len(ids) == 0 {
		err := "Se debe indicar un elemento válido"
		logs.Error(err)
		eval := "actaId <= 0 || len(ids) == 0"
		return nil, errorCtrl.Error(funcion+eval, err, "400")
	}

	// Solicita información elementos acta
	var query string
	if actaId > 0 {
		query += "Activo:True,ActaRecibidoId__Id:" + strconv.Itoa(actaId)
	} else {
		query += "Id__in:" + utilsHelper.ArrayToString(ids, "|")
	}

	elementos, err := actaRecibido.GetAllElemento(query, "", "Id", "desc", "", "-1")
	if err != nil {
		return nil, err
	}

	var auxE *models.DetalleElemento
	subgrupos := make(map[int]*models.DetalleSubgrupo)
	tiposBien := make(map[int]*models.TipoBien)

	payload := "fields=Id,SubgrupoId,TipoBienId,Depreciacion,Amortizacion,ValorResidual,VidaUtil&sortby=Id&order=desc" +
		"&query=Activo:true,SubgrupoId__Id:"
	for _, elemento := range elementos {
		auxE = new(models.DetalleElemento)

		if elemento.SubgrupoCatalogoId > 0 {
			if val, ok := subgrupos[elemento.SubgrupoCatalogoId]; !ok || val == nil {
				if detalleSubgrupo_, err := catalogoElementos.GetAllDetalleSubgrupo(payload + strconv.Itoa(elemento.SubgrupoCatalogoId)); err != nil {
					return nil, err
				} else if len(detalleSubgrupo_) == 1 {
					subgrupos[elemento.SubgrupoCatalogoId] = detalleSubgrupo_[0]
				}
			}
		}

		if elemento.TipoBienId > 0 {
			if val, ok := tiposBien[elemento.TipoBienId]; !ok || val == nil {
				var tipoBien_ models.TipoBien
				if err := catalogoElementos.GetTipoBienById(elemento.TipoBienId, &tipoBien_); err != nil {
					return nil, err
				} else if tipoBien_.Id > 0 {
					tiposBien[elemento.TipoBienId] = &tipoBien_
				}
			}
		}

		auxE.Id = elemento.Id
		auxE.Nombre = elemento.Nombre
		auxE.Cantidad = elemento.Cantidad
		auxE.Marca = elemento.Marca
		auxE.Serie = elemento.Serie
		auxE.UnidadMedida = elemento.UnidadMedida
		auxE.ValorUnitario = elemento.ValorUnitario
		auxE.Subtotal = elemento.Subtotal
		auxE.Descuento = elemento.Descuento
		auxE.ValorTotal = elemento.ValorTotal
		auxE.PorcentajeIvaId = elemento.PorcentajeIvaId
		auxE.ValorIva = elemento.ValorIva
		auxE.ValorFinal = elemento.ValorFinal
		auxE.SubgrupoCatalogoId = subgrupos[elemento.SubgrupoCatalogoId]
		auxE.TipoBienId = tiposBien[elemento.TipoBienId]
		auxE.EstadoElementoId = elemento.EstadoElementoId
		auxE.ActaRecibidoId = elemento.ActaRecibidoId
		auxE.Placa = elemento.Placa
		auxE.Activo = elemento.Activo
		auxE.FechaCreacion = elemento.FechaCreacion
		auxE.FechaModificacion = elemento.FechaModificacion

		elementosActa = append(elementosActa, auxE)

	}

	return

}
