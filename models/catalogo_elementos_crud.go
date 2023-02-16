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
	TipoMovimientoId    int
	SubtipoMovimientoId int
	SubgrupoId          *Subgrupo
	TipoBienId          *TipoBien
	Activo              bool
}

type DetalleSubgrupo struct {
	Id            int
	Depreciacion  bool
	Valorizacion  bool
	Amortizacion  bool
	VidaUtil      float64
	ValorResidual float64
	Activo        bool
	SubgrupoId    *Subgrupo
	TipoBienId    *TipoBien
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
	Id              int
	Nombre          string
	Descripcion     string
	Activo          bool
	Reglas          string
	LimiteInferior  float64
	LimiteSuperior  float64
	NecesitaPlaca   bool
	NecesitaPoliza  bool
	BodegaConsumo   bool
	TipoBienPadreId *TipoBien
}

type TipoNivel struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Orden             float64
	Activo            bool
}

type DetalleCuentasSubgrupo struct {
	Id                  int
	CuentaCreditoId     *DetalleCuenta
	CuentaDebitoId      *DetalleCuenta
	TipoMovimientoId    *FormatoTipoMovimiento
	SubtipoMovimientoId *FormatoTipoMovimiento
	TipoBienId          TipoBien
	SubgrupoId          int
}
