package actaRecibido

import (
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
				"funcion": "/GetActasRecibidoTipo",
				"err":     err,
				"status":  "502",
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
		logs.Debug(urlcrud)
		if response, err := request.GetJsonTest(urlcrud, &historicoActa); err == nil && response.StatusCode == 200 { // (2) error servicio caido
			logs.Debug(historicoActa[0].EstadoActaId)

			if len(historicoActa) == 0 || historicoActa[0].Id == 0 {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "/GetActasRecibidoTipo",
					"err":     err,
					"status":  "200",
				}
				return nil, outputError
			}

			if response.StatusCode == 200 { // (3) error estado de la solicitud
				for _, acta := range historicoActa {
					// UBICACION
					ubicacion, err := ubicacionHelper.GetUbicacion(acta.ActaRecibidoId.UbicacionId)

					if err != nil {
						panic(err)
					}

					logs.Debug(ubicacion)

					actaRecibidoAux := models.ActaRecibidoUbicacion{
						Id:                acta.ActaRecibidoId.Id,
						RevisorId:         acta.ActaRecibidoId.RevisorId,
						FechaCreacion:     acta.ActaRecibidoId.FechaCreacion,
						FechaModificacion: acta.ActaRecibidoId.FechaModificacion,
						FechaVistoBueno:   acta.ActaRecibidoId.FechaVistoBueno,
						Observaciones:     acta.ActaRecibidoId.Observaciones,
						Activo:            acta.ActaRecibidoId.Activo,
						EstadoActaId:      acta.EstadoActaId,
						UbicacionId:       ubicacion[0],
					}

					actasRecibido = append(actasRecibido, actaRecibidoAux)
				}
				return actasRecibido, nil
			} else {
				logs.Info("Error (3) estado de la solicitud")
				outputError = map[string]interface{}{"Function": "GetActasRecibidoTipo:GetActasRecibidoTipo", "Error": response.Status}
				return nil, outputError
			}
		} else {
			logs.Info("Error (2) servicio caido")
			outputError = map[string]interface{}{"Function": "GetActasRecibidoTipo", "Error": err}
			return nil, outputError
		}
	} else {
		logs.Info("Error (1) Parametro")
		outputError = map[string]interface{}{"Function": "FuncionalidadMidController:getUserAgora", "Error": "null parameter"}
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
