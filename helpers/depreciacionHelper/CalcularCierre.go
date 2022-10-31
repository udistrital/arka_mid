package depreciacionHelper

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/catalogoElementos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/terceros"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// calcularCierre Calcula la transacción contable que se generará una vez se liquide el cierre a una fecha determinada
func calcularCierre(fechaCorte string, elementos *[]int, cuentas map[string]models.CuentaContable, transaccion *models.TransaccionMovimientos, resultado *models.ResultadoMovimiento) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("calcularCierre - Unhandled Error!", "500")

	var (
		infoCorte    []models.DepreciacionElemento
		formtatoDp   int
		formtatoAm   int
		payload      string
		terceroUD    int
		elementosAmt []*models.Elemento
		elementosDep []*models.Elemento
	)

	if err := movimientosArka.GetCorteDepreciacion(fechaCorte, &infoCorte); err != nil {
		return err
	}

	if len(infoCorte) == 0 {
		return
	}

	if elementos == nil {
		elementos = new([]int)
	}

	payload = "query=TipoDocumentoId__Nombre:NIT,Numero:" + terceros.GetDocUD()
	if terceroUD_, err := terceros.GetAllDatosIdentificacion(payload); err != nil {
		return err
	} else if len(terceroUD_) != 1 || terceroUD_[0].TerceroId == nil {
		resultado.Error = "No se pudo consultar el tercero para asociar al movimiento contable. Contacte soporte."
		return
	} else {
		terceroUD = terceroUD_[0].TerceroId.Id
	}

	subgrupos := make(map[int]models.DetalleSubgrupo)
	for _, val := range infoCorte {
		payload = "Id:" + strconv.Itoa(val.ElementoActaId)
		if elemento, err := actaRecibido.GetAllElemento(payload, "Id,ValorUnitario,ValorTotal,SubgrupoCatalogoId,TipoBienId", "", "", "", ""); err != nil {
			return err
		} else if len(elemento) == 1 {

			payload = "limit=1&fields=TipoBienId,Amortizacion,Depreciacion&sortby=Id&order=desc&query=Activo:true,SubgrupoId__Id:"
			if _, ok := subgrupos[elemento[0].SubgrupoCatalogoId]; !ok {
				if sg, err := catalogoElementos.GetAllDetalleSubgrupo(payload + strconv.Itoa(elemento[0].SubgrupoCatalogoId)); err != nil {
					return err
				} else if len(sg) == 1 {
					subgrupos[elemento[0].SubgrupoCatalogoId] = *sg[0]
				} else {
					resultado.Error = "No se pudo consultar la parametrización de las clases. Contacte soporte."
					return
				}
			}

			elemento[0].ValorTotal = val.DeltaValor
			if subgrupos[elemento[0].SubgrupoCatalogoId].Depreciacion {
				*elementos = append(*elementos, val.ElementoMovimientoId)
				elementosDep = append(elementosDep, elemento...)
			} else if subgrupos[elemento[0].SubgrupoCatalogoId].Amortizacion {
				*elementos = append(*elementos, val.ElementoMovimientoId)
				elementosAmt = append(elementosAmt, elemento...)
			}

		} else {
			resultado.Error = "No se pudo consultar el detalle de los elementos. Contacte soporte."
			return
		}

	}

	if len(elementosDep) == 0 && len(elementosAmt) == 0 {
		return
	}

	if len(elementosAmt) > 0 {
		if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&formtatoAm, "AMT"); err != nil {
			return err
		}

		if msg, err := asientoContable.CalcularMovimientosContables(elementosAmt, getDescripcionMovmientoCierre(), formtatoAm,
			terceroUD, terceroUD, cuentas, subgrupos, &transaccion.Movimientos); err != nil || msg != "" {
			resultado.Error = msg
			return err
		}
	}

	if len(elementosDep) > 0 {
		if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&formtatoDp, "DEP"); err != nil {
			return err
		}

		if msg, err := asientoContable.CalcularMovimientosContables(elementosDep, getDescripcionMovmientoCierre(), formtatoDp,
			terceroUD, terceroUD, cuentas, subgrupos, &transaccion.Movimientos); err != nil || msg != "" {
			resultado.Error = msg
			return err
		}
	}

	return
}

func getTipoComprobanteCierre() string {
	return "H22"
}

func getDescripcionMovmientoCierre() string {
	return "Mediciones posteriores"
}
