package actaRecibido

import "time"

type Impuesto struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
}

type VigenciaImpuesto struct {
	Id                   int
	Activo               bool
	Tarifa               int64
	PorcentajeAplicacion int
	ImpuestoId           Impuesto
}

type Imp struct {
	PorcentajeAplicacion int
	Tarifa               int
	BasePesos            int
	BaseUvt              int
	CodigoAbreviacion    string
}

type Unidad struct {
	Id          int
	Unidad      string
	Tipo        string
	Descripcion string
	Estado      bool
}

type Subgrupo struct {
	Id                int
	Nombre            string
	Descripcion       string
	Activo            bool
	Codigo            string
	Estado            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
}
