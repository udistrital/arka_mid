package models

import "time"

type ElementosMovimiento struct {
	Id                int         `orm:"column(id);pk;auto"`
	ElementoActaId    int         `orm:"column(elemento_acta_id)"`
	Unidad            float64     `orm:"column(unidad)"`
	ValorUnitario     float64     `orm:"column(valor_unitario)"`
	ValorTotal        float64     `orm:"column(valor_total)"`
	SaldoCantidad     float64     `orm:"column(saldo_cantidad)"`
	SaldoValor        float64     `orm:"column(saldo_valor)"`
	Activo            bool        `orm:"column(activo)"`
	FechaCreacion     time.Time   `orm:"auto_now_add;column(fecha_creacion);type(timestamp without time zone)"`
	FechaModificacion time.Time   `orm:"auto_now;column(fecha_modificacion);type(timestamp without time zone)"`
	MovimientoId      *Movimiento `orm:"column(movimiento_id);rel(fk)"`
}
