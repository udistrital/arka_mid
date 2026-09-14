package movimientosArka

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/udistrital/arka_mid/models"
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

func TestGetAllMovimientoControlaErrorHTTP(t *testing.T) {
	original := getAllMovimientoRequest
	getAllMovimientoRequest = func(string, interface{}) (*http.Response, error) {
		return nil, errors.New("timeout")
	}
	t.Cleanup(func() { getAllMovimientoRequest = original })

	movimientos, count, outputError := GetAllMovimiento("limit=-1")
	if outputError == nil {
		t.Fatal("se esperaba un error controlado")
	}
	if outputError["status"] != "502" {
		t.Fatalf("status inesperado: %v", outputError["status"])
	}
	if movimientos != nil || count != "" {
		t.Fatalf("no se esperaba respuesta parcial: movimientos=%v count=%q", movimientos, count)
	}
}

func TestGetAllMovimientoControlaRespuestaNula(t *testing.T) {
	original := getAllMovimientoRequest
	getAllMovimientoRequest = func(string, interface{}) (*http.Response, error) {
		return nil, nil
	}
	t.Cleanup(func() { getAllMovimientoRequest = original })

	_, _, outputError := GetAllMovimiento("limit=-1")
	if outputError == nil || outputError["status"] != "502" {
		t.Fatalf("se esperaba error 502 para respuesta nula: %v", outputError)
	}
}

func TestNormalizarConsultaCentroCostosPorID(t *testing.T) {
	t.Parallel()

	got := normalizarConsultaCentroCostos("query=Id:422")
	want := "query=Id__in:422"
	if got != want {
		t.Fatalf("consulta normalizada = %q; se esperaba %q", got, want)
	}
}

func TestNormalizarConsultaCentroCostosConservaFiltroPorCodigo(t *testing.T) {
	t.Parallel()

	query := "query=Codigo:A1302040205"
	if got := normalizarConsultaCentroCostos(query); got != query {
		t.Fatalf("consulta normalizada = %q; se esperaba %q", got, query)
	}
}

func TestCentroCostosAceptaCodigoString(t *testing.T) {
	t.Parallel()

	var centros []models.CentroCostos
	data := []byte(`[{"Id":422,"Codigo":"A1302040205","Nombre":"Ingenieria"}]`)
	if err := json.Unmarshal(data, &centros); err != nil {
		t.Fatalf("no se pudo decodificar Codigo como string: %v", err)
	}
	if len(centros) != 1 || centros[0].Codigo != "A1302040205" {
		t.Fatalf("centros de costos inesperados: %#v", centros)
	}
}
