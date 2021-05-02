package ubicacionHelper_test

import (
	"flag"
	"os"
	"testing"

	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers/ubicacionHelper"
)

var parameters struct {
	OIKOS2_CRUD string
}

func TestMain(m *testing.M) {
	parameters.OIKOS2_CRUD = os.Getenv("OIKOS2_CRUD")
	beego.AppConfig.Set("oikos2Service", parameters.OIKOS2_CRUD)
	flag.Parse()
	os.Exit(m.Run())
}

// GetAsignacionSedeDependencia ...
func TestGetAsignacionSedeDependencia(t *testing.T) {

	if valor, err := ubicacionHelper.GetAsignacionSedeDependencia("2"); err != nil {
		t.Error("No se pudo consultar la ubicacion", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetAsignacionSedeDependencia Finalizado Correctamente")
	}
}

// GetAsignacionSedeDependencia ...
func TestGetSedeDependenciaUbicacion(t *testing.T) {

	if s, d, u, err := ubicacionHelper.GetSedeDependenciaUbicacion("2"); err != nil {
		t.Error("No se pudo consultar la ubicacion", err)
		t.Fail()
	} else {
		t.Log(s, d, u)
		t.Log("TestGetSedeDependenciaUbicacion Finalizado Correctamente")
	}
}

func TestEndPointGetOikosService(t *testing.T) {
	t.Log("Testing EndPoint OIKOS_CRUD")
	t.Log(parameters.OIKOS2_CRUD)
}
