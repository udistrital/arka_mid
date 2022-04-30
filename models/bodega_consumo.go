package models

type BodegaConsumoSolicitud struct {
	Elementos *ElementoSolicitud
	Solicitud *Movimiento
}

type ElementoSolicitud struct {
	Cantidad           uint
	CantidadAprobada   uint
	ElementoCatalogoId *ElementoCatalogo
	Id                 uint
	Nombre             string
	SaldoCantidad      uint
	SaldoValor         uint
	Sede               *EspacioFisico
	Dependencia        *Dependencia
	Ubicacion          *AsignacionEspacioFisicoDependencia
}

type ElementoSolicitud_ struct {
	Cantidad           int
	Ubicacion          int
	ElementoCatalogoId int
	CantidadAprobada   int
}
type FormatoSolicitudBodega struct {
	Funcionario   int
	ConsecutivoId int
	Consecutivo   string
	Elementos     []ElementoSolicitud_
}

type DetalleSolicitudBodega struct {
	Movimiento
	Solicitante IdentificacionTercero
}
