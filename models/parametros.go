package models

import "time"

type Parametro struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
	NumeroOrden       float64
	TipoParametroId   *TipoParametro
	ParametroPadreId  *Parametro
}

type TipoParametro struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
	NumeroOrden       float64
	FechaCreacion     string
	FechaModificacion string
	AreaTipoId        *AreaTipo
}

type AreaTipo struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
	NumeroOrden       float64
}

type ParametroPeriodo struct {
	Id          int
	ParametroId *Parametro
	PeriodoId   *Periodo
	Valor       string
	Activo      bool
}

type Periodo struct {
	Id                int
	Nombre            string
	Descripcion       string
	Year              float64
	Ciclo             string
	CodigoAbreviacion string
	Activo            bool
	AplicacionId      int
	InicioVigencia    time.Time
	FinVigencia       time.Time
}
