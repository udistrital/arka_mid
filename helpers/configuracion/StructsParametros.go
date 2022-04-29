// Documentar en este archivo las estructuras de datos
// a almacenar en registros de configuracion_crud.parametro

package configuracion

// RolesArka es la estructura de los datos a almacenar en el
// campo Valor de un regitro de configuracion_api.parametro,
// cuyo nombre es el de la constante NombreParametroRoles.
//
// El objetivo es guardar los id de registros de
// configuracion_api.perfil asociados a cada rol necesario
// desde el modelo de negocios de Arka II
//
// Adicionalmente, esto es para permitir que los roles se puedan
// renombrar directamente desde el cliente de configuración
// sin afectar el funcionamiento del sistema
type RolesArka struct {
	Administrador uint
	Secretaria    uint
	AuxiliarUno   uint
	AuxiliarDos   uint
	Proveedor     uint
	Contabilidad  uint
}

// TiposContablesArka es la estructura a guardar en el campo
// Valor de un registro en configuracion_api.parametro cuyo
// Nombre es el de la constante NombreParametroTiposDeComprobante.
//
// Su objetivo es guardar los id de los comprobantes contables
// de cuentas_contables_crud.comprobante a ser especificados
// en el campo de etiquetas de cada transaccion contable
// enviada a movimientos_contables_mid
type TiposContablesArka struct {
	Entrada      string // P8
	Salida       string // H21
	Traslado     string // N39 (ajuste)
	Baja         string // H23
	Amortizacion string // H22 (cierre)
	Depreciacion string // H22 (cierre)
	Ajuste       string // N39
	Avaluo       string // N40
}
