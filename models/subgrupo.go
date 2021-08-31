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
	Id          int
	Nombre      string
	Descripcion string
	Activo      bool
	Codigo      string
	TipoNivelId *TipoNivel
}

type TipoNivel struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Orden             float64
	Activo            bool
	FechaCreacion     string
	FechaModificacion string
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
	CuentasAsociadas  []CuentaSubgrupo
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

type DetalleSubgrupo struct {
	Id           int
	Depreciacion bool
	Valorizacion bool
	Deterioro    bool
	Activo       bool
	SubgrupoId   *Subgrupo
	TipoBienId   *TipoBien
}
