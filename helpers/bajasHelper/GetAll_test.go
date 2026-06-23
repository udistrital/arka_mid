package bajasHelper

import "testing"

func TestEsTerceroNoEncontrado(t *testing.T) {
	errMap := map[string]interface{}{
		"err":    "http 404: {\"Message\":\"Not found resource\"}",
		"status": "502",
	}

	if !esTerceroNoEncontrado(errMap) {
		t.Fatal("se esperaba detectar 404 de terceros")
	}
}

func TestCargarNombreTerceroBajaIdInvalido(t *testing.T) {
	buffer := make(map[int]string)

	if err := cargarNombreTerceroBaja(0, buffer); err != nil {
		t.Fatalf("no se esperaba error para tercero 0: %v", err)
	}

	if got := buffer[0]; got != "" {
		t.Fatalf("se esperaba valor vacío para tercero 0, se obtuvo %q", got)
	}
}
