package bajasHelper

import (
	"net/url"
	"strconv"

	"github.com/astaxie/beego"

	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/oikos"
	midTerceros "github.com/udistrital/arka_mid/helpers/mid/terceros"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// RegistrarBaja Crea registro de baja
func RegistrarBaja(baja *models.TrSoporteMovimiento) (bajaR *models.Movimiento, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("RegistrarBaja - Unhandled Error!", "500")

	var consecutivo models.Consecutivo
	ctxConsecutivo, _ := beego.AppConfig.Int("contxtBajaCons")
	outputError = consecutivos.Get(ctxConsecutivo, "Registro Baja Arka", &consecutivo)
	if outputError != nil {
		return
	}

	baja.Movimiento.Consecutivo = utilsHelper.String(consecutivos.Format("%05d", getTipoComprobanteBajas(), &consecutivo))
	baja.Movimiento.ConsecutivoId = &consecutivo.Id

	// Crea registro en api movimientos_arka_crud
	outputError = movimientosArka.PostMovimiento(baja.Movimiento)
	if outputError != nil {
		return
	}

	// Crea registro en table soporte_movimiento si es necesario
	baja.Soporte.MovimientoId = baja.Movimiento
	outputError = movimientosArka.PostSoporteMovimiento(baja.Soporte)
	if outputError != nil {
		return
	}

	return baja.Movimiento, nil
}

// ActualizarBaja Actualiza información de baja
func ActualizarBaja(baja *models.TrSoporteMovimiento, bajaId int) (bajaR *models.Movimiento, outputError map[string]interface{}) {

	funcion := "ActualizarBaja"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		movimiento *models.Movimiento
		soporte    *models.SoporteMovimiento
	)

	// Actualiza registro en api movimientos_arka_crud
	if movimiento_, err := movimientosArka.PutMovimiento(baja.Movimiento, bajaId); err != nil {
		return nil, err
	} else {
		movimiento = movimiento_
	}

	// Actualiza el documento soporte en la tabla soporte_movimiento
	query := "query=MovimientoId__Id:" + strconv.Itoa(bajaId)
	if soporte_, err := movimientosArka.GetAllSoporteMovimiento(query); err != nil {
		return nil, err
	} else {
		soporte = soporte_[0]
		soporte.DocumentoId = baja.Soporte.DocumentoId
	}

	if _, err := movimientosArka.PutSoporteMovimiento(soporte, soporte.Id); err != nil {
		return nil, err
	}

	return movimiento, nil

}

// GetDetalleElemento Consulta historial de un elemento dado el id del elemento en el api acta_recibido_crud
func GetDetalleElemento(id int, Elemento *models.DetalleElementoBaja) (outputError map[string]interface{}) {

	funcion := "GetDetalleElemento - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	var (
		elemento           models.DetalleElemento
		elementoMovimiento models.ElementosMovimiento
	)

	if err := movimientosArka.GetElementosMovimientoById(id, &elementoMovimiento); err != nil {
		return err
	} else if elementoMovimiento.Id == 0 {
		return
	}

	if historial_, err := movimientosArka.GetHistorialElemento(elementoMovimiento.Id, true); err != nil {
		return err
	} else {
		Elemento.Historial = historial_
	}

	// Consulta de Marca, Nombre, Serie y Subgrupo se hace mediante el actaRecibidoHelper
	ids := []int{elementoMovimiento.ElementoActaId}
	if elementos, err := actaRecibido.GetElementos(0, ids); err != nil {
		return err
	} else if len(elementos) == 1 {
		elemento = *elementos[0]
	} else {
		return
	}

	if fc, ub, err := GetEncargado(Elemento.Historial); err != nil {
		return err
	} else {
		if ubicacion_, err := oikos.GetSedeDependenciaUbicacion(ub); err != nil {
			return err
		} else {
			Elemento.Ubicacion = ubicacion_
		}

		if funcionario_, err := midTerceros.GetInfoTerceroById(fc); err != nil {
			return err
		} else {
			Elemento.Funcionario = funcionario_
		}
	}

	Elemento.Id = elementoMovimiento.Id
	Elemento.Placa = elemento.Placa
	Elemento.Nombre = elemento.Nombre
	Elemento.Marca = elemento.Marca
	Elemento.Serie = elemento.Serie
	Elemento.SubgrupoCatalogoId = elemento.SubgrupoCatalogoId

	return
}

