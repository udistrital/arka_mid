package models

import "time"

type NivelClasificacion struct {
	Id                int
	Nombre            string
	Longitud          int
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
	NumeroOrden       int
}

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

type TipoComprobante struct {
	Id            int
	Activo        bool
	Codigo        string
	TipoDocumento string
	Descripcion   string
}

type Comprobante struct {
	Id              string `json:"_id" bson:"_id,omitempty"`
	Codigo          int
	Descripcion     string
	Comprobante     string
	Numero          int
	TipoComprobante *TipoComprobante
}

type Etiquetas struct {
	ComprobanteId string
}

type DetalleCuenta struct {
	Id              string
	Codigo          string
	Nombre          string
	RequiereTercero bool
}

type DetalleMovimientoContable struct {
	Cuenta      *DetalleCuenta
	Debito      float64
	Credito     float64
	Descripcion string
	TerceroId   *IdentificacionTercero
}

type InfoTransaccionContable struct {
	Movimientos []*DetalleMovimientoContable `json:"movimientos"`
	Concepto    string
	Fecha       time.Time
}
