package models


type Impuesto struct {
	Id						int
    Nombre					string
    Descripcion				string
    CodigoAbreviacion		string
    Activo					bool
}

type VigenciaImpuesto struct {
	Id						int
    Activo					bool
    Tarifa					int
    PorcentajeAplicacion	int
    ImpuestoId				Impuesto
}
