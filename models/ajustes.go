package models

type PreMovAjuste struct {
	Cuenta      string
	Debito      float64
	Credito     float64
	Descripcion string
	TerceroId   int
}

type PreTrAjuste struct {
	Descripcion string
	Movimientos []*PreMovAjuste
}

type FormatoAjuste struct {
	PreTrAjuste  *PreTrAjuste
	RazonRechazo string
}

type DetalleAjuste struct {
	Movimiento *Movimiento
	TrContable []*DetalleMovimientoContable
}
