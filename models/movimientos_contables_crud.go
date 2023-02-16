package models

import "time"

type TransaccionMovimientos struct {
	ConsecutivoId    int
	Etiquetas        string
	Descripcion      string
	FechaTransaccion time.Time
	Activo           bool
	Movimientos      []*MovimientoTransaccion
}

type MovimientoTransaccion struct {
	TerceroId        *int
	CuentaId         string
	NombreCuenta     string
	TipoMovimientoId int
	Valor            float64
	Descripcion      string
	Activo           bool
}

type DetalleMovimientoTransaccion struct {
	TerceroId        string
	CuentaId         *CuentaContable
	TipoMovimientoId int
	Valor            float64
	Descripcion      string
}

type DetalleTrContable struct {
	Movimientos      []*DetalleMovimientoTransaccion
	ConsecutivoId    int
	Etiquetas        string
	Descripcion      string
	FechaTransaccion time.Time
}
