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
	"github.com/udistrital/arka_mid/helpers/catalogoElementosHelper"
	crud_actas "github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/tercerosHelper"
	"github.com/udistrital/arka_mid/helpers/ubicacionHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"

	// "github.com/udistrital/utils_oas/formatdata"
	"net/url"

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
	urlEstados := "http://" + beego.AppConfig.String("actaRecibidoService") + "historico_acta?limit=-1&sortby=ActaRecibidoId__Id&order=desc"
	urlEstados += "&query=Activo:true"
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
				"funcion": "GetAllActasRecibidoActivas - request.GetJsonTest(urlContProv, &hists)",
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
					if Tercero, err := tercerosHelper.GetTerceroById(id); err == nil {
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

			idUbStr := strconv.Itoa(int(historicos["UbicacionId"].(float64)))
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

			idAsignado := int(historicos["PersonaAsignadaId"].(float64))
			if v, err := utilsHelper.BufferGeneric(idAsignado, Terceros, reqTercero(idAsignado), &consultasTerceros, &evTerceros); err == nil {
				if v2, ok := v.(*models.Tercero); ok {
					asignado = v2
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

		logs.Info("consultasTerceros:", consultasTerceros, " - Evitadas: ", evTerceros)
		logs.Info("consultasUbicaciones:", consultasUbicaciones, " - Evitadas: ", evUbicaciones)
		logs.Info("consultasProveedores:", consultasProveedores, " - Evitadas: ", evProveedores)
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

	var (
		Unidades  []Unidad
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
	tipoBien := new(models.TipoBien)
	subgrupoId := new(models.Subgrupo)
	var subgrupo = map[string]interface{}{
		"SubgrupoId": &subgrupoId,
		"TipoBienId": &tipoBien,
	}

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

					var vlrcantidad int64
					var tarifaIva float64
					var vlrsubtotal float64
					var vlrdcto float64
					var vlrunitario float64
					var vlrIva = float64(-1)

					if elementos[0] != "Totales" {
						if vlrcantidad, err = strconv.ParseInt(elementos[6], 10, 64); err != nil {
							vlrcantidad = 0
						}

						if vlrunitario, err = strconv.ParseFloat(elementos[8], 64); err != nil {
							vlrunitario = float64(0)
						}

						if vlrdcto, err = strconv.ParseFloat(elementos[10], 64); err != nil {
							vlrdcto = float64(0)
						}

						vlrsubtotal = float64(vlrcantidad) * (vlrunitario - vlrdcto)

						if tarifaIva, err = strconv.ParseFloat(strings.ReplaceAll(elementos[11], "%", ""), 64); err == nil {
							for _, valor_iva := range IvaTest {
								if tarifaIva == float64(valor_iva.Tarifa) {
									vlrIva = (vlrsubtotal) * float64(tarifaIva) / 100
								}
							}
							if vlrIva == -1 {
								tarifaIva = 0
								vlrIva = 0
							}
						} else {
							tarifaIva = 0
							vlrIva = 0
						}

						vlrtotal := vlrsubtotal + vlrIva

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

						Elemento = append(Elemento, map[string]interface{}{
							"Id":                 0,
							"SubgrupoCatalogoId": subgrupo,
							"Nombre":             elementos[3],
							"Marca":              elementos[4],
							"Serie":              elementos[5],
							"Cantidad":           vlrcantidad,
							"UnidadMedida":       elementos[7],
							"ValorUnitario":      vlrunitario,
							"Subtotal":           vlrsubtotal,
							"Descuento":          vlrdcto,
							"PorcentajeIvaId":    tarifaIva,
							"ValorIva":           vlrIva,
							"ValorTotal":         vlrtotal,
						})
					} else {
						Respuesta = append(Respuesta, map[string]interface{}{
							"Hoja":      hojas,
							"Campos":    campos,
							"Elementos": Elemento,
						})
						break
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
func GetElementos(actaId int, ids []int) (elementosActa []*models.DetalleElemento, outputError map[string]interface{}) {

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
		urlcrud string
		auxE    *models.DetalleElemento
	)

	subgrupos := make(map[int]interface{})
	consultasSubgrupos := 0
	evSubgrupos := 0

	if actaId > 0 || len(ids) > 0 { // (1) error parametro
		// Solicita información elementos acta

		var query string
		if actaId > 0 {
			query += "Activo:True,ActaRecibidoId__Id:" + strconv.Itoa(actaId)
		} else {
			query += "Id__in:" + utilsHelper.ArrayToString(ids, "|")
		}

		if elementos, err := crud_actas.GetAllElemento(query, "", "Id", "desc", "", "-1"); err != nil {
			return nil, err
		} else {

			if len(elementos) == 0 || elementos[0].Id == 0 {
				return nil, nil
			}

			for _, elemento := range elementos {

				var subgrupoId *models.Subgrupo
				subgrupoId = new(models.Subgrupo)
				var tipoBienId *models.TipoBien
				tipoBienId = new(models.TipoBien)
				auxE = new(models.DetalleElemento)
				subgrupo := *&models.DetalleSubgrupo{
					SubgrupoId: subgrupoId,
					TipoBienId: tipoBienId,
				}

				subgrupo.TipoBienId = tipoBienId
				subgrupo.SubgrupoId = subgrupoId

				idSubgrupo := elemento.SubgrupoCatalogoId
				reqSubgrupo := func() (interface{}, map[string]interface{}) {
					urlcrud = "query=Activo:true,SubgrupoId__Id:" + strconv.Itoa(idSubgrupo)
					urlcrud += "&fields=SubgrupoId,TipoBienId,Depreciacion,Amortizacion,ValorResidual,VidaUtil&sortby=Id&order=desc"
					if detalleSubgrupo_, err := catalogoElementosHelper.GetAllDetalleSubgrupo(urlcrud); err == nil && len(detalleSubgrupo_) > 0 {
						return detalleSubgrupo_[0], nil
					} else if err != nil {
						return nil, err
					} else {
						logs.Error(err)
						return nil, map[string]interface{}{
							"funcion": "GetElementos - catalogoElementosHelper.GetDetalleSubgrupo(idSubgrupo)",
							"err":     err,
							"status":  "500",
						}
					}
				}

				if idSubgrupo > 0 {
					if v, err := utilsHelper.BufferGeneric(idSubgrupo, subgrupos, reqSubgrupo, &consultasSubgrupos, &evSubgrupos); err == nil {
						if v != nil {
							if jsonString, err := json.Marshal(v); err == nil {
								if err := json.Unmarshal(jsonString, &subgrupo); err != nil {
									logs.Error(err)
									outputError = map[string]interface{}{
										"funcion": "GetElementos - json.Unmarshal(jsonString, &subgrupo)",
										"err":     err,
										"status":  "500",
									}
									return nil, outputError
								}
							}
						}
					}
				}

				auxE.Id = elemento.Id
				auxE.Nombre = elemento.Nombre
				auxE.Cantidad = elemento.Cantidad
				auxE.Marca = elemento.Marca
				auxE.Serie = elemento.Serie
				auxE.UnidadMedida = elemento.UnidadMedida
				auxE.ValorUnitario = elemento.ValorUnitario
				auxE.Subtotal = elemento.Subtotal
				auxE.Descuento = elemento.Descuento
				auxE.ValorTotal = elemento.ValorTotal
				auxE.PorcentajeIvaId = elemento.PorcentajeIvaId
				auxE.ValorIva = elemento.ValorIva
				auxE.ValorFinal = elemento.ValorFinal
				auxE.SubgrupoCatalogoId = &subgrupo
				auxE.EstadoElementoId = elemento.EstadoElementoId
				auxE.ActaRecibidoId = elemento.ActaRecibidoId
				auxE.Placa = elemento.Placa
				auxE.Activo = elemento.Activo
				auxE.FechaCreacion = elemento.FechaCreacion
				auxE.FechaModificacion = elemento.FechaModificacion

				elementosActa = append(elementosActa, auxE)

			}

			logs.Info("consultasSubgrupos:", consultasSubgrupos, " - Evitadas: ", evSubgrupos)
			return elementosActa, nil
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
