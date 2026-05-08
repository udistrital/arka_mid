package reportesHelper

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/tealeg/xlsx"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	crudActaRecibido "github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

const (
	dateLayout     = "2006-01-02"
	dateTimeLayout = "2006-01-02 15:04:05"
	sheetName      = "Elementos"
)

var (
	reporteElementosHeaders = []string{
		"Id",
		"Nombre",
		"Cantidad",
		"Marca",
		"Serie",
		"UnidadMedida",
		"ValorUnitario",
		"Subtotal",
		"Descuento",
		"ValorTotal",
		"PorcentajeIvaId",
		"ValorIva",
		"ValorFinal",
		"SubgrupoCatalogoId",
		"TipoBienId",
		"EstadoElementoId",
		"ActaRecibidoId",
		"Placa",
		"Activo",
		"FechaCreacion",
		"FechaModificacion",
	}

	consultarElementosReporte = consultarElementosReporteDefault
)

// GenerarReporteElementos genera un archivo Excel en base64 con la estructura
// de acta_recibido/elementos para el rango de fechas solicitado.
func GenerarReporteElementos(req *models.ReporteFechasRequest) (respuesta *models.ReporteExcelBase64Response, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("GenerarReporteElementos - Unhandled Error!", "500")

	if req == nil {
		return nil, errorCtrl.Error("GenerarReporteElementos - req", "request nil", "400")
	}

	fechaInicial, err := time.Parse(dateLayout, req.FechaInicial)
	if err != nil {
		return nil, errorCtrl.Error("GenerarReporteElementos - time.Parse(fecha_inicial)", err, "400")
	}

	fechaFinal, err := time.Parse(dateLayout, req.FechaFinal)
	if err != nil {
		return nil, errorCtrl.Error("GenerarReporteElementos - time.Parse(fecha_final)", err, "400")
	}

	if fechaFinal.Before(fechaInicial) {
		return nil, errorCtrl.Error("GenerarReporteElementos - rango_fechas", "fecha_final debe ser mayor o igual a fecha_inicial", "400")
	}

	elementos, outputError := consultarElementosReporte(fechaInicial, fechaFinal)
	if outputError != nil {
		return nil, outputError
	}

	archivo := xlsx.NewFile()
	hoja, err := archivo.AddSheet(sheetName)
	if err != nil {
		return nil, errorCtrl.Error("GenerarReporteElementos - archivo.AddSheet", err, "500")
	}

	filaEncabezado := hoja.AddRow()
	for _, header := range reporteElementosHeaders {
		filaEncabezado.AddCell().SetString(header)
	}

	for _, elemento := range elementos {
		addElementoRow(hoja, elemento)
	}

	setColumnWidths(hoja)

	var buffer bytes.Buffer
	if err := archivo.Write(&buffer); err != nil {
		return nil, errorCtrl.Error("GenerarReporteElementos - archivo.Write", err, "500")
	}

	respuesta = &models.ReporteExcelBase64Response{
		ArchivoBase64: base64.StdEncoding.EncodeToString(buffer.Bytes()),
		NombreArchivo: fmt.Sprintf("reporte_elementos_%s_%s.xlsx", fechaInicial.Format("20060102"), fechaFinal.Format("20060102")),
		TipoArchivo:   "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}

	return respuesta, nil
}

func consultarElementosReporteDefault(fechaInicial, fechaFinal time.Time) (elementos []*models.DetalleElemento, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarElementosReporteDefault - Unhandled Error!", "500")

	fechaInicial = time.Date(fechaInicial.Year(), fechaInicial.Month(), fechaInicial.Day(), 0, 0, 0, 0, time.UTC)
	fechaFinal = time.Date(fechaFinal.Year(), fechaFinal.Month(), fechaFinal.Day(), 23, 59, 59, 0, time.UTC)

	query := "Activo:true,FechaCreacion__gte:" + fechaInicial.Format(time.RFC3339) +
		",FechaCreacion__lte:" + fechaFinal.Format(time.RFC3339)

	elementosBase, outputError := crudActaRecibido.GetAllElemento(query, "Id,FechaCreacion", "FechaCreacion", "asc", "", "-1")
	if outputError != nil {
		return nil, outputError
	}
	if len(elementosBase) == 0 {
		return []*models.DetalleElemento{}, nil
	}

	ids := make([]int, 0, len(elementosBase))
	for _, elemento := range elementosBase {
		if elemento != nil && elemento.Id > 0 {
			ids = append(ids, elemento.Id)
		}
	}

	if len(ids) == 0 {
		return []*models.DetalleElemento{}, nil
	}

	elementos, outputError = actaRecibido.GetElementos(0, ids)
	if outputError != nil {
		return nil, outputError
	}

	return ordenarElementosPorIds(ids, elementos), nil
}

func ordenarElementosPorIds(ids []int, elementos []*models.DetalleElemento) []*models.DetalleElemento {
	elementosPorId := make(map[int]*models.DetalleElemento)
	for _, elemento := range elementos {
		if elemento != nil {
			elementosPorId[elemento.Id] = elemento
		}
	}

	ordenados := make([]*models.DetalleElemento, 0, len(elementos))
	for _, id := range ids {
		if elemento, ok := elementosPorId[id]; ok {
			ordenados = append(ordenados, elemento)
		}
	}

	return ordenados
}

func addElementoRow(hoja *xlsx.Sheet, elemento *models.DetalleElemento) {
	if elemento == nil {
		return
	}

	row := hoja.AddRow()
	values := []string{
		strconv.Itoa(elemento.Id),
		elemento.Nombre,
		strconv.Itoa(elemento.Cantidad),
		elemento.Marca,
		elemento.Serie,
		strconv.Itoa(elemento.UnidadMedida),
		floatCell(elemento.ValorUnitario),
		floatCell(elemento.Subtotal),
		floatCell(elemento.Descuento),
		floatCell(elemento.ValorTotal),
		strconv.Itoa(elemento.PorcentajeIvaId),
		floatCell(elemento.ValorIva),
		floatCell(elemento.ValorFinal),
		jsonCell(elemento.SubgrupoCatalogoId),
		jsonCell(elemento.TipoBienId),
		jsonCell(elemento.EstadoElementoId),
		jsonCell(elemento.ActaRecibidoId),
		elemento.Placa,
		strconv.FormatBool(elemento.Activo),
		timeCell(elemento.FechaCreacion),
		timeCell(elemento.FechaModificacion),
	}

	for _, value := range values {
		row.AddCell().SetString(value)
	}
}

func floatCell(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func timeCell(value time.Time) string {
	if value.IsZero() {
		return ""
	}

	return value.Format(dateTimeLayout)
}

func jsonCell(value interface{}) string {
	if value == nil {
		return ""
	}

	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%v", value)
	}

	return string(encoded)
}

func setColumnWidths(hoja *xlsx.Sheet) {
	_ = hoja.SetColWidth(0, 0, 10)
	_ = hoja.SetColWidth(1, 4, 22)
	_ = hoja.SetColWidth(5, 12, 16)
	_ = hoja.SetColWidth(13, 16, 48)
	_ = hoja.SetColWidth(17, 20, 22)
}
