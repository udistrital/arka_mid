// Este archivo se borró erróneamente en el commit 74f0c97f7
// Para ver la historia antes de haberlo borrado, referirse
// al commit bfaf5a315, que fue el último antes de borrarlo

package actaRecibido

import (
	"github.com/udistrital/arka_mid/models"
)

// Mapeo con los roles que pueden ver TODAS las actas en dicho estado
// Llaves: estados válidos de Actas
// Valores: mapeables desde models.RolesArka
var reglasVerTodas = map[string][]string{
	"Registrada":         {models.RolesArka["Secretaria"]},
	"EnElaboracion":      {models.RolesArka["Proveedor"], models.RolesArka["Contratista"]},
	"EnModificacion":     {models.RolesArka["Contratista"]},
	"EnVerificacion":     {},
	"Aceptada":           {},
	"Asociada a Entrada": {},
	"Anulada":            {},
}

// Arreglo de roles que pueden ver actas en cualquier estado
var verCualquierEstado = []string{
	models.RolesArka["Admin"],
	models.RolesArka["Revisor"],
}
