package configuracion

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/crud/configuracion"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	modelsConfiguracion "github.com/udistrital/configuracion_api/models"
	e "github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
)

// ActualizaRolesArka se puede llamar periodicamente. Un candidato podría ser
// el healthcheck
func ActualizaRolesArka() {
	const funcion = "ActualizaRolesArka - "
	defer e.ErrorControlFunction(funcion+"unhandled error", fmt.Sprint(http.StatusInternalServerError))

	// parametro de roles registrados
	getParametrosArka("RolesRegistrados", 1, &roles)
}

// ActualizaTiposDeComprobante carga los tipos de comprobante
func ActualizaTiposDeComprobante() {
	const funcion = "ActualizaTiposDeComprobante - "
	defer e.ErrorControlFunction(funcion+"unhandled error!", fmt.Sprint(http.StatusInternalServerError))

	getParametrosArka("TiposDeComprobante", 1, &comprobantes)
}

func getParametrosArka(parametro string, resultadosEsperados uint, out interface{}) {
	const funcion = "parametroArka - "
	defer e.ErrorControlFunction(funcion+"unhandled error!", fmt.Sprint(http.StatusInternalServerError))

	var parametros []modelsConfiguracion.Parametro
	query := utilsHelper.Query{
		Query: map[string]string{
			"Aplicacion__Nombre": beego.AppConfig.String("nombreAplicacion"),
			"Nombre":             parametro,
		},
		Limit: -1,
	}
	if err := configuracion.GetParametros(query, &parametros); err != nil {
		logs.Critical(err)
		panic(err)
	}
	if len(parametros) != 1 {
		cond := ""
		if len(parametros) >= 10 {
			cond = " (o más)"
		}
		err := fmt.Errorf("se esperaba encontrar %d registro(s) con Nombre:%s en configuracion_crud/parametros, hay: %d%s",
			resultadosEsperados, parametro, len(parametros), cond)
		logs.Critical(err)
		panic(err)
	}
	if err := formatdata.FillStruct(parametros[0].Valor, &out); err != nil {
		logs.Critical(err)
		panic(err)
	}
}
