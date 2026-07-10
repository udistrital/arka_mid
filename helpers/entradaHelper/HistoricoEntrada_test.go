package entradaHelper

import (
	"strings"
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
		TransaccionEntrada: models.TransaccionEntrada{
			FormatoTipoMovimientoId: "ENT_ADQ",
			Detalle: models.FormatoBaseEntrada{
				ActaRecibidoId: 3348,
			},
		},
		ConsecutivoId: 10,
		Year:          2024,
		FechaCreacion: time.Date(2024, 12, 31, 8, 0, 0, 0, time.UTC),
		FechaCorte:    time.Date(2024, 12, 31, 9, 0, 0, 0, time.UTC),
	}
	if err := validarAprobacionHistorica(payload); err != nil {
		t.Fatalf("expected payload to be valid, got %v", err)
	}
}

func TestValidarAprobacionHistoricaErroresDetallados(t *testing.T) {
	t.Parallel()

	err := validarAprobacionHistorica(&models.TransaccionEntradaHistorica{})
	if err == nil {
		t.Fatal("expected validation error")
	}

	msg := err.Error()
	expected := []string{
		"FormatoTipoMovimientoId es obligatorio",
		"ConsecutivoId es obligatorio y debe ser mayor a 0",
		"Year es obligatorio y debe ser mayor a 0",
		"FechaCreacion es obligatoria y debe venir en formato RFC3339",
		"FechaCorte es obligatoria y debe venir en formato RFC3339",
		"Detalle.acta_recibido_id o Detalle.elementos es obligatorio",
	}

	for _, fragment := range expected {
		if !strings.Contains(msg, fragment) {
			t.Fatalf("expected validation error to contain %q, got %q", fragment, msg)
		}
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

func TestNormalizarFechaCorteHistorica(t *testing.T) {
	t.Parallel()

	creacion := time.Date(2024, 12, 30, 8, 0, 0, 0, time.UTC)

	payload := &models.TransaccionEntradaHistorica{
		FechaCreacion: creacion,
	}

	if got := normalizarFechaCorteHistorica(payload); !got.Equal(creacion) {
		t.Fatalf("expected fecha creacion fallback, got %v", got)
	}
}

func TestAplicarConsecutivoHistoricoEntradaConsecutivoNoEncontrado(t *testing.T) {
	t.Parallel()

	original := getConsecutivoByIDEntradaHistorica
	getConsecutivoByIDEntradaHistorica = func(id int, consecutivo *models.Consecutivo) map[string]interface{} {
		return nil
	}
	defer func() {
		getConsecutivoByIDEntradaHistorica = original
	}()

	movimiento := &models.Movimiento{}
	err := aplicarConsecutivoHistoricoEntrada(movimiento, 10764, 1997)
	if err == nil {
		t.Fatal("expected error for missing consecutivo")
	}

	data, ok := err["err"].(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected error payload: %#v", err)
	}
	if got, ok := data["detalle"].(string); !ok || !strings.Contains(got, "ConsecutivoId 10764 no existe o no está disponible") {
		t.Fatalf("unexpected error payload: %#v", err)
	}
	if got, ok := data["paso"].(string); !ok || got != "validar consecutivo histórico" {
		t.Fatalf("unexpected paso payload: %#v", err)
	}
	if got, ok := err["status"].(string); !ok || got != "404" {
		t.Fatalf("expected status 404, got %#v", err["status"])
	}
}

func TestMensajeEstadoActaHistoricaInvalido(t *testing.T) {
	t.Parallel()

	msg := mensajeEstadoActaHistoricaInvalido(3348, &models.HistoricoActa{
		EstadoActaId: &models.EstadoActa{
			Nombre:            "En revisión",
			CodigoAbreviacion: "REV",
		},
	})

	if !strings.Contains(msg, `acta 3348`) || !strings.Contains(msg, `"En revisión"`) || !strings.Contains(msg, "Aceptada") {
		t.Fatalf("unexpected message: %q", msg)
	}
}

func TestRegistrarEntradaHistoricaEnviaFechasAntesDelPost(t *testing.T) {
	creacion := time.Date(1997, 1, 15, 8, 0, 0, 0, time.UTC)
	corte := time.Date(1997, 1, 16, 9, 0, 0, 0, time.UTC)

	originalGetConsecutivo := getConsecutivoByIDEntradaHistorica
	originalGetActa := getTransaccionActaRecibidoEntradaHistorica
	originalPutActa := putTransaccionActaRecibidoEntradaHistorica
	originalGetEstado := getEstadoMovimientoIdByNombreEntradaHistorica
	originalGetFormato := getFormatoTipoMovimientoIdByCodigoEntradaHistorica
	originalPostMovimiento := postMovimientoEntradaHistorica
	originalPutMovimiento := putMovimientoEntradaHistorica
	originalPostSoporte := postSoporteMovimientoEntradaHistorica
	originalAprobar := aprobarEntradaHistoricaFn

	getEstadoMovimientoIdByNombreEntradaHistorica = func(id *int, nombre string) map[string]interface{} {
		*id = 2
		return nil
	}
	getFormatoTipoMovimientoIdByCodigoEntradaHistorica = func(id *int, codigo string) map[string]interface{} {
		*id = 1
		return nil
	}
	getConsecutivoByIDEntradaHistorica = func(id int, consecutivo *models.Consecutivo) map[string]interface{} {
		consecutivo.Id = id
		consecutivo.Consecutivo = 1
		return nil
	}
	getTransaccionActaRecibidoEntradaHistorica = func(id int, elementos bool, transaccion *models.TransaccionActaRecibido) map[string]interface{} {
		transaccion.UltimoEstado = &models.HistoricoActa{
			EstadoActaId: &models.EstadoActa{CodigoAbreviacion: "Aceptada"},
		}
		transaccion.Elementos = []*models.Elemento{{Id: 1}}
		return nil
	}
	putTransaccionActaRecibidoEntradaHistorica = func(id int, transaccion *models.TransaccionActaRecibido) map[string]interface{} {
		return nil
	}
	postMovimientoEntradaHistorica = func(movimiento *models.Movimiento) map[string]interface{} {
		if movimiento == nil {
			t.Fatal("expected movimiento payload")
		}
		if movimiento.FechaCreacion.IsZero() {
			t.Fatal("expected FechaCreacion to be set before post")
		}
		if movimiento.FechaModificacion.IsZero() {
			t.Fatal("expected FechaModificacion to be set before post")
		}
		if !movimiento.FechaCreacion.Equal(creacion) {
			t.Fatalf("unexpected FechaCreacion: %v", movimiento.FechaCreacion)
		}
		if !movimiento.FechaModificacion.Equal(corte) {
			t.Fatalf("unexpected FechaModificacion: %v", movimiento.FechaModificacion)
		}
		if movimiento.FechaCorte == nil {
			t.Fatal("expected FechaCorte to be set before post")
		}
		if !movimiento.FechaCorte.Equal(corte) {
			t.Fatalf("unexpected FechaCorte: %v", movimiento.FechaCorte)
		}
		movimiento.Id = 99
		return nil
	}
	putMovimientoEntradaHistorica = func(movimiento *models.Movimiento, id int) map[string]interface{} {
		return nil
	}
	postSoporteMovimientoEntradaHistorica = func(soporte *models.SoporteMovimiento) map[string]interface{} {
		return nil
	}
	aprobarEntradaHistoricaFn = func(entradaId int, data *models.TransaccionEntradaHistorica, resultado *models.ResultadoMovimiento) map[string]interface{} {
		return nil
	}

	t.Cleanup(func() {
		getConsecutivoByIDEntradaHistorica = originalGetConsecutivo
		getTransaccionActaRecibidoEntradaHistorica = originalGetActa
		putTransaccionActaRecibidoEntradaHistorica = originalPutActa
		getEstadoMovimientoIdByNombreEntradaHistorica = originalGetEstado
		getFormatoTipoMovimientoIdByCodigoEntradaHistorica = originalGetFormato
		postMovimientoEntradaHistorica = originalPostMovimiento
		putMovimientoEntradaHistorica = originalPutMovimiento
		postSoporteMovimientoEntradaHistorica = originalPostSoporte
		aprobarEntradaHistoricaFn = originalAprobar
	})

	data := &models.TransaccionEntradaHistorica{
		TransaccionEntrada: models.TransaccionEntrada{
			Observacion:             "Migración histórica",
			FormatoTipoMovimientoId: "ENT_ADQ",
			Detalle: models.FormatoBaseEntrada{
				ActaRecibidoId: 3353,
				Factura:        2203,
			},
		},
		ConsecutivoId: 10764,
		Year:          1997,
		FechaCreacion: creacion,
		FechaCorte:    corte,
	}

	resultado := &models.ResultadoMovimiento{}
	err := RegistrarEntradaHistorica(data, resultado)
	if err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}
}

func TestRegistrarEntradaHistoricaConsultaActaConElementos(t *testing.T) {
	originalGetConsecutivo := getConsecutivoByIDEntradaHistorica
	originalGetActa := getTransaccionActaRecibidoEntradaHistorica
	originalPutActa := putTransaccionActaRecibidoEntradaHistorica
	originalGetEstado := getEstadoMovimientoIdByNombreEntradaHistorica
	originalGetFormato := getFormatoTipoMovimientoIdByCodigoEntradaHistorica
	originalPostMovimiento := postMovimientoEntradaHistorica
	originalPutMovimiento := putMovimientoEntradaHistorica
	originalPostSoporte := postSoporteMovimientoEntradaHistorica
	originalAprobar := aprobarEntradaHistoricaFn

	getEstadoMovimientoIdByNombreEntradaHistorica = func(id *int, nombre string) map[string]interface{} {
		*id = 2
		return nil
	}
	getFormatoTipoMovimientoIdByCodigoEntradaHistorica = func(id *int, codigo string) map[string]interface{} {
		*id = 1
		return nil
	}
	getConsecutivoByIDEntradaHistorica = func(id int, consecutivo *models.Consecutivo) map[string]interface{} {
		consecutivo.Id = id
		consecutivo.Consecutivo = 1
		return nil
	}
	getTransaccionActaRecibidoEntradaHistorica = func(id int, elementos bool, transaccion *models.TransaccionActaRecibido) map[string]interface{} {
		if !elementos {
			t.Fatal("expected consultar acta histórica con elementos=true")
		}
		transaccion.UltimoEstado = &models.HistoricoActa{
			EstadoActaId: &models.EstadoActa{CodigoAbreviacion: "Aceptada"},
		}
		transaccion.Elementos = []*models.Elemento{{Id: 1}}
		return nil
	}
	putTransaccionActaRecibidoEntradaHistorica = func(id int, transaccion *models.TransaccionActaRecibido) map[string]interface{} {
		return nil
	}
	postMovimientoEntradaHistorica = func(movimiento *models.Movimiento) map[string]interface{} {
		movimiento.Id = 99
		return nil
	}
	putMovimientoEntradaHistorica = func(movimiento *models.Movimiento, id int) map[string]interface{} {
		return nil
	}
	postSoporteMovimientoEntradaHistorica = func(soporte *models.SoporteMovimiento) map[string]interface{} {
		return nil
	}
	aprobarEntradaHistoricaFn = func(entradaId int, data *models.TransaccionEntradaHistorica, resultado *models.ResultadoMovimiento) map[string]interface{} {
		return nil
	}

	t.Cleanup(func() {
		getConsecutivoByIDEntradaHistorica = originalGetConsecutivo
		getTransaccionActaRecibidoEntradaHistorica = originalGetActa
		putTransaccionActaRecibidoEntradaHistorica = originalPutActa
		getEstadoMovimientoIdByNombreEntradaHistorica = originalGetEstado
		getFormatoTipoMovimientoIdByCodigoEntradaHistorica = originalGetFormato
		postMovimientoEntradaHistorica = originalPostMovimiento
		putMovimientoEntradaHistorica = originalPutMovimiento
		postSoporteMovimientoEntradaHistorica = originalPostSoporte
		aprobarEntradaHistoricaFn = originalAprobar
	})

	data := &models.TransaccionEntradaHistorica{
		TransaccionEntrada: models.TransaccionEntrada{
			Observacion:             "Migración histórica",
			FormatoTipoMovimientoId: "ENT_ADQ",
			Detalle: models.FormatoBaseEntrada{
				ActaRecibidoId: 3353,
			},
		},
		ConsecutivoId: 10764,
		Year:          1997,
		FechaCreacion: time.Date(1997, 1, 15, 8, 0, 0, 0, time.UTC),
		FechaCorte:    time.Date(1997, 1, 16, 9, 0, 0, 0, time.UTC),
	}

	resultado := &models.ResultadoMovimiento{}
	if err := RegistrarEntradaHistorica(data, resultado); err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}
}
