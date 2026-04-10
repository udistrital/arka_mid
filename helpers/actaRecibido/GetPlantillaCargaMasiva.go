package actaRecibido

import (
	"bytes"
	"encoding/base64"
	"strconv"
	"time"

	"github.com/tealeg/xlsx"

	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// GetPlantillaCargaMasiva genera la plantilla Excel para cargue masivo de elementos
// y la retorna serializada en base64.
func GetPlantillaCargaMasiva() (respuesta *models.PlantillaArchivoResponse, outputError map[string]interface{}) {
	funcion := "GetPlantillaCargaMasiva - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	// Consultar catálogos necesarios
	detalles, err := catalogoElementos.GetAllDetalleSubgrupo("limit=-1&query=Activo:true")
	if err != nil {
		return nil, err
	}

	var tiposBien []models.TipoBien
	if err := catalogoElementos.GetAllTipoBien("limit=-1&query=Activo:true", &tiposBien); err != nil {
		return nil, err
	}

	var ivas []models.Iva
	if err := parametros.GetAllIVAByPeriodo(strconv.Itoa(time.Now().Year()), &ivas); err != nil {
		return nil, err
	}

	const payloadUnidades = "limit=-1&fields=Id,Nombre&sortby=Nombre&order=asc&query=TipoParametroId__CodigoAbreviacion__in:L|M|T|C|S"
	unidades, err := parametros.GetAllParametro(payloadUnidades)
	if err != nil {
		return nil, err
	}

	// Crear workbook
	file := xlsx.NewFile()

	// Hoja principal
	if err := addHojaCargaMasiva(file); err != nil {
		return nil, err
	}

	// Hoja de catálogos
	if err := addHojaCatalogos(file, detalles, tiposBien, ivas, unidades); err != nil {
		return nil, err
	}

	// Serializar a bytes
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

// addHojaCargaMasiva crea la hoja principal con los encabezados de cargue.
func addHojaCargaMasiva(file *xlsx.File) map[string]interface{} {
	funcion := "addHojaCargaMasiva - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	sheet, err := file.AddSheet("CargaMasiva")
	if err != nil {
		return errorCtrl.Error(funcion+"file.AddSheet(CargaMasiva)", err, "500")
	}

	// Encabezados compatibles con el parser actual + nuevos campos para seriales
	headers := []string{
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

	headerRow := sheet.AddRow()
	for _, h := range headers {
		cell := headerRow.AddCell()
		cell.Value = h
	}

	// Fila ejemplo
	exampleRow := sheet.AddRow()
	exampleValues := []string{
		"010001 - COMPUTO - (DEV)",
		"CONSUMO",
		"Portátil",
		"Lenovo",
		"ABC123456",
		"1",
		"UNIDAD",
		"3500000",
		"3500000",
		"0",
		"0.19",
		"665000",
		"4165000",
	}

	for _, v := range exampleValues {
		cell := exampleRow.AddCell()
		cell.Value = v
	}

	// Ajuste básico de ancho de columnas
	_ = sheet.SetColWidth(0, 0, 30)  // Nombre
	_ = sheet.SetColWidth(1, 2, 20)  // Marca, Serie
	_ = sheet.SetColWidth(3, 10, 18) // Numéricas
	_ = sheet.SetColWidth(11, 14, 28)

	return nil
}

// addHojaCatalogos crea la hoja con catálogos de apoyo para el diligenciamiento.
func addHojaCatalogos(file *xlsx.File, detalles []*models.DetalleSubgrupo, tiposBien []models.TipoBien, ivas []models.Iva, unidades []*models.Parametro) map[string]interface{} {
	funcion := "addHojaCatalogos - "
	defer errorCtrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	sheet, err := file.AddSheet("Catalogos")
	if err != nil {
		return errorCtrl.Error(funcion+"file.AddSheet(Catalogos)", err, "500")
	}

	// =========================
	// Sección: Clases
	// =========================
	rowTituloClases := sheet.AddRow()
	rowTituloClases.AddCell().Value = "CLASES"

	rowHeaderClases := sheet.AddRow()
	rowHeaderClases.AddCell().Value = "Clase"

	for _, d := range detalles {
		if d == nil {
			continue
		}

		row := sheet.AddRow()

		clase := ""

		if d.SubgrupoId != nil {
			clase = d.SubgrupoId.Codigo + " - " + d.SubgrupoId.Nombre
		}
		row.AddCell().Value = clase
	}

	// Separación visual
	sheet.AddRow()
	sheet.AddRow()

	// =========================
	// Sección: Tipos de bien
	// =========================
	rowTituloTipos := sheet.AddRow()
	rowTituloTipos.AddCell().Value = "TIPOS DE BIEN"

	rowHeaderTipos := sheet.AddRow()
	rowHeaderTipos.AddCell().Value = "Tipo Bien"

	for _, tb := range tiposBien {
		row := sheet.AddRow()
		row.AddCell().Value = strconv.Itoa(tb.Id)
		row.AddCell().Value = tb.Nombre
	}

	// Separación visual
	sheet.AddRow()
	sheet.AddRow()

	// =========================
	// Sección: IVA
	// =========================
	rowTituloIVA := sheet.AddRow()
	rowTituloIVA.AddCell().Value = "IVA"

	rowHeaderIVA := sheet.AddRow()
	rowHeaderIVA.AddCell().Value = "Porcentaje IVA"
	rowHeaderIVA.AddCell().Value = "Tarifa"

	for _, iva := range ivas {
		row := sheet.AddRow()
		// El parser actual interpreta el valor como float y lo multiplica por 100,
		// por eso en plantilla conviene mostrar 0.19, 0.05, etc.
		row.AddCell().Value = formatIVADecimal(iva.Tarifa)
		row.AddCell().Value = strconv.Itoa(iva.Tarifa)
	}

	// Separación visual
	sheet.AddRow()
	sheet.AddRow()

	// =========================
	// Sección: Unidades de medida
	// =========================
	rowTituloUnidades := sheet.AddRow()
	rowTituloUnidades.AddCell().Value = "UNIDADES DE MEDIDA"

	rowHeaderUnidades := sheet.AddRow()
	rowHeaderUnidades.AddCell().Value = "Unidad de Medida"

	for _, unidad := range unidades {
		if unidad == nil {
			continue
		}

		row := sheet.AddRow()
		row.AddCell().Value = unidad.Nombre
	}

	// Ajuste básico de ancho
	_ = sheet.SetColWidth(0, 0, 22)
	_ = sheet.SetColWidth(1, 1, 45)
	_ = sheet.SetColWidth(2, 3, 24)

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
