package entradaHelper

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/inventarioHelper"
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

	logs.Info("==== INICIO entradaHelper.AprobarEntrada ====")
	logs.Info("AprobarEntrada -> entradaId=%d", entradaId)

	if resultado_ == nil {
		outputError = map[string]interface{}{
			"funcion": "AprobarEntrada - resultado_",
			"err":     "resultado_ nil",
			"status":  "500",
		}
		logs.Error("AprobarEntrada -> resultado_ nil")
		return
	}

	formato, outputError := getFormato(entradaId, resultado_)
	logs.Info("AprobarEntrada -> getFormato outputError=%v resultado.Error=%q formato=%+v", outputError, resultado_.Error, formato)
	if outputError != nil || resultado_.Error != "" {
		logs.Error("AprobarEntrada -> aborta en getFormato")
		return
	}

	terceroId, outputError := getTerceroEntrada(formato, resultado_)
	logs.Info("AprobarEntrada -> getTerceroEntrada outputError=%v resultado.Error=%q terceroId=%d", outputError, resultado_.Error, terceroId)
	if outputError != nil || resultado_.Error != "" {
		logs.Error("AprobarEntrada -> aborta en getTerceroEntrada")
		return
	}

	elementos, novedades, outputError := getElementosEntrada(formato, entradaId, resultado_)
	logs.Info("AprobarEntrada -> getElementosEntrada outputError=%v resultado.Error=%q len(elementos)=%d len(novedades)=%d",
		outputError, resultado_.Error, len(elementos), len(novedades))
	if outputError != nil || resultado_.Error != "" {
		logs.Error("AprobarEntrada -> aborta en getElementosEntrada")
		return
	}

	outputError = contabilidadEntrada(resultado_, formato, elementos, terceroId)
	logs.Info("AprobarEntrada -> contabilidadEntrada outputError=%v resultado.Error=%q transaccionContable=%+v",
		outputError, resultado_.Error, resultado_.TransaccionContable)
	if outputError != nil || resultado_.Error != "" {
		logs.Error("AprobarEntrada -> aborta en contabilidadEntrada")
		return
	}

	for i, nov := range novedades {
		logs.Info("AprobarEntrada -> registrando novedad %d/%d: %+v", i+1, len(novedades), nov)
		outputError = movimientosArka.PostNovedadElemento(&nov)
		if outputError != nil {
			logs.Error("AprobarEntrada -> error en PostNovedadElemento: %v", outputError)
			return
		}
	}

	resultado_.Movimiento.FechaCorte = utilsHelper.Time(timebogota.TiempoBogota())
	logs.Info("AprobarEntrada -> actualizando movimiento Id=%d con FechaCorte=%v y EstadoMovimientoId=%+v",
		resultado_.Movimiento.Id, resultado_.Movimiento.FechaCorte, resultado_.Movimiento.EstadoMovimientoId)

	outputError = movimientosArka.PutMovimiento(&resultado_.Movimiento, resultado_.Movimiento.Id)
	logs.Info("AprobarEntrada -> PutMovimiento outputError=%v", outputError)
	logs.Info("==== FIN entradaHelper.AprobarEntrada ====")
	return
}

