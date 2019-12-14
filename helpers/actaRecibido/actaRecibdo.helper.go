package actaRecibido

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"encoding/json"
	"mime/multipart"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/tealeg/xlsx"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/arka_mid/models"
)
type Impuesto struct {
	Id						int
    Nombre					string
    Descripcion				string
    CodigoAbreviacion		string
    Activo					bool
}

type VigenciaImpuesto struct {
	Id						int
    Activo					bool
    Tarifa					int64
    PorcentajeAplicacion	int
    ImpuestoId				Impuesto
}

type Unidad struct {
	Id				int
	Unidad			string
    Tipo			string
	Descripcion		string
    Estado			bool
}


// GetAllActasRecibido ...
func GetAllActasRecibido() (historicoActa interface{}, outputError map[string]interface{}) {
	// if idUser != 0 { // (1) error parametro
	if response, err := request.GetJsonTest("http://"+beego.AppConfig.String("actaRecibidoService")+"historico_acta?query=ActaRecibidoId.Activo:True", &historicoActa); err == nil { // (2) error servicio caido
		if response.StatusCode == 200 { // (3) error estado de la solicitud
			// c.Data["json"] = response
			// c.ServeJSON()
			return historicoActa, nil
		} else {
			logs.Info("Error (3) estado de la solicitud")
			outputError = map[string]interface{}{"Function": "GetAllActasRecibido:GetAllActasRecibido", "Error": response.Status}
			return outputError, nil
		}
	} else {
		logs.Info("Error (2) servicio caido")
		outputError = map[string]interface{}{"Function": "GetAllActasRecibido", "Error": err}
		return outputError, nil
	}
	// c.ServeJSON()
	// } else {
	// 	logs.Info("Error (1) Parametro")
	// 	outputError = map[string]interface{}{"Function": "FuncionalidadMidController:getUserAgora", "Error": "null parameter"}
	// 	return nil, outputError
	// }
}

// GetActasRecibidoTipo ...
func GetAllParametrosActa() (Parametros []map[string]interface{}, outputError map[string]interface{}) {

	var Unidades interface{}
	var IVA interface{}
	var TipoBien interface{}
	var EstadoActa interface{}
	var EstadoElemento interface{}
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
	if _, err := request.GetJsonTest("http://"+beego.AppConfig.String("parametrosGobiernoService")+"vigencia_impuesto?limit=-1", &IVA); err == nil { // (2) error servicio caido

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
		"IVA":            IVA,
		"TipoBien":       TipoBien,
		"EstadoActa":     EstadoActa,
		"EstadoElemento": EstadoElemento,
	})

	return parametros, nil
}

// "PostDecodeXlsx2Json ..."
func DecodeXlsx2Json(c multipart.File) (Archivo []map[string]interface{}, outputError map[string]interface{}) {

	var IVA []VigenciaImpuesto
	var Unidades []Unidad

	if _, err := request.GetJsonTest("http://"+beego.AppConfig.String("parametrosGobiernoService")+"vigencia_impuesto?limit=-1", &IVA); err == nil { // (2) error servicio caido
		logs.Info(IVA)
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

	validar_campos := []string{"Nivel Inventarios",	"Tipo de Bien", "Subgrupo Catalogo",	"Nombre",	"Marca", "Serie",	"Cantidad",	"Unidad de Medida", "Valor Unitario", "Subtotal",	"Descuento", "Tipo IVA", "Valor IVA",	"Valor Total",}

	for s, sheet := range xlFile.Sheets {

		if s == 0 {
			hojas = append(hojas, sheet.Name)
			for r, row := range sheet.Rows {
				if r == 0 {
					for i, cell := range row.Cells {
						campos = append(campos, cell.String())
						if campos[i] != validar_campos[i] {
							logs.Info("Error Dependencia servicio caido")
							outputError = map[string]interface{}{"Function": "GetAllActasRecibido","Error": 403}
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
						convertir := strings.Split(elementos[11],".")
						if err == nil {
							logs.Info(convertir)
							valor, err := strconv.ParseInt(convertir[0], 10, 64) 
							if err == nil {
								for _, valor_iva := range IVA{
									if valor == valor_iva.Tarifa {
										elementos[11] = strconv.Itoa(valor_iva.Id)
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
								for _, unidad := range Unidades{
									if convertir2 == unidad.Unidad {
										elementos[7] = strconv.Itoa(unidad.Id)
									}
								}
						} else {
							logs.Info(err)
						}

						Elemento = append(Elemento, map[string]interface{}{
							"NivelInventariosId": elementos[0],
							"TipoBienId":         elementos[1],
							"SubgrupoCatalogoId": elementos[2],
							"Nombre":        	  elementos[3],
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

// GetActasRecibidoTipo ...
func GetAllParametrosSoporte() (Parametros []map[string]interface{}, outputError map[string]interface{}) {

	var Dependencias interface{}
	var Sedes interface{}
	var Ubicaciones interface{}
	parametros := make([]map[string]interface{}, 0)

	if _, err := request.GetJsonTest("http://"+beego.AppConfig.String("oikosService")+"dependencia?limit=-1", &Dependencias); err == nil { // (2) error servicio caido

	} else {
		logs.Info("Error Dependencia servicio caido")
		outputError = map[string]interface{}{"Function": "GetAllActasRecibido", "Error": err}
		return nil, outputError
	}

	if _, err := request.GetJsonTest("http://"+beego.AppConfig.String("oikosService")+"asignacion_espacio_fisico_dependencia?limit=-1", &Ubicaciones); err == nil { // (2) error servicio caido

	} else {
		logs.Info("Error Ubicaciones servicio caido")
		outputError = map[string]interface{}{"Function": "GetAllActasRecibido", "Error": err}
		return nil, outputError
	}
	if _, err := request.GetJsonTest("http://"+beego.AppConfig.String("oikosService")+"espacio_fisico?query=TipoEspacio.Id:1&limit=-1", &Sedes); err == nil { // (2) error servicio caido

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

// GetActasRecibidoTipo ...
func GetAsignacionSedeDependencia(Datos models.GetSedeDependencia) (Parametros []map[string]interface{}, outputError map[string]interface{}) {

	var Ubicaciones []map[string]interface{}
	var Parametros2 []map[string]interface{}
	fmt.Println(Datos.Sede)
	fmt.Println(Datos.Dependencia)
	if _, err := request.GetJsonTest("http://"+beego.AppConfig.String("oikosService") +
		"asignacion_espacio_fisico_dependencia?query=DependenciaId.Id:" + strconv.Itoa(Datos.Dependencia.Id) +
		"&limit=-1", &Ubicaciones); err == nil { // (2) error servicio caido
			fmt.Println(Ubicaciones)
		for _, relacion := range Ubicaciones {
			var data map[string]interface{}
			if jsonString, err := json.Marshal(relacion["EspacioFisicoId"]); err == nil {
				if err2 := json.Unmarshal(jsonString, &data); err2 == nil {
						if number := strings.Index(fmt.Sprintf("%v",data["Codigo"]), Datos.Sede.Codigo); number != -1 {
							Parametros2 = append( Parametros2,  map[string]interface{}{
								"Id":					relacion["Id"],
								"DependenciaId":		relacion["DependenciaId"],
								"EspacioFisicoId":		relacion["EspacioFisicoId"],
								"Estado":				relacion["Estado"],
								"FechaFin":				relacion["FechaFin"],
								"FechaInicio":			relacion["FechaInicio"],
								"Nombre":				data["Nombre"],
							})
						}
						Parametros = append( Parametros,  map[string]interface{}{
							"Relaciones":	Parametros2,
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
