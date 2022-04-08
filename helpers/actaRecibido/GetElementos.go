package actaRecibido

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/astaxie/beego/logs"

	crud_actas "github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	// "github.com/udistrital/utils_oas/formatdata"
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

	subgrupos := make(map[int]interface{})
	consultasSubgrupos := 0
	evSubgrupos := 0

	if actaId > 0 || len(ids) > 0 { // (1) error parametro
		// Solicita información elementos acta

		var query string
		if actaId > 0 {
			query += "Activo:True,ActaRecibidoId__Id:" + strconv.Itoa(actaId)
		} else {
			query += "Id__in:" + utilsHelper.ArrayToString(ids, "|")
		}

		if elementos, err := crud_actas.GetAllElemento(query, "", "Id", "desc", "", "-1"); err != nil {
			return nil, err
		} else {

			if len(elementos) == 0 || elementos[0].Id == 0 {
				return nil, nil
			}

			for _, elemento := range elementos {

				var subgrupoId *models.Subgrupo
				subgrupoId = new(models.Subgrupo)
				var tipoBienId *models.TipoBien
				tipoBienId = new(models.TipoBien)
				auxE = new(models.DetalleElemento)
				subgrupo := *&models.DetalleSubgrupo{
					SubgrupoId: subgrupoId,
					TipoBienId: tipoBienId,
				}

				subgrupo.TipoBienId = tipoBienId
				subgrupo.SubgrupoId = subgrupoId

				idSubgrupo := elemento.SubgrupoCatalogoId
				reqSubgrupo := func() (interface{}, map[string]interface{}) {
					urlcrud = "query=Activo:true,SubgrupoId__Id:" + strconv.Itoa(idSubgrupo)
					urlcrud += "&fields=SubgrupoId,TipoBienId,Depreciacion,Amortizacion,ValorResidual,VidaUtil&sortby=Id&order=desc"
					if detalleSubgrupo_, err := catalogoElementos.GetAllDetalleSubgrupo(urlcrud); err == nil && len(detalleSubgrupo_) > 0 {
						return detalleSubgrupo_[0], nil
					} else if err != nil {
						return nil, err
					} else {
						logs.Error(err)
						return nil, map[string]interface{}{
							"funcion": "GetElementos - catalogoElementosHelper.GetDetalleSubgrupo(idSubgrupo)",
							"err":     err,
							"status":  "500",
						}
					}
				}

				if idSubgrupo > 0 {
					if v, err := utilsHelper.BufferGeneric(idSubgrupo, subgrupos, reqSubgrupo, &consultasSubgrupos, &evSubgrupos); err == nil {
						if v != nil {
							if jsonString, err := json.Marshal(v); err == nil {
								if err := json.Unmarshal(jsonString, &subgrupo); err != nil {
									logs.Error(err)
									outputError = map[string]interface{}{
										"funcion": "GetElementos - json.Unmarshal(jsonString, &subgrupo)",
										"err":     err,
										"status":  "500",
									}
									return nil, outputError
								}
							}
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
				auxE.SubgrupoCatalogoId = &subgrupo
				auxE.EstadoElementoId = elemento.EstadoElementoId
				auxE.ActaRecibidoId = elemento.ActaRecibidoId
				auxE.Placa = elemento.Placa
				auxE.Activo = elemento.Activo
				auxE.FechaCreacion = elemento.FechaCreacion
				auxE.FechaModificacion = elemento.FechaModificacion

				elementosActa = append(elementosActa, auxE)

			}

			logs.Info("consultasSubgrupos:", consultasSubgrupos, " - Evitadas: ", evSubgrupos)
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
