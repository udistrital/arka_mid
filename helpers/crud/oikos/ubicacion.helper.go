package oikos

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	e "github.com/udistrital/utils_oas/errorctrl"
)

func GetAsignacionSedeDependencia(Id int) (Relacion models.AsignacionEspacioFisicoDependencia, outputError map[string]interface{}) {
	const funcion = "GetAsignacionSedeDependencia - "
	defer e.ErrorControlFunction(funcion+"Unhandled Error!", fmt.Sprint(http.StatusInternalServerError))

	var res []*models.AsignacionEspacioFisicoDependencia
	if res, outputError = GetAllAsignacion(fmt.Sprintf("?query=Id:%d", Id)); outputError != nil {
		return
	} else {
		return *res[0], nil
	}
}

func GetSedeDependenciaUbicacion(ubicacionId int) (DetalleUbicacion *models.DetalleSedeDependencia, outputError map[string]interface{}) {
	const funcion = "GetSedeDependenciaUbicacion - "
	defer e.ErrorControlFunction(funcion+"Unhandled Error!", fmt.Sprint(http.StatusInternalServerError))

	var (
		urlcrud   string
		ubicacion []*models.AsignacionEspacioFisicoDependencia
	)
	resultado := new(models.DetalleSedeDependencia)

	urlcrud = "?query=Id:" + strconv.Itoa(ubicacionId)
	if ubicacion, outputError = GetAllAsignacion(urlcrud); outputError != nil {
		return
	}

	resultado.Dependencia = ubicacion[0].DependenciaId
	resultado.Ubicacion = ubicacion[0]

	if espFisico, err := utilsHelper.ConvertirInterfaceMap(ubicacion[0].EspacioFisicoId); err != nil {
		return nil, err
	} else {
		rgxp := regexp.MustCompile("[0-9]")
		strSede := espFisico["CodigoAbreviacion"].(string)
		strSede = rgxp.ReplaceAllString(strSede, "")
		urlcrud = "?query=CodigoAbreviacion:" + strSede
	}

	if sede_, err := GetAllEspacioFisico(urlcrud); err != nil {
		return nil, err
	} else {
		resultado.Sede = sede_[0]
	}

	return resultado, nil
}
