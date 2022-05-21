package actaRecibido

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	// "github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/arka_mid/helpers/crud/oikos"
	crudTerceros "github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/mid/autenticacion"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	e "github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

// GetAllActasRecibido ...
func GetAllActasRecibidoActivas(states []string, usrWSO2 string, limit, offset int) (historicoActa []map[string]interface{}, outputError map[string]interface{}) {

	const funcion = "GetAllActasRecibidoActivas - "
	defer e.ErrorControlFunction(funcion+"Unhandled Error!", fmt.Sprint(http.StatusInternalServerError))

	// PARTE "0": Buffers, para evitar repetir consultas...
	var hists []map[string]interface{}
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
					idTercero = data.TerceroId.Id
				} else {
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
	base := "http://" + beego.AppConfig.String("actaRecibidoService") + "historico_acta?"
	params := url.Values{}
	params.Add("limit", fmt.Sprint(limit))
	params.Add("offset", fmt.Sprint(offset))
	params.Add("sortby", "ActaRecibidoId__Id")
	params.Add("order", "desc")
	query := "Activo:true,ActaRecibidoId__TipoActaId__Nombre__in:Regular|Especial"
	if verTodasLasActas {
		params.Add("query", query)
		urlTodas := base + params.Encode()
		logs.Debug("urlTodas:", urlTodas)
		if resp, err := request.GetJsonTest(urlTodas, &hists); err == nil && resp.StatusCode == 200 {
			if len(hists) == 0 || len(hists[0]) == 0 {
				return nil, nil
			}
		} else {
			if err == nil {
				err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
			}
			logs.Error(err)
			outputError = e.Error(funcion+"request.GetJsonTest(urlTodas, &hists)", err, fmt.Sprint(http.StatusBadGateway))
			return nil, outputError
		}

	} else if len(algunosEstados) > 0 {
		query += ",EstadoActaId__Nombre__in:" + strings.Join(algunosEstados, "|")
		params.Add("query", query)
		urlEstados := base + params.Encode()
		logs.Debug("urlEstados:", urlEstados)
		if resp, err := request.GetJsonTest(urlEstados, &hists); err == nil && resp.StatusCode == 200 {
			if len(hists) == 0 || len(hists[0]) == 0 {
				return nil, nil
			}
		} else {
			if err == nil {
				err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
			}
			logs.Error(err)
			outputError = e.Error(funcion+"request.GetJsonTest(urlEstados, &hists)", err, fmt.Sprint(http.StatusBadGateway))
			return nil, outputError
		}

	} else if contratista || proveedor {
		query += ",EstadoActaId__Nombre"
		if contratista {
			query += "__in:En Elaboracion|En Modificacion"
			query += ",PersonaAsignadaId:" + fmt.Sprint(idTercero)
		} else if proveedor {
			query += ":En Elaboracion"
			query += ",ProveedorId:" + fmt.Sprint(idTercero)
		}
		params.Add("query", query)

		urlContProv := base + params.Encode()
		logs.Debug("urlContProv:", urlContProv)
		if resp, err := request.GetJsonTest(urlContProv, &hists); err == nil && resp.StatusCode == 200 {
			if len(hists) == 0 || len(hists[0]) == 0 {
				return nil, nil
			}
		} else {
			if err == nil {
				err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
			}
			logs.Error(err)
			outputError = e.Error(funcion+"request.GetJsonTest(urlContProv, &hists)", err, fmt.Sprint(http.StatusBadGateway))
			return nil, outputError
		}

	}

	// PARTE 3: Completar data faltante
	if len(hists) > 0 {

		for _, historicos := range hists {

			var acta map[string]interface{}
			var estado map[string]interface{}
			var ubicacionData map[string]interface{}
			var editor *models.Tercero
			var preUbicacion map[string]interface{}
			var asignado *models.Tercero

			preUbicacion = nil

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
					return crudTerceros.GetTerceroById(id)
				}
			}

			idRev := int(historicos["RevisorId"].(float64))
			if v, err := utilsHelper.BufferGeneric(idRev, Terceros, reqTercero(idRev), &consultasTerceros, &evTerceros); err == nil {
				if v2, ok := v.(*models.Tercero); ok {
					editor = v2
				}
			}

			idUbStr := strconv.Itoa(int(historicos["UbicacionId"].(float64)))
			reqUbicacion := func() (interface{}, map[string]interface{}) {
				return oikos.GetAsignacionSedeDependencia(idUbStr)
			}
			if idUb, err := strconv.Atoi(idUbStr); err == nil {
				if v, err := utilsHelper.BufferGeneric(idUb, Ubicaciones, reqUbicacion, &consultasUbicaciones, &evUbicaciones); err == nil {
					if v2, ok := v.(map[string]interface{}); ok {
						preUbicacion = v2
					}
				}
			}

			idAsignado := int(historicos["PersonaAsignadaId"].(float64))
			if v, err := utilsHelper.BufferGeneric(idAsignado, Terceros, reqTercero(idAsignado), &consultasTerceros, &evTerceros); err == nil {
				if v2, ok := v.(*models.Tercero); ok {
					asignado = v2
				}
			}

			if v, ok := preUbicacion["EspacioFisicoId"]; ok {
				if err := formatdata.FillStruct(v, &ubicacionData); err != nil {
					logs.Error(err)
					outputError = e.Error(funcion+"error al obtener información del espacio fisico", err, fmt.Sprint(http.StatusBadGateway))
					return nil, outputError
				}
			} else {
				ubicacionData = map[string]interface{}{
					"Nombre": "",
				}
			}

			fVistoBueno := historicos["FechaVistoBueno"].(string)
			if fVistoBueno == "0001-01-01T00:00:00Z" {
				fVistoBueno = ""
			}

			Acta := map[string]interface{}{
				"Id":                acta["Id"],
				"UbicacionId":       ubicacionData["Nombre"],
				"FechaCreacion":     acta["FechaCreacion"],
				"FechaVistoBueno":   fVistoBueno,
				"FechaModificacion": historicos["FechaModificacion"],
				"Observaciones":     historicos["Observaciones"],
				"RevisorId":         editor.NombreCompleto,
				"PersonaAsignada":   asignado.NombreCompleto,
				"Estado":            estado["Nombre"],
				"EstadoActaId":      estado,
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
		return historicoActa, nil

	} else {
		return nil, nil
	}
}
