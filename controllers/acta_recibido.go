package controllers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/actaRecibido"
)

// ActaRecibidoController operations for ActaRecibido
type ActaRecibidoController struct {
	beego.Controller
}

// URLMapping ...
func (c *ActaRecibidoController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// GetAllActasRecibido ...
// @Title GetAllActasRecibido
// @Description Devuelve las todas las actas de recibido
// @Success 200 {object} models.Acta_recibido
// @Failure 403
// @router /get_actas_recibido/ [get]
func (c *ActaRecibidoController) GetAllActasRecibido() {
	// idStr := c.Ctx.Input.Param(":id")
	// id, _ := strconv.Atoi(idStr)
	v, err := actaRecibido.GetAllActasRecibido()
	if err != nil {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("404")
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetActasRecibidoTipo ...
// @Title GetActasRecibidoTipo
// @Description Devuelve las todas las actas de recibido
// @Success 200 {object} models.Acta_recibido
// @Failure 403
// @router /get_actas_recibido_tipo/:id [get]
func (c *ActaRecibidoController) GetActasRecibidoTipo() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := actaRecibido.GetActasRecibidoTipo(id)
	if err != nil {
		logs.Error(err)
		c.Data["system"] = err
		c.Abort("404")
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// Post ...
// @Title Create
// @Description create Acta_recibido
// @Param	body		body 	models.Acta_recibido	true		"body for Acta_recibido content"
// @Success 201 {object} models.Acta_recibido
// @Failure 403 body is empty
// @router / [post]
func (c *ActaRecibidoController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Acta_recibido by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Acta_recibido
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ActaRecibidoController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Acta_recibido
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Acta_recibido
// @Failure 403
// @router / [get]
func (c *ActaRecibidoController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Acta_recibido
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Acta_recibido	true		"body for Acta_recibido content"
// @Success 200 {object} models.Acta_recibido
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ActaRecibidoController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Acta_recibido
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *ActaRecibidoController) Delete() {

}
