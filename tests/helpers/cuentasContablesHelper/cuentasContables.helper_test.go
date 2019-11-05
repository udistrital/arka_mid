package cuentasContablesHelper_test

import (
	"flag"
	"os"
	"testing"

	"github.com/udistrital/arka_mid/helpers/cuentasContablesHelper"
)

var parameters struct {
	GetCuentaContable string
}

func TestMain(m *testing.M) {
	parameters.GetCuentaContable = os.Getenv("GetCuentaContable")
	flag.Parse()
	os.Exit(m.Run())
}

// GetCuentaContable ...
func TestGetCuentaContable(t *testing.T) {
	valor, err := cuentasContablesHelper.GetCuentaContable(1)
	if err != nil {
		t.Error("No se pudo consultar las cuentas contables", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetCuentaContable Finalizado Correctamente (OK)")
	}
}

func TestEndPointGetCuentaContableo(t *testing.T) {
	t.Log("Testing EndPoint GetCuentaContable")
	t.Log(parameters.GetCuentaContable)
}
