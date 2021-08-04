package tercerosHelper_test

import (
	"flag"
	"os"
	"testing"

	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers/tercerosHelper"
)

var parameters struct {
	PARAMETROS_CRUD  string
	TERCEROS_SERVICE string
	OIKOS2_CRUD      string
}

func TestMain(m *testing.M) {
	parameters.PARAMETROS_CRUD = os.Getenv("PARAMETROS_CRUD")
	beego.AppConfig.Set("parametrosService", parameters.PARAMETROS_CRUD)
	parameters.TERCEROS_SERVICE = os.Getenv("TERCEROS_SERVICE")
	beego.AppConfig.Set("tercerosService", parameters.TERCEROS_SERVICE)
	parameters.OIKOS2_CRUD = os.Getenv("OIKOS2_CRUD")
	beego.AppConfig.Set("oikos2Service", parameters.OIKOS2_CRUD)
	flag.Parse()
	os.Exit(m.Run())
}

// TestGetContratista ...
func TestGetContratista(t *testing.T) {

	if valor, err := tercerosHelper.GetContratista(9825); err != nil {
		t.Error("No se pudo consultar el contratista", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetContratista Finalizado Correctamente")
	}
}

// TestGetFuncionariosPlanta ...
func TestGetFuncionariosPlanta(t *testing.T) {

	if valor, err := tercerosHelper.GetFuncionariosPlanta(9801); err != nil {
		t.Error("No se pudo consultar el funcionario", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetFuncionariosPlanta Finalizado Correctamente")
	}
}

// TestGetFuncionariosPlanta ...
func TestGetOrdenadores(t *testing.T) {

	if valor, err := tercerosHelper.GetOrdenadores(9804); err != nil {
		t.Error("No se pudo consultar el ordenador", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetOrdenadores Finalizado Correctamente")
	}
}

// TestGetFuncionariosPlanta ...
func TestGetProveedor(t *testing.T) {

	if valor, err := tercerosHelper.GetProveedor(9769); err != nil {
		t.Error("No se pudo consultar el proveedor", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetProveedor Finalizado Correctamente")
	}
}

// TestGetFuncionariosPlanta ...
func TestGetNombreTerceroById(t *testing.T) {

	if valor, err := tercerosHelper.GetNombreTerceroById("81"); err != nil {
		t.Error("No se pudo consultar el tercero", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetNombreTerceroById Finalizado Correctamente")
	}
}

// TestGetFuncionariosPlanta ...
func TestGetTerceroByUsuarioWSO2(t *testing.T) {

	if valor, err := tercerosHelper.GetTerceroByUsuarioWSO2("utest01"); err != nil {
		t.Error("No se pudo consultar el tercero", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetTerceroByUsuarioWSO2 Finalizado Correctamente")
	}
}

// TestGetTerceroByDoc ...
func TestGetTerceroByDoc(t *testing.T) {

	if valor, err := tercerosHelper.GetTerceroByDoc("80000000"); err != nil {
		t.Error("No se pudo consultar el tercero", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetTerceroByDoc Finalizado Correctamente")
	}
}

// TestGetTipos ...
func TestGetTipos(t *testing.T) {

	if valor, err := tercerosHelper.GetTipos(); err != nil {
		t.Error("No se pudo consultar los tipos", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetTipos Finalizado Correctamente")
	}
}

func TestEndPointParametrosService(t *testing.T) {
	t.Log("Testing EndPoint parametrosService")
	t.Log(parameters.PARAMETROS_CRUD)
	t.Log("Testing EndPoint tercerosService")
	t.Log(parameters.TERCEROS_SERVICE)
	t.Log("Testing EndPoint oikos2Service")
	t.Log(parameters.OIKOS2_CRUD)
}
