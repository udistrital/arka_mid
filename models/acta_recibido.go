package models

import (
	"bytes"
	"encoding/gob"
	"time"
)

type ActaRecibido struct {
	Id                int       `orm:"column(id);pk;auto"`
	UbicacionId       int       `orm:"column(ubicacion_id)"`
	FechaVistoBueno   time.Time `orm:"column(fecha_visto_bueno);type(date);null"`
	RevisorId         int       `orm:"column(revisor_id)"`
	Observaciones     string    `orm:"column(observaciones);null"`
	Activo            bool      `orm:"column(activo)"`
	FechaCreacion     time.Time `orm:"column(fecha_creacion);type(timestamp without time zone)"`
	FechaModificacion time.Time `orm:"column(fecha_modificacion);type(timestamp without time zone)"`
}

type ActaRecibidoUbicacion struct {
	Id                int
	UbicacionId       *AsignacionEspacioFisicoDependencia
	FechaVistoBueno   time.Time
	RevisorId         int
	Observaciones     string
	Activo            bool
	EstadoActaId      *EstadoActa
	FechaCreacion     time.Time
	FechaModificacion time.Time
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
