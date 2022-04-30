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

// RespuestaAPI1Arr es un RespuestaAPI1 donde
// los datos es un string
type RespuestaAPI1Str struct {
	RespuestaAPI1
	Data string
}

// RespuestaAPI1Arr es un RespuestaAPI1 donde
// los datos varían según la petición
type RespuestaAPI1Interface struct {
	RespuestaAPI1
	Data interface{}
}

// RespuestaAPI2 es similar a RespuestaAPI1
// Agrupa las estructuras comunes y por lo mismo
// NO se debe usar directamente
type RespuestaAPI2 struct {
	Code    int
	Message string
}

// RespuestaAPI2obj es un RespuestaAPI2 donde
// el Body es un objeto
type RespuestaAPI2obj struct {
	RespuestaAPI2
	Body map[string]interface{}
}

// RespuestaAPI2arr es un RespuestaAPI2 donde
// el Body es un arreglo
type RespuestaAPI2arr struct {
	RespuestaAPI2
	Body []map[string]interface{}
}
