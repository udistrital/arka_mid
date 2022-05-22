package trasladoshelper

import (
	"net/url"
	"strconv"

	"github.com/astaxie/beego"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/oikos"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/mid/terceros"
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
		detalle    models.FormatoTraslado
	)
	Traslado = new(models.TrTraslado)

	// Se consulta el movimiento
	query := "query=Id:" + strconv.Itoa(id)
	if movimientoA, err := movimientosArka.GetAllMovimiento(query); err != nil {
		return nil, err
	} else if len(movimientoA) == 1 {
		movimiento = movimientoA[0]
	} else {
		return nil, nil
	}

	if err := utilsHelper.Unmarshal(movimiento.Detalle, &detalle); err != nil {
		return nil, err
	}

	// Se consulta el detalle del funcionario origen
	if origen, err := terceros.GetDetalleFuncionario(detalle.FuncionarioOrigen); err != nil {
		return nil, err
	} else {
		Traslado.FuncionarioOrigen = origen
	}

	// Se consulta el detalle del funcionario destino
	if destino, err := terceros.GetDetalleFuncionario(detalle.FuncionarioDestino); err != nil {
		return nil, err
	} else {
		Traslado.FuncionarioDestino = destino
	}

	// Se consulta la sede, dependencia correspondiente a la ubicacion
	if ubicacionDetalle, err := oikos.GetSedeDependenciaUbicacion(detalle.Ubicacion); err != nil {
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

	if movimiento.EstadoMovimientoId.Nombre == "Traslado Aprobado" {
		if detalle.ConsecutivoId > 0 {
			if tr, err := movimientosContables.GetTransaccion(detalle.ConsecutivoId, "consecutivo", true); err != nil {
				return nil, err
			} else if len(tr.Movimientos) > 0 {
				if detalleContable, err := asientoContable.GetDetalleContable(tr.Movimientos, nil); err != nil {
					return nil, err
				} else {
					trContable := models.InfoTransaccionContable{
						Movimientos: detalleContable,
						Concepto:    tr.Descripcion,
						Fecha:       tr.FechaTransaccion,
					}
					Traslado.TrContable = &trContable
				}
			}
		}
	}

	Traslado.Movimiento = movimiento
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
	if elementos_, err := movimientosArka.GetAllElementosMovimiento(query); err != nil {
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

	var (
		detalle     models.FormatoTraslado
		consecutivo models.Consecutivo
	)

	if err := utilsHelper.Unmarshal(data.Detalle, &detalle); err != nil {
		return nil, err
	}

	ctxConsecutivo, _ := beego.AppConfig.Int("contxtAjusteCons")
	if err := consecutivos.Get(ctxConsecutivo, "Registro Traslado Arka", &consecutivo); err != nil {
		return nil, err
	}

	detalle.Consecutivo = consecutivos.Format("%05d", getTipoComprobanteTraslados(), &consecutivo)
	detalle.ConsecutivoId = consecutivo.Id

	if err := utilsHelper.Marshal(detalle, &data.Detalle); err != nil {
		return nil, err
	}

	// Crea registro en api movimientos_arka_crud
	if err := movimientosArka.PostMovimiento(data); err != nil {
		return nil, err
	}

	return data, nil

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
	if elemento_, err := movimientosArka.GetElementosFuncionario(id); err != nil {
		return nil, err
	} else {
		elementosF = elemento_
	}

	// Consulta id del elemento en el api acta_recibido_crud
	if len(elementosF) > 0 {
		query := "limit=-1&sortby=ElementoActaId&order=desc&query=Id__in:"
		query += url.QueryEscape(utilsHelper.ArrayToString(elementosF, "|"))
		if elementoMovimiento_, err := movimientosArka.GetAllElementosMovimiento(query); err != nil {
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

func getTipoComprobanteTraslados() string {
	return "N39"
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
