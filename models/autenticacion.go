package models

// UsuarioAutenticacion Modelo retornado por
// la MID API de Autenticacion
//
// Endpoint: /token/userRol
type UsuarioAutenticacion struct {
	Codigo             string
	Estado             string
	FamilyName         string
	Documento          string   `json:"documento"`
	DocumentoCompuesto string   `json:"documento_compuesto"`
	Email              string   `json:"email"`
	Role               []string `json:"role"`
}

// UsuarioDataRequest Modelo requerido en el cuerpo
// de la solicitud GET al endpoint /token/userRol
// de la MID API de Autenticación
type UsuarioDataRequest struct {
	User string `json:"user"`
}

// RolesArka subconjunto de roles que retorna WSO2
// en uso específico de ARKA
var RolesArka = map[string]string{
	"Admin":       "ADMIN_ARKA",
	"Revisor":     "REVISOR",
	"Contratista": "CONTRATISTA",
	"Proveedor":   "PROVEEDOR",
	"Secretaria":  "SECRETARIA_ALMACEN",
}
