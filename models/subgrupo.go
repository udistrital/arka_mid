package models

import "time"

type Subgrupo struct {
	Id                int
	Nombre            string
	Descripcion       string
	FechaCreacion     time.Time
	FechaModificacion time.Time
	Activo            bool
	Codigo            int
}

type SubgrupoTransaccion struct {
	data     *Subgrupo
	children []*Subgrupo
}
