package models

import "time"

type Elemento struct {
	Id                 int
	Nombre             string
	Cantidad           int
	Marca              string
	Serie              string
	UnidadMedida       int
	ValorUnitario      float64
	Subtotal           float64
	Descuento          float64
	ValorTotal         float64
	PorcentajeIvaId    int
	ValorIva           float64
	ValorFinal         float64
	SubgrupoCatalogoId int
	EstadoElementoId   *EstadoElemento
	ActaRecibidoId     *ActaRecibido
	Placa              string
	Activo             bool
	FechaCreacion      time.Time
	FechaModificacion  time.Time
}

type DetalleElemento struct {
	Id                 int
	Nombre             string
	Cantidad           int
	Marca              string
	Serie              string
	UnidadMedida       int
	ValorUnitario      float64
	Subtotal           float64
	Descuento          float64
	ValorTotal         float64
	PorcentajeIvaId    int
	ValorIva           float64
	ValorFinal         float64
	SubgrupoCatalogoId *DetalleSubgrupo
	EstadoElementoId   *EstadoElemento
	ActaRecibidoId     *ActaRecibido
	Placa              string
	Activo             bool
	FechaCreacion      time.Time
	FechaModificacion  time.Time
}
