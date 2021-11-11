package models

type Catalogo struct {
	Id          int
	Nombre      string
	Descripcion string
	FechaInicio string
	FechaFin    string
	Activo      bool
}

type CuentasSubgrupo struct {
	Id                  int
	CuentaCreditoId     string
	CuentaDebitoId      string
	SubtipoMovimientoId int
	SubgrupoId          *Subgrupo
	Activo              bool
}

type DetalleSubgrupo struct {
	Id           int
	Depreciacion bool
	Valorizacion bool
	Deterioro    bool
	Activo       bool
	SubgrupoId   *Subgrupo
	TipoBienId   *TipoBien
}

type ElementoCatalogo struct {
	Id          uint
	Nombre      string
	Descripcion string
	Codigo      string
	Activo      bool
	SubgrupoId  *Subgrupo
}

type RelacionNivel struct {
	Id           int
	Activo       bool
	NivelPadreId *RelacionNivel
	NivelHijoId  *RelacionNivel
}

type SubgrupoCatalogo struct {
	Id         int
	Activo     bool
	CatalogoId *Catalogo
	SubgrupoId *Subgrupo
}

type SubgrupoSubgrupo struct {
	Id              int
	Activo          bool
	SubgrupoPadreId *Subgrupo
	SubgrupoHijoId  *Subgrupo
}

type Subgrupo struct {
	Id          int
	Nombre      string
	Descripcion string
	Activo      bool
	Codigo      string
	TipoNivelId *TipoNivel
}

type TipoBien struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Orden             float64
	NecesitaPlaca     bool
	Reglas            string
	Activo            bool
	Tipo_bien_padre   *TipoBien
}

type TipoNivel struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Orden             float64
	Activo            bool
}
