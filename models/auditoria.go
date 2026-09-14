package models

// EventoAuditoria representa una operación funcional que se desea
// registrar explícitamente en los logs de arka_mid.
//
// El usuario NO se recibe en este objeto.
// utils_oas/v2/auditoria obtiene el usuario desde el token Authorization.
type EventoAuditoria struct {
	Evento            string           `json:"evento"`
	RequestOriginal   RequestAuditoria `json:"request_original"`
	ResultadoOriginal ResultAuditoria  `json:"resultado_original"`
}

// RequestAuditoria contiene la información de la petición funcional
// que originó el evento de auditoría.
type RequestAuditoria struct {
	Method   string      `json:"method"`
	Endpoint string      `json:"endpoint"`
	Body     interface{} `json:"body"`
}

// ResultAuditoria contiene la respuesta REAL obtenida por el cliente
// al ejecutar la petición funcional.
//
// HttpStatus y Body NO son valores quemados.
// Los envía el cliente con base en la respuesta real de la operación.
type ResultAuditoria struct {
	HttpStatus int         `json:"http_status"`
	Body       interface{} `json:"body"`
}
