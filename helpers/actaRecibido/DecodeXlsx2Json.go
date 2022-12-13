package actaRecibido

import (
	"io/ioutil"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/tealeg/xlsx"

	"github.com/udistrital/arka_mid/helpers/crud/administrativa"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// DecodeXlsx2Json Convierte el archivo excel en una lista de elementos
func DecodeXlsx2Json(c multipart.File) (resultado map[string]interface{}, outputError map[string]interface{}) {

	funcion := "GetAllIVAByPeriodo - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	var (
		Unidades []Unidad
		Ivas     []models.Iva
	)

	if err := parametros.GetAllIVAByPeriodo(strconv.Itoa(time.Now().Year()-1), &Ivas); err != nil {
		return nil, err
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

							for _, tarifa_ := range Ivas {
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

						elementos = append(elementos, fila)
					}

				}
			}
		}
	}

	resultado["Elementos"] = elementos

	return resultado, nil
}
