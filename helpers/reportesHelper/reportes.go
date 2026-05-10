package reportesHelper

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/catalogoElementosHelper"
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/entradaHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

const (
	dateLayout     = "2006-01-02"
	dateTimeLayout = "2006-01-02 15:04:05"
	sheetName      = "EntradasElementos"
)

type entradaReporteData struct {
	Movimiento          *models.Movimiento
	Formato             models.FormatoBaseEntrada
	TransaccionContable *models.InfoTransaccionContable
	Elementos           []*models.DetalleElemento
	CuentasPorSubgrupo  map[int]models.CuentasSubgrupo
}

type reporteElementoEntradaRow struct {
	EntradaID                          int
	EntradaConsecutivo                 string
	EntradaEstado                      string
	EntradaFechaCreacion               time.Time
	EntradaFechaCorte                  *time.Time
	EntradaActaRecibidoID              int
	EntradaConceptoTransaccionContable string
	EntradaFechaTransaccionContable    time.Time
	ElementoID                         int
	ElementoNombre                     string
	ElementoCantidad                   int
	ElementoMarca                      string
	ElementoSerie                      string
	ElementoUnidadMedida               int
	ElementoValorUnitario              float64
	ElementoSubtotal                   float64
	ElementoDescuento                  float64
	ElementoValorTotal                 float64
	ElementoPorcentajeIvaID            int
	ElementoValorIva                   float64
	ElementoValorFinal                 float64
	ElementoSubgrupoID                 int
	ElementoSubgrupoCodigo             string
	ElementoSubgrupoNombre             string
	ElementoTipoBienID                 int
	ElementoTipoBienNombre             string
	ElementoActaRecibidoID             int
	ElementoPlaca                      string
	ElementoActivo                     bool
	ElementoFechaCreacion              time.Time
	ElementoFechaModificacion          time.Time
	CuentaDebitoEntrada                string
	CuentaCreditoEntrada               string
}

var (
	reporteElementosHeaders = []string{
		"entrada_id",
		"entrada_consecutivo",
		"entrada_estado",
		"entrada_fecha_creacion",
		"entrada_fecha_corte",
		"entrada_acta_recibido_id",
		"entrada_concepto_transaccion_contable",
		"entrada_fecha_transaccion_contable",
		"elemento_id",
		"elemento_nombre",
		"elemento_cantidad",
		"elemento_marca",
		"elemento_serie",
		"elemento_unidad_medida",
		"elemento_valor_unitario",
		"elemento_subtotal",
		"elemento_descuento",
		"elemento_valor_total",
		"elemento_porcentaje_iva_id",
		"elemento_valor_iva",
		"elemento_valor_final",
		"elemento_subgrupo_id",
		"elemento_subgrupo_codigo",
		"elemento_subgrupo_nombre",
		"elemento_tipo_bien_id",
		"elemento_tipo_bien_nombre",
		"elemento_acta_recibido_id",
		"elemento_placa",
		"elemento_activo",
		"elemento_fecha_creacion",
		"elemento_fecha_modificacion",
		"cuenta_debito_entrada",
		"cuenta_credito_entrada",
	}

	consultarEntradasReporteData = consultarEntradasReporteDataDefault
	consultarCuentaContable      = cuentasContables.GetCuentaContable
)

// GenerarReporteElementos genera un archivo Excel en base64 con una fila por
// elemento, incluyendo la entrada asociada y las cuentas débito/crédito de la entrada.
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

