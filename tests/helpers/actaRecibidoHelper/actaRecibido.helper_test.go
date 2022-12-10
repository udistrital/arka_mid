package actaRecibidoHelper_test

import (
	"flag"
	"mime/multipart"
	"os"
	"testing"

	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
)

var parameters struct {
	ACTA_RECIBIDO_CRUD string
}

// TestMain
//
// Ya estaba! Se desconoce su proposito, se podría eliminar...?
func TestMain(m *testing.M) {
	parameters.ACTA_RECIBIDO_CRUD = os.Getenv("ACTA_RECIBIDO_CRUD")
	if err := beego.AppConfig.Set("ActaRecibidoService", os.Getenv("ACTA_RECIBIDO_CRUD")); err != nil {
		panic(err)
	}
	flag.Parse()
	os.Exit(m.Run())
}

// ARCHIVO: actaRecibdo.helper.go

// TestGetAllActasRecibidoActivas ...
func TestGetAllActasRecibidoActivas(t *testing.T) {
	valor, err := actaRecibido.GetAllActasRecibidoActivas([]string{"Aceptada"}, "ADMIN_ARKA", 1, 0)
	if err != nil {
		t.Error("No se pudo consultar las actas de recibido - err:", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetAllActasRecibidoActivas Finalizado Correctamente (OK)")
	}
}

// TestGetAllParametrosActa ...
func TestGetAllParametrosActa(t *testing.T) {
	valor, err := actaRecibido.GetAllParametrosActa()
	if err != nil {
		t.Error("No se pudo traer los parametros - err:", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetAllParametrosActa Finalizado Correctamente (OK)")
	}
}

// TestDecodeXlsx2Json ...
// TODO: CONVERTIR A PRUEBA UNITARIA!
func TestDecodeXlsx2Json(t *testing.T) {
	// TODO: Traer, de alguna manera, la plantilla que está en Nuxeo
	// y ubicarla de alguna manera en la siguiente variable:
	var file multipart.File

	valor, err := actaRecibido.DecodeXlsx2Json(file)
	if err != nil {
		t.Error("No se pudo procesar la plantilla - err:", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestDecodeXlsx2Json Finalizado Correctamente (OK)")
	}
}

// TestGetElementos ...
func TestGetElementos(t *testing.T) {
	id := 14
	valor, err := actaRecibido.GetElementos(id, nil)
	if err != nil {
		t.Error("No se pudo consultar los elementos del acta de recibido", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetElementos Finalizado Correctamente (OK)")
	}
}

// TestGetIdElementoPlaca ...
func TestGetIdElementoPlaca(t *testing.T) {
	placa := "2021"
	valor, err := actaRecibido.GetIdElementoPlaca(placa)
	if err != nil {
		t.Error("No se pudo consultar el id a partir de la placa - err:", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetIdElementoPlaca Finalizado Correctamente (OK)")
	}
}

// TestEndPointACTA_RECIBIDO_CRUD
//
// Ya estaba! Se desconoce su proposito, se podría eliminar?
func TestEndPointACTA_RECIBIDO_CRUD(t *testing.T) {
	t.Log("Testing EndPoint ACTA_RECIBIDO_CRUD")
	t.Log(parameters.ACTA_RECIBIDO_CRUD)
}
