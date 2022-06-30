// Labores de Procesamiento/Limpieza desde el controlador

package actaRecibido

import (
	"fmt"
	"net/http"

	"github.com/udistrital/arka_mid/helpers/actaRecibido/constantes"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	e "github.com/udistrital/utils_oas/errorctrl"
)

func ProcesaQueryListaActas(in string, out map[string]interface{}) (outputError map[string]interface{}) {
	const funcion = "ProcesaQueryListaActas - "
	defer e.ErrorControlFunction(funcion+"unhandled error!", fmt.Sprint(http.StatusBadRequest))

	// Dividir el query
	raw := make(map[string]string)
	if err := utilsHelper.QuerySplit(in, &raw); err != nil {
		outputError = e.Error(funcion+"utilsHelper.QuerySplit(v,query)",
			err, fmt.Sprint(http.StatusBadRequest))
		return
	}

	// Validar que solo vengan parametros permitidos
	for k := range raw {
		if !constantes.ParametroValidoListaActas(k) {
			err := fmt.Errorf("'%s' nor allowed nor implemented", k)
			outputError = e.Error(funcion+"!actaRecibido.CriterioListaActasPermitido(k)",
				err, fmt.Sprint(http.StatusBadRequest))
			break
		}
	}

	// Validar el tipo/formato de datos
	return constantes.ParserParametrosListaActas(raw, out)
}
