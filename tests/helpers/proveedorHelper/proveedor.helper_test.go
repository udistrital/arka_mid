package proveedorHelper_test

import (
	"flag"
	"os"
	"testing"

	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers/proveedorHelper"
)

var parameters struct {
	ADMINISTRATIVA_SERVICE string
}

func TestMain(m *testing.M) {
	parameters.ADMINISTRATIVA_SERVICE = os.Getenv("ADMINISTRATIVA_SERVICE")
	beego.AppConfig.Set("administrativaService", os.Getenv("ADMINISTRATIVA_SERVICE"))
	flag.Parse()
	os.Exit(m.Run())
}

// GetProveedorById ...
func TestGetProveedorById(t *testing.T) {
	valor, err := proveedorHelper.GetProveedorById(1)
	if err != nil {
		t.Error("No se pudo consultar el proveedor", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetProveedorById Finalizado Correctamente (OK)")
	}
}

func TestEndPointGetAdministrativaService(t *testing.T) {
	t.Log("Testing EndPoint ADMINISTRATIVA_SERVICE")
	t.Log(parameters.ADMINISTRATIVA_SERVICE)
}
