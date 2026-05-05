package models

type ReporteFechasRequest struct {
	FechaInicial string `json:"fecha_inicial"`
	FechaFinal   string `json:"fecha_final"`
}

type ReporteExcelBase64Response struct {
	ArchivoBase64 string `json:"archivo_base64"`
	NombreArchivo string `json:"nombre_archivo"`
	TipoArchivo   string `json:"tipo_archivo"`
}
