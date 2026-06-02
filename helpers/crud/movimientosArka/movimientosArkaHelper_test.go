package movimientosArka

import (
	"strings"
	"testing"
)

func TestBuildCorteDepreciacionURL(t *testing.T) {
	t.Parallel()

	url := buildCorteDepreciacionURL("2026-04-30")

	if strings.Contains(url, "cierre/?") {
		t.Fatalf("la url no debe contener barra antes del query: %s", url)
	}

	if !strings.Contains(url, "cierre?fechaCorte=2026-04-30") {
		t.Fatalf("url inesperada: %s", url)
	}
}
