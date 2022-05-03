package configuracion

const (
	NombreParametroRoles              = "RolesRegistrados"
	NombreParametroTiposDeComprobante = "ComprobantesKronos"
)

var (
	roles        RolesArka
	comprobantes TiposContablesArka
)

func init() {
	ActualizaRolesArka()
	ActualizaTiposDeComprobante()
}
