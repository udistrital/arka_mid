package depreciacionHelper

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/catalogoElementosHelper"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// calcularCierre Calcula la transacción contable que se generará una vez se liquide el cierre a una fecha determinada
func calcularCierre(fechaCorte string, elementos *[]int, transaccion *models.TransaccionMovimientos, resulado *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("calcularCierre - Unhandled Error!", "500")

	var (
		infoCorte        []models.DepreciacionElemento
		subgrupoElemento map[int]int
		detalleSubgrupos map[int]models.DetalleSubgrupo
		totalesDp        map[int]float64
		totalesAm        map[int]float64
		cuentasSubgrupos map[int]models.CuentaSubgrupo
		detalleCuentas   map[string]models.CuentaContable
		idsCuentas       []string
		formtatoDp       int
		formtatoAm       int
		query            string
		terceroUD        int
	)

	if err := movimientosArka.GetCorteDepreciacion(fechaCorte, &infoCorte); err != nil {
		return err
	}

	if len(infoCorte) == 0 {
		return
	}

	// Consulta el subgrupo al que pertenece cada elemento
	ids := []int{}
	if elementos != nil {
		for _, val := range infoCorte {
			ids = append(ids, val.ElementoMovimientoId)
		}
		*elementos = ids
		ids = []int{}
	}

	for _, val := range infoCorte {
		ids = append(ids, val.ElementoActaId)
	}

	query = "Id__in:" + utilsHelper.ArrayToString(ids, "|")
	if elemento_, err := actaRecibido.GetAllElemento(query, "Id,SubgrupoCatalogoId", "Id", "desc", "", strconv.Itoa(len(ids))); err != nil {
		return err
	} else {
		ids = []int{}
		subgrupoElemento = make(map[int]int)
		for _, el := range elemento_ {
			subgrupoElemento[el.Id] = el.SubgrupoCatalogoId
			ids = append(ids, el.SubgrupoCatalogoId)
		}
	}

	if len(ids) == 0 {
		return
	}

	ids = utilsHelper.RemoveDuplicateInt(ids)
	query = "fields=SubgrupoId,Depreciacion,Amortizacion&sortby=FechaCreacion&order=desc&limit=-1"
	query += "&query=Activo:true,SubgrupoId__Id__in:" + utilsHelper.ArrayToString(ids, "|")
	if detalles, err := catalogoElementos.GetAllDetalleSubgrupo(query); err != nil {
		return err
	} else {
		detalleSubgrupos = make(map[int]models.DetalleSubgrupo)
		for _, dt := range detalles {
			if _, ok := detalleSubgrupos[dt.SubgrupoId.Id]; !ok {
				detalleSubgrupos[dt.SubgrupoId.Id] = *dt
			}
		}
	}

	if len(detalleSubgrupos) == 0 {
		return
	}

	// Validar que existan los detalles

	totalesAm = make(map[int]float64)
	totalesDp = make(map[int]float64)
	for _, dt := range infoCorte {
		if val, ok := detalleSubgrupos[subgrupoElemento[dt.ElementoActaId]]; ok {
			if val.Depreciacion {
				totalesDp[subgrupoElemento[dt.ElementoActaId]] += dt.DeltaValor
			} else if val.Amortizacion {
				totalesAm[subgrupoElemento[dt.ElementoActaId]] += dt.DeltaValor
			}
		}
	}

	cuentasSubgrupos = make(map[int]models.CuentaSubgrupo)
	if len(totalesAm) > 0 {
		if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&formtatoAm, "AMT"); err != nil {
			return err
		}

		ids = []int{}
		for key := range totalesAm {
			ids = append(ids, key)
		}
		if err := catalogoElementosHelper.GetCuentasByMovimientoAndSubgrupos(formtatoAm, ids, cuentasSubgrupos); err != nil {
			return err
		}
	}

	if len(totalesDp) > 0 {
		if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&formtatoDp, "DEP"); err != nil {
			return err
		}

		ids = []int{}
		for key := range totalesDp {
			ids = append(ids, key)
		}
		if err := catalogoElementosHelper.GetCuentasByMovimientoAndSubgrupos(formtatoDp, ids, cuentasSubgrupos); err != nil {
			return err
		}
	}

	for _, cta := range cuentasSubgrupos {
		idsCuentas = append(idsCuentas, cta.CuentaCreditoId)
		idsCuentas = append(idsCuentas, cta.CuentaDebitoId)
	}

	idsCuentas = utilsHelper.RemoveDuplicateStr(idsCuentas)
	detalleCuentas = make(map[string]models.CuentaContable)
	if err := cuentasContables.GetDetalleCuentasContables(idsCuentas, detalleCuentas); err != nil {
		return err
	}

	query = "query=TipoDocumentoId__Nombre:NIT,Numero:" + terceros.GetDocUD()
	if terceroUD_, err := terceros.GetAllDatosIdentificacion(query); err != nil {
		return err
	} else {
		terceroUD = terceroUD_[0].TerceroId.Id
	}

	if transaccion == nil {
		transaccion = new(models.TransaccionMovimientos)
	}
	if err_, err := asientoContable.ConstruirMovimientosContables(totalesAm, detalleCuentas, cuentasSubgrupos,
		terceroUD, terceroUD, getDescripcionMovmientoCierre(), false, &transaccion.Movimientos); err != nil {
		return err
	} else if err_ != "" {
		resulado.Error = err_
		return
	}

	if err_, err := asientoContable.ConstruirMovimientosContables(totalesDp, detalleCuentas, cuentasSubgrupos,
		terceroUD, terceroUD, getDescripcionMovmientoCierre(), false, &transaccion.Movimientos); err != nil {
		return err
	} else if err_ != "" {
		resulado.Error = err_
		return
	}

	if len(transaccion.Movimientos) == 0 {
		return
	}

	if detalleContable, err := asientoContable.GetDetalleContable(transaccion.Movimientos, detalleCuentas); err != nil {
		return err
	} else if len(detalleContable) > 0 {
		trContable := models.InfoTransaccionContable{
			Movimientos: detalleContable,
			Concepto:    descAsiento(),
		}
		resulado.TransaccionContable = trContable
	}

	return
}

func getTipoComprobanteCierre() string {
	return "H22"
}

func getDescripcionMovmientoCierre() string {
	return "Mediciones posteriores"
}
