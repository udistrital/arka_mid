package models

type Sede struct {
	Activo            bool
	CodigoAbreviacion string
	Descripcion       string
	Id                int
	Nombre            string
	TipoEspacio       *TipoEspacio
}
