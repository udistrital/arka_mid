package models

import "time"

type EspacioFisico struct {
	Id                int
	Nombre            string
	CodigoAbreviacion string
	Activo            bool
	TipoTerrenoId     int
	TipoEdificacionId int
	TipoEspacio       *TipoEspacio
	FechaCreacion     time.Time
	FechaModificacion time.Time
}
