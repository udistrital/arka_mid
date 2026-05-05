package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/helpers/reportesHelper"
	"github.com/udistrital/arka_mid/helpers/utilsHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// ReportesController operations for reportes
type ReportesController struct {
	beego.Controller
}

// URLMapping ...
func (c *ReportesController) URLMapping() {
	c.Mapping("PostReporteElementos", c.PostReporteElementos)
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
