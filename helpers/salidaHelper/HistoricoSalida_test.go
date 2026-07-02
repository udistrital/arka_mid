package salidaHelper

import (
	"testing"
	"time"

	"github.com/udistrital/arka_mid/models"
)

func TestValidarAprobacionHistoricaSalida(t *testing.T) {
	t.Parallel()

	if err := validarAprobacionHistoricaSalida(nil); err == nil {
		t.Fatal("expected error for nil payload")
	}

	payload := &models.TransaccionSalidaHistorica{
		ConsecutivoId: 10,
		Year:          2024,
		FechaCreacion: time.Date(2024, 12, 31, 8, 0, 0, 0, time.UTC),
		FechaCorte:    time.Date(2024, 12, 31, 9, 0, 0, 0, time.UTC),
		SalidaGeneral: models.SalidaGeneral{
			Salidas: []models.TrSalida{
				{Salida: &models.Movimiento{MovimientoPadreId: &models.Movimiento{Id: 1}, FormatoTipoMovimientoId: &models.FormatoTipoMovimiento{Id: 7}}},
			},
		},
	}
	if err := validarAprobacionHistoricaSalida(payload); err != nil {
		t.Fatalf("expected payload to be valid, got %v", err)
	}
}

func TestStringHistoricoConsecutivoSalida(t *testing.T) {
	t.Parallel()

	consecutivo := stringHistoricoConsecutivoSalida(&models.Consecutivo{
		Id:          99,
		Consecutivo: 665,
		Year:        2024,
	})
	if consecutivo == nil {
		t.Fatal("expected formatted string")
	}
	if *consecutivo != "H21-00665-2024" {
		t.Fatalf("unexpected consecutivo: %q", *consecutivo)
	}
}

func TestNormalizarFechaModificacionHistoricaSalida(t *testing.T) {
	t.Parallel()

	creacion := time.Date(2024, 12, 30, 8, 0, 0, 0, time.UTC)
	corte := time.Date(2024, 12, 31, 10, 0, 0, 0, time.UTC)

	payload := &models.TransaccionSalidaHistorica{
		FechaCreacion: creacion,
		FechaCorte:    corte,
	}

	if got := normalizarFechaModificacionHistoricaSalida(payload); !got.Equal(corte) {
		t.Fatalf("expected fecha corte fallback, got %v", got)
	}
}
