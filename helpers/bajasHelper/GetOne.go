package bajasHelper

import (
	"net/url"
	"strconv"

	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/oikos"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/arka_mid/helpers/inventarioHelper"
	"github.com/udistrital/arka_mid/helpers/mid/terceros"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// GetOne Consulta el detalle de la baja: elementos, revisor, solicitante, soporte, tipo
func GetOne(id int, Baja *models.TrBaja) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GetOne - Unhandled Error!", "500")

	var (
		movimiento  *models.Movimiento
		detalle     models.FormatoBaja
		dependencia models.Parametro
	)

	// Se consulta el movimiento
	query := "query=Id:" + strconv.Itoa(id)
	if movimientoA, err := movimientosArka.GetAllMovimiento(query); err != nil || len(movimientoA) != 1 {
		return err
	} else {
		movimiento = movimientoA[0]
	}

	if err := utilsHelper.Unmarshal(movimiento.Detalle, &detalle); err != nil {
		return err
	}

	// Se consulta el detalle del funcionario solicitante
	if detalle.Funcionario > 0 {
		Baja.Funcionario, outputError = terceros.GetInfoTerceroById(detalle.Funcionario)
		if outputError != nil {
			return
		}
	}

	// Se consulta el detalle del revisor si lo hay
	if detalle.Revisor > 0 {
		Baja.Revisor, outputError = terceros.GetInfoTerceroById(detalle.Revisor)
		if outputError != nil {
			return
		}
	}

	// Se consulta el detalle de los elementos relacionados en la solicitud
	if len(detalle.Elementos) > 0 {
		Baja.Elementos, outputError = getDetalleElementos(detalle.Elementos)
		if outputError != nil {
			return
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
		if err := parametros.GetParametroById(detalle.DependenciaId, &dependencia); err != nil {
			return err
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
	Baja.TrContable = &models.InfoTransaccionContable{}

	if movimiento.EstadoMovimientoId.Nombre != "Baja Aprobada" || movimiento.ConsecutivoId == nil || *movimiento.ConsecutivoId <= 0 {
		return
	}

	*Baja.TrContable, outputError = asientoContable.GetFullDetalleContable(*movimiento.ConsecutivoId)

	return
}

// getDetalleElementos consulta el historial de una serie de elementos dados los ids en el api movimientos_arka_crud
func getDetalleElementos(ids []int) (Elementos []*models.DetalleElementoBaja, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("getDetalleElementos - Unhandled Error!", "500")

	var (
		elementosActa       []*models.DetalleElemento
		elementosMovimiento []*models.ElementosMovimiento
	)
	Elementos = make([]*models.DetalleElementoBaja, 0)

	// Consulta asignación de los elementos
	query := "sortby=ElementoActaId&order=desc&limit=-1&query=Id__in:"
	query += url.QueryEscape(utilsHelper.ArrayToString(ids, "|"))
	if elementoMovimiento_, err := movimientosArka.GetAllElementosMovimiento(query); err != nil {
		return nil, err
	} else {
		elementosMovimiento = elementoMovimiento_
	}

	ids = []int{}
	for _, el := range elementosMovimiento {
		ids = append(ids, el.ElementoActaId)
	}

	// Consulta de Marca, Nombre, Serie y Subgrupo se hace mediante el actaRecibidoHelper
	if elemento_, err := actaRecibido.GetElementos(0, ids); err != nil {
		return nil, err
	} else {
		elementosActa = elemento_
	}

	if len(elementosActa) == len(elementosMovimiento) {

		for i := 0; i < len(elementosActa); i++ {

			elemento := new(models.DetalleElementoBaja)
			elemento.Historial, outputError = movimientosArka.GetHistorialElemento(elementosMovimiento[i].Id, true)
			if outputError != nil {
				return
			}

			funcionario, ubicacion, err := inventarioHelper.GetEncargado(elemento.Historial)
			if err != nil {
				return nil, err
			}

			if ubicacion > 0 {
				elemento.Ubicacion, outputError = oikos.GetSedeDependenciaUbicacion(ubicacion)
				if outputError != nil {
					return
				}
			}

			if funcionario > 0 {
				elemento.Funcionario, outputError = terceros.GetInfoTerceroById(funcionario)
				if outputError != nil {
					return
				}
			}

			elemento.Id = elementosMovimiento[i].Id
			elemento.Placa = elementosActa[i].Placa
			elemento.Nombre = elementosActa[i].Nombre
			elemento.Marca = elementosActa[i].Marca
			elemento.Serie = elementosActa[i].Serie
			elemento.SubgrupoCatalogoId = elementosActa[i].SubgrupoCatalogoId

			Elementos = append(Elementos, elemento)
		}
	}

	return Elementos, nil
}
