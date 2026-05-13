package reportesHelper

import (
	"bytes"
	"encoding/base64"
	"testing"
	"time"

	"github.com/tealeg/xlsx"
	"github.com/udistrital/arka_mid/models"
)

func TestGenerarReporteElementos(t *testing.T) {
	mockConsultarEntradasReporteData(t, []*entradaReporteData{
		{
			Movimiento: &models.Movimiento{
				Id:            7995,
				Consecutivo:   stringPtr("ENT-7995"),
				FechaCreacion: time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
				FormatoTipoMovimientoId: &models.FormatoTipoMovimiento{
					Nombre: "Entrada por compra",
				},
				EstadoMovimientoId: &models.EstadoMovimiento{
					Nombre: "Entrada Aprobada",
				},
			},
			Formato: models.FormatoBaseEntrada{
				ActaRecibidoId: 555,
			},
			TransaccionContable: &models.InfoTransaccionContable{
				Concepto: "Entrada Almacén",
				Fecha:    time.Date(2026, 5, 9, 11, 0, 0, 0, time.UTC),
				Movimientos: []*models.DetalleMovimientoContable{
					{
						Cuenta: &models.DetalleCuenta{
							Id:     "cta-db-1",
							Codigo: "151001",
							Nombre: "Equipo de cómputo",
						},
						Debito: 2500,
					},
					{
						Cuenta: &models.DetalleCuenta{
							Id:     "cta-cr-1",
							Codigo: "240801",
							Nombre: "Bienes recibidos",
						},
						Credito: 2500,
					},
				},
			},
			Elementos: []*models.DetalleElemento{
				{
					Id:            101,
					Nombre:        "Elemento Uno",
					Cantidad:      2,
					Marca:         "Marca A",
					Serie:         "SERIE-001",
					UnidadMedida:  1,
					ValorUnitario: 1250.567,
					Subtotal:      2501.134,
					ValorTotal:    2501.134,
					ValorIva:      475.246,
					ValorFinal:    2976.38,
					Activo:        true,
					Placa:         "PL-001",
					FechaCreacion: time.Date(2026, 5, 9, 10, 15, 0, 0, time.UTC),
					ActaRecibidoId: &models.ActaRecibido{
						Id: 555,
					},
					SubgrupoCatalogoId: &models.DetalleSubgrupo{
						VidaUtil: 5,
						SubgrupoId: &models.Subgrupo{
							Id:     9,
							Codigo: "SG-09",
							Nombre: "Computadores",
						},
					},
					TipoBienId: &models.TipoBien{
						Id:     3,
						Nombre: "Devolutivo",
					},
				},
			},
			CuentasPorSubgrupo: map[int]models.CuentasSubgrupo{
				9: {
					CuentaDebitoId:  "cta-db-1",
					CuentaCreditoId: "cta-cr-1",
				},
			},
			Proveedor:          "900123456 - Proveedor Uno",
			FacturaConsecutivo: "FAC-2026-001",
			FacturaFecha:       time.Date(2026, 5, 8, 8, 30, 0, 0, time.UTC),
			SalidasPorElemento: map[int]*salidaReporteData{
				101: {
					Movimiento: &models.Movimiento{
						Id:            9001,
						Consecutivo:   stringPtr("SAL-9001"),
						FechaCreacion: time.Date(2026, 5, 10, 14, 0, 0, 0, time.UTC),
						EstadoMovimientoId: &models.EstadoMovimiento{
							Nombre: "Salida Aprobada",
						},
					},
					FuncionarioAsignado: "12345 - Funcionario Uno",
					TrasladosAsociados:  "TRS-1001",
					Sede:                "Sede Central",
					Dependencia:         "Almacén General",
					TransaccionContable: &models.InfoTransaccionContable{
						Concepto: "Salida Almacén",
						Fecha:    time.Date(2026, 5, 10, 15, 0, 0, 0, time.UTC),
						Movimientos: []*models.DetalleMovimientoContable{
							{
								Cuenta: &models.DetalleCuenta{
									Id:     "cta-db-sal-1",
									Codigo: "839090",
									Nombre: "Responsabilidades en proceso",
								},
								Debito: 2500,
							},
							{
								Cuenta: &models.DetalleCuenta{
									Id:     "cta-cr-sal-1",
									Codigo: "151001",
									Nombre: "Equipo de cómputo",
								},
								Credito: 2500,
							},
						},
					},
					CuentasPorSubgrupo: map[int]models.CuentasSubgrupo{
						9: {
							CuentaDebitoId:  "cta-db-sal-1",
							CuentaCreditoId: "cta-cr-sal-1",
						},
					},
				},
			},
		},
	})

	respuesta, err := GenerarReporteElementos(&models.ReporteFechasRequest{
		FechaInicial: "2026-05-01",
		FechaFinal:   "2026-05-31",
	})
	if err != nil {
		t.Fatalf("GenerarReporteElementos retornó error: %v", err)
	}

	if respuesta == nil {
		t.Fatal("GenerarReporteElementos retornó respuesta nil")
	}

	contenido, decodeErr := base64.StdEncoding.DecodeString(respuesta.ArchivoBase64)
	if decodeErr != nil {
		t.Fatalf("base64 inválido: %v", decodeErr)
	}

	archivo, openErr := xlsx.OpenBinary(contenido)
	if openErr != nil {
		t.Fatalf("no se pudo abrir el excel generado: %v", openErr)
	}

	if len(archivo.Sheets) != 1 {
		t.Fatalf("se esperaba una hoja, se obtuvieron %d", len(archivo.Sheets))
	}

	if len(archivo.Sheets[0].Rows) != 2 {
		t.Fatalf("se esperaban 2 filas, se obtuvieron %d", len(archivo.Sheets[0].Rows))
	}

	headers := archivo.Sheets[0].Rows[0]
	if headers.Cells[0].String() != "Vigencia" {
		t.Fatalf("encabezado inesperado en la primera columna: %q", headers.Cells[0].String())
	}
	headerIndex := buildHeaderIndex(headers)

	dataRow := archivo.Sheets[0].Rows[1]
	if dataRow.Cells[headerIndex["Nombre / Descripción"]].String() != "Elemento Uno" {
		t.Fatalf("elemento_nombre inesperado: %q", dataRow.Cells[headerIndex["Nombre / Descripción"]].String())
	}
	if dataRow.Cells[headerIndex["Cuenta débito entrada"]].String() != "151001 - Equipo de cómputo" {
		t.Fatalf("cuenta débito entrada inesperada: %q", dataRow.Cells[headerIndex["Cuenta débito entrada"]].String())
	}
	if dataRow.Cells[headerIndex["Cuenta crédito entrada"]].String() != "240801 - Bienes recibidos" {
		t.Fatalf("cuenta crédito entrada inesperada: %q", dataRow.Cells[headerIndex["Cuenta crédito entrada"]].String())
	}
	if dataRow.Cells[headerIndex["Proveedor"]].String() != "900123456 - Proveedor Uno" {
		t.Fatalf("proveedor inesperado: %q", dataRow.Cells[headerIndex["Proveedor"]].String())
	}
	if dataRow.Cells[headerIndex["Consecutivo Factura"]].String() != "FAC-2026-001" {
		t.Fatalf("factura inesperada: %q", dataRow.Cells[headerIndex["Consecutivo Factura"]].String())
	}
	if dataRow.Cells[headerIndex["Tipo de entrada"]].String() != "Entrada por compra" {
		t.Fatalf("tipo de entrada inesperado: %q", dataRow.Cells[headerIndex["Tipo de entrada"]].String())
	}
	if dataRow.Cells[headerIndex["Vida útil (años)"]].Value != "5" {
		t.Fatalf("vida útil inesperada: %q", dataRow.Cells[headerIndex["Vida útil (años)"]].Value)
	}
	if dataRow.Cells[headerIndex["Consecutivo salida"]].String() != "SAL-9001" {
		t.Fatalf("consecutivo salida inesperado: %q", dataRow.Cells[headerIndex["Consecutivo salida"]].String())
	}
	if dataRow.Cells[headerIndex["Funcionario asignado"]].String() != "12345 - Funcionario Uno" {
		t.Fatalf("funcionario asignado inesperado: %q", dataRow.Cells[headerIndex["Funcionario asignado"]].String())
	}
	if dataRow.Cells[headerIndex["Sede"]].String() != "Sede Central" {
		t.Fatalf("sede inesperada: %q", dataRow.Cells[headerIndex["Sede"]].String())
	}
	if dataRow.Cells[headerIndex["Dependencia"]].String() != "Almacén General" {
		t.Fatalf("dependencia inesperada: %q", dataRow.Cells[headerIndex["Dependencia"]].String())
	}
	if dataRow.Cells[headerIndex["Cuenta débito salida"]].String() != "839090 - Responsabilidades en proceso" {
		t.Fatalf("cuenta débito salida inesperada: %q", dataRow.Cells[headerIndex["Cuenta débito salida"]].String())
	}
	if dataRow.Cells[headerIndex["Cuenta crédito salida"]].String() != "151001 - Equipo de cómputo" {
		t.Fatalf("cuenta crédito salida inesperada: %q", dataRow.Cells[headerIndex["Cuenta crédito salida"]].String())
	}
	if dataRow.Cells[headerIndex["Valor unitario"]].Type() != xlsx.CellTypeNumeric {
		t.Fatalf("valor unitario debe ser numérico, se obtuvo tipo %v", dataRow.Cells[headerIndex["Valor unitario"]].Type())
	}
	if dataRow.Cells[headerIndex["Valor unitario"]].GetNumberFormat() != decimalNumFmt {
		t.Fatalf("formato numérico inesperado para valor unitario: %q", dataRow.Cells[headerIndex["Valor unitario"]].GetNumberFormat())
	}
	if dataRow.Cells[headerIndex["Valor unitario"]].Value != "1250.57" {
		t.Fatalf("valor interno inesperado para valor unitario: %q", dataRow.Cells[headerIndex["Valor unitario"]].Value)
	}
	if dataRow.Cells[headerIndex["Subtotal"]].Value != "2501.13" {
		t.Fatalf("valor interno inesperado para subtotal: %q", dataRow.Cells[headerIndex["Subtotal"]].Value)
	}
	if dataRow.Cells[headerIndex["IVA"]].Value != "475.25" {
		t.Fatalf("valor interno inesperado para IVA: %q", dataRow.Cells[headerIndex["IVA"]].Value)
	}
}

