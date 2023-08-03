package entradaHelper

import (
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// RegistrarEntrada Crea registro de entrada en estado en trámite
func RegistrarEntrada(data *models.TransaccionEntrada, etl bool, resultado *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("RegistrarEntrada - Unhandled Error!", "500")

	resultado.Movimiento = models.Movimiento{
		Observacion:             data.Observacion,
		Activo:                  true,
		FormatoTipoMovimientoId: &models.FormatoTipoMovimiento{},
		EstadoMovimientoId:      &models.EstadoMovimiento{},
	}

	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&resultado.Movimiento.EstadoMovimientoId.Id, "Entrada En Trámite")
	if outputError != nil {
		return
	}

	outputError = movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&resultado.Movimiento.FormatoTipoMovimientoId.Id, data.FormatoTipoMovimientoId)
	if outputError != nil {
		return
	}

	outputError = crearDetalleEntrada(data.Detalle, &resultado.Movimiento.Detalle)
	if outputError != nil {
		return
	}

	var acta models.TransaccionActaRecibido
	if data.Detalle.ActaRecibidoId > 0 {
		outputError = actaRecibido.GetTransaccionActaRecibidoById(data.Detalle.ActaRecibidoId, false, &acta)
		if outputError != nil {
			return
		} else if acta.UltimoEstado.EstadoActaId.CodigoAbreviacion != "Aceptada" {
			resultado.Error = "El acta asociada no está en estado aceptada y no se puede continuar."
			return
		}
	}

	outputError = getConsecutivoEntrada(&resultado.Movimiento, etl)
	if outputError != nil {
		return
	}

	if data.Detalle.ActaRecibidoId > 0 {
		resultado.Error, outputError = asignarPlacas(data.Detalle.ActaRecibidoId, &acta.Elementos)
		if outputError != nil || resultado.Error != "" {
			return
		} else if len(acta.Elementos) == 0 {
			resultado.Error = "No se encontraron elementos asociados al acta."
			return
		}
	}

	outputError = movimientosArka.PostMovimiento(&resultado.Movimiento)
	if outputError != nil {
		return
	}

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

	defer errorctrl.ErrorControlFunction("crearDetalleEntrada - Unhandled Error!", "500")

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
