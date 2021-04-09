package models

type CuentaContable struct {
	Activo             bool
	Id                 int
	Ajustable          bool
	Saldo              float64
	Nombre             string
	Naturaleza         string
	Descripcion        string
	Codigo             string
	NivelClasificacion *NivelClasificacion
}
