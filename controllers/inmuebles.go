package controllers

import (
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
	inmuebleshelper "github.com/udistrital/arka_mid/helpers/inmueblesHelper"
	"github.com/udistrital/arka_mid/models"
	"github.com/udistrital/arka_mid/utils_oas/errorCtrl"
)

// InmueblesController operations for Inmuebles
type InmueblesController struct {
	beego.Controller
}

// URLMapping ...
func (c *InmueblesController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
}

// Post ...
// @Title Create
// @Description create Inmuebles
// @Param	body	body	models.Inmuebles	true	"body for Inmuebles content"
// @Success 201 {object} models.Inmueble
// @Failure 403 body is empty
// @router / [post]
func (c *InmueblesController) Post() {

	defer errorCtrl.ErrorControlController(c.Controller, "InmueblesController")

	var v models.Inmueble
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err != nil {
		panic(err)
	}

	data, err_ := inmuebleshelper.Post(&v)
	if err_ == nil {
		c.Data["json"] = data
	} else {
		panic(err_)
	}

	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Inmuebles by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Inmueble
// @Failure 403 :id is empty
// @router /:id [get]
func (c *InmueblesController) GetOne() {

	defer errorCtrl.ErrorControlController(c.Controller, "InmueblesController")

	id, err := c.GetInt64(":id")
	if err != nil {
		panic(err)
	}

	data, err_ := inmuebleshelper.GetOne(int(id))
	if err_ == nil {
		c.Data["json"] = data
	} else {
		panic(err_)
	}

	c.ServeJSON()

}

// GetAll ...
// @Title GetAll
// @Description get Inmuebles
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Inmuebles
// @Failure 403
// @router / [get]
func (c *InmueblesController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Inmuebles
// @Param	id		path 	string			true	"The id you want to update"
// @Param	body	body 	models.Inmueble	true	"body for Inmuebles content"
// @Success 200 {object} models.Inmuebles
// @Failure 403 :id is not int
// @router /:id [put]
func (c *InmueblesController) Put() {

	defer errorCtrl.ErrorControlController(c.Controller, "InmueblesController")

	id, err := c.GetInt64(":id")
	if err != nil || id <= 0 {
		panic(err)
	}

	var v models.Inmueble
	err = json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err != nil {
		panic(err)
	}

	v.Elemento.Id = int(id)
	data, err_ := inmuebleshelper.Update(&v)
	if err_ == nil {
		c.Data["json"] = data
	} else {
		panic(err_)
	}

	c.ServeJSON()
}
