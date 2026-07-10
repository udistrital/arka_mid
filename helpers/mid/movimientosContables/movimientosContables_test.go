package movimientosContables

import (
	"errors"
	"testing"

	"github.com/udistrital/arka_mid/models"
)

func TestPostTrContableRetryOnHTTP404(t *testing.T) {
	originalSendJSON := sendJSONPostTrContable
	t.Cleanup(func() {
		sendJSONPostTrContable = originalSendJSON
	})

	calls := 0
	sendJSONPostTrContable = func(urlp string, trequest string, target interface{}, datajson interface{}) error {
		calls++
		if calls == 1 {
			return errors.New(`http 404: {"Data":null,"Message":null,"Status":"404","Success":false}`)
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
		return nil
	}

	tr := &models.TransaccionMovimientos{ConsecutivoId: 10764}
	res, err := PostTrContable(tr)
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
