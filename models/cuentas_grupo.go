package models

import "time"

type CuentaSubgrupo struct {
	Id                  int
	Activo              bool
	CuentaCreditoId     string
	CuentaDebitoId      string
	SubtipoMovimientoId int
	FechaCreacion       string
	FechaModificacion   string
	SubgrupoId          *Subgrupo
}

type CuentasGrupoModelo struct {
	Id                  int             `orm:"column(id);pk;auto"`
	CuentaCreditoId     int             `orm:"column(cuenta_credito_id)"`
	CuentaDebitoId      int             `orm:"column(cuenta_debito_id)"`
	SubtipoMovimientoId int             `orm:"column(subtipo_movimiento_id)"`
	FechaCreacion       time.Time       `orm:"auto_now;column(fecha_creacion);type(date)"`
	FechaModificacion   time.Time       `orm:"auto_now;column(fecha_modificacion);type(date)"`
	Activo              bool            `orm:"column(activo)"`
	SubgrupoId          *SubgrupoModelo `orm:"column(subgrupo_id);rel(fk)"`
}

type CuentasGrupoMovimiento struct {
	Id              int
	CuentaCreditoId string
	CuentaDebitoId  string
	//SubtipoMovimientoId *TipoMovimiento
	SubtipoMovimientoId int
	FechaCreacion       time.Time
	FechaModificacion   time.Time
	Activo              bool
	SubgrupoId          *SubgrupoModelo
}

type CuentasGrupoTransaccion struct {
	Id                  int
	CuentaCreditoId     *CuentaContable
	CuentaDebitoId      *CuentaContable
	SubtipoMovimientoId int
	FechaCreacion       time.Time
	FechaModificacion   time.Time
	Activo              bool
	SubgrupoId          *Subgrupo
}

type Movimientos_Kronos struct {
	Id                int
	Nombre            string
	Descripcion       string
	Acronimo          string
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
	Parametros        string
}

type TransaccionMovimientos struct {
	ConsecutivoId    int
	Etiquetas        string
	Descripcion      string
	FechaTransaccion time.Time
	Activo           bool
	Movimientos      []MovimientoTransaccion `json:"movimientos"`
}

type MovimientoTransaccion struct {
	TerceroId        int
	CuentaId         string
	NombreCuenta     string
	TipoMovimientoId int
	Valor            float64
	Descripcion      string
	Activo           bool
}
