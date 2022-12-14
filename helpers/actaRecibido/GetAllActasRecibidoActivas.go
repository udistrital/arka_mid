package actaRecibido

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/oikos"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/mid/autenticacion"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// GetAllActasRecibido ...
func GetAllActasRecibidoActivas(states []string, usrWSO2 string, limit int64, offset int64) (historicoActa []map[string]interface{}, outputError map[string]interface{}) {

	funcion := "GetAllActasRecibidoActivas - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	// PARTE "0": Buffers, para evitar repetir consultas...
	Terceros := make(map[int]models.Tercero)
	Ubicaciones := make(map[int]models.AsignacionEspacioFisicoDependencia)

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
			outputError = map[string]interface{}{
				"funcion": funcion + "autenticacion.DataUsuario(usrWSO2)",
				"err":     err,
				"status":  "404",
			}
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
				err := autenticacion.GetTerceroUser(usr, &idTercero)
				if err != nil || idTercero == 0 {
					return nil, err
				}
			}
		}
	}
	logs.Info("u:", usrWSO2, "- t:", verTodasLasActas, "- e:", algunosEstados, "- p:", proveedor, "- c:", contratista, "- i:", idTercero)

	// fmt.Print("Estados Solicitados: ")
	// fmt.Println(states)

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
		logs.Info("t:", verTodasLasActas, "- e:", algunosEstados)
	}

	// PARTE 2: Traer los tipos de actas identificados
	// (con base a la estrategia definida anteriormente)

	// TODO: Por rendimiento, TODO lo relacionado a ...
	// - buscar el historico_acta mas reciente
	// - Filtrar por estados
	// ... debería moverse a una o más función(es) y/o controlador(es) del CRUD
	query := "Activo:true,ActaRecibidoId__TipoActaId__Nombre__in:Regular|Especial"
	if len(algunosEstados) > 0 {
		query += ",EstadoActaId__Nombre__in:" + strings.Join(algunosEstados, "|")
	} else if contratista || proveedor {
		query += ",EstadoActaId__Nombre"
		if contratista {
			query += "__in:En Elaboracion|En Modificacion,PersonaAsignadaId:" + fmt.Sprint(idTercero)
		} else if proveedor {
			query += ":En Elaboracion,ProveedorId:" + fmt.Sprint(idTercero)
		}
	}

	historicos, err := actaRecibido.GetAllHistoricoActa(query, "", "ActaRecibidoId__Id", "desc", fmt.Sprint(offset), fmt.Sprint(limit))
	if err != nil {
		return nil, err
	}

	// PARTE 3: Completar data faltante
	for _, historico := range historicos {

		var editor models.Tercero
		var asignado models.Tercero

		if historico.RevisorId > 0 {
			if val, ok := Terceros[historico.RevisorId]; !ok {
				logs.Info("Consulta revisor: ", historico.RevisorId)
				if revisor, err := terceros.GetTerceroById(historico.RevisorId); err != nil {
					return nil, err
				} else if revisor != nil {
					editor = *revisor
					Terceros[historico.RevisorId] = *revisor
				}
			} else {
				editor = val
			}
		}

		if historico.UbicacionId > 0 {
			if _, ok := Ubicaciones[historico.UbicacionId]; !ok {
				id_ := strconv.Itoa(historico.UbicacionId)
				if asignacion, err := oikos.GetAllAsignacion("query=Id:" + id_); err != nil {
					return nil, err
				} else if len(asignacion) == 1 {
					Ubicaciones[historico.UbicacionId] = asignacion[0]
				}
			}
		}

		if historico.PersonaAsignadaId > 0 {
			if val, ok := Terceros[historico.PersonaAsignadaId]; !ok {
				if revisor, err := terceros.GetTerceroById(historico.PersonaAsignadaId); err != nil {
					return nil, err
				} else if revisor != nil {
					asignado = *revisor
					Terceros[historico.PersonaAsignadaId] = *revisor
				}
			} else {
				asignado = val
			}
		}

		Acta := map[string]interface{}{
			"Id":                historico.ActaRecibidoId.Id,
			"UbicacionId":       "",
			"FechaCreacion":     historico.FechaCreacion,
			"FechaVistoBueno":   historico.FechaVistoBueno,
			"FechaModificacion": historico.FechaModificacion,
			"Observaciones":     historico.Observaciones,
			"RevisorId":         editor.NombreCompleto,
			"PersonaAsignada":   asignado.NombreCompleto,
			"Estado":            historico.EstadoActaId.Nombre,
			"EstadoActaId":      historico.EstadoActaId,
		}

		if val, ok := Ubicaciones[historico.UbicacionId]; ok && val.EspacioFisicoId != nil {
			Acta["UbicacionId"] = val.EspacioFisicoId.Nombre
		}

		historicoActa = append(historicoActa, Acta)
	}

	logs.Info(len(historicoActa), "actas")

	return historicoActa, nil

}
