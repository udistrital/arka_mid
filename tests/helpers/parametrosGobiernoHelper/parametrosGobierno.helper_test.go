package parametrosGobiernoHelper_test

import (
	"flag"
	"os"
	"testing"

	"github.com/udistrital/arka_mid/helpers/parametrosGobiernoHelper"
)

var parameters struct {
	GetIva string
}

func TestMain(m *testing.M) {
	parameters.GetIva = os.Getenv("GetIva")
	flag.Parse()
	os.Exit(m.Run())
}

// GetIva ...
func TestGetIva(t *testing.T) {
	valor, err := parametrosGobiernoHelper.GetIva(1)
	if err != nil {
		t.Error("No se pudo consultar el IVA", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetIva Finalizado Correctamente (OK)")
	}
}

func TestEndPointGetIva(t *testing.T) {
	t.Log("Testing EndPoint GetIva")
	t.Log(parameters.GetIva)
}
