package reportesHelper

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
	crudActaRecibido "github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
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

type reporteContabilizacionGrupo struct {
	TipoMovimientoID  int
	Consecutivo       string
	Fecha             time.Time
	Observacion       string
	ActaID            int
	CentroCostoNombre string
	CentroCostoCodigo string
	SubgrupoID        int
	SubgrupoNombre    string
	ValorTotal        float64
	CuentaDebito      string
	CuentaCredito     string
}

type reporteContabilizacionCuentaCache struct {
	CuentaDebitoID  string
	CuentaCreditoID string
	CuentaDebito    string
	CuentaCredito   string
}

type reporteContabilizacionRow struct {
	Cuenta            string
	Naturaleza        string
	Consecutivo       string
	Fecha             time.Time
	Observacion       string
	Valor             float64
	Clase             string
	CentroCostoNombre string
	CentroCostoCodigo string
}

var (
	reporteContabilizacionHeaders = []string{
		"Cuenta",
		"Naturaleza",
		"Consecutivo",
		"Fecha",
		"Observacion",
		"Valor",
		"Clase",
		"CentroCostoNombre",
		"CentroCostoCodigo",
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

	addReporteContabilizacionHeaders(hojaEntradas)
	addReporteContabilizacionHeaders(hojaSalidas)

	for _, row := range expandirFilasReporteContabilizacion(entradas) {
		addReporteContabilizacionRow(hojaEntradas, row)
	}
	for _, row := range expandirFilasReporteContabilizacion(salidas) {
		addReporteContabilizacionRow(hojaSalidas, row)
	}

	setReporteContabilizacionColumnWidths(hojaEntradas)
	setReporteContabilizacionColumnWidths(hojaSalidas)

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

func consultarEntradasContabilizacionReporteDataDefault(fechaInicial, fechaFinal time.Time) (entradas []*reporteContabilizacionGrupo, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarEntradasContabilizacionReporteDataDefault - Unhandled Error!", "500")

	codigosEntrada, outputError := consultarCodigosEntrada()
	if outputError != nil {
		return nil, outputError
	}
	if len(codigosEntrada) == 0 {
		return []*reporteContabilizacionGrupo{}, nil
	}

	formatosEntrada, outputError := consultarFormatosPorCodigos(codigosEntrada)
	if outputError != nil {
		return nil, outputError
	}

	cuentasPorMovimientoYSubgrupo, outputError := consultarCuentasContabilizacionPorMovimiento(formatosEntrada)
	if outputError != nil {
		return nil, outputError
	}

	movimientos, outputError := consultarMovimientosContabilizacionPorFechaYEstado(fechaInicial, fechaFinal, codigosEntrada, "Entrada Con Salida")
	if outputError != nil {
		return nil, outputError
	}

	entradas = make([]*reporteContabilizacionGrupo, 0)
	for _, movimiento := range movimientos {
		if movimiento == nil || movimiento.Id <= 0 {
			continue
		}

		var formato models.FormatoBaseEntrada
		if strings.TrimSpace(movimiento.Detalle) != "" {
			if outputError = utilsHelper.Unmarshal(movimiento.Detalle, &formato); outputError != nil {
				return nil, outputError
			}
		}

		var elementos []*models.DetalleElemento
		if formato.ActaRecibidoId > 0 {
			elementos, outputError = consultarElementosActa(formato.ActaRecibidoId, nil)
		} else {
			elementos, outputError = resolverElementosEntrada(formato)
		}
		if outputError != nil {
			return nil, outputError
		}

		centroCostoNombre, centroCostoCodigo, outputError := centroCostoContabilizacionEntradaInfo(movimiento, formato, elementos)
		if outputError != nil {
			return nil, outputError
		}

		agrupados := agruparElementosContabilizacion(
			movimiento,
			formato.ActaRecibidoId,
			centroCostoNombre,
			centroCostoCodigo,
			elementos,
			cuentasPorMovimientoYSubgrupo[movimientoTipoID(movimiento)],
		)
		entradas = append(entradas, agrupados...)
	}

	return entradas, nil
}

func consultarSalidasContabilizacionReporteDataDefault(fechaInicial, fechaFinal time.Time) (salidas []*reporteContabilizacionGrupo, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarSalidasContabilizacionReporteDataDefault - Unhandled Error!", "500")

	codigosSalida, outputError := consultarCodigosSalida()
	if outputError != nil {
		return nil, outputError
	}
	if len(codigosSalida) == 0 {
		return []*reporteContabilizacionGrupo{}, nil
	}

	formatosSalida, outputError := consultarFormatosPorCodigos(codigosSalida)
	if outputError != nil {
		return nil, outputError
	}

	cuentasPorMovimientoYSubgrupo, outputError := consultarCuentasContabilizacionPorMovimiento(formatosSalida)
	if outputError != nil {
		return nil, outputError
	}

	movimientos, outputError := consultarMovimientosContabilizacionPorFechaYEstado(fechaInicial, fechaFinal, codigosSalida, "Salida Aprobada")
	if outputError != nil {
		return nil, outputError
	}

	salidas = make([]*reporteContabilizacionGrupo, 0)
	for _, movimiento := range movimientos {
		if movimiento == nil || movimiento.Id <= 0 {
			continue
		}

		trSalida, outputError := consultarTrSalida(movimiento.Id)
		if outputError != nil {
			return nil, outputError
		}
		if trSalida == nil || trSalida.Salida == nil {
			continue
		}

		entradaPadre, actaID, outputError := entradaPadreYActaIDSalida(movimiento)
		if outputError != nil {
			return nil, outputError
		}
		if entradaPadre != nil {
			trSalida.Salida.MovimientoPadreId = entradaPadre
		}

		elementos, outputError := resolverElementosSalida(trSalida)
		if outputError != nil {
			return nil, outputError
		}

		centroCostoNombre, centroCostoCodigo := salidaUbicacionInfo(trSalida.Salida)

		agrupados := agruparElementosContabilizacion(
			trSalida.Salida,
			actaID,
			centroCostoNombre,
			centroCostoCodigo,
			elementos,
			cuentasPorMovimientoYSubgrupo[movimientoTipoID(trSalida.Salida)],
		)
		salidas = append(salidas, agrupados...)
	}

	return salidas, nil
}

func consultarMovimientosContabilizacionPorFechaYEstado(fechaInicial, fechaFinal time.Time, codigos []string, estado string) (movimientos []*models.Movimiento, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarMovimientosContabilizacionPorFechaYEstado - Unhandled Error!", "500")

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
		"Activo:true,EstadoMovimientoId__Nombre:"+estado+
			",FormatoTipoMovimientoId__CodigoAbreviacion__in:"+strings.Join(codigos, "|")+
			",FechaCreacion__gte:"+fechaInicial.Format(time.RFC3339)+
			",FechaCreacion__lte:"+fechaFinal.Format(time.RFC3339),
	)

	movimientos, _, outputError = consultarMovimientosReporteFn(params.Encode())
	if outputError != nil {
		return nil, outputError
	}

	return movimientos, nil
}

func consultarFormatosPorCodigos(codigos []string) (ids []int, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarFormatosPorCodigos - Unhandled Error!", "500")

	if len(codigos) == 0 {
		return nil, nil
	}

	formatos, outputError := movimientosArka.GetAllFormatoTipoMovimiento("limit=-1&fields=Id,CodigoAbreviacion")
	if outputError != nil {
		return nil, outputError
	}

	codigosSet := make(map[string]struct{}, len(codigos))
	for _, codigo := range codigos {
		codigosSet[codigo] = struct{}{}
	}

	seen := make(map[int]struct{})
	for _, formato := range formatos {
		if formato == nil {
			continue
		}
		if _, ok := codigosSet[formato.CodigoAbreviacion]; !ok {
			continue
		}
		if _, ok := seen[formato.Id]; ok || formato.Id <= 0 {
			continue
		}
		seen[formato.Id] = struct{}{}
		ids = append(ids, formato.Id)
	}

	sort.Ints(ids)
	return ids, nil
}

func consultarCuentasContabilizacionPorMovimiento(movimientoIDs []int) (map[int]map[int]reporteContabilizacionCuentaCache, map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarCuentasContabilizacionPorMovimiento - Unhandled Error!", "500")

	cuentasPorMovimientoYSubgrupo := make(map[int]map[int]reporteContabilizacionCuentaCache)
	if len(movimientoIDs) == 0 {
		return cuentasPorMovimientoYSubgrupo, nil
	}

	query := "limit=-1&fields=Id,CuentaDebitoId,CuentaCreditoId,SubgrupoId,SubtipoMovimientoId" +
		"&sortby=Id&order=desc&query=Activo:true,SubtipoMovimientoId__in:" + utilsHelper.ArrayToString(movimientoIDs, "|")
	cuentasSubgrupo, outputError := catalogoElementos.GetAllCuentasSubgrupo(query)
	if outputError != nil {
		return nil, outputError
	}

	codigosCuenta := make(map[string]string)
	for _, cuentaSubgrupo := range cuentasSubgrupo {
		if cuentaSubgrupo == nil || cuentaSubgrupo.SubgrupoId == nil || cuentaSubgrupo.SubtipoMovimientoId <= 0 || cuentaSubgrupo.SubgrupoId.Id <= 0 {
			continue
		}
		if _, ok := cuentasPorMovimientoYSubgrupo[cuentaSubgrupo.SubtipoMovimientoId]; !ok {
			cuentasPorMovimientoYSubgrupo[cuentaSubgrupo.SubtipoMovimientoId] = make(map[int]reporteContabilizacionCuentaCache)
		}
		if _, ok := cuentasPorMovimientoYSubgrupo[cuentaSubgrupo.SubtipoMovimientoId][cuentaSubgrupo.SubgrupoId.Id]; ok {
			continue
		}

		cuentasPorMovimientoYSubgrupo[cuentaSubgrupo.SubtipoMovimientoId][cuentaSubgrupo.SubgrupoId.Id] = reporteContabilizacionCuentaCache{
			CuentaDebitoID:  strings.TrimSpace(cuentaSubgrupo.CuentaDebitoId),
			CuentaCreditoID: strings.TrimSpace(cuentaSubgrupo.CuentaCreditoId),
			CuentaDebito:    consultarCodigoCuentaContabilizacion(strings.TrimSpace(cuentaSubgrupo.CuentaDebitoId), codigosCuenta),
			CuentaCredito:   consultarCodigoCuentaContabilizacion(strings.TrimSpace(cuentaSubgrupo.CuentaCreditoId), codigosCuenta),
		}
	}

	return cuentasPorMovimientoYSubgrupo, nil
}

func consultarCodigoCuentaContabilizacion(cuentaID string, cache map[string]string) string {
	cuentaID = strings.TrimSpace(cuentaID)
	if cuentaID == "" {
		return ""
	}
	if codigo, ok := cache[cuentaID]; ok {
		return codigo
	}

	codigo := cuentaID
	if cuenta, outputError := consultarCuentaContable(cuentaID); outputError == nil && cuenta != nil {
		if strings.TrimSpace(cuenta.Codigo) != "" {
			codigo = strings.TrimSpace(cuenta.Codigo)
		} else if strings.TrimSpace(cuenta.Id) != "" {
			codigo = strings.TrimSpace(cuenta.Id)
		}
	}

	cache[cuentaID] = codigo
	return codigo
}

func centroCostoEntradaA11Info(formato models.FormatoBaseEntrada, elementos []*models.DetalleElemento) (codigo, nombre string, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("centroCostoEntradaA11Info - Unhandled Error!", "500")

	actaRecibidoID := resolverActaRecibidoIDCentroCostoEntrada(formato, elementos)
	if actaRecibidoID <= 0 {
		return "", "", nil
	}

	historicos, outputError := consultarHistoricosActaReporteFn(
		"ActaRecibidoId__Id:"+strconv.Itoa(actaRecibidoID),
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

func consultarCentroCostoA11ByID(id string) (codigo, nombre string, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarCentroCostoA11ByID - Unhandled Error!", "500")

	return consultarCentroCostoA11ByReferencia(id)
}

func entradaPadreYActaIDSalida(movimiento *models.Movimiento) (entradaPadre *models.Movimiento, actaID int, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("entradaPadreYActaIDSalida - Unhandled Error!", "500")

	if movimiento == nil || movimiento.MovimientoPadreId == nil || movimiento.MovimientoPadreId.Id <= 0 {
		return nil, 0, nil
	}

	entradaPadre = movimiento.MovimientoPadreId
	if strings.TrimSpace(entradaPadre.Detalle) == "" {
		entradaPadre, outputError = consultarMovimientoPorID(entradaPadre.Id)
		if outputError != nil {
			return nil, 0, outputError
		}
	}
	if entradaPadre == nil || strings.TrimSpace(entradaPadre.Detalle) == "" {
		return entradaPadre, 0, nil
	}

	var formato models.FormatoBaseEntrada
	if outputError = utilsHelper.Unmarshal(entradaPadre.Detalle, &formato); outputError != nil {
		return nil, 0, outputError
	}

	return entradaPadre, formato.ActaRecibidoId, nil
}

func centroCostoContabilizacionEntradaInfo(movimiento *models.Movimiento, formato models.FormatoBaseEntrada, elementos []*models.DetalleElemento) (nombre, codigo string, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("centroCostoContabilizacionEntradaInfo - Unhandled Error!", "500")

	codigoRaw, nombre, outputError := centroCostoEntradaA11Info(formato, elementos)
	if outputError == nil && (nombre != "" || codigoRaw != "") {
		return nombre, normalizarCodigoCentroCostoReporte(codigoRaw), nil
	}

	if centroCostoNombre, centroCostoCodigo, movimientoErr := resolverCentroCostoMovimiento(movimiento); movimientoErr != nil {
		if outputError != nil {
			return "", "", outputError
		}
		return "", "", movimientoErr
	} else if centroCostoNombre != "" || centroCostoCodigo != "" {
		return centroCostoNombre, centroCostoCodigo, nil
	}

	if outputError != nil {
		return "", "", outputError
	}

	return "", "", nil
}

func agruparElementosContabilizacion(
	movimiento *models.Movimiento,
	actaID int,
	centroCostoNombre, centroCostoCodigo string,
	elementos []*models.DetalleElemento,
	cuentasPorSubgrupo map[int]reporteContabilizacionCuentaCache,
) []*reporteContabilizacionGrupo {
	if movimiento == nil || len(elementos) == 0 {
		return nil
	}

	type acumulado struct {
		nombre string
		valor  float64
	}

	acumulados := make(map[int]*acumulado)
	for _, elemento := range elementos {
		subgrupoID, _, subgrupoNombre := subgrupoInfo(elemento)
		if subgrupoID <= 0 {
			continue
		}
		if _, ok := acumulados[subgrupoID]; !ok {
			acumulados[subgrupoID] = &acumulado{nombre: subgrupoNombre}
		}
		acumulados[subgrupoID].valor += elemento.ValorFinal
	}

	subgrupos := make([]int, 0, len(acumulados))
	for subgrupoID := range acumulados {
		subgrupos = append(subgrupos, subgrupoID)
	}
	sort.Ints(subgrupos)

	grupos := make([]*reporteContabilizacionGrupo, 0, len(subgrupos))
	for _, subgrupoID := range subgrupos {
		acum := acumulados[subgrupoID]
		if acum == nil {
			continue
		}

		grupo := &reporteContabilizacionGrupo{
			TipoMovimientoID:  movimientoTipoID(movimiento),
			Consecutivo:       stringPtrValue(movimiento.Consecutivo),
			Fecha:             movimiento.FechaCreacion,
			Observacion:       strings.TrimSpace(movimiento.Observacion),
			ActaID:            actaID,
			CentroCostoNombre: centroCostoNombre,
			CentroCostoCodigo: centroCostoCodigo,
			SubgrupoID:        subgrupoID,
			SubgrupoNombre:    acum.nombre,
			ValorTotal:        roundToTwoDecimals(acum.valor),
		}

		if cuentas, ok := cuentasPorSubgrupo[subgrupoID]; ok {
			grupo.CuentaDebito = cuentas.CuentaDebito
			grupo.CuentaCredito = cuentas.CuentaCredito
		}

		grupos = append(grupos, grupo)
	}

	return grupos
}

func movimientoTipoID(movimiento *models.Movimiento) int {
	if movimiento == nil || movimiento.FormatoTipoMovimientoId == nil {
		return 0
	}
	return movimiento.FormatoTipoMovimientoId.Id
}

func expandirFilasReporteContabilizacion(grupos []*reporteContabilizacionGrupo) []*reporteContabilizacionRow {
	rows := make([]*reporteContabilizacionRow, 0, len(grupos)*2)
	for _, grupo := range grupos {
		if grupo == nil {
			continue
		}

		if strings.TrimSpace(grupo.CuentaDebito) != "" {
			rows = append(rows, &reporteContabilizacionRow{
				Cuenta:            strings.TrimSpace(grupo.CuentaDebito),
				Naturaleza:        "Debito",
				Consecutivo:       grupo.Consecutivo,
				Fecha:             grupo.Fecha,
				Observacion:       grupo.Observacion,
				Valor:             grupo.ValorTotal,
				Clase:             grupo.SubgrupoNombre,
				CentroCostoNombre: grupo.CentroCostoNombre,
				CentroCostoCodigo: grupo.CentroCostoCodigo,
			})
		}
		if strings.TrimSpace(grupo.CuentaCredito) != "" {
			rows = append(rows, &reporteContabilizacionRow{
				Cuenta:            strings.TrimSpace(grupo.CuentaCredito),
				Naturaleza:        "Credito",
				Consecutivo:       grupo.Consecutivo,
				Fecha:             grupo.Fecha,
				Observacion:       grupo.Observacion,
				Valor:             grupo.ValorTotal,
				Clase:             grupo.SubgrupoNombre,
				CentroCostoNombre: grupo.CentroCostoNombre,
				CentroCostoCodigo: grupo.CentroCostoCodigo,
			})
		}
	}
	return rows
}

func addReporteContabilizacionHeaders(hoja *xlsx.Sheet) {
	filaEncabezado := hoja.AddRow()
	for _, header := range reporteContabilizacionHeaders {
		filaEncabezado.AddCell().SetString(header)
	}
}

func addReporteContabilizacionRow(hoja *xlsx.Sheet, rowData *reporteContabilizacionRow) {
	if rowData == nil {
		return
	}

	row := hoja.AddRow()
	addStringCell(row, rowData.Cuenta)
	addStringCell(row, rowData.Naturaleza)
	addStringCell(row, rowData.Consecutivo)
	addDateCell(row, rowData.Fecha)
	addStringCell(row, rowData.Observacion)
	addDecimalCell(row, rowData.Valor)
	addStringCell(row, rowData.Clase)
	addStringCell(row, rowData.CentroCostoNombre)
	addStringCell(row, rowData.CentroCostoCodigo)
}

func setReporteContabilizacionColumnWidths(hoja *xlsx.Sheet) {
	_ = hoja.SetColWidth(0, 0, 18)
	_ = hoja.SetColWidth(1, 1, 14)
	_ = hoja.SetColWidth(2, 2, 24)
	_ = hoja.SetColWidth(3, 3, 14)
	_ = hoja.SetColWidth(4, 4, 42)
	_ = hoja.SetColWidth(5, 5, 16)
	_ = hoja.SetColWidth(6, 6, 28)
	_ = hoja.SetColWidth(7, 7, 32)
	_ = hoja.SetColWidth(8, 8, 20)
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
