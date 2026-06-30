package reportesHelper

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tealeg/xlsx"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/catalogoElementosHelper"
	crudActaRecibido "github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

const (
	dateLayout     = "2006-01-02"
	dateTimeLayout = "2006-01-02 15:04:05"
	sheetName      = "EntradasElementos"
	historyWorkers = 8
	salidaWorkers  = 6
	decimalNumFmt  = "#.##0,00"

	estadoEntradaAnuladaReporte = "Entrada anulada"
	estadoSalidaAnuladaReporte  = "Salida anulada"
	rotuloEntradaAnuladaReporte = "ENTRADA ANULADA"
	rotuloSalidaAnuladaReporte  = "SALIDA ANULADA"
)

type entradaReporteData struct {
	Movimiento          *models.Movimiento
	Formato             models.FormatoBaseEntrada
	TransaccionContable *models.InfoTransaccionContable
	Elementos           []*models.DetalleElemento
	CuentasPorSubgrupo  map[int]models.CuentasSubgrupo
	SalidasPorElemento  map[int]*salidaReporteData
	Proveedor           string
	FacturaConsecutivo  string
	FacturaFecha        time.Time
}

type salidaReporteData struct {
	Movimiento          *models.Movimiento
	TransaccionContable *models.InfoTransaccionContable
	FuncionarioAsignado string
	TrasladosAsociados  string
	CuentasPorSubgrupo  map[int]models.CuentasSubgrupo
	CentroCostoNombre   string
	CentroCostoCodigo   string
}

type salidaReporteBaseData struct {
	Movimiento          *models.Movimiento
	TransaccionContable *models.InfoTransaccionContable
	FuncionarioAsignado string
	CuentasPorSubgrupo  map[int]models.CuentasSubgrupo
	CentroCostoNombre   string
	CentroCostoCodigo   string
}

type reporteElementoEntradaRow struct {
	Vigencia                  string
	Periodo                   string
	ElementoNombre            string
	ElementoMarca             string
	ElementoSerie             string
	ElementoPlaca             string
	ElementoCantidad          int
	ElementoValorUnitario     float64
	ElementoValorTotal        float64
	ElementoPorcentajeIvaID   int
	ElementoValorIva          float64
	ElementoValorFinal        float64
	ElementoSubgrupoNombre    string
	EntradaProveedor          string
	EntradaFacturaConsecutivo string
	EntradaFacturaFecha       time.Time
	ElementoVidaUtilCatalogo  float64
	TipoEntrada               string
	EntradaFechaCreacion      time.Time
	EntradaConsecutivo        string
	CuentaDebitoEntrada       string
	CuentaCreditoEntrada      string
	SalidaFechaCreacion       time.Time
	SalidaConsecutivo         string
	SalidaFuncionarioAsignado string
	SalidaCentroCosto         string
	SalidaCodigoCentroCosto   string
	CuentaDebitoSalida        string
	CuentaCreditoSalida       string

	EntradaEstado          string
	EntradaActaRecibidoID  int
	ElementoSubtotal       float64
	ElementoDescuento      float64
	ElementoSubgrupoCodigo string
	ElementoTipoBienNombre string
	SalidaEstado           string
	SalidaFechaCorte       *time.Time
}

var (
	reporteElementosHeaders = []string{
		"Vigencia",
		"Periodo",
		"Nombre / Descripción",
		"Marca",
		"Serie",
		"Placa",
		"Cantidad",
		"Valor unitario",
		"Subtotal",
		"Porcentaje IVA",
		"IVA",
		"Total más IVA",
		"Clase",
		"Proveedor",
		"Consecutivo Factura",
		"Fecha Factura",
		"Vida útil (años)",
		"Tipo de entrada",
		"Fecha entrada",
		"Consecutivo Entrada",
		"Cuenta débito entrada",
		"Cuenta crédito entrada",
		"Fecha salida",
		"Consecutivo salida",
		"Funcionario asignado",
		"Centro de costo",
		"codigo centro de costo",
		"Cuenta débito salida",
		"Cuenta crédito salida",
		"Fecha depreciación",
		"Vida útil",
		"entrada_estado",
		"Tipo de bien",
		"salida_estado",
		"salida_fecha_corte",
	}

	consultarEntradasReporteData             = consultarEntradasReporteDataDefault
	consultarEntradasAnuladasReporteData     = consultarEntradasAnuladasReporteDataDefault
	consultarSalidasAnuladasReporteData      = consultarSalidasAnuladasReporteDataDefault
	consultarCuentaContable                  = cuentasContables.GetCuentaContable
	consultarMovimientoPorConsec             = consultarMovimientoPorConsecutivoDefault
	consultarMovimientoPorID                 = movimientosArka.GetMovimientoById
	consultarTrSalida                        = movimientosArka.GetTrSalida
	consultarCentroCostosFn                  = movimientosArka.GetAllCentroCostos
	consultarElementosActa                   = actaRecibido.GetElementos
	consultarMetadataEntradaFn               = consultarMetadataEntrada
	resolverSalidasPorElementoFn             = resolverSalidasPorElemento
	getCuentasByMovimientoAndSubgrupos       = catalogoElementosHelper.GetCuentasByMovimientoAndSubgrupos
	getNombreTerceroByID                     = terceros.GetNombreTerceroById
	consultarTransaccionContableMovimientoFn = consultarTransaccionContableMovimiento
	consultarMovimientosReporteFn            = movimientosArka.GetAllMovimiento
)

