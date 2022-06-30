// Labores de Procesamiento/Limpieza desde el controlador

package actaRecibido

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/actaRecibido/constantes"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	e "github.com/udistrital/utils_oas/errorctrl"
)

func ProcesaQueryListaActas(in string, out map[string]string) (outputError map[string]interface{}) {
	const funcion = "ProcesaQuery - "
	defer e.ErrorControlFunction(funcion+"unhandled error!", fmt.Sprint(http.StatusBadRequest))

	// Dividir el query
	if err := utilsHelper.QuerySplit(in, out); err != nil {
		logs.Debug(err)
		outputError = e.Error(funcion+"utilsHelper.QuerySplit(v,query)",
			err, fmt.Sprint(http.StatusBadRequest))
		return
	}

	// Validar que solo vengan parametros permitidos
	for k := range out {
		if !constantes.ParametroValidoListaActas(k) {
			err := fmt.Errorf("'%s' nor allowed nor implemented", k)
			outputError = e.Error(funcion+"!actaRecibido.CriterioListaActasPermitido(k)",
				err, fmt.Sprint(http.StatusBadRequest))
			break
		}

		// En este punto el criterio actual es permitido, ahora validar el tipo/formato de datos
		// TODO: Implementar, posiblemente con un switch/case
	}
	return
}
