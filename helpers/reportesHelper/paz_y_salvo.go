package reportesHelper

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/phpdave11/gofpdf"
	crudTerceros "github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/mid/administrativa"
	tercerosmid "github.com/udistrital/arka_mid/helpers/mid/terceros"
	trasladoshelper "github.com/udistrital/arka_mid/helpers/trasladosHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

const pdfMimeType = "application/pdf"
const cargoFirmantePazYSalvo = "JEFE DE SECCIÓN DE ALMACÉN GENERAL E INVENTARIOS"

var (
	consultarTerceroPorDocumentoFn = consultarTerceroPorDocumento
	consultarInventarioTerceroFn   = consultarInventarioTercero
	consultarElaboradorPazYSalvoFn = consultarElaboradorPazYSalvo
	consultarResponsableFirmaFn    = consultarResponsableFirma
)

type pazYSalvoFirmante struct {
	Nombre          string
	Cargo           string
	TipoDocumento   string
	NumeroDocumento string
}

type supervisorContrato struct {
	Id          int       `json:"Id"`
	Nombre      string    `json:"Nombre"`
	Documento   int       `json:"Documento"`
	Cargo       string    `json:"Cargo"`
	FechaInicio time.Time `json:"FechaInicio"`
	FechaFin    time.Time `json:"FechaFin"`
}

func GenerarPazYSalvo(req *models.PazYSalvoRequest) (respuesta *models.PazYSalvoResponse, outputError map[string]interface{}) {
	defer errorCtrl.ErrorControlFunction("GenerarPazYSalvo - Unhandled Error!", "500")

	if req == nil {
		return nil, errorCtrl.Error("GenerarPazYSalvo - req", "request nil", "400")
	}

	numeroDocumento := strings.TrimSpace(req.NumeroDocumento)
	if numeroDocumento == "" {
		return nil, errorCtrl.Error("GenerarPazYSalvo - numero_documento", "se debe indicar un número de documento válido", "400")
	}

	tercero, outputError := consultarTerceroPorDocumentoFn(numeroDocumento)
	if outputError != nil {
		return nil, outputError
	}

	inventario, outputError := consultarInventarioTerceroFn(tercero.Tercero.Id)
	if outputError != nil {
		return nil, outputError
	}

	elaborador, outputError := consultarElaboradorPazYSalvoFn(req.Usuario, req.ElaboroTerceroId)
	if outputError != nil {
		elaborador = nil
	}

	responsableFirma, outputError := consultarResponsableFirmaFn(tercero)
	if outputError != nil {
		return nil, outputError
	}

	terceroResp := construirTerceroResponse(inventario, numeroDocumento)
	puedeGenerar := len(inventario.Elementos) == 0
	mensaje := "Actualmente no se puede generar el paz y salvo porque cuenta con elementos bajo su responsabilidad."
	if puedeGenerar {
		mensaje = "El tercero no cuenta con elementos en inventario. Se generó el paz y salvo."
	}

	nombreArchivo := fmt.Sprintf("paz_y_salvo_%s_%s.pdf", sanitizarNombreArchivo(numeroDocumento), time.Now().Format("20060102"))
	archivoBase64, outputError := construirPDFPazYSalvo(terceroResp, inventario.Elementos, puedeGenerar, elaborador, responsableFirma)
	if outputError != nil {
		return nil, outputError
	}

	respuesta = &models.PazYSalvoResponse{
		ArchivoBase64:         archivoBase64,
		NombreArchivo:         nombreArchivo,
		TipoArchivo:           pdfMimeType,
		Mensaje:               mensaje,
		PuedeGenerarPazYSalvo: puedeGenerar,
		Tercero:               terceroResp,
		Elementos:             inventario.Elementos,
	}

	return respuesta, nil
}

