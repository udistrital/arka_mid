package entradaHelper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	actasCrud "github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// AprobarEntrada Actualiza una entrada a estado aprobada y hace los respectivos registros en kronos y transacciones contables
func AprobarEntrada(entradaId int) (result map[string]interface{}, outputError map[string]interface{}) {

	funcion := "AprobarEntrada"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		historico          *models.HistoricoActa
		elementos          []*models.Elemento
		movimiento         *models.Movimiento
		detalleMovimiento  map[string]interface{}
		estadoMovimientoId int
		consecutivoId      int
		detalleContable    string
	)

	resultado := make(map[string]interface{})

	if mov, err := movimientosArka.GetMovimientoById(entradaId); err != nil {
		return nil, err
	} else {
		movimiento = mov
	}

	if err := utilsHelper.Unmarshal(movimiento.Detalle, &detalleMovimiento); err != nil {
		return nil, err
	}

	if err := movimientosArka.GetEstadoMovimientoIdByNombre(&estadoMovimientoId, "Entrada Aprobada"); err != nil {
		return nil, err
	}

	query := "Activo:true,ActaRecibidoId__Id:" + fmt.Sprint(detalleMovimiento["acta_recibido_id"])
	if ha, err := actasCrud.GetAllHistoricoActa(query, "", "FechaCreacion", "desc", "", "-1"); err != nil {
		return nil, err
	} else {
		historico = ha[0]
	}

	if el_, err := actasCrud.GetAllElemento(query, "", "FechaCreacion", "desc", "", "-1"); err != nil {
		return nil, err
	} else {
		elementos = el_
	}

	var groups = make(map[int]float64)
	for _, elemento := range elementos {
		x := float64(0)
		if val, ok := groups[elemento.SubgrupoCatalogoId]; ok {
			x = val + elemento.ValorTotal
		} else {
			x = elemento.ValorTotal
		}
		groups[elemento.SubgrupoCatalogoId] = x
	}

	if val, ok := detalleMovimiento["ConsecutivoId"]; ok && val != nil {
		consecutivoId = int(val.(float64))
	}

	if err := descripcionMovimientoContable(detalleMovimiento, &detalleContable); err != nil {
		return nil, err
	}

	var trContable map[string]interface{}
	if len(groups) > 0 && historico.ProveedorId > 0 && consecutivoId > 0 {
		if tr_, err := asientoContable.AsientoContable(groups, getTipoComprobanteEntradas(), strconv.Itoa(movimiento.FormatoTipoMovimientoId.Id), detalleContable, "Entrada de almacen", historico.ProveedorId, consecutivoId, true); tr_ == nil || err != nil {
			return nil, err
		} else {
			trContable = tr_
			if tr_["errorTransaccion"].(string) != "" {
				return tr_, nil
			}
		}
	}

	movimiento.EstadoMovimientoId = &models.EstadoMovimiento{Id: estadoMovimientoId}
	if movimiento_, err := movimientosArka.PutMovimiento(movimiento, movimiento.Id); err != nil {
		return nil, err
	} else {
		movimiento = movimiento_
	}

	resultado["movimientoArka"] = movimiento
	resultado["transaccionContable"] = trContable["resultadoTransaccion"]
	resultado["tercero"] = trContable["tercero"]
	resultado["errorTransaccion"] = ""

	return resultado, nil
}

// descripcionMovimientoContable Genera la descipción de cada uno de los movimientos contables asociados a una entrada.
func descripcionMovimientoContable(detalle map[string]interface{}, detalle_ *string) (outputError map[string]interface{}) {

	funcion := "descripcionMovimientoContable"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	for k, v := range detalle {
		if k == "factura" {
			var sop models.SoporteActa

			if err := actasCrud.GetSoporteById(int(v.(float64)), &sop); err != nil {
				return err
			}

			*detalle_ += "Factura: " + sop.Consecutivo + ", "
		} else if k != "consecutivo" && k != "ConsecutivoId" {
			k = strings.TrimSuffix(k, "_id")
			k = strings.ReplaceAll(k, "_", " ")
			k = strings.Title(k)
			*detalle_ += k + ": " + fmt.Sprintf("%v", v) + ", "
		}
	}

	*detalle_ = strings.TrimSuffix(*detalle_, ", ")

	return
}
