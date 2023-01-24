package entradaHelper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/arka_mid/helpers/crud/cuentasContables"
	"github.com/udistrital/arka_mid/helpers/crud/movimientosArka"
	"github.com/udistrital/arka_mid/helpers/crud/parametros"
	"github.com/udistrital/arka_mid/helpers/mid/movimientosContables"
	"github.com/udistrital/arka_mid/models"
	e "github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

// AnularEntrada Anula una entrada y los movimientos posteriores a esta, el acta asociada queda en estado aceptada
func AnularEntrada(movimientoId int) (response map[string]interface{}, outputError map[string]interface{}) {

	const funcion = "AnularEntrada - "
	defer e.ErrorControlFunction(funcion+"unhandled error!", fmt.Sprint(http.StatusInternalServerError))

	var (
		urlcrud                 string
		err                     error
		res                     map[string]interface{}
		resMap                  map[string]interface{}
		movimientoArka          models.Movimiento
		transaccionActaRecibido models.TransaccionActaRecibido
		detalleMovimiento       map[string]interface{}
		estadoActa              models.EstadoActa
		estadoMovimiento        int
		parametroTipoDebito     int
		parametroTipoCredito    int
		tipoComprobante         string
		consecutivoId           int
		consecutivo             int
		transaccion             models.TransaccionMovimientos
		cuentasSubgrupo         []models.CuentasSubgrupo
		jsonString              []byte
	)

	res = make(map[string]interface{})

	query := "query=Id:" + strconv.Itoa(int(movimientoId))
	if mov, err := movimientosArka.GetAllMovimiento(query); err != nil {
		return nil, err
	} else {
		movimientoArka = *mov[0]
	}

	if err = json.Unmarshal([]byte(movimientoArka.Detalle), &detalleMovimiento); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "AnularEntrada - entrada.AnularEntrada(entrada);", "status": "500", "err": err}
		return nil, outputError
	}

	var resTrActa []models.TransaccionActaRecibido
	urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/" + fmt.Sprint(detalleMovimiento["acta_recibido_id"])
	if err = request.GetJson(urlcrud, &resTrActa); err != nil { // Get informacion acta de api acta_recibido_crud
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "AnularEntrada - actaRecibido.TransaccionActaRecibido(acta);", "status": "502", "err": err}
		return nil, outputError
	}

	transaccionActaRecibido = resTrActa[0]
	if err := movimientosArka.GetEstadoMovimientoIdByNombre(&estadoMovimiento, "Entrada Anulada"); err != nil {
		return nil, err
	}

	resMap = make(map[string]interface{})
	var resEstadoActa []models.EstadoActa

	urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "estado_acta?query=Nombre:Aceptada"
	if err = request.GetJson(urlcrud, &resEstadoActa); err != nil { // Get parametrización acta de api acta_recibido_crud
		logs.Error(err)
		return nil, e.Error(funcion+"request.GetJson(urlcrud, &resEstadoActa)", err, fmt.Sprint(http.StatusBadGateway))
	}
	estadoActa = resEstadoActa[0]
	movimientoArka.EstadoMovimientoId.Id = estadoMovimiento
	transaccionActaRecibido.UltimoEstado.EstadoActaId.Id = estadoActa.Id
	transaccionActaRecibido.UltimoEstado.Id = 0

	if dt, cr, err := parametros.GetParametrosDebitoCredito(); err != nil {
		return nil, err
	} else {
		parametroTipoDebito = dt
		parametroTipoCredito = cr
	}

	if err := cuentasContables.GetComprobante("E", &tipoComprobante); err != nil {
		return nil, err
	}

	year, _, _ := time.Now().Date()
	postConsecutivo := models.Consecutivo{
		Id:          0,
		ContextoId:  199,
		Year:        year,
		Consecutivo: 0,
		Descripcion: "Ajustes Arka",
		Activo:      true,
	}
	urlcrud = "http://" + beego.AppConfig.String("consecutivosService") + "consecutivo"
	if err = request.SendJson(urlcrud, "POST", &resMap, &postConsecutivo); err != nil {
		logs.Error(err)
		return nil, e.Error(funcion+`request.SendJson(urlcrud, "POST", &resMap, &postConsecutivo)`, err, fmt.Sprint(http.StatusBadGateway))
	}
	if consecutivoId, err = strconv.Atoi(fmt.Sprint(resMap["Data"].(map[string]interface{})["Id"])); err != nil {
		logs.Error(err)
		return nil, e.Error(funcion+`strconv.Atoi(fmt.Sprint(resMap["Data"].(map[string]interface{})["Id"]))`, err, fmt.Sprint(http.StatusBadGateway))
	}
	if consecutivo, err = strconv.Atoi(fmt.Sprint(resMap["Data"].(map[string]interface{})["Consecutivo"])); err != nil {
		logs.Error(err)
		return nil, e.Error(funcion+`strconv.Atoi(fmt.Sprint(resMap["Data"].(map[string]interface{})["Consecutivo"]))`,
			err, fmt.Sprint(http.StatusInternalServerError))
	}
	resMap = make(map[string]interface{})
	transaccion.ConsecutivoId = consecutivoId

	// Se crea map para agrupar los valores totales según el código del subgrupo
	mapSubgruposTotales := map[int]float64{}
	for _, elemento := range transaccionActaRecibido.Elementos { // Proceso para registrar el movimiento contable para cada elemento
		if mapSubgruposTotales[elemento.SubgrupoCatalogoId] == 0 {
			mapSubgruposTotales[elemento.SubgrupoCatalogoId] = elemento.ValorTotal
		} else {
			mapSubgruposTotales[elemento.SubgrupoCatalogoId] += elemento.ValorTotal
		}
	}

	etiquetas := make(map[string]interface{})
	etiquetas["TipoComprobanteId"] = tipoComprobante
	if jsonString, err = json.Marshal(etiquetas); err != nil {
		logs.Error(err)
		return nil, e.Error(funcion+"json.Marshal(etiquetas)", err, fmt.Sprint(http.StatusInternalServerError))
	}
	transaccion.Etiquetas = string(jsonString)
	transaccion.Activo = true
	transaccion.FechaTransaccion = time_bogota.Tiempo_bogota()
	transaccion.Descripcion = "Transacción para registrar movimientos contables correspondientes a entrada del sistema arka"

	for SubgrupoId, valor := range mapSubgruposTotales {
		var cuentaDebito models.CuentaContable
		var cuentaCredito models.CuentaContable
		var movimientoDebito models.MovimientoTransaccion
		var movimientoCredito models.MovimientoTransaccion

		urlcrud = "http://" + beego.AppConfig.String("catalogoElementosService") + "cuentas_subgrupo?query=SubgrupoId__Id:" + strconv.Itoa(SubgrupoId) + ",SubtipoMovimientoId:" + strconv.Itoa(movimientoArka.EstadoMovimientoId.Id) + ",Activo:true"
		if err = request.GetJson(urlcrud, &cuentasSubgrupo); err != nil { // Obtiene cuentas que deben ser afectadas
			logs.Error(err)
			return nil, e.Error(funcion+"request.GetJson(urlcrud, &cuentasSubgrupo)", err, fmt.Sprint(http.StatusBadGateway))
		}
		urlcrud = "http://" + beego.AppConfig.String("cuentasContablesService") + "nodo_cuenta_contable/" + cuentasSubgrupo[0].CuentaDebitoId
		if err = request.GetJson(urlcrud, &resMap); err != nil { // Se trae información de cuenta débito de api cuentas_contables_crud
			logs.Error(err)
			return nil, e.Error(funcion+"request.GetJson(urlcrud, &resMap)", err, fmt.Sprint(http.StatusBadGateway))
		}
		if jsonString, err = json.Marshal(resMap["Body"]); err != nil {
			logs.Error(err)
			return nil, e.Error(funcion+`json.Marshal(resMap["Body"])`, err, fmt.Sprint(http.StatusInternalServerError))
		}
		if err := json.Unmarshal(jsonString, &cuentaDebito); err != nil {
			logs.Error(err)
			return nil, e.Error(funcion+"json.Unmarshal(jsonString, &cuentaDebito)", err, fmt.Sprint(http.StatusInternalServerError))
		}
		resMap = make(map[string]interface{})

		movimientoDebito.NombreCuenta = cuentaDebito.Nombre
		movimientoDebito.CuentaId = cuentaDebito.Codigo
		movimientoDebito.TipoMovimientoId = parametroTipoCredito
		movimientoDebito.Valor = valor
		movimientoDebito.Descripcion = "Movimiento crédito registrado desde sistema arka"
		movimientoDebito.Activo = true
		movimientoDebito.TerceroId = nil // Provisional
		transaccion.Movimientos = append(transaccion.Movimientos, &movimientoDebito)

		urlcrud = "http://" + beego.AppConfig.String("cuentasContablesService") + "nodo_cuenta_contable/" + cuentasSubgrupo[0].CuentaCreditoId
		if err = request.GetJson(urlcrud, &resMap); err != nil { // Se trae información de cuenta crédito de api cuentas_contables_crud
			logs.Error(err)
			return nil, e.Error(funcion+"request.GetJson(urlcrud, &resMap)", err, fmt.Sprint(http.StatusBadGateway))
		}
		if jsonString, err = json.Marshal(resMap["Body"]); err != nil {
			logs.Error(err)
			return nil, e.Error(funcion+`json.Marshal(resMap["Body"])`, err, fmt.Sprint(http.StatusInternalServerError))
		}
		if err = json.Unmarshal(jsonString, &cuentaCredito); err != nil {
			logs.Error(err)
			return nil, e.Error(funcion+"json.Unmarshal(jsonString, &cuentaCredito)", err, fmt.Sprint(http.StatusInternalServerError))
		}
		movimientoCredito.NombreCuenta = cuentaCredito.Nombre
		movimientoCredito.CuentaId = cuentaCredito.Codigo
		movimientoCredito.TipoMovimientoId = parametroTipoDebito
		movimientoCredito.Valor = valor
		movimientoCredito.Descripcion = "Movimiento débito registrado desde sistema arka"
		movimientoCredito.Activo = true
		movimientoCredito.TerceroId = nil // Provisional
		transaccion.Movimientos = append(transaccion.Movimientos, &movimientoCredito)
	}

	urlcrud = "http://" + beego.AppConfig.String("movimientosContablesmidService") + "transaccion_movimientos/transaccion_movimientos"
	if _, err := movimientosContables.PostTrContable(&transaccion); err != nil {
		return nil, err
	}
	res["transaccion"] = transaccion

	// Anulación de salidas asociadas
	// Si el estado de movimientoArka es Entrada Asociada a una salida, continuar con la anulación de las salidas

	consecutivoAjuste := "H20-" + fmt.Sprintf("%05d", consecutivo) + "-" + strconv.Itoa(year)
	detalleMovimiento["consecutivo_ajuste"] = consecutivoAjuste
	detalleMovimiento["mov_contable_ajuste_consecutivo_id"] = transaccion.ConsecutivoId

	if jsonString, err = json.Marshal(detalleMovimiento); err != nil {
		logs.Error(err)
		return nil, e.Error(funcion+"json.Marshal(detalleMovimiento)", err, fmt.Sprint(http.StatusInternalServerError))
	}
	movimientoArka.Detalle = string(jsonString)
	if _, err := movimientosArka.PutMovimiento(&movimientoArka, movimientoArka.Id); err != nil {
		return nil, err
	}
	res["arka"] = movimientoArka.Detalle
	urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/" + fmt.Sprint(detalleMovimiento["acta_recibido_id"])
	if err = request.SendJson(urlcrud, "PUT", &transaccionActaRecibido, &transaccionActaRecibido); err != nil { // Puesto que se anula la entrada, el acta debe quedar disponible para volver ser asociada a una entrada
		logs.Error(err)
		return nil, e.Error(funcion+`request.SendJson(urlcrud, "PUT", &transaccionActaRecibido, &transaccionActaRecibido)`, err, fmt.Sprint(http.StatusBadGateway))
	}
	res["movArkaId"] = movimientoArka.EstadoMovimientoId.Id
	res["EstadoActaId"] = transaccionActaRecibido.UltimoEstado.EstadoActaId.Id
	return res, nil
}
