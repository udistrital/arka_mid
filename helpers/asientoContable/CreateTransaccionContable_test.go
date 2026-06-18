package asientoContable

import (
	"testing"
	"time"

	"github.com/udistrital/arka_mid/models"
)

func TestCreateTransaccionContablePreservaFechaExistente(t *testing.T) {
	t.Parallel()

	fecha := time.Date(2024, 12, 31, 15, 0, 0, 0, time.UTC)
	transaccion := &models.TransaccionMovimientos{
		FechaTransaccion: fecha,
	}

	original := getComprobanteCreateTransaccionContable
	getComprobanteCreateTransaccionContable = func(tipoComprobante string, comprobanteID *string) map[string]interface{} {
		*comprobanteID = "cmp-1"
		return nil
	}
	defer func() {
		getComprobanteCreateTransaccionContable = original
	}()

	msg, err := CreateTransaccionContable("P8", "Entrada Almacén", transaccion)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg != "" {
		t.Fatalf("unexpected message: %q", msg)
	}
	if !transaccion.FechaTransaccion.Equal(fecha) {
		t.Fatalf("expected fecha to be preserved, got %v", transaccion.FechaTransaccion)
	}
}
