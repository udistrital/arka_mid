package tercerosHelper

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego"

	"github.com/udistrital/utils_oas/request"
)

//GetNombreTerceroById trae el nombre de un encargado por su id
func GetNombreTerceroById(idTercero string) (tercero map[string]interface{}, err error) {
	var urltercero string
	var personas []map[string]interface{}

	urltercero = "http://" + beego.AppConfig.String("tercerosService") + "datos_identificacion/?query=TerceroId__Id:" + idTercero + ",Activo:true"
	if response, err := request.GetJsonTest(urltercero, &personas); err == nil {
		if response.StatusCode == 200 {
			for _, element := range personas {
				if len(element) == 0 {
					return nil, errors.New("No se encontro registro")
				} else {
					fmt.Println("encargado: ", element)
					return map[string]interface{}{
						"Numero":         element["Numero"],
						"NombreCompleto": element["TerceroId"].(map[string]interface{})["NombreCompleto"],
					}, nil
				}

			}
		} else if response.StatusCode == 400 {
			return nil, err
		}
	} else {
		fmt.Println("error: ", err)
		return nil, err
	}

	return

}
