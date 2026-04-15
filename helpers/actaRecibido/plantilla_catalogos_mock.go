package actaRecibido

// plantillaCatalogosMock concentra los datos de apoyo usados por la plantilla
// de cargue masivo. Están definidos localmente para que la generación del Excel
// no dependa de consultas externas y sea fácil reemplazarlos después por una
// fuente real.
type plantillaCatalogosMock struct {
	Clases   []plantillaClaseMock
	Tipos    []plantillaTipoBienMock
	Ivas     []plantillaIVAMock
	Unidades []plantillaUnidadMock
}

type plantillaClaseMock struct {
	Codigo string
	Nombre string
}

type plantillaTipoBienMock struct {
	Id     int
	Nombre string
}

type plantillaIVAMock struct {
	Tarifa int
}

type plantillaUnidadMock struct {
	Nombre string
}

// getPlantillaCatalogosMock retorna el catálogo mockeado usado únicamente por
// la generación de la plantilla Excel. El punto natural de reemplazo futuro es
// esta función, sustituyendo estas estructuras por un origen real.
func getPlantillaCatalogosMock() plantillaCatalogosMock {
	return plantillaCatalogosMock{
		Clases: []plantillaClaseMock{
			{Codigo: "010001", Nombre: "COMPUTO - (DEV)"},
			{Codigo: "010002", Nombre: "COMUNICACIONES - (DEV)"},
			{Codigo: "010003", Nombre: "MUEBLES Y ENSERES - (DEV)"},
			{Codigo: "010004", Nombre: "EQUIPO DE LABORATORIO - (DEV)"},
			{Codigo: "010005", Nombre: "EQUIPO AUDIOVISUAL - (DEV)"},
			{Codigo: "010006", Nombre: "HERRAMIENTA MENOR - (DEV)"},
		},
		Tipos: []plantillaTipoBienMock{
			{Id: 1, Nombre: "CONSUMO"},
			{Id: 2, Nombre: "DEVOLUTIVO"},
			{Id: 3, Nombre: "CONTROL ADMINISTRATIVO"},
		},
		Ivas: []plantillaIVAMock{
			{Tarifa: 0},
			{Tarifa: 5},
			{Tarifa: 19},
		},
		Unidades: []plantillaUnidadMock{
			{Nombre: "NO APLICA"},
			{Nombre: "UNIDAD"},
		},
	}
}
