package actaRecibido

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/tealeg/xlsx"

	// "github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/arka_mid/helpers/crud/administrativa"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

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

	if outputError = administrativa.GetUnidades(&Unidades); outputError != nil {
		return
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
