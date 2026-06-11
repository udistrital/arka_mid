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
	mockConsultarEntradasAnuladasReporteData(t, []*entradaReporteData{})
	mockConsultarSalidasAnuladasReporteData(t, []*salidaReporteData{})

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
	mockConsultarEntradasAnuladasReporteData(t, []*entradaReporteData{})
	mockConsultarSalidasAnuladasReporteData(t, []*salidaReporteData{})

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

func TestGenerarReporteIncluyeMovimientosAnuladosSinElementos(t *testing.T) {
	mockConsultarEntradasReporteData(t, []*entradaReporteData{})
	mockConsultarEntradasAnuladasReporteData(t, []*entradaReporteData{
		{
			Movimiento: &models.Movimiento{
				Id:            8101,
				Activo:        false,
				Consecutivo:   stringPtr("ENT-ANU-01"),
				FechaCreacion: time.Date(2026, 5, 11, 9, 0, 0, 0, time.UTC),
				FormatoTipoMovimientoId: &models.FormatoTipoMovimiento{
					Nombre: "Entrada por compra",
				},
				EstadoMovimientoId: &models.EstadoMovimiento{
					Nombre: estadoEntradaAnuladaReporte,
				},
			},
			Formato: models.FormatoBaseEntrada{
				ActaRecibidoId: 777,
			},
			Elementos:          []*models.DetalleElemento{},
			CuentasPorSubgrupo: map[int]models.CuentasSubgrupo{},
			SalidasPorElemento: map[int]*salidaReporteData{},
		},
	})
	mockConsultarSalidasAnuladasReporteData(t, []*salidaReporteData{
		{
			Movimiento: &models.Movimiento{
				Id:            9101,
				Activo:        false,
				Consecutivo:   stringPtr("SAL-ANU-01"),
				FechaCreacion: time.Date(2026, 5, 12, 14, 0, 0, 0, time.UTC),
				FechaCorte:    timePtr(time.Date(2026, 5, 12, 15, 0, 0, 0, time.UTC)),
				EstadoMovimientoId: &models.EstadoMovimiento{
					Nombre: estadoSalidaAnuladaReporte,
				},
				MovimientoPadreId: &models.Movimiento{
					Id:            8102,
					Consecutivo:   stringPtr("ENT-BASE-01"),
					FechaCreacion: time.Date(2026, 5, 10, 8, 0, 0, 0, time.UTC),
					FormatoTipoMovimientoId: &models.FormatoTipoMovimiento{
						Nombre: "Entrada por compra",
					},
					EstadoMovimientoId: &models.EstadoMovimiento{
						Nombre: "Entrada Aprobada",
					},
					Detalle: `{"acta_recibido_id":888}`,
				},
			},
			FuncionarioAsignado: "12345 - Funcionario Uno",
			Sede:                "Sede Central",
			Dependencia:         "Almacén General",
		},
	})

	respuesta, err := GenerarReporteElementos(&models.ReporteFechasRequest{
		FechaInicial: "2026-05-01",
		FechaFinal:   "2026-05-31",
	})
	if err != nil {
		t.Fatalf("GenerarReporteElementos retornó error: %v", err)
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

	if len(archivo.Sheets[0].Rows) != 3 {
		t.Fatalf("se esperaban 3 filas, se obtuvieron %d", len(archivo.Sheets[0].Rows))
	}

	headerIndex := buildHeaderIndex(archivo.Sheets[0].Rows[0])
	entradaAnuladaRow := archivo.Sheets[0].Rows[1]
	if entradaAnuladaRow.Cells[headerIndex["Consecutivo Entrada"]].String() != "ENT-ANU-01" {
		t.Fatalf("consecutivo entrada anulada inesperado: %q", entradaAnuladaRow.Cells[headerIndex["Consecutivo Entrada"]].String())
	}
	if entradaAnuladaRow.Cells[headerIndex["entrada_estado"]].String() != estadoEntradaAnuladaReporte {
		t.Fatalf("estado entrada anulada inesperado: %q", entradaAnuladaRow.Cells[headerIndex["entrada_estado"]].String())
	}

	salidaAnuladaRow := archivo.Sheets[0].Rows[2]
	if salidaAnuladaRow.Cells[headerIndex["Consecutivo salida"]].String() != "SAL-ANU-01" {
		t.Fatalf("consecutivo salida anulada inesperado: %q", salidaAnuladaRow.Cells[headerIndex["Consecutivo salida"]].String())
	}
	if salidaAnuladaRow.Cells[headerIndex["salida_estado"]].String() != estadoSalidaAnuladaReporte {
		t.Fatalf("estado salida anulada inesperado: %q", salidaAnuladaRow.Cells[headerIndex["salida_estado"]].String())
	}
	if salidaAnuladaRow.Cells[headerIndex["Consecutivo Entrada"]].String() != "ENT-BASE-01" {
		t.Fatalf("consecutivo entrada padre inesperado: %q", salidaAnuladaRow.Cells[headerIndex["Consecutivo Entrada"]].String())
	}
}

func TestGetDetalleCuentasEntradaPorConsecutivo(t *testing.T) {
	mockConsultarMovimientoPorConsecutivo(t, &models.Movimiento{
		Id:            7995,
		Consecutivo:   stringPtr("ENT-7995"),
		Detalle:       `{"acta_recibido_id":555}`,
		FechaCreacion: time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
		FormatoTipoMovimientoId: &models.FormatoTipoMovimiento{
			Id:     55,
			Nombre: "Entrada por compra",
		},
		EstadoMovimientoId: &models.EstadoMovimiento{
			Nombre: "Entrada Aprobada",
		},
	})

	mockConsultarElementosActa(t, []*models.DetalleElemento{
		{
			Id:         101,
			Nombre:     "Elemento Uno",
			ValorFinal: 2976.38,
			SubgrupoCatalogoId: &models.DetalleSubgrupo{
				SubgrupoId: &models.Subgrupo{Id: 9},
			},
		},
	})
	mockConsultarMetadataEntrada(t, "900123456 - Proveedor Uno", "FAC-2026-001", time.Date(2026, 5, 8, 8, 30, 0, 0, time.UTC))
	mockResolverSalidasPorElemento(t, map[int]*salidaReporteData{
		101: {
			FuncionarioAsignado: "12345 - Funcionario Uno",
		},
	})
	mockGetCuentasByMovimientoAndSubgrupos(t, map[int]models.CuentasSubgrupo{
		9: {
			CuentaDebitoId:  "cta-db-1",
			CuentaCreditoId: "cta-cr-1",
		},
	})
	mockConsultarTransaccionContable(t, &models.InfoTransaccionContable{
		Movimientos: []*models.DetalleMovimientoContable{
			{Cuenta: &models.DetalleCuenta{Id: "cta-db-1", Codigo: "151001", Nombre: "Equipo de cómputo"}, Debito: 2976.38},
			{Cuenta: &models.DetalleCuenta{Id: "cta-cr-1", Codigo: "240801", Nombre: "Bienes recibidos"}, Credito: 2976.38},
		},
	})

	respuesta, err := GetDetalleCuentasEntradaPorConsecutivo("ENT-7995")
	if err != nil {
		t.Fatalf("GetDetalleCuentasEntradaPorConsecutivo retornó error: %v", err)
	}
	if len(respuesta) != 1 {
		t.Fatalf("se esperaba una fila, se obtuvieron %d", len(respuesta))
	}
	if respuesta[0].ElementoNombre != "Elemento Uno" {
		t.Fatalf("ElementoNombre inesperado: %q", respuesta[0].ElementoNombre)
	}
	if respuesta[0].ElementoValorFinal != 2976.38 {
		t.Fatalf("ElementoValorFinal inesperado: %v", respuesta[0].ElementoValorFinal)
	}
	if respuesta[0].SalidaFuncionarioAsignado != "12345 - Funcionario Uno" {
		t.Fatalf("SalidaFuncionarioAsignado inesperado: %q", respuesta[0].SalidaFuncionarioAsignado)
	}
	if respuesta[0].CuentaDebitoEntrada != "151001 - Equipo de cómputo" {
		t.Fatalf("CuentaDebitoEntrada inesperada: %q", respuesta[0].CuentaDebitoEntrada)
	}
	if respuesta[0].CuentaCreditoEntrada != "240801 - Bienes recibidos" {
		t.Fatalf("CuentaCreditoEntrada inesperada: %q", respuesta[0].CuentaCreditoEntrada)
	}
}

func TestGetDetalleCuentasSalidaPorConsecutivo(t *testing.T) {
	mockConsultarMovimientoPorConsecutivo(t, &models.Movimiento{
		Id:          9001,
		Consecutivo: stringPtr("SAL-9001"),
	})
	mockConsultarTrSalida(t, &models.TrSalida{
		Salida: &models.Movimiento{
			Id:          9001,
			Consecutivo: stringPtr("SAL-9001"),
			Detalle:     `{"funcionario":12345}`,
			FormatoTipoMovimientoId: &models.FormatoTipoMovimiento{
				Id: 77,
			},
			EstadoMovimientoId: &models.EstadoMovimiento{
				Nombre: "Salida Aprobada",
			},
		},
		Elementos: []*models.ElementosMovimiento{
			{
				Id:             1,
				ElementoActaId: intPtr(101),
			},
		},
	})
	mockConsultarElementosActa(t, []*models.DetalleElemento{
		{
			Id:         101,
			Nombre:     "Elemento Uno",
			ValorFinal: 2976.38,
			ValorTotal: 2976.38,
			SubgrupoCatalogoId: &models.DetalleSubgrupo{
				SubgrupoId: &models.Subgrupo{Id: 9},
			},
		},
	})
	mockGetCuentasByMovimientoAndSubgrupos(t, map[int]models.CuentasSubgrupo{
		9: {
			CuentaDebitoId:  "cta-db-sal-1",
			CuentaCreditoId: "cta-cr-sal-1",
		},
	})
	mockGetNombreTerceroByID(t, &models.IdentificacionTercero{
		Numero:         "12345",
		NombreCompleto: "Funcionario Uno",
	})
	mockConsultarTransaccionContable(t, &models.InfoTransaccionContable{
		Movimientos: []*models.DetalleMovimientoContable{
			{Cuenta: &models.DetalleCuenta{Id: "cta-db-sal-1", Codigo: "839090", Nombre: "Responsabilidades en proceso"}, Debito: 2976.38},
			{Cuenta: &models.DetalleCuenta{Id: "cta-cr-sal-1", Codigo: "151001", Nombre: "Equipo de cómputo"}, Credito: 2976.38},
		},
	})

	respuesta, err := GetDetalleCuentasSalidaPorConsecutivo("SAL-9001")
	if err != nil {
		t.Fatalf("GetDetalleCuentasSalidaPorConsecutivo retornó error: %v", err)
	}
	if len(respuesta) != 1 {
		t.Fatalf("se esperaba una fila, se obtuvieron %d", len(respuesta))
	}
	if respuesta[0].ElementoNombre != "Elemento Uno" {
		t.Fatalf("ElementoNombre inesperado: %q", respuesta[0].ElementoNombre)
	}
	if respuesta[0].ElementoValorFinal != 2976.38 {
		t.Fatalf("ElementoValorFinal inesperado: %v", respuesta[0].ElementoValorFinal)
	}
	if respuesta[0].SalidaFuncionarioAsignado != "12345 - Funcionario Uno" {
		t.Fatalf("SalidaFuncionarioAsignado inesperado: %q", respuesta[0].SalidaFuncionarioAsignado)
	}
	if respuesta[0].CuentaDebitoSalida != "839090 - Responsabilidades en proceso" {
		t.Fatalf("CuentaDebitoSalida inesperada: %q", respuesta[0].CuentaDebitoSalida)
	}
	if respuesta[0].CuentaCreditoSalida != "151001 - Equipo de cómputo" {
		t.Fatalf("CuentaCreditoSalida inesperada: %q", respuesta[0].CuentaCreditoSalida)
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

func mockConsultarEntradasAnuladasReporteData(t *testing.T, entradas []*entradaReporteData) {
	t.Helper()

	original := consultarEntradasAnuladasReporteData
	consultarEntradasAnuladasReporteData = func(fechaInicial, fechaFinal time.Time) ([]*entradaReporteData, map[string]interface{}) {
		return entradas, nil
	}

	t.Cleanup(func() {
		consultarEntradasAnuladasReporteData = original
	})
}

func mockConsultarSalidasAnuladasReporteData(t *testing.T, salidas []*salidaReporteData) {
	t.Helper()

	original := consultarSalidasAnuladasReporteData
	consultarSalidasAnuladasReporteData = func(fechaInicial, fechaFinal time.Time) ([]*salidaReporteData, map[string]interface{}) {
		return salidas, nil
	}

	t.Cleanup(func() {
		consultarSalidasAnuladasReporteData = original
	})
}

func mockConsultarMovimientoPorConsecutivo(t *testing.T, movimiento *models.Movimiento) {
	t.Helper()

	original := consultarMovimientoPorConsec
	consultarMovimientoPorConsec = func(consecutivo string) (*models.Movimiento, map[string]interface{}) {
		return movimiento, nil
	}

	t.Cleanup(func() {
		consultarMovimientoPorConsec = original
	})
}

func mockConsultarTrSalida(t *testing.T, trSalida *models.TrSalida) {
	t.Helper()

	original := consultarTrSalida
	consultarTrSalida = func(id int) (*models.TrSalida, map[string]interface{}) {
		return trSalida, nil
	}

	t.Cleanup(func() {
		consultarTrSalida = original
	})
}

func mockConsultarElementosActa(t *testing.T, elementos []*models.DetalleElemento) {
	t.Helper()

	original := consultarElementosActa
	consultarElementosActa = func(actaID int, ids []int) ([]*models.DetalleElemento, map[string]interface{}) {
		return elementos, nil
	}

	t.Cleanup(func() {
		consultarElementosActa = original
	})
}

func mockConsultarMetadataEntrada(t *testing.T, proveedor, facturaConsecutivo string, facturaFecha time.Time) {
	t.Helper()

	original := consultarMetadataEntradaFn
	consultarMetadataEntradaFn = func(formato models.FormatoBaseEntrada) (string, string, time.Time, map[string]interface{}) {
		return proveedor, facturaConsecutivo, facturaFecha, nil
	}

	t.Cleanup(func() {
		consultarMetadataEntradaFn = original
	})
}

func mockResolverSalidasPorElemento(t *testing.T, salidas map[int]*salidaReporteData) {
	t.Helper()

	original := resolverSalidasPorElementoFn
	resolverSalidasPorElementoFn = func(elementos []*models.DetalleElemento) (map[int]*salidaReporteData, map[string]interface{}) {
		return salidas, nil
	}

	t.Cleanup(func() {
		resolverSalidasPorElementoFn = original
	})
}

func mockGetCuentasByMovimientoAndSubgrupos(t *testing.T, cuentas map[int]models.CuentasSubgrupo) {
	t.Helper()

	original := getCuentasByMovimientoAndSubgrupos
	getCuentasByMovimientoAndSubgrupos = func(movimientoID int, subgrupos []int, cuentasPorSubgrupo map[int]models.CuentasSubgrupo) map[string]interface{} {
		for id, cuenta := range cuentas {
			cuentasPorSubgrupo[id] = cuenta
		}
		return nil
	}

	t.Cleanup(func() {
		getCuentasByMovimientoAndSubgrupos = original
	})
}

func mockGetNombreTerceroByID(t *testing.T, tercero *models.IdentificacionTercero) {
	t.Helper()

	original := getNombreTerceroByID
	getNombreTerceroByID = func(terceroID int) (*models.IdentificacionTercero, map[string]interface{}) {
		return tercero, nil
	}

	t.Cleanup(func() {
		getNombreTerceroByID = original
	})
}

func mockConsultarTransaccionContable(t *testing.T, transaccion *models.InfoTransaccionContable) {
	t.Helper()

	original := consultarTransaccionContableMovimientoFn
	consultarTransaccionContableMovimientoFn = func(movimiento *models.Movimiento, estadosPermitidos ...string) *models.InfoTransaccionContable {
		return transaccion
	}

	t.Cleanup(func() {
		consultarTransaccionContableMovimientoFn = original
	})
}

func stringPtr(value string) *string {
	return &value
}

func intPtr(value int) *int {
	return &value
}

func timePtr(value time.Time) *time.Time {
	return &value
}

func buildHeaderIndex(headerRow *xlsx.Row) map[string]int {
	index := make(map[string]int, len(headerRow.Cells))
	for i, cell := range headerRow.Cells {
		index[cell.String()] = i
	}
	return index
}
