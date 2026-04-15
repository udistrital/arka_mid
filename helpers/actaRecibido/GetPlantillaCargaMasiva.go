package actaRecibido

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/tealeg/xlsx"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

const (
	nombreHojaCargaMasiva = "CargaMasiva"
	nombreHojaCatalogos   = "Catalogos"
	filasEjemploPlantilla = 12
)

var headersCargaMasiva = []string{
	"Serial Clase",
	"Tipo Bien",
	"Nombre",
	"Marca",
	"Serie",
	"Cantidad",
	"Unidad de Medida",
	"Valor Unitario",
	"Subtotal",
	"Descuento",
	"Porcentaje IVA",
	"Valor IVA",
	"Valor Total",
}

// GetPlantillaCargaMasiva genera la plantilla Excel para cargue masivo de elementos
// y la retorna serializada en base64.
func GetPlantillaCargaMasiva() (respuesta *models.PlantillaArchivoResponse, outputError map[string]interface{}) {
	funcion := "GetPlantillaCargaMasiva - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	catalogos := getPlantillaCatalogosMock()

	file := xlsx.NewFile()

	if err := addHojaCargaMasiva(file, catalogos); err != nil {
		return nil, err
	}

	if err := addHojaCatalogos(file, catalogos); err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	if err := file.Write(&buffer); err != nil {
		return nil, errorCtrl.Error(funcion+"file.Write(&buffer)", err, "500")
	}

	respuesta = &models.PlantillaArchivoResponse{
		FileName: "plantilla_carga_masiva_elementos_" + time.Now().Format("20060102_150405") + ".xlsx",
		MimeType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		File:     base64.StdEncoding.EncodeToString(buffer.Bytes()),
		Version:  "v2",
	}

	return respuesta, nil
}

// addHojaCargaMasiva crea la hoja principal con los encabezados y deja embebidas
// las fórmulas equivalentes a la plantilla de referencia:
//   - Subtotal     = Cantidad * (Valor Unitario - Descuento)
//   - Valor IVA    = Subtotal * Porcentaje IVA
//   - Valor Total  = Subtotal + Valor IVA
//
// En esta plantilla dichas columnas se desplazan por la presencia de
// "Serial Clase" y "Tipo Bien" al inicio, por eso las fórmulas se escriben en:
//   - I: Subtotal
//   - L: Valor IVA
//   - M: Valor Total
func addHojaCargaMasiva(file *xlsx.File, catalogos plantillaCatalogosMock) map[string]interface{} {
	funcion := "addHojaCargaMasiva - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	sheet, err := file.AddSheet(nombreHojaCargaMasiva)
	if err != nil {
		return errorCtrl.Error(funcion+"file.AddSheet(CargaMasiva)", err, "500")
	}

	headerRow := sheet.AddRow()
	for _, header := range headersCargaMasiva {
		headerRow.AddCell().Value = header
	}

	for i := 0; i < filasEjemploPlantilla; i++ {
		filaExcel := i + 2
		row := sheet.AddRow()

		clase := catalogos.Clases[i%len(catalogos.Clases)]
		tipo := catalogos.Tipos[i%len(catalogos.Tipos)]
		unidad := catalogos.Unidades[i%len(catalogos.Unidades)]
		iva := catalogos.Ivas[len(catalogos.Ivas)-1]

		row.AddCell().Value = clase.Codigo + " - " + clase.Nombre
		row.AddCell().Value = tipo.Nombre
		row.AddCell().Value = fmt.Sprintf("Elemento %d", i+1)
		row.AddCell().Value = fmt.Sprintf("Marca %d", (i*2)+1)
		row.AddCell().Value = fmt.Sprintf("SERIE%d", i+1)
		row.AddCell().Value = "1"
		row.AddCell().Value = unidad.Nombre
		row.AddCell().Value = strconv.Itoa((i + 1) * 50000)

		subtotalCell := row.AddCell()
		subtotalCell.SetFormula(fmt.Sprintf("F%d*(H%d-J%d)", filaExcel, filaExcel, filaExcel))

		row.AddCell().Value = "0"
		row.AddCell().Value = formatIVADecimal(iva.Tarifa)

		valorIVACell := row.AddCell()
		valorIVACell.SetFormula(fmt.Sprintf("I%d*K%d", filaExcel, filaExcel))

		valorTotalCell := row.AddCell()
		valorTotalCell.SetFormula(fmt.Sprintf("I%d+L%d", filaExcel, filaExcel))
	}

	_ = sheet.SetColWidth(0, 0, 30)  // Serial Clase
	_ = sheet.SetColWidth(1, 1, 22)  // Tipo Bien
	_ = sheet.SetColWidth(2, 2, 28)  // Nombre
	_ = sheet.SetColWidth(3, 4, 20)  // Marca, Serie
	_ = sheet.SetColWidth(5, 5, 12)  // Cantidad
	_ = sheet.SetColWidth(6, 6, 20)  // Unidad de Medida
	_ = sheet.SetColWidth(7, 12, 16) // Valores

	return nil
}

