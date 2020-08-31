package models

import "time"

type SubgrupoModelo struct {
	Id                int
	Nombre            string
	Descripcion       string
	FechaCreacion     time.Time
	FechaModificacion time.Time
	Activo            bool
	Codigo            int
}
type Subgrupo struct {
	Id                int
	Nombre            string
	Descripcion       string
	FechaCreacion     time.Time
	FechaModificacion time.Time
	Activo            bool
	Codigo            int
	TipoBienId        *TipoBien
}

type SubgrupoTransaccion struct {
	data     *Subgrupo
	children []*Subgrupo
}

//SubgrupoCuentas define la estructura requerida para devolver las cuentas asociadas a un subgrupo especifico
type SubgrupoCuentas struct {
	Id                int
	Nombre            string
	Descripcion       string
	FechaCreacion     time.Time
	FechaModificacion time.Time
	Activo            bool
	Codigo            int
	CuentasAsociadas  []CuentasGrupo
}
type SubgrupoCuentasModelo struct {
	Id                int
	Nombre            string
	Descripcion       string
	FechaCreacion     time.Time
	FechaModificacion time.Time
	Activo            bool
	Codigo            int
	CuentasAsociadas  []CuentasGrupoModelo
}
type SubgrupoCuentasMovimiento struct {
	Id                int
	Nombre            string
	Descripcion       string
	FechaCreacion     time.Time
	FechaModificacion time.Time
	Activo            bool
	Codigo            int
	CuentasAsociadas  []CuentasGrupoMovimiento
}
