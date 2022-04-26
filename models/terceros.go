// Modelos asociados al CRUD API de Terceros

package models

import "time"

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

type GrupoInfoComplementaria struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
}

type InfoComplementariaTercero struct {
	Id                       int
	TerceroId                *Tercero
	InfoComplementariaId     *InfoComplementaria
	Dato                     string
	Activo                   bool
	InfoCompleTerceroPadreId *InfoComplementariaTercero
}

type InfoComplementaria struct {
	Id                        int
	Nombre                    string
	CodigoAbreviacion         string
	Activo                    bool
	TipoDeDato                string
	GrupoInfoComplementariaId *GrupoInfoComplementaria
}

type SeguridadSocialTercero struct {
	Id                     int
	TerceroId              *Tercero
	TerceroEntidadId       *Tercero
	Activo                 bool
	FechaInicioVinculacion *time.Time
	FechaFinVinculacion    *time.Time
}

type TerceroFamiliar struct {
	Id                int
	TerceroId         *Tercero
	TerceroFamiliarId *Tercero
	TipoParentescoId  *TipoParentesco
	CodigoAbreviacion string
	Activo            bool
}

type TipoParentesco struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
}

type DetalleTercero struct {
	Tercero         *Tercero
	TipoVinculacion int
	DependenciaId   int
	Identificacion  *DatosIdentificacion
}
type DetalleFuncionario struct {
	Tercero []*DetalleTercero
	Correo  []*InfoComplementariaTercero
	Cargo   []*Parametro
}

type InfoTercero struct {
	Tercero        *Tercero
	Identificacion *DatosIdentificacion
}

type IdentificacionTercero struct {
	Id             int
	Numero         string
	NombreCompleto string
}
