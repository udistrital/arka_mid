package models

import "time"

type Movimiento struct {
	Id                      int                    `orm:"column(id);pk;auto"`
	Observacion             string                 `orm:"column(observacion);null"`
	Detalle                 string                 `orm:"column(detalle);type(json)"`
	FechaCreacion           time.Time              `orm:"auto_now_add;column(fecha_creacion);type(timestamp without time zone)"`
	FechaModificacion       time.Time              `orm:"auto_now;column(fecha_modificacion);type(timestamp without time zone)"`
	Activo                  bool                   `orm:"column(activo)"`
	MovimientoPadreId       *Movimiento            `orm:"column(movimiento_padre_id);rel(fk)"`
	FormatoTipoMovimientoId *FormatoTipoMovimiento `orm:"column(formato_tipo_movimiento_id);rel(fk)"`
	EstadoMovimientoId      *EstadoMovimiento      `orm:"column(estado_movimiento_id);rel(fk)"`
	IdTipoMovimiento        int
	SoporteMovimientoId     int
}
