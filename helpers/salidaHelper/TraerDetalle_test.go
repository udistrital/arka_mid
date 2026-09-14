package salidaHelper

import (
	"testing"

	"github.com/udistrital/arka_mid/models"
)

func TestTraerDetalleResuelveUbicacionComoCentroCostosPorID(t *testing.T) {
	original := consultarCentroCostosSalida
	consultarCentroCostosSalida = func(payload string) ([]models.CentroCostos, map[string]interface{}) {
		if payload != "query=Id:57" {
			t.Fatalf("consulta inesperada: %q", payload)
		}
		return []models.CentroCostos{{
			Id:          57,
			Sede:        "Aduanilla de Paiba",
			Dependencia: "Almacén e Inventarios",
			Codigo:      "A130101",
			Nombre:      "Centro de costos",
		}}, nil
	}
	t.Cleanup(func() { consultarCentroCostosSalida = original })

	detalle, outputError := traerDetalle(
		&models.Movimiento{Id: 13183},
		models.FormatoSalidaCostos{FormatoSalida: models.FormatoSalida{Ubicacion: 57}},
		nil,
		nil,
	)
	if outputError != nil {
		t.Fatalf("traerDetalle retornó error: %v", outputError)
	}

	sede, ok := detalle["Sede"].(models.EspacioFisico)
	if !ok || sede.Nombre != "Aduanilla de Paiba" {
		t.Fatalf("sede inesperada: %#v", detalle["Sede"])
	}
	dependencia, ok := detalle["Dependencia"].(*models.Dependencia)
	if !ok || dependencia == nil || dependencia.Nombre != "Almacén e Inventarios" {
		t.Fatalf("dependencia inesperada: %#v", detalle["Dependencia"])
	}
	ubicacion, ok := detalle["Ubicacion"].(models.CentroCostos)
	if !ok || ubicacion.Id != 57 || ubicacion.Codigo != "A130101" {
		t.Fatalf("centro de costos inesperado: %#v", detalle["Ubicacion"])
	}
}

func TestTraerDetalleConservaConsultaPorCodigoCentroCostos(t *testing.T) {
	original := consultarCentroCostosSalida
	consultarCentroCostosSalida = func(payload string) ([]models.CentroCostos, map[string]interface{}) {
		if payload != "query=Codigo:A130101" {
			t.Fatalf("consulta inesperada: %q", payload)
		}
		return []models.CentroCostos{{Id: 57, Codigo: "A130101", Nombre: "Centro de costos"}}, nil
	}
	t.Cleanup(func() { consultarCentroCostosSalida = original })

	detalle, outputError := traerDetalle(
		&models.Movimiento{Id: 13183},
		models.FormatoSalidaCostos{CentroCostos: "A130101"},
		nil,
		nil,
	)
	if outputError != nil {
		t.Fatalf("traerDetalle retornó error: %v", outputError)
	}

	dependencia, ok := detalle["Dependencia"].(*models.Dependencia)
	if !ok || dependencia == nil || dependencia.Nombre != "Centro de costos" {
		t.Fatalf("dependencia inesperada: %#v", detalle["Dependencia"])
	}
}

func TestTraerDetalleRetornaMarcadorCuandoCentroCostosNoExiste(t *testing.T) {
	original := consultarCentroCostosSalida
	consultarCentroCostosSalida = func(payload string) ([]models.CentroCostos, map[string]interface{}) {
		if payload != "query=Id:57" {
			t.Fatalf("consulta inesperada: %q", payload)
		}
		return []models.CentroCostos{}, nil
	}
	t.Cleanup(func() { consultarCentroCostosSalida = original })

	detalle, outputError := traerDetalle(
		&models.Movimiento{Id: 13183},
		models.FormatoSalidaCostos{FormatoSalida: models.FormatoSalida{Ubicacion: 57}},
		nil,
		nil,
	)
	if outputError != nil {
		t.Fatalf("traerDetalle retornó error: %v", outputError)
	}

	ubicacion, ok := detalle["Ubicacion"].(models.CentroCostos)
	if !ok {
		t.Fatalf("ubicación inesperada: %#v", detalle["Ubicacion"])
	}
	if ubicacion.Codigo != "0" || ubicacion.Nombre != mensajeCentroCostosNoEncontrado {
		t.Fatalf("marcador de centro de costos inesperado: %#v", ubicacion)
	}

	dependencia, ok := detalle["Dependencia"].(*models.Dependencia)
	if !ok || dependencia == nil || dependencia.Nombre != mensajeCentroCostosNoEncontrado {
		t.Fatalf("dependencia inesperada: %#v", detalle["Dependencia"])
	}
}
