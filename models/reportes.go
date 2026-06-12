package models

type ReporteFechasRequest struct {
	FechaInicial string `json:"fecha_inicial"`
	FechaFinal   string `json:"fecha_final"`
}

type PazYSalvoRequest struct {
	Usuario          string `json:"usuario"`
	ElaboroTerceroId int    `json:"elaboro_tercero_id"`
	NumeroDocumento  string `json:"numero_documento"`
}

type ReporteExcelBase64Response struct {
	ArchivoBase64 string `json:"archivo_base64"`
	NombreArchivo string `json:"nombre_archivo"`
	TipoArchivo   string `json:"tipo_archivo"`
}

type PazYSalvoTerceroResponse struct {
	Id              int    `json:"id"`
	NombreCompleto  string `json:"nombre_completo"`
	NumeroDocumento string `json:"numero_documento"`
	TipoDocumento   string `json:"tipo_documento"`
	Cargo           string `json:"cargo"`
}

type PazYSalvoResponse struct {
	ArchivoBase64         string                    `json:"archivo_base64"`
	NombreArchivo         string                    `json:"nombre_archivo"`
	TipoArchivo           string                    `json:"tipo_archivo"`
	Mensaje               string                    `json:"mensaje"`
	PuedeGenerarPazYSalvo bool                      `json:"puede_generar_paz_y_salvo"`
	Tercero               *PazYSalvoTerceroResponse `json:"tercero"`
	Elementos             []DetalleElementoPlaca    `json:"elementos"`
}

type ReporteDetalleEntradaResponse struct {
	ElementoNombre            string  `json:"ElementoNombre"`
	ElementoValorFinal        float64 `json:"ElementoValorFinal"`
	SalidaFuncionarioAsignado string  `json:"SalidaFuncionarioAsignado"`
	CuentaDebitoEntrada       string  `json:"CuentaDebitoEntrada"`
	CuentaCreditoEntrada      string  `json:"CuentaCreditoEntrada"`
}

type ReporteDetalleSalidaResponse struct {
	ElementoNombre            string  `json:"ElementoNombre"`
	ElementoValorFinal        float64 `json:"ElementoValorFinal"`
	SalidaFuncionarioAsignado string  `json:"SalidaFuncionarioAsignado"`
	CuentaDebitoSalida        string  `json:"CuentaDebitoSalida"`
	CuentaCreditoSalida       string  `json:"CuentaCreditoSalida"`
}
