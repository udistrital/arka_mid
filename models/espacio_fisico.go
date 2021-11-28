package models

type EspacioFisico struct {
	Id                  int
	Nombre              string
	CodigoAbreviacion   string
	Activo              bool
	TipoEspacio         *TipoEspacio
	Descripcion         string
	TipoTerrenoId       int
	TipoEdificacionId   int
	TipoEspacioFisicoId *TipoEspacioFisicoV2
}

type TipoEspacioFisicoV2 struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
}

type DetalleSedeDependencia struct {
	Sede        *EspacioFisico
	Dependencia *Dependencia
	Ubicacion   int
}
