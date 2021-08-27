package actaRecibido

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/tealeg/xlsx"

	"github.com/udistrital/arka_mid/helpers/autenticacion"
	"github.com/udistrital/arka_mid/helpers/proveedorHelper"
	"github.com/udistrital/arka_mid/helpers/tercerosHelper"
	"github.com/udistrital/arka_mid/helpers/ubicacionHelper"
	"github.com/udistrital/arka_mid/helpers/unidadHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"

	// "github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

// GetAllActasRecibido ...
func GetAllActasRecibidoActivas(states []string, usrWSO2 string) (historicoActa []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetAllActasRecibidoActivas - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	// PARTE "0": Buffers, para evitar repetir consultas...
	var Historico []map[string]interface{}
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
			outputError = map[string]interface{}{
				"funcion": "GetAllActasRecibidoActivas - autenticacion.DataUsuario(usrWSO2)",
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
				if data, err := tercerosHelper.GetTerceroByDoc(usr.Documento); err == nil {
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
	urlEstados := "http://" + beego.AppConfig.String("actaRecibidoService") + "historico_acta?limit=-1"
	urlEstados += "&fields=ActaRecibidoId,EstadoActaId&query=Activo:true"
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
				"funcion": "GetAllActasRecibidoActivas - request.GetJsonTest(urlTodas, &hists)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}

	} else if len(algunosEstados) > 0 {
		for _, estado := range algunosEstados {
			var hists []map[string]interface{}
			urlEstado := urlEstados + ",EstadoActaId__Nombre:" + estado
			urlEstado = strings.ReplaceAll(urlEstado, " ", "%20")
			if resp, err := request.GetJsonTest(urlEstado, &hists); err == nil && resp.StatusCode == 200 {
				if len(hists) == 0 || len(hists[0]) == 0 {
					continue
				}
				Historico = append(Historico, hists...)
			} else {
				if err == nil {
					err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
				}
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "GetAllActasRecibidoActivas - request.GetJsonTest(urlEstado, &hists)",
					"err":     err,
					"status":  "502",
				}
				return nil, outputError
			}
		}

	} else if contratista || proveedor {

		histMap := make(map[int](map[string]interface{})) // mapeo "idActa --> historico_acta activo"

		var estados []string
		if contratista {
			estados = append(estados, "En Elaboracion", "En Modificacion")
		} else if proveedor {
			estados = append(estados, "En Elaboracion")
		}

		for _, estado := range estados {
			var hists []map[string]interface{}
			urlContProv := urlEstados + ",EstadoActaId__Nombre:" + estado
			if !proveedor {
				// Si no es proveedor, agregar de una vez el filtro del contratista
				// pues sería la única razón para que se ejecute este "for"
				urlContProv += ",ActaRecibidoId__PersonaAsignada:" + fmt.Sprint(idTercero)
			}
			urlContProv = strings.ReplaceAll(urlContProv, " ", "%20")
			// logs.Debug("urlContProv:", urlContProv, "- estado:", estado)
			if resp, err := request.GetJsonTest(urlContProv, &hists); err == nil && resp.StatusCode == 200 {
				if len(hists) == 0 || len(hists[0]) == 0 {
					continue
				}
				for _, hist := range hists {
					idActa := int(hist["ActaRecibidoId"].(map[string]interface{})["Id"].(float64))
					histMap[idActa] = hist
				}
			} else {
				if err == nil {
					err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
				}
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "GetAllActasRecibidoActivas - request.GetJsonTest(urlContProv, &hists)",
					"err":     err,
					"status":  "502",
				}
				return nil, outputError
			}
		}

		for idActa, hist := range histMap {
			agregar := false
			if proveedor {
				var soportes []map[string]interface{}
				urlSoporteActa := "http://" + beego.AppConfig.String("actaRecibidoService") + "soporte_acta"
				urlSoporteActa += "?fields=Id" // Realmente no importan los campos, lo que importa es la asociacion con el proveedor y el acta
				urlSoporteActa += "&query=Activo:true,ActaRecibidoId__Id:" + fmt.Sprint(idActa)
				urlSoporteActa += ",ProveedorId:" + fmt.Sprint(idTercero)
				// logs.Debug("urlSoporteActa:", urlSoporteActa)
				if resp, err := request.GetJsonTest(urlSoporteActa, &soportes); err == nil && resp.StatusCode == 200 {
					if len(soportes) >= 1 {
						for _, soporte := range soportes {
							if len(soporte) > 0 {
								agregar = true
								break
							}
						}
					}
				} else {
					if err == nil {
						err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
					}
					logs.Error(err)
					outputError = map[string]interface{}{
						"funcion": "GetAllActasRecibidoActivas - request.GetJsonTest(urlSoporteActa, &soportes)",
						"err":     err,
						"status":  "502",
					}
					return nil, outputError
				}
			}
			if !agregar && contratista {
				if proveedor {
					if idTercero == int(hist["ActaRecibidoId"].(map[string]interface{})["PersonaAsignada"].(float64)) {
						agregar = true
					}
				} else {
					// Las actas (historicos activos) ya se trajeron filtradas por contratista
					agregar = true
				}
			}
			if agregar {
				Historico = append(Historico, hist)
			}
		}

	}

	// PARTE 3: Completar data faltante
	if len(Historico) > 0 {

		for _, historicos := range Historico {

			var acta map[string]interface{}
			var estado map[string]interface{}
			var ubicacionData map[string]interface{}
			var editor map[string]interface{}
			var preUbicacion map[string]interface{}
			// var oldAsignado *models.Proveedor // "old": de proveedores, se va a eliminar
			var asignado map[string]interface{}

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

			reqTercero := func(id string) func() (interface{}, map[string]interface{}) {
				return func() (interface{}, map[string]interface{}) {
					if Tercero, err := tercerosHelper.GetNombreTerceroById(id); err == nil {
						return Tercero, nil
					} else {
						return nil, err
					}
				}
			}

			idRevStr := fmt.Sprintf("%v", acta["RevisorId"])
			if idRev, err := strconv.Atoi(idRevStr); err == nil {
				if v, err := utilsHelper.BufferGeneric(idRev, Terceros, reqTercero(idRevStr), &consultasTerceros, &evTerceros); err == nil {
					if v2, ok := v.(map[string]interface{}); ok {
						editor = v2
					}
				}
			}

			idUbStr := fmt.Sprintf("%v", acta["UbicacionId"])
			reqUbicacion := func() (interface{}, map[string]interface{}) {
				if ubicacion, err := ubicacionHelper.GetAsignacionSedeDependencia(idUbStr); err == nil {
					return ubicacion, nil
				} else {
					logs.Error(err)
					return nil, map[string]interface{}{
						"funcion": "findAndAddUbicacion - ubicacionHelper.GetAsignacionSedeDependencia(idStr)",
						"err":     err,
						"status":  "502",
					}
				}
			}
			if idUb, err := strconv.Atoi(idUbStr); err == nil {
				if v, err := utilsHelper.BufferGeneric(idUb, Ubicaciones, reqUbicacion, &consultasUbicaciones, &evUbicaciones); err == nil {
					if v2, ok := v.(map[string]interface{}); ok {
						preUbicacion = v2
					}
				}
			}

			idAsignadoStr := fmt.Sprintf("%v", acta["PersonaAsignada"])
			if idAsignado, err := strconv.Atoi(idAsignadoStr); err == nil {
				if v, err := utilsHelper.BufferGeneric(idAsignado, Terceros, reqTercero(idAsignadoStr), &consultasTerceros, &evTerceros); err == nil {
					if v2, ok := v.(map[string]interface{}); ok {
						asignado = v2
					}
				}
			}

			if preUbicacion != nil {
				if jsonString2, err := json.Marshal(preUbicacion["EspacioFisicoId"]); err == nil {
					if err2 := json.Unmarshal(jsonString2, &ubicacionData); err2 != nil {
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "/GetAllActasRecibidoActivas",
							"err":     err,
							"status":  "500",
						}
						return nil, outputError
					}
				}
			} else {
				ubicacionData = map[string]interface{}{
					"Nombre": "Ubicacion No Especificada",
				}
			}
			// fmt.Println(ubicacionData)
			Acta := map[string]interface{}{
				"UbicacionId":       ubicacionData["Nombre"],
				"Activo":            acta["Activo"],
				"FechaCreacion":     acta["FechaCreacion"],
				"FechaVistoBueno":   acta["FechaVistoBueno"],
				"FechaModificacion": acta["FechaModificacion"],
				"Id":                acta["Id"],
				"Observaciones":     acta["Observaciones"],
				"RevisorId":         editor["NombreCompleto"],
				// "oldAsignada":       oldAsignado.NomProveedor,
				"PersonaAsignada": asignado["NombreCompleto"],
				// "PersonaAsignadaId": tmpAsignadoId,
				"Estado": estado["Nombre"],
			}
			// fmt.Println("Es esto")
			// fmt.Println(Acta)
			historicoActa = append(historicoActa, Acta)
		}

		logs.Info("consultasTerceros:", consultasTerceros, " - Evitadas: ", evTerceros)
		logs.Info("consultasUbicaciones:", consultasUbicaciones, " - Evitadas: ", evUbicaciones)
		logs.Info("consultasProveedores:", consultasProveedores, " - Evitadas: ", evProveedores)
		// formatdata.JsonPrint(Proveedores)
		logs.Info(len(historicoActa), "actas")
		return historicoActa, nil

	} else {
		return nil, nil
	}
}

