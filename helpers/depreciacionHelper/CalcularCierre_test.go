package depreciacionHelper

import (
	"testing"

	"github.com/udistrital/arka_mid/models"
)

func TestElementoValidoParaDepreciacion(t *testing.T) {
	t.Parallel()

	if elementoValidoParaDepreciacion(nil) {
		t.Fatal("expected nil elemento to be skipped")
	}

	if elementoValidoParaDepreciacion(&models.Elemento{Activo: false, TipoBienId: 10}) {
		t.Fatal("expected inactive elemento to be skipped")
	}

	if elementoValidoParaDepreciacion(&models.Elemento{Activo: true, TipoBienId: 0}) {
		t.Fatal("expected elemento without tipo bien to be skipped")
	}

	if !elementoValidoParaDepreciacion(&models.Elemento{Activo: true, TipoBienId: 10}) {
		t.Fatal("expected active elemento with tipo bien to be included")
	}
}
