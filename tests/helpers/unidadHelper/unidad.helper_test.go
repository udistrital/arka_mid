package ubicacionHelper_test

import (
	"flag"
	"os"
	"testing"

	"github.com/udistrital/arka_mid/helpers/ubicacionHelper"
)

var parameters struct {
	GetUnidad string
}

func TestMain(m *testing.M) {
	parameters.GetUnidad = os.Getenv("GetUnidad")
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

func TestEndPointGetUnidad(t *testing.T) {
	t.Log("Testing EndPoint GetUnidad")
	t.Log(parameters.GetUnidad)
}
