package models

import "time"

type AsignacionEspacioFisicoDependencia struct {
	Id               int
	EspacioFisicoId  *EspacioFisico
	DependenciaId    *Dependencia
	Activo           bool
	FechaInicio      time.Time
	FechaFin         time.Time
	DocumentoSoporte int
}