func RemoveIndex(s []byte, index int) []byte {
	return append(s[:index], s[index+1:]...)
}

// GetAllParametrosActa ...
func GetAllParametrosActa() (Parametros []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetAllParametrosActa - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var (
		Unidades       interface{}
		TipoBien       interface{}
		EstadoActa     interface{}
		EstadoElemento interface{}
		ss             map[string]interface{}
		Parametro      []interface{}
		Valor          []interface{}
		IvaTest        []Imp
		Ivas           []Imp
	)

	parametros := make([]map[string]interface{}, 0)

	urlActasTipoBien := "http://" + beego.AppConfig.String("actaRecibidoService") + "tipo_bien?limit=-1"
	if _, err := request.GetJsonTest(urlActasTipoBien, &TipoBien); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllParametrosActa - request.GetJsonTest(urlActasTipoBien, &TipoBien)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	urlActasEstadoActa := "http://" + beego.AppConfig.String("actaRecibidoService") + "estado_acta?limit=-1"
	if _, err := request.GetJsonTest(urlActasEstadoActa, &EstadoActa); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllParametrosActa - request.GetJsonTest(urlActasEstadoActa, &EstadoActa)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	urlACtasEstadoElem := "http://" + beego.AppConfig.String("actaRecibidoService") + "estado_elemento?limit=-1"
	if _, err := request.GetJsonTest(urlACtasEstadoElem, &EstadoElemento); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllParametrosActa - request.GetJsonTest(urlACtasEstadoElem, &EstadoElemento)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	urlParametros := "http://" + beego.AppConfig.String("parametrosService") + "parametro_periodo?query=PeriodoId__Nombre:2021,ParametroId__TipoParametroId__Id:12"
	if _, err := request.GetJsonTest(urlParametros, &ss); err == nil {

		var data []map[string]interface{}
		if jsonString, err := json.Marshal(ss["Data"]); err == nil {
			if err := json.Unmarshal(jsonString, &data); err == nil {
				for _, valores := range data {
					Parametro = append(Parametro, valores["ParametroId"])
					v := []byte(fmt.Sprintf("%v", valores["Valor"]))
					var valorUnm interface{}
					if err := json.Unmarshal(v, &valorUnm); err == nil {
						Valor = append(Valor, valorUnm)
					}
				}
			}
		}

		if jsonbody1, err := json.Marshal(Parametro); err == nil {
			if err := json.Unmarshal(jsonbody1, &Ivas); err != nil {
				fmt.Println(err)
				return
			}
		}

		if jsonbody1, err := json.Marshal(Valor); err == nil {
			if err := json.Unmarshal(jsonbody1, &IvaTest); err != nil {
				fmt.Println(err)
				return
			}
		}

		for i, valores := range IvaTest {
			IvaTest[i].CodigoAbreviacion = valores.CodigoAbreviacion
		}
		for i, valores := range Ivas {
			IvaTest[i].BasePesos = valores.BasePesos
			IvaTest[i].BaseUvt = valores.BaseUvt
			IvaTest[i].PorcentajeAplicacion = valores.PorcentajeAplicacion
			IvaTest[i].CodigoAbreviacion = valores.CodigoAbreviacion
		}

	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllParametrosActa - request.GetJsonTest(urlParametros, &ss)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	urlUnidad := "http://" + beego.AppConfig.String("AdministrativaService") + "unidad?limit=-1"
	if _, err := request.GetJsonTest(urlUnidad, &Unidades); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllParametrosActa - request.GetJsonTest(urlUnidad, &Unidades)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	parametros = append(parametros, map[string]interface{}{
		"Unidades":       Unidades,
		"TipoBien":       TipoBien,
		"EstadoActa":     EstadoActa,
		"EstadoElemento": EstadoElemento,
		"IVA":            IvaTest,
	})

	return parametros, nil
}

// "DecodeXlsx2Json ..."
func DecodeXlsx2Json(c multipart.File) (Archivo []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "DecodeXlsx2Json - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var Unidades []Unidad
	var SubgruposConsumo []map[string]interface{}
	var SubgruposConsumoControlado []map[string]interface{}
	var SubgruposDevolutivo []map[string]interface{}
	var (
		ss        map[string]interface{}
		Parametro []interface{}
		Valor     []interface{}
		IvaTest   []Imp
		Ivas      []Imp
	)

	urlIva := "http://" + beego.AppConfig.String("parametrosService") + "parametro_periodo?query=PeriodoId__Nombre:2021,ParametroId__TipoParametroId__Id:12"
	// logs.Debug("urlIva:", urlIva)
	if resp, err := request.GetJsonTest(urlIva, &ss); err == nil && resp.StatusCode == 200 {

		var data []map[string]interface{}
		if jsonString, err := json.Marshal(ss["Data"]); err == nil {
			if err := json.Unmarshal(jsonString, &data); err == nil {
				for _, valores := range data {
					Parametro = append(Parametro, valores["ParametroId"])
					v := []byte(fmt.Sprintf("%v", valores["Valor"]))
					var valorUnm interface{}
					if err := json.Unmarshal(v, &valorUnm); err == nil {
						Valor = append(Valor, valorUnm)
					}
				}
			}
		}

		if jsonbody1, err := json.Marshal(Parametro); err == nil {
			if err := json.Unmarshal(jsonbody1, &Ivas); err != nil {
				fmt.Println(err)
				return
			}
		}

		if jsonbody1, err := json.Marshal(Valor); err == nil {
			if err := json.Unmarshal(jsonbody1, &IvaTest); err != nil {
				fmt.Println(err)
				return
			}
		}

		for i, valores := range IvaTest {
			IvaTest[i].CodigoAbreviacion = valores.CodigoAbreviacion
		}
		for i, valores := range Ivas {
			IvaTest[i].BasePesos = valores.BasePesos
			IvaTest[i].BaseUvt = valores.BaseUvt
			IvaTest[i].PorcentajeAplicacion = valores.PorcentajeAplicacion
			IvaTest[i].CodigoAbreviacion = valores.CodigoAbreviacion
		}

	} else {
		if err == nil {
			err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "DecodeXlsx2Json - request.GetJsonTest(urlIva, &ss)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	urlCatT1 := "http://" + beego.AppConfig.String("catalogoElementosService") + "tr_catalogo/tipo_de_bien/1"
	if _, err := request.GetJsonTest(urlCatT1, &SubgruposConsumo); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "DecodeXlsx2Json - request.GetJsonTest(urlCatT1, &SubgruposConsumo)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	urlCatT2 := "http://" + beego.AppConfig.String("catalogoElementosService") + "tr_catalogo/tipo_de_bien/2"
	if _, err := request.GetJsonTest(urlCatT2, &SubgruposConsumoControlado); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "DecodeXlsx2Json - request.GetJsonTest(urlCatT2, &SubgruposConsumoControlado)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	urlCatT3 := "http://" + beego.AppConfig.String("catalogoElementosService") + "tr_catalogo/tipo_de_bien/3"
	if _, err := request.GetJsonTest(urlCatT3, &SubgruposDevolutivo); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "DecodeXlsx2Json - request.GetJsonTest(urlCatT3, &SubgruposDevolutivo)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	urlAdmistrativa := "http://" + beego.AppConfig.String("AdministrativaService") + "unidad?limit=-1"
	if _, err := request.GetJsonTest(urlAdmistrativa, &Unidades); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "DecodeXlsx2Json - request.GetJsonTest(urlAdmistrativa, &Unidades)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	file, err := ioutil.ReadAll(c)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "DecodeXlsx2Json - ioutil.ReadAll(c)",
			"err":     err,
			"status":  "400",
		}
		return nil, outputError
	}

	xlFile, err := xlsx.OpenBinary(file)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "DecodeXlsx2Json - xlsx.OpenBinary(file)",
			"err":     err,
			"status":  "400",
		}
		return nil, outputError
	}

	Respuesta := make([]map[string]interface{}, 0)
	Elemento := make([]map[string]interface{}, 0)

	var hojas []string
	var campos []string
	var elementos [14]string

	validar_campos := []string{"Nivel Inventarios", "Tipo de Bien", "Subgrupo Catalogo", "Nombre", "Marca", "Serie", "Cantidad", "Unidad de Medida", "Valor Unitario", "Subtotal", "Descuento", "Tipo IVA", "Valor IVA", "Valor Total"}

	for s, sheet := range xlFile.Sheets {

		if s == 0 {
			hojas = append(hojas, sheet.Name)
			for r, row := range sheet.Rows {
				if r == 0 {
					for i, cell := range row.Cells {
						campos = append(campos, cell.String())
						if campos[i] != validar_campos[i] {
							err := fmt.Errorf("el formato no corresponde a las columnas necesarias")
							logs.Error(err)
							outputError = map[string]interface{}{
								"funcion": "DecodeXlsx2Json - campos[i] != validar_campos[i]",
								"err":     err,
								"status":  "400",
							}
							return nil, outputError
						}
					}
				} else {

					for i, cell := range row.Cells {
						elementos[i] = cell.String()
					}
					if elementos[0] != "Totales" {
						vlrcantidad, err := strconv.ParseInt(elementos[6], 10, 64)
						if err == nil {
						} else {
							vlrcantidad = 0
							logs.Warn(err)
						}

						vlrunitario, err := strconv.ParseFloat(elementos[8], 64)
						if err == nil {
						} else {
							vlrunitario = float64(0)
							logs.Warn(err)
						}

						vlrsubtotal := float64(0)
						vlrsubtotal = float64(vlrunitario) * float64(vlrcantidad)
						elementos[9] = strconv.FormatFloat(vlrsubtotal, 'f', 2, 64)

						vlrdcto, err := strconv.ParseFloat(elementos[10], 64)
						if err == nil {
							vlrdcto = vlrsubtotal - vlrdcto
						} else {
							vlrdcto = float64(0)
							logs.Warn(err)
						}

						convertir := strings.Split(elementos[11], ".")
						if err == nil {
							valor, err := strconv.ParseInt(convertir[0], 10, 64)
							if err == nil {
								for _, valor_iva := range IvaTest {
									if valor == int64(valor_iva.Tarifa) {
										elementos[12] = strconv.FormatFloat(vlrdcto*float64(valor)/100, 'f', 2, 64)
										elementos[11] = strconv.Itoa(valor_iva.Tarifa)
									}
								}
							} else {
								logs.Warn(err)
							}
						} else {
							logs.Warn(err)
						}

						vlrtotal, err := strconv.ParseFloat(elementos[12], 64)
						if err == nil {
							vlrtotal = vlrdcto + vlrtotal
							elementos[13] = strconv.FormatFloat(vlrtotal, 'f', 2, 64)
						} else {
							vlrtotal = float64(0)
							logs.Warn(err)
						}

						convertir2 := strings.ToUpper(elementos[7])
						if err == nil {
							for _, unidad := range Unidades {
								if convertir2 == unidad.Unidad {
									elementos[7] = strconv.Itoa(unidad.Id)
								}
							}
						} else {
							logs.Warn(err)
						}

						convertir3 := elementos[2]
						if err == nil {
							logs.Info(convertir3)
							for _, consumo := range SubgruposConsumo {
								if convertir3 == consumo["Nombre"] {
									elementos[2] = fmt.Sprintf("%v", consumo["Id"])
									elementos[1] = strconv.Itoa(1)
								}
							}
							for _, consumoC := range SubgruposConsumoControlado {
								if convertir3 == consumoC["Nombre"] {
									elementos[2] = fmt.Sprintf("%v", consumoC["Id"])
									elementos[1] = strconv.Itoa(2)
								}
							}
							for _, devolutivo := range SubgruposDevolutivo {
								if convertir3 == devolutivo["Nombre"] {
									elementos[2] = fmt.Sprintf("%v", devolutivo["Id"])
									elementos[1] = strconv.Itoa(3)
								}
							}
						} else {
							logs.Warn(err)
						}

						Elemento = append(Elemento, map[string]interface{}{
							"NivelInventariosId": elementos[0],
							"TipoBienId":         elementos[1],
							"SubgrupoCatalogoId": elementos[2],
							"Nombre":             elementos[3],
							"Marca":              elementos[4],
							"Serie":              elementos[5],
							"Cantidad":           elementos[6],
							"UnidadMedida":       elementos[7],
							"ValorUnitario":      elementos[8],
							"Subtotal":           elementos[9],
							"Descuento":          elementos[10],
							"PorcentajeIvaId":    elementos[11],
							"ValorIva":           elementos[12],
							"ValorTotal":         elementos[13],
						})
					} else {
						Respuesta = append(Respuesta, map[string]interface{}{
							"Hoja":      hojas,
							"Campos":    campos,
							"Elementos": Elemento,
						})

					}
				}
			}
		}
	}
	return Respuesta, nil
}

