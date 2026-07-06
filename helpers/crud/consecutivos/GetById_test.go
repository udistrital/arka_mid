package consecutivos

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/udistrital/arka_mid/models"
)

func TestGetByIdDecodesWrappedDataResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/consecutivo/10764" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"Data": {
				"Id": 10764,
				"ContextoId": 8484,
				"Year": 1997,
				"Consecutivo": 1,
				"Descripcion": "Entradas Arka",
				"Activo": true
			},
			"Message": "Request successful",
			"Status": "200",
			"Success": true
		}`))
	}))
	defer server.Close()

	original := ConsecutivosCRUD
	ConsecutivosCRUD = server.URL + "/"
	defer func() {
		ConsecutivosCRUD = original
	}()

	var consecutivo models.Consecutivo
	if err := GetById(10764, &consecutivo); err != nil {
		t.Fatalf("GetById() error = %#v", err)
	}

	if consecutivo.Id != 10764 {
		t.Fatalf("expected Id 10764, got %d", consecutivo.Id)
	}
	if consecutivo.ContextoId != 8484 {
		t.Fatalf("expected ContextoId 8484, got %d", consecutivo.ContextoId)
	}
	if consecutivo.Year != 1997 {
		t.Fatalf("expected Year 1997, got %d", consecutivo.Year)
	}
	if consecutivo.Consecutivo != 1 {
		t.Fatalf("expected Consecutivo 1, got %d", consecutivo.Consecutivo)
	}
}
