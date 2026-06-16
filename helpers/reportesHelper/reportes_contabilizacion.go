package reportesHelper

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/tealeg/xlsx"
	crudActaRecibido "github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

const (
	sheetNameContabilizacionEntradas = "entradas"
	sheetNameContabilizacionSalidas  = "Salidas"
	a11DateNumFmt                    = "dd/mm/yyyy"
)

type reporteContabilizacionEntradaData struct {
	Movimiento          *models.Movimiento
	TransaccionContable *models.InfoTransaccionContable
	ProveedorLabel      string
	FacturaConsecutivo  string
	CentroCostoCodigo   string
	CentroCostoNombre   string
}

type reporteContabilizacionSalidaData struct {
	Movimiento          *models.Movimiento
	EntradaPadre        *models.Movimiento
	TransaccionContable *models.InfoTransaccionContable
	ProveedorLabel      string
	FacturaConsecutivo  string
	CentroCostoCodigo   string
	CentroCostoNombre   string
}

type reporteA11ContabilizacionRow struct {
	Renglon                        int
	CuentaContable                 string
	NaturalezaCuenta               string
	Descripcion                    string
	Valor                          float64
	IdentificacionTercero          string
	NombreTercero                  string
	CentroCosto                    string
	DescripcionCentroCosto         string
	Proyecto                       string
	DescripcionProyecto            string
	ClaseDocumentoLinea            string
	Empresa                        string
	TipoDocumentoCodigo            string
	TipoDocumentoNombre            string
	ClaseDocumentoMovimientoCodigo string
	ClaseDocumentoMovimientoNombre string
	EstadoCodigo                   string
	EstadoNombre                   string
	Detalle                        string
	DocumentoSalida                string
	FechaDocumento                 time.Time
	DocumentoEntrada               string
	Factura                        string
	Vigencia                       string
	Usuario                        string
	JefeAlmacen                    string
}

var (
	reporteA11Headers = []string{
		"Renglón",
		"Cuenta_Contable",
		"Naturaleza_Cuenta",
		"Descripción",
		"Valor",
		"Identificación_Tercero",
		"Nombre_Tercero",
		"Centro_Costo",
		"Descripción_Centro_Costo",
		"Proyecto",
		"Descripción_Proyecto",
		"Clase_Documento",
		"Empresa",
		"Tipo_Documento",
		"Tipo_Documento",
		"Clase_Documento",
		"Clase_Documento",
		"Estado",
		"Estado",
		"Detalle",
		"Documento_Salida",
		"Fecha_Documento",
		"Documento_Entrada",
		"Factura",
		"Vigencia",
		"Usuario",
		"Jefe_Almacén",
	}

	consultarEntradasContabilizacionReporteData = consultarEntradasContabilizacionReporteDataDefault
	consultarSalidasContabilizacionReporteData  = consultarSalidasContabilizacionReporteDataDefault
	consultarHistoricosActaReporteFn            = crudActaRecibido.GetAllHistoricoActa
)

