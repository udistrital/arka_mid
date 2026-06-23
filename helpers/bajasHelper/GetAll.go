package bajasHelper

import (
	"fmt"
	"net/url"
	"strconv"
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

		if err := cargarNombreTerceroBaja(detalle.Funcionario, bufferTerceros); err != nil {
			return err
		}

		if err := cargarNombreTerceroBaja(detalle.Revisor, bufferTerceros); err != nil {
			return err
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

func cargarNombreTerceroBaja(terceroID int, buffer map[int]string) (outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("cargarNombreTerceroBaja - Unhandled Error!", "500")

	if buffer == nil {
		return nil
	}
	if terceroID <= 0 {
		buffer[terceroID] = ""
		return nil
	}
	if _, ok := buffer[terceroID]; ok {
		return nil
	}

	tercero, err := terceros.GetTerceroById(terceroID)
	if err != nil {
		if esTerceroNoEncontrado(err) {
			buffer[terceroID] = strconv.Itoa(terceroID)
			return nil
		}
		return err
	}

	buffer[terceroID] = tercero.NombreCompleto
	if buffer[terceroID] == "" {
		buffer[terceroID] = strconv.Itoa(terceroID)
	}
	return nil
}

func esTerceroNoEncontrado(err map[string]interface{}) bool {
	if err == nil {
		return false
	}
	return strings.Contains(fmt.Sprint(err["err"]), "http 404:")
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

		if solicitudes_, _, err := movimientosArka.GetAllMovimiento(payload); err != nil {
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
		if tr_, _, err := movimientosArka.GetAllMovimiento(query); err != nil {
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
