package catalogoElementosHelper_test

import (
	"flag"
	"os"
	"testing"

	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers/catalogoElementosHelper"
)

var parameters struct {
	CATALOGO_ELEMENTOS_SERVICE string
}

func TestMain(m *testing.M) {
	parameters.CATALOGO_ELEMENTOS_SERVICE = os.Getenv("CATALOGO_ELEMENTOS_SERVICE")
	beego.AppConfig.Set("catalogoElementosService", os.Getenv("CATALOGO_ELEMENTOS_SERVICE"))
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

func TestEndPointGetCatalogoElementosCrud(t *testing.T) {
	t.Log("Testing EndPoint CATALOGO_ELEMENTOS_SERVICE")
	t.Log(parameters.CATALOGO_ELEMENTOS_SERVICE)
}
