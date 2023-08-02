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
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	timebogota "github.com/udistrital/arka_mid/utils_oas/timeBogota"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const errNoElementos = "No se encontraron elementos para asociar a la entrada."

// AprobarEntrada Actualiza una entrada a estado aprobada, calcula la transacción contable y genera las novedades correspondientes
func AprobarEntrada(entradaId int, resultado_ *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("AprobarEntrada - Unhandled Error!", "500")

	formato, outputError := getFormato(entradaId, resultado_)
	if outputError != nil || resultado_.Error != "" {
		return
	}

	terceroId, outputError := getTerceroEntrada(formato, resultado_)
	if outputError != nil || resultado_.Error != "" {
		return
	}

	elementos, novedades, outputError := getElementosEntrada(formato, entradaId, resultado_)
	if outputError != nil || resultado_.Error != "" {
		return
	}

	outputError = contabilidadEntrada(resultado_, formato, elementos, terceroId)
	if outputError != nil || resultado_.Error != "" {
		return
	}

	for _, nov := range novedades {
		outputError = movimientosArka.PostNovedadElemento(&nov)
		if outputError != nil {
			return
		}
	}

	resultado_.Movimiento.FechaCorte = utilsHelper.Time(timebogota.TiempoBogota())
	_, outputError = movimientosArka.PutMovimiento(&resultado_.Movimiento, resultado_.Movimiento.Id)
	return
}

func getFormato(entradaId int, resultado *models.ResultadoMovimiento) (formato models.FormatoBaseEntrada, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("getFormato - Unhandled Error!", "500")

	movimiento, outputError := movimientosArka.GetAllMovimiento("query=Id:" + strconv.Itoa(entradaId))
	if outputError != nil || len(movimiento) != 1 {
		return
	}

	resultado.Movimiento = *movimiento[0]
	if resultado.Movimiento.ConsecutivoId == nil || *resultado.Movimiento.ConsecutivoId == 0 {
		resultado.Error = "No se puede continuar con el cálculo de la transaccón contable. Contacte soporte."
		return
	}

	outputError = utilsHelper.Unmarshal(resultado.Movimiento.Detalle, &formato)
	if outputError != nil {
		return
	}

	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&resultado.Movimiento.EstadoMovimientoId.Id, "Entrada Aprobada")
	return
}

func getTerceroEntrada(detalle models.FormatoBaseEntrada, resutado *models.ResultadoMovimiento) (terceroId int, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("getTerceroEntrada - Unhandled Error!", "500")

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

func getElementosEntrada(detalle models.FormatoBaseEntrada, movimientoId int, resultado *models.ResultadoMovimiento) (elementos []*models.Elemento, novedades []models.NovedadElemento, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("getElementosEntrada - Unhandled Error!", "500")

	if detalle.ActaRecibidoId == 0 && len(detalle.Elementos) == 0 {
		resultado.Error = errNoElementos
		return
	}

	if detalle.ActaRecibidoId > 0 {
		query := "Activo:true,ActaRecibidoId__Id:" + strconv.Itoa(detalle.ActaRecibidoId)
		elementos, outputError = actaRecibido.GetAllElemento(query, "ValorUnitario,ValorTotal,SubgrupoCatalogoId,TipoBienId", "SubgrupoCatalogoId", "desc", "", "-1")
		if len(elementos) == 0 {
			resultado.Error = errNoElementos
		}
	} else if len(detalle.Elementos) > 0 {
		for _, el := range detalle.Elementos {
			if el.VidaUtil == nil || el.ValorLibros == nil || el.ValorResidual == nil {
				resultado.Error = "No se indicó correctamente el nuevo valor de los elementos. Rechace la entrada y haga la respectiva edición."
				return
			}

			var novedad = models.NovedadElemento{
				VidaUtil:             *el.VidaUtil,
				ValorLibros:          *el.ValorLibros,
				ValorResidual:        *el.ValorResidual * *el.ValorLibros,
				MovimientoId:         &models.Movimiento{Id: movimientoId},
				ElementoMovimientoId: &models.ElementosMovimiento{Id: el.Id},
				Activo:               true,
			}
			novedades = append(novedades, novedad)

			if *el.ValorLibros > 0 {
				var elementoMovimiento models.ElementosMovimiento
				outputError = movimientosArka.GetElementosMovimientoById(el.Id, &elementoMovimiento)
				if outputError != nil {
					return
				}

				var elementoActa models.Elemento
				outputError = actaRecibido.GetElementoById(*elementoMovimiento.ElementoActaId, &elementoActa)
				if outputError != nil {
					return
				}

				elementoActa.ValorUnitario = *el.ValorLibros
				elementoActa.ValorTotal = *el.ValorLibros
				elementos = append(elementos, &elementoActa)
			}
		}
	}

	return
}

func contabilidadEntrada(resultado_ *models.ResultadoMovimiento, formatoEntrada models.FormatoBaseEntrada, elementos []*models.Elemento, terceroId int) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("contabilidadEntrada - Unhandled Error!", "500")

	if len(elementos) == 0 {
		return
	}

	detalleContable, outputError := descripcionMovimientoContable(resultado_.Movimiento.Detalle)
	if outputError != nil {
		return
	}

	var transaccion = models.TransaccionMovimientos{ConsecutivoId: *resultado_.Movimiento.ConsecutivoId}
	bufferCuentas := make(map[string]models.CuentaContable)
	resultado_.Error, outputError = asientoContable.CalcularMovimientosContables(elementos, detalleContable, 0, resultado_.Movimiento.FormatoTipoMovimientoId.Id, terceroId, terceroId, bufferCuentas, nil, &transaccion.Movimientos)
	if outputError != nil || resultado_.Error != "" {
		return
	}

	resultado_.Error, outputError = asientoContable.CreateTransaccionContable(getTipoComprobanteEntradas(), "Entrada Almacén", &transaccion)
	if outputError != nil || resultado_.Error != "" {
		return
	}

	resultado_.TransaccionContable.Concepto = transaccion.Descripcion
	resultado_.TransaccionContable.Fecha = transaccion.FechaTransaccion
	resultado_.TransaccionContable.Movimientos, outputError = asientoContable.GetDetalleContable(transaccion.Movimientos, bufferCuentas)
	_, outputError = movimientosContables.PostTrContable(&transaccion)
	return
}

// descripcionMovimientoContable Genera la descipción de cada uno de los movimientos contables asociados a una entrada.
func descripcionMovimientoContable(detalle string) (detalle_ string, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("descripcionMovimientoContable - Unhandled Error!", "500")

	var mapDetalle map[string]interface{}
	outputError = utilsHelper.Unmarshal(detalle, &mapDetalle)
	if outputError != nil {
		return
	}

	for k, v := range mapDetalle {
		if k == "factura" {
			var sop models.SoporteActa
			outputError = actaRecibido.GetSoporteById(int(v.(float64)), &sop)
			if outputError != nil {
				return
			}

			detalle_ += "Factura: " + sop.Consecutivo + ", "
		} else if k != "consecutivo" && k != "ConsecutivoId" && k != "elementos" {
			k = strings.TrimSuffix(k, "_id")
			k = strings.ReplaceAll(k, "_", " ")
			caser := cases.Title(language.Spanish)
			k = caser.String(k)
			detalle_ += k + ": " + fmt.Sprintf("%v", v) + ", "
		}
	}

	detalle_ = strings.TrimSuffix(detalle_, ", ")
	return
}
