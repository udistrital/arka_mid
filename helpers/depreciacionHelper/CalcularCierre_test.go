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

func TestConsultarElementosParaCierre(t *testing.T) {
	originalGetter := getAllElementoDepreciacionCierre
	getAllElementoDepreciacionCierre = func(query string, fields string, sortby string, order string, offset string, limit string) (elementos []*models.Elemento, outputError map[string]interface{}) {
		switch query {
		case "Id:1":
			return []*models.Elemento{{Id: 1, Activo: true, TipoBienId: 10, SubgrupoCatalogoId: 100}}, nil
		case "Id:2":
			return []*models.Elemento{{Id: 2, Activo: false, TipoBienId: 10, SubgrupoCatalogoId: 100}}, nil
		case "Id:3":
			return []*models.Elemento{{Id: 3, Activo: true, TipoBienId: 0, SubgrupoCatalogoId: 100}}, nil
		default:
			return nil, nil
		}
	}
	defer func() { getAllElementoDepreciacionCierre = originalGetter }()

	detalles, errMsg, outputError := consultarElementosParaCierre([]models.DepreciacionElemento{
		{ElementoActaId: 1, DeltaValor: 15},
		{ElementoActaId: 2, DeltaValor: 20},
		{ElementoActaId: 3, DeltaValor: 25},
		{ElementoActaId: 4, DeltaValor: 30},
		{ElementoActaId: 5, DeltaValor: 0},
	})

	if outputError != nil {
		t.Fatalf("expected no output error, got %#v", outputError)
	}

	if errMsg == "" {
		t.Fatal("expected missing element to produce business error message")
	}

	if len(detalles) != 1 {
		t.Fatalf("expected only one valid elemento to be collected, got %d", len(detalles))
	}

	if detalles[0].elemento == nil || detalles[0].elemento.Id != 1 || detalles[0].deltaValor != 15 {
		t.Fatalf("unexpected collected detalle %+v", detalles[0])
	}
}
