package models

type Proveedor struct {
	Id                      int
	Tipopersona             string
	NumDocumento            string
	IdCiudadContacto        int
	Direccion               string
	Correo                  string
	Web                     string
	NomAsesor               string
	TelAsesor               string
	Descripcion             string
	PuntajeEvaluacion       int
	ClasificacionEvaluacion string
	Estado                  *Estado
	TipoCuentaBancaria      string
	NumCuentaBancaria       string
	IdEntidadBancaria       int
	FechaRegistro           string
	FechaUltimaModificacion string
	NomProveedor            string
	Anexorut                string
	Anexorup                string
	RegimenContributivo     string
}

type Estado struct {
	Id                   int
	ClaseParametro       string
	ValorParametro       string
	DescripcionParametro string
	Abreviatura          string
}
