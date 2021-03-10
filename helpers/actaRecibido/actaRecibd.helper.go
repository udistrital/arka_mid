package actaRecibido

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/ubicacionHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetActasRecibidoTipo ...
func GetActasRecibidoTipo(tipoActa int) (actasRecibido []models.ActaRecibidoUbicacion, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/GetActasRecibidoTipo - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var (
		urlcrud       string
		historicoActa []*models.HistoricoActa
	)
	if tipoActa != 0 { // (1) error parametro
		urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "historico_acta?query=EstadoActaId.Id:" + strconv.Itoa(tipoActa) + ",Activo:True&limit=-1"
		// logs.Debug(urlcrud)
		if response, err := request.GetJsonTest(urlcrud, &historicoActa); err == nil && response.StatusCode == 200 { // (2) error servicio caido
			// logs.Debug(historicoActa[0].EstadoActaId)
			// logs.Debug("Estado:", tipoActa, "- Actas:", len(historicoActa), "- Id_acta[0]:", historicoActa[0].Id)

			if len(historicoActa) == 0 || historicoActa[0].Id == 0 {
				return nil, nil
			}
			// logs.Debug(historicoActa, "- len:", len(historicoActa))

			for _, acta := range historicoActa {
				var ubicacion *models.AsignacionEspacioFisicoDependencia
				if id := acta.ActaRecibidoId.UbicacionId; id > 0 {
					if ubicaciones, err := ubicacionHelper.GetAsignacionSedeDependencia(strconv.Itoa(id)); err == nil {
						if jsonString, err := json.Marshal(ubicaciones); err == nil {
							if err := json.Unmarshal(jsonString, &ubicacion); err != nil {
								logs.Error(err)
								return nil, map[string]interface{}{
									"funcion": "GetActasRecibidoTipo",
									"err":     err,
									"status":  "500",
								}
							}
						} else {
							logs.Error(err)
							return nil, map[string]interface{}{
								"funcion": "GetActasRecibidoTipo",
								"err":     err,
								"status":  "500",
							}
						}

					} else {
						logs.Error(err)
						return nil, map[string]interface{}{
							"funcion": "GetActasRecibidoTipo",
							"err":     err,
							"status":  "502",
						}
					}
				}

				actaRecibidoAux := models.ActaRecibidoUbicacion{
					Id:                acta.ActaRecibidoId.Id,
					RevisorId:         acta.ActaRecibidoId.RevisorId,
					FechaCreacion:     acta.ActaRecibidoId.FechaCreacion,
					FechaModificacion: acta.ActaRecibidoId.FechaModificacion,
					FechaVistoBueno:   acta.ActaRecibidoId.FechaVistoBueno,
					Observaciones:     acta.ActaRecibidoId.Observaciones,
					Activo:            acta.ActaRecibidoId.Activo,
					EstadoActaId:      acta.EstadoActaId,
					UbicacionId:       ubicacion,
				}

				actasRecibido = append(actasRecibido, actaRecibidoAux)
			}
			return actasRecibido, nil
		} else {
			if err == nil {
				err = fmt.Errorf("Error (3) estado de la solicitud: %s", response.Status)
			}
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "/GetActasRecibidoTipo",
				"err":     err,
				"status:": "502",
			}
			return nil, outputError
		}
	} else {
		err := fmt.Errorf("Error (1) Parametro")
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/GetActasRecibidoTipo",
			"err":     err,
			"status":  "400",
		}
		return nil, outputError
	}

}

func GetElementoById(Id string) (Elemento map[string]interface{}, outputError map[string]interface{}) {
	var urlelemento string
	var elemento []map[string]interface{}
	urlelemento = "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento/?query=Id:" + Id + "&limit=1"
	if response, err := request.GetJsonTest(urlelemento, &elemento); err == nil {

		if response.StatusCode == 200 {
			for _, element := range elemento {
				if len(element) == 0 {
					return nil, map[string]interface{}{"Function": "GetElementoById", "Error": "Sin Elementos"}
				} else {
					return element, nil
				}

			}
		}
	} else {
		return nil, map[string]interface{}{"Function": "GetElementoById", "Error": err}
	}
	return
}
