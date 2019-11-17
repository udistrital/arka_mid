package actaRecibidoHelper_test

import (
	"flag"
	"os"
	"testing"

	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers/actaRecibidoHelper"
)

var parameters struct {
	ACTA_RECIBIDO_CRUD string
}

func TestMain(m *testing.M) {
	parameters.ACTA_RECIBIDO_CRUD = os.Getenv("ACTA_RECIBIDO_CRUD")
	beego.AppConfig.Set("ActaRecibidoService", os.Getenv("ACTA_RECIBIDO_CRUD"))
	flag.Parse()
	os.Exit(m.Run())
}

// GetAllActasRecibido ...
func TestGetAllActasRecibido(t *testing.T) {
	valor, err := actaRecibidoHelper.GetAllActasRecibido()
	if err != nil {
		t.Error("No se pudo consultar las actas de recibido", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetAllActasRecibido Finalizado Correctamente (OK)")
	}
}

// GetActasRecibidoTipo ...
func TestGetActasRecibidoTipo(t *testing.T) {
	valor, err := actaRecibidoHelper.GetActasRecibidoTipo(5)
	if err != nil {
		t.Error("No se pudo consultar las actas de recibido por tipo", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetActasRecibidoTipo Finalizado Correctamente (OK)")
	}
}

// GetElementos ...
func TestGetElementos(t *testing.T) {
	valor, err := actaRecibidoHelper.GetElementos(14)
	if err != nil {
		t.Error("No se pudo consultar los elementos del acta de recibido", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetElementos Finalizado Correctamente (OK)")
	}
}

func TestEndPointACTA_RECIBIDO_CRUD(t *testing.T) {
	t.Log("Testing EndPoint ACTA_RECIBIDO_CRUD")
	t.Log(parameters.ACTA_RECIBIDO_CRUD)
}