func GenerarReporteElementos(req *models.ReporteFechasRequest) (respuesta *models.ReporteExcelBase64Response, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("GenerarReporteElementos - Unhandled Error!", "500")

	if req == nil {
		return nil, errorCtrl.Error("GenerarReporteElementos - req", "request nil", "400")
	}

	fechaInicial, err := time.Parse(dateLayout, req.FechaInicial)
	if err != nil {
		return nil, errorCtrl.Error("GenerarReporteElementos - time.Parse(fecha_inicial)", err, "400")
	}

	fechaFinal, err := time.Parse(dateLayout, req.FechaFinal)
	if err != nil {
		return nil, errorCtrl.Error("GenerarReporteElementos - time.Parse(fecha_final)", err, "400")
	}

	if fechaFinal.Before(fechaInicial) {
		return nil, errorCtrl.Error("GenerarReporteElementos - rango_fechas", "fecha_final debe ser mayor o igual a fecha_inicial", "400")
	}

	entradas, outputError := consultarEntradasReporteData(fechaInicial, fechaFinal)
	if outputError != nil {
		return nil, outputError
	}

	rows := construirFilasReporteEntradas(entradas)

	entradasAnuladas, outputError := consultarEntradasAnuladasReporteData(fechaInicial, fechaFinal)
	if outputError != nil {
		return nil, outputError
	}
	rows = appendRowsEntradasAnuladasUnicas(rows, construirFilasMovimientosEntradaSinElementos(entradasAnuladas)...)

	salidasAnuladas, outputError := consultarSalidasAnuladasReporteData(fechaInicial, fechaFinal)
	if outputError != nil {
		return nil, outputError
	}
	rows = appendRowsSalidasAnuladasUnicas(rows, construirFilasMovimientosSalidaSinElementos(salidasAnuladas)...)

	archivo := xlsx.NewFile()
	hoja, err := archivo.AddSheet(sheetName)
	if err != nil {
		return nil, errorCtrl.Error("GenerarReporteElementos - archivo.AddSheet", err, "500")
	}

	filaEncabezado := hoja.AddRow()
	for _, header := range reporteElementosHeaders {
		filaEncabezado.AddCell().SetString(header)
	}

	for _, row := range rows {
		addElementoEntradaRow(hoja, row)
	}

	setColumnWidths(hoja)

	var buffer bytes.Buffer
	if err := archivo.Write(&buffer); err != nil {
		return nil, errorCtrl.Error("GenerarReporteElementos - archivo.Write", err, "500")
	}

	respuesta = &models.ReporteExcelBase64Response{
		ArchivoBase64: base64.StdEncoding.EncodeToString(buffer.Bytes()),
		NombreArchivo: fmt.Sprintf("reporte_elementos_entradas_%s_%s.xlsx", fechaInicial.Format("20060102"), fechaFinal.Format("20060102")),
		TipoArchivo:   "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	}

	return respuesta, nil
}

func GetDetalleCuentasEntradaPorConsecutivo(consecutivo string) (respuesta []*models.ReporteDetalleEntradaResponse, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("GetDetalleCuentasEntradaPorConsecutivo - Unhandled Error!", "500")

	movimiento, outputError := consultarMovimientoPorConsec(consecutivo)
	if outputError != nil {
		return nil, outputError
	}

	entrada, outputError := construirEntradaReporteData(movimiento)
	if outputError != nil {
		return nil, outputError
	}
	if entrada == nil {
		return []*models.ReporteDetalleEntradaResponse{}, nil
	}

	rows := construirFilasReporteEntradas([]*entradaReporteData{entrada})
	respuesta = make([]*models.ReporteDetalleEntradaResponse, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		respuesta = append(respuesta, &models.ReporteDetalleEntradaResponse{
			ElementoNombre:            row.ElementoNombre,
			ElementoValorFinal:        row.ElementoValorFinal,
			SalidaFuncionarioAsignado: row.SalidaFuncionarioAsignado,
			CuentaDebitoEntrada:       row.CuentaDebitoEntrada,
			CuentaCreditoEntrada:      row.CuentaCreditoEntrada,
		})
	}

	return respuesta, nil
}

func GetDetalleCuentasSalidaPorConsecutivo(consecutivo string) (respuesta []*models.ReporteDetalleSalidaResponse, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("GetDetalleCuentasSalidaPorConsecutivo - Unhandled Error!", "500")

	movimiento, outputError := consultarMovimientoPorConsec(consecutivo)
	if outputError != nil {
		return nil, outputError
	}

	trSalida, outputError := consultarTrSalida(movimiento.Id)
	if outputError != nil {
		return nil, outputError
	}
	if trSalida == nil || trSalida.Salida == nil {
		return []*models.ReporteDetalleSalidaResponse{}, nil
	}

	elementos, outputError := resolverElementosSalida(trSalida)
	if outputError != nil {
		return nil, outputError
	}

	cuentasPorSubgrupo := make(map[int]models.CuentasSubgrupo)
	if trSalida.Salida.FormatoTipoMovimientoId != nil && trSalida.Salida.FormatoTipoMovimientoId.Id > 0 {
		subgrupos := collectSubgrupoIDs(elementos)
		if len(subgrupos) > 0 {
			outputError = getCuentasByMovimientoAndSubgrupos(trSalida.Salida.FormatoTipoMovimientoId.Id, subgrupos, cuentasPorSubgrupo)
			if outputError != nil {
				return nil, outputError
			}
		}
	}

	transaccion := consultarTransaccionContableSalida(trSalida.Salida)
	funcionario := funcionarioSalidaLabel(trSalida.Salida)

	respuesta = make([]*models.ReporteDetalleSalidaResponse, 0, len(elementos))
	for _, elemento := range elementos {
		if elemento == nil {
			continue
		}

		subgrupoID, _, _ := subgrupoInfo(elemento)
		cuentas := cuentasPorSubgrupo[subgrupoID]
		respuesta = append(respuesta, &models.ReporteDetalleSalidaResponse{
			ElementoNombre:            elemento.Nombre,
			ElementoValorFinal:        elemento.ValorFinal,
			SalidaFuncionarioAsignado: funcionario,
			CuentaDebitoSalida:        resolveCuentaMovimientoLabel(cuentas.CuentaDebitoId, transaccion, true),
			CuentaCreditoSalida:       resolveCuentaMovimientoLabel(cuentas.CuentaCreditoId, transaccion, false),
		})
	}

	return respuesta, nil
}

func consultarEntradasReporteDataDefault(fechaInicial, fechaFinal time.Time) (entradas []*entradaReporteData, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarEntradasReporteDataDefault - Unhandled Error!", "500")

	codigosEntrada, outputError := consultarCodigosEntrada()
	if outputError != nil {
		return nil, outputError
	}
	if len(codigosEntrada) == 0 {
		return []*entradaReporteData{}, nil
	}

	movimientos, outputError := consultarEntradasPorFecha(fechaInicial, fechaFinal, codigosEntrada)
	if outputError != nil {
		return nil, outputError
	}

	entradas = make([]*entradaReporteData, 0, len(movimientos))
	for _, movimiento := range movimientos {
		if movimiento == nil || movimiento.Id <= 0 {
			continue
		}

		entrada, outputError := construirEntradaReporteData(movimiento)
		if outputError != nil {
			return nil, outputError
		}
		if entrada != nil {
			entradas = append(entradas, entrada)
		}
	}

	return entradas, nil
}

func consultarEntradasAnuladasReporteDataDefault(fechaInicial, fechaFinal time.Time) (entradas []*entradaReporteData, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarEntradasAnuladasReporteDataDefault - Unhandled Error!", "500")

	codigosEntrada, outputError := consultarCodigosEntrada()
	if outputError != nil {
		return nil, outputError
	}
	if len(codigosEntrada) == 0 {
		return []*entradaReporteData{}, nil
	}

	movimientos, outputError := consultarMovimientosPorFechaYEstado(fechaInicial, fechaFinal, codigosEntrada, estadoEntradaAnuladaReporte)
	if outputError != nil {
		return nil, outputError
	}

	entradas = make([]*entradaReporteData, 0, len(movimientos))
	for _, movimiento := range movimientos {
		entrada, outputError := construirEntradaReporteSinElementosData(movimiento)
		if outputError != nil {
			return nil, outputError
		}
		if entrada != nil {
			entradas = append(entradas, entrada)
		}
	}

	return entradas, nil
}

func consultarSalidasAnuladasReporteDataDefault(fechaInicial, fechaFinal time.Time) (salidas []*salidaReporteData, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarSalidasAnuladasReporteDataDefault - Unhandled Error!", "500")

	codigosSalida, outputError := consultarCodigosSalida()
	if outputError != nil {
		return nil, outputError
	}
	if len(codigosSalida) == 0 {
		return []*salidaReporteData{}, nil
	}

	movimientos, outputError := consultarMovimientosPorEstado(codigosSalida, estadoSalidaAnuladaReporte)
	if outputError != nil {
		return nil, outputError
	}

	salidas = make([]*salidaReporteData, 0, len(movimientos))
	for _, movimiento := range movimientos {
		if movimiento == nil || movimiento.Id <= 0 {
			continue
		}

		base, outputError := construirSalidaReporteBaseData(movimiento, nil)
		if outputError != nil {
			return nil, outputError
		}
		if base == nil {
			continue
		}

		salida := &salidaReporteData{
			Movimiento:          base.Movimiento,
			TransaccionContable: base.TransaccionContable,
			FuncionarioAsignado: base.FuncionarioAsignado,
			CuentasPorSubgrupo:  base.CuentasPorSubgrupo,
			CentroCostoNombre:   base.CentroCostoNombre,
			CentroCostoCodigo:   base.CentroCostoCodigo,
		}

		if movimiento.MovimientoPadreId != nil && movimiento.MovimientoPadreId.Id > 0 {
			entradaPadre, outputError := consultarMovimientoPorID(movimiento.MovimientoPadreId.Id)
			if outputError != nil {
				return nil, outputError
			}
			if entradaPadre != nil {
				salida.Movimiento.MovimientoPadreId = entradaPadre
			}
		}

		if !movimientoEnRangoPorFechaReporteSalida(salida.Movimiento, fechaInicial, fechaFinal) {
			continue
		}

		salidas = append(salidas, salida)
	}

	return salidas, nil
}

func consultarMovimientoPorConsecutivoDefault(consecutivo string) (movimiento *models.Movimiento, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarMovimientoPorConsecutivoDefault - Unhandled Error!", "500")

	if strings.TrimSpace(consecutivo) == "" {
		return nil, errorCtrl.Error("consultarMovimientoPorConsecutivoDefault - consecutivo", "consecutivo vacío", "400")
	}

	params := url.Values{}
	params.Add("limit", "1")
	params.Add("sortby", "Id")
	params.Add("order", "desc")
	params.Add("query", "Consecutivo:"+strings.TrimSpace(consecutivo))

	movimientos, _, outputError := movimientosArka.GetAllMovimiento(params.Encode())
	if outputError != nil {
		return nil, outputError
	}
	if len(movimientos) == 0 || movimientos[0] == nil {
		return nil, errorCtrl.Error("consultarMovimientoPorConsecutivoDefault - movimientosArka.GetAllMovimiento", "no se encontró movimiento con el consecutivo indicado", "404")
	}

	return movimientos[0], nil
}

func consultarCodigosEntrada() (codigos []string, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarCodigosEntrada - Unhandled Error!", "500")

	formatos, outputError := movimientosArka.GetAllFormatoTipoMovimiento("limit=-1&fields=CodigoAbreviacion")
	if outputError != nil {
		return nil, outputError
	}

	return extraerCodigosEntrada(formatos), nil
}

func extraerCodigosEntrada(formatos []*models.FormatoTipoMovimiento) (codigos []string) {
	seen := make(map[string]struct{})
	for _, formato := range formatos {
		if formato == nil {
			continue
		}
		codigo := formato.CodigoAbreviacion
		if strings.Contains(codigo, "ENT_") && !strings.Contains(codigo, "KDX") {
			if _, ok := seen[codigo]; ok {
				continue
			}
			seen[codigo] = struct{}{}
			codigos = append(codigos, codigo)
		}
	}

	return codigos
}

func consultarCodigosSalida() (codigos []string, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarCodigosSalida - Unhandled Error!", "500")

	formatos, outputError := movimientosArka.GetAllFormatoTipoMovimiento("limit=-1&fields=CodigoAbreviacion")
	if outputError != nil {
		return nil, outputError
	}

	return extraerCodigosSalida(formatos), nil
}

func extraerCodigosSalida(formatos []*models.FormatoTipoMovimiento) (codigos []string) {
	seen := make(map[string]struct{})
	for _, formato := range formatos {
		if formato == nil {
			continue
		}
		codigo := formato.CodigoAbreviacion
		if (codigo == "SAL" || codigo == "SAL_CONS") && !strings.Contains(codigo, "KDX") {
			if _, ok := seen[codigo]; ok {
				continue
			}
			seen[codigo] = struct{}{}
			codigos = append(codigos, codigo)
		}
	}

	return codigos
}

func consultarEntradasPorFecha(fechaInicial, fechaFinal time.Time, codigosEntrada []string) (movimientos []*models.Movimiento, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarEntradasPorFecha - Unhandled Error!", "500")

	fechaInicial = time.Date(fechaInicial.Year(), fechaInicial.Month(), fechaInicial.Day(), 0, 0, 0, 0, time.UTC)
	fechaFinal = time.Date(fechaFinal.Year(), fechaFinal.Month(), fechaFinal.Day(), 23, 59, 59, 0, time.UTC)

	params := url.Values{}
	params.Add("limit", "-1")
	params.Add("sortby", "FechaCreacion")
	params.Add("order", "asc")
	params.Add(
		"query",
		"Activo:true,FormatoTipoMovimientoId__CodigoAbreviacion__in:"+strings.Join(codigosEntrada, "|")+
			",FechaCreacion__gte:"+fechaInicial.Format(time.RFC3339)+
			",FechaCreacion__lte:"+fechaFinal.Format(time.RFC3339),
	)

	movimientos, _, outputError = consultarMovimientosReporteFn(params.Encode())
	if outputError != nil {
		return nil, outputError
	}

	return movimientos, nil
}

func consultarMovimientosPorEstado(codigos []string, estado string) (movimientos []*models.Movimiento, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarMovimientosPorEstado - Unhandled Error!", "500")

	if len(codigos) == 0 {
		return []*models.Movimiento{}, nil
	}

	params := url.Values{}
	params.Add("limit", "-1")
	params.Add("sortby", "FechaCreacion")
	params.Add("order", "asc")
	params.Add(
		"query",
		"EstadoMovimientoId__Nombre:"+estado+
			",FormatoTipoMovimientoId__CodigoAbreviacion__in:"+strings.Join(codigos, "|"),
	)

	movimientos, _, outputError = consultarMovimientosReporteFn(params.Encode())
	if outputError != nil {
		return nil, outputError
	}

	return movimientos, nil
}

func consultarMovimientosPorFechaYEstado(fechaInicial, fechaFinal time.Time, codigos []string, estado string) (movimientos []*models.Movimiento, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarMovimientosPorFechaYEstado - Unhandled Error!", "500")

	if len(codigos) == 0 {
		return []*models.Movimiento{}, nil
	}

	fechaInicial = time.Date(fechaInicial.Year(), fechaInicial.Month(), fechaInicial.Day(), 0, 0, 0, 0, time.UTC)
	fechaFinal = time.Date(fechaFinal.Year(), fechaFinal.Month(), fechaFinal.Day(), 23, 59, 59, 0, time.UTC)

	params := url.Values{}
	params.Add("limit", "-1")
	params.Add("sortby", "FechaCreacion")
	params.Add("order", "asc")
	params.Add(
		"query",
		"EstadoMovimientoId__Nombre:"+estado+
			",FormatoTipoMovimientoId__CodigoAbreviacion__in:"+strings.Join(codigos, "|")+
			",FechaCorte__gte:"+fechaInicial.Format(time.RFC3339)+
			",FechaCorte__lte:"+fechaFinal.Format(time.RFC3339),
	)

	movimientos, _, outputError = consultarMovimientosReporteFn(params.Encode())
	if outputError != nil {
		return nil, outputError
	}

	return movimientos, nil
}

func construirEntradaReporteData(movimiento *models.Movimiento) (entrada *entradaReporteData, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("construirEntradaReporteData - Unhandled Error!", "500")

	if movimiento == nil {
		return nil, nil
	}

	var formato models.FormatoBaseEntrada
	outputError = utilsHelper.Unmarshal(movimiento.Detalle, &formato)
	if outputError != nil {
		return nil, outputError
	}

	elementos, outputError := resolverElementosEntrada(formato)
	if outputError != nil {
		return nil, outputError
	}

	proveedor, facturaConsecutivo, facturaFecha, outputError := consultarMetadataEntradaFn(formato)
	if outputError != nil {
		return nil, outputError
	}

	salidasPorElemento, outputError := resolverSalidasPorElementoFn(elementos)
	if outputError != nil {
		return nil, outputError
	}

	cuentasPorSubgrupo := make(map[int]models.CuentasSubgrupo)
	if movimiento.FormatoTipoMovimientoId != nil && movimiento.FormatoTipoMovimientoId.Id > 0 {
		subgrupos := collectSubgrupoIDs(elementos)
		if len(subgrupos) > 0 {
			outputError = getCuentasByMovimientoAndSubgrupos(movimiento.FormatoTipoMovimientoId.Id, subgrupos, cuentasPorSubgrupo)
			if outputError != nil {
				return nil, outputError
			}
		}
	}

	entrada = &entradaReporteData{
		Movimiento:          movimiento,
		Formato:             formato,
		TransaccionContable: consultarTransaccionContableEntrada(movimiento),
		Elementos:           elementos,
		CuentasPorSubgrupo:  cuentasPorSubgrupo,
		SalidasPorElemento:  salidasPorElemento,
		Proveedor:           proveedor,
		FacturaConsecutivo:  facturaConsecutivo,
		FacturaFecha:        facturaFecha,
	}

	return entrada, nil
}

func construirEntradaReporteSinElementosData(movimiento *models.Movimiento) (entrada *entradaReporteData, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("construirEntradaReporteSinElementosData - Unhandled Error!", "500")

	if movimiento == nil {
		return nil, nil
	}

	var formato models.FormatoBaseEntrada
	if strings.TrimSpace(movimiento.Detalle) != "" {
		outputError = utilsHelper.Unmarshal(movimiento.Detalle, &formato)
		if outputError != nil {
			return nil, outputError
		}
	}

	proveedor, facturaConsecutivo, facturaFecha, outputError := consultarMetadataEntradaFn(formato)
	if outputError != nil {
		return nil, outputError
	}

	return &entradaReporteData{
		Movimiento:          movimiento,
		Formato:             formato,
		TransaccionContable: consultarTransaccionContableEntrada(movimiento),
		Elementos:           []*models.DetalleElemento{},
		CuentasPorSubgrupo:  map[int]models.CuentasSubgrupo{},
		SalidasPorElemento:  map[int]*salidaReporteData{},
		Proveedor:           proveedor,
		FacturaConsecutivo:  facturaConsecutivo,
		FacturaFecha:        facturaFecha,
	}, nil
}

func consultarMetadataEntrada(formato models.FormatoBaseEntrada) (proveedor, facturaConsecutivo string, facturaFecha time.Time, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarMetadataEntrada - Unhandled Error!", "500")

	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		firstErr map[string]interface{}
	)

	setErr := func(err map[string]interface{}) {
		if err == nil {
			return
		}
		mu.Lock()
		defer mu.Unlock()
		if firstErr == nil {
			firstErr = err
		}
	}

	if formato.ActaRecibidoId > 0 {
		wg.Add(1)
		go func(actaRecibidoID int) {
			defer wg.Done()
			proveedor_, err := consultarProveedorActa(actaRecibidoID)
			if err != nil {
				setErr(err)
				return
			}
			mu.Lock()
			proveedor = proveedor_
			mu.Unlock()
		}(formato.ActaRecibidoId)
	}

	if formato.Factura > 0 {
		wg.Add(1)
		go func(facturaID int) {
			defer wg.Done()
			consecutivo, fecha, err := consultarFacturaSoporte(facturaID)
			if err != nil {
				setErr(err)
				return
			}
			mu.Lock()
			facturaConsecutivo = consecutivo
			facturaFecha = fecha
			mu.Unlock()
		}(formato.Factura)
	}

	wg.Wait()
	if firstErr != nil {
		return "", "", time.Time{}, firstErr
	}

	return proveedor, facturaConsecutivo, facturaFecha, nil
}