func getFormato(entradaId int, resultado *models.ResultadoMovimiento) (formato models.FormatoBaseEntrada, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("getFormato - Unhandled Error!", "500")

	logs.Info("==== INICIO getFormato ====")
	logs.Info("getFormato -> entradaId=%d", entradaId)

	if resultado == nil {
		logs.Error("getFormato -> resultado nil")
		outputError = map[string]interface{}{
			"funcion": "getFormato - resultado",
			"err":     "resultado nil",
			"status":  "500",
		}
		return
	}

	movimiento, outputError := movimientosArka.GetMovimientoById(entradaId)
	logs.Info("getFormato -> GetMovimientoById outputError=%v movimiento=%+v", outputError, movimiento)

	if outputError != nil {
		logs.Error("getFormato -> error consultando movimiento")
		return
	}

	if movimiento == nil {
		logs.Error("getFormato -> movimiento nil")
		outputError = map[string]interface{}{
			"funcion": "getFormato - movimientosArka.GetMovimientoById",
			"err":     "movimiento nil",
			"status":  "404",
		}
		return
	}

	if movimiento.EstadoMovimientoId == nil {
		logs.Error("getFormato -> EstadoMovimientoId nil en movimiento Id=%d", movimiento.Id)
		outputError = map[string]interface{}{
			"funcion": "getFormato - EstadoMovimientoId",
			"err":     "EstadoMovimientoId nil",
			"status":  "500",
		}
		return
	}

	if movimiento.FormatoTipoMovimientoId == nil {
		logs.Error("getFormato -> FormatoTipoMovimientoId nil en movimiento Id=%d", movimiento.Id)
		outputError = map[string]interface{}{
			"funcion": "getFormato - FormatoTipoMovimientoId",
			"err":     "FormatoTipoMovimientoId nil",
			"status":  "500",
		}
		return
	}

	if movimiento.EstadoMovimientoId.Nombre != "Entrada En Trámite" {
		logs.Error("getFormato -> estado inválido: %q", movimiento.EstadoMovimientoId.Nombre)
		outputError = map[string]interface{}{
			"funcion": "getFormato - EstadoMovimientoId.Nombre",
			"err":     "el movimiento no está en estado 'Entrada En Trámite'",
			"status":  "400",
		}
		return
	}

	resultado.Movimiento = *movimiento
	logs.Info("getFormato -> resultado.Movimiento=%+v", resultado.Movimiento)

	if resultado.Movimiento.ConsecutivoId == nil || *resultado.Movimiento.ConsecutivoId == 0 {
		resultado.Error = "No se puede continuar con el cálculo de la transaccón contable. Contacte soporte."
		logs.Error("getFormato -> ConsecutivoId nil o 0")
		return
	}

	logs.Info("getFormato -> Detalle crudo=%s", resultado.Movimiento.Detalle)
	outputError = utilsHelper.Unmarshal(resultado.Movimiento.Detalle, &formato)
	logs.Info("getFormato -> formato parseado=%+v outputError=%v", formato, outputError)
	if outputError != nil {
		return
	}

	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&resultado.Movimiento.EstadoMovimientoId.Id, "Entrada Aprobada")
	logs.Info("getFormato -> GetEstadoMovimientoIdByNombre outputError=%v nuevoEstadoId=%d",
		outputError, resultado.Movimiento.EstadoMovimientoId.Id)
	if outputError != nil {
		return
	}

	logs.Info("==== FIN getFormato ====")
	return
}

func getTerceroEntrada(detalle models.FormatoBaseEntrada, resutado *models.ResultadoMovimiento) (terceroId int, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("getTerceroEntrada - Unhandled Error!", "500")

	logs.Info("==== INICIO getTerceroEntrada ====")
	logs.Info("getTerceroEntrada -> detalle=%+v", detalle)

	var historico []models.HistoricoActa
	query := "Activo:true,ActaRecibidoId__Id:" + strconv.Itoa(detalle.ActaRecibidoId)

	if detalle.ActaRecibidoId > 0 {
		logs.Info("getTerceroEntrada -> consultando historico con query=%s", query)
		historico, outputError = actaRecibido.GetAllHistoricoActa(query, "", "FechaCreacion", "desc", "", "1")
		logs.Info("getTerceroEntrada -> historico len=%d outputError=%v historico=%+v", len(historico), outputError, historico)

		if outputError != nil || len(historico) != 1 {
			if len(historico) != 1 {
				resutado.Error = "No se pudo consultar la información del acta. Contacte soporte."
			}
			logs.Error("getTerceroEntrada -> error/longitud inesperada en historico")
			return
		}
	}

	if detalle.ActaRecibidoId > 0 && historico[0].ActaRecibidoId.TipoActaId.CodigoAbreviacion == "REG" {
		terceroId = historico[0].ProveedorId
		logs.Info("getTerceroEntrada -> tercero tomado del historico: %d", terceroId)
	} else {
		terceroId, outputError = terceros.GetTerceroUD()
		logs.Info("getTerceroEntrada -> terceros.GetTerceroUD terceroId=%d outputError=%v", terceroId, outputError)
	}

	if terceroId == 0 {
		resutado.Error = "No se pudo consultar el tercero para asociar a la transacción contable. Contacte soporte."
		logs.Error("getTerceroEntrada -> terceroId quedó en 0")
	}

	logs.Info("==== FIN getTerceroEntrada ====")
	return
}

