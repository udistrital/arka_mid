// Guardar en este archivo los modelos/structs no asociados
// a un crud u otra API en particular

package models

import "time"

type PorDefinir struct{}

type ActaResumen struct {
	Estado            string
	EstadoActaId      *EstadoActa
	FechaCreacion     time.Time
	FechaModificacion time.Time
	FechaVistoBueno   time.Time
	Id                int
	Observaciones     string
	PersonaAsignada   string
	RevisorId         string
	UbicacionId       string
}

type ElementoActaCargado struct {
	Id                 int
	Nombre             string
	Marca              string
	Serie              string
	Cantidad           float64
	UnidadMedida       int
	ValorUnitario      float64
	Subtotal           float64
	Descuento          float64
	PorcentajeIvaId    float64
	ValorIva           float64
	ValorTotal         float64
	SubgrupoCatalogoId *PorDefinir
}

type CargaMasivaElementosActa struct {
	Elementos []*ElementoActaCargado
}

type ElementoAperturaKardex struct {
	CantidadMaxima     int
	CantidadMinima     int
	ElementoCatalogoId *PorDefinir
	FechaCreacion      time.Time
	Id                 int
	MetodoValoracion   int
	MovimientoPadreId  *PorDefinir
	Observaciones      string
}

type ElementoSinAsignar struct {
	Activo             bool
	ElementoActaId     int
	ElementoCatalogoId int
	FechaCreacion      time.Time
	FechaModificacion  time.Time
	Id                 int
	Marca              string
	MovimientoId       *PorDefinir
	Nombre             string
	SaldoCantidad      int
	SaldoValor         float64
	Serie              string
	SubgrupoCatalogoId *PorDefinir
	Unidad             int
	ValorResidual      float64
	ValorTotal         float64
	ValorUnitario      float64
	VidaUtil           float64
}

type ExistenciasKardex struct {
	ElementoCatalogoId *ElementoCatalogo
	Id                 int
	SaldoCantidad      int
	SaldoValor         float64
	SubgrupoCatalogo   *SubgrupoCatalogo
}

type DetalleEntrada struct {
	Contrato   *PorDefinir `json:"contrato"`
	Movimiento *PorDefinir `json:"movimiento"`
	Proveedor  *PorDefinir `json:"proveedor"`
}

type MvtoArkaMasTransaccion struct {
	TrContable *PorDefinir
	Movimiento *PorDefinir `json:"omitempty"`
}
