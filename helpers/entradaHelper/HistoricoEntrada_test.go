package entradaHelper

import (
	"testing"
	"time"

	"github.com/udistrital/arka_mid/models"
)

func TestValidarAprobacionHistorica(t *testing.T) {
	t.Parallel()

	if err := validarAprobacionHistorica(nil); err == nil {
		t.Fatal("expected error for nil payload")
	}

	payload := &models.TransaccionEntradaHistorica{
		ConsecutivoId: 10,
		Year:          2024,
		FechaCreacion: time.Date(2024, 12, 31, 8, 0, 0, 0, time.UTC),
		FechaCorte:    time.Date(2024, 12, 31, 9, 0, 0, 0, time.UTC),
	}
	if err := validarAprobacionHistorica(payload); err != nil {
		t.Fatalf("expected payload to be valid, got %v", err)
	}
}

func TestStringHistoricoConsecutivoEntrada(t *testing.T) {
	t.Parallel()

	consecutivo := stringHistoricoConsecutivoEntrada(&models.Consecutivo{
		Id:          99,
		Consecutivo: 649,
		Year:        2024,
	})
	if consecutivo == nil {
		t.Fatal("expected formatted string")
	}
	if *consecutivo != "P8-00649-2024" {
		t.Fatalf("unexpected consecutivo: %q", *consecutivo)
	}
}

func TestNormalizarFechaModificacionHistorica(t *testing.T) {
	t.Parallel()

	creacion := time.Date(2024, 12, 30, 8, 0, 0, 0, time.UTC)
	corte := time.Date(2024, 12, 31, 10, 0, 0, 0, time.UTC)

	payload := &models.TransaccionEntradaHistorica{
		FechaCreacion: creacion,
		FechaCorte:    corte,
	}

	if got := normalizarFechaModificacionHistorica(payload); !got.Equal(corte) {
		t.Fatalf("expected fecha corte fallback, got %v", got)
	}
}