func consultarProveedorActa(actaRecibidoID int) (proveedor string, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarProveedorActa - Unhandled Error!", "500")

	if actaRecibidoID <= 0 {
		return "", nil
	}

	historicos, outputError := consultarHistoricosActaReporteFn("ActaRecibidoId__Id:"+strconv.Itoa(actaRecibidoID), "", "Id", "desc", "", "1")
	if outputError != nil || len(historicos) == 0 || historicos[0].ProveedorId <= 0 {
		return "", outputError
	}

	tercero, outputError := getNombreTerceroByID(historicos[0].ProveedorId)
	if outputError != nil {
		if esTerceroNoEncontradoReporte(outputError) {
			return strconv.Itoa(historicos[0].ProveedorId), nil
		}
		return "", outputError
	}

	return identificacionTerceroLabel(tercero), nil
}

func esTerceroNoEncontradoReporte(err map[string]interface{}) bool {
	if err == nil {
		return false
	}
	return strings.Contains(fmt.Sprint(err["err"]), "http 404:")
}

func consultarFacturaSoporte(facturaID int) (consecutivo string, fecha time.Time, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarFacturaSoporte - Unhandled Error!", "500")

	if facturaID <= 0 {
		return "", time.Time{}, nil
	}

	var soporte models.SoporteActa
	outputError = crudActaRecibido.GetSoporteById(facturaID, &soporte)
	if outputError != nil {
		return "", time.Time{}, outputError
	}

	return soporte.Consecutivo, soporte.FechaSoporte, nil
}

