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
		detalle          models.FormatoTraslado
		elementosActa    []*models.Elemento
		novedades        map[int]models.NovedadElemento
		ids              []int
		query            string
		tipoMovimientoId int
		transaccion      models.TransaccionMovimientos
	)

	if movimiento_, err := movimientosArka.GetMovimientoById(id); err != nil {
		return
	} else {
		response.Movimiento = *movimiento_
	}

	if err := utilsHelper.Unmarshal(response.Movimiento.Detalle, &detalle); err != nil {
		return err
	}

	query = "limit=-1&fields=Id,ElementoActaId,ValorTotal&sortby=ElementoActaId&order=desc"
	query += "&query=Id__in:" + url.QueryEscape(utilsHelper.ArrayToString(detalle.Elementos, "|"))
	if elementos_, err := movimientosArka.GetAllElementosMovimiento(query); err != nil {
		return err
	} else if len(elementos_) == len(detalle.Elementos) {
		for _, el := range elementos_ {
			ids = append(ids, el.ElementoActaId)
		}
	} else {
		response.Error = "No se pudo consultar la parametrización de los elementos. Contacte soporte"
		return
	}

	fields := "Id,SubgrupoCatalogoId,TipoBienId,ValorUnitario,ValorTotal"
	query = "Id__in:" + utilsHelper.ArrayToString(ids, "|")
	if elementos_, err := actaRecibido.GetAllElemento(query, fields, "Id", "desc", "", "-1"); err != nil {
		return err
	} else if len(elementos_) == len(detalle.Elementos) {
		elementosActa = elementos_
	} else {
		response.Error = "No se pudo consultar la parametrización de los elementos. Contacte soporte"
		return
	}

	query = "limit=-1&fields=ElementoMovimientoId,ValorLibros&sortby=MovimientoId,FechaCreacion&order=asc,asc&query=Activo:true,ElementoMovimientoId__Id__in:"
	query += utilsHelper.ArrayToString(detalle.Elementos, "|")
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

	if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&tipoMovimientoId, "SAL"); err != nil {
		return err
	}

	if err := movimientosArka.GetEstadoMovimientoIdByNombre(&response.Movimiento.EstadoMovimientoId.Id, "Traslado Aprobado"); err != nil {
		return err
	}

	bufferCuentas := make(map[string]models.CuentaContable)
	if msg, err := asientoContable.CalcularMovimientosContables(elementosActa, descMovDestino(), tipoMovimientoId, detalle.FuncionarioDestino, detalle.FuncionarioOrigen, bufferCuentas,
		nil, &transaccion.Movimientos); err != nil || msg != "" {
		response.Error = msg
		return err
	}

	transaccion.ConsecutivoId = detalle.ConsecutivoId
	if msg, err := asientoContable.CreateTransaccionContable(getTipoComprobanteTraslados(), "Traslado de elementos", &transaccion); err != nil || msg != "" {
		response.Error = msg
		return err
	}

	if _, err := movimientosContables.PostTrContable(&transaccion); err != nil {
		return err
	}

	if detalleContable, err := asientoContable.GetDetalleContable(transaccion.Movimientos, bufferCuentas); err != nil {
		return err
	} else {
		response.TransaccionContable.Movimientos = detalleContable
		response.TransaccionContable.Concepto = transaccion.Descripcion
		response.TransaccionContable.Fecha = transaccion.FechaTransaccion
	}

	if movimiento_, err := movimientosArka.PutMovimiento(&response.Movimiento, response.Movimiento.Id); err != nil {
		return err
	} else {
		response.Movimiento = *movimiento_
	}

	return
}

func descMovDestino() string {
	return "Traslado de elementos"
}
