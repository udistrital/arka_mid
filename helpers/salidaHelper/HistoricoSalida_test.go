package salidaHelper

import (
	"strings"
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

func TestNormalizarElementoMovimientoHistoricoSalida(t *testing.T) {
	t.Parallel()

	elementoActaID := 46676
	movimientoID := 8192
	creacion := time.Date(1997, 12, 31, 8, 0, 0, 0, time.UTC)
	corte := time.Date(1997, 12, 31, 10, 0, 0, 0, time.UTC)

	elemento := &models.ElementosMovimiento{
		Id:                 41409,
		ElementoActaId:     &elementoActaID,
		ElementoCatalogoId: 0,
		Unidad:             1,
		ValorUnitario:      100,
		ValorTotal:         100,
		SaldoCantidad:      0,
		SaldoValor:         0,
		VidaUtil:           0,
		ValorResidual:      0,
		Activo:             false,
		MovimientoId:       &models.Movimiento{Id: movimientoID, Observacion: "debe limpiarse"},
	}
	payload := &models.TransaccionSalidaHistorica{
		FechaCreacion: creacion,
		FechaCorte:    corte,
	}

	got := normalizarElementoMovimientoHistoricoSalida(elemento, payload, nil, nil)
	if got == nil {
		t.Fatal("expected normalized element")
	}
	if got.Id != elemento.Id {
		t.Fatalf("unexpected Id: %d", got.Id)
	}
	if got.MovimientoId == nil || got.MovimientoId.Id != movimientoID {
		t.Fatalf("unexpected MovimientoId: %+v", got.MovimientoId)
	}
	if got.MovimientoId.Observacion != "" {
		t.Fatalf("expected sanitized MovimientoId, got %+v", got.MovimientoId)
	}
	if got.ElementoActaId == nil || *got.ElementoActaId != elementoActaID {
		t.Fatalf("unexpected ElementoActaId: %+v", got.ElementoActaId)
	}
	if !got.FechaCreacion.Equal(creacion) {
		t.Fatalf("unexpected FechaCreacion: %v", got.FechaCreacion)
	}
	if !got.FechaModificacion.Equal(corte) {
		t.Fatalf("unexpected FechaModificacion: %v", got.FechaModificacion)
	}
}

func TestNormalizarElementoMovimientoHistoricoSalidaSinMovimientoId(t *testing.T) {
	t.Parallel()

	elementoActaID := 46676
	creacion := time.Date(1997, 12, 31, 8, 0, 0, 0, time.UTC)
	corte := time.Date(1997, 12, 31, 10, 0, 0, 0, time.UTC)

	elemento := &models.ElementosMovimiento{
		Id:                 41412,
		ElementoActaId:     &elementoActaID,
		ElementoCatalogoId: 0,
		Activo:             false,
	}
	payload := &models.TransaccionSalidaHistorica{
		FechaCreacion: creacion,
		FechaCorte:    corte,
	}

	got := normalizarElementoMovimientoHistoricoSalida(elemento, payload, nil, nil)
	if got == nil {
		t.Fatal("expected normalized element")
	}
	if got.MovimientoId != nil {
		t.Fatalf("expected MovimientoId to stay nil in normalizer, got %+v", got.MovimientoId)
	}
	if got.ElementoActaId == nil || *got.ElementoActaId != elementoActaID {
		t.Fatalf("unexpected ElementoActaId: %+v", got.ElementoActaId)
	}
}

func TestHidratarElementoMovimientoHistoricoSalidaCompletaValoresFaltantes(t *testing.T) {
	t.Parallel()

	elemento := &models.ElementosMovimiento{
		Id:     41412,
		Activo: true,
	}
	seed := &elementoMovimientoHistoricoSalidaSeed{
		ValorResidual: 15,
	}
	elementoActa := &models.DetalleElemento{
		Cantidad:      2,
		ValorUnitario: 1250,
		ValorTotal:    2500,
		SubgrupoCatalogoId: &models.DetalleSubgrupo{
			VidaUtil:      5,
			ValorResidual: 0.1,
		},
	}

	hidratarElementoMovimientoHistoricoSalida(elemento, seed, elementoActa)

	if elemento.Unidad != 2 {
		t.Fatalf("unexpected Unidad: %v", elemento.Unidad)
	}
	if elemento.ValorUnitario != 1250 {
		t.Fatalf("unexpected ValorUnitario: %v", elemento.ValorUnitario)
	}
	if elemento.ValorTotal != 2500 {
		t.Fatalf("unexpected ValorTotal: %v", elemento.ValorTotal)
	}
	if elemento.SaldoCantidad != 2 {
		t.Fatalf("unexpected SaldoCantidad: %v", elemento.SaldoCantidad)
	}
	if elemento.SaldoValor != 2500 {
		t.Fatalf("unexpected SaldoValor: %v", elemento.SaldoValor)
	}
	if elemento.VidaUtil != 5 {
		t.Fatalf("unexpected VidaUtil: %v", elemento.VidaUtil)
	}
	if elemento.ValorResidual != 15 {
		t.Fatalf("expected seed ValorResidual to be preserved, got %v", elemento.ValorResidual)
	}
}

func TestResolverValorResidualHistoricoSalidaUsaProporcion(t *testing.T) {
	t.Parallel()

	if got := resolverValorResidualHistoricoSalida(0.1, 2500); got != 250 {
		t.Fatalf("unexpected proportional residual: %v", got)
	}
	if got := resolverValorResidualHistoricoSalida(300, 2500); got != 300 {
		t.Fatalf("unexpected absolute residual: %v", got)
	}
}

func TestWrapPostTrSalidaHistoricaErrorDuplicateElementoActa(t *testing.T) {
	t.Parallel()

	elementoActaID := 46676
	got := wrapPostTrSalidaHistoricaError(&models.SalidaGeneral{
		Salidas: []models.TrSalida{
			{
				Salida: &models.Movimiento{
					MovimientoPadreId: &models.Movimiento{Id: 8188},
				},
				Elementos: []*models.ElementosMovimiento{
					{ElementoActaId: &elementoActaID},
				},
			},
		},
	}, map[string]interface{}{
		"err": `http 400: {"Data":"pq: duplicate key value violates unique constraint \"uq_elemento_acta_id\""}`,
	})

	if got == nil {
		t.Fatal("expected wrapped error")
	}
	if got["status"] != "409" {
		t.Fatalf("unexpected status: %#v", got["status"])
	}
	errPayload, ok := got["err"].(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected error payload: %#v", got)
	}
	msg, _ := errPayload["detalle"].(string)
	if !strings.Contains(msg, "46676") || !strings.Contains(msg, "8188") {
		t.Fatalf("unexpected error message: %q", msg)
	}
}
