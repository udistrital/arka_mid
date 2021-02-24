// Modelos asociados al CRUD API de Terceros

package models

// "time"

type DatosIdentificacion struct {
	Id                 int
	TipoDocumentoId    *TipoDocumento
	TerceroId          *Tercero
	Numero             string
	DigitoVerificacion int
	CiudadExpedicion   int
	FechaExpedicion    string
	Activo             bool
	DocumentoSoporte   int
	FechaCreacion      string
	FechaModificacion  string
}

type TipoContribuyente struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
	FechaCreacion     string
	FechaModificacion string
}

type Tercero struct {
	Id                  int
	NombreCompleto      string
	PrimerNombre        string
	SegundoNombre       string
	PrimerApellido      string
	SegundoApellido     string
	LugarOrigen         int
	FechaNacimiento     string
	Activo              bool
	TipoContribuyenteId *TipoContribuyente
	FechaCreacion       string
	FechaModificacion   string
	UsuarioWSO2         string
}

type TerceroTipoTercero struct {
	Id                int
	TerceroId         *Tercero
	TipoTerceroId     *TipoTercero
	Activo            bool
	FechaCreacion     string
	FechaModificacion string
}

type TipoDocumento struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
	FechaCreacion     string
	FechaModificacion string
	NumeroOrden       int
}

type TipoTercero struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
	FechaCreacion     string
	FechaModificacion string
}

type Vinculacion struct {
	Id                     int
	TerceroPrincipalId     *Tercero
	TerceroRelacionadoId   *Tercero
	TipoVinculacionId      int
	CargoId                int
	DependenciaId          int
	Soporte                int
	PeriodoId              int
	FechaInicioVinculacion string
	FechaFinVinculacion    string
	Activo                 bool
	FechaCreacion          string
	FechaModificacion      string
}
