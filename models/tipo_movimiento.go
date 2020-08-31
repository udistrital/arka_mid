package models

import "time"

type TipoMovimiento struct {
	Id                int       `orm:"column(id);pk;auto"`
	Nombre            string    `orm:"column(nombre)"`
	Descripcion       string    `orm:"column(descripcion);null"`
	Acronimo          string    `orm:"column(acronimo)"`
	Activo            bool      `orm:"column(activo);null"`
	FechaCreacion     time.Time `orm:"column(fecha_creacion);type(date);null"`
	FechaModificacion time.Time `orm:"column(fecha_modificacion);type(date);null"`
	Parametros        string    `orm:"column(parametros);type(json);null"`
}
