package ubicacionHelper_test

import (
	"flag"
	"os"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/helpers/crud/oikos"
)

var parameters struct {
	OIKOS_CRUD string
}

func TestMain(m *testing.M) {
	parameters.OIKOS_CRUD = os.Getenv("OIKOS_CRUD")
	if err := beego.AppConfig.Set("oikosService", parameters.OIKOS_CRUD); err != nil {
		panic(err)
	}
	flag.Parse()
	os.Exit(m.Run())
}

// GetAsignacionSedeDependencia ...
func TestGetSedeDependenciaUbicacion(t *testing.T) {

	if s, err := oikos.GetSedeDependenciaUbicacion(2); err != nil {
		t.Error("No se pudo consultar la ubicacion", err)
		t.Fail()
	} else {
		t.Log(s)
		t.Log("TestGetSedeDependenciaUbicacion Finalizado Correctamente")
	}
}

func TestEndPointGetOikosService(t *testing.T) {
	t.Log("Testing EndPoint OIKOS_CRUD")
	t.Log(parameters.OIKOS_CRUD)
}
