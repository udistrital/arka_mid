package models

import (
	"time"
)

type ActaRecibido struct {
	Id                int
	Activo            bool
	TipoActaId        *TipoActa
	UnidadEjecutoraId int
	FechaCreacion     time.Time
	FechaModificacion time.Time
}

type ActaRecibidoUbicacion struct {
	Id                int
	UbicacionId       *AsignacionEspacioFisicoDependencia
	FechaVistoBueno   time.Time
	RevisorId         int
	Observaciones     string
	Activo            bool
	EstadoActaId      *EstadoActa
	FechaCreacion     time.Time
	FechaModificacion time.Time
}

type TipoActa struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
}

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