// GetAllParametrosSoporte ...
func GetAllParametrosSoporte() (Parametros []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetAllParametrosSoporte - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var Dependencias interface{}
	var Sedes interface{}
	var Ubicaciones interface{}
	parametros := make([]map[string]interface{}, 0)

	urlOikosDependencia := "http://" + beego.AppConfig.String("oikos2Service") + "dependencia?limit=-1"
	if _, err := request.GetJsonTest(urlOikosDependencia, &Dependencias); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllParametrosSoporte - request.GetJsonTest(urlOikosDependencia, &Dependencias)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	urlOikosAsignacion := "http://" + beego.AppConfig.String("oikos2Service") + "asignacion_espacio_fisico_dependencia?limit=-1"
	if _, err := request.GetJsonTest(urlOikosAsignacion, &Ubicaciones); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllParametrosSoporte - request.GetJsonTest(urlOikosAsignacion, &Ubicaciones)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	urlOikosEspFis := "http://" + beego.AppConfig.String("oikos2Service") + "espacio_fisico?query=TipoEspacioFisicoId.Id:1&limit=-1"
	if _, err := request.GetJsonTest(urlOikosEspFis, &Sedes); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllParametrosSoporte - request.GetJsonTest(urlOikosEspFis, &Sedes)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	parametros = append(parametros, map[string]interface{}{
		"Dependencias": Dependencias,
		"Ubicaciones":  Ubicaciones,
		"Sedes":        Sedes,
	})

	return parametros, nil
}

