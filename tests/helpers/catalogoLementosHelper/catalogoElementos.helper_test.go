package catalogoElementosHelper_test

import (
	"flag"
	"os"
	"testing"

	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers/catalogoElementosHelper"
	"github.com/udistrital/arka_mid/models"
)

var parameters struct {
	CATALOGO_ELEMENTOS_SERVICE string
}

func TestMain(m *testing.M) {
	parameters.CATALOGO_ELEMENTOS_SERVICE = os.Getenv("CATALOGO_ELEMENTOS_SERVICE")
	if err := beego.AppConfig.Set("catalogoElementosService", os.Getenv("CATALOGO_ELEMENTOS_SERVICE")); err != nil {
		panic(err)
	}
	flag.Parse()
	os.Exit(m.Run())
}

// GetCuentasContablesGrupo ...
func TestGetCuentasContablesSubgrupo(t *testing.T) {
	var ctas []models.DetalleCuentasSubgrupo
	err := catalogoElementosHelper.GetCuentasContablesSubgrupo(1, 1, &ctas)
	if err != nil || len(ctas) == 0 {
		if err != nil {
			t.Error("No se pudo consultar las cuentas contables del subgrupo", err)
		} else {
			t.Error("No se pudo consultar las cuentas contables del subgrupo", err)
		}

		t.Fail()
	} else {
		t.Log("TestGetCuentasContablesGrupo Finalizado Correctamente (OK)")
	}
}

func TestEndPointGetCatalogoElementosCrud(t *testing.T) {
	t.Log("Testing EndPoint CATALOGO_ELEMENTOS_SERVICE")
	t.Log(parameters.CATALOGO_ELEMENTOS_SERVICE)
}