func resolverElementosEntrada(formato models.FormatoBaseEntrada) (elementos []*models.DetalleElemento, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("resolverElementosEntrada - Unhandled Error!", "500")

	if formato.ActaRecibidoId > 0 {
		return consultarElementosActa(formato.ActaRecibidoId, nil)
	}

	if len(formato.Elementos) == 0 {
		return []*models.DetalleElemento{}, nil
	}

	idsElementosMovimiento := make([]int, 0, len(formato.Elementos))
	for _, elemento := range formato.Elementos {
		if elemento.Id > 0 {
			idsElementosMovimiento = append(idsElementosMovimiento, elemento.Id)
		}
	}
	if len(idsElementosMovimiento) == 0 {
		return []*models.DetalleElemento{}, nil
	}

	params := url.Values{}
	params.Add("limit", "-1")
	params.Add("fields", "Id,ElementoActaId")
	params.Add("query", "Id__in:"+utilsHelper.ArrayToString(idsElementosMovimiento, "|"))
	elementosMovimiento, outputError := movimientosArka.GetAllElementosMovimiento(params.Encode())
	if outputError != nil {
		return nil, outputError
	}

	elementosActaIds := make([]int, 0, len(idsElementosMovimiento))
	elementoActaPorMovimiento := make(map[int]int)
	for _, elementoMovimiento := range elementosMovimiento {
		if elementoMovimiento == nil || elementoMovimiento.ElementoActaId == nil || *elementoMovimiento.ElementoActaId <= 0 {
			continue
		}
		elementoActaPorMovimiento[elementoMovimiento.Id] = *elementoMovimiento.ElementoActaId
	}

	for _, idMovimiento := range idsElementosMovimiento {
		if idActa, ok := elementoActaPorMovimiento[idMovimiento]; ok {
			elementosActaIds = append(elementosActaIds, idActa)
		}
	}
	if len(elementosActaIds) == 0 {
		return []*models.DetalleElemento{}, nil
	}

	elementos, outputError = consultarElementosActa(0, elementosActaIds)
	if outputError != nil {
		return nil, outputError
	}

	return ordenarElementosPorIds(elementosActaIds, elementos), nil
}

func resolverElementosSalida(trSalida *models.TrSalida) (elementos []*models.DetalleElemento, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("resolverElementosSalida - Unhandled Error!", "500")

	if trSalida == nil || len(trSalida.Elementos) == 0 {
		return []*models.DetalleElemento{}, nil
	}

	ids := make([]int, 0, len(trSalida.Elementos))
	for _, elemento := range trSalida.Elementos {
		if elemento == nil || elemento.ElementoActaId == nil || *elemento.ElementoActaId <= 0 {
			continue
		}
		ids = append(ids, *elemento.ElementoActaId)
	}
	if len(ids) == 0 {
		return []*models.DetalleElemento{}, nil
	}

	elementos, outputError = consultarElementosActa(0, ids)
	if outputError != nil {
		return nil, outputError
	}

	return ordenarElementosPorIds(ids, elementos), nil
}

func consultarTransaccionContableEntrada(movimiento *models.Movimiento) *models.InfoTransaccionContable {
	return consultarTransaccionContableMovimientoFn(movimiento, "Entrada Aprobada", "Entrada Con Salida")
}

func consultarTransaccionContableSalida(movimiento *models.Movimiento) *models.InfoTransaccionContable {
	return consultarTransaccionContableMovimientoFn(movimiento, "Salida Aprobada")
}

func consultarTransaccionContableMovimiento(movimiento *models.Movimiento, estadosPermitidos ...string) *models.InfoTransaccionContable {
	if movimiento == nil || movimiento.EstadoMovimientoId == nil || movimiento.ConsecutivoId == nil || *movimiento.ConsecutivoId <= 0 {
		return nil
	}

	if !movimientoTieneEstado(movimiento, estadosPermitidos...) {
		return nil
	}

	transaccion, outputError := asientoContable.GetFullDetalleContable(*movimiento.ConsecutivoId)
	if outputError != nil {
		return nil
	}

	return &transaccion
}

func resolverSalidasPorElemento(elementos []*models.DetalleElemento) (salidasPorElemento map[int]*salidaReporteData, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("resolverSalidasPorElemento - Unhandled Error!", "500")

	salidasPorElemento = make(map[int]*salidaReporteData)
	if len(elementos) == 0 {
		return salidasPorElemento, nil
	}

	subgrupoPorElemento := make(map[int]int)
	elementosActaIDs := make([]int, 0, len(elementos))
	for _, elemento := range elementos {
		if elemento == nil || elemento.Id <= 0 {
			continue
		}
		elementosActaIDs = append(elementosActaIDs, elemento.Id)
		subgrupoID, _, _ := subgrupoInfo(elemento)
		subgrupoPorElemento[elemento.Id] = subgrupoID
	}
	if len(elementosActaIDs) == 0 {
		return salidasPorElemento, nil
	}

	params := url.Values{}
	params.Add("limit", "-1")
	params.Add("sortby", "Id")
	params.Add("order", "desc")
	params.Add("fields", "Id,ElementoActaId")
	params.Add("query", "ElementoActaId__in:"+utilsHelper.ArrayToString(elementosActaIDs, "|"))
	elementosMovimiento, outputError := movimientosArka.GetAllElementosMovimiento(params.Encode())
	if outputError != nil {
		return nil, outputError
	}

	ultimoMovimientoPorElemento := make(map[int]*models.ElementosMovimiento)
	for _, elementoMovimiento := range elementosMovimiento {
		if elementoMovimiento == nil || elementoMovimiento.ElementoActaId == nil || *elementoMovimiento.ElementoActaId <= 0 {
			continue
		}
		elementoActaID := *elementoMovimiento.ElementoActaId
		if _, ok := ultimoMovimientoPorElemento[elementoActaID]; ok {
			continue
		}
		ultimoMovimientoPorElemento[elementoActaID] = elementoMovimiento
	}

	historialesPorElemento, outputError := consultarHistorialesElementos(ultimoMovimientoPorElemento)
	if outputError != nil {
		return nil, outputError
	}

	salidasDescriptor := make(map[int]*salidaReporteBaseData)
	subgruposPorSalida := make(map[int]map[int]struct{})
	for _, elementoActaID := range elementosActaIDs {
		historial := historialesPorElemento[elementoActaID]
		if historial == nil || historial.Salida == nil {
			continue
		}

		salidaID := historial.Salida.Id
		if salidaID <= 0 {
			continue
		}

		if _, ok := salidasDescriptor[salidaID]; !ok {
			salidasDescriptor[salidaID] = &salidaReporteBaseData{
				Movimiento: historial.Salida,
			}
		}
		if _, ok := subgruposPorSalida[salidaID]; !ok {
			subgruposPorSalida[salidaID] = make(map[int]struct{})
		}
		if subgrupoID := subgrupoPorElemento[elementoActaID]; subgrupoID > 0 {
			subgruposPorSalida[salidaID][subgrupoID] = struct{}{}
		}
	}

	salidasBasePorID, outputError := construirSalidasReporteBaseData(salidasDescriptor, subgruposPorSalida)
	if outputError != nil {
		return nil, outputError
	}

	for _, elementoActaID := range elementosActaIDs {
		historial := historialesPorElemento[elementoActaID]
		if historial == nil || historial.Salida == nil {
			continue
		}

		base := salidasBasePorID[historial.Salida.Id]
		if base == nil {
			continue
		}

		salidasPorElemento[elementoActaID] = &salidaReporteData{
			Movimiento:          base.Movimiento,
			TransaccionContable: base.TransaccionContable,
			FuncionarioAsignado: base.FuncionarioAsignado,
			TrasladosAsociados:  trasladosAsociadosLabel(historial.Traslados),
			CuentasPorSubgrupo:  base.CuentasPorSubgrupo,
			CentroCostoNombre:   base.CentroCostoNombre,
			CentroCostoCodigo:   base.CentroCostoCodigo,
		}
	}

	return salidasPorElemento, nil
}

