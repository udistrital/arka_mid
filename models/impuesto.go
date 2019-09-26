package models

import "time"

type Impuesto struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
	NumeroOrden       int
	FechaCreacion     time.Time
	FechaModificacion time.Time
	TipoImpuestoId    *TipoImpuesto
}
