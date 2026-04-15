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
			{Codigo: "010101", Nombre: "COMPUTO - (DEV)"},
			{Codigo: "020202", Nombre: "COMUNICACIÓN - (DEV)"},
			{Codigo: "030303", Nombre: "EDIFICACIONES - (DEV)"},
			{Codigo: "040404", Nombre: "ELEMENTOS DE ARTE DEVOLUTIVO - (DEV)"},
			{Codigo: "050505", Nombre: "EQU Y MAQ PARA DEPORTES - (DEV)"},
			{Codigo: "060606", Nombre: "EQU Y MAQ PARA TRANSPORTES - (DEV)"},
			{Codigo: "070707", Nombre: "EQUI Y MAQ PARA CONSTRUCCION - (DEV)"},
			{Codigo: "080808", Nombre: "EQUI Y MAQUI COMEDOR Y COCINA - (DEV)"},
			{Codigo: "090909", Nombre: "EQUIPOS PARA MEDICINA Y ODONTOL - (DEV)"},
			{Codigo: "101010", Nombre: "INSTRUMENTOS MUSICALES - (DEV)"},
			{Codigo: "111111", Nombre: "LABORATORIO (DEV)"},
			{Codigo: "121212", Nombre: "MUEBLES Y ENSERES - (DEV)"},
			{Codigo: "131313", Nombre: "SOFTWARE - (DEV)"},
			{Codigo: "141414", Nombre: "LIBROS Y BIBLIOTECAS - (DEV)"},
			{Codigo: "212121", Nombre: "ELEMENTOS DE ARTE - (CTRL)"},
			{Codigo: "222222", Nombre: "EQUIPO DE RECREACION Y DEPORTE - (CTRL)"},
			{Codigo: "232323", Nombre: "EQUIPO Y MAQUINARIA PARA COMEDOR Y COCINA - (CTRL)"},
			{Codigo: "242424", Nombre: "EQUIPO Y MAQUINARIA PARA COMPUTACION - (CTRL)"},
			{Codigo: "262626", Nombre: "EQUIPO Y MAQUINARIA PARA CONSTRUCCION - (CTRL)"},
			{Codigo: "272727", Nombre: "EQUIPO Y MAQUINARIA PARA LABORATORIO - (CTRL)"},
			{Codigo: "282828", Nombre: "EQUIPO Y MAQUINARIA PARA OFICINA - (CTRL)"},
			{Codigo: "292929", Nombre: "EQUIPO Y MAQUINARIA PARA TRANSPORTE - (CTRL)"},
			{Codigo: "303030", Nombre: "EQUIPOS PARA MEDICINA Y ODONTOLOGIA - (CTRL)"},
			{Codigo: "313131", Nombre: "INSTRUMENTOS MUSICALES - (CTRL)"},
			{Codigo: "323232", Nombre: "LIBROS - (CTRL)"},
			{Codigo: "333333", Nombre: "MUEBLES Y ENSERES - (CTRL)"},
			{Codigo: "343434", Nombre: "SOFTWARE - (CTRL)"},
			{Codigo: "353535", Nombre: "EQUIPO Y MAQUINARIA PARA COMUNICACIÓN - (CTRL)"},
			{Codigo: "363636", Nombre: "DOTACION A TRABAJADORES - (CTRL)"},
			{Codigo: "373737", Nombre: "ELEMENTOS DE CONSTRUCCIÓN - (CTRL)"},
			{Codigo: "404040", Nombre: "TERRENOS (DEV)"},
			{Codigo: "505050", Nombre: "ELEMENTO CONSUMO ALMACEN - (CONS)"},
			{Codigo: "515151", Nombre: "LUBRICANTES Y COMBUSTIBLES (CONS CONT)"},
			{Codigo: "616161", Nombre: "MANTENIMIENTO"},
			{Codigo: "626262", Nombre: "INSTALACIONES"},
			{Codigo: "636363", Nombre: "REPARACIONES"},
			{Codigo: "646464", Nombre: "TRANSPORTE"},
		},
		Tipos: []plantillaTipoBienMock{
			{Id: 12, Nombre: "CONSUMO"},
			{Id: 10, Nombre: "DEVOLUTIVO"},
			{Id: 9, Nombre: "CONSUMO CONTROLADO"},
		},
		Ivas: []plantillaIVAMock{
			{Tarifa: 0},
			{Tarifa: 5},
			{Tarifa: 10},
			{Tarifa: 16},
			{Tarifa: 19},
		},
		Unidades: []plantillaUnidadMock{
			{Nombre: "NO APLICA"},
			{Nombre: "UNIDAD"},
		},
	}
}
