package actaRecibido

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/tealeg/xlsx"
	"github.com/udistrital/arka_mid/helpers/proveedorHelper"
	"github.com/udistrital/arka_mid/helpers/tercerosHelper"
	"github.com/udistrital/arka_mid/helpers/ubicacionHelper"
	"github.com/udistrital/arka_mid/helpers/unidadHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

type Impuesto struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
}

type VigenciaImpuesto struct {
	Id                   int
	Activo               bool
	Tarifa               int64
	PorcentajeAplicacion int
	ImpuestoId           Impuesto
}

type Imp struct {
	PorcentajeAplicacion int
	Tarifa               int
	BasePesos            int
	BaseUvt              int
	CodigoAbreviacion    string
}

type Unidad struct {
	Id          int
	Unidad      string
	Tipo        string
	Descripcion string
	Estado      bool
}

type Subgrupo struct {
	Id                int
	Nombre            string
	Descripcion       string
	Activo            bool
	Codigo            string
	Estado            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
}

// GetAllActasRecibido ...
func GetAllActasRecibidoActivas(states []string, usrWSO2 string) (historicoActa []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/GetAllActasRecibidoActivas",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var Historico []map[string]interface{}
	var Terceros []map[string]interface{}
	var Ubicaciones []map[string]interface{}
	var asignado []*models.Proveedor

	// fmt.Print("Estados Solicitados: ")
	// fmt.Println(states)

	url := "http://" + beego.AppConfig.String("actaRecibidoService") + "historico_acta?limit=-1&query=Activo:true"
	fmt.Println(url)
	// url += ",EstadoActaId__Id:3"
	// TODO: Por rendimiento, TODO lo relacionado a ...
	// - buscar el historico_acta mas reciente
	// - Filtrar por estados
	// ... debería moverse a una o más función(es) y/o controlador(es) del CRUD

	if _, err := request.GetJsonTest(url, &Historico); err == nil { // (2) error servicio caido

		// fmt.Print("historicos:")
		// fmt.Println(len(Historico))

		for _, historicos := range Historico {

			var data_ map[string]interface{}
			var data2_ map[string]interface{}
			var data3_ map[string]interface{}
			var Tercero_ map[string]interface{}
			var Ubicacion_ map[string]interface{}
			var nombreAsignado string

			Ubicacion_ = nil

			if data, err := utilsHelper.ConvertirInterfaceMap(historicos["ActaRecibidoId"]); err == nil {
				data_ = data
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "/GetAllActasRecibidoActivas",
					"err":     err,
					"status":  "500",
				}
				return nil, outputError
			}

			if data, err := utilsHelper.ConvertirInterfaceMap(historicos["EstadoActaId"]); err == nil {
				data2_ = data
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "/GetAllActasRecibidoActivas",
					"err":     err,
					"status":  "500",
				}
				return nil, outputError
			}

			if Terceros == nil {
				if Tercero, err := tercerosHelper.GetNombreTerceroById2(fmt.Sprintf("%v", data_["RevisorId"])); err == nil {
					Tercero_ = Tercero
					Terceros = append(Terceros, Tercero)
				} else {
					logs.Error(err)
					outputError = map[string]interface{}{
						"funcion": "/GetAllActasRecibidoActivas",
						"err":     err,
						"status":  "502",
					}
					return nil, outputError
				}
			} else {
				if keys := len(Terceros[0]); keys != 0 {
					if Tercero, err := utilsHelper.ArrayFind(Terceros, "Id", fmt.Sprintf("%v", data_["RevisorId"])); err == nil {
						if keys := len(Tercero); keys == 0 {
							if Tercero, err := tercerosHelper.GetNombreTerceroById2(fmt.Sprintf("%v", data_["RevisorId"])); err == nil {
								Tercero_ = Tercero
								Terceros = append(Terceros, Tercero)
							} else {
								logs.Error(err)
								outputError = map[string]interface{}{
									"funcion": "/GetAllActasRecibidoActivas",
									"err":     err,
									"status":  "502",
								}
								return nil, outputError
							}
						} else {
							Tercero_ = Tercero
						}
					} else {
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "/GetAllActasRecibidoActivas",
							"err":     err,
							"status":  "500",
						}
						return nil, outputError
					}
				} else {
					if Tercero, err := tercerosHelper.GetNombreTerceroById2(fmt.Sprintf("%v", data_["RevisorId"])); err == nil {
						Tercero_ = Tercero
						Terceros = append(Terceros, Tercero)
					} else {
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "/GetAllActasRecibidoActivas",
							"err":     err,
							"status":  "502",
						}
						return nil, outputError
					}
				}
			}

			if Ubicaciones == nil {
				if ubicacion, err := ubicacionHelper.GetAsignacionSedeDependencia(fmt.Sprintf("%v", data_["UbicacionId"])); err == nil {
					// fmt.Println(ubicacion)
					if keys := len(ubicacion); keys != 0 {
						Ubicacion_ = ubicacion
						Ubicaciones = append(Ubicaciones, ubicacion)
					}

				} else {
					logs.Error(err)
					outputError = map[string]interface{}{
						"funcion": "/GetAllActasRecibidoActivas",
						"err":     err,
						"status":  "502",
					}
					return nil, outputError
				}
			} else {
				if keys := len(Ubicaciones[0]); keys != 0 {
					if ubicacion, err := utilsHelper.ArrayFind(Ubicaciones, "Id", fmt.Sprintf("%v", data_["UbicacionId"])); err == nil {
						if keys := len(ubicacion); keys == 0 {
							if ubicacion, err := ubicacionHelper.GetAsignacionSedeDependencia(fmt.Sprintf("%v", data_["UbicacionId"])); err == nil {
								// fmt.Println(ubicacion)
								if keys := len(ubicacion); keys != 0 {
									Ubicacion_ = ubicacion
									Ubicaciones = append(Ubicaciones, ubicacion)
								}
							} else {
								logs.Error(err)
								outputError = map[string]interface{}{
									"funcion": "/GetAllActasRecibidoActivas",
									"err":     err,
									"status":  "502",
								}
								return nil, outputError
							}
						} else {
							Ubicacion_ = ubicacion
						}
					} else {
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "/GetAllActasRecibidoActivas",
							"err":     err,
							"status":  "500",
						}
						return nil, outputError
					}
				} else {
					if ubicacion, err := ubicacionHelper.GetAsignacionSedeDependencia(fmt.Sprintf("%v", data_["UbicacionId"])); err == nil {
						// fmt.Println(ubicacion)
						if keys := len(ubicacion); keys != 0 {
							Ubicacion_ = ubicacion
							Ubicaciones = append(Ubicaciones, ubicacion)
						}
					} else {
						logs.Error(err)
						outputError = map[string]interface{}{
							"funcion": "/GetAllActasRecibidoActivas",
							"err":     err,
							"status":  "502",
						}
						return nil, outputError
					}
				}
			}

			var tmpAsignadoId = int(data_["PersonaAsignada"].(float64))
			asignado, outputError = proveedorHelper.GetProveedorById(tmpAsignadoId)
			if outputError == nil {
				nombreAsignado = asignado[0].NomProveedor
				// fmt.Println(outputError)
			}

			if Ubicacion_ != nil {
				if jsonString2, err := json.Marshal(Ubicacion_["EspacioFisicoId"]); err == nil {
					if err2 := json.Unmarshal(jsonString2, &data3_); err2 != nil {
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
				data3_ = map[string]interface{}{
					"Nombre": "Ubicacion No Especificada",
				}
			}
			// fmt.Println(data3_)
			Acta := map[string]interface{}{
				"UbicacionId":       data3_["Nombre"],
				"Activo":            data_["Activo"],
				"FechaCreacion":     data_["FechaCreacion"],
				"FechaVistoBueno":   data_["FechaVistoBueno"],
				"FechaModificacion": data_["FechaModificacion"],
				"Id":                data_["Id"],
				"Observaciones":     data_["Observaciones"],
				"RevisorId":         Tercero_["NombreCompleto"],
				"PersonaAsignada":   nombreAsignado,
				"PersonaAsignadaId": int(data_["PersonaAsignada"].(float64)),
				"Estado":            data2_["Nombre"],
			}
			// fmt.Println("Es esto")
			// fmt.Println(Acta)
			historicoActa = append(historicoActa, Acta)
		}

		if len(states) > 0 {
			historicoActa = filtrarActasPorEstados(historicoActa, states)
		}

		// TODO: Manejar concurrencia en las peticiones a otras APIS
		// Referencia: https://www.golang-book.com/books/intro/10
		// TODO: Quitar los parámetros de ID de Proveedor y Contratista
		// de la siguiente función, una vez sea uniforme el espacio de
		// usuarios
		// fmt.Println(usrWSO2)
		// if usrWSO2 != "" {
		// 	if actas, err := filtrarActasSegunRoles(historicoActa, usrWSO2); err == nil {
		// 		historicoActa = actas
		// 	} else {
		// 		logs.Error(err)
		// 		outputError = map[string]interface{}{
		// 			"funcion": "/GetAllActasRecibidoActivas",
		// 			"err":     err,
		// 			"status":  "502",
		// 		}
		// 		return nil, outputError
		// 	}
		// }

		return historicoActa, nil

	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/GetAllActasRecibidoActivas",
			"err":     err,
			"status":  "502",
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
		"IVVAA":          ss,
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

// GetActasRecibidoTipo ...
func GetActasRecibidoTipo(tipoActa int) (actasRecibido []models.ActaRecibidoUbicacion, outputError map[string]interface{}) {
	var (
		urlcrud       string
		historicoActa []*models.HistoricoActa
	)
	if tipoActa != 0 { // (1) error parametro
		urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "historico_acta?query=EstadoActaId.Id:" + strconv.Itoa(tipoActa) + ",Activo:True&limit=-1"
		logs.Debug(urlcrud)
		if response, err := request.GetJsonTest(urlcrud, &historicoActa); err == nil { // (2) error servicio caido
			logs.Debug(historicoActa[0].EstadoActaId)
			if response.StatusCode == 200 { // (3) error estado de la solicitud
				for _, acta := range historicoActa {
					// UBICACION
					ubicacion, err := ubicacionHelper.GetUbicacion(acta.ActaRecibidoId.UbicacionId)

					if err != nil {
						panic(err)
					}

					logs.Debug(ubicacion)

					actaRecibidoAux := models.ActaRecibidoUbicacion{
						Id:                acta.ActaRecibidoId.Id,
						RevisorId:         acta.ActaRecibidoId.RevisorId,
						FechaCreacion:     acta.ActaRecibidoId.FechaCreacion,
						FechaModificacion: acta.ActaRecibidoId.FechaModificacion,
						FechaVistoBueno:   acta.ActaRecibidoId.FechaVistoBueno,
						Observaciones:     acta.ActaRecibidoId.Observaciones,
						Activo:            acta.ActaRecibidoId.Activo,
						EstadoActaId:      acta.EstadoActaId,
						UbicacionId:       ubicacion[0],
					}

					actasRecibido = append(actasRecibido, actaRecibidoAux)
				}
				return actasRecibido, nil
			} else {
				logs.Info("Error (3) estado de la solicitud")
				outputError = map[string]interface{}{"Function": "GetActasRecibidoTipo:GetActasRecibidoTipo", "Error": response.Status}
				return nil, outputError
			}
		} else {
			logs.Info("Error (2) servicio caido")
			outputError = map[string]interface{}{"Function": "GetActasRecibidoTipo", "Error": err}
			return nil, outputError
		}
	} else {
		logs.Info("Error (1) Parametro")
		outputError = map[string]interface{}{"Function": "FuncionalidadMidController:getUserAgora", "Error": "null parameter"}
		return nil, outputError
	}

}

// GetElementos ...
func GetElementos(actaId int) (elementosActa []models.ElementosActa, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/GetElementos",
				"err":     err,
				"status":  "502",
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
	if actaId != 0 { // (1) error parametro
		// Solicita información elementos acta
		urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento?query=SoporteActaId.ActaRecibidoId.Id:" + strconv.Itoa(actaId) +
			",Activo:True&limit=-1"
		if response, err := request.GetJsonTest(urlcrud, &elementos); err == nil && response.StatusCode == 200 {
			// Solicita información unidad elemento
			// urlcrud = "http://" + beego.AppConfig.String("administrativaService") + "/unidad/"
			// fmt.Printf("#Elementos: %v\n", len(elementos))
			for k, elemento := range elementos {
				fmt.Printf("#Elemento: %v\n", k)

				auxE.Id = elemento.Id
				auxE.Nombre = elemento.Nombre
				auxE.Cantidad = elemento.Cantidad
				auxE.Marca = elemento.Marca
				auxE.Serie = elemento.Serie
				// UNIDAD DEMEDIDA
				if unidad, err2 := unidadHelper.GetUnidad(elemento.UnidadMedida); err2 == nil && len(unidad) > 0 {
					auxE.UnidadMedida = unidad[0]
				} else if err2 != nil {
					return nil, err2
				} else {
					outputError = map[string]interface{}{
						"funcion": "/GetElementos",
						"err":     err2,
						"status":  "502",
					}
					logs.Error(outputError)
					return nil, outputError
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
					if proveedor, err2 := proveedorHelper.GetProveedorById(elemento.SoporteActaId.ProveedorId); err2 == nil && len(proveedor) > 0 {
						fmt.Printf("proveedor: %#v\n", proveedor[0])
						soporte.ProveedorId = proveedor[0]
					} else if err2 != nil {
						return nil, err2
					} else {
						outputError = map[string]interface{}{
							"funcion": "/GetElementos",
							"err":     err2,
							"status":  "502",
						}
						logs.Error(outputError)
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
			// logs.Info("Error (2) servicio caido")
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "/GetElementos", "err": err, "status": "502"}
			return nil, outputError
		}
	} else {
		// logs.Info("Error (1) Parametro")
		outputError = map[string]interface{}{
			"funcion": "/GetElementos",
			"err":     "Error (1) Parametro",
			"status":  "502",
		}
		return nil, outputError
	}
}

// GetSoportes ...
func GetSoportes(actaId int) (soportesActa []models.SoporteActaProveedor, outputError map[string]interface{}) {
	var (
		urlcrud   string
		soportes  []models.SoporteActa
		proveedor []*models.Proveedor
		auxS      models.SoporteActaProveedor
	)
	if actaId != 0 { // (1) error parametro
		// Solicita información elementos acta
		urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "soporte_acta?query=ActaRecibidoId:" + strconv.Itoa(actaId) + ",ActaRecibidoId.Activo:True&limit=-1"
		if response, err := request.GetJsonTest(urlcrud, &soportes); err == nil {
			// Solicita información unidad elemento
			urlcrud = "http://" + beego.AppConfig.String("administrativaService") + "/unidad/"
			for _, soporte := range soportes {
				if response.StatusCode == 200 { // (3) error estado de la solicitud
					auxS.Id = soporte.Id
					auxS.Consecutivo = soporte.Consecutivo
					auxS.ActaRecibidoId = soporte.ActaRecibidoId
					auxS.FechaSoporte = soporte.FechaSoporte
					auxS.Activo = soporte.Activo
					// SOPORTE
					proveedor, outputError = proveedorHelper.GetProveedorById(soporte.ProveedorId)
					//soporteAux = new(models.SoporteActaProveedor)
					auxS.ProveedorId = proveedor[0]

					auxS.FechaCreacion = soporte.FechaCreacion
					auxS.FechaModificacion = soporte.FechaModificacion

					soportesActa = append(soportesActa, auxS)
				} else {
					logs.Info("Error (3) estado de la solicitud")
					outputError = map[string]interface{}{"Function": "GetAllActasRecibido:GetAllActasRecibido", "Error": response.Status}
					return nil, outputError
				}
			}
			return soportesActa, nil
		} else {
			logs.Info("Error (2) servicio caido")
			outputError = map[string]interface{}{"Function": "GetIva", "Error": err}
			return nil, outputError
		}
	} else {
		logs.Info("Error (1) Parametro")
		outputError = map[string]interface{}{"Function": "FuncionalidadMidController:GetIva", "Error": "null parameter"}
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
