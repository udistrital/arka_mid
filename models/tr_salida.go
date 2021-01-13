package models

type TrSalida struct {
	Salida    *Movimiento
	Elementos []*ElementosMovimiento
}
type SalidaGeneral struct {
	Salidas []*TrSalida
}

type TrSalida2 struct {
	TrSalida
	MovimientosKronos *MovimientoProcesoExterno
}
