package bodegaConsumoHelper

import (
	"net/url"

	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	crudTerceros "github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

const estadoSolicitudPendiente = "Solicitud Pendiente"

func GetAllSolicitudes(pendientesOnly bool, solictudes_ *[]models.DetalleSolicitudBodega) (outputError map[string]interface{}) {

	funcion := "GetAllSolicitudes"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		solicitudes []*models.Movimiento
		terceros    map[int]models.IdentificacionTercero
	)

	query := "limit=-1&sortby=FechaCreacion&order=desc&query=Activo:true,FormatoTipoMovimientoId__CodigoAbreviacion:SOL_BOD"
	if pendientesOnly {
		query += ",EstadoMovimientoId__Nombre:" + url.QueryEscape(estadoSolicitudPendiente)
	}

	if movimientos, err := movimientosArka.GetAllMovimiento(query); err != nil {
		return err
	} else {
		solicitudes = movimientos
	}

	terceros = make(map[int]models.IdentificacionTercero)
	for _, sol := range solicitudes {

		var (
			detalle   models.FormatoSolicitudBodega
			solicitud models.DetalleSolicitudBodega
		)

		solicitud.Movimiento = *sol

		if err := utilsHelper.Unmarshal(sol.Detalle, &detalle); err != nil {
			return err
		}

		if detalle.Funcionario > 0 {
			if val, ok := terceros[detalle.Funcionario]; ok {
				solicitud.Solicitante = val
			} else {
				if tercero, err := crudTerceros.GetNombreTerceroById(detalle.Funcionario); err != nil {
					return err
				} else if tercero != nil {
					terceros[detalle.Funcionario] = *tercero
					solicitud.Solicitante = *tercero
				}
			}
		}

		*solictudes_ = append(*solictudes_, solicitud)

	}

	return
}