func TestGenerarReporteElementosFechaFinalMenor(t *testing.T) {
	respuesta, err := GenerarReporteElementos(&models.ReporteFechasRequest{
		FechaInicial: "2026-02-01",
		FechaFinal:   "2026-01-31",
	})
	if err == nil {
		t.Fatal("se esperaba error cuando fecha_final es menor a fecha_inicial")
	}

	if respuesta != nil {
		t.Fatal("no se esperaba respuesta cuando el rango es inválido")
	}
}

func TestExcelGeneradoEsBinarioValido(t *testing.T) {
	mockConsultarEntradasReporteData(t, []*entradaReporteData{})

	respuesta, err := GenerarReporteElementos(&models.ReporteFechasRequest{
		FechaInicial: "2026-03-01",
		FechaFinal:   "2026-03-15",
	})
	if err != nil {
		t.Fatalf("GenerarReporteElementos retornó error: %v", err)
	}

	contenido, decodeErr := base64.StdEncoding.DecodeString(respuesta.ArchivoBase64)
	if decodeErr != nil {
		t.Fatalf("base64 inválido: %v", decodeErr)
	}

	if len(bytes.TrimSpace(contenido)) == 0 {
		t.Fatal("el excel generado no debe ser vacío")
	}
}

func mockConsultarEntradasReporteData(t *testing.T, entradas []*entradaReporteData) {
	t.Helper()

	original := consultarEntradasReporteData
	consultarEntradasReporteData = func(fechaInicial, fechaFinal time.Time) ([]*entradaReporteData, map[string]interface{}) {
		return entradas, nil
	}

	t.Cleanup(func() {
		consultarEntradasReporteData = original
	})
}

func stringPtr(value string) *string {
	return &value
}

func buildHeaderIndex(headerRow *xlsx.Row) map[string]int {
	index := make(map[string]int, len(headerRow.Cells))
	for i, cell := range headerRow.Cells {
		index[cell.String()] = i
	}
	return index
}
