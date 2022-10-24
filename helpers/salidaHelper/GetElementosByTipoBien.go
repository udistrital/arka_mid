package salidaHelper

import (
	"time"

	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
)

// GetElementosByTipoBien Consulta la lista de elementos para asociar en una salida determinada agrupando por si son asignables a bodega de consumo.
func GetElementosByTipoBien(entradaId, salidaId int) (elementos_ interface{}, outputError map[string]interface{}) {

	var uvt float64
	if uvt_, err := parametros.GetUVTByVigencia(time.Now().Year()); err != nil {
		return "", err
	} else if uvt_ == 0 {
		return map[string]interface{}{
			"Error": "No se pudo consultar el valor del UVT. Intente más tarde o contacte soporte.",
		}, nil
	} else {
		uvt = uvt_
	}

	bufferTiposBien := make(map[int]models.TipoBien)

	if entradaId > 0 {

		var movimiento models.Movimiento
		var elementos []*models.DetalleElemento
		var consumo = make([]*models.DetalleElemento, 0)
		var devolutivo = make([]*models.DetalleElemento, 0)

		if mov, err := movimientosArka.GetMovimientoById(entradaId); err != nil {
			return nil, err
		} else {
			movimiento = *mov
		}

		var detalle models.FormatoBaseEntrada
		if err := utilsHelper.Unmarshal(movimiento.Detalle, &detalle); err != nil {
			return nil, err
		}

		if detalle.ActaRecibidoId <= 0 {
			return
		}

		if el, err := actaRecibido.GetElementos(detalle.ActaRecibidoId, []int{}); err != nil {
			return nil, err
		} else {
			elementos = el
		}

		for _, el := range elementos {

			if bodega, msg, err := checkBodegaConsumo(el.TipoBienId, el.SubgrupoCatalogoId, int(el.ValorUnitario/uvt), bufferTiposBien); err != nil {
				return nil, err
			} else if msg != "" {
				return map[string]interface{}{
					"Error": msg,
				}, nil
			} else if bodega {
				consumo = append(consumo, el)
			} else {
				devolutivo = append(devolutivo, el)
			}

		}

		return map[string]interface{}{
			"Salida":     nil,
			"Devolutivo": devolutivo,
			"Consumo":    consumo,
		}, nil

	} else if salidaId > 0 {

		if salida, err := GetSalidaById(salidaId); err != nil {
			return nil, err
		} else {
			var elementos []models.DetalleElementoSalida = salida["Elementos"].([]models.DetalleElementoSalida)
			var consumo = make([]models.DetalleElementoSalida, 0)
			var devolutivo = make([]models.DetalleElementoSalida, 0)

			for _, el := range elementos {

				if bodega, msg, err := checkBodegaConsumo(el.TipoBienId, el.SubgrupoCatalogoId, int(el.ValorUnitario/uvt), bufferTiposBien); err != nil {
					return nil, err
				} else if msg != "" {
					return map[string]interface{}{
						"Error": msg,
					}, nil
				} else if bodega {
					consumo = append(consumo, el)
				} else {
					devolutivo = append(devolutivo, el)
				}

			}

			return map[string]interface{}{
				"Salida":     salida["Salida"],
				"Devolutivo": devolutivo,
				"Consumo":    consumo,
			}, nil
		}
	}

	return

}

func checkBodegaConsumo(tipoBienId *models.TipoBien, subgrupo *models.DetalleSubgrupo, valor int, tiposBien map[int]models.TipoBien) (
	bodega bool, msg string, outputError map[string]interface{}) {

	if (tipoBienId == nil || tipoBienId.Id == 0) && (subgrupo != nil && subgrupo.TipoBienId.Id > 0) {
		if tb, err := catalogoElementos.GetTipoBienIdByValor(subgrupo.TipoBienId.Id, valor, tiposBien); err != nil {
			return false, "", err
		} else if tb == 0 {
			return false, "No se pudo determinar el tipo de bien de los elementos. Revise la parametriazación o contacte soporte.", nil
		} else {
			return tiposBien[tb].BodegaConsumo, "", nil
		}

	} else if tipoBienId != nil && tipoBienId.Id > 0 {

		if _, ok := tiposBien[tipoBienId.Id]; !ok {
			tiposBien[tipoBienId.Id] = *tipoBienId
		}

		return tiposBien[tipoBienId.Id].BodegaConsumo, "", nil
	}

	return false, "No se pudo determinar el tipo de bien de los elementos. Revise la parametriazación o contacte soporte.", nil

}
