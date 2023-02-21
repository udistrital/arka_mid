package inmuebleshelper

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/oikos"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

func GetOne(id int) (detalle models.Inmueble, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GetOne - Unhandled Error!", "500")

	outputError = actaRecibido.GetElementoById(id, &detalle.Elemento)
	if outputError != nil {
		return
	}

	elementoMovimiento, outputError := movimientosArka.GetAllElementosMovimiento(getPayloadElementosMovimiento(id))
	if len(elementoMovimiento) == 1 {
		detalle.ElementoMovimiento = *elementoMovimiento[0]
	}

	elementosCampo, outputError := actaRecibido.GetAllElementoCampo(getPayloadElementoCampo(id))
	if outputError != nil {
		return
	}

	for _, campo := range elementosCampo {
		var cuentas models.ParametrizacionContable
		outputError = utilsHelper.Unmarshal(campo.Valor, &cuentas)
		if outputError != nil {
			return
		}

		var cuentas_ models.ParametrizacionContable_
		cuentaCredito, err := cuentasContables.GetCuentaContable(cuentas.CuentaCreditoId)
		if err != nil {
			outputError = err
			return
		} else if cuentaCredito != nil {
			cuentas_.CuentaCreditoId = *cuentaCredito
		}

		cuentaDebito, err := cuentasContables.GetCuentaContable(cuentas.CuentaDebitoId)
		if err != nil {
			outputError = err
			return
		} else if cuentaDebito != nil {
			cuentas_.CuentaDebitoId = *cuentaDebito
		}

		if campo.CampoId.Sigla == "CC_ENT" {
			detalle.Cuentas = cuentas_
		} else if campo.CampoId.Sigla == "CC_MED" {
			detalle.CuentasMediciones = cuentas_
		}
	}

	if detalle.Elemento.EspacioFisicoId > 0 {
		espacioFisico_, err := oikos.GetAllEspacioFisico(getPayloadEspacioFisico(detalle.Elemento.EspacioFisicoId))
		if err != nil {
			return detalle, err
		}

		if len(espacioFisico_) == 1 {
			detalle.EspacioFisico = espacioFisico_[0]
			detalle.Sede, outputError = oikos.GetSedeEspacioFisico(espacioFisico_[0])
			if outputError != nil {
				return
			}
		}

		detalle.Otros, outputError = oikos.GetAllEspacioFisicoCampo(getPayloadEspacioFisico(id))
	}

	return
}

func getPayloadElementosMovimiento(id int) string {
	return "limit=1&query=Activo:true,ElementoActaId:" + strconv.Itoa(id)
}

func getPayloadElementoCampo(id int) string {
	return "query=Activo:true,CampoId__Sigla__in:CC_ENT|CC_MED,ElementoId__Id:" + strconv.Itoa(id)
}

func getPayloadEspacioFisico(id int) string {
	return "query=Id:" + strconv.Itoa(id)
}

func getPayloadEspacioFisicoCampo(id int) string {
	return "query=Activo,EspacioFisicoId__Id:" + strconv.Itoa(id)
}