// GetDetalleElementos consulta el historial de una serie de elementos dados los ids en el api movimientos_arka_crud
func GetDetalleElementos(ids []int) (Elementos []*models.DetalleElementoBaja, outputError map[string]interface{}) {

	funcion := "GetDetalleElementos"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		elementosActa       []*models.DetalleElemento
		elementosMovimiento []*models.ElementosMovimiento
	)
	Elementos = make([]*models.DetalleElementoBaja, 0)

	// Consulta asignación de los elementos
	query := "sortby=ElementoActaId&order=desc&limit=-1&query=Id__in:"
	query += url.QueryEscape(utilsHelper.ArrayToString(ids, "|"))
	if elementoMovimiento_, err := movimientosArka.GetAllElementosMovimiento(query); err != nil {
		return nil, err
	} else {
		elementosMovimiento = elementoMovimiento_
	}

	ids = []int{}
	for _, el := range elementosMovimiento {
		ids = append(ids, el.ElementoActaId)
	}

	// Consulta de Marca, Nombre, Serie y Subgrupo se hace mediante el actaRecibidoHelper
	if elemento_, err := actaRecibido.GetElementos(0, ids); err != nil {
		return nil, err
	} else {
		elementosActa = elemento_
	}

	if len(elementosActa) == len(elementosMovimiento) {

		for i := 0; i < len(elementosActa); i++ {

			elemento := new(models.DetalleElementoBaja)

			if historial_, err := movimientosArka.GetHistorialElemento(elementosMovimiento[i].Id, true); err != nil {
				return nil, err
			} else {
				elemento.Historial = historial_
			}

			if fc, ub, err := GetEncargado(elemento.Historial); err != nil {
				return nil, err
			} else {
				if ubicacion_, err := oikos.GetSedeDependenciaUbicacion(ub); err != nil {
					return nil, err
				} else {
					elemento.Ubicacion = ubicacion_
				}

				if funcionario_, err := midTerceros.GetInfoTerceroById(fc); err != nil {
					return nil, err
				} else {
					elemento.Funcionario = funcionario_
				}
			}

			elemento.Id = elementosMovimiento[i].Id
			elemento.Placa = elementosActa[i].Placa
			elemento.Nombre = elementosActa[i].Nombre
			elemento.Marca = elementosActa[i].Marca
			elemento.Serie = elementosActa[i].Serie
			elemento.SubgrupoCatalogoId = elementosActa[i].SubgrupoCatalogoId

			Elementos = append(Elementos, elemento)
		}
	}

	return Elementos, nil
}

// findIdInArray Retorna la posicion en que se encuentra el id específicado
func FindInArray(cuentasSg []*models.CuentasSubgrupo, subgrupoId int) (i int) {
	for i, cuentaSg := range cuentasSg {
		if int(cuentaSg.SubgrupoId.Id) == subgrupoId {
			return i
		}
	}
	return -1
}

// GetEncargado Retorna el funcionario y ubicacion actual de un elemento de acuerdo a su historial
func GetEncargado(historial *models.Historial) (funcionarioId int, ubicacionId int, outputError map[string]interface{}) {

	funcion := "GetEncargado - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	if historial.Traslados != nil {
		var detalleTr models.DetalleTraslado
		if err := utilsHelper.Unmarshal(historial.Traslados[0].Detalle, &detalleTr); err != nil {
			eval := "utilsHelper.Unmarshal(historial.Traslados[0].Detalle, &detalleTr)"
			return 0, 0, errorctrl.Error(funcion+eval, err, "500")
		}

		return detalleTr.FuncionarioDestino, detalleTr.Ubicacion, nil
	} else {
		var detalleS models.FormatoSalida
		if err := utilsHelper.Unmarshal(historial.Salida.Detalle, &detalleS); err != nil {
			eval := "utilsHelper.Unmarshal(historial.Salida.Detalle, &detalleS)"
			return 0, 0, errorctrl.Error(funcion+eval, err, "500")
		}

		return detalleS.Funcionario, detalleS.Ubicacion, nil
	}
}
