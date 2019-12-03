package models

import "time"

type AsignacionEspacioFisicoDependencia struct {
	Id          		int
	Estado      		string
	FechaInicio 		time.Time
	FechaFin			time.Time
	EspacioFisicoId 	*EspacioFisico
	DependenciaId		*Dependencia
}