// GetAsignacionSedeDependencia ...
func GetAsignacionSedeDependencia(Datos models.GetSedeDependencia) (Parametros []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetAsignacionSedeDependencia - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	if Datos.Sede == nil {
		err := fmt.Errorf("sede no especificada")
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAsignacionSedeDependencia - Datos.Sede == nil",
			"err":     err,
			"status":  "400",
		}
		return nil, outputError
	}

	var Ubicaciones []map[string]interface{}
	var Parametros2 []map[string]interface{}
	// logs.Debug("Datos:")
	// formatdata.JsonPrint(Datos)
	// fmt.Println("")
	oikosUrl := "http://" + beego.AppConfig.String("oikos2Service") + "asignacion_espacio_fisico_dependencia?limit=-1"
	oikosUrl += "&query=DependenciaId.Id:" + strconv.Itoa(Datos.Dependencia.Id)
	// logs.Debug("oikosUrl:", oikosUrl)
	if resp, err := request.GetJsonTest(oikosUrl, &Ubicaciones); err == nil && resp.StatusCode == 200 { // (2) error servicio caido
		for _, relacion := range Ubicaciones {
			var data map[string]interface{}
			if jsonString, err := json.Marshal(relacion["EspacioFisicoId"]); err == nil {
				if err2 := json.Unmarshal(jsonString, &data); err2 == nil {
					if number := strings.Index(fmt.Sprintf("%v", data["CodigoAbreviacion"]), Datos.Sede.CodigoAbreviacion); number != -1 {
						Parametros2 = append(Parametros2, map[string]interface{}{
							"Id":              relacion["Id"],
							"DependenciaId":   relacion["DependenciaId"],
							"EspacioFisicoId": relacion["EspacioFisicoId"],
							"Estado":          relacion["Estado"],
							"FechaFin":        relacion["FechaFin"],
							"FechaInicio":     relacion["FechaInicio"],
							"Nombre":          data["Nombre"],
						})
					}
					Parametros = append(Parametros, map[string]interface{}{
						"Relaciones": Parametros2,
					})

				} else {
					logs.Error(err2)
					outputError = map[string]interface{}{
						"funcion": "GetAsignacionSedeDependencia - json.Unmarshal(jsonString, &data)",
						"err":     err2,
						"status":  "500",
					}
					return nil, outputError
				}
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "GetAsignacionSedeDependencia - json.Marshal(relacion[\"EspacioFisicoId\"])",
					"err":     err,
					"status":  "500",
				}
				return nil, outputError
			}
		}

		return Parametros, nil

	} else {
		if err == nil {
			err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAsignacionSedeDependencia - request.GetJsonTest(oikosUrl, &Ubicaciones)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

}

// GetElementos ...
func GetElementos(actaId int) (elementosActa []models.ElementosActa, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetElementos - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var (
		urlcrud   string
		elementos []models.Elemento
		auxE      models.ElementosActa
		soporte   *models.SoporteActaProveedor
	)
	if actaId > 0 { // (1) error parametro
		// Solicita información elementos acta
		urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento?query=SoporteActaId.ActaRecibidoId.Id:" + strconv.Itoa(actaId) +
			",Activo:True&limit=-1"
		if response, err := request.GetJsonTest(urlcrud, &elementos); err == nil && response.StatusCode == 200 {
			// Solicita información unidad elemento
			// urlcrud = "http://" + beego.AppConfig.String("administrativaService") + "/unidad/"
			// fmt.Printf("#Elementos: %v\n", len(elementos))

			if len(elementos) == 0 || elementos[0].Id == 0 {
				err := fmt.Errorf("no elements for Act #%d (or Act not found)", actaId)
				logs.Warn(err)
				outputError = map[string]interface{}{
					"funcion": "GetElementos - len(elementos) == 0 || elementos[0].Id == 0",
					"err":     err,
					"status":  "204",
				}
				return nil, outputError
			}

			for k, elemento := range elementos {
				fmt.Printf("#Elemento: %v\n", k)

				auxE.Id = elemento.Id
				auxE.Nombre = elemento.Nombre
				auxE.Cantidad = elemento.Cantidad
				auxE.Marca = elemento.Marca
				auxE.Serie = elemento.Serie

				// UNIDAD DE MEDIDA
				if elemento.UnidadMedida > 0 {
					if unidad, err := unidadHelper.GetUnidad(elemento.UnidadMedida); err == nil && len(unidad) > 0 {
						auxE.UnidadMedida = unidad[0]
					} else if err != nil {
						return nil, err
					} else {
						err := fmt.Errorf("UnidadMedida '%d' Not Found", elemento.UnidadMedida)
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "GetElementos - unidadHelper.GetUnidad(elemento.UnidadMedida)",
							"err":     err,
							"status":  "500",
						}
						return nil, outputError
					}
				}

				auxE.ValorUnitario = elemento.ValorUnitario
				auxE.Subtotal = elemento.Subtotal
				auxE.Descuento = elemento.Descuento
				auxE.ValorTotal = elemento.ValorTotal
				auxE.PorcentajeIvaId = elemento.PorcentajeIvaId
				auxE.ValorIva = elemento.ValorIva
				auxE.ValorFinal = elemento.ValorFinal
				auxE.SubgrupoCatalogoId = elemento.SubgrupoCatalogoId
				auxE.Verificado = elemento.Verificado
				auxE.TipoBienId = elemento.TipoBienId
				auxE.EstadoElementoId = elemento.EstadoElementoId
				// SOPORTE
				soporte = new(models.SoporteActaProveedor)

				if elemento.SoporteActaId.ProveedorId > 0 {
					if proveedor, err := proveedorHelper.GetProveedorById(elemento.SoporteActaId.ProveedorId); err == nil && len(proveedor) > 0 {
						fmt.Printf("proveedor: %#v\n", proveedor[0])
						soporte.ProveedorId = proveedor[0]
					} else if err != nil {
						return nil, err
					} else {
						err := fmt.Errorf("ProveedorId '%d' Not Found", elemento.SoporteActaId.ProveedorId)
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "GetElementos - proveedorHelper.GetProveedorById(elemento.SoporteActaId.ProveedorId)",
							"err":     err,
							"status":  "500",
						}
						return nil, outputError
					}
				}

				soporte.Id = elemento.SoporteActaId.Id
				soporte.ActaRecibidoId = elemento.SoporteActaId.ActaRecibidoId
				soporte.Consecutivo = elemento.SoporteActaId.Consecutivo
				soporte.Activo = elemento.SoporteActaId.Activo
				soporte.FechaCreacion = elemento.SoporteActaId.FechaCreacion
				soporte.FechaModificacion = elemento.SoporteActaId.FechaModificacion
				soporte.FechaSoporte = elemento.SoporteActaId.FechaSoporte
				auxE.SoporteActaId = soporte

				auxE.Placa = elemento.Placa
				auxE.Activo = elemento.Activo
				auxE.FechaCreacion = elemento.FechaCreacion
				auxE.FechaModificacion = elemento.FechaModificacion

				elementosActa = append(elementosActa, auxE)

			}

			return elementosActa, nil
		} else {
			if err == nil {
				err = fmt.Errorf("undesired State: %d", response.StatusCode)
			}
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetElementos - request.GetJsonTest(urlcrud, &elementos)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
	} else {
		err := errors.New("ID must be greater than 0")
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetElementos - actaId > 0",
			"err":     err,
			"status":  "400",
		}
		return nil, outputError
	}
}

