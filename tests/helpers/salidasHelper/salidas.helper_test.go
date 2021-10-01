package salidaHelper_test

import (
	"flag"
	"os"
	"testing"

	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers/salidaHelper"
)

var parameters struct {
	MOVIMIENTOS_ARKA_SERVICE   string
	ACTA_RECIBIDO_CRUD         string
	CATALOGO_ELEMENTOS_SERVICE string
	TERCEROS_SERVICE           string
	OIKOS2_CRUD                string
}

func TestMain(m *testing.M) {
	parameters.MOVIMIENTOS_ARKA_SERVICE = os.Getenv("MOVIMIENTOS_ARKA_SERVICE")
	beego.AppConfig.Set("movimientosArkaService", parameters.MOVIMIENTOS_ARKA_SERVICE)
	parameters.ACTA_RECIBIDO_CRUD = os.Getenv("ACTA_RECIBIDO_CRUD")
	beego.AppConfig.Set("actaRecibidoService", parameters.ACTA_RECIBIDO_CRUD)
	parameters.CATALOGO_ELEMENTOS_SERVICE = os.Getenv("CATALOGO_ELEMENTOS_SERVICE")
	beego.AppConfig.Set("catalogoElementosService", parameters.CATALOGO_ELEMENTOS_SERVICE)
	parameters.TERCEROS_SERVICE = os.Getenv("TERCEROS_SERVICE")
	beego.AppConfig.Set("tercerosService", parameters.TERCEROS_SERVICE)
	parameters.OIKOS2_CRUD = os.Getenv("OIKOS2_CRUD")
	beego.AppConfig.Set("oikos2Service", parameters.OIKOS2_CRUD)
	flag.Parse()
	os.Exit(m.Run())
}

// GetAsignacionSedeDependencia ...
func TestGetSalida(t *testing.T) {

	if valor, err := salidaHelper.GetSalida(319); err != nil {
		t.Error("No se pudo consultar la salida", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetSalida Finalizado Correctamente")
	}
}

// GetAsignacionSedeDependencia ...
func TestGetSalidas(t *testing.T) {

	if salidas, err := salidaHelper.GetSalidas(false); err != nil {
		t.Error("No se pudo consultar las salidas", err)
		t.Fail()
	} else {
		t.Log(len(salidas))
		t.Log("TestGetSalidas Finalizado Correctamente")
	}
}

func TestEndPointGetOikosService(t *testing.T) {
	t.Log("Testing EndPoint MOVIMIENTOS_ARKA_SERVICE")
	t.Log(parameters.MOVIMIENTOS_ARKA_SERVICE)
	t.Log("Testing EndPoint ACTA_RECIBIDO_CRUD")
	t.Log(parameters.ACTA_RECIBIDO_CRUD)
	t.Log("Testing EndPoint CATALOGO_ELEMENTOS_SERVICE")
	t.Log(parameters.CATALOGO_ELEMENTOS_SERVICE)
	t.Log("Testing EndPoint TERCEROS_SERVICE")
	t.Log(parameters.TERCEROS_SERVICE)
	t.Log("Testing EndPoint OIKOS2_CRUD")
	t.Log(parameters.OIKOS2_CRUD)
}
