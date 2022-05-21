package trasladoshelper

import (
	"net/url"
	"strings"

	"github.com/udistrital/arka_mid/helpers/crud/configuracion"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/oikos"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/mid/autenticacion"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// GetAllTraslados Consulta información general de todos los traslados asociados a un usuario determinado. Permite filtrar por los que están pendientes por aprobar o confirmar
func GetAllTraslados(user string, confirmar, aprobar bool, traslados_ *[]*models.DetalleTrasladoLista) (outputError map[string]interface{}) {

	funcion := "GetAllTraslados"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var traslados []*models.Movimiento
	if err := getTraslados(user, confirmar, aprobar, &traslados); err != nil {
		return err
	}

	if len(traslados) == 0 {
		return nil
	}

	tercerosBuffer := make(map[int]interface{})
	ubicacionesBuffer := make(map[int]interface{})

	for _, solicitud := range traslados {

		var (
			detalle    *models.FormatoTraslado
			Tercero_   string
			Revisor_   string
			Ubicacion_ string
		)

		if err := utilsHelper.Unmarshal(solicitud.Detalle, &detalle); err != nil {
			return err
		}

		requestTercero := func(id int) func() (interface{}, map[string]interface{}) {
			return func() (interface{}, map[string]interface{}) {
				if Tercero, err := terceros.GetTerceroById(id); err == nil {
					return Tercero, nil
				}
				return nil, nil
			}
		}

		requestUbicacion := func(id int) func() (interface{}, map[string]interface{}) {
			return func() (interface{}, map[string]interface{}) {
				if Ubicacion, err := oikos.GetSedeDependenciaUbicacion(id); err == nil {
					return Ubicacion, nil
				}
				return nil, nil
			}
		}

		if v, err := utilsHelper.BufferGeneric(detalle.FuncionarioDestino, tercerosBuffer, requestTercero(detalle.FuncionarioDestino), nil, nil); err == nil {
			if v2, ok := v.(*models.Tercero); ok {
				Tercero_ = v2.NombreCompleto
			}
		}

		if v, err := utilsHelper.BufferGeneric(detalle.FuncionarioOrigen, tercerosBuffer, requestTercero(detalle.FuncionarioOrigen), nil, nil); err == nil {
			if v2, ok := v.(*models.Tercero); ok {
				Revisor_ = v2.NombreCompleto
			}
		}

		if v, err := utilsHelper.BufferGeneric(detalle.Ubicacion, ubicacionesBuffer, requestUbicacion(detalle.Ubicacion), nil, nil); err == nil {
			if v2, ok := v.(*models.DetalleSedeDependencia); ok {
				Ubicacion_ = v2.Ubicacion.EspacioFisicoId.Nombre
			}
		}

		baja := models.DetalleTrasladoLista{
			Id:                 solicitud.Id,
			Consecutivo:        detalle.Consecutivo,
			FechaCreacion:      solicitud.FechaCreacion.String(),
			FuncionarioOrigen:  Tercero_,
			FuncionarioDestino: Revisor_,
			Ubicacion:          Ubicacion_,
			EstadoMovimientoId: solicitud.EstadoMovimientoId.Id,
		}
		*traslados_ = append(*traslados_, &baja)
	}
	return

}

// getTraslados Consulta lista de traslados asociados a un usuario de acuerdo al filtro y permisos del usuario
func getTraslados(user string, confirmar, aprobar bool, traslados *[]*models.Movimiento) (outputError map[string]interface{}) {

	var (
		terceroId int
		roles     []string
		opciones  []*models.PerfilXMenuOpcion
	)

	if err := getInfoUser(user, &terceroId, &roles); err != nil {
		return err
	}

	if terceroId == 0 {
		return
	}

	if confirmar && !aprobar {
		if err := movimientosArka.GetTrasladosByTerceroId(terceroId, confirmar, traslados); err != nil {
			return err
		}

		return

	} else if !confirmar && !aprobar {

		query := "limit=-1&query=Opcion__Nombre:trasladosVerTodaSolicitud,Perfil__Nombre__in:" + strings.Join(roles, "|")
		if err := configuracion.GetAllPerfilXMenuOpcion(query, &opciones); err != nil {
			return err
		}

		if len(opciones) > 0 {
			query := "limit=-1&query=Activo:true,FormatoTipoMovimientoId__CodigoAbreviacion:SOL_TRD"
			if tr_, err := movimientosArka.GetAllMovimiento(query); err != nil {
				return err
			} else {
				*traslados = tr_
			}
		} else {
			if err := movimientosArka.GetTrasladosByTerceroId(terceroId, confirmar, traslados); err != nil {
				return err
			}
		}

	} else if aprobar {
		query := "limit=-1&query=Activo:true,EstadoMovimientoId__Nombre:" + url.QueryEscape("Traslado Confirmado")
		if tr_, err := movimientosArka.GetAllMovimiento(query); err != nil {
			return err
		} else {
			*traslados = tr_
		}
	}

	return

}

// getInfoUser Consulta los roles y el TerceroId asociado a un usuario determinado
func getInfoUser(usr string, terceroId *int, roles *[]string) (outputError map[string]interface{}) {

	var (
		user    models.UsuarioAutenticacion
		tercero models.DatosIdentificacion
	)

	if data, err := autenticacion.DataUsuario(usr); err != nil {
		return err
	} else {
		user = data
		*roles = user.Role
	}

	if user.Documento == "" {
		return
	}

	if data, err := terceros.GetTerceroByDoc(user.Documento); err != nil {
		return err
	} else {
		tercero = *data
	}

	if tercero.TerceroId != nil && tercero.TerceroId.Id > 0 {
		*terceroId = tercero.TerceroId.Id
	}

	return

}
