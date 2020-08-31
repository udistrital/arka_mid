package models

import "time"

type MovimientoProcesoExterno struct {
	Id                       int             `orm:"column(id);auto"`
	TipoMovimientoId         *TipoMovimiento `orm:"column(tipo_movimiento_id);rel(fk)"`
	ProcesoExterno           int64           `orm:"column(proceso_externo)"`
	MovimientoProcesoExterno int             `orm:"column(movimiento_proceso_externo);null"`
	Activo                   bool            `orm:"column(activo);null"`
	FechaCreacion            time.Time       `orm:"auto_now_add;column(fecha_creacion);type(date)";null`
	FechaModificacion        time.Time       `orm:"auto_now;column(fecha_modificacion);type(date)";null`
}
