package models

// "time"

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
