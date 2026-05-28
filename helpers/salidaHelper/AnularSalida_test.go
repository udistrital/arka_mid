package salidaHelper

import (
	"testing"

	"github.com/udistrital/arka_mid/models"
)

func TestHistorialPermiteAnularSalida(t *testing.T) {
	t.Parallel()

	if ok, _ := historialPermiteAnularSalida(nil, 10); ok {
		t.Fatal("expected nil historial to be rejected")
	}

	historialBase := &models.Historial{
		Salida: &models.Movimiento{Id: 10},
	}
	if ok, msg := historialPermiteAnularSalida(historialBase, 10); !ok || msg != "" {
		t.Fatalf("expected historial base to be valid, got ok=%v msg=%q", ok, msg)
	}

	historialConTraslado := &models.Historial{
		Salida:    &models.Movimiento{Id: 10},
		Traslados: []*models.Movimiento{{Id: 20}},
	}
	if ok, _ := historialPermiteAnularSalida(historialConTraslado, 10); ok {
		t.Fatal("expected traslado posterior to block annulment")
	}

	historialConBaja := &models.Historial{
		Salida: &models.Movimiento{Id: 10},
		Baja:   &models.Movimiento{Id: 30},
	}
	if ok, _ := historialPermiteAnularSalida(historialConBaja, 10); ok {
		t.Fatal("expected baja posterior to block annulment")
	}

	historialConNovedad := &models.Historial{
		Salida:    &models.Movimiento{Id: 10},
		Novedades: []models.NovedadElemento{{Id: 40}},
	}
	if ok, _ := historialPermiteAnularSalida(historialConNovedad, 10); ok {
		t.Fatal("expected novedad posterior to block annulment")
	}

	historialOtraSalida := &models.Historial{
		Salida: &models.Movimiento{Id: 11},
	}
	if ok, _ := historialPermiteAnularSalida(historialOtraSalida, 10); ok {
		t.Fatal("expected different current salida to block annulment")
	}
}

func TestDebeRestaurarEntradaAprobada(t *testing.T) {
	t.Parallel()

	salidasSoloActual := []*models.Movimiento{
		{Id: 1, EstadoMovimientoId: &models.EstadoMovimiento{Nombre: estadoSalidaAprobada}},
	}
	if !debeRestaurarEntradaAprobada(salidasSoloActual, 1) {
		t.Fatal("expected single current salida to restore entry state")
	}

	salidasRestantesAnuladas := []*models.Movimiento{
		{Id: 1, EstadoMovimientoId: &models.EstadoMovimiento{Nombre: estadoSalidaAprobada}},
		{Id: 2, EstadoMovimientoId: &models.EstadoMovimiento{Nombre: estadoSalidaAnulada}},
		{Id: 3, EstadoMovimientoId: &models.EstadoMovimiento{Nombre: estadoSalidaRechazada}},
	}
	if !debeRestaurarEntradaAprobada(salidasRestantesAnuladas, 1) {
		t.Fatal("expected only annulled/rejected sibling salidas to restore entry state")
	}

	salidasConOtraActiva := []*models.Movimiento{
		{Id: 1, EstadoMovimientoId: &models.EstadoMovimiento{Nombre: estadoSalidaAprobada}},
		{Id: 2, EstadoMovimientoId: &models.EstadoMovimiento{Nombre: "Salida En Trámite"}},
	}
	if debeRestaurarEntradaAprobada(salidasConOtraActiva, 1) {
		t.Fatal("did not expect another active salida to restore entry state")
	}
}

func TestNormalizarTerceroIdSalida(t *testing.T) {
	t.Parallel()

	if normalizarTerceroIdSalida(nil) != nil {
		t.Fatal("expected nil tercero for nil pointer")
	}

	cero := 0
	if normalizarTerceroIdSalida(&cero) != nil {
		t.Fatal("expected nil tercero for zero value")
	}

	negativo := -2
	if normalizarTerceroIdSalida(&negativo) != nil {
		t.Fatal("expected nil tercero for negative value")
	}

	valido := 55
	tercero := normalizarTerceroIdSalida(&valido)
	if tercero == nil || *tercero != valido {
		t.Fatalf("expected tercero id %d to be preserved", valido)
	}
}

func TestSincronizarEstadoMovimiento(t *testing.T) {
	t.Parallel()

	estado := &models.EstadoMovimiento{
		Id:     3,
		Nombre: estadoEntradaConSalida,
	}

	sincronizarEstadoMovimiento(estado, 7, estadoEntradaAprobada)

	if estado.Id != 7 {
		t.Fatalf("expected estado id 7, got %d", estado.Id)
	}

	if estado.Nombre != estadoEntradaAprobada {
		t.Fatalf("expected estado nombre %q, got %q", estadoEntradaAprobada, estado.Nombre)
	}

	sincronizarEstadoMovimiento(nil, 9, estadoSalidaAnulada)
}
