package actaRecibido

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/oikos"
	crudTerceros "github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/mid/autenticacion"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	e "github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
)

const (
	FormatoFecha = time.RFC3339
	FechaCero    = "0001-01-01T00:00:00Z"
)

var (
	zero time.Time
)

func init() {
	var err error
	zero, err = time.Parse(FormatoFecha, FechaCero)
	if err != nil {
		logs.Critical(err)
	}
	logs.Debug("time.Parse de fecha zero:", zero)
}

// GetAllActasRecibido ...
func GetAllActasRecibidoActivas(states []string, usrWSO2 string, limit, offset int,
	customQuery map[string]string) (historicoActa []models.ActaResumen, outputError map[string]interface{}) {
	const funcion = "GetAllActasRecibidoActivas - "
	defer e.ErrorControlFunction(funcion+"Unhandled Error!", fmt.Sprint(http.StatusInternalServerError))

	// PARTE "0": Buffers, para evitar repetir consultas...
	var hists []*models.HistoricoActa
	Terceros := make(map[int]interface{})
	Ubicaciones := make(map[int]interface{})

	consultasTerceros := 0
	consultasUbicaciones := 0
	consultasProveedores := 0
	evTerceros := 0
	evUbicaciones := 0
	evProveedores := 0

	// PARTE 1 - Identificar los tipos de actas que hay que traer
	// (y así definir la estrategia para traer las actas)
	verTodasLasActas := false
	algunosEstados := []string{}
	proveedor := false
	contratista := false
	idTercero := 0

	// De especificarse un usuario, hay que definir las actas que puede ver
	if usrWSO2 != "" {

		// Traer la información de Autenticación MID para obtener los roles
		var usr models.UsuarioAutenticacion
		if data, err := autenticacion.DataUsuario(usrWSO2); err == nil && data.Role != nil && len(data.Role) > 0 {
			// logs.Debug(data)
			usr = data
		} else if err != nil {
			// formatdata.JsonPrint(data)
			return nil, err
		} else { // data.Role == nil || len(data.Role) == 0
			err := fmt.Errorf("el usuario '%s' no está registrado en WSO2 y/o no tiene roles asignados", usrWSO2)
			logs.Warn(err)
			outputError = e.Error(funcion+"autenticacion.DataUsuario(usrWSO2)", err, fmt.Sprint(http.StatusNotFound))
			return nil, outputError
		}

		// Averiguar si el usuario puede ver todas las actas en todos los estados
		for _, rol := range usr.Role {
			if verTodasLasActas {
				break
			}
			for _, rolSuficiente := range verCualquierEstado {
				if rol == rolSuficiente {
					verTodasLasActas = true
					break
				}
			}
		}

		// Si no puede ver actas en cualquier estado, averiguar en qué estados puede ver
		if !verTodasLasActas {
			for estado, roles := range reglasVerTodas {
				verEstado := false
				for _, rolSuficiente := range roles {
					if verEstado {
						break
					}
					for _, rol := range usr.Role {
						if rol == rolSuficiente {
							verEstado = true
							break
						}
					}
				}
				if verEstado {
					algunosEstados = append(algunosEstados, estado)
				}
			}
		}

		// Si no puede ver todas las actas de al menos un estado, únicamente se
		// traerán las asignadas como contratista o proveedor
		if len(algunosEstados) == 0 {
			for _, rol := range usr.Role {
				if proveedor && contratista {
					break
				}
				if rol == models.RolesArka["Proveedor"] {
					proveedor = true
				} else if rol == models.RolesArka["Contratista"] {
					contratista = true
				}
			}
			if proveedor || contratista {
				// fmt.Println(usr.Documento)
				if data, err := crudTerceros.GetTerceroByDoc(usr.Documento); err == nil {
					// fmt.Println(data.TerceroId.Id)
					if data.TerceroId != nil {
						idTercero = data.TerceroId.Id
					} else {
						return nil, err
					}
				} else {
					return nil, err
				}
			}
		}
	}
	logs.Debug("u:", usrWSO2, "- t:", verTodasLasActas, "- e:", algunosEstados, "- p:", proveedor, "- c:", contratista, "- i:", idTercero)
	logs.Debug("Estados Solicitados:", states)

	// Si se pasaron estados
	if len(states) > 0 {
		if usrWSO2 == "" || verTodasLasActas {
			algunosEstados = states
			verTodasLasActas = false
		} else if idTercero == 0 { // len(algunosEstados) > 0
			estFinales := []string{}
			for _, estUsuario := range algunosEstados {
				for _, est := range states {
					if est == estUsuario {
						estFinales = append(estFinales, estUsuario)
						break
					}
				}
			}
			algunosEstados = estFinales
		}
		logs.Debug("t:", verTodasLasActas, "- e:", algunosEstados)
	}

	// PARTE 2: Traer los tipos de actas identificados
	// (con base a la estrategia definida anteriormente)
	limitStr := fmt.Sprint(limit)
	offsetStr := fmt.Sprint(offset)
	const (
		sortby = "ActaRecibidoId__Id"
		order  = "desc"
	)
	query := "Activo:true,ActaRecibidoId__TipoActaId__Nombre__in:Regular|Especial"
	if verTodasLasActas {
	} else if len(algunosEstados) > 0 {
		query += ",EstadoActaId__Nombre__in:" + strings.Join(algunosEstados, "|")
	} else if contratista || proveedor {
		query += ",EstadoActaId__Nombre"
		if contratista {
			query += "__in:En Elaboracion|En Modificacion"
			query += ",PersonaAsignadaId:" + fmt.Sprint(idTercero)
		} else if proveedor {
			query += ":En Elaboracion"
			query += ",ProveedorId:" + fmt.Sprint(idTercero)
		}
	}
	if hists, outputError = actaRecibido.GetAllHistoricoActa(query, "", sortby, order, offsetStr, limitStr); outputError != nil {
		return
	}

	// PARTE 3: Completar data faltante
	for _, historicos := range hists {

		var ubicacionData map[string]interface{}
		var editor *models.Tercero
		var preUbicacion map[string]interface{}
		var asignado *models.Tercero

		preUbicacion = nil

		reqTercero := func(id int) func() (interface{}, map[string]interface{}) {
			return func() (interface{}, map[string]interface{}) {
				return crudTerceros.GetTerceroById(id)
			}
		}
		idRev := historicos.RevisorId
		if v, err := utilsHelper.BufferGeneric(idRev, Terceros, reqTercero(idRev), &consultasTerceros, &evTerceros); err == nil {
			if v2, ok := v.(*models.Tercero); ok {
				editor = v2
			}
		}

		idUb := historicos.UbicacionId
		reqUbicacion := func() (interface{}, map[string]interface{}) {
			return oikos.GetAsignacionSedeDependencia(fmt.Sprint(idUb))
		}
		if v, err := utilsHelper.BufferGeneric(idUb, Ubicaciones, reqUbicacion, &consultasUbicaciones, &evUbicaciones); err == nil {
			if v2, ok := v.(map[string]interface{}); ok {
				preUbicacion = v2
			}
		}

		idAsignado := historicos.PersonaAsignadaId
		if v, err := utilsHelper.BufferGeneric(idAsignado, Terceros, reqTercero(idAsignado), &consultasTerceros, &evTerceros); err == nil {
			if v2, ok := v.(*models.Tercero); ok {
				asignado = v2
			}
		}

		if v, ok := preUbicacion["EspacioFisicoId"]; ok {
			if err := formatdata.FillStruct(v, &ubicacionData); err != nil {
				logs.Error(err)
				outputError = e.Error(funcion+"error al obtener información del espacio fisico", err, fmt.Sprint(http.StatusBadGateway))
				return
			}
		} else {
			ubicacionData = map[string]interface{}{
				"Nombre": "",
			}
		}

		fVistoBueno := ""
		if historicos.FechaVistoBueno.After(zero) {
			fVistoBueno = historicos.FechaVistoBueno.Format(FormatoFecha)
		}

		Acta := models.ActaResumen{
			Id:                historicos.ActaRecibidoId.Id,
			UbicacionId:       ubicacionData["Nombre"].(string),
			FechaCreacion:     historicos.ActaRecibidoId.FechaCreacion,
			FechaVistoBueno:   fVistoBueno,
			FechaModificacion: historicos.FechaModificacion,
			Observaciones:     historicos.Observaciones,
			RevisorId:         editor.NombreCompleto,
			PersonaAsignada:   asignado.NombreCompleto,
			Estado:            historicos.EstadoActaId.Nombre,
			EstadoActaId:      (*models.EstadoActa)(historicos.EstadoActaId),
		}

		historicoActa = append(historicoActa, Acta)
	}
	logs.Debug(map[string]interface{}{
		"consultasTerceros":    consultasTerceros,
		"evTerceros":           evTerceros,
		"consultasUbicaciones": consultasUbicaciones,
		"evUbicaciones":        evUbicaciones,
		"consultasProveedores": consultasProveedores,
		"evProveedores":        evProveedores,
		"actas":                len(historicoActa),
	})
	return
}
