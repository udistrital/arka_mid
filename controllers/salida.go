package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/salidaHelper"
	"github.com/udistrital/arka_mid/models"
)

// SalidaController operations for Salida
type SalidaController struct {
	beego.Controller
}

// URLMapping ...
func (c *SalidaController) URLMapping() {
	c.Mapping("Post", c.Post)
}

// Post ...
// @Title Create
// @Description create Salida
// @Param	body		body 	models.Salida	true		"body for Salida content"
// @Success 201 {object} models.Salida
// @Failure 403 body is empty
// @router / [post]
func (c *SalidaController) Post() {
	var v models.TrSalida
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if respuesta := salidaHelper.AddSalida(&v); respuesta != nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = respuesta
		} else {
			logs.Error(respuesta)
			//c.Data["development"] = map[string]interface{}{"Code": "000", "Body": err.Error(), "Type": "error"}
			c.Data["system"] = err
			c.Abort("400")
		}
	} else {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "000", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = err
		c.Abort("400")
	}
	c.ServeJSON()
}