func consultarHistorialesElementos(ultimoMovimientoPorElemento map[int]*models.ElementosMovimiento) (historialesPorElemento map[int]*models.Historial, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarHistorialesElementos - Unhandled Error!", "500")

	historialesPorElemento = make(map[int]*models.Historial, len(ultimoMovimientoPorElemento))
	if len(ultimoMovimientoPorElemento) == 0 {
		return historialesPorElemento, nil
	}

	type historialResult struct {
		elementoActaID int
		historial      *models.Historial
		outputError    map[string]interface{}
	}

	results := make(chan historialResult, len(ultimoMovimientoPorElemento))
	sem := make(chan struct{}, historyWorkers)
	var wg sync.WaitGroup

	for elementoActaID, elementoMovimiento := range ultimoMovimientoPorElemento {
		if elementoMovimiento == nil || elementoMovimiento.Id <= 0 {
			continue
		}

		wg.Add(1)
		go func(elementoActaID int, elementoMovimientoID int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			historial, outputError := movimientosArka.GetHistorialElemento(elementoMovimientoID, true)
			results <- historialResult{
				elementoActaID: elementoActaID,
				historial:      historial,
				outputError:    outputError,
			}
		}(elementoActaID, elementoMovimiento.Id)
	}

	wg.Wait()
	close(results)

	for result := range results {
		if result.outputError != nil {
			if outputError == nil {
				outputError = result.outputError
			}
			continue
		}
		historialesPorElemento[result.elementoActaID] = result.historial
	}

	if outputError != nil {
		return nil, outputError
	}

	return historialesPorElemento, nil
}

func construirSalidasReporteBaseData(
	salidasDescriptor map[int]*salidaReporteBaseData,
	subgruposPorSalida map[int]map[int]struct{},
) (salidasBasePorID map[int]*salidaReporteBaseData, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("construirSalidasReporteBaseData - Unhandled Error!", "500")

	salidasBasePorID = make(map[int]*salidaReporteBaseData, len(salidasDescriptor))
	if len(salidasDescriptor) == 0 {
		return salidasBasePorID, nil
	}

	type salidaResult struct {
		salidaID    int
		salida      *salidaReporteBaseData
		outputError map[string]interface{}
	}

	results := make(chan salidaResult, len(salidasDescriptor))
	sem := make(chan struct{}, salidaWorkers)
	var wg sync.WaitGroup

	for salidaID, descriptor := range salidasDescriptor {
		if descriptor == nil || descriptor.Movimiento == nil {
			continue
		}

		subgrupos := subgrupoSetToSlice(subgruposPorSalida[salidaID])
		wg.Add(1)
		go func(salidaID int, movimiento *models.Movimiento, subgrupos []int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			salida, outputError := construirSalidaReporteBaseData(movimiento, subgrupos)
			results <- salidaResult{
				salidaID:    salidaID,
				salida:      salida,
				outputError: outputError,
			}
		}(salidaID, descriptor.Movimiento, subgrupos)
	}

	wg.Wait()
	close(results)

	for result := range results {
		if result.outputError != nil {
			if outputError == nil {
				outputError = result.outputError
			}
			continue
		}
		if result.salida != nil {
			salidasBasePorID[result.salidaID] = result.salida
		}
	}

	if outputError != nil {
		return nil, outputError
	}

	return salidasBasePorID, nil
}

func construirSalidaReporteBaseData(movimiento *models.Movimiento, subgrupos []int) (salida *salidaReporteBaseData, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("construirSalidaReporteBaseData - Unhandled Error!", "500")

	if movimiento == nil {
		return nil, nil
	}

	centroCostoNombre, centroCostoCodigo := salidaUbicacionInfo(movimiento)

	cuentasPorSubgrupo := make(map[int]models.CuentasSubgrupo)
	if movimiento.FormatoTipoMovimientoId != nil && movimiento.FormatoTipoMovimientoId.Id > 0 && len(subgrupos) > 0 {
		outputError = getCuentasByMovimientoAndSubgrupos(movimiento.FormatoTipoMovimientoId.Id, subgrupos, cuentasPorSubgrupo)
		if outputError != nil {
			return nil, outputError
		}
	}

	salida = &salidaReporteBaseData{
		Movimiento:          movimiento,
		TransaccionContable: consultarTransaccionContableSalida(movimiento),
		FuncionarioAsignado: funcionarioSalidaLabel(movimiento),
		CuentasPorSubgrupo:  cuentasPorSubgrupo,
		CentroCostoNombre:   centroCostoNombre,
		CentroCostoCodigo:   centroCostoCodigo,
	}

	return salida, nil
}

func collectSubgrupoIDs(elementos []*models.DetalleElemento) []int {
	subgrupos := make([]int, 0, len(elementos))
	seen := make(map[int]struct{})

	for _, elemento := range elementos {
		if elemento == nil || elemento.SubgrupoCatalogoId == nil || elemento.SubgrupoCatalogoId.SubgrupoId == nil {
			continue
		}
		subgrupoID := elemento.SubgrupoCatalogoId.SubgrupoId.Id
		if subgrupoID <= 0 {
			continue
		}
		if _, ok := seen[subgrupoID]; ok {
			continue
		}
		seen[subgrupoID] = struct{}{}
		subgrupos = append(subgrupos, subgrupoID)
	}

	return subgrupos
}

func ordenarElementosPorIds(ids []int, elementos []*models.DetalleElemento) []*models.DetalleElemento {
	elementosPorID := make(map[int]*models.DetalleElemento)
	for _, elemento := range elementos {
		if elemento != nil {
			elementosPorID[elemento.Id] = elemento
		}
	}

	ordenados := make([]*models.DetalleElemento, 0, len(ids))
	for _, id := range ids {
		if elemento, ok := elementosPorID[id]; ok {
			ordenados = append(ordenados, elemento)
		}
	}

	return ordenados
}

func subgrupoSetToSlice(subgrupoSet map[int]struct{}) []int {
	if len(subgrupoSet) == 0 {
		return nil
	}

	subgrupos := make([]int, 0, len(subgrupoSet))
	for subgrupoID := range subgrupoSet {
		if subgrupoID > 0 {
			subgrupos = append(subgrupos, subgrupoID)
		}
	}

	return subgrupos
}

