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

func TestPostValidaSalidasVacias(t *testing.T) {
	t.Parallel()

	if _, err := Post(nil, false); err == nil {
		t.Fatal("expected nil request to be rejected")
	}

	if _, err := Post(&models.SalidaGeneral{}, false); err == nil {
		t.Fatal("expected empty salidas request to be rejected")
	}
}

func TestNormalizarNuevaSalida(t *testing.T) {
	t.Parallel()

	elementoActaID := 41801
	trSalida := &models.TrSalida{
		Salida: &models.Movimiento{
			Id:            25,
			Consecutivo:   normalizarString("H21-00001-2026"),
			ConsecutivoId: normalizarInt(101),
			MovimientoPadreId: &models.Movimiento{
				Id:     8079,
				Activo: false,
			},
			FormatoTipoMovimientoId: &models.FormatoTipoMovimiento{
				Id: 7,
			},
		},
		Elementos: []*models.ElementosMovimiento{
			{
				Id:             88,
				ElementoActaId: &elementoActaID,
				MovimientoId:   &models.Movimiento{Id: 25},
			},
		},
	}

	if err := normalizarNuevaSalida(trSalida, 3); err != nil {
		t.Fatalf("expected salida to be normalized, got %#v", err)
	}

	if trSalida.Salida.Id != 0 {
		t.Fatalf("expected salida id 0, got %d", trSalida.Salida.Id)
	}

	if trSalida.Salida.MovimientoPadreId == nil || trSalida.Salida.MovimientoPadreId.Id != 8079 {
		t.Fatalf("expected parent id 8079, got %+v", trSalida.Salida.MovimientoPadreId)
	}

	if trSalida.Salida.MovimientoPadreId.Activo {
		t.Fatal("expected parent payload to be reduced to id-only reference")
	}

	if trSalida.Salida.FormatoTipoMovimientoId == nil || trSalida.Salida.FormatoTipoMovimientoId.Id != 7 {
		t.Fatalf("expected format id 7, got %+v", trSalida.Salida.FormatoTipoMovimientoId)
	}

	if trSalida.Salida.EstadoMovimientoId == nil || trSalida.Salida.EstadoMovimientoId.Id != 3 || trSalida.Salida.EstadoMovimientoId.Nombre != "Salida En Trámite" {
		t.Fatalf("expected estado salida en tramite, got %+v", trSalida.Salida.EstadoMovimientoId)
	}

	if trSalida.Elementos[0].Id != 0 || trSalida.Elementos[0].MovimientoId != nil {
		t.Fatalf("expected element to be detached from previous movement, got %+v", trSalida.Elementos[0])
	}
}

func normalizarString(v string) *string { return &v }

func normalizarInt(v int) *int { return &v }
