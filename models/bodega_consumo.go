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

type ElementoCatalogo struct {
	Id          uint
	Nombre      string
	Descripcion string
	Codigo      string
	Activo      bool
}