func construirFilasReporteEntradas(entradas []*entradaReporteData) []*reporteElementoEntradaRow {
	rows := make([]*reporteElementoEntradaRow, 0)

	for _, entrada := range entradas {
		if entrada == nil || entrada.Movimiento == nil {
			continue
		}

		if movimientoAnuladoReporte(entrada.Movimiento, estadoEntradaAnuladaReporte) {
			rows = append(rows, construirFilaMovimientoEntradaSinElementos(entrada))
			continue
		}

		for _, elemento := range entrada.Elementos {
			if elemento == nil {
				continue
			}

			subgrupoID, subgrupoCodigo, subgrupoNombre := subgrupoInfo(elemento)
			_, tipoBienNombre := tipoBienInfo(elemento)

			movimientoCuenta := entrada.CuentasPorSubgrupo[subgrupoID]
			salida := entrada.SalidasPorElemento[elemento.Id]
			if salida != nil && movimientoAnuladoReporte(salida.Movimiento, estadoSalidaAnuladaReporte) {
				continue
			}
			salidaCuenta := models.CuentasSubgrupo{}
			if salida != nil {
				salidaCuenta = salida.CuentasPorSubgrupo[subgrupoID]
			}
			row := &reporteElementoEntradaRow{
				Vigencia:                  strconv.Itoa(entrada.Movimiento.FechaCreacion.Year()),
				Periodo:                   strconv.Itoa(int(entrada.Movimiento.FechaCreacion.Month())),
				EntradaConsecutivo:        stringPtrValue(entrada.Movimiento.Consecutivo),
				EntradaEstado:             estadoMovimientoNombre(entrada.Movimiento),
				EntradaFechaCreacion:      entrada.Movimiento.FechaCreacion,
				EntradaActaRecibidoID:     entrada.Formato.ActaRecibidoId,
				EntradaProveedor:          entrada.Proveedor,
				EntradaFacturaConsecutivo: entrada.FacturaConsecutivo,
				EntradaFacturaFecha:       entrada.FacturaFecha,
				TipoEntrada:               tipoEntradaNombre(entrada.Movimiento),
				ElementoNombre:            elemento.Nombre,
				ElementoCantidad:          elemento.Cantidad,
				ElementoMarca:             elemento.Marca,
				ElementoSerie:             elemento.Serie,
				ElementoValorUnitario:     elemento.ValorUnitario,
				ElementoSubtotal:          elemento.Subtotal,
				ElementoDescuento:         elemento.Descuento,
				ElementoValorTotal:        elemento.ValorTotal,
				ElementoPorcentajeIvaID:   elemento.PorcentajeIvaId,
				ElementoValorIva:          elemento.ValorIva,
				ElementoValorFinal:        elemento.ValorFinal,
				ElementoSubgrupoCodigo:    subgrupoCodigo,
				ElementoSubgrupoNombre:    subgrupoNombre,
				ElementoTipoBienNombre:    tipoBienNombre,
				ElementoVidaUtilCatalogo:  vidaUtilCatalogo(elemento),
				ElementoPlaca:             elemento.Placa,
				CuentaDebitoEntrada:       resolveCuentaMovimientoLabel(movimientoCuenta.CuentaDebitoId, entrada.TransaccionContable, true),
				CuentaCreditoEntrada:      resolveCuentaMovimientoLabel(movimientoCuenta.CuentaCreditoId, entrada.TransaccionContable, false),
				SalidaConsecutivo:         movimientoConsecutivo(salida),
				SalidaEstado:              movimientoEstado(salida),
				SalidaFechaCreacion:       movimientoFechaCreacion(salida),
				SalidaFechaCorte:          movimientoFechaCorte(salida),
				SalidaFuncionarioAsignado: salidaFuncionario(salida),
				SalidaCentroCosto:         salidaCentroCosto(salida),
				SalidaCodigoCentroCosto:   salidaCodigoCentroCosto(salida),
				CuentaDebitoSalida:        resolveCuentaMovimientoLabel(salidaCuenta.CuentaDebitoId, transaccionContableSalida(salida), true),
				CuentaCreditoSalida:       resolveCuentaMovimientoLabel(salidaCuenta.CuentaCreditoId, transaccionContableSalida(salida), false),
			}
			rows = append(rows, row)
		}
	}

	return rows
}

func construirFilasMovimientosEntradaSinElementos(entradas []*entradaReporteData) []*reporteElementoEntradaRow {
	rows := make([]*reporteElementoEntradaRow, 0, len(entradas))
	for _, entrada := range entradas {
		if entrada == nil || !movimientoAnuladoReporte(entrada.Movimiento, estadoEntradaAnuladaReporte) {
			continue
		}
		if row := construirFilaMovimientoEntradaSinElementos(entrada); row != nil {
			rows = append(rows, row)
		}
	}
	return rows
}

func construirFilasMovimientosSalidaSinElementos(salidas []*salidaReporteData) []*reporteElementoEntradaRow {
	rows := make([]*reporteElementoEntradaRow, 0, len(salidas))
	for _, salida := range salidas {
		if row := construirFilaMovimientoSalidaSinElementos(salida); row != nil {
			rows = append(rows, row)
		}
	}
	return rows
}

func construirFilaMovimientoEntradaSinElementos(entrada *entradaReporteData) *reporteElementoEntradaRow {
	if entrada == nil || entrada.Movimiento == nil {
		return nil
	}

	row := &reporteElementoEntradaRow{
		Vigencia:                  strconv.Itoa(entrada.Movimiento.FechaCreacion.Year()),
		Periodo:                   strconv.Itoa(int(entrada.Movimiento.FechaCreacion.Month())),
		EntradaConsecutivo:        stringPtrValue(entrada.Movimiento.Consecutivo),
		EntradaEstado:             estadoMovimientoNombre(entrada.Movimiento),
		EntradaFechaCreacion:      entrada.Movimiento.FechaCreacion,
		EntradaActaRecibidoID:     entrada.Formato.ActaRecibidoId,
		EntradaProveedor:          entrada.Proveedor,
		EntradaFacturaConsecutivo: entrada.FacturaConsecutivo,
		EntradaFacturaFecha:       entrada.FacturaFecha,
		TipoEntrada:               tipoEntradaNombre(entrada.Movimiento),
	}
	aplicarRotuloEntradaAnulada(row)
	return row
}

func construirFilaMovimientoSalidaSinElementos(salida *salidaReporteData) *reporteElementoEntradaRow {
	if salida == nil || salida.Movimiento == nil {
		return nil
	}

	row := &reporteElementoEntradaRow{
		SalidaConsecutivo:         movimientoConsecutivo(salida),
		SalidaEstado:              movimientoEstado(salida),
		SalidaFechaCreacion:       movimientoFechaCreacion(salida),
		SalidaFechaCorte:          movimientoFechaCorte(salida),
		SalidaFuncionarioAsignado: salidaFuncionario(salida),
		SalidaCentroCosto:         salidaCentroCosto(salida),
		SalidaCodigoCentroCosto:   salidaCodigoCentroCosto(salida),
	}

	if padre := salida.Movimiento.MovimientoPadreId; padre != nil {
		var formato models.FormatoBaseEntrada
		if strings.TrimSpace(padre.Detalle) != "" {
			_ = utilsHelper.Unmarshal(padre.Detalle, &formato)
		}
		row.Vigencia = strconv.Itoa(padre.FechaCreacion.Year())
		row.Periodo = strconv.Itoa(int(padre.FechaCreacion.Month()))
		row.EntradaConsecutivo = stringPtrValue(padre.Consecutivo)
		row.EntradaEstado = estadoMovimientoNombre(padre)
		row.EntradaFechaCreacion = padre.FechaCreacion
		row.EntradaActaRecibidoID = formato.ActaRecibidoId
		row.TipoEntrada = tipoEntradaNombre(padre)
	}

	aplicarRotuloSalidaAnulada(row)
	return row
}

func aplicarRotuloEntradaAnulada(row *reporteElementoEntradaRow) {
	if row == nil {
		return
	}
	aplicarValoresMonetariosCero(row)
	row.ElementoNombre = rotuloEntradaAnuladaReporte
	row.ElementoMarca = rotuloEntradaAnuladaReporte
	row.ElementoSerie = rotuloEntradaAnuladaReporte
	row.ElementoPlaca = rotuloEntradaAnuladaReporte
	row.ElementoSubgrupoNombre = rotuloEntradaAnuladaReporte
	row.ElementoSubgrupoCodigo = rotuloEntradaAnuladaReporte
	row.ElementoTipoBienNombre = rotuloEntradaAnuladaReporte
	row.EntradaProveedor = rotuloEntradaAnuladaReporte
	row.EntradaFacturaConsecutivo = rotuloEntradaAnuladaReporte
	row.TipoEntrada = rotuloEntradaAnuladaReporte
	row.CuentaDebitoEntrada = rotuloEntradaAnuladaReporte
	row.CuentaCreditoEntrada = rotuloEntradaAnuladaReporte
	row.SalidaConsecutivo = rotuloEntradaAnuladaReporte
	row.SalidaFuncionarioAsignado = rotuloEntradaAnuladaReporte
	row.SalidaCentroCosto = rotuloEntradaAnuladaReporte
	row.SalidaCodigoCentroCosto = rotuloEntradaAnuladaReporte
	row.CuentaDebitoSalida = rotuloEntradaAnuladaReporte
	row.CuentaCreditoSalida = rotuloEntradaAnuladaReporte
}