func consultarTerceroPorDocumento(numeroDocumento string) (detalle models.DetalleTercero, outputError map[string]interface{}) {
	funcion := "consultarTerceroPorDocumento"
	payload := "documento=" + url.QueryEscape(numeroDocumento)

	terceros, outputError := crudTerceros.GetAllTrTerceroIdentificacion(payload)
	if outputError != nil {
		return detalle, outputError
	}

	if len(terceros) == 0 || terceros[0].Tercero == nil {
		return detalle, errorCtrl.Error(funcion+" - GetAllTrTerceroIdentificacion", "no se encontró tercero para el documento suministrado", "404")
	}

	if len(terceros) > 1 {
		return detalle, errorCtrl.Error(funcion+" - GetAllTrTerceroIdentificacion", "se encontró más de un tercero para el documento suministrado", "409")
	}

	return terceros[0], nil
}

func consultarInventarioTercero(terceroId int) (inventario *models.InventarioTercero, outputError map[string]interface{}) {
	funcion := "consultarInventarioTercero"

	inventario = new(models.InventarioTercero)
	if err := trasladoshelper.GetElementosTercero(terceroId, inventario); err != nil {
		return nil, errorCtrl.Error(funcion+" - trasladoshelper.GetElementosTercero", err, "502")
	}

	return inventario, nil
}

func consultarElaboradorPazYSalvo(_ string, terceroId int) (*pazYSalvoFirmante, map[string]interface{}) {
	funcion := "consultarElaboradorPazYSalvo"

	if terceroId <= 0 {
		return nil, errorCtrl.Error(funcion+" - elaboro_tercero_id", "no se recibió terceroId válido para resolver Elaboró", "400")
	}

	detalleFuncionario, outputError := tercerosmid.GetDetalleFuncionario(terceroId)
	if outputError != nil || detalleFuncionario == nil || len(detalleFuncionario.Tercero) == 0 {
		return nil, outputError
	}

	detalleTercero := detalleFuncionario.Tercero[0]

	elaborador := &pazYSalvoFirmante{
		Nombre:          strings.TrimSpace(detalleTercero.Tercero.NombreCompleto),
		TipoDocumento:   obtenerTipoDocumento(detalleTercero.Identificacion),
		NumeroDocumento: detalleTercero.Identificacion.Numero,
	}

	if len(detalleFuncionario.Cargo) > 0 && detalleFuncionario.Cargo[0] != nil {
		elaborador.Cargo = strings.TrimSpace(detalleFuncionario.Cargo[0].Nombre)
	}

	return elaborador, nil
}

func consultarResponsableFirma(_ models.DetalleTercero) (*pazYSalvoFirmante, map[string]interface{}) {
	funcion := "consultarResponsableFirma"

	payload := "limit=-1&sortby=Id&order=desc&fields=Id,Nombre,Documento,Cargo,FechaInicio,FechaFin"

	var supervisores []supervisorContrato
	if outputError := administrativa.GetSupervisorByQuery(payload, &supervisores); outputError != nil {
		return nil, outputError
	}

	supervisor := seleccionarSupervisorVigentePorCargo(supervisores, cargoFirmantePazYSalvo, time.Now())

	if supervisor == nil {
		return nil, errorCtrl.Error(funcion+" - supervisor_contrato", "no se encontró supervisor vigente para el cargo requerido", "404")
	}

	if supervisor.Documento <= 0 {
		return nil, errorCtrl.Error(funcion+" - supervisor.Documento", "el supervisor encontrado no tiene documento válido", "502")
	}

	detalleSupervisor, outputError := consultarTerceroPorDocumentoFn(strconv.Itoa(supervisor.Documento))
	if outputError != nil {
		return nil, outputError
	}

	return &pazYSalvoFirmante{
		Nombre:          strings.TrimSpace(supervisor.Nombre),
		Cargo:           strings.TrimSpace(supervisor.Cargo),
		TipoDocumento:   obtenerTipoDocumento(detalleSupervisor.Identificacion),
		NumeroDocumento: strconv.Itoa(supervisor.Documento),
	}, nil
}

