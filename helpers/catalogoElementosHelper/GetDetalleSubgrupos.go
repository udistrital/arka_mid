package catalogoElementosHelper

import (
	"net/url"

	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// GetCuentasByMovimientoSubgrupos Consulta las cuentas para una serie de subgrupos y las almacena en una estructura de fácil acceso
func GetDetalleSubgrupos(subgrupos []int, detalleSubgrupo map[int]models.DetalleSubgrupo) (
	outputError map[string]interface{}) {

	funcion := "GetDetalleSubgrupos"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	query := "limit=-1&sortby=Id&order=desc&fields=Amortizacion,Depreciacion,SubgrupoId&query=Activo:true"
	query += ",SubgrupoId__Id__in:" + url.QueryEscape(utilsHelper.ArrayToString(subgrupos, "|"))
	if detalles, err := catalogoElementos.GetAllDetalleSubgrupo(query); err != nil {
		return err
	} else {
		for _, detalle := range detalles {
			detalleSubgrupo[detalle.SubgrupoId.Id] = *detalle
		}

	}

	return

}
