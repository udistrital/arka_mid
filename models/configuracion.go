package models

type PerfilXMenuOpcion struct {
	Id         int
	Nombre     string
	Aplicacion *Aplicacion
}

type Aplicacion struct {
	Id          int
	Nombre      string
	Descripcion string
	Dominio     string
	Estado      bool
	Alias       string
	EstiloIcono string
}
