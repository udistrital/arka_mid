package catalogoElementosHelper

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/cuentasContablesHelper"

	"github.com/udistrital/arka_mid/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
)

// GetCatalogoById ...
func GetCatalogoById(catalogoId int) (catalogo *[]map[string]interface{}, outputError map[string]interface{}) {
	var (
		urlcrud string
	)

	urlcrud = "http://" + beego.AppConfig.String("catalogoElementosService") + "tr_catalogo/" + strconv.Itoa(int(catalogoId))

	if response, err := request.GetJsonTest(urlcrud, &catalogo); err == nil { // (2) error servicio caido
		if response.StatusCode == 200 { // (3) error estado de la solicitud
			logs.Debug(catalogo)
			return catalogo, nil
		} else {
			logs.Info("Error (3) estado de la solicitud")
			outputError = map[string]interface{}{"Function": "GetCatalogoById:GetCatalogoById", "Error": response.Status}
			return nil, outputError
		}
	} else {
		logs.Info("Error (2) servicio caido")
		outputError = map[string]interface{}{"Function": "GetCatalogoById", "Error": err}
		return nil, outputError
	}
}

// GetCuentasContablesGrupo ...
func GetCuentasContablesSubgrupo(subgrupoId int) (cuentasSubgrupoTransaccion []models.CuentasGrupoTransaccion, outputError map[string]interface{}) {
	var (
		urlcrud         string
		cuentasSubgrupo []*models.CuentasGrupo
		cuentaCredito   *models.CuentaContable
		cuentaDebito    *models.CuentaContable
	)

	urlcrud = "http://" + beego.AppConfig.String("catalogoElementosService") + "cuentas_grupo?query=SubgrupoId.Id:" + strconv.Itoa(int(subgrupoId)) + ",Activo:True&limit=-1"

	if response, err := request.GetJsonTest(urlcrud, &cuentasSubgrupo); err == nil { // (2) error servicio caido
		if response.StatusCode == 200 { // (3) error estado de la solicitud
			for _, cuenta := range cuentasSubgrupo {
				cuentaCredito, outputError = cuentasContablesHelper.GetCuentaContable(cuenta.CuentaCreditoId)
				cuentaDebito, outputError = cuentasContablesHelper.GetCuentaContable(cuenta.CuentaDebitoId)

				cuentaContableAux := models.CuentasGrupoTransaccion{
					Id:                  cuenta.Id,
					CuentaCreditoId:     cuentaCredito,
					CuentaDebitoId:      cuentaDebito,
					SubtipoMovimientoId: cuenta.SubtipoMovimientoId,
					FechaCreacion:       cuenta.FechaCreacion,
					FechaModificacion:   cuenta.FechaModificacion,
					Activo:              cuenta.Activo,
					SubgrupoId:          cuenta.SubgrupoId,
				}

				cuentasSubgrupoTransaccion = append(cuentasSubgrupoTransaccion, cuentaContableAux)
			}

			return cuentasSubgrupoTransaccion, nil
		} else {
			logs.Info("Error (3) estado de la solicitud")
			outputError = map[string]interface{}{"Function": "GetCuentasContablesGrupo:GetCuentasContablesGrupo", "Error": response.Status}
			return nil, outputError
		}
	} else {
		logs.Info("Error (2) servicio caido")
		outputError = map[string]interface{}{"Function": "GetCuentasContablesGrupo", "Error": err}
		return nil, outputError
	}
}
