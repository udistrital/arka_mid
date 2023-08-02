package entradaHelper_test

import (
	"flag"
	"os"
	"testing"

	"github.com/astaxie/beego"
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

func TestEndPointGetMovimientos_Arka_Service_Crud(t *testing.T) {
	t.Log("Testing EndPoint MOVIMIENTOS_ARKA_SERVICE")
	t.Log(parameters.MOVIMIENTOS_ARKA_SERVICE)
}
