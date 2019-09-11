package controllers

import (
	"github.com/udistrital/acta_recibido_crud/models"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// ActaRecibidoController operations for ActaRecibido
type ActaRecibidoController struct {
	beego.Controller
}

// URLMapping ...
func (c *ActaRecibidoController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetAll", c.GetAll)
}

// Post ...
// @Title Create
// @Description create Acta_recibido
// @Param	body		body 	models.Acta_recibido	true		"body for Acta_recibido content"
// @Success 201 {object} models.Acta_recibido
// @Failure 403 body is empty
// @router / [post]
func (c *ActaRecibidoController) Post() {

	if multipartFile, _, err := c.GetFile("archivo"); err == nil {
		if Archivo, err := actaRecibido.DecodeXlsx2Json(multipartFile); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = Archivo
		} else {
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

// GetAll ...
// @Title Get All
// @Description get ActaRecibido
// @Success 200 {object} models.ActaRecibido
// @Failure 404 not found resource
// @router / [get]
func (c *ActaRecibidoController) GetAll() {

	l, err := models.GetAllParametrosActa()
	if err != nil {
		logs.Error(err)
		//c.Data["development"] = map[string]interface{}{"Code": "000", "Body": err.Error(), "Type": "error"}
		c.Data["system"] = err
		c.Abort("404")
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}
