package bajasHelper

import (
	"encoding/json"
	"net/url"

	"github.com/astaxie/beego/logs"

	crud_actas "github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	crudMovimientosArka "github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
)

// AprobarBajas Aprobación masiva de bajas, transacción contable y actualización de movmientos
func AprobarBajas(data *models.TrRevisionBaja) (ids []int, outputError map[string]interface{}) {

	funcion := "AprobarBajas"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		bajas           []*models.Movimiento
		idsMov          []int
		idsActa         []int
		idsSubgrupos    []int
		elementosMov    []*models.ElementosMovimiento
		elementosActa   []*models.Elemento
		cuentasSubgrupo []*models.CuentaSubgrupo
	)

	// Paso 1: Transacción contable
	query := "fields=Detalle&limit=-1&query=Id__in:"
	query += url.QueryEscape(utilsHelper.ArrayToString(data.Bajas, "|"))
	if bajas_, err := crudMovimientosArka.GetAllMovimiento(query); err != nil {
		return nil, err
	} else {
		bajas = bajas_
	}

	for _, mov := range bajas {

		var detalle *models.FormatoBaja

		if err := json.Unmarshal([]byte(mov.Detalle), &detalle); err != nil {
			logs.Error(err)
			eval := " - json.Unmarshal([]byte(mov.Detalle), &detalle)"
			return nil, errorctrl.Error(funcion+eval, err, "500")
		}

		idsMov = append(idsMov, detalle.Elementos...)
	}

	query = "fields=ElementoActaId&limit=-1&query=Id__in:"
	query += url.QueryEscape(utilsHelper.ArrayToString(idsMov, "|"))
	if elementos_, err := crudMovimientosArka.GetAllElementosMovimiento(query); err != nil {
		return nil, err
	} else {
		elementosMov = elementos_
	}

	for _, el := range elementosMov {
		idsActa = append(idsActa, el.ElementoActaId)
	}

	query = "Id__in:" + utilsHelper.ArrayToString(idsActa, "|")
	if elementos_, err := crud_actas.GetAllElemento(query, "SubgrupoCatalogoId", "", "", "", "-1"); err != nil {
		return nil, err
	} else {
		elementosActa = elementos_
	}

	for _, el := range elementosActa {
		idsSubgrupos = append(idsSubgrupos, el.SubgrupoCatalogoId)
	}

	query = "limit=-1&fields=CuentaDebitoId,CuentaCreditoId,SubgrupoId&sortby=Id&order=desc&"
	query += "query=SubtipoMovimientoId:32,Activo:true,SubgrupoId__Id__in:"
	query += url.QueryEscape(utilsHelper.ArrayToString(idsSubgrupos, "|"))
	if elementos_, err := catalogoElementos.GetAllCuentasSubgrupo(query); err != nil {
		return nil, err
	} else {
		cuentasSubgrupo = elementos_
	}

	infoCuentas := make(map[int]*InfoCuentasSubgrupos)
	for _, idSubgrupo := range idsSubgrupos {

		var (
			ctaCr *models.CuentaContable
			ctaDb *models.CuentaContable
		)

		if idx := FindInArray(cuentasSubgrupo, idSubgrupo); idx > -1 {
			if ctaCr_, err := cuentasContables.GetCuentaContable(cuentasSubgrupo[idx].CuentaCreditoId); err != nil {
				return nil, err
			} else {
				if err := formatdata.FillStruct(ctaCr_, &ctaCr); err != nil {
					logs.Error(err)
					eval := " - formatdata.FillStruct(ctaCr_, &ctaCr)"
					return nil, errorctrl.Error(funcion+eval, err, "500")
				}
			}

			if ctaDb_, err := cuentasContables.GetCuentaContable(cuentasSubgrupo[idx].CuentaDebitoId); err != nil {
				return nil, err
			} else {
				if err := formatdata.FillStruct(ctaDb_, &ctaDb); err != nil {
					logs.Error(err)
					eval := " - formatdata.FillStruct(ctaDb_, &ctaDb)"
					return nil, errorctrl.Error(funcion+eval, err, "500")
				}
			}

			infoCuentas[idSubgrupo] = new(InfoCuentasSubgrupos)
			infoCuentas[idSubgrupo].CuentaCredito = ctaCr
			infoCuentas[idSubgrupo].CuentaDebito = ctaDb
		} else {
			// Se llega acá cuando la clase de algún elemento no tiene la cuenta contable registrada
			// Queda pendiente qué se debe mostrar al usuario en ese caso
		}

	}

	// Paso 2: Actualiza el estado de las bajas en api movimientos_arka_crud
	if ids_, err := crudMovimientosArka.PutRevision(data); err != nil {
		return nil, err
	} else {
		ids = ids_
	}

	return ids, nil
}
