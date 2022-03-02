package parametrosGobiernoHelper_test

import (
	"flag"
	"os"
	"testing"

	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers/crud/parametrosGobierno"
)

var parameters struct {
	PARAMETROS_GOBIERNO_CRUD string
}

func TestMain(m *testing.M) {
	parameters.PARAMETROS_GOBIERNO_CRUD = os.Getenv("PARAMETROS_GOBIERNO_CRUD")
	beego.AppConfig.Set("parametrosGobiernoService", os.Getenv("PARAMETROS_GOBIERNO_CRUD"))
	flag.Parse()
	os.Exit(m.Run())
}

// GetIva ...
func TestGetIva(t *testing.T) {
	valor, err := parametrosGobierno.GetIva(1)
	if err != nil {
		t.Error("No se pudo consultar el IVA", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetIva Finalizado Correctamente (OK)")
	}
}

func TestEndPointGetParametrosGobiernoService(t *testing.T) {
	t.Log("Testing EndPoint PARAMETROS_GOBIERNO_CRUD")
	t.Log(parameters.PARAMETROS_GOBIERNO_CRUD)
}
