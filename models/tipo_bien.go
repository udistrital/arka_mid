package models

import "time"

type TipoBien struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Orden             float64
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
}
