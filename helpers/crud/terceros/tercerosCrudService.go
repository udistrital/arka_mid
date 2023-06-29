package terceros

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	"github.com/udistrital/arka_mid/utils_oas/request"
)

var basePath, _ = beego.AppConfig.String("tercerosService")

// GetCorreo Consulta el correo de un tercero
func GetCorreo(id int) (DetalleFuncionario []*models.InfoComplementariaTercero, outputError map[string]interface{}) {

	funcion := "GetCorreo"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		urlcrud string
		correo  []*models.InfoComplementariaTercero
	)

	// Consulta correo
	urlcrud = "http://" + basePath + "info_complementaria_tercero?limit=1&fields=Dato&sortby=Id&order=desc"
	urlcrud += "&query=Activo%3Atrue,InfoComplementariaId__Nombre__icontains%3Acorreo,TerceroId__Id%3A" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &correo); err != nil {
		logs.Error(err)
		eval := " - request.GetJson(urlcrud, &correo)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}

	return correo, nil
}

// GetAllDatosIdentificacion get controlador datos_identificacion de api terceros_crud
func GetAllDatosIdentificacion(query string) (datosId []models.DatosIdentificacion, outputError map[string]interface{}) {

	funcion := "GetAllDatosIdentificacion"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	// Consulta correo
	urlcrud := "http://" + basePath + "datos_identificacion?" + query
	if err := request.GetJson(urlcrud, &datosId); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetCorreo - request.GetJson(urlcrud, &response2)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	return datosId, nil
}

// DocumentosValidos recibe un arreglo de documentos activos
// y retorna true solo si son únicos en su tipo (si solo hay 1 CC, 1 NIT, 1...)
//
// si los documentos incluyen Activo=true/false, se pueden filtrar con soloActivos = true
//
// si los documentos incluyen el tercero, se puede validar que correspondan
// al mismo con validaTercero = true
func DocumentosValidos(documentos []models.DatosIdentificacion,
	soloActivos, validaTercero bool) bool {
	dicc := make(map[int]bool)
	tercero := -1
	for _, documento := range documentos {
		if soloActivos {
			if !documento.Activo {
				continue
			}
		}
		if _, ok := dicc[documento.TipoDocumentoId.Id]; ok {
			return false
		}
		dicc[documento.TipoDocumentoId.Id] = true
		if validaTercero { // Que estén asociados al mismo tercero
			terceroActual := documento.TerceroId.Id
			if tercero >= 0 {
				if terceroActual != tercero {
					return false
				}
			} else {
				tercero = terceroActual
			}
		}
	}
	return true
}

// GetTerceroById get controlador tercero/{id} del api terceros_crud
func GetTerceroById(id int) (tercero *models.Tercero, outputError map[string]interface{}) {

	funcion := "GetTerceroById"
	defer errorCtrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + basePath + "tercero/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &tercero); err != nil {
		eval := " - request.GetJson(urlcrud, &tercero)"
		return nil, errorCtrl.Error(funcion+eval, err, "502")
	}
	return tercero, nil
}

// GetDocUD Get documento de identificación UD
func GetDocUD() string {
	return "899999230"
}
