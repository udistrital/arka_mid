package controllers

import (
	"encoding/json"
	"strconv"

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
	c.Mapping("Get", c.GetSalida)
	c.Mapping("GetAll", c.GetSalidas)
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

// GetSalida ...
// @Title Get User
// @Description get Salida by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Salida
// @Failure 404 not found resource
// @router /:id [get]
func (c *SalidaController) GetSalida() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := salidaHelper.GetSalida(id)
	if err != nil {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "000", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = err
		c.Abort("404")
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetSalidas ...
// @Title Get User
// @Description get Entradas
// @Success 200 {object} models.Salida
// @Failure 404 not found resource
// @router / [get]
func (c *SalidaController) GetSalidas() {
	v, err := salidaHelper.GetSalidas()
	if err != nil {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "000", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = err
		c.Abort("404")
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}