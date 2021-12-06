package models

import "time"

type Movimiento struct {
	Id                      int
	Observacion             string
	Detalle                 string
	FechaCreacion           time.Time
	FechaModificacion       time.Time
	Activo                  bool
	MovimientoPadreId       *Movimiento
	FormatoTipoMovimientoId *FormatoTipoMovimiento
	EstadoMovimientoId      *EstadoMovimiento
	SoporteMovimientoId     int
}

type ElementosMovimiento struct {
	Id                int
	ElementoActaId    int
	Unidad            float64
	ValorUnitario     float64
	ValorTotal        float64
	SaldoCantidad     float64
	SaldoValor        float64
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
	MovimientoId      *Movimiento
}

type FormatoTipoMovimiento struct {
	Id                int
	Nombre            string
	Formato           string
	Descripcion       string
	CodigoAbreviacion string
	NumeroOrden       float64
	FechaCreacion     time.Time
	FechaModificacion time.Time
	Activo            bool
}

type EstadoMovimiento struct {
	Id                int
	Nombre            string
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
	Descripcion       string
}

type SoporteMovimiento struct {
	Id                int
	DocumentoId       int
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
	MovimientoId      *Movimiento
}

type TrSoporteMovimiento struct {
	Movimiento *Movimiento
	Soporte    *SoporteMovimiento
}

type TrSalida struct {
	Salida    *Movimiento
	Elementos []*ElementosMovimiento
}
type SalidaGeneral struct {
	Salidas []*TrSalida
}

type DetalleTraslado struct {
	Id                 int
	FuncionarioOrigen  int
	FuncionarioDestino int
	Elementos          []int
	Ubicacion          int
	Observaciones      string
	MovimientoId       int
	Consecutivo        string
}

type DetalleElementoPlaca struct {
	Id             int
	ElementoActaId int
	Placa          string
	Nombre         string
	Marca          string
}

type TrTraslado struct {
	Detalle            string
	Observaciones      string
	Elementos          []*DetalleElementoPlaca
	FuncionarioOrigen  *DetalleFuncionario
	FuncionarioDestino *DetalleFuncionario
	Ubicacion          *DetalleSedeDependencia
}

type FormatoBaja struct {
	Consecutivo    string
	Elementos      []int
	FechaRevisionA string
	FechaRevisionC string
	Funcionario    int
	Revisor        int
}

type DetalleBaja struct {
	Id                 int
	Consecutivo        string
	FechaCreacion      string
	FechaRevisionA     string
	FechaRevisionC     string
	Funcionario        string
	Revisor            string
	TipoBaja           string
	EstadoMovimientoId int
}

type TrBaja struct {
	Id            int
	Soporte       int
	Funcionario   *InfoTercero
	Revisor       *InfoTercero
	TipoBaja      *FormatoTipoMovimiento
	Elementos     []*DetalleElementoBaja
	Observaciones string
	Consecutivo   string
}

type DetalleElementoBaja struct {
	Id                 int
	Placa              string
	Nombre             string
	Marca              string
	Serie              string
	SubgrupoCatalogoId *DetalleSubgrupo
	Salida             *Movimiento
	Ubicacion          *DetalleSedeDependencia
	Funcionario        *InfoTercero
}

type TrRevisionBaja struct {
	Bajas         []int
	Aprobacion    bool
	Observaciones string
}
