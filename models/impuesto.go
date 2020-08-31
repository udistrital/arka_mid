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

type VigenciaImpuesto struct {
	Id                   int
	Activo               bool
	Tarifa               int
	PorcentajeAplicacion int
	ImpuestoId           Impuesto
}
