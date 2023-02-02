package actaRecibido

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/oikos"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/mid/autenticacion"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// GetAllActasRecibidoActivas ...
func GetAllActasRecibidoActivas(usrWSO2 string,
	id_, tipos string, estados []string, fechaCreacion_, fechaModificacion_ string,
	sortby, order string, limit int64, offset int64) (
	historicoActa []map[string]interface{}, count string, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GetAllActasRecibidoActivas - Unhandled Error!", "500")

	// PARTE 1 - Identificar los tipos de actas que hay que traer

	verTodasLasActas, algunosEstados, user, outputError := getEstados(estados, usrWSO2)
	if outputError != nil {
		return nil, "", outputError
	}

	if !verTodasLasActas && len(algunosEstados) == 0 {
		return
	}

	proveedor, contratista, idTercero, outputError := getTereroId(verTodasLasActas, algunosEstados, user)
	if outputError != nil {
		return nil, "", outputError
	}

	// PARTE 2: Traer los tipos de actas identificados
	query := "Activo:true"
	if len(algunosEstados) != 0 {
		query += ",EstadoActaId__CodigoAbreviacion__in:" + strings.Join(algunosEstados, "|")
	}

	if contratista || proveedor {
		if contratista {
			query += ",PersonaAsignadaId:" + fmt.Sprint(idTercero)
		} else if proveedor {
			query += ",ProveedorId:" + fmt.Sprint(idTercero)
		}
	}

	if limit > 0 && offset > 0 {
		offset--
		offset *= limit
	}

	if id_ != "" {
		query += ",ActaRecibidoId__Id__icontains:" + id_
	}

	if fechaCreacion_ != "" {
		fechaCreacion_ = strings.ReplaceAll(fechaCreacion_, "/", "-")
		query += ",ActaRecibidoId__FechaCreacion__icontains:" + fechaCreacion_
	}

	if fechaModificacion_ != "" {
		fechaModificacion_ = strings.ReplaceAll(fechaModificacion_, "/", "-")
		query += ",FechaModificacion__icontains:" + fechaModificacion_
	}

	query += ",ActaRecibidoId__TipoActaId__Id__in:"
	if tipos != "" {
		query += tipos
	} else {
		query += "1|2"
	}

	if order != "" && (sortby == "Id" || sortby == "FechaCreacion" || sortby == "FechaModificacion" || sortby == "FechaVistoBueno" || sortby == "EstadoActaId") {
		if sortby == "FechaCreacion" {
			sortby = "ActaRecibidoId__FechaCreacion"
		} else if sortby == "Id" {
			sortby = "ActaRecibidoId__Id"
		} else if sortby == "EstadoActaId" {
			sortby = "EstadoActaId__Nombre"
		}
		order = strings.ToLower(order)
	} else {
		sortby = "ActaRecibidoId__Id"
		order = "desc"
	}

	historicos, count, err := actaRecibido.GetAllHistoricoActas(query, "", sortby, order, fmt.Sprint(offset), fmt.Sprint(limit))
	if err != nil {
		return nil, "", err
	}

	// PARTE 3: Completar data faltante

	Terceros := make(map[int]models.Tercero)
	Ubicaciones := make(map[int]models.AsignacionEspacioFisicoDependencia)
	for _, historico := range historicos {

		var editor models.Tercero
		var asignado models.Tercero

		if historico.RevisorId > 0 {
			if val, ok := Terceros[historico.RevisorId]; !ok {
				if revisor, err := terceros.GetTerceroById(historico.RevisorId); err != nil {
					return nil, "", err
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
					return nil, "", err
				} else if len(asignacion) == 1 {
					Ubicaciones[historico.UbicacionId] = asignacion[0]
				}
			}
		}

		if historico.PersonaAsignadaId > 0 {
			if val, ok := Terceros[historico.PersonaAsignadaId]; !ok {
				if revisor, err := terceros.GetTerceroById(historico.PersonaAsignadaId); err != nil {
					return nil, "", err
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
			"FechaCreacion":     historico.ActaRecibidoId.FechaCreacion,
			"FechaVistoBueno":   historico.FechaVistoBueno,
			"FechaModificacion": historico.FechaModificacion,
			"Observaciones":     historico.Observaciones,
			"RevisorId":         editor.NombreCompleto,
			"PersonaAsignada":   asignado.NombreCompleto,
			"EstadoActaId":      historico.EstadoActaId.Id,
		}

		if historico.EstadoActaId.CodigoAbreviacion == "Aceptada" || historico.EstadoActaId.CodigoAbreviacion == "AsociadoEntrada" {
			Acta["AceptadaPor"] = editor.NombreCompleto
		}

		if val, ok := Ubicaciones[historico.UbicacionId]; ok && val.EspacioFisicoId != nil {
			Acta["DependenciaId"] = val.DependenciaId.Nombre
		}

		historicoActa = append(historicoActa, Acta)
	}

	return

}

func getEstados(estados []string, user string) (verTodas bool, estados_ []string, usr models.UsuarioAutenticacion, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("getEstados - Unhandled Error!", "500")

	if user != "" {
		// Consulta de roles
		usr, outputError = autenticacion.DataUsuario(user)
		if outputError != nil || usr.Role == nil || len(usr.Role) == 0 {
			return
		}

		for _, rol := range usr.Role {
			if verTodas {
				break
			}
			for _, rolSuficiente := range verCualquierEstado {
				if rol == rolSuficiente {
					verTodas = true
					break
				}
			}
		}

		// Si no puede ver actas en cualquier estado, averiguar en qué estados puede ver
		if !verTodas {
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

				if !verEstado {
					continue
				} else if len(estados) == 0 {
					estados_ = append(estados_, estado)
					continue
				}

				for _, st := range estados {
					if estado == st {
						estados_ = append(estados_, estado)
						break
					}
				}
			}
		} else if len(estados) > 0 {
			estados_ = estados
		}

	} else if len(estados) > 0 {
		estados_ = estados
	} else {
		verTodas = true
	}

	return
}

func getTereroId(verTodas bool, estados []string, usr models.UsuarioAutenticacion) (proveedor, contratista bool, tercero int, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("getTereroId - Unhandled Error!", "500")

	if verTodas {
		return
	}

	for _, rol := range usr.Role {
		if rol == models.RolesArka["Contratista"] {
			contratista = true
			break
		} else if rol == models.RolesArka["Proveedor"] {
			proveedor = true
			break
		}
	}

	if proveedor || contratista {
		outputError = autenticacion.GetTerceroUser(usr, &tercero)
		if outputError != nil {
			return
		} else if tercero == 0 {
			contratista = false
			proveedor = false
		}
	}

	return
}
