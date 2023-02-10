package bodegaConsumoHelper

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

func PostSolicitud(solicitud *models.FormatoSolicitudBodega, movimiento *models.Movimiento) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("PostSolicitud - Unhandled Error!", "500")

	movimiento.EstadoMovimientoId = &models.EstadoMovimiento{}
	movimiento.FormatoTipoMovimientoId = &models.FormatoTipoMovimiento{}

	outputError = movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&movimiento.FormatoTipoMovimientoId.Id, "SOL_BOD")
	if outputError != nil {
		return
	}

	outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&movimiento.EstadoMovimientoId.Id, estadoSolicitudPendiente)
	if outputError != nil {
		return
	}

	detalle_ := models.FormatoSolicitudBodega{
		Funcionario: solicitud.Funcionario,
		Elementos:   solicitud.Elementos,
	}

	outputError = utilsHelper.Marshal(detalle_, &movimiento.Detalle)
	if outputError != nil {
		return
	}

	ctxConsecutivo, _ := beego.AppConfig.Int("contxtSolicitudBodega")
	consecutivo := models.Consecutivo{
		ContextoId:  ctxConsecutivo,
		Year:        0,
		Descripcion: "Solicitud elementos bodega de consumo almacén de inventarios.",
		Activo:      true,
	}

	outputError = consecutivos.Post(&consecutivo)
	if outputError != nil {
		return
	}

	movimiento.Activo = true
	movimiento.ConsecutivoId = &consecutivo.Id
	movimiento.Consecutivo = utilsHelper.String(strconv.Itoa(consecutivo.Consecutivo))

	outputError = movimientosArka.PostMovimiento(movimiento)

	return
}
