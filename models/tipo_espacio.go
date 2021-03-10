package models

import "time"

type TipoEspacio struct {
	Id                int
	Nombre            string
	CodigoAbreviacion string
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
}
