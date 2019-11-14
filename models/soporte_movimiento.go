package models

import "time"

type SoporteMovimiento struct {
	Id                int         `orm:"column(id);pk;auto"`
	DocumentoId       int         `orm:"column(documento_id)"`
	Activo            bool        `orm:"column(activo)"`
	FechaCreacion     time.Time   `orm:"auto_now_add;column(fecha_creacion);type(timestamp without time zone)"`
	FechaModificacion time.Time   `orm:"auto_now;column(fecha_modificacion);type(timestamp without time zone)"`
	MovimientoId      *Movimiento `orm:"column(movimiento_id);rel(fk)"`
}
