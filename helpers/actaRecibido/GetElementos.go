package actaRecibido

import (
	"errors"
	"strconv"

	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
)

// GetElementos ...
func GetElementos(actaId int, ids []int) (elementosActa []*models.DetalleElemento, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetElementos - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var (
		urlcrud string
		auxE    *models.DetalleElemento
	)

	subgrupos := make(map[int]*models.DetalleSubgrupo)
	tiposBien := make(map[int]*models.TipoBien)

	if actaId > 0 || len(ids) > 0 { // (1) error parametro
		// Solicita información elementos acta

		var query string
		if actaId > 0 {
			query += "Activo:True,ActaRecibidoId__Id:" + strconv.Itoa(actaId)
		} else {
			query += "Id__in:" + utilsHelper.ArrayToString(ids, "|")
		}

		if elementos, err := actaRecibido.GetAllElemento(query, "", "Id", "desc", "", "-1"); err != nil {
			return nil, err
		} else {

			if len(elementos) == 0 || elementos[0].Id == 0 {
				return nil, nil
			}

			for _, elemento := range elementos {
				auxE = new(models.DetalleElemento)
				if elemento.SubgrupoCatalogoId > 0 {
					if val, ok := subgrupos[elemento.SubgrupoCatalogoId]; !ok || val == nil {
						urlcrud = "query=Activo:true,SubgrupoId__Id:" + strconv.Itoa(elemento.SubgrupoCatalogoId)
						urlcrud += "&fields=Id,SubgrupoId,TipoBienId,Depreciacion,Amortizacion,ValorResidual,VidaUtil&sortby=Id&order=desc"
						if detalleSubgrupo_, err := catalogoElementos.GetAllDetalleSubgrupo(urlcrud); err != nil {
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

			return elementosActa, nil
		}
	} else {
		err := errors.New("ID must be greater than 0")
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetElementos - actaId > 0",
			"err":     err,
			"status":  "400",
		}
		return nil, outputError
	}
}
