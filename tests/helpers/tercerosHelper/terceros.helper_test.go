package tercerosHelper_test

import (
	"flag"
	"os"
	"testing"

	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/arka_mid/helpers/crud/terceros"
)

var parameters struct {
	PARAMETROS_CRUD  string
	TERCEROS_SERVICE string
	OIKOS_CRUD       string
}

func TestMain(m *testing.M) {
	parameters.PARAMETROS_CRUD = os.Getenv("PARAMETROS_CRUD")
	if err := beego.AppConfig.Set("parametrosService", parameters.PARAMETROS_CRUD); err != nil {
		panic(err)
	}
	parameters.TERCEROS_SERVICE = os.Getenv("TERCEROS_SERVICE")
	if err := beego.AppConfig.Set("tercerosService", parameters.TERCEROS_SERVICE); err != nil {
		panic(err)
	}
	parameters.OIKOS_CRUD = os.Getenv("OIKOS_CRUD")
	if err := beego.AppConfig.Set("oikosService", parameters.OIKOS_CRUD); err != nil {
		panic(err)
	}
	flag.Parse()
	os.Exit(m.Run())
}

// TestGetFuncionariosPlanta ...
func TestGetNombreTerceroById(t *testing.T) {

	if valor, err := terceros.GetNombreTerceroById(81); err != nil {
		t.Error("No se pudo consultar el tercero", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetNombreTerceroById Finalizado Correctamente")
	}
}

// TestGetFuncionariosPlanta ...
func TestGetTerceroByUsuarioWSO2(t *testing.T) {

	if valor, err := terceros.GetTerceroByUsuarioWSO2("utest01"); err != nil {
		t.Error("No se pudo consultar el tercero", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetTerceroByUsuarioWSO2 Finalizado Correctamente")
	}
}

// TestGetTerceroByDoc ...
func TestGetTerceroByDoc(t *testing.T) {

	if valor, err := terceros.GetTerceroByDoc("80000000"); err != nil {
		t.Error("No se pudo consultar el tercero", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetTerceroByDoc Finalizado Correctamente")
	}
}

func TestEndPointParametrosService(t *testing.T) {
	t.Log("Testing EndPoint parametrosService")
	t.Log(parameters.PARAMETROS_CRUD)
	t.Log("Testing EndPoint tercerosService")
	t.Log(parameters.TERCEROS_SERVICE)
	t.Log("Testing EndPoint oikosService")
	t.Log(parameters.OIKOS_CRUD)
}
