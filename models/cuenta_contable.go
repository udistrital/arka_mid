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
	Padre              string
	Hijos              []string
	RequiereTercero    bool
	NivelClasificacion *NivelClasificacion
}

type TipoComprobanteContable struct {
	Id            int
	Activo        bool
	Codigo        string
	TipoDocumento string
}
