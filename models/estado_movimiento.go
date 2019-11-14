package models

import "time"

type EstadoMovimiento struct {
	Id                int       `orm:"column(id);pk;auto"`
	Nombre            string    `orm:"column(nombre)"`
	Activo            bool      `orm:"column(activo)"`
	FechaCreacion     time.Time `orm:"auto_now_add;column(fecha_creacion);type(timestamp without time zone)"`
	FechaModificacion time.Time `orm:"auto_now;column(fecha_modificacion);type(timestamp without time zone)"`
	Descripcion       string    `orm:"column(descripcion)"`
}
