package models

import "time"

type FormatoTipoMovimiento struct {
	Id                int       `orm:"column(id);pk;auto"`
	Nombre            string    `orm:"column(nombre)"`
	Formato           string    `orm:"column(formato);type(json)"`
	Descripcion       string    `orm:"column(descripcion);null"`
	CodigoAbreviacion string    `orm:"column(codigo_abreviacion);null"`
	NumeroOrden       float64   `orm:"column(numero_orden);null"`
	FechaCreacion     time.Time `orm:"auto_now_add;column(fecha_creacion);type(timestamp without time zone)"`
	FechaModificacion time.Time `orm:"auto_now;column(fecha_modificacion);type(timestamp without time zone)"`
	Activo            bool      `orm:"column(activo)"`
}
