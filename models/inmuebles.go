package models

type Inmueble struct {
	Elemento           Elemento
	ElementoMovimiento ElementosMovimiento
	EspacioFisico      EspacioFisico
	Sede               EspacioFisico
	SubgrupoId         Subgrupo
	Cuentas            ParametrizacionContable_
	CuentasMediciones  ParametrizacionContable_
	Otros              []EspacioFisicoCampo
}

type ParametrizacionContable struct {
	CuentaCreditoId string
	CuentaDebitoId  string
}

type ParametrizacionContable_ struct {
	CuentaCreditoId CuentaContable
	CuentaDebitoId  CuentaContable
}
