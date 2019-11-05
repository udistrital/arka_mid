package ubicacionHelper_test

import (
	"flag"
	"os"
	"testing"

	"github.com/udistrital/arka_mid/helpers/ubicacionHelper"
)

var parameters struct {
	GetUbicacion string
}

func TestMain(m *testing.M) {
	parameters.GetUbicacion = os.Getenv("GetUbicacion")
	flag.Parse()
	os.Exit(m.Run())
}

// GetUbicacion ...
func TestGetUbicacion(t *testing.T) {
	valor, err := ubicacionHelper.GetUbicacion(1)
	if err != nil {
		t.Error("No se pudo consultar la ubicacion", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetUbicacion Finalizado Correctamente (OK)")
	}
}

func TestEndPointGetUbicacion(t *testing.T) {
	t.Log("Testing EndPoint GetUbicacion")
	t.Log(parameters.GetUbicacion)
}
