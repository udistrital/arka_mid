package models

import "time"

type ActaRecibido struct {
	Id                int       `orm:"column(id);pk;auto"`
	UbicacionId       int       `orm:"column(ubicacion_id)"`
	FechaVistoBueno   time.Time `orm:"column(fecha_visto_bueno);type(date);null"`
	RevisorId         int       `orm:"column(revisor_id)"`
	Observaciones     string    `orm:"column(observaciones);null"`
	Activo            bool      `orm:"column(activo)"`
	FechaCreacion     time.Time `orm:"column(fecha_creacion);type(timestamp without time zone)"`
	FechaModificacion time.Time `orm:"column(fecha_modificacion);type(timestamp without time zone)"`
}
