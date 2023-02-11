package trasladoshelper

import (
	"net/url"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// AprobarTraslado Actualiza el estado del traslado y genera la transaccion contable correspondiente
func AprobarTraslado(id int, response *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	funcion := "AprobarTraslado - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	var (
		detalle      models.FormatoTraslado
		tipoSalida   int
		query        string
		mapElementos = make(map[int][]models.ElementosMovimiento)
		transaccion  models.TransaccionMovimientos
	)

	movimiento_, outputError := movimientosArka.GetMovimientoById(id)
	if outputError != nil {
		return
	}

	response.Movimiento = *movimiento_
	if err := utilsHelper.Unmarshal(response.Movimiento.Detalle, &detalle); err != nil {
		return err
	}

	if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&tipoSalida, "SAL"); err != nil {
		return err
	}

	if err := movimientosArka.GetEstadoMovimientoIdByNombre(&response.Movimiento.EstadoMovimientoId.Id, "Traslado Aprobado"); err != nil {
		return err
	}

	query = "limit=-1&fields=Id,ElementoActaId,ValorTotal,MovimientoId&sortby=ElementoActaId&order=desc" +
		"&query=Id__in:" + url.QueryEscape(utilsHelper.ArrayToString(detalle.Elementos, "|"))
	if elementos_, err := movimientosArka.GetAllElementosMovimiento(query); err != nil {
		return err
	} else if len(elementos_) == len(detalle.Elementos) {
		for _, el := range elementos_ {
			mapElementos[el.MovimientoId.MovimientoPadreId.FormatoTipoMovimientoId.Id] = append(mapElementos[el.MovimientoId.MovimientoPadreId.FormatoTipoMovimientoId.Id], *el)
		}
	} else {
		response.Error = "No se pudo consultar la parametrización de los elementos. Contacte soporte"
		return
	}

	bufferCuentas := make(map[string]models.CuentaContable)
	fields := "Id,SubgrupoCatalogoId,TipoBienId,ValorUnitario,ValorTotal"
	for tipoEntr, el_ := range mapElementos {

		var (
			ids           []int
			ids_          []int
			elementosActa []*models.Elemento
			novedades     map[int]models.NovedadElemento
		)

		for _, el := range el_ {
			ids = append(ids, el.Id)
			ids_ = append(ids_, el.ElementoActaId)
		}

		query = "Id__in:" + utilsHelper.ArrayToString(ids_, "|")
		elementosActa, outputError = actaRecibido.GetAllElemento(query, fields, "Id", "desc", "", "-1")
		if outputError != nil {
			return
		} else if len(elementosActa) != len(el_) {
			response.Error = "No se pudo consultar la parametrización de los elementos. Contacte soporte"
			return
		}

		query = "limit=-1&fields=ElementoMovimientoId,ValorLibros&sortby=MovimientoId,FechaCreacion&order=asc,asc&query=Activo:true,ElementoMovimientoId__Id__in:" +
			utilsHelper.ArrayToString(ids, "|")
		if novedades_, err := movimientosArka.GetAllNovedadElemento(query); err != nil {
			return err
		} else {
			novedades = make(map[int]models.NovedadElemento)
			for _, nov := range novedades_ {
				novedades[nov.ElementoMovimientoId.ElementoActaId] = *nov
			}
		}

		for _, el := range elementosActa {
			if val, ok := novedades[el.Id]; ok {
				el.ValorTotal = val.ValorLibros
			}
		}

		response.Error, outputError = asientoContable.CalcularMovimientosContables(elementosActa, descMovDestino(), tipoEntr, tipoSalida, detalle.FuncionarioDestino, detalle.FuncionarioOrigen, bufferCuentas, nil, &transaccion.Movimientos)
		if outputError != nil || response.Error != "" {
			return
		}
	}

	transaccion.ConsecutivoId = *response.Movimiento.ConsecutivoId
	response.Error, outputError = asientoContable.CreateTransaccionContable(getTipoComprobanteTraslados(), "Traslado de elementos", &transaccion)
	if outputError != nil || response.Error != "" {
		return
	}

	response.TransaccionContable.Concepto = transaccion.Descripcion
	response.TransaccionContable.Fecha = transaccion.FechaTransaccion
	response.TransaccionContable.Movimientos, outputError = asientoContable.GetDetalleContable(transaccion.Movimientos, bufferCuentas)
	if outputError != nil {
		return
	}

	_, outputError = movimientosContables.PostTrContable(&transaccion)
	if outputError != nil {
		return
	}

	_, outputError = movimientosArka.PutMovimiento(&response.Movimiento, response.Movimiento.Id)

	return
}

func descMovDestino() string {
	return "Traslado de elementos"
}
