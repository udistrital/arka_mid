package models

import "time"

type SoporteActa struct {
	Id                int
	Consecutivo       string
	DocumentoId       int
	FechaSoporte      time.Time
	ActaRecibidoId    *ActaRecibido
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
}
