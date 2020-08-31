package models

type CuentaContable struct {
	Id                 int
	Saldo              float64
	Nombre             string
	Naturaleza         string
	Descripcion        string
	Codigo             string
	NivelClasificacion *NivelClasificacion
}
