package entradaHelper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// AprobarEntrada Actualiza una entrada a estado aprobada y hace los respectivos registros en kronos y transacciones contables
func AprobarEntrada(entradaId int, resultado_ *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	funcion := "AprobarEntrada - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	var (
		historico         models.HistoricoActa
		elementos         []*models.Elemento
		detalleMovimiento map[string]interface{}
		detalleContable   string
		transaccion       models.TransaccionMovimientos
		bufferCuentas     map[string]models.CuentaContable
	)

	query := "query=Id:" + strconv.Itoa(entradaId)
	if mov, err := movimientosArka.GetAllMovimiento(query); err != nil {
		return err
	} else if len(mov) == 1 && mov[0].EstadoMovimientoId.Nombre == "Entrada En Trámite" {
		resultado_.Movimiento = *mov[0]
	} else {
		return
	}

	if err := movimientosArka.GetEstadoMovimientoIdByNombre(&resultado_.Movimiento.EstadoMovimientoId.Id, "Entrada Aprobada"); err != nil {
		return err
	}

	if err := utilsHelper.Unmarshal(resultado_.Movimiento.Detalle, &detalleMovimiento); err != nil {
		return err
	}

	query = "Activo:true,ActaRecibidoId__Id:" + fmt.Sprint(detalleMovimiento["acta_recibido_id"])
	if ha, err := actaRecibido.GetAllHistoricoActa(query, "", "FechaCreacion", "desc", "", "-1"); err != nil {
		return err
	} else {
		historico = *ha[0]
	}

	if el_, err := actaRecibido.GetAllElemento(query, "ValorUnitario,ValorTotal,SubgrupoCatalogoId,TipoBienId", "SubgrupoCatalogoId", "desc", "", "-1"); err != nil {
		return err
	} else {
		elementos = el_
	}

	if val, ok := detalleMovimiento["ConsecutivoId"]; ok && val != nil && val.(float64) > 0 {
		transaccion.ConsecutivoId = int(val.(float64))
	} else {
		resultado_.Error = "No se puede continuar con el cálculo de la transaccón contable. Contacte soporte."
		return
	}

	if err := descripcionMovimientoContable(detalleMovimiento, &detalleContable); err != nil {
		return err
	}

	bufferCuentas = make(map[string]models.CuentaContable)
	if msg, err := asientoContable.CalcularMovimientosContables(elementos, detalleContable, resultado_.Movimiento.FormatoTipoMovimientoId.Id, historico.ProveedorId, historico.ProveedorId, bufferCuentas,
		nil, &transaccion.Movimientos); err != nil || msg != "" {
		resultado_.Error = msg
		return err
	}

	if msg, err := asientoContable.CreateTransaccionContable(getTipoComprobanteEntradas(), "Entrada Almacén", &transaccion); err != nil || msg != "" {
		resultado_.Error = msg
		return err
	}

	if _, err := movimientosContables.PostTrContable(&transaccion); err != nil {
		return err
	}

	if detalleContable, err := asientoContable.GetDetalleContable(transaccion.Movimientos, bufferCuentas); err != nil {
		return err
	} else {
		resultado_.TransaccionContable.Movimientos = detalleContable
		resultado_.TransaccionContable.Concepto = transaccion.Descripcion
		resultado_.TransaccionContable.Fecha = transaccion.FechaTransaccion
	}

	if movimiento_, err := movimientosArka.PutMovimiento(&resultado_.Movimiento, resultado_.Movimiento.Id); err != nil {
		return err
	} else {
		resultado_.Movimiento = *movimiento_
	}

	return
}

// descripcionMovimientoContable Genera la descipción de cada uno de los movimientos contables asociados a una entrada.
func descripcionMovimientoContable(detalle map[string]interface{}, detalle_ *string) (outputError map[string]interface{}) {

	funcion := "descripcionMovimientoContable"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	for k, v := range detalle {
		if k == "factura" {
			var sop models.SoporteActa

			if err := actaRecibido.GetSoporteById(int(v.(float64)), &sop); err != nil {
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
