package proveedorHelper_test

import (
	"flag"
	"os"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

var parameters struct {
	ADMINISTRATIVA_SERVICE string
}

func TestMain(m *testing.M) {
	parameters.ADMINISTRATIVA_SERVICE = os.Getenv("ADMINISTRATIVA_SERVICE")
	if err := beego.AppConfig.Set("administrativaService", os.Getenv("ADMINISTRATIVA_SERVICE")); err != nil {
		panic(err)
	}
	flag.Parse()
	os.Exit(m.Run())
}

func TestEndPointGetAdministrativaService(t *testing.T) {
	t.Log("Testing EndPoint ADMINISTRATIVA_SERVICE")
	t.Log(parameters.ADMINISTRATIVA_SERVICE)
}
