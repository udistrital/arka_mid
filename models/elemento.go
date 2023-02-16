package models

import "time"

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
	TipoBienId         *TipoBien
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

type DetalleElementoSalida struct {
	Cantidad           int
	ElementoActaId     int
	Id                 int
	Marca              string
	Nombre             string
	Placa              string
	Serie              string
	SubgrupoCatalogoId *DetalleSubgrupo
	TipoBienId         *TipoBien
	ValorUnitario      float64
	ValorResidual      float64
	ValorTotal         float64
	VidaUtil           float64
}