func construirTerceroResponse(inventario *models.InventarioTercero, numeroDocumento string) *models.PazYSalvoTerceroResponse {
	resp := &models.PazYSalvoTerceroResponse{
		NumeroDocumento: numeroDocumento,
		TipoDocumento:   "Documento",
	}

	if inventario == nil || len(inventario.Tercero.Tercero) == 0 {
		return resp
	}

	detalle := inventario.Tercero.Tercero[0]
	if detalle.Tercero != nil {
		resp.Id = detalle.Tercero.Id
		resp.NombreCompleto = detalle.Tercero.NombreCompleto
	}
	if detalle.Identificacion != nil {
		if detalle.Identificacion.Numero != "" {
			resp.NumeroDocumento = detalle.Identificacion.Numero
		}
		if detalle.Identificacion.TipoDocumentoId != nil && detalle.Identificacion.TipoDocumentoId.CodigoAbreviacion != "" {
			resp.TipoDocumento = detalle.Identificacion.TipoDocumentoId.CodigoAbreviacion
		}
	}
	if len(inventario.Tercero.Cargo) > 0 {
		resp.Cargo = inventario.Tercero.Cargo[0].Nombre
	}

	return resp
}

func construirPDFPazYSalvo(tercero *models.PazYSalvoTerceroResponse, elementos []models.DetalleElementoPlaca, puedeGenerar bool, elaborador, responsableFirma *pazYSalvoFirmante) (archivoBase64 string, outputError map[string]interface{}) {
	funcion := "construirPDFPazYSalvo"

	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.SetMargins(25, 20, 25)
	pdf.SetAutoPageBreak(true, 35)
	pdf.SetFooterFunc(func() {
		renderFooterPazYSalvo(pdf)
	})
	pdf.SetHeaderFunc(func() {
		renderMarcaDeAguaPazYSalvo(pdf)
	})
	pdf.AddPage()

	tr := pdf.UnicodeTranslatorFromDescriptor("")
	renderHeaderPazYSalvo(pdf, tr)

	pdf.SetFont("Times", "B", 12)
	pdf.CellFormat(0, 7, tr("CERTIFICA"), "", 1, "C", false, 0, "")
	pdf.Ln(14)

	pdf.SetFont("Times", "", 11)
	pdf.MultiCell(0, 6, tr(construirParrafoPrincipal(tercero, puedeGenerar)), "", "J", false)
	pdf.Ln(2)

	if puedeGenerar {
		pdf.MultiCell(0, 6, tr("Se expide la presente a solicitud del interesado, para los fines que estime pertinentes."), "", "J", false)
	} else {
		pdf.MultiCell(0, 6, tr("Actualmente no se puede generar el paz y salvo porque el funcionario cuenta con los siguientes elementos bajo su responsabilidad:"), "", "J", false)
		pdf.Ln(2)
		construirTablaElementos(pdf, tr, elementos)
	}

	pdf.Ln(8)
	pdf.MultiCell(0, 6, tr("Expedido en Bogotá D.C., a los "+fechaLargaEs(time.Now())+"."), "", "J", false)
	renderBloqueFirmas(pdf, tr, elaborador, responsableFirma)

	var buffer bytes.Buffer
	if err := pdf.Output(&buffer); err != nil {
		return "", errorCtrl.Error(funcion+" - pdf.Output", err, "500")
	}

	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

func construirParrafoPrincipal(tercero *models.PazYSalvoTerceroResponse, puedeGenerar bool) string {
	nombre := "la persona solicitante"
	if tercero != nil && strings.TrimSpace(tercero.NombreCompleto) != "" {
		nombre = tercero.NombreCompleto
	}

	tipoDocumento := "documento"
	if tercero != nil && strings.TrimSpace(tercero.TipoDocumento) != "" {
		tipoDocumento = tercero.TipoDocumento
	}

	numeroDocumento := ""
	if tercero != nil {
		numeroDocumento = tercero.NumeroDocumento
	}

	base := fmt.Sprintf(
		"Que una vez revisada la aplicación ARKA II Sistema Gestión Almacén e Inventarios, se evidenció que %s identificado con %s número %s",
		nombre,
		tipoDocumento,
		numeroDocumento,
	)

	if puedeGenerar {
		return base + " no tiene ningún elemento bajo su responsabilidad registrado en el aplicativo."
	}

	return base + " cuenta con elementos bajo su responsabilidad registrados en el aplicativo."
}

func construirTablaElementos(pdf *gofpdf.Fpdf, tr func(string) string, elementos []models.DetalleElementoPlaca) {
	headers := []string{"Placa", "Descripción", "Marca", "Serie"}
	widths := []float64{24, 82, 28, 32}

	imprimirCabeceraTabla(pdf, tr, headers, widths)

	pdf.SetFont("Times", "", 9)
	for _, elemento := range elementos {
		valores := []string{
			elemento.Placa,
			elemento.Nombre,
			elemento.Marca,
			elemento.Serie,
		}
		alturaFila := calcularAlturaFilaTabla(pdf, tr, widths, valores, 4.5)
		_, pageHeight := pdf.GetPageSize()
		if pdf.GetY()+alturaFila > pageHeight-20 {
			pdf.AddPage()
			imprimirCabeceraTabla(pdf, tr, headers, widths)
			pdf.SetFont("Times", "", 9)
		}

		dibujarFilaTabla(pdf, tr, widths, valores, 4.5, alturaFila)
	}
}

func fechaLargaEs(fecha time.Time) string {
	meses := map[time.Month]string{
		time.January:   "1 de enero de 2006",
		time.February:  "1 de febrero de 2006",
		time.March:     "1 de marzo de 2006",
		time.April:     "1 de abril de 2006",
		time.May:       "1 de mayo de 2006",
		time.June:      "1 de junio de 2006",
		time.July:      "1 de julio de 2006",
		time.August:    "1 de agosto de 2006",
		time.September: "1 de septiembre de 2006",
		time.October:   "1 de octubre de 2006",
		time.November:  "1 de noviembre de 2006",
		time.December:  "1 de diciembre de 2006",
	}
	layout, ok := meses[fecha.Month()]
	if !ok {
		layout = "2 de January de 2006"
	}
	return strings.ToLower(fecha.Format(layout))
}

func obtenerConfig(key, fallback string) string {
	if value, err := beego.AppConfig.String(key); err == nil && strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	return fallback
}

func sanitizarNombreArchivo(value string) string {
	replacer := strings.NewReplacer(" ", "_", ".", "", "/", "_", "\\", "_")
	return replacer.Replace(strings.TrimSpace(value))
}

func seleccionarSupervisorPorCargo(supervisores []supervisorContrato, cargoBuscado string) *supervisorContrato {
	var mejorCoincidencia *supervisorContrato

	for i := range supervisores {
		cargoSupervisor := strings.TrimSpace(supervisores[i].Cargo)
		if strings.EqualFold(cargoSupervisor, cargoBuscado) {
			if mejorCoincidencia == nil || supervisores[i].Id > mejorCoincidencia.Id {
				mejorCoincidencia = &supervisores[i]
			}
		}
	}

	if mejorCoincidencia != nil {
		return mejorCoincidencia
	}

	for i := range supervisores {
		cargoSupervisor := strings.ToLower(strings.TrimSpace(supervisores[i].Cargo))
		if strings.Contains(cargoSupervisor, strings.ToLower(cargoBuscado)) {
			if mejorCoincidencia == nil || supervisores[i].Id > mejorCoincidencia.Id {
				mejorCoincidencia = &supervisores[i]
			}
		}
	}

	return mejorCoincidencia
}

func obtenerTipoDocumento(identificacion *models.DatosIdentificacion) string {
	if identificacion != nil && identificacion.TipoDocumentoId != nil && strings.TrimSpace(identificacion.TipoDocumentoId.CodigoAbreviacion) != "" {
		return strings.TrimSpace(identificacion.TipoDocumentoId.CodigoAbreviacion)
	}
	return ""
}

func construirTextoElaborador(elaborador *pazYSalvoFirmante) string {
	if elaborador == nil {
		return ""
	}

	return strings.TrimSpace(elaborador.Nombre)
}

func obtenerTextoFirmante(firmante *pazYSalvoFirmante, campo, fallback string) string {
	if firmante == nil {
		return strings.TrimSpace(fallback)
	}

	switch campo {
	case "nombre":
		if strings.TrimSpace(firmante.Nombre) != "" {
			return strings.TrimSpace(firmante.Nombre)
		}
	case "cargo":
		if strings.TrimSpace(firmante.Cargo) != "" {
			return strings.TrimSpace(firmante.Cargo)
		}
	case "tipo_documento":
		if strings.TrimSpace(firmante.TipoDocumento) != "" {
			return strings.TrimSpace(firmante.TipoDocumento)
		}
	case "numero_documento":
		if strings.TrimSpace(firmante.NumeroDocumento) != "" {
			return strings.TrimSpace(firmante.NumeroDocumento)
		}
	}

	return strings.TrimSpace(fallback)
}

func seleccionarSupervisorVigentePorCargo(supervisores []supervisorContrato, cargoBuscado string, fechaReferencia time.Time) *supervisorContrato {
	cargoBuscado = normalizarTextoComparacion(cargoBuscado)
	for i := range supervisores {
		cargoSupervisor := normalizarTextoComparacion(supervisores[i].Cargo)
		if cargoSupervisor != cargoBuscado {
			continue
		}
		if !supervisorVigente(supervisores[i], fechaReferencia) {
			continue
		}
		return &supervisores[i]
	}

	for i := range supervisores {
		cargoSupervisor := normalizarTextoComparacion(supervisores[i].Cargo)
		if !strings.Contains(cargoSupervisor, cargoBuscado) {
			continue
		}
		if !supervisorVigente(supervisores[i], fechaReferencia) {
			continue
		}
		return &supervisores[i]
	}

	return nil
}

func supervisorVigente(supervisor supervisorContrato, fechaReferencia time.Time) bool {
	if !supervisor.FechaInicio.IsZero() && fechaReferencia.Before(supervisor.FechaInicio) {
		return false
	}
	if !supervisor.FechaFin.IsZero() && fechaReferencia.After(supervisor.FechaFin) {
		return false
	}
	return true
}

func normalizarTextoComparacion(valor string) string {
	valor = strings.ToUpper(strings.TrimSpace(valor))
	replacer := strings.NewReplacer(
		"Á", "A",
		"É", "E",
		"Í", "I",
		"Ó", "O",
		"Ú", "U",
		"Ü", "U",
	)
	return replacer.Replace(valor)
}

func renderHeaderPazYSalvo(pdf *gofpdf.Fpdf, tr func(string) string) {
	agregarLogoInstitucional(pdf)

	pdf.SetY(44)
	pdf.SetFont("Times", "BI", 13)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 6, tr("ALMACEN GENERAL DE LA UNIVERSIDAD DISTRITAL"), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 6, tr("FRANCISCO JOSE DE CALDAS"), "", 1, "C", false, 0, "")
	pdf.Ln(8)
	pdf.SetFont("Times", "BI", 13)
	pdf.CellFormat(0, 6, tr("PAZ Y SALVO"), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 6, tr("INVENTARIO"), "", 1, "C", false, 0, "")
	pdf.Ln(30)
}

