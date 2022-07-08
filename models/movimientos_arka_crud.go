package models

import (
	"time"
)

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

type TransaccionEntrada struct {
	Id                      int
	Observacion             string
	Detalle                 FormatoBaseEntrada
	FormatoTipoMovimientoId string
	SoporteMovimientoId     int
}

type FormatoBaseEntrada struct {
	ActaRecibidoId      int    `json:"acta_recibido_id"`
	Consecutivo         string `json:"consecutivo"`
	ConsecutivoId       int
	ContratoId          int     `json:"contrato_id"`
	Divisa              string  `json:"divisa"`
	EncargadoId         int     `json:"encargado_id"`
	Factura             int     `json:"factura"`
	OrdenadorGastoId    int     `json:"ordenador_gasto_id"`
	Placa               string  `json:"placa_id"`
	RegistroImportacion string  `json:"num_reg_importacion"`
	SupervisorId        int     `json:"supervisor"`
	TRM                 float64 `json:"TRM"`
	Vigencia            string  `json:"vigencia"`
	VigenciaContrato    string  `json:"vigencia_contrato"`
	VigenciaOrdenador   string  `json:"vigencia_ordenador"`
	VigenciaSolicitante string  `json:"vigencia_solicitante"`
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

type ResultadoMovimiento struct {
	Error               string
	Movimiento          Movimiento
	TransaccionContable InfoTransaccionContable
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
	Consecutivo    string
	ConsecutivoId  int
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
	Salida       *Movimiento
	Traslados    []*Movimiento
	Baja         *Movimiento
	Depreciacion *Movimiento
}

type FormatoDepreciacion struct {
	ConsecutivoId int
	Consecutivo   string
	FechaCorte    string
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

type DepreciacionElemento struct {
	DeltaValor           float64
	ElementoMovimientoId int
	ElementoActaId       int
}

type TransaccionCierre struct {
	MovimientoId         int
	ElementoMovimientoId []int
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
	Consecutivo   string
	ConsecutivoId int
	Elementos     []int
}

type ConsecutivoMovimiento struct {
	Consecutivo   string
	ConsecutivoId int
}

type FormatoSalida struct {
	Consecutivo   string `json:"consecutivo"`
	ConsecutivoId int
	Funcionario   int `json:"funcionario"`
	Ubicacion     int `json:"ubicacion"`
}
