package bajasHelper

import (
	"net/url"
	"strings"

	"github.com/udistrital/arka_mid/helpers/crud/configuracion"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/mid/autenticacion"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// GetAll Consulta información general de todas las bajas filtrando por usuario o las que están pendientes por revisar.
func GetAll(user string, revComite, revAlmacen bool, bajas *[]models.DetalleBaja) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("GetAll - Unhandled Error!", "500")

	var solicitudes []*models.Movimiento

	if err := loadBajas(user, revAlmacen, revComite, &solicitudes); err != nil {
		return err
	}

	if len(solicitudes) == 0 {
		return
	}

	bufferTerceros := make(map[int]string)

	for _, solicitud := range solicitudes {

		var detalle *models.FormatoBaja

		if err := utilsHelper.Unmarshal(solicitud.Detalle, &detalle); err != nil {
			return err
		}

		if _, ok := bufferTerceros[detalle.Funcionario]; !ok {
			if tercero, err := terceros.GetTerceroById(detalle.Funcionario); err != nil {
				return err
			} else {
				bufferTerceros[detalle.Funcionario] = tercero.NombreCompleto
			}
		}

		if _, ok := bufferTerceros[detalle.Revisor]; !ok {
			if tercero, err := terceros.GetTerceroById(detalle.Revisor); err != nil {
				return err
			} else {
				bufferTerceros[detalle.Revisor] = tercero.NombreCompleto
			}
		}

		baja := models.DetalleBaja{
			Id:                 solicitud.Id,
			Consecutivo:        *solicitud.Consecutivo,
			FechaCreacion:      solicitud.FechaCreacion.String(),
			FechaRevisionA:     detalle.FechaRevisionA,
			FechaRevisionC:     detalle.FechaRevisionC,
			Funcionario:        bufferTerceros[detalle.Funcionario],
			Revisor:            bufferTerceros[detalle.Revisor],
			TipoBaja:           solicitud.FormatoTipoMovimientoId.Nombre,
			EstadoMovimientoId: solicitud.EstadoMovimientoId.Id,
		}
		*bajas = append(*bajas, baja)
	}

	return

}

// loadBajas Consulta lista de bajas asociadas a un usuario de acuerdo a las revisiones y permisos del usuario
func loadBajas(user string, revAlmacen, revComite bool, bajas *[]*models.Movimiento) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("loadBajas - Unhandled Error!", "500")

	var (
		terceroId int
		roles     []string
		opciones  []*models.PerfilXMenuOpcion
	)

	if revAlmacen || revComite {

		payload := "limit=-1&sortby=Id&order=desc&query=Activo:true,EstadoMovimientoId__Nombre:"
		if revComite {
			payload += url.QueryEscape("Baja En Comité")
		} else if revAlmacen {
			payload += url.QueryEscape("Baja En Trámite")
		}

		if solicitudes_, err := movimientosArka.GetAllMovimiento(payload); err != nil {
			return err
		} else {
			*bajas = solicitudes_
		}

		return

	}

	if err := autenticacion.GetInfoUser(user, &terceroId, &roles); err != nil {
		return err
	}

	if terceroId == 0 {
		return
	}

	query := "limit=-1&query=Opcion__Nombre:bajasVerTodaSolicitud,Perfil__Nombre__in:" + strings.Join(roles, "|")
	if err := configuracion.GetAllPerfilXMenuOpcion(query, &opciones); err != nil {
		return err
	}

	if len(opciones) > 0 {
		query := "limit=-1&query=Activo:true,EstadoMovimientoId__Nombre__startswith:Baja"
		if tr_, err := movimientosArka.GetAllMovimiento(query); err != nil {
			return err
		} else {
			*bajas = tr_
		}
	} else {
		if err := movimientosArka.GetBajasByTerceroId(terceroId, bajas); err != nil {
			return err
		}
	}

	return

}