func renderBloqueFirmas(pdf *gofpdf.Fpdf, tr func(string) string, elaborador, responsableFirma *pazYSalvoFirmante) {
	const firmaY = 200.0
	const textoElaboroOffset = 4.0

	if pdf.GetY() > firmaY-18 {
		pdf.AddPage()
	}

	pdf.Line(25, firmaY-5, 190, firmaY-5)
	pdf.SetY(firmaY)
	pdf.SetFont("Times", "B", 11)
	pdf.CellFormat(0, 6, tr(formatearNombreFirmante(obtenerTextoFirmante(responsableFirma, "nombre", ""))), "", 1, "C", false, 0, "")
	pdf.SetFont("Times", "", 11)
	pdf.CellFormat(0, 6, tr(obtenerTextoFirmante(responsableFirma, "cargo", "")), "", 1, "C", false, 0, "")

	if elaboradoPor := construirTextoElaborador(elaborador); elaboradoPor != "" {
		pdf.SetFont("Times", "", 8)
		pdf.SetY(pdf.GetY() + textoElaboroOffset)
		pdf.CellFormat(0, 4.2, tr("Elaboró: "+elaboradoPor), "", 1, "L", false, 0, "")
	}
}

func agregarLogoInstitucional(pdf *gofpdf.Fpdf) {
	logoPath := obtenerRutaLogoPazYSalvo()
	if _, err := os.Stat(logoPath); err != nil {
		return
	}

	logoWidth := 64.0
	posicionX := 18.0
	logoType := strings.TrimPrefix(strings.ToLower(filepath.Ext(logoPath)), ".")

	pdf.ImageOptions(logoPath, posicionX, 12, logoWidth, 0, false, gofpdf.ImageOptions{
		ImageType: logoType,
		ReadDpi:   true,
	}, 0, "")
}

