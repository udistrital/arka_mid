package catalogoElementosHelper

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// GetCuentasContablesSubgrupo ...
func GetCuentasContablesSubgrupo(subgrupoId, movimientoId int, cuentas *[]models.DetalleCuentasSubgrupo) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("GetCuentasContablesSubgrupo - Unhandled Error!", "500")

	query := "limit=-1&sortby=CodigoAbreviacion&order=asc&fields=Id,CodigoAbreviacion,Nombre&query=Activo:true"
	tipos, outputError := movimientosArka.GetAllFormatoTipoMovimiento(query)
	if outputError != nil {
		return
	}

	ctas, outputError := consultarCuentasSubgrupoRecientes(subgrupoId, movimientoId)
	if outputError != nil {
		return
	}

	formatosPorID := indexarFormatosPorID(tipos)
	detalleCtas := make(map[string]models.DetalleCuenta)
	for _, cta := range ctas {
		detalle := models.DetalleCuentasSubgrupo{
			Id:         cta.Id,
			SubgrupoId: subgrupoId,
		}
		if cta.TipoBienId != nil {
			detalle.TipoBienId = *cta.TipoBienId
		}

		detalle.TipoMovimientoId = formatoMovimientoByID(cta.TipoMovimientoId, formatosPorID)
		detalle.SubtipoMovimientoId = formatoMovimientoByID(cta.SubtipoMovimientoId, formatosPorID)
		detalle.CuentaCreditoId = new(models.DetalleCuenta)
		detalle.CuentaDebitoId = new(models.DetalleCuenta)

		outputError = findCuentaSubgrupo(detalle.CuentaCreditoId, cta.CuentaCreditoId, detalleCtas)
		if outputError != nil {
			return
		}

		outputError = findCuentaSubgrupo(detalle.CuentaDebitoId, cta.CuentaDebitoId, detalleCtas)
		if outputError != nil {
			return
		}

		*cuentas = append(*cuentas, detalle)
	}

	return
}

func consultarCuentasSubgrupoRecientes(subgrupoId, movimientoId int) (ctas []models.CuentasSubgrupo, outputError map[string]interface{}) {
	query := "limit=-1&fields=Id,CuentaDebitoId,CuentaCreditoId,TipoMovimientoId,SubtipoMovimientoId,SubgrupoId,TipoBienId" +
		"&sortby=Id&order=desc&query=Activo:true,SubgrupoId__Id:" + strconv.Itoa(subgrupoId)
	cuentasSubgrupo, outputError := catalogoElementos.GetAllCuentasSubgrupo(query)
	if outputError != nil {
		return nil, outputError
	}

	return seleccionarCuentasSubgrupoRecientes(cuentasSubgrupo, movimientoId), nil
}

func seleccionarCuentasSubgrupoRecientes(cuentasSubgrupo []*models.CuentasSubgrupo, movimientoId int) []models.CuentasSubgrupo {
	seleccionadas := make([]models.CuentasSubgrupo, 0, len(cuentasSubgrupo))
	vistas := make(map[string]struct{}, len(cuentasSubgrupo))
	for _, cuenta := range cuentasSubgrupo {
		if cuenta == nil || cuenta.TipoBienId == nil {
			continue
		}

		if movimientoId > 0 && cuenta.TipoMovimientoId != movimientoId && cuenta.SubtipoMovimientoId != movimientoId {
			continue
		}

		key := cuentaSubgrupoKey(*cuenta)
		if _, ok := vistas[key]; ok {
			continue
		}

		vistas[key] = struct{}{}
		seleccionadas = append(seleccionadas, *cuenta)
	}

	return seleccionadas
}

func cuentaSubgrupoKey(cuenta models.CuentasSubgrupo) string {
	tipoBienID := 0
	if cuenta.TipoBienId != nil {
		tipoBienID = cuenta.TipoBienId.Id
	}

	return fmt.Sprintf("%d|%d|%d", cuenta.TipoMovimientoId, cuenta.SubtipoMovimientoId, tipoBienID)
}

func indexarFormatosPorID(tipos []*models.FormatoTipoMovimiento) map[int]models.FormatoTipoMovimiento {
	index := make(map[int]models.FormatoTipoMovimiento, len(tipos))
	for _, fm := range tipos {
		if fm == nil || fm.Id <= 0 {
			continue
		}

		index[fm.Id] = *fm
	}

	return index
}

func formatoMovimientoByID(id int, formatos map[int]models.FormatoTipoMovimiento) *models.FormatoTipoMovimiento {
	if id == 0 {
		return &models.FormatoTipoMovimiento{Id: 0}
	}

	if fm, ok := formatos[id]; ok {
		formato := fm
		return &formato
	}

	return &models.FormatoTipoMovimiento{Id: id}
}

func findCuentaSubgrupo(ctaSg *models.DetalleCuenta, cuentaId string, cuentas map[string]models.DetalleCuenta) (outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("findCuentaSubgrupo - Unhandled Error!", "500")

	if val, ok := cuentas[cuentaId]; ok {
		*ctaSg = val
		return
	}

	cta, outputError := cuentasContables.GetCuentaContable(cuentaId)
	if outputError != nil {
		return
	} else if cta != nil {
		var dcta models.DetalleCuenta
		outputError = utilsHelper.FillStruct(cta, &dcta)
		if outputError != nil {
			return
		}

		*ctaSg = dcta
		cuentas[cuentaId] = dcta
	}

	return
}

// GetCuentasByMovimientoSubgrupos Consulta las cuentas para una serie de subgrupos y las almacena en una estructura de fácil acceso
func GetCuentasByMovimientoAndSubgrupos(movimientoId int, subgrupos []int, cuentasSubgrupo map[int]models.CuentasSubgrupo) (
	outputError map[string]interface{}) {

	funcion := "GetCuentasByMovimientoSubgrupos"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var subgrupos_ []int
	for _, sg := range subgrupos {
		if _, ok := cuentasSubgrupo[sg]; !ok {
			subgrupos_ = append(subgrupos_, sg)
		}
	}

	if len(subgrupos_) == 0 {
		return
	}

	query := "limit=-1&fields=CuentaDebitoId,CuentaCreditoId,SubgrupoId&sortby=Id&order=desc&" +
		"query=Activo:true,SubtipoMovimientoId:" + strconv.Itoa(movimientoId) +
		",SubgrupoId__Id__in:" + url.QueryEscape(utilsHelper.ArrayToString(subgrupos_, "|"))
	if cuentas_, err := catalogoElementos.GetAllCuentasSubgrupo(query); err != nil {
		return err
	} else {
		for _, cuenta := range cuentas_ {
			cuentasSubgrupo[cuenta.SubgrupoId.Id] = *cuenta
		}

	}

	return

}
