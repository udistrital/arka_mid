package models

type EspacioFisico struct {
	Id          int
	Estado      string
	TipoEspacio *TipoEspacio
	Nombre      string
	Codigo      string
}
