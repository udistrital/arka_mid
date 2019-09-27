package models

import "time"

type Contrato struct {
	FechaSuscripcion       time.Time
	Justificacion          string
	TipoContrato           int
	UnidadEjecucion        int
	OrdenadorGasto         *OrdenadorGasto
	DescripcionFormaPago   string
	FechaRegistro          time.Time
	Observaciones          string
	ObjetoContrato         string
	Contratista            int
	NumeroContratoSuscrito int
	Supervisor             *Supervisor
	LugarEjecucion         int
	Rubro                  string
	Actividades            string
	UnidadEjecutora        int
	NumeroContrato         int
	PlazoEjecucion         int
	ValorContrato          int
}
