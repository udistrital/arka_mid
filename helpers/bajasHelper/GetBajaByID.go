package bajasHelper

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/oikos"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	midTerceros "github.com/udistrital/arka_mid/helpers/mid/terceros"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// GetBajaByID Consulta el detalle de la baja: elementos, revisor, solicitante, soporte, tipo
func GetBajaByID(id int, Baja *models.TrBaja) (outputError map[string]interface{}) {

	funcion := "GetBajaByID"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		movimiento  *models.Movimiento
		detalle     models.FormatoBaja
		dependencia models.Dependencia
	)

	// Se consulta el movimiento
	query := "query=Id:" + strconv.Itoa(id)
	if movimientoA, err := movimientosArka.GetAllMovimiento(query); err != nil {
		return err
	} else if len(movimientoA) > 0 {
		movimiento = movimientoA[0]
	}

	if err := utilsHelper.Unmarshal(movimiento.Detalle, &detalle); err != nil {
		return err
	}

	// Se consulta el detalle del funcionario solicitante
	if detalle.Funcionario > 0 {
		if funcionario, err := midTerceros.GetInfoTerceroById(detalle.Funcionario); err != nil {
			return err
		} else {
			Baja.Funcionario = funcionario
		}
	}

	// Se consulta el detalle del revisor si lo hay
	if detalle.Revisor > 0 {
		if revisor, err := midTerceros.GetInfoTerceroById(detalle.Revisor); err != nil {
			return err
		} else {
			Baja.Revisor = revisor
		}
	}

	// Se consulta el detalle de los elementos relacionados en la solicitud
	if len(detalle.Elementos) > 0 {
		if elementos, err := GetDetalleElementos(detalle.Elementos); err != nil {
			return err
		} else {
			Baja.Elementos = elementos
		}
	}

	// Se consulta el detalle de los elementos relacionados en la solicitud
	query = "query=MovimientoId__Id:" + strconv.Itoa(id)
	if soportes, err := movimientosArka.GetAllSoporteMovimiento(query); err != nil {
		return err
	} else if len(soportes) > 0 {
		Baja.Soporte = soportes[0].DocumentoId
	}

	if detalle.DependenciaId > 0 {
		if err := oikos.GetDependenciaById(detalle.DependenciaId, &dependencia); err != nil {
			return err
		}
	}

	if movimiento.EstadoMovimientoId.Nombre == "Baja Aprobada" {
		if detalle.ConsecutivoId > 0 {
			if tr, err := movimientosContables.GetTransaccion(detalle.ConsecutivoId, "consecutivo", true); err != nil {
				return err
			} else if len(tr.Movimientos) > 0 {
				if detalleContable, err := asientoContable.GetDetalleContable(tr.Movimientos, nil); err != nil {
					return err
				} else {
					trContable := models.InfoTransaccionContable{
						Movimientos: detalleContable,
						Concepto:    tr.Descripcion,
						Fecha:       tr.FechaTransaccion,
					}
					Baja.TrContable = &trContable
				}
			}
		}
	}

	Baja.Id = movimiento.Id
	Baja.TipoBaja = movimiento.FormatoTipoMovimientoId
	Baja.Movimiento = movimiento
	Baja.Observaciones = movimiento.Observacion
	Baja.RazonRechazo = detalle.RazonRechazo
	Baja.Resolucion = detalle.Resolucion
	Baja.FechaRevisionC = detalle.FechaRevisionC
	Baja.DependenciaId = dependencia.Nombre

	return

}
