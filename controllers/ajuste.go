package controllers

import (
	"encoding/json"
	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/arka_mid/helpers/ajustesHelper"
	"github.com/udistrital/utils_oas/errorctrl"
)

// AjusteController operations for Ajuste
type AjusteController struct {
	beego.Controller
}

// URLMapping ...
func (c *AjusteController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
}

// Post ...
// @Title Create
// @Description create Ajuste
// @Param	body		body 	ajustesHelper.PreTrAjuste	true		"body for Ajuste content"
// @Success 201 {object} models.Movimiento
// @Failure 403 body is empty
// @router / [post]
func (c *AjusteController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Ajuste by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Ajuste
// @Failure 403 :id is empty
// @router /:id [get]
func (c *AjusteController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Ajuste
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Ajuste
// @Failure 403
// @router / [get]
func (c *AjusteController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Ajuste
// @Param	id		path 	string	true		"The id you want to update"
// @Success 200 {object} models.Ajuste
// @Failure 403 :id is not int
// @router /:id [put]
func (c *AjusteController) Put() {

}
