package ubicacionHelper_test

import (
	"flag"
	"os"
	"testing"

	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers/unidadHelper"
)

var parameters struct {
	ADMINISTRATIVA_SERVICE string
}

func TestMain(m *testing.M) {
	parameters.ADMINISTRATIVA_SERVICE = os.Getenv("ADMINISTRATIVA_SERVICE")
	beego.AppConfig.Set("administrativaService", parameters.ADMINISTRATIVA_SERVICE)
	flag.Parse()
	os.Exit(m.Run())
}

// GetUnidad ...
func TestGetUnidad(t *testing.T) {

	if valor, err := unidadHelper.GetUnidad(1); err != nil {
		t.Error("No se pudo consultar la ubicacion", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetUnidad Finalizado Correctamente")
	}
}

func TestEndPointAdministrativaService(t *testing.T) {
	t.Log("Testing EndPoint administrativaService")
	t.Log(parameters.ADMINISTRATIVA_SERVICE)
}
