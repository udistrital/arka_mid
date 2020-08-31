package models

import "time"

type ParametrosGobierno struct {
	Id                   int
	Activo               bool
	Tarifa               int
	PorcentajeAplicacion int
	BaseUvt              int
	BasePesos            int
	InicioVigencia       time.Time
	FinVigencia          time.Time
	Decreto              string
	FechaCreacion        time.Time
	FechaModificacion    time.Time
	ImpuestoId           *Impuesto
}
