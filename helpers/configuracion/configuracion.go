package configuracion

const (
	NombreParametroRoles              = "RolesRegistrados"
	NombreParametroTiposDeComprobante = "TiposDeComprobante"
)

var (
	roles        RolesArka
	comprobantes TiposContablesArka
)

func init() {
	ActualizaRolesArka()
}
