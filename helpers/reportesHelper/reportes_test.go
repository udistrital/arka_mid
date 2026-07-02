package reportesHelper

import (
	"bytes"
	"encoding/base64"
	"net/url"
	"strings"
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
					CentroCostoNombre:   "Almacén General",
					CentroCostoCodigo:   "CC-001",
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
	if dataRow.Cells[headerIndex["Centro de costo"]].String() != "Almacén General" {
		t.Fatalf("centro de costo inesperado: %q", dataRow.Cells[headerIndex["Centro de costo"]].String())
	}
	if dataRow.Cells[headerIndex["codigo centro de costo"]].String() != "CC-001" {
		t.Fatalf("codigo centro de costo inesperado: %q", dataRow.Cells[headerIndex["codigo centro de costo"]].String())
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

func TestConsultarEntradasPorFechaFiltraPorFechaCreacion(t *testing.T) {
	original := consultarMovimientosReporteFn
	var capturedQuery string
	consultarMovimientosReporteFn = func(query string) ([]*models.Movimiento, string, map[string]interface{}) {
		capturedQuery = query
		return []*models.Movimiento{}, "0", nil
	}
	t.Cleanup(func() {
		consultarMovimientosReporteFn = original
	})

	_, err := consultarEntradasPorFecha(
		time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC),
		[]string{"ENT_ADQ"},
	)
	if err != nil {
		t.Fatalf("consultarEntradasPorFecha retornó error: %v", err)
	}

	decodedQuery, decodeErr := url.QueryUnescape(capturedQuery)
	if decodeErr != nil {
		t.Fatalf("no se pudo decodificar el query capturado: %v", decodeErr)
	}
	if !strings.Contains(decodedQuery, "FechaCreacion__gte:2026-04-01T00:00:00Z") {
		t.Fatalf("el query no filtró por FechaCreacion inicial: %q", decodedQuery)
	}
	if !strings.Contains(decodedQuery, "FechaCreacion__lte:2026-04-30T23:59:59Z") {
		t.Fatalf("el query no filtró por FechaCreacion final: %q", decodedQuery)
	}
	if strings.Contains(decodedQuery, "FechaCorte__gte") || strings.Contains(decodedQuery, "FechaCorte__lte") {
		t.Fatalf("el query no debe filtrar entradas normales por FechaCorte: %q", decodedQuery)
	}
}

func TestFechaReporteSalidaUsaFechaEntradaSiSalidaEsAnterior(t *testing.T) {
	fechaEntrada := time.Date(2026, 5, 4, 14, 43, 1, 0, time.UTC)
	fechaSalida := time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)

	fechaReporte := fechaReporteSalidaMovimiento(&models.Movimiento{
		FechaCorte: &fechaSalida,
		MovimientoPadreId: &models.Movimiento{
			FechaCreacion: fechaEntrada,
		},
	})

	if fechaReporte == nil {
		t.Fatal("se esperaba fecha de reporte")
	}
	if !fechaReporte.Equal(fechaEntrada) {
		t.Fatalf("se esperaba usar la fecha de la entrada %v, se obtuvo %v", fechaEntrada, *fechaReporte)
	}
}

