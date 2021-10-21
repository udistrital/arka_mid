package models

import "time"

type MovimientoProcesoExterno struct {
	Id                       int
	TipoMovimientoId         *TipoMovimiento
	ProcesoExterno           int64
	MovimientoProcesoExterno int
	Activo                   bool
	FechaCreacion            time.Time
	FechaModificacion        time.Time
}

type TipoMovimiento struct {
	Id                int
	Nombre            string
	Descripcion       string
	Acronimo          string
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
	Parametros        string
}
