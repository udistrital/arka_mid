package cuentasContablesHelper_test

import (
	"flag"
	"os"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
)

var parameters struct {
	CUENTAS_CONTABLES_SERVICE string
}

func TestMain(m *testing.M) {
	parameters.CUENTAS_CONTABLES_SERVICE = os.Getenv("CUENTAS_CONTABLES_SERVICE")
	if err := beego.AppConfig.Set("cuentasContablesService", os.Getenv("CUENTAS_CONTABLES_SERVICE")); err != nil {
		panic(err)
	}
	flag.Parse()
	os.Exit(m.Run())
}

// GetCuentaContable ...
func TestGetCuentaContable(t *testing.T) {
	valor, err := cuentasContables.GetCuentaContable("4")
	if err != nil {
		t.Error("No se pudo consultar las cuentas contables", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetCuentaContable Finalizado Correctamente (OK)")
	}
}

func TestEndPointGetCuentasContablesService(t *testing.T) {
	t.Log("Testing EndPoint cuentasContablesService")
	t.Log(parameters.CUENTAS_CONTABLES_SERVICE)
}
