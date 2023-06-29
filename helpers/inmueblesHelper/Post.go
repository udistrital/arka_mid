package inmuebleshelper

import (
	"strconv"
	"time"

	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

func Post(inmueble *models.Inmueble) (resultado models.ResultadoMovimiento, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("Post - Unhandled Error!", "500")

	if inmueble.Cuentas.CuentaCreditoId.Id == "" || inmueble.Cuentas.CuentaDebitoId.Id == "" {
		resultado.Error = "No se indicaron las cuentas para el inmueble."
	} else if inmueble.SubgrupoId.Id <= 0 {
		resultado.Error = "No se indicó la clase del inmueble."
	} else if inmueble.ElementoMovimiento.VidaUtil > 0 {
		if inmueble.CuentasMediciones.CuentaCreditoId.Id == "" || inmueble.CuentasMediciones.CuentaDebitoId.Id == "" {
			resultado.Error = "No se indicaron las cuentas para la depreciación del inmueble."
		} else if inmueble.ElementoMovimiento.ValorResidual > inmueble.ElementoMovimiento.ValorTotal {
			resultado.Error = "El valor residual no puede ser mayor al valor del inmueble."
		}
	}

	if resultado.Error != "" {
		return
	}

	payload := "limit=1&sortby=Id&order=desc&query=TipoActaId__CodigoAbreviacion:INM,Activo:true"
	actas, outputError := actaRecibido.GetAllActaRecibido(payload)
	if outputError != nil {
		return
	}

	var acta models.ActaRecibido
	if len(actas) == 1 {
		acta = actas[0]
	} else {
		acta = models.ActaRecibido{
			Activo:     true,
			TipoActaId: &models.TipoActa{Id: 3},
		}

		outputError = actaRecibido.PostActaRecibido(&acta)
		if outputError != nil {
			return
		}
	}

	inmueble.Elemento.ActaRecibidoId = &acta
	inmueble.Elemento.EspacioFisicoId = inmueble.EspacioFisico.Id
	inmueble.Elemento.EstadoElementoId = &models.EstadoElemento{Id: 2}
	inmueble.Elemento.SubgrupoCatalogoId = inmueble.SubgrupoId.Id
	inmueble.Elemento.Activo = true

	outputError = actaRecibido.PostElemento(&inmueble.Elemento)
	if outputError != nil {
		return
	}

	resultado.Error, outputError = registrarCuentas(*inmueble)
	if resultado.Error != "" && outputError != nil {
		return
	}

	var fechaCorte *time.Time
	payload = "limit=1&sortby=Id&order=desc&query=FormatoTipoMovimientoId__CodigoAbreviacion:INM_REG,FechaCorte"
	if inmueble.ElementoMovimiento.MovimientoId != nil && inmueble.ElementoMovimiento.MovimientoId.FechaCorte == nil {
		payload += "__isnull:true"
	} else {
		payload += ":" + inmueble.ElementoMovimiento.MovimientoId.FechaCorte.UTC().Format("2006-01-02")
		fechaCorte = inmueble.ElementoMovimiento.MovimientoId.FechaCorte
	}

	movimientos, outputError := movimientosArka.GetAllMovimiento(payload)
	if outputError != nil {
		return
	}

	var movimiento models.Movimiento
	if len(movimientos) == 1 {
		movimiento = *movimientos[0]
	} else {
		var formato, estado int
		outputError = movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&formato, "INM_REG")
		if outputError != nil {
			return
		}

		outputError = movimientosArka.GetEstadoMovimientoIdByNombre(&estado, "Bienes inmuebles registrados")
		if outputError != nil {
			return
		}

		movimiento = models.Movimiento{
			Detalle:                 "{}",
			FechaCorte:              fechaCorte,
			Activo:                  true,
			FormatoTipoMovimientoId: &models.FormatoTipoMovimiento{Id: formato},
			EstadoMovimientoId:      &models.EstadoMovimiento{Id: estado},
		}

		outputError = movimientosArka.PostMovimiento(&movimiento)
		if outputError != nil {
			return
		}
	}

	inmueble.ElementoMovimiento.MovimientoId = &movimiento
	inmueble.ElementoMovimiento.ElementoActaId = utilsHelper.Int(inmueble.Elemento.Id)
	inmueble.ElementoMovimiento.Activo = true

	outputError = movimientosArka.PostElementosMovimiento(&inmueble.ElementoMovimiento)

	return
}

func registrarCuentas(inmueble models.Inmueble) (mensaje string, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("registrarCuentas - Unhandled Error!", "500")

	mensaje, outputError = registrarCuentas_(inmueble.Elemento.Id, inmueble.Cuentas, "CC_ENT")
	if mensaje != "" || outputError != nil || inmueble.CuentasMediciones.CuentaCreditoId.Id == "" {
		return
	}

	mensaje, outputError = registrarCuentas_(inmueble.Elemento.Id, inmueble.Cuentas, "CC_MED")

	return
}

func registrarCuentas_(elementoId int, cuentas_ models.ParametrizacionContable_, tipoCuentas string) (mensaje string, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("registrarCuentas_ - Unhandled Error!", "500")

	payload := "limit=1&sortby=Id&order=desc&query=Activo:true,CampoId__Sigla:" + tipoCuentas + ",ElementoId__Id:" + strconv.Itoa(elementoId)
	cuentas, outputError := actaRecibido.GetAllElementoCampo(payload)
	if outputError != nil {
		return
	}

	var elementoCampo models.ElementoCampo
	var valor models.ParametrizacionContable
	var campoId int

	if len(cuentas) == 1 {
		var detalle models.ParametrizacionContable
		outputError = utilsHelper.Unmarshal(cuentas[0].Valor, &detalle)
		if outputError != nil || (detalle.CuentaCreditoId == cuentas_.CuentaCreditoId.Id && detalle.CuentaDebitoId == cuentas_.CuentaDebitoId.Id) {
			return
		}

		campoId = cuentas[0].CampoId.Id
		valor = models.ParametrizacionContable{
			CuentaCreditoId: cuentas_.CuentaCreditoId.Id,
			CuentaDebitoId:  cuentas_.CuentaDebitoId.Id,
		}

		cuentas[0].Activo = false
		outputError = actaRecibido.PutElementoCampo(&cuentas[0], cuentas[0].Id)
		if outputError != nil {
			return
		}

	} else {
		payload := "query=Sigla:" + tipoCuentas
		campo, outputError_ := actaRecibido.GetAllCampo(payload)
		if outputError_ != nil {
			outputError = outputError_
			return
		} else if len(campo) == 0 {
			mensaje = "No se pudieron registrar las cuentas contables. Contacte soporte."
			return
		}

		campoId = campo[0].Id
		valor = models.ParametrizacionContable{
			CuentaCreditoId: cuentas_.CuentaCreditoId.Id,
			CuentaDebitoId:  cuentas_.CuentaDebitoId.Id,
		}

	}

	elementoCampo = models.ElementoCampo{
		ElementoId: &models.Elemento{Id: elementoId},
		CampoId:    &models.Campo{Id: campoId},
		Activo:     true,
	}

	outputError = utilsHelper.Marshal(valor, &elementoCampo.Valor)
	if outputError != nil {
		return
	}

	outputError = actaRecibido.PostElementoCampo(&elementoCampo)

	return

}
