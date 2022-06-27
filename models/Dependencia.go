package models

import "time"

type Dependencia struct {
	Id                  int
	Nombre              string
	TelefonoDependencia string
	CorreoElectronico   string
	Activo              bool
	FechaCreacion       time.Time
	FechaModificacion   time.Time
}
