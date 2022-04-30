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

	funcion := "PostSolicitud"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		estadoId  int
		formatoId int
		detalle   string
	)

	ctxConsecutivo, _ := beego.AppConfig.Int("contxtSolicitudBodega")
	consecutivo := models.Consecutivo{
		ContextoId:  ctxConsecutivo,
		Year:        0,
		Descripcion: "Solicitud elementos bodega de consumo almacén de inventarios.",
		Activo:      true,
	}

	if err := consecutivos.Post(&consecutivo); err != nil {
		return err
	}

	if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&formatoId, "SOL_BOD"); err != nil {
		return err
	}

	if err := movimientosArka.GetEstadoMovimientoIdByNombre(&estadoId, estadoSolicitudPendiente); err != nil {
		return err
	}

	detalle_ := models.FormatoSolicitudBodega{
		ConsecutivoId: consecutivo.Id,
		Consecutivo:   strconv.Itoa(consecutivo.Consecutivo),
		Funcionario:   solicitud.Funcionario,
		Elementos:     solicitud.Elementos,
	}

	if err := utilsHelper.Marshal(detalle_, &detalle); err != nil {
		return err
	}

	*movimiento = models.Movimiento{
		Detalle:                 detalle,
		Activo:                  true,
		FormatoTipoMovimientoId: &models.FormatoTipoMovimiento{Id: formatoId},
		EstadoMovimientoId:      &models.EstadoMovimiento{Id: estadoId},
	}

	if err := movimientosArka.PostMovimiento(movimiento); err != nil {
		return err
	}

	return

}