func renderMarcaDeAguaPazYSalvo(pdf *gofpdf.Fpdf) {
	isotipoPath := obtenerRutaIsotipoArkaII()
	if _, err := os.Stat(isotipoPath); err != nil {
		return
	}

	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(isotipoPath)), ".")
	if ext == "" {
		return
	}

	pdf.SetAlpha(0.08, "Normal")
	pdf.ImageOptions(isotipoPath, 55, 90, 100, 0, false, gofpdf.ImageOptions{
		ImageType: ext,
		ReadDpi:   true,
	}, 0, "")
	pdf.SetAlpha(1.0, "Normal")
}

func renderFooterPazYSalvo(pdf *gofpdf.Fpdf) {
	pdf.SetY(-28)
	pdf.SetFont("Times", "", 8)
	pdf.CellFormat(95, 4, "PBX 57(1)3239300 Ext. 1621", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 4, "Linea de atencion gratuita", "", 1, "R", false, 0, "")
	pdf.CellFormat(95, 4, "Carrera 7 No. 40B-53 Piso 6, Bogota D.C. - Colombia", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 4, "01 800 091 44 10", "", 1, "R", false, 0, "")
	pdf.CellFormat(95, 4, "Acreditacion Institucional de Alta Calidad. Resolucion No. 23096 del 15 diciembre de 2016", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 4, "www.udistrital.edu.co", "", 1, "R", false, 0, "")
	pdf.CellFormat(95, 4, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 4, "almacen@udistrital.edu.co", "", 1, "R", false, 0, "")
}

