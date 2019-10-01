package actaRecibidoHelper

import (
	"strconv"

	"github.com/udistrital/arka_mid/helpers/proveedorHelper"
	"github.com/udistrital/arka_mid/helpers/ubicacionHelper"

	"github.com/udistrital/arka_mid/helpers/parametrosGobiernoHelper"
	"github.com/udistrital/arka_mid/helpers/unidadHelper"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetAllActasRecibido ...
func GetAllActasRecibido() (historicoActa interface{}, outputError map[string]interface{}) {
	if response, err := request.GetJsonTest("http://"+beego.AppConfig.String("actaRecibidoService")+"historico_acta?query=ActaRecibidoId.Activo:True", &historicoActa); err == nil { // (2) error servicio caido
		if response.StatusCode == 200 { // (3) error estado de la solicitud
			return historicoActa, nil
		} else {
			logs.Info("Error (3) estado de la solicitud")
			outputError = map[string]interface{}{"Function": "GetAllActasRecibido:GetAllActasRecibido", "Error": response.Status}
			return outputError, nil
		}
	} else {
		logs.Info("Error (2) servicio caido")
		outputError = map[string]interface{}{"Function": "GetAllActasRecibido", "Error": err}
		return outputError, nil
	}
}

// GetActasRecibidoTipo ...
func GetActasRecibidoTipo(tipoActa int) (actasRecibido []models.ActaRecibidoUbicacion, outputError map[string]interface{}) {
	var (
		urlcrud string
		// elementos []models.Elemento
		historicoActa []*models.HistoricoActa
		//actaRecibidoAux *models.ActaRecibidoUbicacion
		//ubicacion []*models.EspacioFisico
		// auxE      models.ElementosActa
		// soporte   *models.SoporteActaProveedor
	)
	if tipoActa != 0 { // (1) error parametro
		urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "historico_acta?query=EstadoActaId.Id:" + strconv.Itoa(tipoActa) + ",Activo:True&limit=-1"
		logs.Debug(urlcrud)
		if response, err := request.GetJsonTest(urlcrud, &historicoActa); err == nil { // (2) error servicio caido
			logs.Debug(historicoActa[0].EstadoActaId)
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

// GetElementos ...
func GetElementos(actaId int) (elementosActa []models.ElementosActa, outputError map[string]interface{}) {
	var (
		urlcrud   string
		elementos []models.Elemento
		unidad    []*models.Unidad
		iva       []*models.ParametrosGobierno
		proveedor []*models.Proveedor
		auxE      models.ElementosActa
		soporte   *models.SoporteActaProveedor
	)
	if actaId != 0 { // (1) error parametro
		// Solicita información elementos acta
		urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "elemento?query=SoporteActaId.ActaRecibidoId.Id:" + strconv.Itoa(actaId) +
			",Activo:True&limit=-1"
		if response, err := request.GetJsonTest(urlcrud, &elementos); err == nil {
			// Solicita información unidad elemento
			urlcrud = "http://" + beego.AppConfig.String("administrativaService") + "/unidad/"
			for _, elemento := range elementos {
				if response.StatusCode == 200 { // (3) error estado de la solicitud
					auxE.Id = elemento.Id
					auxE.Nombre = elemento.Nombre
					auxE.Cantidad = elemento.Cantidad
					auxE.Marca = elemento.Marca
					auxE.Serie = elemento.Serie
					// UNIDAD DEMEDIDA
					unidad, outputError = unidadHelper.GetUnidad(elemento.UnidadMedida)
					auxE.UnidadMedida = unidad[0]

					auxE.ValorUnitario = elemento.ValorUnitario
					auxE.Subtotal = elemento.Subtotal
					auxE.Descuento = elemento.Descuento
					auxE.ValorTotal = elemento.ValorTotal
					// PORCENTAJE IVA
					iva, outputError = parametrosGobiernoHelper.GetIva(elemento.PorcentajeIvaId)
					auxE.PorcentajeIvaId = iva[0]

					auxE.ValorIva = elemento.ValorIva
					auxE.ValorFinal = elemento.ValorFinal
					auxE.SubgrupoCatalogoId = elemento.SubgrupoCatalogoId
					auxE.Verificado = elemento.Verificado
					auxE.TipoBienId = elemento.TipoBienId
					auxE.EstadoElementoId = elemento.EstadoElementoId
					// SOPORTE
					proveedor, outputError = proveedorHelper.GetProveedorById(elemento.SoporteActaId.ProveedorId)
					soporte = new(models.SoporteActaProveedor)
					soporte.Id = elemento.SoporteActaId.Id
					soporte.ActaRecibidoId = elemento.SoporteActaId.ActaRecibidoId
					soporte.Consecutivo = elemento.SoporteActaId.Consecutivo
					soporte.Activo = elemento.SoporteActaId.Activo
					soporte.FechaCreacion = elemento.SoporteActaId.FechaCreacion
					soporte.FechaModificacion = elemento.SoporteActaId.FechaModificacion
					soporte.FechaSoporte = elemento.SoporteActaId.FechaSoporte
					soporte.ProveedorId = proveedor[0]
					auxE.SoporteActaId = soporte

					auxE.Placa = elemento.Placa
					auxE.Activo = elemento.Activo
					auxE.FechaCreacion = elemento.FechaCreacion
					auxE.FechaModificacion = elemento.FechaModificacion

					elementosActa = append(elementosActa, auxE)
				} else {
					logs.Info("Error (3) estado de la solicitud")
					outputError = map[string]interface{}{"Function": "GetAllActasRecibido:GetAllActasRecibido", "Error": response.Status}
					return nil, outputError
				}
			}
			return elementosActa, nil
		} else {
			logs.Info("Error (2) servicio caido")
			outputError = map[string]interface{}{"Function": "GetIva", "Error": err}
			return nil, outputError
		}
	} else {
		logs.Info("Error (1) Parametro")
		outputError = map[string]interface{}{"Function": "FuncionalidadMidController:GetIva", "Error": "null parameter"}
		return nil, outputError
	}
}

// GetSoportes ...
func GetSoportes(actaId int) (soportesActa []models.SoporteActaProveedor, outputError map[string]interface{}) {
	var (
		urlcrud   string
		soportes  []models.SoporteActa
		proveedor []*models.Proveedor
		auxS      models.SoporteActaProveedor
	)
	if actaId != 0 { // (1) error parametro
		// Solicita información elementos acta
		urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "soporte_acta?query=ActaRecibidoId:" + strconv.Itoa(actaId) + ",ActaRecibidoId.Activo:True&limit=-1"
		if response, err := request.GetJsonTest(urlcrud, &soportes); err == nil {
			// Solicita información unidad elemento
			urlcrud = "http://" + beego.AppConfig.String("administrativaService") + "/unidad/"
			for _, soporte := range soportes {
				if response.StatusCode == 200 { // (3) error estado de la solicitud
					auxS.Id = soporte.Id
					auxS.Consecutivo = soporte.Consecutivo
					auxS.ActaRecibidoId = soporte.ActaRecibidoId
					auxS.FechaSoporte = soporte.FechaSoporte
					auxS.Activo = soporte.Activo
					// SOPORTE
					proveedor, outputError = proveedorHelper.GetProveedorById(soporte.ProveedorId)
					//soporteAux = new(models.SoporteActaProveedor)
					auxS.ProveedorId = proveedor[0]

					auxS.FechaCreacion = soporte.FechaCreacion
					auxS.FechaModificacion = soporte.FechaModificacion

					soportesActa = append(soportesActa, auxS)
				} else {
					logs.Info("Error (3) estado de la solicitud")
					outputError = map[string]interface{}{"Function": "GetAllActasRecibido:GetAllActasRecibido", "Error": response.Status}
					return nil, outputError
				}
			}
			return soportesActa, nil
		} else {
			logs.Info("Error (2) servicio caido")
			outputError = map[string]interface{}{"Function": "GetIva", "Error": err}
			return nil, outputError
		}
	} else {
		logs.Info("Error (1) Parametro")
		outputError = map[string]interface{}{"Function": "FuncionalidadMidController:GetIva", "Error": "null parameter"}
		return nil, outputError
	}
}