func GenerarReporteContabilizacion(req *models.ReporteFechasRequest) (respuesta *models.ReporteExcelBase64Response, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("GenerarReporteContabilizacion - Unhandled Error!", "500")

	if req == nil {
		return nil, errorCtrl.Error("GenerarReporteContabilizacion - req", "request nil", "400")
	}

	fechaInicial, err := time.Parse(dateLayout, req.FechaInicial)
	if err != nil {
		return nil, errorCtrl.Error("GenerarReporteContabilizacion - time.Parse(fecha_inicial)", err, "400")
	}

	fechaFinal, err := time.Parse(dateLayout, req.FechaFinal)
	if err != nil {
		return nil, errorCtrl.Error("GenerarReporteContabilizacion - time.Parse(fecha_final)", err, "400")
	}

	if fechaFinal.Before(fechaInicial) {
		return nil, errorCtrl.Error("GenerarReporteContabilizacion - rango_fechas", "fecha_final debe ser mayor o igual a fecha_inicial", "400")
	}

	entradas, outputError := consultarEntradasContabilizacionReporteData(fechaInicial, fechaFinal)
	if outputError != nil {
		return nil, outputError
	}

	salidas, outputError := consultarSalidasContabilizacionReporteData(fechaInicial, fechaFinal)
	if outputError != nil {
		return nil, outputError
	}

	archivo := xlsx.NewFile()
	hojaEntradas, err := archivo.AddSheet(sheetNameContabilizacionEntradas)
	if err != nil {
		return nil, errorCtrl.Error("GenerarReporteContabilizacion - archivo.AddSheet(entradas)", err, "500")
	}

	hojaSalidas, err := archivo.AddSheet(sheetNameContabilizacionSalidas)
	if err != nil {
		return nil, errorCtrl.Error("GenerarReporteContabilizacion - archivo.AddSheet(Salidas)", err, "500")
	}

	addReporteA11Headers(hojaEntradas)
	addReporteA11Headers(hojaSalidas)

	for _, row := range construirFilasReporteA11Entradas(entradas) {
		addReporteA11Row(hojaEntradas, row)
	}
	for _, row := range construirFilasReporteA11Salidas(salidas) {
		addReporteA11Row(hojaSalidas, row)
	}

	setA11ColumnWidths(hojaEntradas)
	setA11ColumnWidths(hojaSalidas)

	var buffer bytes.Buffer
	if err := archivo.Write(&buffer); err != nil {
		return nil, errorCtrl.Error("GenerarReporteContabilizacion - archivo.Write", err, "500")
	}

	return &models.ReporteExcelBase64Response{
		ArchivoBase64: base64.StdEncoding.EncodeToString(buffer.Bytes()),
		NombreArchivo: fmt.Sprintf("reporte_contabilizacion_%s_%s.xlsx", fechaInicial.Format("20060102"), fechaFinal.Format("20060102")),
		TipoArchivo:   "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}, nil
}

func consultarEntradasContabilizacionReporteDataDefault(fechaInicial, fechaFinal time.Time) (entradas []*reporteContabilizacionEntradaData, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarEntradasContabilizacionReporteDataDefault - Unhandled Error!", "500")

	entradasBase, outputError := consultarEntradasReporteDataDefault(fechaInicial, fechaFinal)
	if outputError != nil {
		return nil, outputError
	}

	entradas = make([]*reporteContabilizacionEntradaData, 0, len(entradasBase))
	for _, entrada := range entradasBase {
		if entrada == nil || entrada.Movimiento == nil {
			continue
		}
		if !movimientoTieneEstado(entrada.Movimiento, "Entrada Aprobada", "Entrada Con Salida") {
			continue
		}
		if !entradaTieneSalidas(entrada) {
			continue
		}
		if entrada.TransaccionContable == nil || len(entrada.TransaccionContable.Movimientos) == 0 {
			continue
		}

		validarCuadreTransaccionA11("entrada", stringPtrValue(entrada.Movimiento.Consecutivo), entrada.TransaccionContable)

		centroCostoCodigo, centroCostoNombre, outputError := centroCostoEntradaA11Info(entrada.Formato)
		if outputError != nil {
			return nil, outputError
		}

		entradas = append(entradas, &reporteContabilizacionEntradaData{
			Movimiento:          entrada.Movimiento,
			TransaccionContable: entrada.TransaccionContable,
			ProveedorLabel:      entrada.Proveedor,
			FacturaConsecutivo:  strings.TrimSpace(entrada.FacturaConsecutivo),
			CentroCostoCodigo:   centroCostoCodigo,
			CentroCostoNombre:   centroCostoNombre,
		})
	}

	return entradas, nil
}

func consultarSalidasContabilizacionReporteDataDefault(fechaInicial, fechaFinal time.Time) (salidas []*reporteContabilizacionSalidaData, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarSalidasContabilizacionReporteDataDefault - Unhandled Error!", "500")

	codigosSalida, outputError := consultarCodigosSalida()
	if outputError != nil {
		return nil, outputError
	}
	if len(codigosSalida) == 0 {
		return []*reporteContabilizacionSalidaData{}, nil
	}

	movimientos, outputError := consultarSalidasPorFecha(fechaInicial, fechaFinal, codigosSalida)
	if outputError != nil {
		return nil, outputError
	}

	salidas = make([]*reporteContabilizacionSalidaData, 0, len(movimientos))
	for _, movimiento := range movimientos {
		if movimiento == nil || movimiento.Id <= 0 {
			continue
		}

		transaccion := consultarTransaccionContableSalida(movimiento)
		if transaccion == nil || len(transaccion.Movimientos) == 0 {
			continue
		}

		validarCuadreTransaccionA11("salida", stringPtrValue(movimiento.Consecutivo), transaccion)

		entradaPadre, proveedorLabel, facturaConsecutivo, outputError := metadataEntradaPadreSalida(movimiento)
		if outputError != nil {
			return nil, outputError
		}

		centroCostoCodigo, centroCostoNombre, outputError := centroCostoMovimientoA11Info(movimiento)
		if outputError != nil {
			return nil, outputError
		}

		salidas = append(salidas, &reporteContabilizacionSalidaData{
			Movimiento:          movimiento,
			EntradaPadre:        entradaPadre,
			TransaccionContable: transaccion,
			ProveedorLabel:      proveedorLabel,
			FacturaConsecutivo:  strings.TrimSpace(facturaConsecutivo),
			CentroCostoCodigo:   centroCostoCodigo,
			CentroCostoNombre:   centroCostoNombre,
		})
	}

	return salidas, nil
}

func consultarSalidasPorFecha(fechaInicial, fechaFinal time.Time, codigosSalida []string) (movimientos []*models.Movimiento, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarSalidasPorFecha - Unhandled Error!", "500")

	fechaInicial = time.Date(fechaInicial.Year(), fechaInicial.Month(), fechaInicial.Day(), 0, 0, 0, 0, time.UTC)
	fechaFinal = time.Date(fechaFinal.Year(), fechaFinal.Month(), fechaFinal.Day(), 23, 59, 59, 0, time.UTC)

	params := url.Values{}
	params.Add("limit", "-1")
	params.Add("sortby", "FechaCreacion")
	params.Add("order", "asc")
	params.Add(
		"query",
		"Activo:true,EstadoMovimientoId__Nombre:Salida Aprobada,FormatoTipoMovimientoId__CodigoAbreviacion__in:"+strings.Join(codigosSalida, "|")+
			",FechaCreacion__gte:"+fechaInicial.Format(time.RFC3339)+
			",FechaCreacion__lte:"+fechaFinal.Format(time.RFC3339),
	)

	movimientos, _, outputError = movimientosArka.GetAllMovimiento(params.Encode())
	if outputError != nil {
		return nil, outputError
	}

	return movimientos, nil
}

func metadataEntradaPadreSalida(movimiento *models.Movimiento) (entradaPadre *models.Movimiento, proveedorLabel, facturaConsecutivo string, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("metadataEntradaPadreSalida - Unhandled Error!", "500")

	if movimiento == nil || movimiento.MovimientoPadreId == nil || movimiento.MovimientoPadreId.Id <= 0 {
		return nil, "", "", nil
	}

	entradaPadre = movimiento.MovimientoPadreId
	if strings.TrimSpace(entradaPadre.Detalle) == "" || entradaPadre.FormatoTipoMovimientoId == nil {
		entradaPadre, outputError = consultarMovimientoPorID(movimiento.MovimientoPadreId.Id)
		if outputError != nil {
			return nil, "", "", outputError
		}
	}
	if entradaPadre == nil {
		return nil, "", "", nil
	}

	var formato models.FormatoBaseEntrada
	if strings.TrimSpace(entradaPadre.Detalle) != "" {
		outputError = utilsHelper.Unmarshal(entradaPadre.Detalle, &formato)
		if outputError != nil {
			return nil, "", "", outputError
		}
	}

	proveedorLabel, facturaConsecutivo, _, outputError = consultarMetadataEntradaFn(formato)
	if outputError != nil {
		return nil, "", "", outputError
	}

	return entradaPadre, proveedorLabel, facturaConsecutivo, nil
}

func construirFilasReporteA11Entradas(entradas []*reporteContabilizacionEntradaData) []*reporteA11ContabilizacionRow {
	rows := make([]*reporteA11ContabilizacionRow, 0)
	for _, entrada := range entradas {
		rows = append(rows, construirFilasA11PorEntrada(entrada)...)
	}
	return rows
}

func construirFilasReporteA11Salidas(salidas []*reporteContabilizacionSalidaData) []*reporteA11ContabilizacionRow {
	rows := make([]*reporteA11ContabilizacionRow, 0)
	for _, salida := range salidas {
		rows = append(rows, construirFilasA11PorSalida(salida)...)
	}
	return rows
}

func construirFilasA11PorEntrada(entrada *reporteContabilizacionEntradaData) []*reporteA11ContabilizacionRow {
	if entrada == nil || entrada.Movimiento == nil || entrada.TransaccionContable == nil {
		return nil
	}

	documentoEntrada := stringPtrValue(entrada.Movimiento.Consecutivo)
	fechaDocumento := fechaDocumentoA11(entrada.TransaccionContable, entrada.Movimiento)
	vigencia := strconv.Itoa(fechaDocumento.Year())
	if fechaDocumento.IsZero() {
		vigencia = strconv.Itoa(entrada.Movimiento.FechaCreacion.Year())
	}

	rows := make([]*reporteA11ContabilizacionRow, 0, len(entrada.TransaccionContable.Movimientos))
	renglon := 1
	for _, movimientoContable := range entrada.TransaccionContable.Movimientos {
		naturaleza, valor := naturalezaYValorMovimientoA11(movimientoContable)
		if naturaleza == "" {
			continue
		}

		identificacion, nombre := terceroMovimientoA11(movimientoContable, entrada.ProveedorLabel)
		rows = append(rows, &reporteA11ContabilizacionRow{
			Renglon:                renglon,
			CuentaContable:         codigoCuentaContableA11(movimientoContable),
			NaturalezaCuenta:       naturaleza,
			Descripcion:            "",
			Valor:                  valor,
			IdentificacionTercero:  identificacion,
			NombreTercero:          nombre,
			CentroCosto:            entrada.CentroCostoCodigo,
			DescripcionCentroCosto: entrada.CentroCostoNombre,
			Detalle:                "",
			FechaDocumento:         fechaDocumento,
			DocumentoEntrada:       documentoEntrada,
			Factura:                strings.TrimSpace(entrada.FacturaConsecutivo),
			Vigencia:               vigencia,
		})
		renglon++
	}

	return rows
}

func construirFilasA11PorSalida(salida *reporteContabilizacionSalidaData) []*reporteA11ContabilizacionRow {
	if salida == nil || salida.Movimiento == nil || salida.TransaccionContable == nil {
		return nil
	}

	documentoSalida := stringPtrValue(salida.Movimiento.Consecutivo)
	documentoEntrada := ""
	if salida.EntradaPadre != nil {
		documentoEntrada = stringPtrValue(salida.EntradaPadre.Consecutivo)
	}
	fechaDocumento := fechaDocumentoA11(salida.TransaccionContable, salida.Movimiento)
	vigencia := strconv.Itoa(fechaDocumento.Year())
	if fechaDocumento.IsZero() {
		vigencia = strconv.Itoa(salida.Movimiento.FechaCreacion.Year())
	}

	rows := make([]*reporteA11ContabilizacionRow, 0, len(salida.TransaccionContable.Movimientos))
	renglon := 1
	for _, movimientoContable := range salida.TransaccionContable.Movimientos {
		naturaleza, valor := naturalezaYValorMovimientoA11(movimientoContable)
		if naturaleza == "" {
			continue
		}

		identificacion, nombre := terceroMovimientoA11(movimientoContable, salida.ProveedorLabel)
		rows = append(rows, &reporteA11ContabilizacionRow{
			Renglon:                renglon,
			CuentaContable:         codigoCuentaContableA11(movimientoContable),
			NaturalezaCuenta:       naturaleza,
			Descripcion:            "",
			Valor:                  valor,
			IdentificacionTercero:  identificacion,
			NombreTercero:          nombre,
			CentroCosto:            salida.CentroCostoCodigo,
			DescripcionCentroCosto: salida.CentroCostoNombre,
			Detalle:                "",
			DocumentoSalida:        documentoSalida,
			FechaDocumento:         fechaDocumento,
			DocumentoEntrada:       documentoEntrada,
			Factura:                strings.TrimSpace(salida.FacturaConsecutivo),
			Vigencia:               vigencia,
		})
		renglon++
	}

	return rows
}

func addReporteA11Headers(hoja *xlsx.Sheet) {
	filaEncabezado := hoja.AddRow()
	for _, header := range reporteA11Headers {
		filaEncabezado.AddCell().SetString(header)
	}
}

func addReporteA11Row(hoja *xlsx.Sheet, rowData *reporteA11ContabilizacionRow) {
	if rowData == nil {
		return
	}

	row := hoja.AddRow()
	addStringCell(row, strconv.Itoa(rowData.Renglon))
	addStringCell(row, rowData.CuentaContable)
	addStringCell(row, rowData.NaturalezaCuenta)
	addStringCell(row, rowData.Descripcion)
	addDecimalCell(row, rowData.Valor)
	addStringCell(row, rowData.IdentificacionTercero)
	addStringCell(row, rowData.NombreTercero)
	addStringCell(row, rowData.CentroCosto)
	addStringCell(row, rowData.DescripcionCentroCosto)
	addStringCell(row, rowData.Proyecto)
	addStringCell(row, rowData.DescripcionProyecto)
	addStringCell(row, rowData.ClaseDocumentoLinea)
	addStringCell(row, rowData.Empresa)
	addStringCell(row, rowData.TipoDocumentoCodigo)
	addStringCell(row, rowData.TipoDocumentoNombre)
	addStringCell(row, rowData.ClaseDocumentoMovimientoCodigo)
	addStringCell(row, rowData.ClaseDocumentoMovimientoNombre)
	addStringCell(row, rowData.EstadoCodigo)
	addStringCell(row, rowData.EstadoNombre)
	addStringCell(row, rowData.Detalle)
	addStringCell(row, rowData.DocumentoSalida)
	addDateCell(row, rowData.FechaDocumento)
	addStringCell(row, rowData.DocumentoEntrada)
	addStringCell(row, rowData.Factura)
	addStringCell(row, rowData.Vigencia)
	addStringCell(row, rowData.Usuario)
	addStringCell(row, rowData.JefeAlmacen)
}

func setA11ColumnWidths(hoja *xlsx.Sheet) {
	_ = hoja.SetColWidth(0, 0, 10)
	_ = hoja.SetColWidth(1, 1, 18)
	_ = hoja.SetColWidth(2, 2, 16)
	_ = hoja.SetColWidth(3, 3, 18)
	_ = hoja.SetColWidth(4, 4, 14)
	_ = hoja.SetColWidth(5, 6, 24)
	_ = hoja.SetColWidth(7, 8, 26)
	_ = hoja.SetColWidth(9, 19, 20)
	_ = hoja.SetColWidth(20, 24, 18)
	_ = hoja.SetColWidth(25, 26, 18)
}

func addDateCell(row *xlsx.Row, value time.Time) {
	cell := row.AddCell()
	if value.IsZero() {
		cell.SetString("")
		return
	}
	cell.SetDate(value)
	cell.NumFmt = a11DateNumFmt
}

func splitTerceroLabel(label string) (identificacion, nombre string) {
	label = strings.TrimSpace(label)
	if label == "" {
		return "", ""
	}

	partes := strings.SplitN(label, " - ", 2)
	if len(partes) == 2 {
		return strings.TrimSpace(partes[0]), strings.TrimSpace(partes[1])
	}

	return label, ""
}

func documentoMovimientoA11(movimiento *models.Movimiento) string {
	if movimiento == nil {
		return ""
	}
	return strings.TrimSpace(stringPtrValue(movimiento.Consecutivo))
}

func fechaDocumentoA11(transaccion *models.InfoTransaccionContable, movimiento *models.Movimiento) time.Time {
	if transaccion != nil && !transaccion.Fecha.IsZero() {
		return transaccion.Fecha
	}
	if movimiento != nil && movimiento.FechaCorte != nil && !movimiento.FechaCorte.IsZero() {
		return *movimiento.FechaCorte
	}
	if movimiento != nil {
		return movimiento.FechaCreacion
	}
	return time.Time{}
}

func validarCuadreTransaccionA11(tipoDocumento, documento string, transaccion *models.InfoTransaccionContable) bool {
	if transaccion == nil {
		return true
	}

	var debito, credito float64
	for _, movimiento := range transaccion.Movimientos {
		if movimiento == nil {
			continue
		}
		debito += movimiento.Debito
		credito += movimiento.Credito
	}

	if diff := roundToTwoDecimals(debito - credito); diff > 0.01 || diff < -0.01 {
		logs.Warn("A11 %s %s descuadrado: debito=%f credito=%f diferencia=%f", tipoDocumento, documento, debito, credito, diff)
		return false
	}

	return true
}

func entradaTieneSalidas(entrada *entradaReporteData) bool {
	if entrada == nil || len(entrada.SalidasPorElemento) == 0 {
		return false
	}
	for _, salida := range entrada.SalidasPorElemento {
		if salida != nil && salida.Movimiento != nil {
			return true
		}
	}
	return false
}

func centroCostoEntradaA11Info(formato models.FormatoBaseEntrada) (codigo, nombre string, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("centroCostoEntradaA11Info - Unhandled Error!", "500")

	if formato.ActaRecibidoId <= 0 {
		return "", "", nil
	}

	historicos, outputError := consultarHistoricosActaReporteFn(
		"ActaRecibidoId__Id:"+strconv.Itoa(formato.ActaRecibidoId),
		"UbicacionId",
		"Id",
		"desc",
		"",
		"1",
	)
	if outputError != nil {
		return "", "", outputError
	}
	if len(historicos) == 0 || historicos[0].UbicacionId <= 0 {
		return "", "", nil
	}

	return consultarCentroCostoA11ByID(strconv.Itoa(historicos[0].UbicacionId))
}

func centroCostoMovimientoA11Info(movimiento *models.Movimiento) (codigo, nombre string, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("centroCostoMovimientoA11Info - Unhandled Error!", "500")

	if movimiento == nil {
		return "", "", nil
	}

	var detalle models.FormatoSalidaCostos
	if err := utilsHelper.Unmarshal(movimiento.Detalle, &detalle); err != nil {
		return "", "", err
	}

	switch {
	case strings.TrimSpace(detalle.CentroCostos) != "":
		return consultarCentroCostoA11ByID(strings.TrimSpace(detalle.CentroCostos))
	case detalle.Ubicacion > 0:
		return consultarCentroCostoA11ByID(strconv.Itoa(detalle.Ubicacion))
	default:
		return "", "", nil
	}
}

func consultarCentroCostoA11ByID(id string) (codigo, nombre string, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarCentroCostoA11ByID - Unhandled Error!", "500")

	id = strings.TrimSpace(id)
	if id == "" {
		return "", "", nil
	}

	centrosCostos, outputError := consultarCentroCostosFn("query=Id:" + id)
	if outputError != nil {
		return "", "", outputError
	}
	if len(centrosCostos) == 0 {
		return "", "", nil
	}

	return strings.TrimSpace(centrosCostos[0].Codigo), strings.TrimSpace(centrosCostos[0].Nombre), nil
}

func naturalezaYValorMovimientoA11(movimiento *models.DetalleMovimientoContable) (naturaleza string, valor float64) {
	if movimiento == nil {
		return "", 0
	}

	if movimiento.Debito > 0 && movimiento.Credito > 0 {
		logs.Warn("A11 movimiento contable con debito y credito simultaneo para cuenta %s", codigoCuentaContableA11(movimiento))
		return "D", movimiento.Debito
	}
	if movimiento.Debito > 0 {
		return "D", movimiento.Debito
	}
	if movimiento.Credito > 0 {
		return "C", movimiento.Credito
	}
	return "", 0
}

func codigoCuentaContableA11(movimiento *models.DetalleMovimientoContable) string {
	if movimiento == nil || movimiento.Cuenta == nil {
		return ""
	}
	if strings.TrimSpace(movimiento.Cuenta.Codigo) != "" {
		return strings.TrimSpace(movimiento.Cuenta.Codigo)
	}
	return strings.TrimSpace(movimiento.Cuenta.Id)
}

func terceroMovimientoA11(movimiento *models.DetalleMovimientoContable, proveedorLabel string) (identificacion, nombre string) {
	if movimiento != nil && movimiento.TerceroId != nil {
		return strings.TrimSpace(movimiento.TerceroId.Numero), strings.TrimSpace(movimiento.TerceroId.NombreCompleto)
	}
	return splitTerceroLabel(proveedorLabel)
}
