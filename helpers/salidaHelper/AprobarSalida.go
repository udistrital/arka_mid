package salidaHelper

import (
	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// AprobarSalida Aprobacion de una salida
func AprobarSalida(salidaId int, res *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	funcion := "AprobarSalida - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	var (
		entrada        models.FormatoBaseEntrada
		salida         models.FormatoSalida
		trSalida       models.TrSalida
		elementosActa  []*models.Elemento
		transaccion    models.TransaccionMovimientos
		tipoMovimiento int
	)

	if tr_, err := movimientosArka.GetTrSalida(salidaId); err != nil {
		return err
	} else if tr_.Salida.EstadoMovimientoId.Nombre == "Salida En Trámite" {
		trSalida = *tr_
		res.Movimiento = *trSalida.Salida
	} else {
		return
	}

	if err := utilsHelper.Unmarshal(trSalida.Salida.Detalle, &salida); err != nil {
		return err
	}

	if len(trSalida.Elementos) == 0 || salida.ConsecutivoId == 0 || salida.Funcionario == 0 {
		res.Error = "No se pudo continuar calcular la transacción contable. Contacte soporte."
		return
	}

	if err := utilsHelper.Unmarshal(trSalida.Salida.MovimientoPadreId.Detalle, &entrada); err != nil {
		return err
	}

	if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&tipoMovimiento, "SAL"); err != nil {
		return err
	}

	if err := movimientosArka.GetEstadoMovimientoIdByNombre(&trSalida.Salida.EstadoMovimientoId.Id, "Salida Aprobada"); err != nil {
		return err
	}

	var idsElementos []int
	for _, el := range trSalida.Elementos {
		idsElementos = append(idsElementos, el.ElementoActaId)
	}

	query := "Id__in:" + utilsHelper.ArrayToString(idsElementos, "|")
	if el_, err := actaRecibido.GetAllElemento(query, "ValorUnitario,ValorTotal,SubgrupoCatalogoId,TipoBienId", "SubgrupoCatalogoId", "desc", "", "-1"); err != nil {
		return err
	} else {
		elementosActa = el_
	}

	dsc := "Entrada: " + entrada.Consecutivo

	bufferCuentas := make(map[string]models.CuentaContable)
	if msg, err := asientoContable.CalcularMovimientosContables(elementosActa, dsc, res.Movimiento.MovimientoPadreId.FormatoTipoMovimientoId.Id, tipoMovimiento, salida.Funcionario, salida.Funcionario, bufferCuentas,
		nil, &transaccion.Movimientos); err != nil || msg != "" {
		res.Error = msg
		return err
	}

	if msg, err := asientoContable.CreateTransaccionContable(getTipoComprobanteSalidas(), "Salida de Almacén", &transaccion); err != nil || msg != "" {
		res.Error = msg
		return err
	}

	transaccion.ConsecutivoId = salida.ConsecutivoId
	if _, err := movimientosContables.PostTrContable(&transaccion); err != nil {
		return err
	}

	if detalleContable, err := asientoContable.GetDetalleContable(transaccion.Movimientos, bufferCuentas); err != nil {
		return err
	} else {
		res.TransaccionContable.Movimientos = detalleContable
		res.TransaccionContable.Concepto = transaccion.Descripcion
		res.TransaccionContable.Fecha = transaccion.FechaTransaccion
	}

	if movimiento_, err := movimientosArka.PutMovimiento(trSalida.Salida, trSalida.Salida.Id); err != nil {
		return err
	} else {
		trSalida.Salida = movimiento_
	}

	return
}
