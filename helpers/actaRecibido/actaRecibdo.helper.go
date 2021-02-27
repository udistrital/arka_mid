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
				"funcion": "/GetAllActasRecibidoActivas - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	// PARTE "0": Buffers, para evitar repetir consultas...
	var Historico []map[string]interface{}
	Terceros := make(map[int](map[string]interface{}))
	Ubicaciones := make(map[int](map[string]interface{}))
	Proveedores := make(map[int](*models.Proveedor))

	consultasTerceros := 0
	consultasUbicaciones := 0
	consultasProveedores := 0
	evTerceros := 0
	evUbicaciones := 0
	evProveedores := 0

	// findAndAddTercero trae la información de un tercero y la agrega
	// al buffer de terceros
	findAndAddTercero := func(TerceroId int) (map[string]interface{}, map[string]interface{}) {

		if TerceroId == 0 {
			return nil, nil
		}

		idStr := fmt.Sprint(TerceroId)

		if Tercero, ok := Terceros[TerceroId]; ok {
			evTerceros++
			return Tercero, nil
		}

		consultasTerceros++
		if Tercero, err := tercerosHelper.GetNombreTerceroById(idStr); err == nil {
			if keys := len(Tercero); keys != 0 {
				Terceros[TerceroId] = Tercero
			}
			return Tercero, nil
		} else {
			logs.Error(err)
			return nil, map[string]interface{}{
				"funcion": "/GetAllActasRecibidoActivas/findAndAddTercero",
				"err":     err,
				"status":  "502",
			}
		}
	}

	// findAndAddUbicacion trae la información de una ubicación y la agrega
	// al buffer de ubicaciones
	findAndAddUbicacion := func(UbicacionId int) (map[string]interface{}, map[string]interface{}) {

		if UbicacionId == 0 {
			return nil, nil
		}

		idStr := fmt.Sprint(UbicacionId)

		if ubicacion, ok := Ubicaciones[UbicacionId]; ok {
			evUbicaciones++
			return ubicacion, nil
		}

		consultasUbicaciones++
		if ubicacion, err := ubicacionHelper.GetAsignacionSedeDependencia(idStr); err == nil {
			if keys := len(ubicacion); keys != 0 {
				Ubicaciones[UbicacionId] = ubicacion
			}
			return ubicacion, nil
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "/GetAllActasRecibidoActivas/findAndAddUbicacion",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
	}

	// findAndAddUbicacion trae la información de un proveedor y la agrega
	// al buffer de proveedores
	// (Nota: Evitar usar, se va a usar terceros en vez de Agora)
	findAndAddProveedor := func(ProveedorId int) (*models.Proveedor, map[string]interface{}) {

		var vacio models.Proveedor
		if ProveedorId == 0 {
			return &vacio, nil
		}

		if proveedor, ok := Proveedores[ProveedorId]; ok {
			evProveedores++
			return proveedor, nil
		}

		// if ProveedorId > 0 {
		consultasProveedores++
		if provs, err := proveedorHelper.GetProveedorById(ProveedorId); err == nil {
			if len(provs) == 1 && provs[0].Id > 0 {
				proveedor := provs[0]
				Proveedores[ProveedorId] = proveedor
				return proveedor, nil
			}
			if len(provs) > 1 {
				logs.Warn("Proveedor", ProveedorId, "tiene más de un resultado")
			}
		} else {
			return nil, err
		}

		// }

		return &vacio, nil
	}

	// PARTE 1 - Identificar los tipos de actas que hay que traer
	// (y así definir la estrategia para traer las actas)
	verTodasLasActas := false
	algunosEstados := []string{}
	proveedor := false
	contratista := false
	idTercero := 0
	idProveedor := 0 // TODO: Eliminar esta variable cuando se pasen los proveedores a Terceros, usar solo idTercero
	// De especificarse un usuario, hay que definir las actas que puede ver
	if usrWSO2 != "" {

		// Traer la información de Autenticación MID para obtener los roles
		var usr models.UsuarioAutenticacion
		if data, err := autenticacion.DataUsuario(usrWSO2); err == nil && data.Role != nil {
			// logs.Debug(data)
			usr = data
		} else if err == nil { // data.Role == nil
			err := fmt.Errorf("El usuario '%s' no está registrado en WSO2 y/o no tiene roles asignados", usrWSO2)
			logs.Warn(err)
			outputError = map[string]interface{}{
				"funcion": "GetAllActasRecibidoActivas - autenticacion.DataUsuario(usrWSO2)",
				"err":     err,
				"status":  "404",
			}
			return nil, outputError
		} else {
			return nil, err
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
			// TODO: Debería ser suficiente cambiar la condicion por 'proveedor || contratista'
			// una vez se decida cargar los proveedores también de terceros
			if contratista && idTercero == 0 {
				if data, err := tercerosHelper.GetTerceroByUsuarioWSO2(usrWSO2); err == nil {
					if v, ok := data["Id"].(int); ok {
						idTercero = v
						// Terceros[v] = data // Se podría agregar de una vez, pero no incluye el documento
					}
				} else {
					return nil, err
				}
			}
			// TODO: Eliminar el siguiente bloque cuando se pasen los proveedores a Terceros
			// (Debería ser suficiente con el bloque anterior, idTercero)
			if proveedor && idProveedor == 0 {
				// 1. A partir del documento obtenido de WSO2, buscar el proveedor
				consultasProveedores++
				if data, err := proveedorHelper.GetProveedorByDoc(usr.Documento); err == nil {
					idProveedor = data.Id
					Proveedores[idProveedor] = data
				} else {
					return nil, err
				}
			}
		}
	}
	logs.Info("u:", usrWSO2, "- t:", verTodasLasActas, "- e:", algunosEstados, "- p:", proveedor, "- c:", contratista, "- i:", idTercero)
	logs.Info("iP:", idProveedor)

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

	return nil, nil

	// PARTE 2: Traer los tipos de actas identificados
	// (con base a la estrategia definida anteriormente)

	url := "http://" + beego.AppConfig.String("actaRecibidoService") + "historico_acta?limit=-1&query=Activo:true"
	// fmt.Println(url)
	// url += ",EstadoActaId__Id:3"
	// TODO: Por rendimiento, TODO lo relacionado a ...
	// - buscar el historico_acta mas reciente
	// - Filtrar por estados
	// ... debería moverse a una o más función(es) y/o controlador(es) del CRUD

	if resp, err := request.GetJsonTest(url, &Historico); err == nil && resp.StatusCode == 200 { // (2) error servicio caido

		// fmt.Print("historicos:")
		// fmt.Println(len(Historico))

		if len(Historico) == 0 || len(Historico[0]) == 0 {
			err := errors.New("There's currently no act records")
			logs.Warn(err)
			outputError = map[string]interface{}{
				"funcion": "/GetAllActasRecibidoActivas",
				"err":     err,
				"status":  "200", // TODO: Debería ser un 204 pero el cliente (Angular) se ofende... (hay que hacer varios ajustes)
			}
			return nil, outputError
		}

		for _, historicos := range Historico {

			var acta map[string]interface{}
			var estado map[string]interface{}
			var ubicacionData map[string]interface{}
			var editor map[string]interface{}
			var preUbicacion map[string]interface{}
			// var nombreAsignado string
			var oldAsignado *models.Proveedor
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

			idRev := int(acta["RevisorId"].(float64))
			if v, err := findAndAddTercero(idRev); err == nil {
				editor = v
			} else {
				return nil, err
			}

			idUb := int(acta["UbicacionId"].(float64))
			if v, err := findAndAddUbicacion(idUb); err == nil {
				preUbicacion = v
			} else {
				return nil, err
			}

			var tmpAsignadoId = int(acta["PersonaAsignada"].(float64))
			if v, err := findAndAddProveedor(tmpAsignadoId); err == nil {
				oldAsignado = v
			}
			if v, err := findAndAddTercero(tmpAsignadoId); err == nil {
				asignado = v
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
				"oldAsignada":       oldAsignado.NomProveedor,
				"PersonaAsignada":   asignado["NombreCompleto"],
				"PersonaAsignadaId": tmpAsignadoId,
				"Estado":            estado["Nombre"],
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

	} else if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/GetAllActasRecibidoActivas - request.GetJsonTest(url, &Historico)",
			"err":     err,
			"status":  "502", // (2) error servicio caido
		}
		return nil, outputError
	} else {
		err := fmt.Errorf("Undesired Status Code: %d", resp.StatusCode)
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/GetAllActasRecibidoActivas - request.GetJsonTest(url, &Historico)",
			"err":     err,
			"status":  "502", // (2) error servicio caido
		}
		return nil, outputError
	}
}

