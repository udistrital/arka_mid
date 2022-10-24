package salidaHelper

import (
	"net/url"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

// PostTrSalidas Completa los detalles de las salidas y hace el respectivo registro en api movimientos_arka_crud
func PostTrSalidas(m *models.SalidaGeneral, etl bool) (resultado map[string]interface{}, outputError map[string]interface{}) {

	funcion := "PostTrSalidas"
	defer errorctrl.ErrorControlFunction(funcion+" - Unhandled Error!", "500")

	var (
		res                map[string][](map[string]interface{})
		estadoMovimientoId int
	)

	resultado = make(map[string]interface{})

	if err := movimientosArka.GetEstadoMovimientoIdByNombre(&estadoMovimientoId, "Salida En Trámite"); err != nil {
		return nil, err
	}

	for _, salida := range m.Salidas {

		if !etl {
			if err := setDetalleSalida("", 0, &salida.Salida.Detalle); err != nil {
				return nil, err
			}
		}

		salida.Salida.EstadoMovimientoId.Id = estadoMovimientoId
	}

	urlcrud := "http://" + beego.AppConfig.String("movimientosArkaService") + "tr_salida"

	// Crea registros en api movimientos_arka_crud
	if err := request.SendJson(urlcrud, "POST", &res, &m); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "PostTrSalidas - request.SendJson(movArka, \"POST\", &res, &m)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	resultado["trSalida"] = res

	return resultado, nil
}

func GetSalidas(tramiteOnly bool) (Salidas []map[string]interface{}, outputError map[string]interface{}) {

	funcion := "GetSalidas - "
	defer errorctrl.ErrorControlFunction(funcion+"Unhandled Error!", "500")

	asignaciones := make(map[int]models.AsignacionEspacioFisicoDependencia)
	sedes := make(map[string]models.EspacioFisico)
	funcionarios := make(map[int]models.Tercero)

	query := "limit=20&sortby=Id&order=desc&query=Activo:true,FormatoTipoMovimientoId__CodigoAbreviacion__in:SAL|SAL_CONS,EstadoMovimientoId__Nombre"
	if tramiteOnly {
		query += url.QueryEscape(":Salida En Trámite")
	} else {
		query += url.QueryEscape("__startswith:Salida")
	}

	if salidas_, err := movimientosArka.GetAllMovimiento(query); err != nil {
		return nil, err
	} else {
		if len(salidas_) == 0 {
			return nil, nil
		}

		for _, salida := range salidas_ {

			var formato models.FormatoSalida

			if err := utilsHelper.Unmarshal(salida.Detalle, &formato); err != nil {
				return nil, err
			}

			if salida__, err := TraerDetalle(salida, formato, asignaciones, sedes, funcionarios); err != nil {
				return nil, err
			} else {
				Salidas = append(Salidas, salida__)
			}

		}
	}

	return Salidas, nil
}

// GetInfoSalida Retorna el funcionario y el consecutivo de una salida a partir del detalle del movimiento
func GetInfoSalida(detalle string) (funcionarioId int, consecutivo string, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("GetInfoSalida - Unhandled Error!", "500")

	var detalle_ models.FormatoSalida
	if err := utilsHelper.Unmarshal(detalle, &detalle_); err != nil {
		return 0, "", err
	}

	return detalle_.Funcionario, detalle_.Consecutivo, nil
}

func getTipoComprobanteSalidas() string {
	return "H21"
}