func aplicarRotuloSalidaAnulada(row *reporteElementoEntradaRow) {
	if row == nil {
		return
	}
	aplicarValoresMonetariosCero(row)
	row.ElementoNombre = rotuloSalidaAnuladaReporte
	row.ElementoMarca = rotuloSalidaAnuladaReporte
	row.ElementoSerie = rotuloSalidaAnuladaReporte
	row.ElementoPlaca = rotuloSalidaAnuladaReporte
	row.ElementoSubgrupoNombre = rotuloSalidaAnuladaReporte
	row.ElementoSubgrupoCodigo = rotuloSalidaAnuladaReporte
	row.ElementoTipoBienNombre = rotuloSalidaAnuladaReporte
	row.SalidaFuncionarioAsignado = rotuloSalidaAnuladaReporte
	row.SalidaCentroCosto = rotuloSalidaAnuladaReporte
	row.SalidaCodigoCentroCosto = rotuloSalidaAnuladaReporte
	row.CuentaDebitoSalida = rotuloSalidaAnuladaReporte
	row.CuentaCreditoSalida = rotuloSalidaAnuladaReporte
}

func aplicarValoresMonetariosCero(row *reporteElementoEntradaRow) {
	if row == nil {
		return
	}
	row.ElementoCantidad = 0
	row.ElementoValorUnitario = 0
	row.ElementoSubtotal = 0
	row.ElementoDescuento = 0
	row.ElementoValorTotal = 0
	row.ElementoPorcentajeIvaID = 0
	row.ElementoValorIva = 0
	row.ElementoValorFinal = 0
	row.ElementoVidaUtilCatalogo = 0
}

func appendRowsEntradasAnuladasUnicas(rows []*reporteElementoEntradaRow, extras ...*reporteElementoEntradaRow) []*reporteElementoEntradaRow {
	seen := make(map[string]struct{})
	for _, row := range rows {
		if row == nil || row.EntradaEstado != estadoEntradaAnuladaReporte {
			continue
		}
		seen[row.EntradaConsecutivo] = struct{}{}
	}
	for _, row := range extras {
		if row == nil {
			continue
		}
		if _, ok := seen[row.EntradaConsecutivo]; ok {
			continue
		}
		seen[row.EntradaConsecutivo] = struct{}{}
		rows = append(rows, row)
	}
	return rows
}

func appendRowsSalidasAnuladasUnicas(rows []*reporteElementoEntradaRow, extras ...*reporteElementoEntradaRow) []*reporteElementoEntradaRow {
	seen := make(map[string]struct{})
	for _, row := range rows {
		if row == nil || row.SalidaEstado != estadoSalidaAnuladaReporte {
			continue
		}
		seen[row.SalidaConsecutivo] = struct{}{}
	}
	for _, row := range extras {
		if row == nil {
			continue
		}
		if _, ok := seen[row.SalidaConsecutivo]; ok {
			continue
		}
		seen[row.SalidaConsecutivo] = struct{}{}
		rows = append(rows, row)
	}
	return rows
}

func movimientoAnuladoReporte(movimiento *models.Movimiento, estado string) bool {
	return movimiento != nil &&
		movimiento.EstadoMovimientoId != nil &&
		movimiento.EstadoMovimientoId.Nombre == estado
}

func addElementoEntradaRow(hoja *xlsx.Sheet, rowData *reporteElementoEntradaRow) {
	if rowData == nil {
		return
	}

	row := hoja.AddRow()

	addStringCell(row, rowData.Vigencia)                              // Vigencia
	addStringCell(row, rowData.Periodo)                               // Periodo
	addStringCell(row, rowData.ElementoNombre)                        // Nombre / Descripción
	addStringCell(row, rowData.ElementoMarca)                         // Marca
	addStringCell(row, rowData.ElementoSerie)                         // Serie
	addStringCell(row, rowData.ElementoPlaca)                         // Placa
	addStringCell(row, strconv.Itoa(rowData.ElementoCantidad))        // Cantidad
	addDecimalCell(row, rowData.ElementoValorUnitario)                // Valor unitario
	addDecimalCell(row, rowData.ElementoSubtotal)                     // Subtotal
	addStringCell(row, strconv.Itoa(rowData.ElementoPorcentajeIvaID)) // Porcentaje IVA
	addDecimalCell(row, rowData.ElementoValorIva)                     // IVA
	addDecimalCell(row, rowData.ElementoValorFinal)                   // Total más IVA
	addStringCell(row, rowData.ElementoSubgrupoNombre)                // Clase
	addStringCell(row, rowData.EntradaProveedor)                      // Proveedor
	addStringCell(row, rowData.EntradaFacturaConsecutivo)             // Consecutivo Factura
	addStringCell(row, timeCell(rowData.EntradaFacturaFecha))         // Fecha Factura
	addDecimalCell(row, rowData.ElementoVidaUtilCatalogo)             // Vida útil (años)
	addStringCell(row, rowData.TipoEntrada)                           // Tipo de entrada
	addStringCell(row, timeCell(rowData.EntradaFechaCreacion))        // Fecha entrada
	addStringCell(row, rowData.EntradaConsecutivo)                    // Consecutivo Entrada
	addStringCell(row, rowData.CuentaDebitoEntrada)                   // Cuenta débito entrada
	addStringCell(row, rowData.CuentaCreditoEntrada)                  // Cuenta crédito entrada
	addStringCell(row, timeCell(rowData.SalidaFechaCreacion))         // Fecha salida
	addStringCell(row, rowData.SalidaConsecutivo)                     // Consecutivo salida
	addStringCell(row, rowData.SalidaFuncionarioAsignado)             // Funcionario asignado
	addStringCell(row, rowData.SalidaCentroCosto)                     // Centro de costo
	addStringCell(row, rowData.SalidaCodigoCentroCosto)               // codigo centro de costo
	addStringCell(row, rowData.CuentaDebitoSalida)                    // Cuenta débito salida
	addStringCell(row, rowData.CuentaCreditoSalida)                   // Cuenta crédito salida
	addStringCell(row, "")                                            // Fecha depreciación
	addStringCell(row, "")                                            // Vida útil

	addStringCell(row, rowData.EntradaEstado)
	addStringCell(row, rowData.ElementoTipoBienNombre)
	addStringCell(row, rowData.SalidaEstado)
	addStringCell(row, timePtrCell(rowData.SalidaFechaCorte))
}

func addStringCell(row *xlsx.Row, value string) {
	row.AddCell().SetString(value)
}

func addDecimalCell(row *xlsx.Row, value float64) {
	cell := row.AddCell()
	cell.SetFloatWithFormat(roundToTwoDecimals(value), decimalNumFmt)
}

func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}

func resolveCuentaMovimientoLabel(cuentaID string, transaccion *models.InfoTransaccionContable, debito bool) string {
	if cuentaID == "" {
		return ""
	}

	if transaccion != nil {
		for _, movimiento := range transaccion.Movimientos {
			if movimiento == nil || movimiento.Cuenta == nil || movimiento.Cuenta.Id != cuentaID {
				continue
			}
			if debito && movimiento.Debito > 0 {
				return detalleCuentaLabel(movimiento.Cuenta)
			}
			if !debito && movimiento.Credito > 0 {
				return detalleCuentaLabel(movimiento.Cuenta)
			}
		}

		for _, movimiento := range transaccion.Movimientos {
			if movimiento != nil && movimiento.Cuenta != nil && movimiento.Cuenta.Id == cuentaID {
				return detalleCuentaLabel(movimiento.Cuenta)
			}
		}
	}

	if cuenta, outputError := consultarCuentaContable(cuentaID); outputError == nil && cuenta != nil {
		return detalleCuentaLabel(&models.DetalleCuenta{
			Id:     cuenta.Id,
			Codigo: cuenta.Codigo,
			Nombre: cuenta.Nombre,
		})
	}

	return cuentaID
}

func detalleCuentaLabel(cuenta *models.DetalleCuenta) string {
	if cuenta == nil {
		return ""
	}
	if cuenta.Codigo != "" && cuenta.Nombre != "" {
		return cuenta.Codigo + " - " + cuenta.Nombre
	}
	if cuenta.Codigo != "" {
		return cuenta.Codigo
	}
	if cuenta.Nombre != "" {
		return cuenta.Nombre
	}
	return cuenta.Id
}

func funcionarioSalidaLabel(movimiento *models.Movimiento) string {
	if movimiento == nil {
		return ""
	}

	var formato models.FormatoSalida
	if outputError := utilsHelper.Unmarshal(movimiento.Detalle, &formato); outputError != nil {
		return ""
	}
	if formato.Funcionario <= 0 {
		return ""
	}

	tercero, outputError := getNombreTerceroByID(formato.Funcionario)
	if outputError != nil || tercero == nil {
		return strconv.Itoa(formato.Funcionario)
	}

	return identificacionTerceroLabel(tercero)
}

