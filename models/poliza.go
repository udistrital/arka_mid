package models

import "time"

type Campo struct {
	Id                int
	Nombre            string
	Sigla             string
	Descripcion       string
	Metadato          string
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
}

type Elemento_campo struct {
	Id                int
	ElementoId        *Elemento
	CampoId           *Campo
	Valor             string
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
}
