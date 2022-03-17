package bajasHelper

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/actaRecibido"
	crud_actas "github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	crudMovimientosArka "github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	crudTerceros "github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/cuentasContablesHelper"
	midTerceros "github.com/udistrital/arka_mid/helpers/mid/terceros"
	"github.com/udistrital/arka_mid/helpers/salidaHelper"
	"github.com/udistrital/arka_mid/helpers/ubicacionHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

type InfoCuentasSubgrupos struct {
	CuentaDebito  *models.CuentaContable
	CuentaCredito *models.CuentaContable
}

// RegistrarBaja Crea registro de baja
func RegistrarBaja(baja *models.TrSoporteMovimiento) (bajaR *models.Movimiento, outputError map[string]interface{}) {

	funcion := "RegistrarBaja"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		movimiento *models.Movimiento
		detalle    *models.FormatoBaja
	)

	if err := json.Unmarshal([]byte(baja.Movimiento.Detalle), &detalle); err != nil {
		eval := " - json.Unmarshal([]byte(baja.Movimiento.Detalle), &detalle)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	ctxConsecutivo, _ := beego.AppConfig.Int("contxtBajaCons")
	if consecutivo, _, err := utilsHelper.GetConsecutivo("%05.0f", ctxConsecutivo, "Registro Baja Arka"); err != nil {
		return nil, err
	} else {
		consecutivo = utilsHelper.FormatConsecutivo(getTipoComprobanteBajas()+"-", consecutivo, fmt.Sprintf("%s%04d", "-", time.Now().Year()))
		detalle.Consecutivo = consecutivo
	}

	if jsonData, err := json.Marshal(detalle); err != nil {
		eval := " - json.Marshal(detalle)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	} else {
		baja.Movimiento.Detalle = string(jsonData[:])
	}

	// Crea registro en api movimientos_arka_crud
	if movimiento_, err := crudMovimientosArka.PostMovimiento(baja.Movimiento); err != nil {
		return nil, err
	} else {
		movimiento = movimiento_
	}

	// Crea registro en table soporte_movimiento si es necesario
	baja.Soporte.MovimientoId = movimiento
	if _, err := crudMovimientosArka.PostSoporteMovimiento(baja.Soporte); err != nil {
		return nil, err
	}

	return movimiento, nil
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
	if movimiento_, err := crudMovimientosArka.PutMovimiento(baja.Movimiento, bajaId); err != nil {
		return nil, err
	} else {
		movimiento = movimiento_
	}

	// Actualiza el documento soporte en la tabla soporte_movimiento
	query := "query=MovimientoId__Id:" + strconv.Itoa(bajaId)
	if soporte_, err := crudMovimientosArka.GetAllSoporteMovimiento(query); err != nil {
		return nil, err
	} else {
		soporte = soporte_[0]
		soporte.DocumentoId = baja.Soporte.DocumentoId
	}

	if _, err := crudMovimientosArka.PutSoporteMovimiento(soporte, soporte.Id); err != nil {
		return nil, err
	}

	return movimiento, nil
}

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
			if ctaCr_, err := cuentasContablesHelper.GetCuentaContable(cuentasSubgrupo[idx].CuentaCreditoId); err != nil {
				return nil, err
			} else {
				if err := formatdata.FillStruct(ctaCr_, &ctaCr); err != nil {
					logs.Error(err)
					eval := " - formatdata.FillStruct(ctaCr_, &ctaCr)"
					return nil, errorctrl.Error(funcion+eval, err, "500")
				}
			}

			if ctaDb_, err := cuentasContablesHelper.GetCuentaContable(cuentasSubgrupo[idx].CuentaDebitoId); err != nil {
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

func TraerDatosElemento(id int) (Elemento map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/TraerDatosElemento", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	var elemento_movimiento_ []map[string]interface{}
	// var movimiento_ map[string]interface{}

	url := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento/?query=ElementoActaId:" + fmt.Sprintf("%v", id) + ",Activo:true"
	if _, err := request.GetJsonTest(url, &elemento_movimiento_); err == nil {

		if v, err := salidaHelper.TraerDetalle(elemento_movimiento_[0]["MovimientoId"]); err == nil {

			fmt.Println("Elemento Movimiento: ", elemento_movimiento_)

			var movimiento_ map[string]interface{}
			if jsonString3, err := json.Marshal(elemento_movimiento_[0]["MovimientoId"]); err == nil {
				if err2 := json.Unmarshal(jsonString3, &movimiento_); err2 == nil {
					movimiento_["MovimientoPadreId"] = nil
				}
			}

			var elemento_ []map[string]interface{}

			urlcrud_ := "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento?query=Id:" + fmt.Sprintf("%v", elemento_movimiento_[0]["ElementoActaId"]) + "&fields=Id,Nombre,TipoBienId,Marca,Serie,Placa,SubgrupoCatalogoId"
			if _, err := request.GetJsonTest(urlcrud_, &elemento_); err == nil {

				fmt.Println("Elemento: ", elemento_)

				var subgrupo_ map[string]interface{}
				urlcrud_2 := "http://" + beego.AppConfig.String("catalogoElementosService") + "subgrupo/" + fmt.Sprintf("%v", elemento_[0]["SubgrupoCatalogoId"])
				if _, err := request.GetJsonTest(urlcrud_2, &subgrupo_); err == nil {
					Elemento := map[string]interface{}{
						"Id":                 elemento_[0]["Id"],
						"Placa":              elemento_[0]["Placa"],
						"Nombre":             elemento_[0]["Nombre"],
						"TipoBienId":         elemento_[0]["TipoBienId"],
						"Entrada":            v["MovimientoPadreId"],
						"Salida":             movimiento_,
						"SubgrupoCatalogoId": subgrupo_,
						"Marca":              elemento_[0]["Marca"],
						"Serie":              elemento_[0]["Serie"],
						"Funcionario":        v["Funcionario"],
						"Sede":               v["Sede"],
						"Dependencia":        v["Dependencia"],
						"Ubicacion":          v["Ubicacion"],
					}

					// elemento_[0]["SubgrupoCatalogoId"] = subgrupo_
					// elemento_movimiento_[0]["ElementoActaId"] = elemento_[0]
					// Elemento = elemento_movimiento_[0]
					return Elemento, nil

				} else {
					logs.Error(err)
					outputError = map[string]interface{}{
						"funcion": "/TraerDatosElemento - request.GetJsonTest(urlcrud_2, &subgrupo_)",
						"err":     err,
						"status":  "502",
					}
					return nil, outputError
				}
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "/TraerDatosElemento - request.GetJsonTest(urlcrud_, &elemento_)",
					"err":     err,
					"status":  "502",
				}
				return nil, outputError
			}
		} else {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "/TraerDatosElemento - salidaHelper.TraerDetalle(elemento_movimiento_[0][\"MovimientoId\"])",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/TraerDatosElemento - request.GetJsonTest(url, &elemento_movimiento_)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}

// GetAllSolicitudes Consulta información general de todas las bajas filtrando por las que están pendientes por revisar en almacén y en comité
func GetAllSolicitudes(revComite bool, revAlmacen bool) (listBajas []*models.DetalleBaja, outputError map[string]interface{}) {

	funcion := "GetAllSolicitudes"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "limit=-1&sortby=Id&order=desc&query=Activo:true,EstadoMovimientoId__Nombre"

	if revComite {
		urlcrud += url.QueryEscape(":Baja En Comité")
	} else if revAlmacen {
		urlcrud += url.QueryEscape(":Baja En Trámite")
	} else {
		urlcrud += "__startswith:Baja"
	}

	if Solicitudes, err := crudMovimientosArka.GetAllMovimiento(urlcrud); err == nil {

		if len(Solicitudes) == 0 {
			return nil, nil
		}

		tercerosBuffer := make(map[int]interface{})

		for _, solicitud := range Solicitudes {

			var detalle *models.FormatoBaja
			var Tercero_ string
			var Revisor_ string

			if err := json.Unmarshal([]byte(solicitud.Detalle), &detalle); err != nil {
				eval := " - json.Unmarshal([]byte(solicitud.Detalle), &detalle)"
				return nil, errorctrl.Error(funcion+eval, err, "500")
			}

			requestTercero := func(id int) func() (interface{}, map[string]interface{}) {
				return func() (interface{}, map[string]interface{}) {
					if Tercero, err := crudTerceros.GetTerceroById(id); err == nil {
						return Tercero, nil
					}
					return nil, nil
				}
			}

			funcionarioID := detalle.Funcionario
			if v, err := utilsHelper.BufferGeneric(funcionarioID, tercerosBuffer, requestTercero(funcionarioID), nil, nil); err == nil {
				if v2, ok := v.(*models.Tercero); ok {
					Tercero_ = v2.NombreCompleto
				}
			}

			revisorID := detalle.Revisor
			if v, err := utilsHelper.BufferGeneric(revisorID, tercerosBuffer, requestTercero(revisorID), nil, nil); err == nil {
				if v2, ok := v.(*models.Tercero); ok {
					Revisor_ = v2.NombreCompleto
				}
			}

			baja := models.DetalleBaja{
				Id:                 solicitud.Id,
				Consecutivo:        detalle.Consecutivo,
				FechaCreacion:      solicitud.FechaCreacion.String(),
				FechaRevisionA:     detalle.FechaRevisionA,
				FechaRevisionC:     detalle.FechaRevisionC,
				Funcionario:        Tercero_,
				Revisor:            Revisor_,
				TipoBaja:           solicitud.FormatoTipoMovimientoId.Nombre,
				EstadoMovimientoId: solicitud.EstadoMovimientoId.Id,
			}
			listBajas = append(listBajas, &baja)
		}
		return listBajas, nil

	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/GetAllSolicitudes",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}

// TraerDetalle Consulta el detalle de la baja, elementos, revisor, solicitante, soporte, tipo
func TraerDetalle(id int) (Baja *models.TrBaja, outputError map[string]interface{}) {

	funcion := "TraerDetalle"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		movimiento *models.Movimiento
		detalle    models.FormatoBaja
	)
	Baja = new(models.TrBaja)

	// Se consulta el movimiento
	if movimientoA, err := crudMovimientosArka.GetMovimientoById(id); err != nil {
		return nil, err
	} else {
		movimiento = movimientoA
	}

	if err := json.Unmarshal([]byte(movimiento.Detalle), &detalle); err != nil {
		eval := " - json.Unmarshal([]byte(movimiento.Detalle), &detalle)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	// Se consulta el detalle del funcionario solicitante
	if detalle.Funcionario > 0 {
		if funcionario, err := midTerceros.GetInfoTerceroById(detalle.Funcionario); err != nil {
			return nil, err
		} else {
			Baja.Funcionario = funcionario
		}
	}

	// Se consulta el detalle del revisor si lo hay
	if detalle.Revisor > 0 {
		if revisor, err := midTerceros.GetInfoTerceroById(detalle.Revisor); err != nil {
			return nil, err
		} else {
			Baja.Revisor = revisor
		}
	}

	// Se consulta el detalle de los elementos relacionados en la solicitud
	if len(detalle.Elementos) > 0 {
		if elementos, err := GetDetalleElementos(detalle.Elementos); err != nil {
			return nil, err
		} else {
			Baja.Elementos = elementos
		}
	}

	// Se consulta el detalle de los elementos relacionados en la solicitud
	query := "query=MovimientoId__Id:" + strconv.Itoa(id)
	if soportes, err := crudMovimientosArka.GetAllSoporteMovimiento(query); err != nil {
		return nil, err
	} else if len(soportes) > 0 {
		Baja.Soporte = soportes[0].DocumentoId
	}

	Baja.Id = movimiento.Id
	Baja.TipoBaja = movimiento.FormatoTipoMovimientoId
	Baja.Consecutivo = detalle.Consecutivo
	Baja.Observaciones = movimiento.Observacion
	Baja.RazonRechazo = detalle.RazonRechazo
	Baja.Resolucion = detalle.Resolucion
	Baja.FechaRevisionC = detalle.FechaRevisionC

	return Baja, nil

}

// GetDetalleElemento Consulta historial de un elemento dado el id del elemento en el api acta_recibido_crud
func GetDetalleElemento(id int) (Elemento *models.DetalleElementoBaja, outputError map[string]interface{}) {

	funcion := "GetDetalleElemento"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		elemento           []*models.DetalleElemento
		elementoMovimiento *models.ElementosMovimiento
	)
	Elemento = new(models.DetalleElementoBaja)

	// Consulta de Marca, Nombre, Serie y Subgrupo se hace mediante el actaRecibidoHelper
	ids := []int{id}
	if elemento_, err := actaRecibido.GetElementos(0, ids); err != nil {
		return nil, err
	} else {
		elemento = elemento_
	}

	query := "sortby=Id&order=desc&query=ElementoActaId:" + strconv.Itoa(id)
	if elementoMovimiento_, err := crudMovimientosArka.GetAllElementosMovimiento(query); err != nil {
		return nil, err
	} else if len(elementoMovimiento_) > 0 {
		elementoMovimiento = elementoMovimiento_[0]
	} else {
		return Elemento, nil
	}

	if historial_, err := crudMovimientosArka.GetHistorialElemento(elementoMovimiento.Id, true); err != nil {
		return nil, err
	} else {
		Elemento.Historial = historial_
	}

	if fc, ub, err := GetEncargado(Elemento.Historial); err != nil {
		return nil, err
	} else {
		if ubicacion_, err := ubicacionHelper.GetSedeDependenciaUbicacion(ub); err != nil {
			return nil, err
		} else {
			Elemento.Ubicacion = ubicacion_
		}

		if funcionario_, err := midTerceros.GetInfoTerceroById(fc); err != nil {
			return nil, err
		} else {
			Elemento.Funcionario = funcionario_
		}
	}

	Elemento.Id = elementoMovimiento.Id
	Elemento.Placa = elemento[0].Placa
	Elemento.Nombre = elemento[0].Nombre
	Elemento.Marca = elemento[0].Marca
	Elemento.Serie = elemento[0].Serie
	Elemento.SubgrupoCatalogoId = elemento[0].SubgrupoCatalogoId

	return Elemento, nil
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
	if elementoMovimiento_, err := crudMovimientosArka.GetAllElementosMovimiento(query); err != nil {
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

			if historial_, err := crudMovimientosArka.GetHistorialElemento(elementosMovimiento[i].Id, true); err != nil {
				return nil, err
			} else {
				elemento.Historial = historial_
			}

			if fc, ub, err := GetEncargado(elemento.Historial); err != nil {
				return nil, err
			} else {
				if ubicacion_, err := ubicacionHelper.GetSedeDependenciaUbicacion(ub); err != nil {
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

func getTipoComprobanteBajas() string {
	return "B"
}

// findIdInArray Retorna la posicion en que se encuentra el id específicado
func FindInArray(cuentasSg []*models.CuentaSubgrupo, subgrupoId int) (i int) {
	for i, cuentaSg := range cuentasSg {
		if int(cuentaSg.SubgrupoId.Id) == subgrupoId {
			return i
		}
	}
	return -1
}

// Retorna el actual encargado de un elemento de acuerdo a su historial
func GetEncargado(historial *models.Historial) (funcionarioId int, ubicacionId int, outputError map[string]interface{}) {

	funcion := "GetEncargado"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	if historial.Traslados != nil {
		var detalleTr models.DetalleTraslado
		if err := json.Unmarshal([]byte(historial.Traslados[0].Detalle), &detalleTr); err != nil {
			eval := " - json.Unmarshal([]byte(historial.Traslados[0].Detalle), &detalleTr)"
			return 0, 0, errorctrl.Error(funcion+eval, err, "500")
		}
		funcionarioId = detalleTr.FuncionarioDestino
		ubicacionId = detalleTr.Ubicacion
		return funcionarioId, ubicacionId, nil
	} else {
		detalleS := map[string]interface{}{}
		if err := json.Unmarshal([]byte(historial.Salida.Detalle), &detalleS); err != nil {
			eval := " - json.Unmarshal([]byte(historial.Salida.Detalle), &detalleS)"
			return 0, 0, errorctrl.Error(funcion+eval, err, "500")
		}
		funcionarioId = int(detalleS["funcionario"].(float64))
		ubicacionId = int(detalleS["ubicacion"].(float64))
		return funcionarioId, ubicacionId, nil
	}
}
