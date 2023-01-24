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

type TransaccionActaRecibido struct {
	ActaRecibido *ActaRecibido
	UltimoEstado *HistoricoActa
	SoportesActa *[]SoporteActa
	Elementos    []*Elemento
}

type SoporteActa struct {
	Id                int
	Consecutivo       string
	DocumentoId       int
	FechaSoporte      time.Time
	ActaRecibidoId    *ActaRecibido
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
}

type HistoricoActa struct {
	Id                int
	ProveedorId       int
	UbicacionId       int
	RevisorId         int
	PersonaAsignadaId int
	Observaciones     string
	FechaVistoBueno   time.Time
	ActaRecibidoId    *ActaRecibido
	EstadoActaId      *EstadoActa
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
}

type EstadoActa struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
	NumeroOrden       float64
	FechaCreacion     time.Time
	FechaModificacion time.Time
}

type EstadoElemento struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
	NumeroOrden       float64
	FechaCreacion     time.Time
	FechaModificacion time.Time
}

type Elemento struct {
	Id                 int
	Nombre             string
	Cantidad           int
	Marca              string
	Serie              string
	UnidadMedida       int
	ValorUnitario      float64
	Subtotal           float64
	Descuento          float64
	ValorTotal         float64
	PorcentajeIvaId    int
	ValorIva           float64
	ValorFinal         float64
	SubgrupoCatalogoId int
	TipoBienId         int
	EstadoElementoId   *EstadoElemento
	ActaRecibidoId     *ActaRecibido
	Placa              string
	Activo             bool
	FechaCreacion      time.Time
	FechaModificacion  time.Time
}
