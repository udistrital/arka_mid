package bodegaConsumoHelper

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	// "github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/utils_oas/request"
)

func GetElementosSinAsignar() (Elementos []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetElementosSinAsignar - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	// fmt.Println("aaaaaaaaaaaaaaaaaaaaa")
	url := "http://" + beego.AppConfig.String("movimientosArkaService") + "elementos_movimiento?limit=-1"
	url += "&query=Activo:true,MovimientoId__FormatoTipoMovimientoId__Nombre:Salida%20de%20Consumo,MovimientoId__EstadoMovimientoId__Nombre:Salida%20Aprobada"
	// logs.Debug("url:", url)
	if res, err := request.GetJsonTest(url, &Elementos); err == nil && res.StatusCode == 200 {

		if keys := len(Elementos[0]); keys != 0 {

			elementosActaBuffer := make(map[int]interface{})
			subgruposCatalogoBuffer := make(map[int]interface{})

			for i, elemento := range Elementos {
				void := true

				var elementoActaId int
				if v, err := strconv.Atoi(fmt.Sprintf("%v", elemento["ElementoActaId"])); err == nil && v > 0 {
					elementoActaId = v
				} else {
					err = fmt.Errorf("ElementoActaId='%v', erroneo para 'elementos_movimiento.Id=%v'", elemento["ElementoActaId"], elemento["Id"])
					logs.Warn(err)
					// TODO: revisar si esto es suficiente
					continue
				}

				urlElemento := "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento"
				urlElemento += "?query=Id:" + strconv.Itoa(elementoActaId)
				urlElemento += "&fields=Id,Nombre,Marca,Serie,SubgrupoCatalogoId"
				if detalle, err := utilsHelper.BufferGet(elementoActaId, elementosActaBuffer, urlElemento); err == nil && detalle != nil {

					var subgrupoCatalogoId int
					if v, err := strconv.Atoi(fmt.Sprintf("%v", detalle["SubgrupoCatalogoId"])); err == nil && v > 0 {
						subgrupoCatalogoId = v
					} else {
						err = fmt.Errorf("SubgrupoCatalogoId='%v', erroneo para 'elemento(Acta).Id=%d'", detalle["SubgrupoCatalogoId"], elementoActaId)
						logs.Warn(err)
						// TODO: revisar si esto es suficiente
						continue
					}

					urlSubgrupo := "http://" + beego.AppConfig.String("catalogoElementosService") + "subgrupo"
					urlSubgrupo += "?query=Id:" + strconv.Itoa(subgrupoCatalogoId)
					if subgrupo, err := utilsHelper.BufferGet(subgrupoCatalogoId, subgruposCatalogoBuffer, urlSubgrupo); err == nil && subgrupo != nil {
						Elementos[i]["Nombre"] = detalle["Nombre"]
						Elementos[i]["Marca"] = detalle["Marca"]
						Elementos[i]["Serie"] = detalle["Serie"]
						Elementos[i]["SubgrupoCatalogoId"] = subgrupo

						void = false
					} else {
						if err == nil {
							logs.Warn(fmt.Errorf("no hay subgrupo_catalogo.Id=%d (CRUD catalogo) asociado al elemento.Id=%d (CRUD Actas)", subgrupoCatalogoId, elementoActaId))
						} else {
							logs.Warn(err)
						}
					}

				} else {
					if err == nil {
						logs.Warn(fmt.Errorf("no hay elemento.Id=%d (CRUD Actas) asociado al elemento.Id=%v (CRUD movimientos_arka)", elementoActaId, elemento["Id"]))
					} else {
						logs.Warn(err)
					}
				}

				if void {
					Elementos[i] = nil
				}
			}

			// Quitar elementos inconsistentes
			fin := len(Elementos)
			// logs.Debug("fin(antes):", fin)
			for i := 0; i < fin; {
				if Elementos[i] != nil {
					i++
				} else {
					Elementos[i] = Elementos[fin-1]
					fin--
				}
			}
			// logs.Debug("fin(despues):", fin)

		}
		return Elementos, nil

	} else {
		if err == nil {
			err = fmt.Errorf("undesired Status Code: %d", res.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetElementosSinAsignar - request.GetJsonTest(url, &Elementos)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}
