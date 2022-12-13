package actaRecibido

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/crud/oikos"
	crudTerceros "github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/mid/autenticacion"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

// GetAllActasRecibido ...
func GetAllActasRecibidoActivas(states []string, usrWSO2 string, limit int64, offset int64) (historicoActa []map[string]interface{}, outputError map[string]interface{}) {

	funcion := "GetAllActasRecibidoActivas - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	// PARTE "0": Buffers, para evitar repetir consultas...
	var Historico []map[string]interface{}
	Terceros := make(map[int]interface{})
	Ubicaciones := make(map[int]models.AsignacionEspacioFisicoDependencia)

	consultasTerceros := 0
	consultasProveedores := 0
	evTerceros := 0
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
	urlEstados := "http://" + beego.AppConfig.String("actaRecibidoService") + "historico_acta?limit=" + fmt.Sprint(limit) + "&offset=" + fmt.Sprint(offset)
	urlEstados += "&sortby=ActaRecibidoId__Id&order=desc&query=Activo:true,ActaRecibidoId__TipoActaId__Nombre__in:Regular|Especial"
	if verTodasLasActas {
		var hists []map[string]interface{}
		if resp, err := request.GetJsonTest(urlEstados, &hists); err == nil && resp.StatusCode == 200 {
			if len(hists) == 0 || len(hists[0]) == 0 {
				return nil, nil
			}
			Historico = append(Historico, hists...)
		} else {
			if err == nil {
				err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
			}
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": funcion + "request.GetJsonTest(urlTodas, &hists)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}

	} else if len(algunosEstados) > 0 {
		estados := strings.Join(algunosEstados, "|")
		urlEstado := urlEstados + ",EstadoActaId__Nombre__in:" + estados
		urlEstado = strings.ReplaceAll(urlEstado, " ", "%20")
		var hists = make([]map[string]interface{}, 0)
		if resp, err := request.GetJsonTest(urlEstado, &hists); err != nil || resp.StatusCode != 200 {
			if err == nil {
				err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
			}
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": funcion + "request.GetJsonTest(urlEstado, &hists)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
		Historico = hists
	} else if contratista || proveedor {

		urlEstados += ",EstadoActaId__Nombre"
		if contratista {
			urlEstados += "__in:" + url.QueryEscape("En Elaboracion|En Modificacion")
			urlEstados += ",PersonaAsignadaId:" + fmt.Sprint(idTercero)
		} else if proveedor {
			urlEstados += ":" + url.QueryEscape("En Elaboracion")
			urlEstados += ",ProveedorId:" + fmt.Sprint(idTercero)
		}

		var hists []map[string]interface{}
		if resp, err := request.GetJsonTest(urlEstados, &hists); err == nil && resp.StatusCode == 200 {
			if len(hists) == 0 || len(hists[0]) == 0 {
				return nil, nil
			}
			Historico = append(Historico, hists...)
		} else {
			if err == nil {
				err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
			}
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": funcion + "request.GetJsonTest(urlContProv, &hists)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}

	}

	// PARTE 3: Completar data faltante
	if len(Historico) > 0 {

		for _, historicos := range Historico {

			var acta map[string]interface{}
			var estado map[string]interface{}
			var editor *models.Tercero
			var asignado *models.Tercero

			if data, err := utilsHelper.ConvertirInterfaceMap(historicos["ActaRecibidoId"]); err == nil {
				acta = data
			} else {
				return nil, err
			}

			if data, err := utilsHelper.ConvertirInterfaceMap(historicos["EstadoActaId"]); err == nil {
				estado = data
			} else {
				return nil, err
			}

			reqTercero := func(id int) func() (interface{}, map[string]interface{}) {
				return func() (interface{}, map[string]interface{}) {
					if Tercero, err := crudTerceros.GetTerceroById(id); err == nil {
						return Tercero, nil
					} else {
						return nil, err
					}
				}
			}

			idRev := int(historicos["RevisorId"].(float64))
			if v, err := utilsHelper.BufferGeneric(idRev, Terceros, reqTercero(idRev), &consultasTerceros, &evTerceros); err == nil {
				if v2, ok := v.(*models.Tercero); ok {
					editor = v2
				}
			}

			if _, ok := Ubicaciones[int(historicos["UbicacionId"].(float64))]; int(historicos["UbicacionId"].(float64)) > 0 && !ok {
				id_ := strconv.Itoa(int(historicos["UbicacionId"].(float64)))
				if asignacion, err := oikos.GetAllAsignacion("query=Id:" + id_); err != nil {
					return nil, err
				} else if len(asignacion) == 1 {
					Ubicaciones[int(historicos["UbicacionId"].(float64))] = asignacion[0]
				}
			}

			idAsignado := int(historicos["PersonaAsignadaId"].(float64))
			if v, err := utilsHelper.BufferGeneric(idAsignado, Terceros, reqTercero(idAsignado), &consultasTerceros, &evTerceros); err == nil {
				if v2, ok := v.(*models.Tercero); ok {
					asignado = v2
				}
			}

			fVistoBueno := historicos["FechaVistoBueno"].(string)
			if fVistoBueno == "0001-01-01T00:00:00Z" {
				fVistoBueno = ""
			}

			Acta := map[string]interface{}{
				"Id":                acta["Id"],
				"UbicacionId":       "",
				"FechaCreacion":     acta["FechaCreacion"],
				"FechaVistoBueno":   fVistoBueno,
				"FechaModificacion": historicos["FechaModificacion"],
				"Observaciones":     historicos["Observaciones"],
				"RevisorId":         editor.NombreCompleto,
				"PersonaAsignada":   asignado.NombreCompleto,
				"Estado":            estado["Nombre"],
				"EstadoActaId":      estado,
			}

			if val, ok := Ubicaciones[int(historicos["UbicacionId"].(float64))]; ok {
				Acta["UbicacionId"] = val.EspacioFisicoId.Nombre
			}

			historicoActa = append(historicoActa, Acta)
		}

		logs.Info("consultasTerceros:", consultasTerceros, " - Evitadas: ", evTerceros)
		logs.Info("consultasProveedores:", consultasProveedores, " - Evitadas: ", evProveedores)
		logs.Info(len(historicoActa), "actas")
		return historicoActa, nil

	} else {
		return nil, nil
	}
}
