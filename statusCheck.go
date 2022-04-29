package main

import (
	"fmt"
	"net/http"

	"github.com/udistrital/arka_mid/helpers/configuracion"
	e "github.com/udistrital/utils_oas/errorctrl"
)

func statusCheck() interface{} {
	const funcion = "statusCheck - "
	defer e.ErrorControlFunction(funcion+"unhandled error!", fmt.Sprint(http.StatusInternalServerError))

	configuracion.ActualizaRolesArka()
	configuracion.ActualizaTiposDeComprobante()
	return nil
}
