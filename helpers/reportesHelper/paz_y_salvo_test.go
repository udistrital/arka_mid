package reportesHelper

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/udistrital/arka_mid/models"
)

func TestGenerarPazYSalvoSinElementos(t *testing.T) {
	mockConsultarElaboradorPazYSalvo(t)
	mockConsultarResponsableFirma(t)
	mockConsultarTerceroPorDocumento(t, models.DetalleTercero{
		Tercero: &models.Tercero{Id: 99, NombreCompleto: "Carmenza Moreno Roa"},
	})
	mockConsultarInventarioTercero(t, &models.InventarioTercero{
		Tercero: models.DetalleFuncionario{
			Tercero: []models.DetalleTercero{
				{
					Tercero: &models.Tercero{Id: 99, NombreCompleto: "Carmenza Moreno Roa"},
					Identificacion: &models.DatosIdentificacion{
						Numero: "52083089",
						TipoDocumentoId: &models.TipoDocumento{
							CodigoAbreviacion: "CC",
						},
					},
				},
			},
			Cargo: []*models.Parametro{
				{Nombre: "Profesional Universitario"},
			},
		},
		Elementos: []models.DetalleElementoPlaca{},
	})

	respuesta, err := GenerarPazYSalvo(&models.PazYSalvoRequest{
		Usuario:          "usuario@udistrital.edu.co",
		ElaboroTerceroId: 999,
		NumeroDocumento:  "52083089",
	})
	if err != nil {
		t.Fatalf("GenerarPazYSalvo retornó error: %v", err)
	}

	if respuesta == nil {
		t.Fatal("GenerarPazYSalvo retornó respuesta nil")
	}

	if !respuesta.PuedeGenerarPazYSalvo {
		t.Fatal("se esperaba paz y salvo generable")
	}

	if respuesta.TipoArchivo != pdfMimeType {
		t.Fatalf("tipo_archivo inesperado: %q", respuesta.TipoArchivo)
	}

	if respuesta.Tercero == nil || respuesta.Tercero.NombreCompleto != "Carmenza Moreno Roa" {
		t.Fatalf("tercero inesperado: %#v", respuesta.Tercero)
	}

	contenido, decodeErr := base64.StdEncoding.DecodeString(respuesta.ArchivoBase64)
	if decodeErr != nil {
		t.Fatalf("base64 inválido: %v", decodeErr)
	}

	if !strings.HasPrefix(string(contenido), "%PDF") {
		t.Fatalf("se esperaba un PDF, prefijo obtenido: %q", string(contenido[:4]))
	}
}

func TestGenerarPazYSalvoConElementos(t *testing.T) {
	mockConsultarElaboradorPazYSalvo(t)
	mockConsultarResponsableFirma(t)
	mockConsultarTerceroPorDocumento(t, models.DetalleTercero{
		Tercero: &models.Tercero{Id: 101, NombreCompleto: "Funcionario Con Elementos"},
	})
	mockConsultarInventarioTercero(t, &models.InventarioTercero{
		Tercero: models.DetalleFuncionario{
			Tercero: []models.DetalleTercero{
				{
					Tercero: &models.Tercero{Id: 101, NombreCompleto: "Funcionario Con Elementos"},
					Identificacion: &models.DatosIdentificacion{
						Numero: "123456",
						TipoDocumentoId: &models.TipoDocumento{
							CodigoAbreviacion: "CC",
						},
					},
				},
			},
		},
		Elementos: []models.DetalleElementoPlaca{
			{
				Placa:  "PLA-001",
				Nombre: "Portátil",
				Marca:  "Dell",
				Serie:  "SER-001",
				Valor:  3500000,
			},
		},
	})

	respuesta, err := GenerarPazYSalvo(&models.PazYSalvoRequest{
		Usuario:          "usuario@udistrital.edu.co",
		ElaboroTerceroId: 999,
		NumeroDocumento:  "123456",
	})
	if err != nil {
		t.Fatalf("GenerarPazYSalvo retornó error: %v", err)
	}

	if respuesta == nil {
		t.Fatal("GenerarPazYSalvo retornó respuesta nil")
	}

	if respuesta.PuedeGenerarPazYSalvo {
		t.Fatal("no se esperaba paz y salvo generable")
	}

	if len(respuesta.Elementos) != 1 {
		t.Fatalf("cantidad de elementos inesperada: %d", len(respuesta.Elementos))
	}

	if respuesta.Elementos[0].Placa != "PLA-001" {
		t.Fatalf("placa inesperada: %q", respuesta.Elementos[0].Placa)
	}
}

