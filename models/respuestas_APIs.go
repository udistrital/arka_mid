package models

// RespuestaAPI1 es típica de una API que encapsula
// su respuesta en una propiedad "Data".
// NO USAR DIRECTAMENTE: La propiedad Data NO hace parte de
// este struct, usar RespuestaAPI1Obj o RespuestaAPI1Arr
type RespuestaAPI1 struct {
	Message string
	Status  string
	Success bool
}

// RespuestaAPI1Obj es un RespuestaAPI1 donde
// los datos son un único objeto
type RespuestaAPI1Obj struct {
	RespuestaAPI1
	Data map[string]interface{}
}

// RespuestaAPI1Arr es un RespuestaAPI1 donde
// los datos son un arreglo de objetos
type RespuestaAPI1Arr struct {
	RespuestaAPI1
	Data []map[string]interface{}
}
