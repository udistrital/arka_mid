package asientoContable

import (
	"strings"
	"testing"

	"github.com/udistrital/arka_mid/models"
)

func TestPayloadCuentasPrefiereSubtipoMovimiento(t *testing.T) {
	t.Parallel()

	payload := payloadCuentas(53504, 1, 7)

	if !strings.Contains(payload, "SubgrupoId__Id:53504") {
		t.Fatalf("payload sin subgrupo esperado: %s", payload)
	}

	if !strings.Contains(payload, "SubtipoMovimientoId:7") {
		t.Fatalf("payload sin subtipo esperado: %s", payload)
	}

	if strings.Contains(payload, "TipoMovimientoId:1") {
		t.Fatalf("payload no debe usar tipo movimiento cuando hay subtipo: %s", payload)
	}
}

func TestPayloadCuentasUsaTipoMovimientoCuandoNoHaySubtipo(t *testing.T) {
	t.Parallel()

	payload := payloadCuentas(53504, 1, 0)

	if !strings.Contains(payload, "TipoMovimientoId:1") {
		t.Fatalf("payload sin tipo movimiento esperado: %s", payload)
	}
}

func TestMensajeFaltaParametrizacionSubgrupoIncluyeCodigoYNombre(t *testing.T) {
	t.Parallel()

	msg := mensajeFaltaParametrizacionSubgrupo(models.DetalleSubgrupo{
		SubgrupoId: &models.Subgrupo{
			Codigo: "123",
			Nombre: "Equipos de prueba",
		},
	})

	if !strings.Contains(msg, "123") || !strings.Contains(msg, "Equipos de prueba") {
		t.Fatalf("mensaje sin contexto de subgrupo: %s", msg)
	}
}
