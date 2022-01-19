package models

type CuentaContable struct {
	Activo             bool
	Id                 string
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

type DetalleCuenta struct {
	Codigo          string
	Nombre          string
	RequiereTercero bool
}

type DetalleMovimientoContable struct {
	Cuenta      *DetalleCuenta
	Debito      float64
	Credito     float64
	Descripcion string
	TerceroId   *Tercero
}
