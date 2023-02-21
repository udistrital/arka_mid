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

type EspacioFisico struct {
	Id                  int
	Nombre              string
	CodigoAbreviacion   string
	Activo              bool
	Descripcion         string
	TipoTerrenoId       int
	TipoEdificacionId   int
	TipoEspacioFisicoId *TipoEspacioFisico
}

type TipoEspacioFisico struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
}

type DetalleSedeDependencia struct {
	Sede        *EspacioFisico
	Dependencia *Dependencia
	Ubicacion   *AsignacionEspacioFisicoDependencia
}

type Dependencia struct {
	Id                  int
	Nombre              string
	TelefonoDependencia string
	CorreoElectronico   string
}

type EspacioFisicoCampo struct {
	Id            int
	Valor         string
	EspacioFisico *EspacioFisico
	Campo         *Campo
}
