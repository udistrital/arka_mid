package models

import "time"

type Elemento struct {
	Id                 int             `orm:"column(id);pk;auto"`
	Nombre             string          `orm:"column(nombre)"`
	Cantidad           int             `orm:"column(cantidad)"`
	Marca              string          `orm:"column(marca);null"`
	Serie              string          `orm:"column(serie);null"`
	UnidadMedida       int             `orm:"column(unidad_medida)"`
	ValorUnitario      float64         `orm:"column(valor_unitario)"`
	Subtotal           float64         `orm:"column(subtotal);null"`
	Descuento          float64         `orm:"column(descuento);null"`
	ValorTotal         float64         `orm:"column(valor_total);null"`
	PorcentajeIvaId    int             `orm:"column(porcentaje_iva_id)"`
	ValorIva           float64         `orm:"column(valor_iva);null"`
	ValorFinal         float64         `orm:"column(valor_final);null"`
	SubgrupoCatalogoId int             `orm:"column(subgrupo_catalogo_id)"`
	Verificado         bool            `orm:"column(verificado)"`
	TipoBienId         *TipoBien       `orm:"column(tipo_bien_id);rel(fk)"`
	EstadoElementoId   *EstadoElemento `orm:"column(estado_elemento_id);rel(fk)"`
	EspacioFisicoId    int             `orm:"column(espacio_fisico_id)"`
	SoporteActaId      *SoporteActa    `orm:"column(soporte_acta_id);rel(fk)"`
	Placa              string          `orm:"column(placa);null"`
	Activo             bool            `orm:"column(activo)"`
	FechaCreacion      time.Time       `orm:"column(fecha_creacion);type(timestamp without time zone)"`
	FechaModificacion  time.Time       `orm:"column(fecha_modificacion);type(timestamp without time zone)"`
}
