package models

import "time"

type EspacioFisico struct {
	Id                int
	Nombre            string
	CodigoAbreviacion string
	Activo            bool
	TipoEspacio       *TipoEspacio
	FechaCreacion     time.Time
	FechaModificacion time.Time
}
