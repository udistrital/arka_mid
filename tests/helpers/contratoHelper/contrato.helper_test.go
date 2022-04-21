package contratoHelper_test

import (
	"flag"
	"os"
	"testing"

	"github.com/astaxie/beego"

	"github.com/udistrital/arka_mid/helpers/crud/administrativa"
)

var parameters struct {
	ADMINISTRATIVA_JBPM string
}

func TestMain(m *testing.M) {
	parameters.ADMINISTRATIVA_JBPM = os.Getenv("ADMINISTRATIVA_JBPM")
	beego.AppConfig.Set("administrativaJbpmService", os.Getenv("ADMINISTRATIVA_JBPM"))
	flag.Parse()
	os.Exit(m.Run())
}

// GetCatalogoById ...
func TestGetContrato(t *testing.T) {
	valor, err := administrativa.GetContrato(15, "2020")
	if err != nil || valor == nil {
		if err != nil {
			t.Error("No se pudo consultar el contrato", err)
		} else {
			t.Error("No se pudo consultar el contrato", err)
		}
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetCatalogoById Finalizado Correctamente (OK)")
	}
}

func TestEndPointGetAdministrativa_jbpm_crud(t *testing.T) {
	t.Log("Testing EndPoint ADMINISTRATIVA_JBPM")
	t.Log(parameters.ADMINISTRATIVA_JBPM)
}
