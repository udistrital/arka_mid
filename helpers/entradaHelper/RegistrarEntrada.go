package entradaHelper

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// RegistrarEntrada Crea registro de entrada en estado en trámite
func RegistrarEntrada(data *models.TransaccionEntrada, etl bool, resultado *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("RegistrarEntrada - Unhandled Error!", "500")

	resultado.Movimiento = models.Movimiento{
		Observacion:             data.Observacion,
		Activo:                  true,
		FormatoTipoMovimientoId: &models.FormatoTipoMovimiento{},
		EstadoMovimientoId:      &models.EstadoMovimiento{},
	}

	logs.Info("DEBUG [RegistrarEntrada] PASO 1: GetEstadoMovimientoIdByNombre")
	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&resultado.Movimiento.EstadoMovimientoId.Id, "Entrada En Trámite")
	if outputError != nil {
		logs.Error("DEBUG [RegistrarEntrada] FALLO EN PASO 1: %v", outputError)
		return
	}
	logs.Info("DEBUG [RegistrarEntrada] PASO 1 OK - EstadoId: %d", resultado.Movimiento.EstadoMovimientoId.Id)

	logs.Info("DEBUG [RegistrarEntrada] PASO 2: GetFormatoTipoMovimientoIdByCodigoAbreviacion(%s)", data.FormatoTipoMovimientoId)
	outputError = movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&resultado.Movimiento.FormatoTipoMovimientoId.Id, data.FormatoTipoMovimientoId)
	if outputError != nil {
		logs.Error("DEBUG [RegistrarEntrada] FALLO EN PASO 2: %v", outputError)
		return
	}
	logs.Info("DEBUG [RegistrarEntrada] PASO 2 OK - FormatoId: %d", resultado.Movimiento.FormatoTipoMovimientoId.Id)

	logs.Info("DEBUG [RegistrarEntrada] PASO 3: crearDetalleEntrada")
	outputError = crearDetalleEntrada(data.Detalle, &resultado.Movimiento.Detalle)
	if outputError != nil {
		logs.Error("DEBUG [RegistrarEntrada] FALLO EN PASO 3: %v", outputError)
		return
	}
	logs.Info("DEBUG [RegistrarEntrada] PASO 3 OK")

	var acta models.TransaccionActaRecibido
	if data.Detalle.ActaRecibidoId > 0 {
		logs.Info("DEBUG [RegistrarEntrada] PASO 4: GetTransaccionActaRecibidoById(%d)", data.Detalle.ActaRecibidoId)
		outputError = actaRecibido.GetTransaccionActaRecibidoById(data.Detalle.ActaRecibidoId, false, &acta)
		if outputError != nil {
			logs.Error("DEBUG [RegistrarEntrada] FALLO EN PASO 4: %v", outputError)
			return
		} else if acta.UltimoEstado.EstadoActaId.CodigoAbreviacion != "Aceptada" {
			logs.Warn("DEBUG [RegistrarEntrada] PASO 4: Acta no aceptada, estado: %s", acta.UltimoEstado.EstadoActaId.CodigoAbreviacion)
			resultado.Error = "El acta asociada no está en estado aceptada y no se puede continuar."
			return
		}
		logs.Info("DEBUG [RegistrarEntrada] PASO 4 OK")
	}

	logs.Info("DEBUG [RegistrarEntrada] PASO 5: getConsecutivoEntrada")
	outputError = getConsecutivoEntrada(&resultado.Movimiento, etl)
	if outputError != nil {
		logs.Error("DEBUG [RegistrarEntrada] FALLO EN PASO 5: %v", outputError)
		return
	}
	logs.Info("DEBUG [RegistrarEntrada] PASO 5 OK - Consecutivo: %v", resultado.Movimiento.Consecutivo)

	if data.Detalle.ActaRecibidoId > 0 {
		logs.Info("DEBUG [RegistrarEntrada] PASO 6: asignarPlacas")
		resultado.Error, outputError = asignarPlacas(data.Detalle.ActaRecibidoId, &acta.Elementos)
		if outputError != nil || resultado.Error != "" {
			logs.Error("DEBUG [RegistrarEntrada] FALLO EN PASO 6: err=%v, error=%s", outputError, resultado.Error)
			return
		} else if len(acta.Elementos) == 0 {
			resultado.Error = "No se encontraron elementos asociados al acta."
			logs.Error("DEBUG [RegistrarEntrada] PASO 6: Sin elementos")
			return
		}
		logs.Info("DEBUG [RegistrarEntrada] PASO 6 OK")
	}

	logs.Info("DEBUG [RegistrarEntrada] PASO 7: PostMovimiento")
	outputError = movimientosArka.PostMovimiento(&resultado.Movimiento)
	if outputError != nil {
		logs.Error("DEBUG [RegistrarEntrada] FALLO EN PASO 7: %v", outputError)
		return
	}
	logs.Info("DEBUG [RegistrarEntrada] PASO 7 OK - MovimientoId: %d", resultado.Movimiento.Id)

	if data.SoporteMovimientoId > 0 {
		soporte := models.SoporteMovimiento{
			DocumentoId:  data.SoporteMovimientoId,
			Activo:       true,
			MovimientoId: &models.Movimiento{Id: resultado.Movimiento.Id},
		}

		outputError = movimientosArka.PostSoporteMovimiento(&soporte)
		if outputError != nil {
			return
		}
	}

	if data.Detalle.ActaRecibidoId > 0 {
		acta.UltimoEstado.EstadoActaId.Id = 6
		acta.UltimoEstado.Id = 0
		outputError = actaRecibido.PutTransaccionActaRecibido(data.Detalle.ActaRecibidoId, &acta)
	}

	return
}

