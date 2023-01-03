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
		infoCorte   []models.DepreciacionElemento
		formtatoCrr int
		elementos_  []*models.Elemento
		payload     string
		terceroUD   int
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
		if val.DeltaValor == 0 {
			continue
		}

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
			if subgrupos[elemento[0].SubgrupoCatalogoId].Depreciacion || subgrupos[elemento[0].SubgrupoCatalogoId].Amortizacion {
				*elementos = append(*elementos, val.ElementoMovimientoId)
				elementos_ = append(elementos_, elemento...)
			}

		} else {
			resultado.Error = "No se pudo consultar el detalle de los elementos. Contacte soporte."
			return
		}

	}

	if len(elementos_) == 0 {
		return
	}

	if err := movimientosArka.GetFormatoTipoMovimientoIdByCodigoAbreviacion(&formtatoCrr, "CRR"); err != nil {
		return err
	}

	if msg, err := asientoContable.CalcularMovimientosContables(elementos_, getDescripcionMovmientoCierre(), 0, formtatoCrr,
		terceroUD, terceroUD, cuentas, subgrupos, &transaccion.Movimientos); err != nil || msg != "" {
		resultado.Error = msg
		return err
	}

	return
}

func getTipoComprobanteCierre() string {
	return "H22"
}

func getDescripcionMovmientoCierre() string {
	return "Mediciones posteriores"
}
