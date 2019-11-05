package proveedorHelper_test

import (
	"flag"
	"os"
	"testing"

	"github.com/udistrital/arka_mid/helpers/proveedorHelper"
)

var parameters struct {
	GetProveedorById string
}

func TestMain(m *testing.M) {
	parameters.GetProveedorById = os.Getenv("GetProveedorById")
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

func TestEndPointGetProveedorById(t *testing.T) {
	t.Log("Testing EndPoint GetProveedorById")
	t.Log(parameters.GetProveedorById)
}