// addHojaCatalogos crea la hoja con catálogos de apoyo para el diligenciamiento,
// alimentada exclusivamente desde estructuras mockeadas en código.
func addHojaCatalogos(file *xlsx.File, catalogos plantillaCatalogosMock) map[string]interface{} {
	funcion := "addHojaCatalogos - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	sheet, err := file.AddSheet(nombreHojaCatalogos)
	if err != nil {
		return errorCtrl.Error(funcion+"file.AddSheet(Catalogos)", err, "500")
	}

	rowTituloClases := sheet.AddRow()
	rowTituloClases.AddCell().Value = "CLASES"

	rowHeaderClases := sheet.AddRow()
	rowHeaderClases.AddCell().Value = "Clase"

	for _, clase := range catalogos.Clases {
		row := sheet.AddRow()
		row.AddCell().Value = clase.Codigo + " - " + clase.Nombre
	}

	sheet.AddRow()
	sheet.AddRow()

	rowTituloTipos := sheet.AddRow()
	rowTituloTipos.AddCell().Value = "TIPOS DE BIEN"

	rowHeaderTipos := sheet.AddRow()
	rowHeaderTipos.AddCell().Value = "Id"
	rowHeaderTipos.AddCell().Value = "Tipo Bien"

	for _, tipo := range catalogos.Tipos {
		row := sheet.AddRow()
		row.AddCell().Value = strconv.Itoa(tipo.Id)
		row.AddCell().Value = tipo.Nombre
	}

	sheet.AddRow()
	sheet.AddRow()

	rowTituloIVA := sheet.AddRow()
	rowTituloIVA.AddCell().Value = "IVA"

	rowHeaderIVA := sheet.AddRow()
	rowHeaderIVA.AddCell().Value = "Porcentaje IVA"
	rowHeaderIVA.AddCell().Value = "Tarifa"

	for _, iva := range catalogos.Ivas {
		row := sheet.AddRow()
		row.AddCell().Value = formatIVADecimal(iva.Tarifa)
		row.AddCell().Value = strconv.Itoa(iva.Tarifa)
	}

	sheet.AddRow()
	sheet.AddRow()

	rowTituloUnidades := sheet.AddRow()
	rowTituloUnidades.AddCell().Value = "UNIDADES DE MEDIDA"

	rowHeaderUnidades := sheet.AddRow()
	rowHeaderUnidades.AddCell().Value = "Unidad de Medida"

	for _, unidad := range catalogos.Unidades {
		row := sheet.AddRow()
		row.AddCell().Value = unidad.Nombre
	}

	_ = sheet.SetColWidth(0, 0, 24)
	_ = sheet.SetColWidth(1, 1, 45)

	return nil
}

// formatIVADecimal convierte una tarifa entera (19, 5, 0)
// al formato decimal esperado por el parser actual ("0.19", "0.05", "0.00").
func formatIVADecimal(tarifa int) string {
	if tarifa <= 0 {
		return "0.00"
	}
	if tarifa < 10 {
		return "0.0" + strconv.Itoa(tarifa)
	}
	return "0." + strconv.Itoa(tarifa)
}