func TestMovimientoEnRangoPorFechaReporteSalidaUsaFechaEntrada(t *testing.T) {
	movimiento := &models.Movimiento{
		FechaCorte: timePtr(time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)),
		MovimientoPadreId: &models.Movimiento{
			FechaCreacion: time.Date(2026, 5, 4, 14, 43, 1, 0, time.UTC),
		},
	}

	if !movimientoEnRangoPorFechaReporteSalida(
		movimiento,
		time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC),
	) {
		t.Fatal("se esperaba que la salida entrara en mayo usando la fecha de la entrada")
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
			CentroCostoNombre:   "Almacén General",
			CentroCostoCodigo:   "CC-001",
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
	if entradaAnuladaRow.Cells[headerIndex["Nombre / Descripción"]].String() != rotuloEntradaAnuladaReporte {
		t.Fatalf("rótulo entrada anulada inesperado: %q", entradaAnuladaRow.Cells[headerIndex["Nombre / Descripción"]].String())
	}
	if entradaAnuladaRow.Cells[headerIndex["Valor unitario"]].Value != "0" {
		t.Fatalf("valor unitario entrada anulada inesperado: %q", entradaAnuladaRow.Cells[headerIndex["Valor unitario"]].Value)
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
	if salidaAnuladaRow.Cells[headerIndex["Nombre / Descripción"]].String() != rotuloSalidaAnuladaReporte {
		t.Fatalf("rótulo salida anulada inesperado: %q", salidaAnuladaRow.Cells[headerIndex["Nombre / Descripción"]].String())
	}
	if salidaAnuladaRow.Cells[headerIndex["Subtotal"]].Value != "0" {
		t.Fatalf("subtotal salida anulada inesperado: %q", salidaAnuladaRow.Cells[headerIndex["Subtotal"]].Value)
	}
}

func TestGenerarReporteEntradaAnuladaConElementosGeneraSoloUnRenglon(t *testing.T) {
	mockConsultarEntradasReporteData(t, []*entradaReporteData{
		{
			Movimiento: &models.Movimiento{
				Id:            8103,
				Consecutivo:   stringPtr("ENT-ANU-02"),
				FechaCreacion: time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC),
				EstadoMovimientoId: &models.EstadoMovimiento{
					Nombre: estadoEntradaAnuladaReporte,
				},
				FormatoTipoMovimientoId: &models.FormatoTipoMovimiento{
					Nombre: "Entrada por compra",
				},
			},
			Formato: models.FormatoBaseEntrada{ActaRecibidoId: 999},
			Elementos: []*models.DetalleElemento{
				{Id: 1, Nombre: "Elemento que no debe salir"},
				{Id: 2, Nombre: "Segundo elemento que no debe salir"},
			},
			SalidasPorElemento: map[int]*salidaReporteData{},
			CuentasPorSubgrupo: map[int]models.CuentasSubgrupo{},
		},
	})
	mockConsultarEntradasAnuladasReporteData(t, []*entradaReporteData{
		{
			Movimiento: &models.Movimiento{
				Id:            8103,
				Consecutivo:   stringPtr("ENT-ANU-02"),
				FechaCreacion: time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC),
				EstadoMovimientoId: &models.EstadoMovimiento{
					Nombre: estadoEntradaAnuladaReporte,
				},
			},
		},
	})
	mockConsultarSalidasAnuladasReporteData(t, []*salidaReporteData{})

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

	if len(archivo.Sheets[0].Rows) != 2 {
		t.Fatalf("se esperaban 2 filas, se obtuvieron %d", len(archivo.Sheets[0].Rows))
	}

	headerIndex := buildHeaderIndex(archivo.Sheets[0].Rows[0])
	row := archivo.Sheets[0].Rows[1]
	if row.Cells[headerIndex["Consecutivo Entrada"]].String() != "ENT-ANU-02" {
		t.Fatalf("consecutivo entrada anulada inesperado: %q", row.Cells[headerIndex["Consecutivo Entrada"]].String())
	}
	if row.Cells[headerIndex["Nombre / Descripción"]].String() != rotuloEntradaAnuladaReporte {
		t.Fatalf("rotulo entrada anulada inesperado: %q", row.Cells[headerIndex["Nombre / Descripción"]].String())
	}
}

func TestExtraerCodigosEntradaIncluyeFormatosInactivos(t *testing.T) {
	codigos := extraerCodigosEntrada([]*models.FormatoTipoMovimiento{
		{CodigoAbreviacion: "ENT_RP", Activo: false},
		{CodigoAbreviacion: "ENT_ADQ", Activo: true},
		{CodigoAbreviacion: "ENT_KDX", Activo: true},
		{CodigoAbreviacion: "SAL", Activo: true},
	})

	if !containsString(codigos, "ENT_RP") {
		t.Fatalf("expected ENT_RP to be included, got %v", codigos)
	}
	if !containsString(codigos, "ENT_ADQ") {
		t.Fatalf("expected ENT_ADQ to be included, got %v", codigos)
	}
	if containsString(codigos, "ENT_KDX") {
		t.Fatalf("did not expect ENT_KDX to be included, got %v", codigos)
	}
	if containsString(codigos, "SAL") {
		t.Fatalf("did not expect SAL to be included as entrada, got %v", codigos)
	}
}

