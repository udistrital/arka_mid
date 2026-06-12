package controllers

import (
	"errors"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/helpers/reportesHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
	"github.com/udistrital/arka_mid/utils_oas/request"
)

// ReportesController operations for reportes
type ReportesController struct {
	beego.Controller
}

// URLMapping ...
func (c *ReportesController) URLMapping() {
	c.Mapping("PostReporteElementos", c.PostReporteElementos)
	c.Mapping("PostPazYSalvo", c.PostPazYSalvo)
	c.Mapping("GetDetalleCuentasEntrada", c.GetDetalleCuentasEntrada)
	c.Mapping("GetDetalleCuentasSalida", c.GetDetalleCuentasSalida)
}

// PostReporteElementos ...
// @Title PostReporteElementos
// @Description Genera un reporte base64 en formato Excel para un rango de fechas.
// @Param	body	body	models.ReporteFechasRequest	true	"Rango de fechas del reporte"
// @Success 200 {object} models.ReporteExcelBase64Response
// @Failure 400 error en los datos de entrada
// @router /elementos [post]
func (c *ReportesController) PostReporteElementos() {
	defer errorCtrl.ErrorControlController(c.Controller, "ReportesController")

	var request models.ReporteFechasRequest
	if err := utilsHelper.Unmarshal(string(c.Ctx.Input.RequestBody), &request); err != nil {
		panic(errorCtrl.Error("PostReporteElementos - utilsHelper.Unmarshal(RequestBody, &request)", err, "400"))
	}

	respuesta, outputError := reportesHelper.GenerarReporteElementos(&request)
	if outputError != nil {
		panic(outputError)
	}

	c.Data["json"] = respuesta
	c.ServeJSON()
}

// PostPazYSalvo ...
// @Title PostPazYSalvo
// @Description Genera un paz y salvo de inventario en PDF base64 a partir del número de documento consultado.
// @Param	body	body	models.PazYSalvoRequest	true	"Datos de consulta y firma opcional"
// @Success 200 {object} models.PazYSalvoResponse
// @Failure 400 error en los datos de entrada
// @router /pazysalvo [post]
func (c *ReportesController) PostPazYSalvo() {
	defer errorCtrl.ErrorControlController(c.Controller, "ReportesController")

	var payload models.PazYSalvoRequest
	if err := utilsHelper.Unmarshal(string(c.Ctx.Input.RequestBody), &payload); err != nil {
		panic(errorCtrl.Error("PostPazYSalvo - utilsHelper.Unmarshal(RequestBody, &payload)", err, "400"))
	}

	headerAnterior := request.GetHeader()
	request.SetHeader(c.Ctx.Request.Header.Get("Authorization"))
	defer request.SetHeader(headerAnterior)

	respuesta, outputError := reportesHelper.GenerarPazYSalvo(&payload)
	if outputError != nil {
		panic(outputError)
	}

	c.Data["json"] = respuesta
	c.ServeJSON()
}

// GetDetalleCuentasEntrada ...
// @Title GetDetalleCuentasEntrada
// @Description Consulta el detalle contable por elemento de una entrada usando su consecutivo.
// @Param	EntradaConsecutivo	query	string	true	"Consecutivo de la entrada"
// @Success 200 {object} []models.ReporteDetalleEntradaResponse
// @Failure 400 error en los datos de entrada
// @router /detalle_cuentas_entrada [get]
func (c *ReportesController) GetDetalleCuentasEntrada() {
	defer errorCtrl.ErrorControlController(c.Controller, "ReportesController")

	consecutivo := c.GetString("EntradaConsecutivo")
	if consecutivo == "" {
		panic(errorCtrl.Error("GetDetalleCuentasEntrada - c.GetString(EntradaConsecutivo)", errors.New("se debe especificar EntradaConsecutivo"), "400"))
	}

	respuesta, outputError := reportesHelper.GetDetalleCuentasEntradaPorConsecutivo(consecutivo)
	if outputError != nil {
		panic(outputError)
	}

	c.Data["json"] = respuesta
	c.ServeJSON()
}

// GetDetalleCuentasSalida ...
// @Title GetDetalleCuentasSalida
// @Description Consulta el detalle contable por elemento de una salida usando su consecutivo.
// @Param	SalidaConsecutivo	query	string	true	"Consecutivo de la salida"
// @Success 200 {object} []models.ReporteDetalleSalidaResponse
// @Failure 400 error en los datos de entrada
// @router /detalle_cuentas_salida [get]
func (c *ReportesController) GetDetalleCuentasSalida() {
	defer errorCtrl.ErrorControlController(c.Controller, "ReportesController")

	consecutivo := c.GetString("SalidaConsecutivo")
	if consecutivo == "" {
		panic(errorCtrl.Error("GetDetalleCuentasSalida - c.GetString(SalidaConsecutivo)", errors.New("se debe especificar SalidaConsecutivo"), "400"))
	}

	respuesta, outputError := reportesHelper.GetDetalleCuentasSalidaPorConsecutivo(consecutivo)
	if outputError != nil {
		panic(outputError)
	}

	c.Data["json"] = respuesta
	c.ServeJSON()
}
