package entradaHelper

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/administrativa"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"

	tercerosCRUD "github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	tercerosMID "github.com/udistrital/arka_mid/helpers/mid/terceros"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// DetalleEntrada Consulta el detalle de una entrada incluyendo la transaccion contable (si aplica)
func DetalleEntrada(entradaId int) (result map[string]interface{}, outputError map[string]interface{}) {

	funcion := "DetalleEntrada"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		detalle    models.FormatoBaseEntrada
		movimiento models.Movimiento
		query      string
	)

	resultado := make(map[string]interface{})

	query = "query=Id:" + strconv.Itoa(entradaId)
	if mov, err := movimientosArka.GetAllMovimiento(query); err != nil {
		return nil, err
	} else if len(mov) > 0 {
		movimiento = *mov[0]
	}

	if err := utilsHelper.Unmarshal(movimiento.Detalle, &detalle); err != nil {
		return nil, err
	}

	if detalle.ContratoId > 0 {
		if vigencia, err := strconv.Atoi(detalle.VigenciaContrato); err == nil && vigencia > 0 {
			if contrato, err := administrativa.GetContrato(detalle.ContratoId, detalle.VigenciaContrato); err != nil {
				return nil, err
			} else {
				resultado["contrato"] = contrato["contrato"]
			}
		}
	}

	if movimiento.EstadoMovimientoId.Nombre == "Entrada Aprobada" || movimiento.EstadoMovimientoId.Nombre == "Entrada Con Salida" {
		if detalle.ConsecutivoId > 0 {
			if tr, err := movimientosContables.GetTransaccion(detalle.ConsecutivoId, "consecutivo", true); err != nil {
				return nil, err
			} else if len(tr.Movimientos) > 0 {
				if detalleContable, err := asientoContable.GetDetalleContable(tr.Movimientos); err != nil {
					return nil, err
				} else {
					trContable := map[string]interface{}{
						"movimientos": detalleContable,
						"concepto":    tr.Descripcion,
						"fecha":       tr.FechaTransaccion,
					}
					resultado["trContable"] = trContable
				}
			}
		}
	}

	if detalle.ActaRecibidoId > 0 {
		query = "ActaRecibidoId__Id:" + strconv.Itoa(detalle.ActaRecibidoId)
		var acta *models.HistoricoActa
		if tr, err := actaRecibido.GetAllHistoricoActa(query, "", "Id", "desc", "", "1"); err != nil {
			return nil, err
		} else {
			acta = tr[0]
		}

		if acta.ProveedorId > 0 {
			if tercero, err := tercerosCRUD.GetNombreTerceroById(acta.ProveedorId); err != nil {
				return nil, err
			} else {
				resultado["proveedor"] = tercero
			}
		}
	}

	if detalle.Factura > 0 {
		soporte := *new(models.SoporteActa)
		if err := actaRecibido.GetSoporteById(detalle.Factura, &soporte); err != nil {
			return nil, err
		}
		resultado["factura"] = soporte
	}

	if detalle.SupervisorId > 0 {
		supervisor := make([]map[string]interface{}, 0)
		if err := tercerosMID.GetTercerosByTipo("funcionarioPlanta", detalle.SupervisorId, &supervisor); err != nil {
			return nil, err
		} else if len(supervisor) > 0 {
			resultado["supervisor"] = supervisor[0]
		}
	}

	if detalle.OrdenadorGastoId > 0 {
		ordenadores := make([]map[string]interface{}, 0)
		if err := tercerosMID.GetTercerosByTipo("ordenadoresGasto", detalle.OrdenadorGastoId, &ordenadores); err != nil {
			return nil, err
		} else if len(ordenadores) > 0 {
			resultado["ordenador"] = ordenadores[0]
		}
	}

	resultado["movimiento"] = movimiento

	return resultado, nil
}
