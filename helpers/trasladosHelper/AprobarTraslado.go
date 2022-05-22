package trasladoshelper

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/catalogoElementosHelper"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// AprobarTraslado Actualiza el estado del traslado y genera la transaccion contable correspondiente
func AprobarTraslado(id int) (resultado map[string]interface{}, outputError map[string]interface{}) {

	funcion := "AprobarTraslado"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		movimiento          models.Movimiento
		detalle             models.FormatoTraslado
		ids                 []int
		query               string
		elementosMovimiento map[int]models.ElementosMovimiento
		elementosActa       map[int]models.Elemento
		novedades           map[int]models.NovedadElemento
		cuentasSubgrupo     map[int]models.CuentaSubgrupo
		detalleCuentas      map[string]models.CuentaContable
		idsCuentas          []string
		parDebito           int
		parCredito          int
		transaccion         models.TransaccionMovimientos
		tipoMovimientoId    int
	)

	resultado = make(map[string]interface{})

	if movimiento_, err := movimientosArka.GetMovimientoById(id); err != nil {
		return
	} else {
		movimiento = *movimiento_
	}

	if err := json.Unmarshal([]byte(movimiento.Detalle), &detalle); err != nil {
		logs.Error(err)
		eval := " - json.Unmarshal([]byte(movimiento.Detalle), &detalle)"
		return nil, errorctrl.Error(funcion+eval, err, "500")
	}

	ids = append(ids, detalle.Elementos...)

	query = "limit=-1&fields=Id,ElementoActaId,ValorTotal&sortby=ElementoActaId&order=desc"
	query += "&query=Id__in:" + url.QueryEscape(utilsHelper.ArrayToString(ids, "|"))
	if elementos_, err := movimientosArka.GetAllElementosMovimiento(query); err != nil {
		return nil, err
	} else {
		elementosMovimiento = make(map[int]models.ElementosMovimiento)
		for _, el := range elementos_ {
			elementosMovimiento[el.Id] = *el
		}

	}

	query = "limit=-1&fields=ElementoMovimientoId,ValorLibros&sortby=MovimientoId,FechaCreacion&order=asc,asc&query=Activo:true,ElementoMovimientoId__Id__in:"
	query += utilsHelper.ArrayToString(ids, "|")
	if novedades_, err := movimientosArka.GetAllNovedadElemento(query); err != nil {
		return nil, err
	} else {
		novedades = make(map[int]models.NovedadElemento)
		for _, nov := range novedades_ {
			novedades[nov.ElementoMovimientoId.Id] = *nov
		}

	}

	ids = []int{}
	for _, el := range elementosMovimiento {
		ids = append(ids, el.ElementoActaId)
	}

	fields := "Id,SubgrupoCatalogoId"
	query = "Id__in:" + utilsHelper.ArrayToString(ids, "|")
	if elementos_, err := actaRecibido.GetAllElemento(query, fields, "", "", "", "-1"); err != nil {
		return nil, err
	} else {
		elementosActa = make(map[int]models.Elemento)
		for _, el_ := range elementos_ {
			elementosActa[el_.Id] = *el_
		}

	}

	var totales = make(map[int]float64)
	for _, el_ := range detalle.Elementos {
		var sg int
		if val, ok := elementosActa[elementosMovimiento[el_].ElementoActaId]; ok {
			sg = val.SubgrupoCatalogoId
		}

		if val, ok := novedades[el_]; ok {
			totales[sg] += val.ValorLibros
			continue
		}

		if val, ok := elementosMovimiento[el_]; ok {
			totales[sg] += val.ValorTotal
		}

	}

	if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&tipoMovimientoId, "SAL"); err != nil {
		return nil, err
	}

	if err := movimientosArka.GetEstadoMovimientoIdByNombre(&movimiento.EstadoMovimientoId.Id, "Traslado Aprobado"); err != nil {
		return nil, err
	}

	ids = []int{}
	for sg := range totales {
		ids = append(ids, sg)
	}

	cuentasSubgrupo = make(map[int]models.CuentaSubgrupo)
	if err := catalogoElementosHelper.GetCuentasByMovimientoAndSubgrupos(tipoMovimientoId, ids, cuentasSubgrupo); err != nil {
		return nil, err
	}

	for _, cta := range cuentasSubgrupo {
		idsCuentas = append(idsCuentas, cta.CuentaCreditoId)
		idsCuentas = append(idsCuentas, cta.CuentaDebitoId)
	}

	detalleCuentas = make(map[string]models.CuentaContable)
	if err := cuentasContables.GetDetalleCuentasContables(idsCuentas, detalleCuentas); err != nil {
		return nil, err
	}

	if db_, cr_, err := parametros.GetParametrosDebitoCredito(); err != nil {
		return nil, err
	} else {
		parDebito = db_
		parCredito = cr_
	}

	transaccion = *new(models.TransaccionMovimientos)
	asientoContable.GenerarMovimientosContables(totales, detalleCuentas, cuentasSubgrupo, parDebito, parCredito, detalle.FuncionarioDestino, detalle.FuncionarioOrigen, descMovDestino(), false, &transaccion.Movimientos)

	if len(transaccion.Movimientos) > 0 {
		transaccion.Activo = true
		transaccion.ConsecutivoId = detalle.ConsecutivoId
		transaccion.Descripcion = "Traslado de elementos"
		transaccion.Etiquetas = ""
		transaccion.FechaTransaccion = time.Now()

		if _, err := movimientosContables.PostTrContable(&transaccion); err != nil {
			return nil, err
		}
	}

	if detalleContable, err := asientoContable.GetDetalleContable(transaccion.Movimientos, detalleCuentas); err != nil {
		return nil, err
	} else if len(transaccion.Movimientos) > 0 {
		trContable := models.InfoTransaccionContable{
			Movimientos: detalleContable,
			Concepto:    transaccion.Descripcion,
			Fecha:       transaccion.FechaTransaccion,
		}
		resultado["trContable"] = trContable
	}

	if movimiento_, err := movimientosArka.PutMovimiento(&movimiento, movimiento.Id); err != nil {
		return nil, err
	} else {
		resultado["movimiento"] = movimiento_
	}

	return
}

func descMovDestino() string {
	return "Traslado de elementos"
}

func descMovOrigen() string {
	return "Movimiento tercero origen de traslado"
}
