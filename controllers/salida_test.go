package controllers

import (
	"testing"

	"github.com/udistrital/arka_mid/models"
)

func TestDecodeSalidaGeneralRequestDirecto(t *testing.T) {
	t.Parallel()

	body := []byte(`{"Salidas":[{"Salida":{"MovimientoPadreId":{"Id":8079}}}]}`)
	var target models.SalidaGeneral

	if err := decodeSalidaGeneralRequest(body, &target); err != nil {
		t.Fatalf("expected direct payload to decode, got %v", err)
	}

	if len(target.Salidas) != 1 {
		t.Fatalf("expected 1 salida, got %d", len(target.Salidas))
	}

	if target.Salidas[0].Salida == nil || target.Salidas[0].Salida.MovimientoPadreId == nil || target.Salidas[0].Salida.MovimientoPadreId.Id != 8079 {
		t.Fatal("expected parent movement id 8079 to be preserved")
	}
}

func TestDecodeSalidaGeneralRequestEnvuelto(t *testing.T) {
	t.Parallel()

	body := []byte(`{"trSalida":{"Salidas":[{"Salida":{"MovimientoPadreId":{"Id":8079}}}]}}`)
	var target models.SalidaGeneral

	if err := decodeSalidaGeneralRequest(body, &target); err != nil {
		t.Fatalf("expected wrapped payload to decode, got %v", err)
	}

	if len(target.Salidas) != 1 {
		t.Fatalf("expected 1 salida, got %d", len(target.Salidas))
	}
}

func TestDecodeSalidaGeneralRequestVacio(t *testing.T) {
	t.Parallel()

	var target models.SalidaGeneral
	if err := decodeSalidaGeneralRequest([]byte(`{}`), &target); err == nil {
		t.Fatal("expected empty payload to be rejected")
	}
}
