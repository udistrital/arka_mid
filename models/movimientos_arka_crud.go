package models

import "time"

type Movimiento struct {
	Id                      int
	Observacion             string
	Detalle                 string
	FechaCreacion           time.Time
	FechaModificacion       time.Time
	Activo                  bool
	MovimientoPadreId       *Movimiento
	FormatoTipoMovimientoId *FormatoTipoMovimiento
	EstadoMovimientoId      *EstadoMovimiento
	SoporteMovimientoId     int
}

type ElementosMovimiento struct {
	Id                int
	ElementoActaId    int
	Unidad            float64
	ValorUnitario     float64
	ValorTotal        float64
	SaldoCantidad     float64
	SaldoValor        float64
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
	MovimientoId      *Movimiento
}

type FormatoTipoMovimiento struct {
	Id                int
	Nombre            string
	Formato           string
	Descripcion       string
	CodigoAbreviacion string
	NumeroOrden       float64
	FechaCreacion     time.Time
	FechaModificacion time.Time
	Activo            bool
}

type EstadoMovimiento struct {
	Id                int
	Nombre            string
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
	Descripcion       string
}

type SoporteMovimiento struct {
	Id                int
	DocumentoId       int
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
	MovimientoId      *Movimiento
}
