package models

import "time"

type EntradaElemento struct {
	Id                  int          `orm:"column(id);pk;auto"`
	Solicitante         int          `orm:"column(solicitante);null"`
	Observacion         string       `orm:"column(observacion);null"`
	Importacion         bool         `orm:"column(importacion);null"`
	FechaCreacion       time.Time    `orm:"auto_now;column(fecha_creacion);type(timestamp without time zone)"`
	FechaModificacion   time.Time    `orm:"auto_now;column(fecha_modificacion);type(timestamp without time zone)"`
	Activo              bool         `orm:"column(activo)"`
	TipoEntradaId       *TipoEntrada `orm:"column(tipo_entrada_id);rel(fk)"`
	ActaRecibidoId      int          `orm:"column(acta_recibido_id)"`
	ContratoId          int          `orm:"column(contrato_id);null"`
	ElementoId          int          `orm:"column(elemento_id);null"`
	DocumentoContableId int          `orm:"column(documento_contable_id)"`
	Consecutivo         string       `orm:"column(consecutivo)"`
	Vigencia            string       `orm:"column(vigencia)"`
}
