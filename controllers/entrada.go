package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/entradaHelper"
	"github.com/udistrital/arka_mid/models"
)

// EntradaController operations for Entrada
type EntradaController struct {
	beego.Controller
}

// URLMapping ...
func (c *EntradaController) URLMapping() {
	c.Mapping("Post", c.Post)
}

// Post ...
// @Title Post
// @Description create Entrada
// @Param	body		body 	models.Entrada	true		"body for Entrada content"
// @Success 201 {object} models.Entrada
// @Failure 403 body is empty
// @Failure 400 the request contains incorrect syntax
// @router / [post]
func (c *EntradaController) Post() {
	var v models.Movimiento
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if respuesta := entradaHelper.AddEntrada(v); respuesta != nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = respuesta
		} else {
			c.Data["system"] = respuesta
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

// GetEntrada ...
// @Title Get User
// @Description get Entrada by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ConsultaEntrada
// @Failure 404 not found resource
// @router /:id [get]
func (c *EntradaController) GetEntrada() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := entradaHelper.GetEntrada(id)
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

// GetEntradas ...
// @Title Get User
// @Description get Entradas
// @Success 200 {object} models.ConsultaEntrada
// @Failure 404 not found resource
// @router / [get]
func (c *EntradaController) GetEntradas() {
	v, err := entradaHelper.GetEntradas()
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