// creaDetalleEntrada construye la data que será almacenada en la columna detalle según se requiera.
func crearDetalleEntrada(completo models.FormatoBaseEntrada, necesario *string) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("crearDetalleEntrada - Unhandled Error!", "500")

	var detalle map[string]interface{}
	outputError = utilsHelper.FillStruct(completo, &detalle)
	if outputError != nil {
		return
	}

	if completo.ContratoId == 0 {
		delete(detalle, "contrato_id")
	}

	if completo.Divisa == "" {
		delete(detalle, "divisa")
	}

	if completo.Factura == 0 {
		delete(detalle, "factura")
	}

	if completo.OrdenadorGastoId == 0 {
		delete(detalle, "ordenador_gasto_id")
	}

	if len(completo.Elementos) == 0 {
		delete(detalle, "elementos")
	} else {
		elementos_, _ := detalle["elementos"].([]interface{})
		for _, elemento_ := range elementos_ {
			el, _ := elemento_.(map[string]interface{})
			if el["AprovechadoId"] == nil {
				delete(el, "AprovechadoId")
			}

			if el["ValorLibros"] == nil {
				delete(el, "ValorLibros")
			}

			if el["VidaUtil"] == nil {
				delete(el, "VidaUtil")
			}

			if el["ValorResidual"] == nil {
				delete(el, "ValorResidual")
			}
		}
	}

	if completo.RegistroImportacion == "" {
		delete(detalle, "num_reg_importacion")
	}

	if completo.SupervisorId == 0 {
		delete(detalle, "supervisor")
	}

	if completo.TRM == 0 {
		delete(detalle, "TRM")
	}

	if completo.VigenciaContrato == "" {
		delete(detalle, "vigencia_contrato")
	}

	outputError = utilsHelper.Marshal(detalle, necesario)
	return
}

func getConsecutivoEntrada(entrada *models.Movimiento, etl bool) (outputError map[string]interface{}) {

	if etl {
		return
	}

	if entrada.ConsecutivoId == nil || *entrada.ConsecutivoId <= 0 {
		var consecutivo models.Consecutivo
		outputError = consecutivos.Get("contxtEntradaCons", "Entradas Arka", &consecutivo)
		if outputError != nil {
			return
		}

		entrada.Consecutivo = utilsHelper.String(consecutivos.Format("%05d", getTipoComprobanteEntradas(), &consecutivo))
		entrada.ConsecutivoId = &consecutivo.Id
	}

	return
}

func getTipoComprobanteEntradas() string {
	return "P8"
}
