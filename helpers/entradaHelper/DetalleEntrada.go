package entradaHelper

import (
	"strconv"

	administrativa_ "github.com/udistrital/administrativa_mid_api/models"
	"github.com/udistrital/arka_mid/helpers/asientoContable"
	"github.com/udistrital/arka_mid/helpers/crud/actaRecibido"
	"github.com/udistrital/arka_mid/helpers/crud/administrativa"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	tercerosCRUD "github.com/udistrital/arka_mid/helpers/crud/terceros"
	administrativaAMAZON "github.com/udistrital/arka_mid/helpers/mid/administrativa"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

// DetalleEntrada Consulta el detalle de una entrada incluyendo la transaccion contable (si aplica)
func DetalleEntrada(entradaId int) (result map[string]interface{}, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("DetalleEntrada - Unhandled Error!", "500")

	var (
		detalle         models.FormatoBaseEntrada
		movimiento      models.Movimiento
		unidadEjecutora models.Parametro
		query           string
	)

	resultado := make(map[string]interface{})

	query = "query=Id:" + strconv.Itoa(entradaId)
	if mov, err := movimientosArka.GetAllMovimiento(query); err != nil {
		return nil, err
	} else if len(mov) > 0 {
		movimiento = *mov[0]
	}

	if err := utilsHelper.Unmarshal(movimiento.Detalle, &detalle); err != nil {
		return nil, err
	}

	if detalle.ActaRecibidoId > 0 {
		query = "ActaRecibidoId__Id:" + strconv.Itoa(detalle.ActaRecibidoId)
		var acta models.HistoricoActa
		if tr, err := actaRecibido.GetAllHistoricoActa(query, "", "Id", "desc", "", "1"); err != nil {
			return nil, err
		} else {
			acta = tr[0]
		}

		if acta.ProveedorId > 0 {
			if tercero, err := tercerosCRUD.GetNombreTerceroById(acta.ProveedorId); err != nil {
				return nil, err
			} else {
				resultado["proveedor"] = tercero
			}
		}

		if acta.ActaRecibidoId.UnidadEjecutoraId > 0 {
			if err := parametros.GetParametroById(acta.ActaRecibidoId.UnidadEjecutoraId, &unidadEjecutora); err != nil {
				return nil, err
			}
			resultado["unidadEjecutora"] = unidadEjecutora
		}
	}

	if detalle.ContratoId > 0 && detalle.VigenciaContrato != "" {
		var contrato administrativa_.InformacionContrato
		if unidadEjecutora.CodigoAbreviacion == "UD" {
			outputError = administrativa.GetContrato(detalle.ContratoId, detalle.VigenciaContrato, &contrato)
			if outputError != nil {
				return
			}

			if contrato.Contrato.NumeroContratoSuscrito != "" {
				resultado["contrato"] = contrato.Contrato
				if contrato.Contrato.TipoContrato != "" {
					var tipoContrato administrativa_.TipoContrato
					outputError = administrativa.GetTipoContratoById(contrato.Contrato.TipoContrato, &tipoContrato)
					if outputError != nil {
						return
					}
					resultado["tipo_contrato_id"] = tipoContrato
				}
			}
		} else {
			contrato.Contrato.NumeroContratoSuscrito = strconv.Itoa(detalle.ContratoId)
			contrato.Contrato.Vigencia = detalle.VigenciaContrato
			resultado["contrato"] = contrato.Contrato
		}
	}

	if (movimiento.EstadoMovimientoId.Nombre == "Entrada Aprobada" || movimiento.EstadoMovimientoId.Nombre == "Entrada Con Salida") && movimiento.ConsecutivoId != nil && *movimiento.ConsecutivoId > 0 {
		resultado["TransaccionContable"] = models.InfoTransaccionContable{}

		resultado["TransaccionContable"], outputError = asientoContable.GetFullDetalleContable(*movimiento.ConsecutivoId)
		if outputError != nil {
			return
		}
	}

	if detalle.Factura > 0 {
		soporte := *new(models.SoporteActa)
		if err := actaRecibido.GetSoporteById(detalle.Factura, &soporte); err != nil {
			return nil, err
		}
		resultado["factura"] = soporte
	}

	if detalle.SupervisorId > 0 {
		supervisor := make(map[string]interface{})
		if err := administrativaAMAZON.GetSupervisor(detalle.SupervisorId, &supervisor); err != nil {
			return nil, err
		} else if len(supervisor) > 0 {
			resultado["supervisor"] = supervisor
		}

		if val, ok := supervisor["DependenciaSupervisor"]; ok && val != nil && val.(string) != "" {
			var dependencia []interface{}
			if err := administrativaAMAZON.GetAllDependenciaSIC("query=ESFCODIGODEP:"+val.(string), &dependencia); err != nil {
				return nil, err
			}

			if len(dependencia) > 0 {
				supervisor["DependenciaSupervisor"] = dependencia[0]
				resultado["supervisor"] = supervisor
			}
		}

	}

	if detalle.OrdenadorGastoId > 0 {
		ordenadores := make(map[string]interface{})
		if err := administrativaAMAZON.GetOrdenadores(detalle.OrdenadorGastoId, &ordenadores); err != nil {
			return nil, err
		} else if len(ordenadores) > 0 {
			resultado["ordenador"] = ordenadores
		}
	}

	if len(detalle.Elementos) > 0 {
		var detalleElementos = make([]map[string]interface{}, 0)
		for _, el := range detalle.Elementos {
			query = "limit=1&query=Id:" + strconv.Itoa(el.Id)
			detalleMov, err := movimientosArka.GetAllElementosMovimiento(query)
			if err != nil {
				return nil, err
			} else if len(detalleMov) != 1 {
				continue
			}

			var detalleElemento map[string]interface{}
			outputError = utilsHelper.FillStruct(el, &detalleElemento)
			if outputError != nil {
				return
			}

			var elemento_ models.Elemento
			outputError = actaRecibido.GetElementoById(detalleMov[0].ElementoActaId, &elemento_)
			if outputError != nil {
				return
			}

			detalleElemento["Salida"] = detalleMov[0].MovimientoId
			detalleElemento["Placa"] = elemento_.Placa
			detalleElemento["ValorTotal"] = elemento_.ValorTotal

			if el.ValorLibros != nil && el.ValorResidual != nil && el.VidaUtil != nil {
				detalleElemento["ValorLibros"] = el.ValorLibros
				detalleElemento["ValorResidual"] = el.ValorResidual
				detalleElemento["VidaUtil"] = el.VidaUtil
			}

			if el.AprovechadoId != nil && *el.AprovechadoId > 0 {
				var elemento__ models.Elemento
				outputError = actaRecibido.GetElementoById(*el.AprovechadoId, &elemento__)
				if outputError != nil {
					return
				}
				detalleElemento["AprovechadoId"] = elemento__.Placa
			}

			detalleElementos = append(detalleElementos, detalleElemento)

		}

		resultado["elementos"] = detalleElementos
	}

	if soporte, err := movimientosArka.GetAllSoporteMovimiento("fields=DocumentoId&query=MovimientoId__Id:" + strconv.Itoa(entradaId)); err != nil {
		return nil, err
	} else if len(soporte) > 0 {
		resultado["documentoId"] = soporte[0].DocumentoId
	}

	resultado["movimiento"] = movimiento

	return resultado, nil
}
