package actaRecibido

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/autenticacion"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// Mapeo con los roles que pueden ver TODAS las actas en dicho estado
// Llaves: estados válidos de Actas
// Valores: mapeables desde models.RolesArka
var reglasVerTodas = map[string][]string{
	"Registrada":         []string{models.RolesArka["Admin"], models.RolesArka["Revisor"], models.RolesArka["Secretaria"]},
	"En Elaboracion":     []string{models.RolesArka["Admin"], models.RolesArka["Revisor"]},
	"En Modificacion":    []string{models.RolesArka["Admin"], models.RolesArka["Revisor"]},
	"En verificacion":    []string{models.RolesArka["Admin"], models.RolesArka["Revisor"]},
	"Aceptada":           []string{models.RolesArka["Admin"], models.RolesArka["Revisor"]},
	"Asociada a Entrada": []string{models.RolesArka["Admin"], models.RolesArka["Revisor"]},
	"Anulada":            []string{models.RolesArka["Admin"], models.RolesArka["Revisor"]},
}

func filtrarActasSegunRoles(actas []map[string]interface{}, usuarioWSO2 string,
	contratista int, proveedor int) (actasFiltradas []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/filtrarActasSegunRoles",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	if data, err := autenticacion.DataUsuario(usuarioWSO2); err == nil {
		// fmt.Printf("DATA_USUARIO: %+v\n", data)

		// Filtrar los tipos de actas que se pueden ver en
		// su totalidad (sin que tenga que estar asignada)
		ver := make(map[string]bool)
		for estado := range reglasVerTodas {
			ver[estado] = puedeVerActa(data.Role, estado)
		}
		// Si no se puede ver por completo TODAS las actas
		// de al menos un estado, se filtrarán las actas asignadas
		soloAsignadas := true
		for _, verTodasEnEsteEstado := range ver {
			if verTodasEnEsteEstado {
				soloAsignadas = false
				break
			}
		}
		// fmt.Printf("Solo Actas Asignadas: %t\n", soloAsignadas)
		// Consultar y Guardar los IDs de las actas asociadas al proveedor (si aplica)
		var idActasProveedor []int
		if soloAsignadas && proveedor > 0 {
			// Hacer consulta adicional para mirar cuales actas son de este proveedor
			url := "http://" + beego.AppConfig.String("actaRecibidoService") + "soporte_acta?limit=-1"
			url += "&query=Activo:true,ProveedorId:" + strconv.Itoa(proveedor)
			url += "&fields=ActaRecibidoId"

			var soportes []models.SoporteActa

			if _, err := request.GetJsonTest(url, &soportes); err == nil {
				// fmt.Printf("Proveedor: %d - #Soportes: %d\n", proveedor, len(soportes))
				for _, soporte := range soportes {
					// infoActa := soporte.ActaRecibidoId.Id
					// idActasProveedor[soporte.ActaRecibidoId.Id] = true
					idActasProveedor = append(idActasProveedor, soporte.ActaRecibidoId.Id)
				}
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "/filtrarActasSegunRoles",
					"err":     err,
					"status":  "502",
				}
				return nil, outputError
			}
			// fmt.Printf("Actas del proveedor: %v\n", idActasProveedor)
		}

		fin := len(actas)
		for i := 0; i < fin; {
			actaID := actas[i]["Id"]
			actaAsig := actas[i]["PersonaAsignadaId"]
			estado := actas[i]["Estado"]
			// fmt.Printf("actaId: %#v - asignada: %v\n", actaID, actaAsig)

			dejar := false
			if soloAsignadas {
				if contratista > 0 && actaAsig == contratista {
					// fmt.Printf("Acta %d - contratista %d\n", actaID, contratista)
					dejar = true
				} else if proveedor > 0 {
					// fmt.Printf("%#T\n", actaID)
					if idF, ok := actaID.(float64); ok {
						id := int(idF)
						for _, actaPrueba := range idActasProveedor {
							// fmt.Printf("%v-%#v\n", k, actaPrueba)
							if actaPrueba == id {
								// fmt.Printf("Acta %v - proveedor %v\n", actaID, proveedor)
								dejar = true
								break
							}
						}
					}
				}
			} else {
				if estadoActa, ok := estado.(string); ok {
					dejar = ver[estadoActa]
				}
			}

			if dejar {
				i++
			} else {
				actas[i] = actas[fin-1]
				fin--
			}
		}
		actas = actas[:fin]

		return actas, nil
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/filtrarActasSegunRoles",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
}

func filtrarActasPorEstados(actas []map[string]interface{}, states []string) (filtradas []map[string]interface{}) {
	fin := len(actas)
	for i := 0; i < fin; {
		dejar := false
		for _, reqState := range states {
			if actas[i]["Estado"] == reqState {
				dejar = true
				break
			}
		}
		if dejar {
			i++
		} else {
			actas[i] = actas[fin-1]
			fin--
		}
	}
	filtradas = actas[:fin]

	return filtradas
}

func puedeVerActa(rolesUsuario []string, estado string) bool {
	for _, rolSuficiente := range reglasVerTodas[estado] {
		for _, rolUsuario := range rolesUsuario {
			if rolUsuario == rolSuficiente {
				return true
			}
		}
	}
	return false
}