func getElementosEntrada(detalle models.FormatoBaseEntrada, movimientoId int, resultado *models.ResultadoMovimiento) (elementos []*models.Elemento, novedades []models.NovedadElemento, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("getElementosEntrada - Unhandled Error!", "500")

	logs.Info("==== INICIO getElementosEntrada ====")
	logs.Info("getElementosEntrada -> movimientoId=%d detalle=%+v", movimientoId, detalle)

	if detalle.ActaRecibidoId == 0 && len(detalle.Elementos) == 0 {
		resultado.Error = errNoElementos
		logs.Error("getElementosEntrada -> sin ActaRecibidoId y sin elementos")
		return
	}

	if detalle.ActaRecibidoId > 0 {
		query := "Activo:true,ActaRecibidoId__Id:" + strconv.Itoa(detalle.ActaRecibidoId)
		logs.Info("getElementosEntrada -> consultando elementos por acta query=%s", query)

		elementos, outputError = actaRecibido.GetAllElemento(query, "ValorUnitario,ValorTotal,SubgrupoCatalogoId,TipoBienId", "SubgrupoCatalogoId", "desc", "", "-1")
		logs.Info("getElementosEntrada -> GetAllElemento len=%d outputError=%v", len(elementos), outputError)

		if len(elementos) == 0 {
			resultado.Error = errNoElementos
			logs.Error("getElementosEntrada -> no se encontraron elementos")
		}
	} else if len(detalle.Elementos) > 0 {
		for i, el := range detalle.Elementos {
			logs.Info("getElementosEntrada -> procesando elemento %d/%d: %+v", i+1, len(detalle.Elementos), el)

			if el.VidaUtil == nil || el.ValorLibros == nil || el.ValorResidual == nil {
				resultado.Error = "No se indicó correctamente el nuevo valor de los elementos. Rechace la entrada y haga la respectiva edición."
				logs.Error("getElementosEntrada -> VidaUtil/ValorLibros/ValorResidual nil")
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
			logs.Info("getElementosEntrada -> novedad generada: %+v", novedad)

			historial, err := movimientosArka.GetHistorialElemento(el.Id, true)
			logs.Info("getElementosEntrada -> GetHistorialElemento el.Id=%d historial=%+v err=%v", el.Id, historial, err)
			if err != nil {
				outputError = err
				return
			} else if historial == nil {
				resultado.Error = "No se pudo consultar la parametrización de los elementos. Contacte soporte."
				logs.Error("getElementosEntrada -> historial nil")
				return
			}

			valor, _, _, _, err := inventarioHelper.GetUltimoValor(*historial)
			logs.Info("getElementosEntrada -> GetUltimoValor valor=%v err=%v ValorLibros=%v", valor, err, *el.ValorLibros)
			if err != nil || math.Abs(*el.ValorLibros-valor) == 0 {
				outputError = err
				logs.Error("getElementosEntrada -> error o diferencia cero. err=%v diferencia=%v", err, math.Abs(*el.ValorLibros-valor))
				return
			}

			if historial.Elemento.ElementoActaId == nil {
				logs.Error("getElementosEntrada -> historial.Elemento.ElementoActaId nil")
				outputError = map[string]interface{}{
					"funcion": "getElementosEntrada - historial.Elemento.ElementoActaId",
					"err":     "ElementoActaId nil",
					"status":  "500",
				}
				return
			}

			var elementoActa models.Elemento
			outputError = actaRecibido.GetElementoById(*historial.Elemento.ElementoActaId, &elementoActa)
			logs.Info("getElementosEntrada -> GetElementoById id=%d outputError=%v elementoActa=%+v",
				*historial.Elemento.ElementoActaId, outputError, elementoActa)
			if outputError != nil {
				return
			}

			if elementoActa.Cantidad == 0 {
				logs.Error("getElementosEntrada -> elementoActa.Cantidad es 0")
				outputError = map[string]interface{}{
					"funcion": "getElementosEntrada - elementoActa.Cantidad",
					"err":     "cantidad del elemento acta es 0",
					"status":  "500",
				}
				return
			}

			elementoActa.ValorUnitario = *el.ValorLibros / float64(elementoActa.Cantidad)
			elementoActa.ValorTotal = math.Abs(*el.ValorLibros - valor)
			elementos = append(elementos, &elementoActa)

			logs.Info("getElementosEntrada -> elementoActa ajustado: %+v", elementoActa)
		}
	}

	logs.Info("getElementosEntrada -> retorna len(elementos)=%d len(novedades)=%d outputError=%v resultado.Error=%q",
		len(elementos), len(novedades), outputError, resultado.Error)
	logs.Info("==== FIN getElementosEntrada ====")
	return
}

func contabilidadEntrada(resultado_ *models.ResultadoMovimiento, formatoEntrada models.FormatoBaseEntrada, elementos []*models.Elemento, terceroId int) (outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("contabilidadEntrada - Unhandled Error!", "500")

	logs.Info("==== INICIO contabilidadEntrada ====")
	logs.Info("contabilidadEntrada -> terceroId=%d len(elementos)=%d formatoEntrada=%+v", terceroId, len(elementos), formatoEntrada)

	if resultado_ == nil {
		outputError = map[string]interface{}{
			"funcion": "contabilidadEntrada - resultado_",
			"err":     "resultado_ nil",
			"status":  "500",
		}
		return
	}

	if resultado_.Movimiento.ConsecutivoId == nil {
		outputError = map[string]interface{}{
			"funcion": "contabilidadEntrada - ConsecutivoId",
			"err":     "ConsecutivoId nil",
			"status":  "500",
		}
		return
	}

	if resultado_.Movimiento.FormatoTipoMovimientoId == nil {
		logs.Error("contabilidadEntrada -> FormatoTipoMovimientoId nil. Movimiento=%+v", resultado_.Movimiento)
		outputError = map[string]interface{}{
			"funcion": "contabilidadEntrada - FormatoTipoMovimientoId",
			"err":     "FormatoTipoMovimientoId nil",
			"status":  "500",
		}
		return
	}

	if resultado_.Movimiento.FormatoTipoMovimientoId.Id == 0 {
		logs.Error("contabilidadEntrada -> FormatoTipoMovimientoId.Id=0. Movimiento=%+v", resultado_.Movimiento)
		outputError = map[string]interface{}{
			"funcion": "contabilidadEntrada - FormatoTipoMovimientoId.Id",
			"err":     "FormatoTipoMovimientoId.Id en 0",
			"status":  "500",
		}
		return
	}

	if len(elementos) == 0 {
		logs.Error("contabilidadEntrada -> len(elementos)=0")
		outputError = map[string]interface{}{
			"funcion": "contabilidadEntrada - elementos",
			"err":     "no hay elementos para contabilizar",
			"status":  "400",
		}
		return
	}

	logs.Info("contabilidadEntrada -> Movimiento=%+v", resultado_.Movimiento)
	logs.Info("contabilidadEntrada -> detalle movimiento crudo=%s", resultado_.Movimiento.Detalle)

	detalleContable, outputError := descripcionMovimientoContable(resultado_.Movimiento.Detalle)
	logs.Info("contabilidadEntrada -> descripcionMovimientoContable detalleContable=%q outputError=%v", detalleContable, outputError)
	if outputError != nil {
		return
	}

	var transaccion = models.TransaccionMovimientos{
		ConsecutivoId: *resultado_.Movimiento.ConsecutivoId,
	}
	bufferCuentas := make(map[string]models.CuentaContable)

	logs.Info("contabilidadEntrada -> FormatoTipoMovimientoId=%+v", resultado_.Movimiento.FormatoTipoMovimientoId)
	logs.Info("contabilidadEntrada -> antes de CalcularMovimientosContables, transaccion=%+v", transaccion)

	resultado_.Error, outputError = asientoContable.CalcularMovimientosContables(
		elementos,
		detalleContable,
		0,
		resultado_.Movimiento.FormatoTipoMovimientoId.Id,
		terceroId,
		terceroId,
		bufferCuentas,
		nil,
		&transaccion.Movimientos,
	)
	logs.Info("contabilidadEntrada -> CalcularMovimientosContables resultado.Error=%q outputError=%v len(movimientos)=%d bufferCuentas=%+v",
		resultado_.Error, outputError, len(transaccion.Movimientos), bufferCuentas)
	if outputError != nil || resultado_.Error != "" {
		logs.Error("contabilidadEntrada -> aborta en CalcularMovimientosContables")
		return
	}

	resultado_.Error, outputError = asientoContable.CreateTransaccionContable(getTipoComprobanteEntradas(), "Entrada Almacén", &transaccion)
	logs.Info("contabilidadEntrada -> CreateTransaccionContable resultado.Error=%q outputError=%v transaccion=%+v",
		resultado_.Error, outputError, transaccion)
	if outputError != nil || resultado_.Error != "" {
		logs.Error("contabilidadEntrada -> aborta en CreateTransaccionContable")
		return
	}

	resultado_.TransaccionContable.Concepto = transaccion.Descripcion
	resultado_.TransaccionContable.Fecha = transaccion.FechaTransaccion

	resultado_.TransaccionContable.Movimientos, outputError = asientoContable.GetDetalleContable(transaccion.Movimientos, bufferCuentas)
	logs.Info("contabilidadEntrada -> GetDetalleContable outputError=%v movimientosDetalle=%+v",
		outputError, resultado_.TransaccionContable.Movimientos)
	if outputError != nil {
		logs.Error("contabilidadEntrada -> aborta en GetDetalleContable")
		return
	}

	logs.Info("contabilidadEntrada -> ANTES de PostTrContable transaccion=%+v", transaccion)
	postRes, outputError := movimientosContables.PostTrContable(&transaccion)
	logs.Info("contabilidadEntrada -> DESPUÉS de PostTrContable response=%+v outputError=%v", postRes, outputError)

	if outputError != nil {
		logs.Error("contabilidadEntrada -> error en PostTrContable: %v", outputError)
		return
	}

	logs.Info("==== FIN contabilidadEntrada ====")
	return
}

// descripcionMovimientoContable Genera la descipción de cada uno de los movimientos contables asociados a una entrada.
func descripcionMovimientoContable(detalle string) (detalle_ string, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("descripcionMovimientoContable - Unhandled Error!", "500")

	logs.Info("==== INICIO descripcionMovimientoContable ====")
	logs.Info("descripcionMovimientoContable -> detalle crudo=%s", detalle)

	var mapDetalle map[string]interface{}
	outputError = utilsHelper.Unmarshal(detalle, &mapDetalle)
	logs.Info("descripcionMovimientoContable -> mapDetalle=%+v outputError=%v", mapDetalle, outputError)
	if outputError != nil {
		return
	}

	for k, v := range mapDetalle {
		logs.Info("descripcionMovimientoContable -> key=%s value=%#v tipo=%T", k, v, v)

		if k == "factura" {
			if v == nil {
				logs.Error("descripcionMovimientoContable -> factura viene nil")
				outputError = map[string]interface{}{
					"funcion": "descripcionMovimientoContable - factura",
					"err":     "factura nil",
					"status":  "400",
				}
				return
			}

			facturaFloat, ok := v.(float64)
			if !ok {
				logs.Error("descripcionMovimientoContable -> factura no es float64. value=%#v tipo=%T", v, v)
				outputError = map[string]interface{}{
					"funcion": "descripcionMovimientoContable - factura type assertion",
					"err":     fmt.Sprintf("se esperaba float64 en factura y llegó %T", v),
					"status":  "400",
				}
				return
			}

			var sop models.SoporteActa
			outputError = actaRecibido.GetSoporteById(int(facturaFloat), &sop)
			logs.Info("descripcionMovimientoContable -> GetSoporteById factura=%v outputError=%v soporte=%+v", facturaFloat, outputError, sop)
			if outputError != nil {
				return
			}

			detalle_ += "Factura: " + sop.Consecutivo + ", "
		} else if k != "consecutivo" && k != "ConsecutivoId" && k != "elementos" {
			k = strings.TrimSuffix(k, "_id")
			caser := cases.Title(language.Spanish)
			k = caser.String(k)
			detalle_ += k + ": " + fmt.Sprintf("%v", v) + ", "
		}
	}

	detalle_ = strings.TrimSuffix(detalle_, ", ")
	logs.Info("descripcionMovimientoContable -> detalle final=%q", detalle_)
	logs.Info("==== FIN descripcionMovimientoContable ====")
	return
}
