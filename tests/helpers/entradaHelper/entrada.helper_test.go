package entradaHelper_test

import (
	"flag"
	"os"
	"testing"

	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers/entradaHelper"
)

var parameters struct {
	MOVIMIENTOS_ARKA_SERVICE string
}

func TestMain(m *testing.M) {
	parameters.MOVIMIENTOS_ARKA_SERVICE = os.Getenv("MOVIMIENTOS_ARKA_SERVICE")
	if err := beego.AppConfig.Set("movimientosArkaService", os.Getenv("MOVIMIENTOS_ARKA_SERVICE")); err != nil {
		panic(err)
	}
	flag.Parse()
	os.Exit(m.Run())
}

func TestAnularEntrada(t *testing.T) {
	valor, err := entradaHelper.AnularEntrada(1)
	if err != nil || valor == nil {
		if err != nil {
			t.Error("No se pudo anular entrada err", err)
		} else {
			t.Error("No se pudo anular el valor de la entrada", err)
		}
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestAnularEntrada Finalizado Correctamente (OK)")
	}
}

func TestGetMovimientosByActa(t *testing.T) {
	valor, err := entradaHelper.GetMovimientosByActa(2)
	if err != nil || valor == nil {
		if err != nil {
			t.Error("No se pudo anular entrada err", err)
		} else {
			t.Error("No se pudo anular el valor de la entrada", err)
		}
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("GetMovimientosByActa Finalizado Correctamente (OK)")
	}
}

func TestEndPointGetMovimientos_Arka_Service_Crud(t *testing.T) {
	t.Log("Testing EndPoint MOVIMIENTOS_ARKA_SERVICE")
	t.Log(parameters.MOVIMIENTOS_ARKA_SERVICE)
}
