package models

import "time"

type BodegaConsumoSolicitud struct {
	Elementos *ElementoSolicitud
	Solicitud *Movimiento
}

type ElementoSolicitud struct {
	Cantidad           uint
	CantidadAprobada   uint
	ElementoCatalogoId *ElementoCatalogo
	Id                 uint
	Nombre             string
	SaldoCantidad      uint
	SaldoValor         uint
	Sede               *EspacioFisico
	Dependencia        *Dependencia
	Ubicacion          *AsignacionEspacioFisicoDependencia
}

type ElementoSolicitud_ struct {
	Cantidad           int
	Ubicacion          int
	ElementoCatalogoId int
	CantidadAprobada   int
}
type FormatoSolicitudBodega struct {
	Funcionario int
	Elementos   []ElementoSolicitud_
}

type DetalleSolicitudBodega struct {
	Movimiento
	Solicitante IdentificacionTercero
}

type Apertura struct {
	CantidadMinima     int
	CantidadMaxima     int
	ElementoCatalogoId int
	FechaCreacion      time.Time
	MetodoValoracion   int
	SaldoCantidad      float64
	SaldoValor         float64
	Unidad             float64
	ValorUnitario      float64
	ValorTotal         float64
}