func RemoveIndex(s []byte, index int) []byte {
	return append(s[:index], s[index+1:]...)
}

// GetAllParametrosActa ...
func GetAllParametrosActa() (Parametros []map[string]interface{}, outputError map[string]interface{}) {

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

	if _, err := request.GetJsonTest("http://"+beego.AppConfig.String("actaRecibidoService")+"tipo_bien?limit=-1", &TipoBien); err == nil { // (2) error servicio caido

	} else {
		logs.Info("Error TipoBien servicio Acta caido")
		outputError = map[string]interface{}{"Function": "GetAllActasRecibido", "Error": err}
		return nil, outputError
	}
	if _, err := request.GetJsonTest("http://"+beego.AppConfig.String("actaRecibidoService")+"estado_acta?limit=-1", &EstadoActa); err == nil { // (2) error servicio caido

	} else {
		logs.Info("Error EstadoActa servicio caido")
		outputError = map[string]interface{}{"Function": "GetAllActasRecibido", "Error": err}
		return nil, outputError
	}
	if _, err := request.GetJsonTest("http://"+beego.AppConfig.String("actaRecibidoService")+"estado_elemento?limit=-1", &EstadoElemento); err == nil { // (2) error servicio caido

	} else {
		logs.Info("Error EstadoElemento servicio caido")
		outputError = map[string]interface{}{"Function": "GetAllActasRecibido", "Error": err}
		return nil, outputError
	}

	if _, err := request.GetJsonTest("http://"+beego.AppConfig.String("parametrosService")+"parametro_periodo?query=PeriodoId__Nombre:2021,ParametroId__TipoParametroId__Id:12", &ss); err == nil { // (2) error servicio caido

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
		logs.Info("Error IVA servicio caido")
		outputError = map[string]interface{}{"Function": "GetAllActasRecibido", "Error": err}
		return nil, outputError
	}

	if _, err := request.GetJsonTest("http://"+beego.AppConfig.String("AdministrativaService")+"unidad?limit=-1", &Unidades); err == nil { // (2) error servicio caido

	} else {
		logs.Info("Error Unidades servicio caido")
		outputError = map[string]interface{}{"Function": "GetAllActasRecibido", "Error": err}
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

	if _, err := request.GetJsonTest("http://"+beego.AppConfig.String("parametrosService")+"parametro_periodo?query=PeriodoId__Nombre:2021,ParametroId__TipoParametroId__Id:12", &ss); err == nil { // (2) error servicio caido

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
		logs.Info("Error IVA servicio caido")
		outputError = map[string]interface{}{"Function": "GetAllActasRecibido", "Error": err}
		return nil, outputError
	}

	if _, err := request.GetJsonTest("http://"+beego.AppConfig.String("catalogoElementosService")+"tr_catalogo/tipo_de_bien/1", &SubgruposConsumo); err == nil { // (2) error servicio caido

	} else {
		logs.Info("Error IVA servicio caido")
		outputError = map[string]interface{}{"Function": "GetAllActasRecibido", "Error": err}
		return nil, outputError
	}

	if _, err := request.GetJsonTest("http://"+beego.AppConfig.String("catalogoElementosService")+"tr_catalogo/tipo_de_bien/2", &SubgruposConsumoControlado); err == nil { // (2) error servicio caido

	} else {
		logs.Info("Error IVA servicio caido")
		outputError = map[string]interface{}{"Function": "GetAllActasRecibido", "Error": err}
		return nil, outputError
	}

	if _, err := request.GetJsonTest("http://"+beego.AppConfig.String("catalogoElementosService")+"tr_catalogo/tipo_de_bien/3", &SubgruposDevolutivo); err == nil { // (2) error servicio caido

	} else {
		logs.Info("Error IVA servicio caido")
		outputError = map[string]interface{}{"Function": "GetAllActasRecibido", "Error": err}
		return nil, outputError
	}
	if _, err := request.GetJsonTest("http://"+beego.AppConfig.String("AdministrativaService")+"unidad?limit=-1", &Unidades); err == nil { // (2) error servicio caido

	} else {
		logs.Info("Error Unidades servicio caido")
		outputError = map[string]interface{}{"Function": "GetAllActasRecibido", "Error": err}
		return nil, outputError
	}

	file, err := ioutil.ReadAll(c)
	if err != nil {
		fmt.Println("err reading file", err)
		logs.Info("Error (1) error de recepcion")
		outputError = map[string]interface{}{"Function": "PostDecodeXlsx2Json", "Error": 400}
		return nil, outputError
	}
	xlFile, err := xlsx.OpenBinary(file)
	if err != nil {
		fmt.Println("err reading file", err)
		logs.Info("Error (1) error de recepcion")
		outputError = map[string]interface{}{"Function": "PostDecodeXlsx2Json", "Error": 400}
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
							logs.Info("Error Dependencia servicio caido")
							outputError = map[string]interface{}{"Function": "GetAllActasRecibido", "Error": 403}
							Respuesta2 := append(Respuesta, map[string]interface{}{
								"Mensaje": "El formato no corresponde a las columnas necesarias",
							})
							return Respuesta2, outputError
						}
					}
				} else {

					for i, cell := range row.Cells {
						elementos[i] = cell.String()
					}
					if elementos[0] != "Totales" {
						convertir := strings.Split(elementos[11], ".")
						if err == nil {
							logs.Info(convertir)
							valor, err := strconv.ParseInt(convertir[0], 10, 64)
							if err == nil {
								for _, valor_iva := range IvaTest {
									if valor == int64(valor_iva.Tarifa) {
										elementos[11] = strconv.Itoa(valor_iva.Tarifa)
									}
								}
							} else {
								logs.Info(err)
							}
						} else {
							logs.Info(err)
						}

						convertir2 := strings.ToUpper(elementos[7])
						if err == nil {
							logs.Info(convertir2)
							for _, unidad := range Unidades {
								if convertir2 == unidad.Unidad {
									elementos[7] = strconv.Itoa(unidad.Id)
								}
							}
						} else {
							logs.Info(err)
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
							logs.Info(err)
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

	var Dependencias interface{}
	var Sedes interface{}
	var Ubicaciones interface{}
	parametros := make([]map[string]interface{}, 0)

	if _, err := request.GetJsonTest("http://"+beego.AppConfig.String("oikos2Service")+"dependencia?limit=-1", &Dependencias); err == nil { // (2) error servicio caido

	} else {
		logs.Info("Error Dependencia servicio caido")
		outputError = map[string]interface{}{"Function": "GetAllActasRecibido", "Error": err}
		return nil, outputError
	}

	if _, err := request.GetJsonTest("http://"+beego.AppConfig.String("oikos2Service")+"asignacion_espacio_fisico_dependencia?limit=-1", &Ubicaciones); err == nil { // (2) error servicio caido

	} else {
		logs.Info("Error Ubicaciones servicio caido")
		outputError = map[string]interface{}{"Function": "GetAllActasRecibido", "Error": err}
		return nil, outputError
	}
	if _, err := request.GetJsonTest("http://"+beego.AppConfig.String("oikos2Service")+"espacio_fisico?query=TipoEspacioFisicoId.Id:1&limit=-1", &Sedes); err == nil { // (2) error servicio caido

	} else {
		logs.Info("Error Sedes servicio caido")
		outputError = map[string]interface{}{"Function": "GetAllActasRecibido", "Error": err}
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

	var Ubicaciones []map[string]interface{}
	var Parametros2 []map[string]interface{}
	fmt.Println(Datos.Sede)
	fmt.Println(Datos.Dependencia)
	if _, err := request.GetJsonTest("http://"+beego.AppConfig.String("oikos2Service")+
		"asignacion_espacio_fisico_dependencia?query=DependenciaId.Id:"+strconv.Itoa(Datos.Dependencia.Id)+
		"&limit=-1", &Ubicaciones); err == nil { // (2) error servicio caido
		fmt.Println(Ubicaciones)
		for _, relacion := range Ubicaciones {
			var data map[string]interface{}
			if jsonString, err := json.Marshal(relacion["EspacioFisicoId"]); err == nil {
				if err2 := json.Unmarshal(jsonString, &data); err2 == nil {
					if number := strings.Index(fmt.Sprintf("%v", data["Codigo"]), Datos.Sede.Codigo); number != -1 {
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
					logs.Info("Error asignacion_espacio_fisico_dependencia servicio caido")
					outputError = map[string]interface{}{"Function": "GetAsignacionSedeDependencia", "Error": err2}
					return nil, outputError
				}
			} else {
				logs.Info("Error asignacion_espacio_fisico_dependencia servicio caido")
				outputError = map[string]interface{}{"Function": "GetAsignacionSedeDependencia", "Error": err}
				return nil, outputError
			}
		}

		return Parametros, nil

	} else {
		logs.Info("Error asignacion_espacio_fisico_dependencia servicio caido")
		outputError = map[string]interface{}{"Function": "GetAsignacionSedeDependencia", "Error": err}
		return nil, outputError
	}

}

// GetElementos ...
func GetElementos(actaId int) (elementosActa []models.ElementosActa, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/GetElementos - Unhandled Error!",
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
				err := fmt.Errorf("No elements for Act #%d (or Act not found)", actaId)
				logs.Warn(err)
				outputError = map[string]interface{}{
					"funcion": "/GetElementos - len(elementos) == 0 || elementos[0].Id == 0",
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
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "/GetElementos - unidadHelper.GetUnidad(elemento.UnidadMedida)",
							"err":     err,
							"status":  "502",
						}
						return nil, outputError
					} else {
						err := fmt.Errorf("UnidadMedida '%d' Not Found", elemento.UnidadMedida)
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "/GetElementos - unidadHelper.GetUnidad(elemento.UnidadMedida) / len(unidad) > 0",
							"err":     err,
							"status":  "502",
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
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "/GetElementos - proveedorHelper.GetProveedorById(elemento.SoporteActaId.ProveedorId)",
							"err":     err,
							"status":  "502",
						}
						return nil, outputError
					} else {
						err := fmt.Errorf("ProveedorId '%d' Not Found", elemento.SoporteActaId.ProveedorId)
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "/GetElementos - proveedorHelper.GetProveedorById(elemento.SoporteActaId.ProveedorId)",
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
		} else if err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "/GetElementos - request.GetJsonTest(urlcrud, &elementos)",
				"err":     err,
				"status":  "502", // Error (2) servicio caido
			}
			return nil, outputError
		} else {
			err := fmt.Errorf("Undesired State: %d", response.StatusCode)
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "/GetElementos - request.GetJsonTest(urlcrud, &elementos)",
				"err":     err,
				"status":  "500",
			}
			return nil, outputError
		}
	} else {
		err := errors.New("ID must be greater than 0")
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/GetElementos - actaId > 0",
			"err":     err,
			"status":  "400",
		}
		return nil, outputError
	}
}