func TestGenerarPazYSalvoContinuaSinElaborador(t *testing.T) {
	mockConsultarElaboradorPazYSalvoConError(t)
	mockConsultarResponsableFirma(t)
	mockConsultarTerceroPorDocumento(t, models.DetalleTercero{
		Tercero: &models.Tercero{Id: 77, NombreCompleto: "Funcionario Sin Elaborador"},
	})
	mockConsultarInventarioTercero(t, &models.InventarioTercero{
		Tercero: models.DetalleFuncionario{
			Tercero: []models.DetalleTercero{
				{
					Tercero: &models.Tercero{Id: 77, NombreCompleto: "Funcionario Sin Elaborador"},
					Identificacion: &models.DatosIdentificacion{
						Numero: "770001",
						TipoDocumentoId: &models.TipoDocumento{
							CodigoAbreviacion: "CC",
						},
					},
				},
			},
		},
		Elementos: []models.DetalleElementoPlaca{},
	})

	respuesta, err := GenerarPazYSalvo(&models.PazYSalvoRequest{
		Usuario:          "usuario@udistrital.edu.co",
		ElaboroTerceroId: 999,
		NumeroDocumento:  "770001",
	})
	if err != nil {
		t.Fatalf("GenerarPazYSalvo no debe fallar si no se resuelve Elaboró: %v", err)
	}

	if respuesta == nil || strings.TrimSpace(respuesta.ArchivoBase64) == "" {
		t.Fatal("se esperaba PDF generado aun sin elaborador")
	}
}

func TestFechaLargaEsUsaDiaReal(t *testing.T) {
	fecha := time.Date(2026, time.July, 15, 10, 30, 0, 0, time.UTC)

	resultado := fechaLargaEs(fecha)

	if resultado != "15 dias de julio de 2026" {
		t.Fatalf("fecha larga inesperada: %q", resultado)
	}
}

func mockConsultarTerceroPorDocumento(t *testing.T, detalle models.DetalleTercero) {
	t.Helper()

	original := consultarTerceroPorDocumentoFn
	consultarTerceroPorDocumentoFn = func(numeroDocumento string) (models.DetalleTercero, map[string]interface{}) {
		return detalle, nil
	}
	t.Cleanup(func() {
		consultarTerceroPorDocumentoFn = original
	})
}

func mockConsultarInventarioTercero(t *testing.T, inventario *models.InventarioTercero) {
	t.Helper()

	original := consultarInventarioTerceroFn
	consultarInventarioTerceroFn = func(terceroId int) (*models.InventarioTercero, map[string]interface{}) {
		return inventario, nil
	}
	t.Cleanup(func() {
		consultarInventarioTerceroFn = original
	})
}

func mockConsultarElaboradorPazYSalvo(t *testing.T) {
	t.Helper()

	original := consultarElaboradorPazYSalvoFn
	consultarElaboradorPazYSalvoFn = func(string, int) (*pazYSalvoFirmante, map[string]interface{}) {
		return &pazYSalvoFirmante{
			Nombre: "Usuario Prueba",
			Cargo:  "Cargo Prueba",
		}, nil
	}
	t.Cleanup(func() {
		consultarElaboradorPazYSalvoFn = original
	})
}

func mockConsultarElaboradorPazYSalvoConError(t *testing.T) {
	t.Helper()

	original := consultarElaboradorPazYSalvoFn
	consultarElaboradorPazYSalvoFn = func(string, int) (*pazYSalvoFirmante, map[string]interface{}) {
		return nil, map[string]interface{}{"status": "500", "err": "sin usuario"}
	}
	t.Cleanup(func() {
		consultarElaboradorPazYSalvoFn = original
	})
}

func mockConsultarResponsableFirma(t *testing.T) {
	t.Helper()

	original := consultarResponsableFirmaFn
	consultarResponsableFirmaFn = func(models.DetalleTercero) (*pazYSalvoFirmante, map[string]interface{}) {
		return &pazYSalvoFirmante{
			Nombre:          "Firmante Prueba",
			Cargo:           "Cargo Firmante",
			TipoDocumento:   "CC",
			NumeroDocumento: "123456789",
		}, nil
	}
	t.Cleanup(func() {
		consultarResponsableFirmaFn = original
	})
}
