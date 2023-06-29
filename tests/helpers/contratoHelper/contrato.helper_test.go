package contratoHelper_test

import (
	"flag"
	"os"
	"testing"

	beego "github.com/beego/beego/v2/server/web"

	administrativa_ "github.com/udistrital/administrativa_mid_api/models"
	"github.com/udistrital/arka_mid/helpers/crud/administrativa"
)

var parameters struct {
	ADMINISTRATIVA_JBPM string
}

func TestMain(m *testing.M) {
	parameters.ADMINISTRATIVA_JBPM = os.Getenv("ADMINISTRATIVA_JBPM")
	if err := beego.AppConfig.Set("administrativaJbpmService", os.Getenv("ADMINISTRATIVA_JBPM")); err != nil {
		panic(err)
	}
	flag.Parse()
	os.Exit(m.Run())
}

// GetCatalogoById ...
func TestGetContrato(t *testing.T) {
	var contrato administrativa_.InformacionContrato
	err := administrativa.GetContrato(15, "2020", &contrato)
	if err != nil {
		if err != nil {
			t.Error("No se pudo consultar el contrato", err)
		} else {
			t.Error("No se pudo consultar el contrato", err)
		}
		t.Fail()
	} else {
		t.Log(contrato)
		t.Log("TestGetCatalogoById Finalizado Correctamente (OK)")
	}
}

func TestEndPointGetAdministrativa_jbpm_crud(t *testing.T) {
	t.Log("Testing EndPoint ADMINISTRATIVA_JBPM")
	t.Log(parameters.ADMINISTRATIVA_JBPM)
}
