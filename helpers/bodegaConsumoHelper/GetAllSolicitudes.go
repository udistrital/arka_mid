package bodegaConsumoHelper

import (
	"net/url"
	"strings"

	"github.com/udistrital/arka_mid/helpers/crud/configuracion"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	crudTerceros "github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/mid/autenticacion"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

const estadoSolicitudPendiente = "Solicitud Pendiente"

func GetAllSolicitudes(user string, revision bool, solictudes_ *[]models.DetalleSolicitudBodega) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GetAllSolicitudes - Unhandled Error!", "500")

	var (
		solicitudes []*models.Movimiento
		terceros    map[int]models.IdentificacionTercero
	)

	if err := loadSolicitudes(user, revision, &solicitudes); err != nil {
		return err
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

func loadSolicitudes(user string, revision bool, solicitudes *[]*models.Movimiento) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("loadSolicitudes - Unhandled Error!", "500")

	var (
		terceroId int
		roles     []string
		opciones  []*models.PerfilXMenuOpcion
	)

	payload := "limit=-1&sortby=Id&order=desc&query=Activo:true,FormatoTipoMovimientoId__CodigoAbreviacion:SOL_BOD"

	if revision {

		payload += ",EstadoMovimientoId__Nombre:" + url.QueryEscape(estadoSolicitudPendiente)

		if solicitudes_, err := movimientosArka.GetAllMovimiento(payload); err != nil {
			return err
		} else {
			*solicitudes = solicitudes_
		}

		return

	}

	if err := autenticacion.GetInfoUser(user, &terceroId, &roles); err != nil {
		return err
	}

	if terceroId == 0 {
		return
	}

	query := "limit=-1&query=Opcion__Nombre:bodegaVerTodasLasSolicitudes,Perfil__Nombre__in:" + strings.Join(roles, "|")
	if err := configuracion.GetAllPerfilXMenuOpcion(query, &opciones); err != nil {
		return err
	}

	if len(opciones) > 0 {
		if sol_, err := movimientosArka.GetAllMovimiento(payload); err != nil {
			return err
		} else {
			*solicitudes = sol_
		}
	} else {
		if err := movimientosArka.GetBodegaByTerceroId(terceroId, solicitudes); err != nil {
			return err
		}
	}

	return
}
