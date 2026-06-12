package customErrorv2

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/udistrital/arka_mid/utils_oas/auditoria"
	"github.com/udistrital/arka_mid/utils_oas/xray"
)

type CustomErrorController struct {
	web.Controller
}

func genericError(c *CustomErrorController, status string) {
	outputError := map[string]interface{}{"Success": false, "Status": status, "Message": c.Data["mesaage"], "Data": c.Data["data"]}
	xray.EndSegment(c.Ctx)
	auditoria.LogRequest(c.Ctx)
	c.Data["json"] = outputError
	c.ServeJSON()
}

func (c *CustomErrorController) Error400() {
	genericError(c, "400")
}

func (c *CustomErrorController) Error404() {
	genericError(c, "404")
}

func (c *CustomErrorController) Error500() {
	genericError(c, "500")
}

func (c *CustomErrorController) Error501() {
	genericError(c, "501")
}

func (c *CustomErrorController) Error502() {
	genericError(c, "502")
}

func (c *CustomErrorController) Error509() {
	genericError(c, "509")
}