func identificacionTerceroLabel(tercero *models.IdentificacionTercero) string {
	if tercero == nil {
		return ""
	}
	if tercero.Numero != "" && tercero.NombreCompleto != "" {
		return tercero.Numero + " - " + tercero.NombreCompleto
	}
	if tercero.NombreCompleto != "" {
		return tercero.NombreCompleto
	}
	if tercero.Numero != "" {
		return tercero.Numero
	}
	if tercero.Id > 0 {
		return strconv.Itoa(tercero.Id)
	}
	return ""
}

func trasladosAsociadosLabel(traslados []*models.Movimiento) string {
	if len(traslados) == 0 {
		return ""
	}

	labels := make([]string, 0, len(traslados))
	for _, traslado := range traslados {
		if traslado == nil {
			continue
		}

		label := stringPtrValue(traslado.Consecutivo)
		if label == "" && traslado.Id > 0 {
			label = strconv.Itoa(traslado.Id)
		}
		if label != "" {
			labels = append(labels, label)
		}
	}

	return strings.Join(labels, " | ")
}

func movimientoTieneEstado(movimiento *models.Movimiento, estadosPermitidos ...string) bool {
	if movimiento == nil || movimiento.EstadoMovimientoId == nil {
		return false
	}

	for _, estado := range estadosPermitidos {
		if movimiento.EstadoMovimientoId.Nombre == estado {
			return true
		}
	}

	return false
}

func movimientoConsecutivo(salida *salidaReporteData) string {
	if salida == nil || salida.Movimiento == nil {
		return ""
	}
	return stringPtrValue(salida.Movimiento.Consecutivo)
}

func movimientoEstado(salida *salidaReporteData) string {
	if salida == nil || salida.Movimiento == nil {
		return ""
	}
	return estadoMovimientoNombre(salida.Movimiento)
}

func movimientoFechaCreacion(salida *salidaReporteData) time.Time {
	if salida == nil || salida.Movimiento == nil {
		return time.Time{}
	}
	return salida.Movimiento.FechaCreacion
}

func movimientoFechaCorte(salida *salidaReporteData) *time.Time {
	if salida == nil {
		return nil
	}
	return fechaReporteSalidaMovimiento(salida.Movimiento)
}

func movimientoEnRangoPorFechaReporteSalida(movimiento *models.Movimiento, fechaInicial, fechaFinal time.Time) bool {
	fecha := fechaReporteSalidaMovimiento(movimiento)
	if fecha == nil {
		return false
	}

	inicio := time.Date(fechaInicial.Year(), fechaInicial.Month(), fechaInicial.Day(), 0, 0, 0, 0, time.UTC)
	fin := time.Date(fechaFinal.Year(), fechaFinal.Month(), fechaFinal.Day(), 23, 59, 59, 0, time.UTC)
	return !fecha.Before(inicio) && !fecha.After(fin)
}

func fechaReporteSalidaMovimiento(movimiento *models.Movimiento) *time.Time {
	if movimiento == nil {
		return nil
	}

	fechaSalida := movimiento.FechaCorte
	if fechaSalida == nil {
		if movimiento.MovimientoPadreId != nil {
			return &movimiento.MovimientoPadreId.FechaCreacion
		}
		return nil
	}

	if movimiento.MovimientoPadreId != nil {
		fechaEntrada := movimiento.MovimientoPadreId.FechaCreacion
		if fechaSalida.Before(fechaEntrada) {
			return &fechaEntrada
		}
	}

	return fechaSalida
}

func salidaFuncionario(salida *salidaReporteData) string {
	if salida == nil {
		return ""
	}
	return salida.FuncionarioAsignado
}

func salidaCentroCosto(salida *salidaReporteData) string {
	if salida == nil {
		return ""
	}
	return salida.CentroCostoNombre
}

func salidaCodigoCentroCosto(salida *salidaReporteData) string {
	if salida == nil {
		return ""
	}
	return salida.CentroCostoCodigo
}

func transaccionContableSalida(salida *salidaReporteData) *models.InfoTransaccionContable {
	if salida == nil {
		return nil
	}
	return salida.TransaccionContable
}

func tipoEntradaNombre(movimiento *models.Movimiento) string {
	if movimiento == nil || movimiento.FormatoTipoMovimientoId == nil {
		return ""
	}
	return movimiento.FormatoTipoMovimientoId.Nombre
}

func vidaUtilCatalogo(elemento *models.DetalleElemento) float64 {
	if elemento == nil || elemento.SubgrupoCatalogoId == nil {
		return 0
	}
	return elemento.SubgrupoCatalogoId.VidaUtil
}

func salidaUbicacionInfo(movimiento *models.Movimiento) (centroCostoNombre, centroCostoCodigo string) {
	if movimiento == nil {
		return "", ""
	}

	var detalle models.FormatoSalidaCostos
	if outputError := utilsHelper.Unmarshal(movimiento.Detalle, &detalle); outputError != nil {
		return "", ""
	}

	if detalle.CentroCostos == "" {
		return "", ""
	}

	if centrosCostos, outputError := consultarCentroCostosFn("query=Id:" + detalle.CentroCostos); outputError == nil && len(centrosCostos) > 0 {
		centroCostoNombre = centrosCostos[0].Nombre
		centroCostoCodigo = centrosCostos[0].Codigo
	}

	return centroCostoNombre, centroCostoCodigo
}

func subgrupoInfo(elemento *models.DetalleElemento) (id int, codigo, nombre string) {
	if elemento == nil || elemento.SubgrupoCatalogoId == nil || elemento.SubgrupoCatalogoId.SubgrupoId == nil {
		return 0, "", ""
	}

	return elemento.SubgrupoCatalogoId.SubgrupoId.Id,
		elemento.SubgrupoCatalogoId.SubgrupoId.Codigo,
		elemento.SubgrupoCatalogoId.SubgrupoId.Nombre
}

func tipoBienInfo(elemento *models.DetalleElemento) (id int, nombre string) {
	if elemento == nil || elemento.TipoBienId == nil {
		return 0, ""
	}

	return elemento.TipoBienId.Id, elemento.TipoBienId.Nombre
}

func estadoMovimientoNombre(movimiento *models.Movimiento) string {
	if movimiento == nil || movimiento.EstadoMovimientoId == nil {
		return ""
	}
	return movimiento.EstadoMovimientoId.Nombre
}

func transaccionConcepto(transaccion *models.InfoTransaccionContable) string {
	if transaccion == nil {
		return ""
	}
	return transaccion.Concepto
}

func transaccionFecha(transaccion *models.InfoTransaccionContable) time.Time {
	if transaccion == nil {
		return time.Time{}
	}
	return transaccion.Fecha
}

func stringPtrValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func optionalIntCell(value int) string {
	if value <= 0 {
		return ""
	}
	return strconv.Itoa(value)
}

func floatCell(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func timeCell(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(dateTimeLayout)
}

func timePtrCell(value *time.Time) string {
	if value == nil {
		return ""
	}
	return timeCell(*value)
}

func setColumnWidths(hoja *xlsx.Sheet) {
	_ = hoja.SetColWidth(0, 1, 14)   // Vigencia, Periodo
	_ = hoja.SetColWidth(2, 2, 30)   // Nombre / Descripción
	_ = hoja.SetColWidth(3, 5, 18)   // Marca, Serie, Placa
	_ = hoja.SetColWidth(6, 6, 12)   // Cantidad
	_ = hoja.SetColWidth(7, 11, 18)  // Valores monetarios
	_ = hoja.SetColWidth(12, 12, 20) // Clase
	_ = hoja.SetColWidth(13, 13, 30) // Proveedor
	_ = hoja.SetColWidth(14, 15, 20) // Factura
	_ = hoja.SetColWidth(16, 17, 18) // Vida útil, Tipo entrada
	_ = hoja.SetColWidth(18, 19, 22) // Fecha entrada, Consecutivo Entrada
	_ = hoja.SetColWidth(20, 21, 36) // Cuentas entrada
	_ = hoja.SetColWidth(22, 24, 22) // Fecha salida, Consecutivo salida, Funcionario
	_ = hoja.SetColWidth(25, 26, 24) // Centro de costo, codigo centro de costo
	_ = hoja.SetColWidth(27, 28, 36) // Cuentas salida
	_ = hoja.SetColWidth(29, 31, 18) // Depreciación, Meses, Vida útil
	_ = hoja.SetColWidth(32, 42, 22) // Columnas extras al final
}
