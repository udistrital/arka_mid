package models

import "time"

type HistoricoActa struct {
	Id                int
	ProveedorId       int
	UbicacionId       int
	RevisorId         int
	PersonaAsignadaId int
	Observaciones     string
	FechaVistoBueno   time.Time
	ActaRecibidoId    *ActaRecibido
	EstadoActaId      *EstadoActa
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
}
