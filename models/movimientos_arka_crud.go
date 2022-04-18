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
	VidaUtil          float64
	ValorResidual     float64
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

type NovedadElemento struct {
	Id                   int
	VidaUtil             float64
	ValorLibros          float64
	ValorResidual        float64
	Metadata             string
	MovimientoId         *Movimiento
	ElementoMovimientoId *ElementosMovimiento
	Activo               bool
	FechaCreacion        time.Time
	FechaModificacion    time.Time
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

type FormatoTraslado struct {
	Consecutivo        string
	ConsecutivoId      int
	Ubicacion          int
	FuncionarioOrigen  int
	FuncionarioDestino int
	Elementos          []int
	RazonRechazo       string
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

type DetalleTrasladoLista struct {
	Id                 int
	Consecutivo        string
	FuncionarioOrigen  string
	FuncionarioDestino string
	FechaCreacion      string
	Ubicacion          string
	EstadoMovimientoId int
}

type DetalleElementoPlaca struct {
	Id             int
	ElementoActaId int
	Placa          string
	Nombre         string
	Marca          string
	Serie          string
	Valor          float64
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
	ConsecutivoId  int
	Elementos      []int
	FechaRevisionA string
	FechaRevisionC string
	Funcionario    int
	Revisor        int
	RazonRechazo   string
	Resolucion     string
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
	Id             int
	Soporte        int
	Funcionario    *InfoTercero
	Revisor        *InfoTercero
	TipoBaja       *FormatoTipoMovimiento
	Elementos      []*DetalleElementoBaja
	Observaciones  string
	Consecutivo    string
	RazonRechazo   string
	Resolucion     string
	FechaRevisionC string
}

type DetalleElementoBaja struct {
	Id                 int
	Placa              string
	Nombre             string
	Marca              string
	Serie              string
	SubgrupoCatalogoId *DetalleSubgrupo
	Historial          *Historial
	Ubicacion          *DetalleSedeDependencia
	Funcionario        *InfoTercero
}

type TrRevisionBaja struct {
	Bajas          []int
	Aprobacion     bool
	RazonRechazo   string
	FechaRevisionC string
	Resolucion     string
}

type Historial struct {
	Salida       *Movimiento
	Traslados    []*Movimiento
	Baja         *Movimiento
	Depreciacion *Movimiento
}

type FormatoDepreciacion struct {
	ConsecutivoId int
	FechaCorte    string
	Totales       map[int]float64
	RazonRechazo  string
}

type DetalleCorteDepreciacion struct {
	ValorPresente        float64
	ElementoMovimientoId int
	VidaUtil             float64
	ElementoActaId       int
	ValorResidual        float64
	NovedadElementoId    int
	FechaRef             time.Time
}
type InfoDepreciacion struct {
	Id            int
	RazonRechazo  string
	FechaCorte    time.Time
	Observaciones string
	Tipo          string
}

type DetalleAjusteAutomatico struct {
	Movimiento *Movimiento
	TrContable []*DetalleMovimientoContable
	Elementos  []*DetalleElemento__
}

type ElementosPorActualizarSalida struct {
	Salida    *Movimiento
	UpdateSg  []*DetalleElemento_
	UpdateVls []*DetalleElemento_
}

type FormatoAjusteAutomatico struct {
	Consecutivo string
	Elementos   []int
	TrContable  int
}
