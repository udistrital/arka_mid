package actaRecibido

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/tealeg/xlsx"

	// "github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/arka_mid/helpers/crud/administrativa"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// DecodeXlsx2Json Convierte el archivo excel en una lista de elementos
func DecodeXlsx2Json(c multipart.File) (resultado map[string]interface{}, outputError map[string]interface{}) {

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

	resultado = make(map[string]interface{})

	var hojas []string
	subgrupo := &models.DetalleSubgrupo{
		SubgrupoId: &models.Subgrupo{},
		TipoBienId: &models.TipoBien{},
	}

	validar_campos := []string{"Nombre", "Marca", "Serie", "Cantidad", "Unidad de Medida", "Valor Unitario", "Subtotal", "Descuento", "Porcentaje IVA", "Valor IVA", "Valor Total"}
	elementos := make([]*models.PlantillaActa, 0)
	for s, sheet := range xlFile.Sheets {

		if s == 0 {
			indexes := make(map[string]int)
			hojas = append(hojas, sheet.Name)
			for r, row := range sheet.Rows {
				if r == 0 {
					for _, label := range validar_campos {
						index := -1
						for i, cell := range row.Cells {
							if label == cell.String() {
								index = i
								break
							}
						}

						if index > -1 {
							indexes[label] = index
						} else {
							resultado["Mensaje"] = "errorPlantillaActa"
							return resultado, nil
						}
					}

				} else {
					emptyRow := true
					end := false
					fila := new(models.PlantillaActa)

					for i, cell := range row.Cells {
						if i == 0 && cell.String() == "Subtotal" {
							end = true
							break
						}

						if emptyRow && cell.String() != "" {
							emptyRow = false
						}

						if i == indexes["Nombre"] {
							fila.Nombre = cell.String()
						}

						if i == indexes["Marca"] {
							fila.Marca = cell.String()
						}

						if i == indexes["Serie"] {
							fila.Serie = cell.String()
						}

						if i == indexes["Cantidad"] {
							var cant int
							if cant, err = cell.Int(); err != nil {
								cant = 0
							}
							fila.Cantidad = cant
						}

						if i == indexes["Valor Unitario"] {
							var unit float64
							if unit, err = cell.Float(); err != nil {
								unit = 0.0
							}
							fila.ValorUnitario = unit
						}

						if i == indexes["Descuento"] {
							var dcto float64
							if dcto, err = cell.Float(); err != nil {
								dcto = 0.0
							}
							fila.Descuento = dcto
						}

						if i == indexes["Porcentaje IVA"] {
							var tarifa int
							if tarifa_, err := cell.Float(); err != nil {
								tarifa = 0
							} else {
								tarifa = int(tarifa_ * 100)
							}

							for _, tarifa_ := range IvaTest {
								if tarifa == tarifa_.Tarifa {
									fila.PorcentajeIvaId = &tarifa
									break
								}
							}

						}

						if i == indexes["Unidad de Medida"] {
							if cell.String() != "" {
								for _, unidad := range Unidades {
									if cell.String() == unidad.Unidad {
										fila.UnidadMedida = unidad.Id
										break
									}
								}
							}
						}

					}

					if end {
						break
					} else if emptyRow {
						continue
					} else {
						fila.Subtotal = float64(fila.Cantidad) * (fila.ValorUnitario - fila.Descuento)
						if fila.PorcentajeIvaId != nil {
							fila.ValorIva = float64(*fila.PorcentajeIvaId) * fila.Subtotal / 100
						}
						fila.ValorTotal = fila.Subtotal + fila.ValorIva

						fila.SubgrupoCatalogoId = subgrupo
						elementos = append(elementos, fila)
					}

				}
			}
		}
	}

	resultado["Elementos"] = elementos

	return resultado, nil
}
