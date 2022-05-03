package entradaHelper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

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
		movimientosKronos       models.MovimientoProcesoExterno
		detalleMovimiento       map[string]interface{}
		tipoMovimiento          models.TipoMovimiento
		estadoActa              models.EstadoActa
		estadoMovimiento        models.EstadoMovimiento
		parametroTipoDebito     models.Parametro
		parametroTipoCredito    models.Parametro
		tipoComprobanteContable models.TipoComprobante
		consecutivoId           int
		consecutivo             int
		transaccion             models.TransaccionMovimientos
		cuentasSubgrupo         []models.CuentaSubgrupo
		TipoEntradaKronos       models.TipoMovimiento
	)

	res = make(map[string]interface{})

	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento?query=Id:" + strconv.Itoa(int(movimientoId))
	var resMovArka []models.Movimiento
	if err = request.GetJson(urlcrud, &resMovArka); err != nil { // Get movimiento de api movimientos_arka_crud
		logs.Error(err)
		return nil, e.Error(funcion+"request.GetJson(urlcrud, &resMovArka)", err, fmt.Sprint(http.StatusBadGateway))
	}
	movimientoArka = resMovArka[0]
	if movimientoArka.Id <= 0 {
		logs.Error(err)
		return nil, e.Error(funcion+"movimientoArka.Id <= 0", err, fmt.Sprint(http.StatusBadGateway))
	}
	urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo?query=ProcesoExterno:" + strconv.Itoa(int(movimientoId))
	if err = request.GetJson(urlcrud, &resMap); err != nil { // Get movimiento de api movimientos_crud
		logs.Error(err)
		return nil, e.Error(funcion+"request.GetJson(urlcrud, &resMap)", err, fmt.Sprint(http.StatusBadGateway))
	}
	var resMovKronos []models.MovimientoProcesoExterno
	var jsonString []byte
	if jsonString, err = json.Marshal(resMap["Body"]); err != nil {
		logs.Error(err)
		return nil, e.Error(funcion+`json.Marshal(resMap["Body"])`, err, fmt.Sprint(http.StatusInternalServerError))
	}
	if err = json.Unmarshal(jsonString, &resMovKronos); err != nil {
		logs.Error(err)
		return nil, e.Error(funcion+"json.Unmarshal(jsonString, &resMovKronos)", err, fmt.Sprint(http.StatusInternalServerError))
	}

	resMap = make(map[string]interface{})
	movimientosKronos = resMovKronos[0]

	if err = json.Unmarshal([]byte(movimientoArka.Detalle), &detalleMovimiento); err != nil {
		logs.Error(err)
		return nil, e.Error(funcion+"json.Unmarshal([]byte(movimientoArka.Detalle), &detalleMovimiento)", err, fmt.Sprint(http.StatusInternalServerError))
	}

	var resTrActa []models.TransaccionActaRecibido

	urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/" + fmt.Sprint(detalleMovimiento["acta_recibido_id"])
	if err = request.GetJson(urlcrud, &resTrActa); err != nil { // Get informacion acta de api acta_recibido_crud
		logs.Error(err)
		return nil, e.Error(funcion+"request.GetJson(urlcrud, &resTrActa)", err, fmt.Sprint(http.StatusBadGateway))
	}
	transaccionActaRecibido = resTrActa[0]
	var resEstadoMovimiento []models.EstadoMovimiento

	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "estado_movimiento?query=Nombre:Entrada%20Anulada"
	if err = request.GetJson(urlcrud, &resEstadoMovimiento); err != nil { // Get parametrización estado_movimiento de api movimientos_arka_crud
		logs.Error(err)
		return nil, e.Error(funcion+"request.GetJson(urlcrud, &resEstadoMovimiento)", err, fmt.Sprint(http.StatusBadGateway))
	}
	estadoMovimiento = resEstadoMovimiento[0]

	urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "tipo_movimiento?query=Nombre:Entrada%20Anulada"
	if err = request.GetJson(urlcrud, &resMap); err != nil { // Get parametrización tipo_movimiento de api movimientos_crud
		logs.Error(err)
		return nil, e.Error(funcion+"request.GetJson(urlcrud, &resMap)", err, fmt.Sprint(http.StatusBadGateway))
	}
	var resTipoMovimiento []models.TipoMovimiento
	if jsonString, err = json.Marshal(resMap["Body"]); err != nil {
		logs.Error(err)
		return nil, e.Error(funcion+`json.Marshal(resMap["Body"])`, err, fmt.Sprint(http.StatusInternalServerError))
	}
	if err = json.Unmarshal(jsonString, &resTipoMovimiento); err != nil {
		logs.Error(err)
		return nil, e.Error(funcion+"json.Unmarshal(jsonString, &resTipoMovimiento)", err, fmt.Sprint(http.StatusInternalServerError))
	}
	resMap = make(map[string]interface{})
	tipoMovimiento = resTipoMovimiento[0]

	urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "tipo_movimiento?query=Nombre__iexact:" + strings.ReplaceAll(movimientoArka.FormatoTipoMovimientoId.Nombre, " ", "%20")
	if err = request.GetJson(urlcrud, &resMap); err != nil { // Get parametrización tipo_movimiento de api movimientos_crud
		logs.Error(err)
		return nil, e.Error(funcion+"request.GetJson(urlcrud, &resMap)", err, fmt.Sprint(http.StatusBadGateway))
	}
	if jsonString, err = json.Marshal(resMap["Body"]); err != nil {
		logs.Error(err)
		return nil, e.Error(funcion+`json.Marshal(resMap["Body"])`, err, fmt.Sprint(http.StatusInternalServerError))
	}
	if err = json.Unmarshal(jsonString, &resTipoMovimiento); err != nil {
		logs.Error(err)
		return nil, e.Error(funcion+"json.Unmarshal(jsonString, &resTipoMovimiento)", err, fmt.Sprint(http.StatusInternalServerError))
	}
	resMap = make(map[string]interface{})
	TipoEntradaKronos = resTipoMovimiento[0]
	var resEstadoActa []models.EstadoActa

	urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "estado_acta?query=Nombre:Aceptada"
	if err = request.GetJson(urlcrud, &resEstadoActa); err != nil { // Get parametrización acta de api acta_recibido_crud
		logs.Error(err)
		return nil, e.Error(funcion+"request.GetJson(urlcrud, &resEstadoActa)", err, fmt.Sprint(http.StatusBadGateway))
	}
	estadoActa = resEstadoActa[0]
	movimientoArka.EstadoMovimientoId.Id = estadoMovimiento.Id
	movimientosKronos.TipoMovimientoId.Id = tipoMovimiento.Id
	transaccionActaRecibido.UltimoEstado.EstadoActaId.Id = estadoActa.Id
	transaccionActaRecibido.UltimoEstado.Id = 0

	urlcrud = "http://" + beego.AppConfig.String("parametrosService") + "parametro?query=CodigoAbreviacion:MCC"
	if err = request.GetJson(urlcrud, &resMap); err != nil { // Get parámetro tipo movimiento contable crédito
		logs.Error(err)
		return nil, e.Error(funcion+"request.GetJson(urlcrud, &resMap)", err, fmt.Sprint(http.StatusBadGateway))
	}
	if jsonString, err = json.Marshal(resMap["Data"]); err != nil {
		logs.Error(err)
		return nil, e.Error(funcion+`json.Marshal(resMap["Data"])`, err, fmt.Sprint(http.StatusInternalServerError))
	}
	var parametro []models.Parametro
	if err = json.Unmarshal(jsonString, &parametro); err != nil {
		logs.Error(err)
		return nil, e.Error(funcion+"json.Unmarshal(jsonString, &parametro)", err, fmt.Sprint(http.StatusInternalServerError))
	}
	resMap = make(map[string]interface{})
	parametroTipoDebito = parametro[0]

	urlcrud = "http://" + beego.AppConfig.String("parametrosService") + "parametro?query=CodigoAbreviacion:MCD"
	if err = request.GetJson(urlcrud, &resMap); err != nil { // Get parámetro tipo movimiento contable débito
		logs.Error(err)
		return nil, e.Error(funcion+"request.GetJson(urlcrud, &resMap)", err, fmt.Sprint(http.StatusBadGateway))
	}
	if jsonString, err = json.Marshal(resMap["Data"]); err != nil {
		logs.Error(err)
		return nil, e.Error(funcion+`json.Marshal(resMap["Data"])`, err, fmt.Sprint(http.StatusInternalServerError))
	}
	if err = json.Unmarshal(jsonString, &parametro); err != nil {
		logs.Error(err)
		return nil, e.Error(funcion+"json.Unmarshal(jsonString, &parametro)", err, fmt.Sprint(http.StatusInternalServerError))
	}
	resMap = make(map[string]interface{})
	parametroTipoCredito = parametro[0]

	urlcrud = "http://" + beego.AppConfig.String("cuentasContablesService") + "tipo_comprobante"
	if err = request.GetJson(urlcrud, &resMap); err != nil { // Para obtener código del tipo de comprobante
		logs.Error(err)
		return nil, e.Error(funcion+"request.GetJson(urlcrud, &resMap)", err, fmt.Sprint(http.StatusBadGateway))
	}
	for _, sliceTipoComprobante := range resMap["Body"].([]interface{}) {
		if sliceTipoComprobante.(map[string]interface{})["TipoDocumento"] == "E" {
			if jsonString, err = json.Marshal(sliceTipoComprobante); err != nil {
				logs.Error(err)
				return nil, e.Error(funcion+"json.Marshal(sliceTipoComprobante)", err, fmt.Sprint(http.StatusInternalServerError))
			}
			if err = json.Unmarshal(jsonString, &tipoComprobanteContable); err != nil {
				logs.Error(err)
				return nil, e.Error(funcion+"json.Unmarshal(jsonString, &tipoComprobanteContable)", err, fmt.Sprint(http.StatusInternalServerError))
			}
			resMap = make(map[string]interface{})
		}
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
	etiquetas["TipoComprobanteId"] = tipoComprobanteContable.Codigo
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

		urlcrud = "http://" + beego.AppConfig.String("catalogoElementosService") + "cuentas_subgrupo?query=SubgrupoId__Id:" + strconv.Itoa(SubgrupoId) + ",SubtipoMovimientoId:" + strconv.Itoa(TipoEntradaKronos.Id) + ",Activo:true"
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
		movimientoDebito.TipoMovimientoId = parametroTipoCredito.Id
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
		movimientoCredito.TipoMovimientoId = parametroTipoDebito.Id
		movimientoCredito.Valor = valor
		movimientoCredito.Descripcion = "Movimiento débito registrado desde sistema arka"
		movimientoCredito.Activo = true
		movimientoCredito.TerceroId = nil // Provisional
		transaccion.Movimientos = append(transaccion.Movimientos, &movimientoCredito)
	}

	res["transaccion"] = transaccion
	var resMovmientoContable interface{}

	urlcrud = "http://" + beego.AppConfig.String("movimientosContablesmidService") + "transaccion_movimientos/transaccion_movimientos"
	if err = request.SendJson(urlcrud, "POST", &resMovmientoContable, &transaccion); err != nil {
		logs.Error(err)
		return nil, e.Error(funcion+`request.SendJson(urlcrud, "POST", &resMovmientoContable, &transaccion)`, err, fmt.Sprint(http.StatusBadGateway))
	}
	if resMovmientoContable.(map[string]interface{})["Status"] != "201" {
		res["Respuesta movimientos contables Entradas"] = resMovmientoContable.(map[string]interface{})["Data"]
		return res, e.Error(funcion+`resMovmientoContable.(map[string]interface{})["Status"] != "201"`, err, fmt.Sprint(http.StatusBadGateway))
	}
	res["Respuesta movimientos contables Entradas"] = resMovmientoContable

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
	urlcrud = "http://" + beego.AppConfig.String("movimientosArkaService") + "movimiento/" + strconv.Itoa(int(movimientoId))
	if err = request.SendJson(urlcrud, "PUT", &movimientoArka, &movimientoArka); err != nil { // Put en el api movimientos_arka_crud
		logs.Error(err)
		return nil, e.Error(funcion+`request.SendJson(urlcrud, "PUT", &movimientoArka, &movimientoArka)`, err, fmt.Sprint(http.StatusBadGateway))
	}
	res["arka"] = movimientoArka.Detalle
	urlcrud = "http://" + beego.AppConfig.String("movimientosKronosService") + "movimiento_proceso_externo/" + strconv.Itoa(movimientoArka.Id)
	if err = request.SendJson(urlcrud, "PUT", &movimientosKronos, &movimientosKronos); err != nil { // Put api movimientos_crud
		logs.Error(err)
		return nil, e.Error(funcion+`request.SendJson(urlcrud, "PUT", &movimientosKronos, &movimientosKronos)`, err, fmt.Sprint(http.StatusBadGateway))
	}
	urlcrud = "http://" + beego.AppConfig.String("actaRecibidoService") + "transaccion_acta_recibido/" + fmt.Sprint(detalleMovimiento["acta_recibido_id"])
	if err = request.SendJson(urlcrud, "PUT", &transaccionActaRecibido, &transaccionActaRecibido); err != nil { // Puesto que se anula la entrada, el acta debe quedar disponible para volver ser asociada a una entrada
		logs.Error(err)
		return nil, e.Error(funcion+`request.SendJson(urlcrud, "PUT", &transaccionActaRecibido, &transaccionActaRecibido)`, err, fmt.Sprint(http.StatusBadGateway))
	}
	res["movArkaId"] = movimientoArka.EstadoMovimientoId.Id
	res["movKronosId"] = movimientosKronos.TipoMovimientoId.Id
	res["EstadoActaId"] = transaccionActaRecibido.UltimoEstado.EstadoActaId.Id
	return res, nil
}
