package models

type InformacionContrato struct {
	Contrato struct {
		FechaSuscripcion string `json:"fecha_suscripcion"`
		Justificacion    string `json:"justificacion"`
		TipoContrato     string `json:"tipo_contrato"`
		UnidadEjecucion  string `json:"unidad_ejecucion"`
		Vigencia         string `json:"vigencia"`
		OrdenadorGasto   struct {
			Id           string `json:"id"`
			RolOrdenador string `json:"rol_ordenador"`
			Nombre       string `json:"nombre"`
		} `json:"ordenador_gasto"`

		DescripcionFormaPago   string `json:"descripcion_forma_pago"`
		FechaRegistro          string `json:"fecha_registro"`
		Observaciones          string `json:"observaciones"`
		ObjetoContrato         string `json:"objeto_contrato"`
		Contratista            string `json:"contratista"`
		NumeroContratoSuscrito string `json:"numero_contrato_suscrito"`
		Supervisor             struct {
			Nombre                  string `json:"nombre"`
			Id                      string `json:"id"`
			DocumentoIdentificacion string `json:"documento_identificacion"`
			Cargo                   string `json:"cargo"`
		} `json:"supervisor"`
		LugarEjecucion  string `json:"lugar_ejecucion"`
		Actividades     string `json:"actividades"`
		UnidadEjecutora string `json:"unidad_ejecutora"`
		NumeroContrato  string `json:"numero_contrato"`
		PlazoEjecucion  string `json:"plazo_ejecucion"`
		ValorContrato   string `json:"valor_contrato"`
		Rubro           string `json:"rubro"`
	} `json:"contrato"`
}

type TipoContrato struct {
	Id           int
	TipoContrato string
}
