package trasladoshelper

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	crudMovimientosArka "github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	crudTerceros "github.com/udistrital/arka_mid/helpers/crud/terceros"
	midTerceros "github.com/udistrital/arka_mid/helpers/mid/terceros"
	"github.com/udistrital/arka_mid/helpers/ubicacionHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// GetDetalle Consulta los funcionarios, ubicación y elementos asociados a un traslado
func GetDetalleTraslado(id int) (Traslado *models.TrTraslado, outputError map[string]interface{}) {

	funcion := "GetDetalleTraslado"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		movimiento *models.Movimiento
		detalle    models.DetalleTraslado
	)
	Traslado = new(models.TrTraslado)

	// Se consulta el movimiento
	if movimientoA, err := crudMovimientosArka.GetMovimientoById(id); err != nil {
		return nil, err
	} else {
		movimiento = movimientoA
	}

	if err := json.Unmarshal([]byte(movimiento.Detalle), &detalle); err != nil {
		logs.Error(err)
		eval := " - json.Unmarshal([]byte(movimiento.Detalle), &detalle)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	// Se consulta el detalle del funcionario origen
	if origen, err := midTerceros.GetDetalleFuncionario(detalle.FuncionarioOrigen); err != nil {
		return nil, err
	} else {
		Traslado.FuncionarioOrigen = origen
	}

	// Se consulta el detalle del funcionario destino
	if destino, err := midTerceros.GetDetalleFuncionario(detalle.FuncionarioDestino); err != nil {
		return nil, err
	} else {
		Traslado.FuncionarioDestino = destino
	}

	// Se consulta la sede, dependencia correspondiente a la ubicacion
	if ubicacionDetalle, err := ubicacionHelper.GetSedeDependenciaUbicacion(detalle.Ubicacion); err != nil {
		return nil, err
	} else {
		Traslado.Ubicacion = ubicacionDetalle
	}

	// Se consultan los detalles de los elementos del traslado
	if elementos, err := GetElementosTraslado(detalle.Elementos); err != nil {
		return nil, err
	} else {
		Traslado.Elementos = elementos
	}
	Traslado.Detalle = movimiento.Detalle
	Traslado.Observaciones = movimiento.Observacion

	return Traslado, nil

}

func GetElementosTraslado(ids []int) (Elementos []*models.DetalleElementoPlaca, outputError map[string]interface{}) {

	funcion := "GetElementosTraslado"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		query     string
		elementos []*models.ElementosMovimiento
	)

	query = "limit=-1&fields=Id,ElementoActaId&sortby=ElementoActaId&order=desc"
	query += "&query=Id__in:" + url.QueryEscape(utilsHelper.ArrayToString(ids, "|"))
	if elementos_, err := crudMovimientosArka.GetAllElementosMovimiento(query); err != nil {
		return nil, err
	} else {
		elementos = elementos_
	}

	idsActa := []int{}
	for _, val := range elementos {
		idsActa = append(idsActa, int(val.ElementoActaId))
	}

	query = "Id__in:" + utilsHelper.ArrayToString(idsActa, "|")
	if response, err := actaRecibido.GetAllElemento(query, "", "Id", "desc", "", "-1"); err != nil {
		return nil, err
	} else {
		if len(response) == len(elementos) {
			for i := 0; i < len(response); i++ {
				elemento := new(models.DetalleElementoPlaca)

				elemento.Id = elementos[i].Id
				elemento.Nombre = response[i].Nombre
				elemento.Placa = response[i].Placa
				elemento.Marca = response[i].Marca
				elemento.Serie = response[i].Serie
				elemento.Valor = response[i].ValorTotal

				Elementos = append(Elementos, elemento)
			}
		}
	}

	return Elementos, nil
}

