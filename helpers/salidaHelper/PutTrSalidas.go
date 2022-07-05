package salidaHelper

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/arka_mid/helpers/crud/consecutivos"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/utils_oas/errorctrl"
)

func PutTrSalidas(m *models.SalidaGeneral, salidaId int) (resultado map[string]interface{}, outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("PutTrSalidas - Unhandled Error!", "500")

	var (
		detalleOriginal    models.FormatoSalida
		estadoMovimientoId int
		query              string
	)

	resultado = make(map[string]interface{})

	// El objetivo es generar los respectivos consecutivos en caso de generarse más de una salida a partir de la original

	// Se consulta la salida original
	query = "limit=1&query=Id:" + strconv.Itoa(salidaId)
	if salida_, err := movimientosArka.GetAllMovimiento(query); err != nil {
		return nil, err
	} else if len(salida_) == 1 && salida_[0].EstadoMovimientoId.Nombre == "Salida Rechazada" {
		if err := utilsHelper.Unmarshal(salida_[0].Detalle, &detalleOriginal); err != nil {
			return nil, err
		}
	} else {
		return
	}

	if err := movimientosArka.GetEstadoMovimientoIdByNombre(&estadoMovimientoId, "Salida En Trámite"); err != nil {
		return nil, err
	}

	if len(m.Salidas) == 1 {
		// Si no se generan nuevas salidas, tan solo se debe actualizar el funcionario y ubicación de la salida original así como la vida útil y valor residual de los elementos

		if err := setDetalleSalida(detalleOriginal.Consecutivo, detalleOriginal.ConsecutivoId, &m.Salidas[0].Salida.Detalle); err != nil {
			return nil, err
		}

		m.Salidas[0].Salida.EstadoMovimientoId.Id = estadoMovimientoId
		m.Salidas[0].Salida.Id = salidaId
		if trRes, err := movimientosArka.PutTrSalida(m); err != nil {
			return nil, err
		} else {
			resultado["trSalida"] = trRes
		}

	} else {
		// Si se generaron salidas a partir de la original, se debe asignar un consecutivo a cada una y una de ellas debe tener el original

		// Se debe decidir a cuál de las nuevas asignarle el id y el consecutivo original
		index := -1
		var detalleNuevo models.FormatoSalida
		for idx, l := range m.Salidas {
			if err := utilsHelper.Unmarshal(l.Salida.Detalle, &detalleNuevo); err != nil {
				return nil, err
			}

			if detalleNuevo.Funcionario == detalleOriginal.Funcionario && detalleNuevo.Ubicacion == detalleOriginal.Ubicacion {
				index = idx
				break
			} else if detalleNuevo.Funcionario == detalleOriginal.Funcionario {
				index = idx
				break
			} else if detalleNuevo.Ubicacion == detalleOriginal.Ubicacion {
				index = idx
				break
			}
		}

		if index == -1 {
			index = 0
		}

		for idx, salida := range m.Salidas {
			var (
				id            int
				consecutivo   string
				consecutivoId int
			)

			if idx == index {
				id = salidaId
				consecutivoId = detalleOriginal.ConsecutivoId
				consecutivo = detalleOriginal.Consecutivo
			}

			if err := setDetalleSalida(consecutivo, consecutivoId, &salida.Salida.Detalle); err != nil {
				return nil, err
			}

			salida.Salida.Id = id
			salida.Salida.EstadoMovimientoId.Id = estadoMovimientoId
		}

		// Hace el put api movimientos_arka_crud
		if trRes, err := movimientosArka.PutTrSalida(m); err != nil {
			return nil, err
		} else {
			resultado["trSalida"] = trRes
		}
	}

	return resultado, nil
}

func setDetalleSalida(consecutivo string, consecutivoId int, detalle *string) (outputError map[string]interface{}) {

	defer errorctrl.ErrorControlFunction("setDetalleSalida - Unhandled Error!", "500")

	var detalle_ models.FormatoSalida
	if err := utilsHelper.Unmarshal(*detalle, &detalle_); err != nil {
		return err
	}

	if consecutivo == "" || consecutivoId <= 0 {
		var consecutivo_ models.Consecutivo
		ctxSalida, _ := beego.AppConfig.Int("contxtSalidaCons")
		if err := consecutivos.Get(ctxSalida, "Registro Salida Arka", &consecutivo_); err != nil {
			return err
		}
		consecutivo = consecutivos.Format("%05d", getTipoComprobanteSalidas(), &consecutivo_)
		consecutivoId = consecutivo_.Id
	}

	detalle_.Consecutivo = consecutivo
	detalle_.ConsecutivoId = consecutivoId
	if err := utilsHelper.Marshal(detalle_, detalle); err != nil {
		return err
	}

	return
}