func obtenerRutaLogoPazYSalvo() string {
	_, archivoActual, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	return filepath.Clean(filepath.Join(filepath.Dir(archivoActual), "..", "..", "assets", "logo_universidad_acreditacion.png"))
}

func obtenerRutaIsotipoArkaII() string {
	_, archivoActual, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	return filepath.Clean(filepath.Join(filepath.Dir(archivoActual), "..", "..", "assets", "isotipo_arkaII.png"))
}

func imprimirCabeceraTabla(pdf *gofpdf.Fpdf, tr func(string) string, headers []string, widths []float64) {
	pdf.SetFont("Times", "B", 10)
	for i, header := range headers {
		pdf.CellFormat(widths[i], 8, tr(header), "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
}

func calcularAlturaFilaTabla(pdf *gofpdf.Fpdf, tr func(string) string, widths []float64, valores []string, altoLinea float64) float64 {
	maxLineas := 1
	for i, valor := range valores {
		lineas := pdf.SplitLines([]byte(tr(valor)), widths[i]-2)
		if len(lineas) > maxLineas {
			maxLineas = len(lineas)
		}
	}

	return float64(maxLineas)*altoLinea + 2
}

func dibujarFilaTabla(pdf *gofpdf.Fpdf, tr func(string) string, widths []float64, valores []string, altoLinea, alturaFila float64) {
	inicioX := pdf.GetX()
	inicioY := pdf.GetY()
	posicionX := inicioX

	for i, valor := range valores {
		pdf.Rect(posicionX, inicioY, widths[i], alturaFila, "D")
		pdf.SetXY(posicionX+1, inicioY+1)
		pdf.MultiCell(widths[i]-2, altoLinea, tr(valor), "", "L", false)
		posicionX += widths[i]
		pdf.SetXY(posicionX, inicioY)
	}

	pdf.SetXY(inicioX, inicioY+alturaFila)
}

func formatearNombreFirmante(nombre string) string {
	partes := strings.Fields(strings.TrimSpace(nombre))
	if len(partes) == 3 {
		return strings.ToUpper(strings.Join([]string{partes[2], partes[0], partes[1]}, " "))
	}
	if len(partes) >= 4 {
		nombres := partes[len(partes)-2:]
		apellidos := partes[:len(partes)-2]
		return strings.ToUpper(strings.Join(append(nombres, apellidos...), " "))
	}
	return strings.ToUpper(strings.TrimSpace(nombre))
}
