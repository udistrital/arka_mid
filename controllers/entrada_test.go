package controllers

import (
	"testing"
	"time"

	"github.com/udistrital/arka_mid/models"
)

func TestDecodeEntradaHistoricaRequestSoportaSnakeCase(t *testing.T) {
	t.Parallel()

	body := []byte(`{
		"Observacion": "Migracion historica",
		"FormatoTipoMovimientoId": "ENT_ADQ",
		"consecutivo_id": 10765,
		"year": 1997,
		"fecha_creacion": "1997-12-31T08:00:00Z",
		"fecha_modificacion": "2026-07-06T16:18:13Z",
		"fecha_corte": "1997-12-31T10:00:00Z"
	}`)

	var payload models.TransaccionEntradaHistorica
	if err := decodeEntradaHistoricaRequest(body, &payload); err != nil {
		t.Fatalf("decodeEntradaHistoricaRequest() error = %v", err)
	}

	if payload.ConsecutivoId != 10765 {
		t.Fatalf("expected consecutivo_id 10765, got %d", payload.ConsecutivoId)
	}
	if payload.Year != 1997 {
		t.Fatalf("expected year 1997, got %d", payload.Year)
	}

	fechaCreacion := time.Date(1997, 12, 31, 8, 0, 0, 0, time.UTC)
	if !payload.FechaCreacion.Equal(fechaCreacion) {
		t.Fatalf("unexpected fecha_creacion: %v", payload.FechaCreacion)
	}
}

func TestDecodeEntradaHistoricaRequestSoportaPascalCaseSwagger(t *testing.T) {
	t.Parallel()

	body := []byte(`{
		"Observacion": "Migracion historica",
		"FormatoTipoMovimientoId": "ENT_ADQ",
		"ConsecutivoId": 10765,
		"Year": 1997,
		"FechaCreacion": "1997-12-31T08:00:00Z",
		"FechaModificacion": "2026-07-06T16:18:13Z",
		"FechaCorte": "1997-12-31T10:00:00Z"
	}`)

	var payload models.TransaccionEntradaHistorica
	if err := decodeEntradaHistoricaRequest(body, &payload); err != nil {
		t.Fatalf("decodeEntradaHistoricaRequest() error = %v", err)
	}

	if payload.ConsecutivoId != 10765 {
		t.Fatalf("expected ConsecutivoId 10765, got %d", payload.ConsecutivoId)
	}
	if payload.Year != 1997 {
		t.Fatalf("expected Year 1997, got %d", payload.Year)
	}

	fechaCreacion := time.Date(1997, 12, 31, 8, 0, 0, 0, time.UTC)
	fechaCorte := time.Date(1997, 12, 31, 10, 0, 0, 0, time.UTC)

	if !payload.FechaCreacion.Equal(fechaCreacion) {
		t.Fatalf("unexpected FechaCreacion: %v", payload.FechaCreacion)
	}
	if !payload.FechaCorte.Equal(fechaCorte) {
		t.Fatalf("unexpected FechaCorte: %v", payload.FechaCorte)
	}
}
