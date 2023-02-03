package entradaHelper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// AprobarEntrada Actualiza una entrada a estado aprobada, calcula la transacción contable y genera las novedades correspondientes
func AprobarEntrada(entradaId int, resultado_ *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("AprobarEntrada - Unhandled Error!", "500")

	var (
		detalleMovimiento models.FormatoBaseEntrada
		detalleContable   string
		transaccion       models.TransaccionMovimientos
	)

	query := "query=Id:" + strconv.Itoa(entradaId)
	movimiento, outputError := movimientosArka.GetAllMovimiento(query)
	if outputError != nil || len(movimiento) != 1 || movimiento[0].EstadoMovimientoId.Nombre != "Entrada En Trámite" {
		return
	}

	resultado_.Movimiento = *movimiento[0]
	outputError = utilsHelper.Unmarshal(resultado_.Movimiento.Detalle, &detalleMovimiento)
	if outputError != nil {
		return
	}

	if detalleMovimiento.ConsecutivoId == 0 {
		resultado_.Error = "No se puede continuar con el cálculo de la transaccón contable. Contacte soporte."
		return
	}

	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&resultado_.Movimiento.EstadoMovimientoId.Id, "Entrada Aprobada")
	if outputError != nil {
		return
	}

	terceroId, outputError := getTerceroEntrada(detalleMovimiento, resultado_)
	if outputError != nil || resultado_.Error != "" {
		return
	}

	elementos, novedades, outputError := getElementosEntrada(detalleMovimiento, entradaId)
	if outputError != nil || len(elementos) == 0 {
		resultado_.Error = "No se encontraron elementos para asociar a la entrada."
		return
	}

	outputError = descripcionMovimientoContable(resultado_.Movimiento.Detalle, &detalleContable)
	if outputError != nil {
		return
	}

	bufferCuentas := make(map[string]models.CuentaContable)
	msg, outputError := asientoContable.CalcularMovimientosContables(elementos, detalleContable, 0, resultado_.Movimiento.FormatoTipoMovimientoId.Id, terceroId, terceroId, bufferCuentas, nil, &transaccion.Movimientos)
	if outputError != nil || msg != "" {
		resultado_.Error = msg
		return
	}

	msg, outputError = asientoContable.CreateTransaccionContable(getTipoComprobanteEntradas(), "Entrada Almacén", &transaccion)
	if outputError != nil || msg != "" {
		resultado_.Error = msg
		return
	}

	transaccion.ConsecutivoId = detalleMovimiento.ConsecutivoId
	_, outputError = movimientosContables.PostTrContable(&transaccion)
	if outputError != nil {
		return
	}

	for _, nov := range novedades {
		outputError = movimientosArka.PostNovedadElemento(&nov)
		if outputError != nil {
			return
		}
	}

	_, outputError = movimientosArka.PutMovimiento(&resultado_.Movimiento, resultado_.Movimiento.Id)
	if outputError != nil {
		return outputError
	}

	resultado_.TransaccionContable.Concepto = transaccion.Descripcion
	resultado_.TransaccionContable.Fecha = transaccion.FechaTransaccion
	resultado_.TransaccionContable.Movimientos, outputError = asientoContable.GetDetalleContable(transaccion.Movimientos, bufferCuentas)

	return
}

func getTerceroEntrada(detalle models.FormatoBaseEntrada, resutado *models.ResultadoMovimiento) (terceroId int, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("getTerceroEntrada - Unhandled Error!", "500")

	var historico []models.HistoricoActa
	query := "Activo:true,ActaRecibidoId__Id:" + strconv.Itoa(detalle.ActaRecibidoId)
	if detalle.ActaRecibidoId > 0 {
		historico, outputError = actaRecibido.GetAllHistoricoActa(query, "", "FechaCreacion", "desc", "", "1")
		if outputError != nil || len(historico) != 1 {
			if len(historico) != 1 {
				resutado.Error = "No se pudo consultar la información del acta. Contacte soporte."
			}
			return
		}
	}

	if detalle.ActaRecibidoId > 0 && historico[0].ActaRecibidoId.TipoActaId.CodigoAbreviacion == "REG" {
		terceroId = historico[0].ProveedorId
	} else {
		terceroId, outputError = terceros.GetTerceroUD()
	}

	if terceroId == 0 {
		resutado.Error = "No se pudo consultar el tercero para asociar a la transacción contable. Contacte soporte."
	}

	return
}

func getElementosEntrada(detalle models.FormatoBaseEntrada, movimientoId int) (elementos []*models.Elemento, novedades []models.NovedadElemento, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("getElementosEntrada - Unhandled Error!", "500")

	query := "Activo:true,ActaRecibidoId__Id:" + strconv.Itoa(detalle.ActaRecibidoId)
	if detalle.ActaRecibidoId > 0 {
		elementos, outputError = actaRecibido.GetAllElemento(query, "ValorUnitario,ValorTotal,SubgrupoCatalogoId,TipoBienId", "SubgrupoCatalogoId", "desc", "", "-1")
	} else if len(detalle.Elementos) > 0 {
		for _, el := range detalle.Elementos {
			query = "Id:" + strconv.Itoa(el.Id)
			elementoMovimiento, err := movimientosArka.GetAllElementosMovimiento(query)
			if err != nil {
				outputError = err
				return
			}

			var elementoActa models.Elemento
			outputError = actaRecibido.GetElementoById(elementoMovimiento[0].ElementoActaId, &elementoActa)
			if outputError != nil {
				return
			}

			if el.VidaUtil != nil && el.ValorLibros != nil && el.ValorResidual != nil {
				var novedad = models.NovedadElemento{
					Id:                   0,
					VidaUtil:             *el.VidaUtil,
					ValorLibros:          *el.ValorLibros,
					ValorResidual:        *el.ValorResidual,
					Metadata:             "",
					MovimientoId:         &models.Movimiento{Id: movimientoId},
					ElementoMovimientoId: &models.ElementosMovimiento{Id: el.Id},
					Activo:               true,
				}
				novedades = append(novedades, novedad)
			}

			elementos = append(elementos, &elementoActa)
		}
	}

	return
}

// descripcionMovimientoContable Genera la descipción de cada uno de los movimientos contables asociados a una entrada.
func descripcionMovimientoContable(detalle string, detalle_ *string) (outputError map[string]interface{}) {

	funcion := "descripcionMovimientoContable"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var mapDetalle map[string]interface{}
	outputError = utilsHelper.Unmarshal(detalle, &mapDetalle)
	if outputError != nil {
		return
	}

	for k, v := range mapDetalle {
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
