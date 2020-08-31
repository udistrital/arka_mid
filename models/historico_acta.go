package models

import "time"

type HistoricoActa struct {
	Id                int           `orm:"column(id);pk;auto"`
	ActaRecibidoId    *ActaRecibido `orm:"column(acta_recibido_id);rel(fk)"`
	EstadoActaId      *EstadoActa   `orm:"column(estado_acta_id);rel(fk)"`
	Activo            bool          `orm:"column(activo)"`
	FechaCreacion     time.Time     `orm:"column(fecha_creacion);type(timestamp without time zone)"`
	FechaModificacion time.Time     `orm:"column(fecha_modificacion);type(timestamp without time zone)"`
}