func consultarCodigosEntrada() (codigos []string, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("consultarCodigosEntrada - Unhandled Error!", "500")

	formatos, outputError := movimientosArka.GetAllFormatoTipoMovimiento("limit=-1&fields=CodigoAbreviacion&query=Activo:true")
	if outputError != nil {
		return nil, outputError
	}

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

	return codigos, nil
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

	movimientos, _, outputError = movimientosArka.GetAllMovimiento(params.Encode())
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

	detalleEntrada, outputError := entradaHelper.DetalleEntrada(movimiento.Id)
	if outputError != nil {
		return nil, outputError
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

	cuentasPorSubgrupo := make(map[int]models.CuentasSubgrupo)
	if movimiento.FormatoTipoMovimientoId != nil && movimiento.FormatoTipoMovimientoId.Id > 0 {
		subgrupos := collectSubgrupoIDs(elementos)
		if len(subgrupos) > 0 {
			outputError = catalogoElementosHelper.GetCuentasByMovimientoAndSubgrupos(movimiento.FormatoTipoMovimientoId.Id, subgrupos, cuentasPorSubgrupo)
			if outputError != nil {
				return nil, outputError
			}
		}
	}

	entrada = &entradaReporteData{
		Movimiento:          movimiento,
		Formato:             formato,
		TransaccionContable: extractTransaccionContable(detalleEntrada),
		Elementos:           elementos,
		CuentasPorSubgrupo:  cuentasPorSubgrupo,
	}

	return entrada, nil
}

func resolverElementosEntrada(formato models.FormatoBaseEntrada) (elementos []*models.DetalleElemento, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("resolverElementosEntrada - Unhandled Error!", "500")

	if formato.ActaRecibidoId > 0 {
		return actaRecibido.GetElementos(formato.ActaRecibidoId, nil)
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

	elementos, outputError = actaRecibido.GetElementos(0, elementosActaIds)
	if outputError != nil {
		return nil, outputError
	}

	return ordenarElementosPorIds(elementosActaIds, elementos), nil
}

func extractTransaccionContable(detalle map[string]interface{}) *models.InfoTransaccionContable {
	if detalle == nil {
		return nil
	}

	transaccion, ok := detalle["TransaccionContable"]
	if !ok || transaccion == nil {
		return nil
	}

	switch value := transaccion.(type) {
	case models.InfoTransaccionContable:
		result := value
		return &result
	case *models.InfoTransaccionContable:
		return value
	default:
		return nil
	}
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

func construirFilasReporteEntradas(entradas []*entradaReporteData) []*reporteElementoEntradaRow {
	rows := make([]*reporteElementoEntradaRow, 0)

	for _, entrada := range entradas {
		if entrada == nil || entrada.Movimiento == nil {
			continue
		}

		for _, elemento := range entrada.Elementos {
			if elemento == nil {
				continue
			}

			subgrupoID, subgrupoCodigo, subgrupoNombre := subgrupoInfo(elemento)
			tipoBienID, tipoBienNombre := tipoBienInfo(elemento)
			actaRecibidoID := 0
			if elemento.ActaRecibidoId != nil {
				actaRecibidoID = elemento.ActaRecibidoId.Id
			}

			movimientoCuenta := entrada.CuentasPorSubgrupo[subgrupoID]
			row := &reporteElementoEntradaRow{
				EntradaID:                          entrada.Movimiento.Id,
				EntradaConsecutivo:                 stringPtrValue(entrada.Movimiento.Consecutivo),
				EntradaEstado:                      estadoMovimientoNombre(entrada.Movimiento),
				EntradaFechaCreacion:               entrada.Movimiento.FechaCreacion,
				EntradaFechaCorte:                  entrada.Movimiento.FechaCorte,
				EntradaActaRecibidoID:              entrada.Formato.ActaRecibidoId,
				EntradaConceptoTransaccionContable: transaccionConcepto(entrada.TransaccionContable),
				EntradaFechaTransaccionContable:    transaccionFecha(entrada.TransaccionContable),
				ElementoID:                         elemento.Id,
				ElementoNombre:                     elemento.Nombre,
				ElementoCantidad:                   elemento.Cantidad,
				ElementoMarca:                      elemento.Marca,
				ElementoSerie:                      elemento.Serie,
				ElementoUnidadMedida:               elemento.UnidadMedida,
				ElementoValorUnitario:              elemento.ValorUnitario,
				ElementoSubtotal:                   elemento.Subtotal,
				ElementoDescuento:                  elemento.Descuento,
				ElementoValorTotal:                 elemento.ValorTotal,
				ElementoPorcentajeIvaID:            elemento.PorcentajeIvaId,
				ElementoValorIva:                   elemento.ValorIva,
				ElementoValorFinal:                 elemento.ValorFinal,
				ElementoSubgrupoID:                 subgrupoID,
				ElementoSubgrupoCodigo:             subgrupoCodigo,
				ElementoSubgrupoNombre:             subgrupoNombre,
				ElementoTipoBienID:                 tipoBienID,
				ElementoTipoBienNombre:             tipoBienNombre,
				ElementoActaRecibidoID:             actaRecibidoID,
				ElementoPlaca:                      elemento.Placa,
				ElementoActivo:                     elemento.Activo,
				ElementoFechaCreacion:              elemento.FechaCreacion,
				ElementoFechaModificacion:          elemento.FechaModificacion,
				CuentaDebitoEntrada:                resolveCuentaMovimientoLabel(movimientoCuenta.CuentaDebitoId, entrada.TransaccionContable, true),
				CuentaCreditoEntrada:               resolveCuentaMovimientoLabel(movimientoCuenta.CuentaCreditoId, entrada.TransaccionContable, false),
			}
			rows = append(rows, row)
		}
	}

	return rows
}

func addElementoEntradaRow(hoja *xlsx.Sheet, rowData *reporteElementoEntradaRow) {
	if rowData == nil {
		return
	}

	row := hoja.AddRow()
	values := []string{
		strconv.Itoa(rowData.EntradaID),
		rowData.EntradaConsecutivo,
		rowData.EntradaEstado,
		timeCell(rowData.EntradaFechaCreacion),
		timePtrCell(rowData.EntradaFechaCorte),
		strconv.Itoa(rowData.EntradaActaRecibidoID),
		rowData.EntradaConceptoTransaccionContable,
		timeCell(rowData.EntradaFechaTransaccionContable),
		strconv.Itoa(rowData.ElementoID),
		rowData.ElementoNombre,
		strconv.Itoa(rowData.ElementoCantidad),
		rowData.ElementoMarca,
		rowData.ElementoSerie,
		strconv.Itoa(rowData.ElementoUnidadMedida),
		floatCell(rowData.ElementoValorUnitario),
		floatCell(rowData.ElementoSubtotal),
		floatCell(rowData.ElementoDescuento),
		floatCell(rowData.ElementoValorTotal),
		strconv.Itoa(rowData.ElementoPorcentajeIvaID),
		floatCell(rowData.ElementoValorIva),
		floatCell(rowData.ElementoValorFinal),
		strconv.Itoa(rowData.ElementoSubgrupoID),
		rowData.ElementoSubgrupoCodigo,
		rowData.ElementoSubgrupoNombre,
		strconv.Itoa(rowData.ElementoTipoBienID),
		rowData.ElementoTipoBienNombre,
		strconv.Itoa(rowData.ElementoActaRecibidoID),
		rowData.ElementoPlaca,
		strconv.FormatBool(rowData.ElementoActivo),
		timeCell(rowData.ElementoFechaCreacion),
		timeCell(rowData.ElementoFechaModificacion),
		rowData.CuentaDebitoEntrada,
		rowData.CuentaCreditoEntrada,
	}

	for _, value := range values {
		row.AddCell().SetString(value)
	}
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
	_ = hoja.SetColWidth(0, 7, 22)
	_ = hoja.SetColWidth(8, 13, 18)
	_ = hoja.SetColWidth(14, 20, 16)
	_ = hoja.SetColWidth(21, 27, 20)
	_ = hoja.SetColWidth(28, 30, 22)
	_ = hoja.SetColWidth(31, 32, 36)
}
