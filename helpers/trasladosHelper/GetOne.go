package trasladoshelper

import (
	"net/url"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/oikos"
	"github.com/udistrital/arka_mid/helpers/mid/terceros"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// GetOne Consulta los funcionarios, ubicación y elementos asociados a un traslado
func GetOne(id int) (Traslado *models.TrTraslado, outputError map[string]interface{}) {

	defer errorCtrl.ErrorControlFunction("GetOne - Unhandled Error!", "500")

	var detalle models.FormatoTraslado
	Traslado = new(models.TrTraslado)

	// Se consulta el movimiento
	movimientoA, outputError := movimientosArka.GetMovimientoById(id)
	if outputError != nil {
		return
	}

	Traslado.Movimiento = movimientoA
	outputError = utilsHelper.Unmarshal(Traslado.Movimiento.Detalle, &detalle)
	if outputError != nil {
		return
	}

	// Se consulta el detalle del funcionario origen
	Traslado.FuncionarioOrigen, outputError = terceros.GetDetalleFuncionario(detalle.FuncionarioOrigen)
	if outputError != nil {
		return
	}

	// Se consulta el detalle del funcionario destino
	Traslado.FuncionarioDestino, outputError = terceros.GetDetalleFuncionario(detalle.FuncionarioDestino)
	if outputError != nil {
		return
	}

	// Se consulta la sede, dependencia correspondiente a la ubicacion
	Traslado.Ubicacion, outputError = oikos.GetSedeDependenciaUbicacion(detalle.Ubicacion)
	if outputError != nil {
		return
	}

	// Se consultan los detalles de los elementos del traslado
	Traslado.Elementos, outputError = getElementosTraslado(detalle.Elementos)
	if outputError != nil {
		return
	}

	if Traslado.Movimiento.EstadoMovimientoId.Nombre == "Traslado Aprobado" && Traslado.Movimiento.ConsecutivoId != nil && *Traslado.Movimiento.ConsecutivoId > 0 {
		Traslado.TrContable = &models.InfoTransaccionContable{}
		*Traslado.TrContable, outputError = asientoContable.GetFullDetalleContable(*Traslado.Movimiento.ConsecutivoId)
		if outputError != nil {
			return
		}

	}

	Traslado.Observaciones = Traslado.Movimiento.Observacion
	return
}

func getElementosTraslado(ids []int) (Elementos []*models.DetalleElementoPlaca, outputError map[string]interface{}) {

	funcion := "getElementosTraslado"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	query := "limit=-1&fields=Id,ElementoActaId&sortby=ElementoActaId&order=desc"
	query += "&query=Id__in:" + url.QueryEscape(utilsHelper.ArrayToString(ids, "|"))
	elementos, outputError := movimientosArka.GetAllElementosMovimiento(query)
	if outputError != nil {
		return
	}

	idsActa := []int{}
	for _, val := range elementos {
		idsActa = append(idsActa, *val.ElementoActaId)
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

func getTipoComprobanteTraslados() string {
	return "N39"
}
