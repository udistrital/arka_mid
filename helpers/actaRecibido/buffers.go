package actaRecibido

import (
	"fmt"

	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/proveedorHelper"
	"github.com/udistrital/arka_mid/helpers/tercerosHelper"
	"github.com/udistrital/arka_mid/helpers/ubicacionHelper"
	"github.com/udistrital/arka_mid/models"
)

// findAndAddTercero trae la información de un tercero y la agrega
// al buffer de terceros
func findAndAddTercero(TerceroID int, Terceros map[int](map[string]interface{}),
	consultasTerceros *int, evTerceros *int) (ter map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "findAndAddTercero - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	if TerceroID == 0 {
		return nil, nil
	}

	idStr := fmt.Sprint(TerceroID)

	if Tercero, ok := Terceros[TerceroID]; ok {
		*evTerceros++
		return Tercero, nil
	}

	*consultasTerceros++
	if Tercero, err := tercerosHelper.GetNombreTerceroById(idStr); err == nil {
		if keys := len(Tercero); keys != 0 {
			Terceros[TerceroID] = Tercero
		}
		return Tercero, nil
	} else {
		return nil, err
	}
}

// findAndAddUbicacion trae la información de una ubicación y la agrega
// al buffer de ubicaciones
func findAndAddUbicacion(UbicacionID int, Ubicaciones map[int](map[string]interface{}),
	consultasUbicaciones *int, evUbicaciones *int) (ub map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "findAndAddUbicacion - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	if UbicacionID == 0 {
		return nil, nil
	}

	idStr := fmt.Sprint(UbicacionID)

	if ubicacion, ok := Ubicaciones[UbicacionID]; ok {
		*evUbicaciones++
		return ubicacion, nil
	}

	*consultasUbicaciones++
	if ubicacion, err := ubicacionHelper.GetAsignacionSedeDependencia(idStr); err == nil {
		if keys := len(ubicacion); keys != 0 {
			Ubicaciones[UbicacionID] = ubicacion
		}
		return ubicacion, nil
	} else {
		logs.Error(err)
		return nil, map[string]interface{}{
			"funcion": "findAndAddUbicacion",
			"err":     err,
			"status":  "502",
		}
	}
}

// findAndAddUbicacion trae la información de un proveedor y la agrega
// al buffer de proveedores
// (Nota: Evitar usar, se va a usar terceros en vez de Agora)
func findAndAddProveedor(ProveedorID int, Proveedores map[int](*models.Proveedor),
	consultasProveedores *int, evProveedores *int) (prov *models.Proveedor, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "findAndAddProveedor - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var vacio models.Proveedor
	if ProveedorID == 0 {
		return &vacio, nil
	}

	if proveedor, ok := Proveedores[ProveedorID]; ok {
		*evProveedores++
		return proveedor, nil
	}

	*consultasProveedores++
	if provs, err := proveedorHelper.GetProveedorById(ProveedorID); err == nil {
		if len(provs) == 1 && provs[0].Id > 0 {
			proveedor := provs[0]
			Proveedores[ProveedorID] = proveedor
			return proveedor, nil
		}
		if len(provs) > 1 {
			logs.Warn("Proveedor", ProveedorID, "tiene más de un resultado")
		}
	} else {
		return nil, err
	}

	return &vacio, nil
}
