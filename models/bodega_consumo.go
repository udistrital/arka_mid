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

type FormatoSolicitudBodega struct {
	Funcionario int
}

type DetalleSolicitudBodega struct {
	Movimiento
	Solicitante IdentificacionTercero
}
