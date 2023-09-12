package models

import (
	"time"
)

type Movimiento struct {
	Id                      int
	Observacion             string
	ConsecutivoId           *int
	Consecutivo             *string
	FechaCorte              *time.Time
	Detalle                 string
	FechaCreacion           time.Time
	FechaModificacion       time.Time
	Activo                  bool
	MovimientoPadreId       *Movimiento
	FormatoTipoMovimientoId *FormatoTipoMovimiento
	EstadoMovimientoId      *EstadoMovimiento
}

type ElementosMovimiento struct {
	Id                 int
	ElementoActaId     *int
	ElementoCatalogoId int
	Unidad             float64
	ValorUnitario      float64
	ValorTotal         float64
	SaldoCantidad      float64
	SaldoValor         float64
	VidaUtil           float64
	ValorResidual      float64
	Activo             bool
	FechaCreacion      time.Time
	FechaModificacion  time.Time
	MovimientoId       *Movimiento
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

type CentroCostos struct {
	Id          int
	Dependencia string
	Sede        string
	Codigo      string
	Nombre      string
}

type TransaccionEntrada struct {
	Id                      int
	Observacion             string
	Detalle                 FormatoBaseEntrada
	FormatoTipoMovimientoId string
	SoporteMovimientoId     int
}

type FormatoBaseEntrada struct {
	ActaRecibidoId      int                `json:"acta_recibido_id"`
	ContratoId          int                `json:"contrato_id"`
	Divisa              string             `json:"divisa"`
	Factura             int                `json:"factura"`
	OrdenadorGastoId    int                `json:"ordenador_gasto_id"`
	Elementos           []ElementoMejorado `json:"elementos"`
	RegistroImportacion string             `json:"num_reg_importacion"`
	SupervisorId        int                `json:"supervisor"`
	TRM                 float64            `json:"TRM"`
	VigenciaContrato    string             `json:"vigencia_contrato"`
}

type ElementoMejorado struct {
	Id            int
	AprovechadoId *int
	ValorLibros   *float64
	VidaUtil      *float64
	ValorResidual *float64
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
	Salidas []TrSalida
}

type ResultadoMovimiento struct {
	Error               string
	Movimiento          Movimiento
	TransaccionContable InfoTransaccionContable
}

type FormatoTraslado struct {
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
	FuncionarioOrigen  Tercero
	FuncionarioDestino Tercero
	FechaCreacion      string
	Ubicacion          string
	EstadoMovimientoId int
}

type InventarioTercero struct {
	Elementos []DetalleElementoPlaca
	Tercero   DetalleFuncionario
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
	Movimiento         *Movimiento
	Observaciones      string
	Elementos          []*DetalleElementoPlaca
	FuncionarioOrigen  *DetalleFuncionario
	FuncionarioDestino *DetalleFuncionario
	Ubicacion          *DetalleSedeDependencia
	TrContable         *InfoTransaccionContable
}

type FormatoBaja struct {
	Elementos      []int
	FechaRevisionA string
	FechaRevisionC string
	Funcionario    int
	Revisor        int
	RazonRechazo   string
	Resolucion     string
	DependenciaId  int
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
	Elementos      []*DetalleElementoBaja
	FechaRevisionC string
	Funcionario    *InfoTercero
	Id             int
	Movimiento     *Movimiento
	Observaciones  string
	RazonRechazo   string
	Resolucion     string
	Revisor        *InfoTercero
	Soporte        int
	TipoBaja       *FormatoTipoMovimiento
	TrContable     *InfoTransaccionContable
	DependenciaId  string
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
	DependenciaId  int
}

type Historial struct {
	Elemento  *ElementosMovimiento
	Entradas  []*Movimiento
	Salida    *Movimiento
	Traslados []*Movimiento
	Novedades []NovedadElemento
	Baja      *Movimiento
}

type FormatoDepreciacion struct {
	RazonRechazo string
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

type DepreciacionElemento struct {
	DeltaValor           float64
	ElementoMovimientoId int
	ElementoActaId       int
}

type InfoDepreciacion struct {
	Id            int
	RazonRechazo  string
	FechaCorte    time.Time
	Observaciones string
	Rechazar      bool
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
	Elementos []int
}

type FormatoSalida struct {
	Funcionario int `json:"funcionario"`
	Ubicacion   int `json:"ubicacion"`
}

type FormatoSalidaCostos struct {
	FormatoSalida
	CentroCostos string `json:"centro_costos"`
}
