package terceros

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

// GetCorreo Consulta el correo de un tercero
func GetCorreo(id int) (DetalleFuncionario []*models.InfoComplementariaTercero, outputError map[string]interface{}) {

	funcion := "GetCorreo"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		urlcrud string
		correo  []*models.InfoComplementariaTercero
	)

	// Consulta correo
	urlcrud = "http://" + beego.AppConfig.String("tercerosService") + "info_complementaria_tercero?limit=1&fields=Dato&sortby=Id&order=desc"
	urlcrud += "&query=Activo%3Atrue,InfoComplementariaId__Nombre__icontains%3Acorreo,TerceroId__Id%3A" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &correo); err != nil {
		logs.Error(err)
		eval := " - request.GetJson(urlcrud, &correo)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}

	return correo, nil
}

// GetAllDatosIdentificacion get controlador datos_identificacion de api terceros_crud
func GetAllDatosIdentificacion(query string) (datosId []*models.DatosIdentificacion, outputError map[string]interface{}) {

	funcion := "GetAllDatosIdentificacion"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	// Consulta correo
	urlcrud := "http://" + beego.AppConfig.String("tercerosService") + "datos_identificacion?" + query
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
func DocumentosValidos(documentos []models.DatosIdentificacion) bool {
	dicc := make(map[int]bool)
	for _, documento := range documentos {
		if _, ok := dicc[documento.TipoDocumentoId.Id]; ok {
			return false
		}
		// if documento.Activo {
		dicc[documento.TipoDocumentoId.Id] = true
		// }
	}
	return true
}

// GetTerceroById get controlador tercero/{id} del api terceros_crud
func GetTerceroById(id int) (tercero *models.Tercero, outputError map[string]interface{}) {

	funcion := "GetTerceroById"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	urlcrud := "http://" + beego.AppConfig.String("tercerosService") + "tercero/" + strconv.Itoa(id)
	if err := request.GetJson(urlcrud, &tercero); err != nil {
		eval := " - request.GetJson(urlcrud, &tercero)"
		return nil, errorctrl.Error(funcion+eval, err, "502")
	}
	return tercero, nil
}

// GetDocUD Get documento de identificación UD
func GetDocUD() string {
	return "899999230"
}