// GetSoportes ...
func GetSoportes(actaId int) (soportesActa []models.SoporteActaProveedor, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/GetSoportes - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var (
		urlcrud   string
		soportes  []models.SoporteActa
		proveedor []*models.Proveedor
		auxS      models.SoporteActaProveedor
	)
	if actaId > 0 { // (1) error parametro
		// Solicita información elementos acta
		urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "soporte_acta?query=ActaRecibidoId:" + strconv.Itoa(actaId) + ",ActaRecibidoId.Activo:True&limit=-1"
		if response, err := request.GetJsonTest(urlcrud, &soportes); err == nil && response.StatusCode == 200 {

			if len(soportes) == 0 || soportes[0].Id == 0 {
				err = fmt.Errorf("El Acta #%d no existe o no tiene soportes", actaId)
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "/GetSoportes - len(soportes) == 0 || soportes[0].Id == 0",
					"err":     err,
					"status":  "200",
				}
				return nil, outputError
			}
			// Solicita información unidad elemento
			for _, soporte := range soportes {
				auxS.Id = soporte.Id
				auxS.Consecutivo = soporte.Consecutivo
				auxS.ActaRecibidoId = soporte.ActaRecibidoId
				auxS.FechaSoporte = soporte.FechaSoporte
				auxS.Activo = soporte.Activo
				// SOPORTE
				if soporte.ProveedorId > 0 {
					proveedor, outputError = proveedorHelper.GetProveedorById(soporte.ProveedorId)
					//soporteAux = new(models.SoporteActaProveedor)
					auxS.ProveedorId = proveedor[0]
				}

				auxS.FechaCreacion = soporte.FechaCreacion
				auxS.FechaModificacion = soporte.FechaModificacion

				soportesActa = append(soportesActa, auxS)
			}
			return soportesActa, nil
		} else {
			if err == nil {
				err = fmt.Errorf("Undesired Status Code: %d", response.StatusCode)
			}
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "/GetSoportes - request.GetJsonTest(urlcrud, &soportes)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
	} else {
		err := fmt.Errorf("Wrong ActaID. Must be greater than 0 - Got: %d", actaId)
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/GetSoportes - actaId > 0",
			"err":     err,
			"status":  "400",
		}
		return nil, outputError
	}
}

// GetIdElementoPlaca Busca el id de un elemento a partir de su placa
func GetIdElementoPlaca(placa string) (idElemento string, err error) {
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
		return "", err
	}
	return
}

// GetAllElementosConsumo obtiene todos los elementos de consumo
func GetAllElementosConsumo() (elementos []map[string]interface{}, outputError map[string]interface{}) {
	var url string
	url = "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento/?query=TipoBienId:1,Activo:true"
	if response, err := request.GetJsonTest(url, &elementos); err == nil {
		if response.StatusCode == 200 {
			if len(elementos) == 0 {
				return nil, map[string]interface{}{"Function": "GetAllElementosConsumo", "Error": errors.New("No se encontro registro")}
			} else {
				return elementos, nil
			}

		} else if response.StatusCode == 400 {
			return nil, map[string]interface{}{"Function": "GetAllElementosConsumo", "Error": errors.New("No se encontro registro")}
		}
	} else {
		fmt.Println("error: ", err)
		return nil, map[string]interface{}{"Function": "GetAllElementosConsumo", "Error": err}
	}

	return
}
