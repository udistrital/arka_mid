package entradaHelper

import (
	"testing"

	"github.com/udistrital/arka_mid/models"
)

func TestEstadoActaPermiteAnulacion(t *testing.T) {
	t.Parallel()

	if !estadoActaPermiteAnulacion(estadoActaAsociadaEntrada) {
		t.Fatalf("expected %s to allow annulment", estadoActaAsociadaEntrada)
	}

	if !estadoActaPermiteAnulacion(estadoActaEnVerificacion) {
		t.Fatalf("expected %s to allow annulment retry", estadoActaEnVerificacion)
	}

	if estadoActaPermiteAnulacion("Aceptada") {
		t.Fatal("did not expect unrelated state to allow annulment")
	}
}

func TestAplicarEstadoActaEnVerificacionNoOpSiYaEstaEnVerificacion(t *testing.T) {
	t.Parallel()

	transaccion := &models.TransaccionActaRecibido{
		UltimoEstado: &models.HistoricoActa{
			Id: 99,
			EstadoActaId: &models.EstadoActa{
				Id:                7,
				CodigoAbreviacion: estadoActaEnVerificacion,
			},
		},
	}

	if err := aplicarEstadoActaEnVerificacion(transaccion); err != nil {
		t.Fatalf("expected nil error, got %#v", err)
	}

	if transaccion.UltimoEstado.Id != 99 {
		t.Fatalf("expected historico id to remain unchanged, got %d", transaccion.UltimoEstado.Id)
	}

	if transaccion.UltimoEstado.EstadoActaId.Id != 7 {
		t.Fatalf("expected estado id to remain unchanged, got %d", transaccion.UltimoEstado.EstadoActaId.Id)
	}
}

func TestNormalizarTerceroId(t *testing.T) {
	t.Parallel()

	if normalizarTerceroId(nil) != nil {
		t.Fatal("expected nil tercero for nil pointer")
	}

	terceroCero := 0
	if normalizarTerceroId(&terceroCero) != nil {
		t.Fatal("expected nil tercero for zero-value tercero id")
	}

	terceroNegativo := -4
	if normalizarTerceroId(&terceroNegativo) != nil {
		t.Fatal("expected nil tercero for negative tercero id")
	}

	terceroValido := 123
	tercero := normalizarTerceroId(&terceroValido)
	if tercero == nil || *tercero != terceroValido {
		t.Fatalf("expected tercero id %d to be preserved", terceroValido)
	}
}
