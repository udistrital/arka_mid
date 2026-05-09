package reportesHelper

import (
	"bytes"
	"encoding/base64"
	"testing"
	"time"

	"github.com/tealeg/xlsx"
	"github.com/udistrital/arka_mid/models"
)

func TestGenerarReporteElementos(t *testing.T) {
	mockConsultarElementosReporte(t, []*models.DetalleElemento{
		{
			Id:            101,
			Nombre:        "Elemento Uno",
			Cantidad:      2,
			Marca:         "Marca A",
			Serie:         "SERIE-001",
			UnidadMedida:  1,
			ValorUnitario: 1500,
			Subtotal:      3000,
			ValorTotal:    3000,
			Activo:        true,
			FechaCreacion: time.Date(2026, 1, 10, 8, 30, 0, 0, time.UTC),
		},
	})

	respuesta, err := GenerarReporteElementos(&models.ReporteFechasRequest{
		FechaInicial: "2026-01-01",
		FechaFinal:   "2026-01-31",
	})
	if err != nil {
		t.Fatalf("GenerarReporteElementos retornó error: %v", err)
	}

	if respuesta == nil {
		t.Fatal("GenerarReporteElementos retornó respuesta nil")
	}

	if respuesta.ArchivoBase64 == "" {
		t.Fatal("ArchivoBase64 no debe ser vacío")
	}

	contenido, decodeErr := base64.StdEncoding.DecodeString(respuesta.ArchivoBase64)
	if decodeErr != nil {
		t.Fatalf("base64 inválido: %v", decodeErr)
	}

	archivo, openErr := xlsx.OpenBinary(contenido)
	if openErr != nil {
		t.Fatalf("no se pudo abrir el excel generado: %v", openErr)
	}

	if len(archivo.Sheets) != 1 {
		t.Fatalf("se esperaba una hoja, se obtuvieron %d", len(archivo.Sheets))
	}

	if len(archivo.Sheets[0].Rows) == 0 {
		t.Fatal("se esperaba al menos una fila de encabezado")
	}

	valor := archivo.Sheets[0].Rows[0].Cells[0].String()
	if valor != "Id" {
		t.Fatalf("encabezado inesperado: %q", valor)
	}

	if len(archivo.Sheets[0].Rows[0].Cells) != len(reporteElementosHeaders) {
		t.Fatalf("número de columnas inesperado: %d", len(archivo.Sheets[0].Rows[0].Cells))
	}

	if len(archivo.Sheets[0].Rows) != 2 {
		t.Fatalf("se esperaban encabezado y una fila de datos, se obtuvieron %d filas", len(archivo.Sheets[0].Rows))
	}

	if nombre := archivo.Sheets[0].Rows[1].Cells[1].String(); nombre != "Elemento Uno" {
		t.Fatalf("nombre de elemento inesperado: %q", nombre)
	}
}

func TestGenerarReporteElementosFechaFinalMenor(t *testing.T) {
	respuesta, err := GenerarReporteElementos(&models.ReporteFechasRequest{
		FechaInicial: "2026-02-01",
		FechaFinal:   "2026-01-31",
	})
	if err == nil {
		t.Fatal("se esperaba error cuando fecha_final es menor a fecha_inicial")
	}

	if respuesta != nil {
		t.Fatal("no se esperaba respuesta cuando el rango es inválido")
	}
}

func TestExcelGeneradoEsBinarioValido(t *testing.T) {
	mockConsultarElementosReporte(t, []*models.DetalleElemento{})

	respuesta, err := GenerarReporteElementos(&models.ReporteFechasRequest{
		FechaInicial: "2026-03-01",
		FechaFinal:   "2026-03-15",
	})
	if err != nil {
		t.Fatalf("GenerarReporteElementos retornó error: %v", err)
	}

	contenido, decodeErr := base64.StdEncoding.DecodeString(respuesta.ArchivoBase64)
	if decodeErr != nil {
		t.Fatalf("base64 inválido: %v", decodeErr)
	}

	if len(bytes.TrimSpace(contenido)) == 0 {
		t.Fatal("el excel generado no debe ser vacío")
	}
}

func mockConsultarElementosReporte(t *testing.T, elementos []*models.DetalleElemento) {
	t.Helper()

	original := consultarElementosReporte
	consultarElementosReporte = func(fechaInicial, fechaFinal time.Time) ([]*models.DetalleElemento, map[string]interface{}) {
		return elementos, nil
	}

	t.Cleanup(func() {
		consultarElementosReporte = original
	})
}
