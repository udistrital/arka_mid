package reportesHelper

import (
	"bytes"
	"encoding/base64"
	"fmt"
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
	"github.com/udistrital/arka_mid/helpers/crud/oikos"
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
	Sede                string
	Dependencia         string
}

type salidaReporteBaseData struct {
	Movimiento          *models.Movimiento
	TransaccionContable *models.InfoTransaccionContable
	FuncionarioAsignado string
	CuentasPorSubgrupo  map[int]models.CuentasSubgrupo
	Sede                string
	Dependencia         string
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
	EntradaProveedor                   string
	EntradaFacturaConsecutivo          string
	EntradaFacturaFecha                time.Time
	TipoEntrada                        string
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
	ElementoVidaUtilCatalogo           float64
	ElementoActaRecibidoID             int
	ElementoPlaca                      string
	ElementoActivo                     bool
	ElementoFechaCreacion              time.Time
	ElementoFechaModificacion          time.Time
	CuentaDebitoEntrada                string
	CuentaCreditoEntrada               string
	SalidaID                           int
	SalidaConsecutivo                  string
	SalidaEstado                       string
	SalidaFechaCreacion                time.Time
	SalidaFechaCorte                   *time.Time
	SalidaFuncionarioAsignado          string
	SalidaSede                         string
	SalidaDependencia                  string
	SalidaConceptoTransaccionContable  string
	SalidaFechaTransaccionContable     time.Time
	TrasladosAsociados                 string
	CuentaDebitoSalida                 string
	CuentaCreditoSalida                string
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
		"Proveedor",
		"Consecutivo Factura",
		"Fecha Factura",
		"Tipo de entrada",
		"elemento_id",
		"Nombre / Descripción",
		"elemento_cantidad",
		"elemento_marca",
		"elemento_serie",
		"elemento_unidad_medida",
		"elemento_valor_unitario",
		"elemento_subtotal",
		"elemento_descuento",
		"elemento_valor_total",
		"Porcentaje IVA",
		"elemento_valor_iva",
		"elemento_valor_final",
		"elemento_subgrupo_id",
		"elemento_subgrupo_codigo",
		"Clase",
		"elemento_tipo_bien_id",
		"Tipo de bien",
		"Vida útil (años)",
		"elemento_acta_recibido_id",
		"elemento_placa",
		"elemento_activo",
		"elemento_fecha_creacion",
		"elemento_fecha_modificacion",
		"cuenta_debito_entrada",
		"cuenta_credito_entrada",
		"salida_id",
		"salida_consecutivo",
		"salida_estado",
		"salida_fecha_creacion",
		"salida_fecha_corte",
		"salida_funcionario_asignado",
		"Sede",
		"Dependencia",
		"salida_concepto_transaccion_contable",
		"salida_fecha_transaccion_contable",
		"traslados_asociados",
		"cuenta_debito_salida",
		"cuenta_credito_salida",
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

	var formato models.FormatoBaseEntrada
	outputError = utilsHelper.Unmarshal(movimiento.Detalle, &formato)
	if outputError != nil {
		return nil, outputError
	}

	elementos, outputError := resolverElementosEntrada(formato)
	if outputError != nil {
		return nil, outputError
	}

	proveedor, facturaConsecutivo, facturaFecha, outputError := consultarMetadataEntrada(formato)
	if outputError != nil {
		return nil, outputError
	}

	salidasPorElemento, outputError := resolverSalidasPorElemento(elementos)
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

	historicos, outputError := crudActaRecibido.GetAllHistoricoActa("ActaRecibidoId__Id:"+strconv.Itoa(actaRecibidoID), "", "Id", "desc", "", "1")
	if outputError != nil || len(historicos) == 0 || historicos[0].ProveedorId <= 0 {
		return "", outputError
	}

	tercero, outputError := terceros.GetNombreTerceroById(historicos[0].ProveedorId)
	if outputError != nil {
		return "", outputError
	}

	return identificacionTerceroLabel(tercero), nil
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

func consultarTransaccionContableEntrada(movimiento *models.Movimiento) *models.InfoTransaccionContable {
	return consultarTransaccionContableMovimiento(movimiento, "Entrada Aprobada", "Entrada Con Salida")
}

func consultarTransaccionContableSalida(movimiento *models.Movimiento) *models.InfoTransaccionContable {
	return consultarTransaccionContableMovimiento(movimiento, "Salida Aprobada")
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
			Sede:                base.Sede,
			Dependencia:         base.Dependencia,
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

	sede, dependencia := salidaUbicacionInfo(movimiento)

	cuentasPorSubgrupo := make(map[int]models.CuentasSubgrupo)
	if movimiento.FormatoTipoMovimientoId != nil && movimiento.FormatoTipoMovimientoId.Id > 0 && len(subgrupos) > 0 {
		outputError = catalogoElementosHelper.GetCuentasByMovimientoAndSubgrupos(movimiento.FormatoTipoMovimientoId.Id, subgrupos, cuentasPorSubgrupo)
		if outputError != nil {
			return nil, outputError
		}
	}

	salida = &salidaReporteBaseData{
		Movimiento:          movimiento,
		TransaccionContable: consultarTransaccionContableSalida(movimiento),
		FuncionarioAsignado: funcionarioSalidaLabel(movimiento),
		CuentasPorSubgrupo:  cuentasPorSubgrupo,
		Sede:                sede,
		Dependencia:         dependencia,
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
			salida := entrada.SalidasPorElemento[elemento.Id]
			salidaCuenta := models.CuentasSubgrupo{}
			if salida != nil {
				salidaCuenta = salida.CuentasPorSubgrupo[subgrupoID]
			}
			row := &reporteElementoEntradaRow{
				EntradaID:                          entrada.Movimiento.Id,
				EntradaConsecutivo:                 stringPtrValue(entrada.Movimiento.Consecutivo),
				EntradaEstado:                      estadoMovimientoNombre(entrada.Movimiento),
				EntradaFechaCreacion:               entrada.Movimiento.FechaCreacion,
				EntradaFechaCorte:                  entrada.Movimiento.FechaCorte,
				EntradaActaRecibidoID:              entrada.Formato.ActaRecibidoId,
				EntradaConceptoTransaccionContable: transaccionConcepto(entrada.TransaccionContable),
				EntradaFechaTransaccionContable:    transaccionFecha(entrada.TransaccionContable),
				EntradaProveedor:                   entrada.Proveedor,
				EntradaFacturaConsecutivo:          entrada.FacturaConsecutivo,
				EntradaFacturaFecha:                entrada.FacturaFecha,
				TipoEntrada:                        tipoEntradaNombre(entrada.Movimiento),
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
				ElementoVidaUtilCatalogo:           vidaUtilCatalogo(elemento),
				ElementoActaRecibidoID:             actaRecibidoID,
				ElementoPlaca:                      elemento.Placa,
				ElementoActivo:                     elemento.Activo,
				ElementoFechaCreacion:              elemento.FechaCreacion,
				ElementoFechaModificacion:          elemento.FechaModificacion,
				CuentaDebitoEntrada:                resolveCuentaMovimientoLabel(movimientoCuenta.CuentaDebitoId, entrada.TransaccionContable, true),
				CuentaCreditoEntrada:               resolveCuentaMovimientoLabel(movimientoCuenta.CuentaCreditoId, entrada.TransaccionContable, false),
				SalidaID:                           movimientoID(salida),
				SalidaConsecutivo:                  movimientoConsecutivo(salida),
				SalidaEstado:                       movimientoEstado(salida),
				SalidaFechaCreacion:                movimientoFechaCreacion(salida),
				SalidaFechaCorte:                   movimientoFechaCorte(salida),
				SalidaFuncionarioAsignado:          salidaFuncionario(salida),
				SalidaSede:                         salidaSede(salida),
				SalidaDependencia:                  salidaDependencia(salida),
				SalidaConceptoTransaccionContable:  movimientoTransaccionConcepto(salida),
				SalidaFechaTransaccionContable:     movimientoTransaccionFecha(salida),
				TrasladosAsociados:                 movimientoTraslados(salida),
				CuentaDebitoSalida:                 resolveCuentaMovimientoLabel(salidaCuenta.CuentaDebitoId, transaccionContableSalida(salida), true),
				CuentaCreditoSalida:                resolveCuentaMovimientoLabel(salidaCuenta.CuentaCreditoId, transaccionContableSalida(salida), false),
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
	addStringCell(row, strconv.Itoa(rowData.EntradaID))
	addStringCell(row, rowData.EntradaConsecutivo)
	addStringCell(row, rowData.EntradaEstado)
	addStringCell(row, timeCell(rowData.EntradaFechaCreacion))
	addStringCell(row, timePtrCell(rowData.EntradaFechaCorte))
	addStringCell(row, strconv.Itoa(rowData.EntradaActaRecibidoID))
	addStringCell(row, rowData.EntradaConceptoTransaccionContable)
	addStringCell(row, timeCell(rowData.EntradaFechaTransaccionContable))
	addStringCell(row, rowData.EntradaProveedor)
	addStringCell(row, rowData.EntradaFacturaConsecutivo)
	addStringCell(row, timeCell(rowData.EntradaFacturaFecha))
	addStringCell(row, rowData.TipoEntrada)
	addStringCell(row, strconv.Itoa(rowData.ElementoID))
	addStringCell(row, rowData.ElementoNombre)
	addStringCell(row, strconv.Itoa(rowData.ElementoCantidad))
	addStringCell(row, rowData.ElementoMarca)
	addStringCell(row, rowData.ElementoSerie)
	addStringCell(row, strconv.Itoa(rowData.ElementoUnidadMedida))
	addDecimalCell(row, rowData.ElementoValorUnitario)
	addDecimalCell(row, rowData.ElementoSubtotal)
	addDecimalCell(row, rowData.ElementoDescuento)
	addDecimalCell(row, rowData.ElementoValorTotal)
	addStringCell(row, strconv.Itoa(rowData.ElementoPorcentajeIvaID))
	addDecimalCell(row, rowData.ElementoValorIva)
	addDecimalCell(row, rowData.ElementoValorFinal)
	addStringCell(row, strconv.Itoa(rowData.ElementoSubgrupoID))
	addStringCell(row, rowData.ElementoSubgrupoCodigo)
	addStringCell(row, rowData.ElementoSubgrupoNombre)
	addStringCell(row, strconv.Itoa(rowData.ElementoTipoBienID))
	addStringCell(row, rowData.ElementoTipoBienNombre)
	addDecimalCell(row, rowData.ElementoVidaUtilCatalogo)
	addStringCell(row, strconv.Itoa(rowData.ElementoActaRecibidoID))
	addStringCell(row, rowData.ElementoPlaca)
	addStringCell(row, strconv.FormatBool(rowData.ElementoActivo))
	addStringCell(row, timeCell(rowData.ElementoFechaCreacion))
	addStringCell(row, timeCell(rowData.ElementoFechaModificacion))
	addStringCell(row, rowData.CuentaDebitoEntrada)
	addStringCell(row, rowData.CuentaCreditoEntrada)
	addStringCell(row, optionalIntCell(rowData.SalidaID))
	addStringCell(row, rowData.SalidaConsecutivo)
	addStringCell(row, rowData.SalidaEstado)
	addStringCell(row, timeCell(rowData.SalidaFechaCreacion))
	addStringCell(row, timePtrCell(rowData.SalidaFechaCorte))
	addStringCell(row, rowData.SalidaFuncionarioAsignado)
	addStringCell(row, rowData.SalidaSede)
	addStringCell(row, rowData.SalidaDependencia)
	addStringCell(row, rowData.SalidaConceptoTransaccionContable)
	addStringCell(row, timeCell(rowData.SalidaFechaTransaccionContable))
	addStringCell(row, rowData.TrasladosAsociados)
	addStringCell(row, rowData.CuentaDebitoSalida)
	addStringCell(row, rowData.CuentaCreditoSalida)
}

func addStringCell(row *xlsx.Row, value string) {
	row.AddCell().SetString(value)
}

func addDecimalCell(row *xlsx.Row, value float64) {
	cell := row.AddCell()
	cell.SetFloatWithFormat(value, decimalNumFmt)
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

	tercero, outputError := terceros.GetNombreTerceroById(formato.Funcionario)
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

func movimientoID(salida *salidaReporteData) int {
	if salida == nil || salida.Movimiento == nil {
		return 0
	}
	return salida.Movimiento.Id
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
	if salida == nil || salida.Movimiento == nil {
		return nil
	}
	return salida.Movimiento.FechaCorte
}

func salidaFuncionario(salida *salidaReporteData) string {
	if salida == nil {
		return ""
	}
	return salida.FuncionarioAsignado
}

func salidaSede(salida *salidaReporteData) string {
	if salida == nil {
		return ""
	}
	return salida.Sede
}

func salidaDependencia(salida *salidaReporteData) string {
	if salida == nil {
		return ""
	}
	return salida.Dependencia
}

func movimientoTransaccionConcepto(salida *salidaReporteData) string {
	if salida == nil {
		return ""
	}
	return transaccionConcepto(salida.TransaccionContable)
}

func movimientoTransaccionFecha(salida *salidaReporteData) time.Time {
	if salida == nil {
		return time.Time{}
	}
	return transaccionFecha(salida.TransaccionContable)
}

func movimientoTraslados(salida *salidaReporteData) string {
	if salida == nil {
		return ""
	}
	return salida.TrasladosAsociados
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

func salidaSedeLabel(movimiento *models.Movimiento) string {
	sede, _ := salidaUbicacionInfo(movimiento)
	return sede
}

func salidaDependenciaLabel(movimiento *models.Movimiento) string {
	_, dependencia := salidaUbicacionInfo(movimiento)
	return dependencia
}

func salidaUbicacionInfo(movimiento *models.Movimiento) (sede, dependencia string) {
	if movimiento == nil {
		return "", ""
	}

	var detalle models.FormatoSalidaCostos
	if outputError := utilsHelper.Unmarshal(movimiento.Detalle, &detalle); outputError != nil {
		return "", ""
	}

	if detalle.Ubicacion > 0 {
		if ubicacion, outputError := oikos.GetSedeDependenciaUbicacion(detalle.Ubicacion); outputError == nil && ubicacion != nil {
			if ubicacion.Sede != nil {
				sede = ubicacion.Sede.Nombre
			}
			if ubicacion.Dependencia != nil {
				dependencia = ubicacion.Dependencia.Nombre
			}
		}
	}

	if (sede == "" || dependencia == "") && detalle.CentroCostos != "" {
		if centrosCostos, outputError := movimientosArka.GetAllCentroCostos("query=Codigo:" + detalle.CentroCostos); outputError == nil && len(centrosCostos) > 0 {
			if sede == "" {
				sede = centrosCostos[0].Sede
			}
			if dependencia == "" {
				if centrosCostos[0].Dependencia != "" {
					dependencia = centrosCostos[0].Dependencia
				} else {
					dependencia = centrosCostos[0].Nombre
				}
			}
		}
	}

	return sede, dependencia
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
	_ = hoja.SetColWidth(0, 11, 22)
	_ = hoja.SetColWidth(12, 17, 18)
	_ = hoja.SetColWidth(18, 25, 16)
	_ = hoja.SetColWidth(26, 33, 20)
	_ = hoja.SetColWidth(34, 37, 22)
	_ = hoja.SetColWidth(38, 39, 36)
	_ = hoja.SetColWidth(40, 51, 24)
	_ = hoja.SetColWidth(52, 53, 36)
}
