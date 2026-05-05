package reportesHelper

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/tealeg/xlsx"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

const (
	dateLayout   = "2006-01-02"
	sheetName    = "Reporte"
	headerColumn = "elementos"
)

// GenerarReporteElementos genera un archivo Excel en base64 con la estructura
// inicial del reporte de elementos para el rango de fechas solicitado.
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

	archivo := xlsx.NewFile()
	hoja, err := archivo.AddSheet(sheetName)
	if err != nil {
		return nil, errorCtrl.Error("GenerarReporteElementos - archivo.AddSheet", err, "500")
	}

	filaEncabezado := hoja.AddRow()
	filaEncabezado.AddCell().SetString(headerColumn)

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
