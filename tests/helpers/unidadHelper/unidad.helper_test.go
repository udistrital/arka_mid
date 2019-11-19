package ubicacionHelper_test

import (
	"flag"
	"os"
	"testing"

	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers/ubicacionHelper"
)

var parameters struct {
	OIKOS_CRUD string
}

func TestMain(m *testing.M) {
	parameters.OIKOS_CRUD = os.Getenv("OIKOS_CRUD")
	beego.AppConfig.Set("oikosService", os.Getenv("OIKOS_CRUD"))
	flag.Parse()
	os.Exit(m.Run())
}

// GetUnidad ...
func TestGetUnidad(t *testing.T) {
	valor, err := ubicacionHelper.GetUbicacion(1)
	if err != nil {
		t.Error("No se pudo consultar la ubicacion", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetUnidad Finalizado Correctamente (OK)")
	}
}

func TestEndPointGetOikosService(t *testing.T) {
	t.Log("Testing EndPoint OikosService")
	t.Log(parameters.OIKOS_CRUD)
}