func TestSalidaUbicacionInfoUsaCentroCostosPorId(t *testing.T) {
	mockConsultarCentroCostos(t, []models.CentroCostos{
		{
			Id:     422,
			Codigo: "CC-422",
			Nombre: "Centro de costo principal",
		},
	})

	nombre, codigo := salidaUbicacionInfo(&models.Movimiento{
		Detalle: `{"funcionario":12345,"ubicacion":999,"centro_costos":"422"}`,
	})

	if nombre != "Centro de costo principal" {
		t.Fatalf("unexpected centro de costo: %q", nombre)
	}
	if codigo != "CC-422" {
		t.Fatalf("unexpected codigo centro de costo: %q", codigo)
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

func TestGenerarReporteContabilizacionRequestNil(t *testing.T) {
	respuesta, err := GenerarReporteContabilizacion(nil)
	if err == nil {
		t.Fatal("se esperaba error cuando request es nil")
	}
	if respuesta != nil {
		t.Fatal("no se esperaba respuesta cuando request es nil")
	}
}

func TestConsultarProveedorActaFallbackCuandoTerceroNoExiste(t *testing.T) {
	originalHistoricos := consultarHistoricosActaReporteFn
	originalGetNombre := getNombreTerceroByID

	consultarHistoricosActaReporteFn = func(query, fields, sortby, order, offset, limit string) ([]models.HistoricoActa, map[string]interface{}) {
		return []models.HistoricoActa{
			{ProveedorId: 123456},
		}, nil
	}
	getNombreTerceroByID = func(terceroID int) (*models.IdentificacionTercero, map[string]interface{}) {
		return nil, map[string]interface{}{
			"err":    "http 404: {\"Message\":\"Not found resource\"}",
			"status": "502",
		}
	}

	t.Cleanup(func() {
		consultarHistoricosActaReporteFn = originalHistoricos
		getNombreTerceroByID = originalGetNombre
	})

	proveedor, err := consultarProveedorActa(555)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if proveedor != "123456" {
		t.Fatalf("fallback de proveedor inesperado: %q", proveedor)
	}
}

func TestGenerarReporteContabilizacionFechaFinalMenor(t *testing.T) {
	respuesta, err := GenerarReporteContabilizacion(&models.ReporteFechasRequest{
		FechaInicial: "2026-06-16",
		FechaFinal:   "2026-06-01",
	})
	if err == nil {
		t.Fatal("se esperaba error cuando fecha_final es menor a fecha_inicial")
	}
	if respuesta != nil {
		t.Fatal("no se esperaba respuesta cuando el rango es inválido")
	}
}

func TestFechaDocumentoA11PriorizaFechaCorteMovimiento(t *testing.T) {
	fechaTransaccion := time.Date(2026, 6, 5, 10, 0, 0, 0, time.UTC)
	fechaCorte := time.Date(2026, 6, 7, 15, 30, 0, 0, time.UTC)

	fecha := fechaDocumentoA11(
		&models.InfoTransaccionContable{Fecha: fechaTransaccion},
		&models.Movimiento{FechaCorte: &fechaCorte},
	)

	if !fecha.Equal(fechaCorte) {
		t.Fatalf("se esperaba priorizar FechaCorte=%v, se obtuvo %v", fechaCorte, fecha)
	}
}

func TestGenerarReporteContabilizacionEncabezadosYRenglones(t *testing.T) {
	mockConsultarCuentaContableReporte(t, map[string]*models.CuentaContable{
		"cta-db-1": {Id: "cta-db-1", Codigo: "151001"},
		"cta-cr-1": {Id: "cta-cr-1", Codigo: "240801"},
	})
	mockConsultarEntradasContabilizacionReporteData(t, []*reporteContabilizacionEntradaData{
		{
			Movimiento: &models.Movimiento{
				Consecutivo:   stringPtr("P8-00001-2026"),
				Observacion:   "Observación entrada contable",
				FechaCreacion: time.Date(2026, 6, 5, 9, 0, 0, 0, time.UTC),
			},
			ProveedorLabel:     "900866324 - PROVEEDOR UNO",
			FacturaConsecutivo: "FV 434",
			CentroCostoCodigo:  "A1205",
			CentroCostoNombre:  "Oficina Asesora de Planeación",
			TransaccionContable: &models.InfoTransaccionContable{
				Fecha: time.Date(2026, 6, 5, 10, 0, 0, 0, time.UTC),
				Movimientos: []*models.DetalleMovimientoContable{
					{
						Cuenta:    &models.DetalleCuenta{Id: "asiento-db-1", Codigo: "999999"},
						Debito:    1000,
						TerceroId: &models.IdentificacionTercero{Numero: "800123", NombreCompleto: "TERCERO CONTABLE"},
					},
					{
						Cuenta:  &models.DetalleCuenta{Id: "asiento-cr-1", Codigo: "888888"},
						Credito: 1000,
					},
				},
			},
			CuentasPorSubgrupo: map[int]models.CuentasSubgrupo{
				9: {
					CuentaDebitoId:  "cta-db-1",
					CuentaCreditoId: "cta-cr-1",
				},
			},
			Subgrupos:         []int{9},
			CuentasDebitoA11:  []string{"151001"},
			CuentasCreditoA11: []string{"240801"},
		},
	})
	mockConsultarSalidasContabilizacionReporteData(t, []*reporteContabilizacionSalidaData{
		{
			Movimiento: &models.Movimiento{
				Consecutivo:   stringPtr("H21-00001-2026"),
				Observacion:   "Observación salida contable",
				FechaCreacion: time.Date(2026, 6, 6, 11, 0, 0, 0, time.UTC),
			},
			EntradaPadre: &models.Movimiento{
				Consecutivo: stringPtr("P8-00001-2026"),
			},
			ProveedorLabel:     "900866324 - PROVEEDOR UNO",
			FacturaConsecutivo: "FV 434",
			CentroCostoCodigo:  "A1205",
			CentroCostoNombre:  "Oficina Asesora de Planeación",
			TransaccionContable: &models.InfoTransaccionContable{
				Fecha: time.Date(2026, 6, 6, 12, 0, 0, 0, time.UTC),
				Movimientos: []*models.DetalleMovimientoContable{
					{Cuenta: &models.DetalleCuenta{Codigo: "77777777"}, Debito: 300},
					{Cuenta: &models.DetalleCuenta{Codigo: "66666666"}, Credito: 300},
				},
			},
			CuentasDebitoA11:  []string{"51111401"},
			CuentasCreditoA11: []string{"15142401"},
		},
		{
			Movimiento: &models.Movimiento{
				Consecutivo:   stringPtr("H21-00002-2026"),
				FechaCreacion: time.Date(2026, 6, 7, 11, 0, 0, 0, time.UTC),
			},
			EntradaPadre: &models.Movimiento{
				Consecutivo: stringPtr("P8-00002-2026"),
			},
			ProveedorLabel:     "900999111 - PROVEEDOR DOS",
			FacturaConsecutivo: "FV 999",
			CentroCostoCodigo:  "A1302",
			CentroCostoNombre:  "Dependencia Dos",
			TransaccionContable: &models.InfoTransaccionContable{
				Fecha: time.Date(2026, 6, 7, 12, 0, 0, 0, time.UTC),
				Movimientos: []*models.DetalleMovimientoContable{
					{Cuenta: &models.DetalleCuenta{Codigo: "8315100728"}, Debito: 500},
				},
			},
		},
	})

	respuesta, err := GenerarReporteContabilizacion(&models.ReporteFechasRequest{
		FechaInicial: "2026-06-01",
		FechaFinal:   "2026-06-16",
	})
	if err != nil {
		t.Fatalf("GenerarReporteContabilizacion retornó error: %v", err)
	}

	contenido, decodeErr := base64.StdEncoding.DecodeString(respuesta.ArchivoBase64)
	if decodeErr != nil {
		t.Fatalf("base64 inválido: %v", decodeErr)
	}

	archivo, openErr := xlsx.OpenBinary(contenido)
	if openErr != nil {
		t.Fatalf("no se pudo abrir el excel generado: %v", openErr)
	}

	if len(archivo.Sheets) != 2 {
		t.Fatalf("se esperaban 2 hojas, se obtuvieron %d", len(archivo.Sheets))
	}
	if archivo.Sheets[0].Name != sheetNameContabilizacionEntradas {
		t.Fatalf("nombre de hoja entradas inesperado: %q", archivo.Sheets[0].Name)
	}
	if archivo.Sheets[1].Name != sheetNameContabilizacionSalidas {
		t.Fatalf("nombre de hoja salidas inesperado: %q", archivo.Sheets[1].Name)
	}

	headerRow := archivo.Sheets[0].Rows[0]
	if len(headerRow.Cells) != 18 {
		t.Fatalf("se esperaban 18 encabezados, se obtuvieron %d", len(headerRow.Cells))
	}
	if headerRow.Cells[0].String() != "Renglón" || headerRow.Cells[17].String() != "Jefe_Almacén" {
		t.Fatalf("encabezados A11 inesperados")
	}

	entradasIndex := buildHeaderIndex(archivo.Sheets[0].Rows[0])
	entradasDataRow := archivo.Sheets[0].Rows[1]
	if entradasDataRow.Cells[entradasIndex["Documento_Entrada"]].String() != "P8-00001-2026" {
		t.Fatalf("documento entrada inesperado: %q", entradasDataRow.Cells[entradasIndex["Documento_Entrada"]].String())
	}
	if entradasDataRow.Cells[entradasIndex["Documento_Salida"]].String() != "" {
		t.Fatalf("documento salida de entrada debe venir vacío: %q", entradasDataRow.Cells[entradasIndex["Documento_Salida"]].String())
	}
	if entradasDataRow.Cells[entradasIndex["Identificación_Tercero"]].String() != "800123" {
		t.Fatalf("tercero contable inesperado: %q", entradasDataRow.Cells[entradasIndex["Identificación_Tercero"]].String())
	}
	if entradasDataRow.Cells[entradasIndex["Centro_Costo"]].String() != "A1205" {
		t.Fatalf("centro de costo entrada inesperado: %q", entradasDataRow.Cells[entradasIndex["Centro_Costo"]].String())
	}
	if entradasDataRow.Cells[entradasIndex["Cuenta_Contable"]].String() != "151001" {
		t.Fatalf("cuenta contable entrada inesperada: %q", entradasDataRow.Cells[entradasIndex["Cuenta_Contable"]].String())
	}
	if entradasDataRow.Cells[entradasIndex["Descripción"]].String() != "Observación entrada contable" {
		t.Fatalf("descripción entrada inesperada: %q", entradasDataRow.Cells[entradasIndex["Descripción"]].String())
	}
	if archivo.Sheets[0].Rows[2].Cells[entradasIndex["Cuenta_Contable"]].String() != "240801" {
		t.Fatalf("cuenta contable crédito entrada inesperada: %q", archivo.Sheets[0].Rows[2].Cells[entradasIndex["Cuenta_Contable"]].String())
	}
	if entradasDataRow.Cells[entradasIndex["Fecha_Documento"]].GetNumberFormat() != a11DateNumFmt {
		t.Fatalf("formato fecha inesperado: %q", entradasDataRow.Cells[entradasIndex["Fecha_Documento"]].GetNumberFormat())
	}

	salidasIndex := buildHeaderIndex(archivo.Sheets[1].Rows[0])
	if archivo.Sheets[1].Rows[1].Cells[salidasIndex["Renglón"]].String() != "1" {
		t.Fatalf("primer renglón de salida inesperado")
	}
	if archivo.Sheets[1].Rows[2].Cells[salidasIndex["Renglón"]].String() != "2" {
		t.Fatalf("segundo renglón de salida inesperado")
	}
	if archivo.Sheets[1].Rows[3].Cells[salidasIndex["Renglón"]].String() != "1" {
		t.Fatalf("el renglón debe reiniciarse por salida")
	}
	if archivo.Sheets[1].Rows[1].Cells[salidasIndex["Documento_Salida"]].String() != "H21-00001-2026" {
		t.Fatalf("documento salida inesperado: %q", archivo.Sheets[1].Rows[1].Cells[salidasIndex["Documento_Salida"]].String())
	}
	if archivo.Sheets[1].Rows[1].Cells[salidasIndex["Descripción"]].String() != "Observación salida contable" {
		t.Fatalf("descripción salida inesperada: %q", archivo.Sheets[1].Rows[1].Cells[salidasIndex["Descripción"]].String())
	}
	if archivo.Sheets[1].Rows[1].Cells[salidasIndex["Documento_Entrada"]].String() != "P8-00001-2026" {
		t.Fatalf("documento entrada asociado inesperado: %q", archivo.Sheets[1].Rows[1].Cells[salidasIndex["Documento_Entrada"]].String())
	}
	if archivo.Sheets[1].Rows[2].Cells[salidasIndex["Naturaleza_Cuenta"]].String() != "C" {
		t.Fatalf("naturaleza cuenta inesperada: %q", archivo.Sheets[1].Rows[2].Cells[salidasIndex["Naturaleza_Cuenta"]].String())
	}
	if archivo.Sheets[1].Rows[1].Cells[salidasIndex["Cuenta_Contable"]].String() != "51111401" {
		t.Fatalf("cuenta débito salida inesperada: %q", archivo.Sheets[1].Rows[1].Cells[salidasIndex["Cuenta_Contable"]].String())
	}
	if archivo.Sheets[1].Rows[2].Cells[salidasIndex["Cuenta_Contable"]].String() != "15142401" {
		t.Fatalf("cuenta crédito salida inesperada: %q", archivo.Sheets[1].Rows[2].Cells[salidasIndex["Cuenta_Contable"]].String())
	}
	if archivo.Sheets[1].Rows[2].Cells[salidasIndex["Identificación_Tercero"]].String() != "900866324" {
		t.Fatalf("fallback de tercero inesperado: %q", archivo.Sheets[1].Rows[2].Cells[salidasIndex["Identificación_Tercero"]].String())
	}
}

func TestGenerarReporteContabilizacionOmiteDocumentoSinTransaccion(t *testing.T) {
	mockConsultarEntradasContabilizacionReporteData(t, []*reporteContabilizacionEntradaData{
		{
			Movimiento:          &models.Movimiento{Consecutivo: stringPtr("P8-00003-2026")},
			TransaccionContable: nil,
		},
	})
	mockConsultarSalidasContabilizacionReporteData(t, []*reporteContabilizacionSalidaData{
		{
			Movimiento: &models.Movimiento{Consecutivo: stringPtr("H21-00003-2026")},
			TransaccionContable: &models.InfoTransaccionContable{
				Movimientos: []*models.DetalleMovimientoContable{},
			},
		},
	})

	respuesta, err := GenerarReporteContabilizacion(&models.ReporteFechasRequest{
		FechaInicial: "2026-06-01",
		FechaFinal:   "2026-06-16",
	})
	if err != nil {
		t.Fatalf("GenerarReporteContabilizacion retornó error: %v", err)
	}

	contenido, decodeErr := base64.StdEncoding.DecodeString(respuesta.ArchivoBase64)
	if decodeErr != nil {
		t.Fatalf("base64 inválido: %v", decodeErr)
	}

	archivo, openErr := xlsx.OpenBinary(contenido)
	if openErr != nil {
		t.Fatalf("no se pudo abrir el excel generado: %v", openErr)
	}

	if len(archivo.Sheets[0].Rows) != 1 {
		t.Fatalf("la hoja entradas solo debe contener encabezados, tiene %d filas", len(archivo.Sheets[0].Rows))
	}
	if len(archivo.Sheets[1].Rows) != 1 {
		t.Fatalf("la hoja salidas solo debe contener encabezados, tiene %d filas", len(archivo.Sheets[1].Rows))
	}
}

func TestSplitTerceroLabel(t *testing.T) {
	identificacion, nombre := splitTerceroLabel("900866324 - GESTION Y ESTUDIOS AMBIENTALES")
	if identificacion != "900866324" {
		t.Fatalf("identificación inesperada: %q", identificacion)
	}
	if nombre != "GESTION Y ESTUDIOS AMBIENTALES" {
		t.Fatalf("nombre inesperado: %q", nombre)
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

func TestCuentaDesdeDetalleA11MultiplesCuentas(t *testing.T) {
	if cuenta := cuentaDesdeDetalleA11("D", []string{"151001", "151002"}, nil, 0, 0); cuenta != "151001" {
		t.Fatalf("cuenta débito índice 0 inesperada: %q", cuenta)
	}
	if cuenta := cuentaDesdeDetalleA11("D", []string{"151001", "151002"}, nil, 1, 0); cuenta != "151002" {
		t.Fatalf("cuenta débito índice 1 inesperada: %q", cuenta)
	}
	if cuenta := cuentaDesdeDetalleA11("C", nil, []string{"240801", "240802"}, 0, 1); cuenta != "240802" {
		t.Fatalf("cuenta crédito índice 1 inesperada: %q", cuenta)
	}
}

func mockConsultarEntradasContabilizacionReporteData(t *testing.T, entradas []*reporteContabilizacionEntradaData) {
	t.Helper()

	original := consultarEntradasContabilizacionReporteData
	consultarEntradasContabilizacionReporteData = func(fechaInicial, fechaFinal time.Time) ([]*reporteContabilizacionEntradaData, map[string]interface{}) {
		return entradas, nil
	}

	t.Cleanup(func() {
		consultarEntradasContabilizacionReporteData = original
	})
}

func mockConsultarSalidasContabilizacionReporteData(t *testing.T, salidas []*reporteContabilizacionSalidaData) {
	t.Helper()

	original := consultarSalidasContabilizacionReporteData
	consultarSalidasContabilizacionReporteData = func(fechaInicial, fechaFinal time.Time) ([]*reporteContabilizacionSalidaData, map[string]interface{}) {
		return salidas, nil
	}

	t.Cleanup(func() {
		consultarSalidasContabilizacionReporteData = original
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

func mockConsultarCentroCostos(t *testing.T, centros []models.CentroCostos) {
	t.Helper()

	original := consultarCentroCostosFn
	consultarCentroCostosFn = func(payload string) ([]models.CentroCostos, map[string]interface{}) {
		if payload != "query=Id:422" {
			t.Fatalf("unexpected centro costos payload: %s", payload)
		}
		return centros, nil
	}

	t.Cleanup(func() {
		consultarCentroCostosFn = original
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

func mockConsultarCuentaContableReporte(t *testing.T, cuentas map[string]*models.CuentaContable) {
	t.Helper()

	original := consultarCuentaContable
	consultarCuentaContable = func(cuentaID string) (*models.CuentaContable, map[string]interface{}) {
		if cuenta, ok := cuentas[cuentaID]; ok {
			return cuenta, nil
		}
		return nil, nil
	}

	t.Cleanup(func() {
		consultarCuentaContable = original
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

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
