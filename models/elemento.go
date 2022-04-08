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

type PlantillaActa struct {
	Id                 int
	Nombre             string
	Marca              string
	Serie              string
	Cantidad           int
	UnidadMedida       int
	ValorUnitario      float64
	Subtotal           float64
	Descuento          float64
	PorcentajeIvaId    *int
	ValorIva           float64
	ValorTotal         float64
	SubgrupoCatalogoId *DetalleSubgrupo
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

type DetalleElemento_ struct {
	Elemento
	VidaUtil      float64
	ValorResidual float64
}

type DetalleElemento__ struct {
	DetalleElemento
	VidaUtil      float64
	ValorResidual float64
}
