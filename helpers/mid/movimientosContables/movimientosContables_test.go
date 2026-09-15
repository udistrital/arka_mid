package movimientosContables

import (
	"context"
	"errors"
	"testing"

	"github.com/udistrital/arka_mid/models"
)

func TestPostTrContableRetryOnHTTP404(t *testing.T) {
	originalPostWithContext := postWithContext
	t.Cleanup(func() {
		postWithContext = originalPostWithContext
	})

	calls := 0
	postWithContext = func(ctx context.Context, urlp string, body, target any) (int, error) {
		calls++
		if calls == 1 {
			return 404, errors.New("unexpected status code: 404")
		}

		resp, ok := target.(*map[string]interface{})
		if !ok {
			t.Fatalf("unexpected target type: %T", target)
		}

		*resp = map[string]interface{}{
			"Success": true,
			"Status":  "201",
			"Data":    "OK",
		}
		return 201, nil
	}

	tr := &models.TransaccionMovimientos{ConsecutivoId: 10764}
	res, err := PostTrContable(context.Background(), tr)
	if err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}
	if res == nil {
		t.Fatal("expected transaction result")
	}
	if calls != 2 {
		t.Fatalf("expected 2 attempts, got %d", calls)
	}
}

func TestNormalizeMovimientosContablesBasePathPromotesHTTPS(t *testing.T) {
	got := normalizeMovimientosContablesBasePath("http://pruebasapi.intranetoas.udistrital.edu.co/movimientos_contables_mid/v1/")
	want := "https://pruebasapi.intranetoas.udistrital.edu.co/movimientos_contables_mid/v1/"
	if got != want {
		t.Fatalf("unexpected base path: got %q want %q", got, want)
	}
}
