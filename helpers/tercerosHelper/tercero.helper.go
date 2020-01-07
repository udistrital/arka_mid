package tercerosHelper

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

//GetNombreTerceroById trae el nombre de un encargado por su id
func GetNombreTerceroById(idTercero map[string]interface{}) (tercero map[string]interface{}, err error) {
	var urltercero string
	var personas []map[string]interface{}

	urltercero = beego.AppConfig.String("terceros") + "tercero/?query=Id:5,Activo:true&fields=NombreCompleto"
	if response, err := request.GetJsonTest(urltercero, &personas); err == nil {
		if response.StatusCode == 200 {
			for _, element := range personas {
				if len(element) == 0 {
					return nil, errors.New("No se encontro registro")
				} else {
					fmt.Println("encargado: ", element)
					return element, nil
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
