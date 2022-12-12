package entradaHelper

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// RegistrarEntrada Crea registro de entrada en estado en trámite
func RegistrarEntrada(data *models.TransaccionEntrada, etl bool, resultado *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	funcion := "RegistrarEntrada - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	var (
		acta             models.TransaccionActaRecibido
		tipoMovimiento   int
		estadoMovimiento int
		detalle          string
	)

	if data.Detalle.ActaRecibidoId <= 0 {
		err := "Se debe indicar un acta de recibido válida."
		return errorctrl.Error(funcion, err, "400")
	}

	if err := movimientosArka.GetEstadoMovimientoIdByNombre(&estadoMovimiento, "Entrada En Trámite"); err != nil {
		return err
	}

	if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&tipoMovimiento, data.FormatoTipoMovimientoId); err != nil {
		return err
	}

	if err := actaRecibido.GetTransaccionActaRecibidoById(data.Detalle.ActaRecibidoId, false, &acta); err != nil {
		return err
	}

	if err := crearDetalleEntrada(&data.Detalle, etl, nil, &detalle); err != nil {
		return err
	}

	if !etl {
		if msjErr, err := asignarPlacas(data.Detalle.ActaRecibidoId, &acta.Elementos); err != nil {
			return err
		} else if msjErr != "" {
			resultado.Error = msjErr
			return
		} else if len(acta.Elementos) == 0 {
			resultado.Error = "No se encontraron elementos asociados al acta."
			return
		}
	}

	resultado.Movimiento = models.Movimiento{
		Observacion:             data.Observacion,
		Detalle:                 detalle,
		Activo:                  true,
		FormatoTipoMovimientoId: &models.FormatoTipoMovimiento{Id: tipoMovimiento},
		EstadoMovimientoId:      &models.EstadoMovimiento{Id: estadoMovimiento},
	}

	if err := movimientosArka.PostMovimiento(&resultado.Movimiento); err != nil {
		return err
	}

	// Crea registro en table soporte_movimiento si es necesario
	if data.SoporteMovimientoId > 0 {
		soporte := models.SoporteMovimiento{
			DocumentoId:  data.SoporteMovimientoId,
			Activo:       true,
			MovimientoId: &models.Movimiento{Id: resultado.Movimiento.Id},
		}

		if err := movimientosArka.PostSoporteMovimiento(&soporte); err != nil {
			return err
		}

	}

	if !etl {
		acta.UltimoEstado.EstadoActaId.Id = 6
		acta.UltimoEstado.Id = 0

		if err := actaRecibido.PutTransaccionActaRecibido(data.Detalle.ActaRecibidoId, &acta); err != nil {
			return err
		}
	}

	return

}

// creaDetalleEntrada construye la data que será almacenada en la columna detalle según se requiera.
func crearDetalleEntrada(completo *models.FormatoBaseEntrada, etl bool, consecutivo_ *models.ConsecutivoMovimiento, necesario *string) (outputError map[string]interface{}) {

	funcion := "crearDetalleEntrada - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	var (
		detalle     map[string]interface{}
		consecutivo models.Consecutivo
	)

	if err := utilsHelper.FillStruct(completo, &detalle); err != nil {
		return err
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

	if !etl && consecutivo_ == nil {
		ctxConsecutivo, _ := beego.AppConfig.Int("contxtEntradaCons")
		if err := consecutivos.Get(ctxConsecutivo, "Entradas Arka", &consecutivo); err != nil {
			return err
		}

		detalle["consecutivo"] = consecutivos.Format("%05d", getTipoComprobanteEntradas(), &consecutivo)
		detalle["ConsecutivoId"] = consecutivo.Id
	} else if consecutivo_ != nil {
		detalle["consecutivo"] = consecutivo_.Consecutivo
		detalle["ConsecutivoId"] = consecutivo_.ConsecutivoId
	}

	if err := utilsHelper.Marshal(detalle, necesario); err != nil {
		return err
	}

	return

}