// RegistrarEntrada Crea registro de traslado en estado en trámite
func RegistrarTraslado(data *models.Movimiento) (result *models.Movimiento, outputError map[string]interface{}) {

	funcion := "RegistrarTraslado"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	result = new(models.Movimiento)

	detalleJSON := map[string]interface{}{}
	if err := json.Unmarshal([]byte(data.Detalle), &detalleJSON); err != nil {
		panic(err.Error())
	}

	ctxConsecutivo, _ := beego.AppConfig.Int("contxtTrasladoCons")
	if consecutivo, _, err := consecutivos.Get("%05.0f", ctxConsecutivo, "Registro Traslado Arka"); err != nil {
		return nil, err
	} else {
		consecutivo = consecutivos.Format(getTipoComprobanteTraslados()+"-", consecutivo, fmt.Sprintf("%s%04d", "-", time.Now().Year()))
		detalleJSON["Consecutivo"] = consecutivo
	}

	if jsonData, err := json.Marshal(detalleJSON); err != nil {
		logs.Error(err)
		eval := " - json.Marshal(detalleJSON)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	} else {
		data.Detalle = string(jsonData[:])
	}

	// Crea registro en api movimientos_arka_crud
	if res, err := crudMovimientosArka.PostMovimiento(data); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func GetElementosFuncionario(id int) (Elementos []*models.DetalleElementoPlaca, outputError map[string]interface{}) {

	funcion := "GetElementosFuncionario"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		elementosF    []int
		elementosM    []*models.ElementosMovimiento
		elementosActa []*models.Elemento
	)

	Elementos = make([]*models.DetalleElementoPlaca, 0)

	// Consulta lista de elementos asignados al funcionario
	if elemento_, err := crudMovimientosArka.GetElementosFuncionario(id); err != nil {
		return nil, err
	} else {
		elementosF = elemento_
	}

	// Consulta id del elemento en el api acta_recibido_crud
	if len(elementosF) > 0 {
		query := "limit=-1&sortby=ElementoActaId&order=desc&query=Id__in:"
		query += url.QueryEscape(utilsHelper.ArrayToString(elementosF, "|"))
		if elementoMovimiento_, err := crudMovimientosArka.GetAllElementosMovimiento(query); err != nil {
			return nil, err
		} else {
			elementosM = elementoMovimiento_
		}
	} else {
		return Elementos, nil
	}

	ids := []int{}
	elementosM = removeDuplicateInt(elementosM)
	for _, el := range elementosM {
		ids = append(ids, el.ElementoActaId)
	}

	// Consulta de Nombre, Placa, Marca, Serie se hace al api acta_recibido_crud
	query := "Id__in:" + utilsHelper.ArrayToString(ids, "|")
	if elemento_, err := actaRecibido.GetAllElemento(query, "", "Id", "desc", "", "-1"); err != nil {
		return nil, err
	} else {
		elementosActa = elemento_
	}

	if len(elementosActa) == len(elementosM) {
		for i := 0; i < len(elementosActa); i++ {
			if elementosActa[i].Placa != "" {
				elemento := new(models.DetalleElementoPlaca)

				elemento.Id = elementosM[i].Id
				elemento.Nombre = elementosActa[i].Nombre
				elemento.Placa = elementosActa[i].Placa
				elemento.Marca = elementosActa[i].Marca
				elemento.Serie = elementosActa[i].Serie
				elemento.Valor = elementosActa[i].ValorTotal

				Elementos = append(Elementos, elemento)
			}
		}
	}

	return Elementos, nil
}

// GetAllTraslados Consulta información general de todos los traslados filtrando por las que están pendientes por revisar
func GetAllTraslados(tramiteOnly bool) (listBajas []*models.DetalleTrasladoLista, outputError map[string]interface{}) {

	funcion := "GetAllSolicitudes"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "limit=-1&query=Activo:true,EstadoMovimientoId__Nombre"

	if tramiteOnly {
		urlcrud += url.QueryEscape(":Traslado En Trámite")
	} else {
		urlcrud += "__startswith:Traslado"
	}

	if Solicitudes, err := crudMovimientosArka.GetAllMovimiento(urlcrud); err != nil {
		return nil, err
	} else {
		if len(Solicitudes) == 0 {
			return nil, nil
		}

		tercerosBuffer := make(map[int]interface{})

		for _, solicitud := range Solicitudes {

			var (
				detalle    *models.FormatoTraslado
				Tercero_   string
				Revisor_   string
				Ubicacion_ string
			)

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

			requestUbicacion := func(id int) func() (interface{}, map[string]interface{}) {
				return func() (interface{}, map[string]interface{}) {
					if Ubicacion, err := ubicacionHelper.GetSedeDependenciaUbicacion(id); err == nil {
						return Ubicacion, nil
					}
					return nil, nil
				}
			}

			if v, err := utilsHelper.BufferGeneric(detalle.FuncionarioDestino, tercerosBuffer, requestTercero(detalle.FuncionarioDestino), nil, nil); err == nil {
				if v2, ok := v.(*models.Tercero); ok {
					Tercero_ = v2.NombreCompleto
				}
			}

			if v, err := utilsHelper.BufferGeneric(detalle.FuncionarioOrigen, tercerosBuffer, requestTercero(detalle.FuncionarioOrigen), nil, nil); err == nil {
				if v2, ok := v.(*models.Tercero); ok {
					Revisor_ = v2.NombreCompleto
				}
			}

			if v, err := utilsHelper.BufferGeneric(detalle.Ubicacion, tercerosBuffer, requestUbicacion(detalle.Ubicacion), nil, nil); err == nil {
				if v2, ok := v.(*models.DetalleSedeDependencia); ok {
					Ubicacion_ = v2.Ubicacion.EspacioFisicoId.Nombre
				}
			}

			baja := models.DetalleTrasladoLista{
				Id:                 solicitud.Id,
				Consecutivo:        detalle.Consecutivo,
				FechaCreacion:      solicitud.FechaCreacion.String(),
				FuncionarioOrigen:  Tercero_,
				FuncionarioDestino: Revisor_,
				Ubicacion:          Ubicacion_,
				EstadoMovimientoId: solicitud.EstadoMovimientoId.Id,
			}
			listBajas = append(listBajas, &baja)
		}
		return listBajas, nil

	}
}

func getTipoComprobanteTraslados() string {
	return "T"
}

func removeDuplicateInt(intSlice []*models.ElementosMovimiento) []*models.ElementosMovimiento {
	allKeys := make(map[int]bool)
	list := make([]*models.ElementosMovimiento, 0)
	for _, item := range intSlice {
		if _, value := allKeys[item.ElementoActaId]; !value {
			allKeys[item.ElementoActaId] = true
			list = append(list, item)
		}
	}
	return list
}
