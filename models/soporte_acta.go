package models

import "time"

type SoporteActa struct {
	Id                int           `orm:"column(id);pk;auto"`
	Consecutivo       string        `orm:"column(consecutivo)"`
	ProveedorId       int           `orm:"column(proveedor_id)"`
	FechaSoporte      time.Time     `orm:"column(fecha_soporte);type(date)"`
	ActaRecibidoId    *ActaRecibido `orm:"column(acta_recibido_id);rel(fk)"`
	Activo            bool          `orm:"column(activo)"`
	FechaCreacion     time.Time     `orm:"column(fecha_creacion);type(timestamp without time zone)"`
	FechaModificacion time.Time     `orm:"column(fecha_modificacion);type(timestamp without time zone)"`
}
