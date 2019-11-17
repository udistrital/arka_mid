package ubicacionHelper_test

import (
	"flag"
	"os"
	"testing"

	"github.com/astaxie/beego"

	"github.com/udistrital/arka_mid/helpers/ubicacionHelper"
)

var parameters struct {
	OIKOS_CRUD string
}

func TestMain(m *testing.M) {
	// ACTA_RECIBIDO_CRUD=pruebasapi.intranetoas.udistrital.edu.co:8206/v1/ OIKOS_CRUD=pruebasapi.intranetoas.udistrital.edu.co:8087/v1/ ENTRADAS_CRUD=pruebasapi.intranetoas.udistrital.edu.co:8207/v1/ PARAMETROS_GOBIERNO_CRUD=pruebasapi.intranetoas.udistrital.edu.co:8205/v1/ OIKOS_CRUD=pruebasapi.intranetoas.udistrital.edu.co:8087/v1/ ADMINISTRATIVA_SERVICE=pruebasapi.intranetoas.udistrital.edu.co:8104/v1/ CUENTAS_CONTABLES_SERVICE=10.20.2.143:8089/v1/ go test ./... -v
	parameters.OIKOS_CRUD = os.Getenv("OIKOS_CRUD")
	beego.AppConfig.Set("oikosService", os.Getenv("OIKOS_CRUD"))
	flag.Parse()
	os.Exit(m.Run())
}

// GetUbicacion ...
func TestGetUbicacion(t *testing.T) {
	valor, err := ubicacionHelper.GetUbicacion(1)
	if err != nil {
		t.Error("No se pudo consultar la ubicacion", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetUbicacion Finalizado Correctamente (OK)")
	}
}

func TestEndPointGetOikosService(t *testing.T) {
	t.Log("Testing EndPoint OIKOS_CRUD")
	t.Log(parameters.OIKOS_CRUD)
}
