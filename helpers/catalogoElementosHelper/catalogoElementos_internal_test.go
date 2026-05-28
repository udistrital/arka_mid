package catalogoElementosHelper

import (
	"testing"

	"github.com/udistrital/arka_mid/models"
)

func TestConsultarCuentasSubgrupoRecientesFiltraYDeduplica(t *testing.T) {
	t.Parallel()

	input := []*models.CuentasSubgrupo{
		{
			Id:                  10,
			TipoMovimientoId:    1,
			SubtipoMovimientoId: 7,
			TipoBienId:          &models.TipoBien{Id: 10},
		},
		{
			Id:                  9,
			TipoMovimientoId:    1,
			SubtipoMovimientoId: 7,
			TipoBienId:          &models.TipoBien{Id: 10},
		},
		{
			Id:                  8,
			TipoMovimientoId:    0,
			SubtipoMovimientoId: 1,
			TipoBienId:          &models.TipoBien{Id: 10},
		},
		{
			Id:                  7,
			TipoMovimientoId:    1,
			SubtipoMovimientoId: 27,
			TipoBienId:          &models.TipoBien{Id: 10},
		},
		{
			Id:                  6,
			TipoMovimientoId:    22,
			SubtipoMovimientoId: 27,
			TipoBienId:          &models.TipoBien{Id: 10},
		},
	}

	got := seleccionarCuentasSubgrupoRecientes(input, 1)
	if len(got) != 3 {
		t.Fatalf("expected 3 latest cuentas for movimiento 1, got %d", len(got))
	}

	if got[0].Id != 10 || got[1].Id != 8 || got[2].Id != 7 {
		t.Fatalf("unexpected selected cuentas: %+v", got)
	}
}

func TestFormatoMovimientoByID(t *testing.T) {
	t.Parallel()

	formatos := map[int]models.FormatoTipoMovimiento{
		7: {Id: 7, CodigoAbreviacion: "SAL"},
	}

	cero := formatoMovimientoByID(0, formatos)
	if cero == nil || cero.Id != 0 {
		t.Fatalf("expected zero-value formato for id 0, got %+v", cero)
	}

	salida := formatoMovimientoByID(7, formatos)
	if salida == nil || salida.CodigoAbreviacion != "SAL" {
		t.Fatalf("expected SAL formato, got %+v", salida)
	}

	desconocido := formatoMovimientoByID(99, formatos)
	if desconocido == nil || desconocido.Id != 99 {
		t.Fatalf("expected placeholder formato for unknown id, got %+v", desconocido)
	}
}