// GetIdElementoPlaca Busca el id de un elemento a partir de su placa
func GetIdElementoPlaca(placa string) (idElemento string, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetIdElementoPlaca - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var urlelemento string
	var elemento []map[string]interface{}
	urlelemento = "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento/?query=Placa:" + placa + "&fields=Id&limit=1"
	if response, err := request.GetJsonTest(urlelemento, &elemento); err == nil {

		if response.StatusCode == 200 {
			for _, element := range elemento {
				if len(element) == 0 {
					return "", nil
				} else {
					return strconv.Itoa(int((element["Id"]).(float64))), nil
				}

			}
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetSoportes - request.GetJsonTest(urlelemento, &elemento)",
			"err":     err,
			"status":  "502",
		}
		return "", outputError
	}
	return
}

// GetAllElementosConsumo obtiene todos los elementos de consumo
func GetAllElementosConsumo() (elementos []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetAllElementosConsumo - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	url := "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento?query=TipoBienId:1,Activo:true"
	if response, err := request.GetJsonTest(url, &elementos); err == nil && response.StatusCode == 200 {
		if len(elementos) == 0 {
			err := errors.New("no hay elementos")
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetAllElementosConsumo - len(elementos) == 0",
				"err":     err,
				"status":  "404",
			}
			return nil, outputError
		} else {
			return elementos, nil
		}

	} else {
		if err == nil {
			err = fmt.Errorf("undesired Status Code: %d", response.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetAllElementosConsumo - request.GetJsonTest(url, &elementos)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

}
