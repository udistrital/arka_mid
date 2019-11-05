package catalogoElementosHelper_test

import (
	"flag"
	"os"
	"testing"

	"github.com/udistrital/arka_mid/helpers/catalogoElementosHelper"
)

var parameters struct {
	GetCatalogoById          string
	GetCuentasContablesGrupo string
}

func TestMain(m *testing.M) {
	parameters.GetCatalogoById = os.Getenv("GetCatalogoById")
	parameters.GetCuentasContablesGrupo = os.Getenv("GetCuentasContablesGrupo")
	flag.Parse()
	os.Exit(m.Run())
}

// GetCatalogoById ...
func TestGetCatalogoById(t *testing.T) {
	valor, err := catalogoElementosHelper.GetCatalogoById(1)
	if err != nil {
		t.Error("No se pudo consultar el catalogo de elementos", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetCatalogoById Finalizado Correctamente (OK)")
	}
}

// GetCuentasContablesGrupo ...
func TestGetCuentasContablesGrupo(t *testing.T) {
	valor, err := catalogoElementosHelper.GetCuentasContablesSubgrupo(3)
	if err != nil {
		t.Error("No se pudo consultar las cuentas contables del subgrupo", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetCuentasContablesGrupo Finalizado Correctamente (OK)")
	}
}

func TestEndPointGetCatalogoById(t *testing.T) {
	t.Log("Testing EndPoint GetCatalogoById")
	t.Log(parameters.GetCatalogoById)
}

func TestEndPointGetCuentasContablesSubgrupo(t *testing.T) {
	t.Log("Testing EndPoint GetCuentasContablesSubgrupo")
	t.Log(parameters.GetCuentasContablesGrupo)
}
